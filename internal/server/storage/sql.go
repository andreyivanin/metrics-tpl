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

	defer db.Close()

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

func (mS *memSQL) UpdateMetric(ctx context.Context, name, mtype string, m models.Metric) (models.Metric, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	return m, nil
}

func (mS *memSQL) GetMetric(ctx context.Context, mname string) (models.Metric, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	row := mS.db.QueryRowContext(ctx,
		"SELECT * FROM metrics WHERE id = $1", mname)

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

func (mS *memSQL) GetAllMetrics(ctx context.Context) (models.Metrics, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return nil, nil
}

func (mS *memSQL) GetConfig() config.Config {
	return mS.config
}
