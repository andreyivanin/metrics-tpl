package storage

import (
	"context"
	"database/sql"
	"log"
	"metrics-tpl/internal/server/config"
	"metrics-tpl/internal/server/models"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type MetricSQL struct {
	id    string
	mtype string
	value sql.NullFloat64
	delta sql.NullInt64
}

type memSQL struct {
	db     *sql.DB
	config config.Config
}

func newSQL(ctx context.Context, cfg config.Config) (*memSQL, error) {
	// ps := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
	// `localhost`, `video`, `XXXXXXXX`, `video`)

	ps := cfg.DatabaseDSN

	db, err := sql.Open("pgx", ps)
	if err != nil {
		panic(err)
	}

	// defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	createTableQuery := `
		CREATE TABLE IF NOT EXISTS metrics (
			id VARCHAR(50),
			type VARCHAR(50),
			value DOUBLE PRECISION,
			delta BIGINT
		);`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, err
	}

	memSQL := &memSQL{
		db:     db,
		config: cfg,
	}

	return memSQL, nil
}

func (ms *memSQL) UpdateMetric(ctx context.Context, name, mtype string, m models.Metric) (models.Metric, error) {
	var err error

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	tx, err := ms.db.Begin()
	if err != nil {
		return nil, err
	}

	var mUpdated models.Metric

	mCurrent, err := ms.GetMetric(ctx, name)

	switch mtype {
	case _GAUGE:
		if err != nil {
			_, err = tx.ExecContext(ctx,
				"INSERT INTO metrics (id, type, value)"+
					" VALUES($1,$2,$3)", name, mtype, m)
			return m, tx.Commit()
		}

		mUpdated = m.(models.Gauge)

		_, err = tx.ExecContext(ctx,
			"UPDATE metrics SET value = $2"+
				" WHERE id = $1", name, mUpdated)

	case _COUNTER:
		if err != nil {
			_, err = tx.ExecContext(ctx,
				"INSERT INTO metrics (id, type, delta)"+
					" VALUES($1,$2,$3)", name, mtype, m)
			return m, tx.Commit()
		}

		mUpdated = mCurrent.(models.Counter) + m.(models.Counter)

		_, err = tx.ExecContext(ctx,
			"UPDATE metrics SET delta = $2"+
				" WHERE id = $1", name, mUpdated)

	default:
		log.Println("wrong metric type")
	}

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return mUpdated, tx.Commit()

}

func (ms *memSQL) GetMetric(ctx context.Context, name string) (models.Metric, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	row := ms.db.QueryRowContext(ctx,
		"SELECT * FROM metrics WHERE id = $1", name)

	m := MetricSQL{}

	err := row.Scan(&m.id, &m.mtype, &m.value, &m.delta)
	if err != nil {
		return "", err
	}

	var metric interface{}

	switch m.mtype {
	case _GAUGE:
		metric = models.Gauge(m.value.Float64)

	case _COUNTER:
		metric = models.Counter(m.delta.Int64)

	default:
		log.Println("wrong metric type")
	}

	return metric, nil

}

func (ms *memSQL) GetAllMetrics(ctx context.Context) (models.Metrics, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	metrics := make(models.Metrics)

	rows, err := ms.db.QueryContext(ctx, "SELECT * from metrics")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var m MetricSQL
		err := rows.Scan(&m.id, &m.mtype, &m.value, &m.delta)
		if err != nil {
			return nil, err
		}

		var metric models.Metric

		switch m.mtype {
		case _GAUGE:
			metric = models.Gauge(m.value.Float64)

		case _COUNTER:
			metric = models.Counter(m.delta.Int64)

		default:
			log.Println("wrong metric type")
		}

		metrics[m.id] = metric
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return metrics, nil
}

func (ms *memSQL) GetConfig() config.Config {
	return ms.config
}
