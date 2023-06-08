package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRead(t *testing.T) {

	const (
		SERVERADDRPORT = "localhost:8080"
		STOREINTERVAL  = 100
		STOREFILE      = "metrics-db-default.json"
		RESTORE        = true
	)

	tests := []struct {
		name    string
		flags   []string
		envs    map[string]string
		want    Config
		wantErr bool
	}{
		{
			name: "varOrderCheck",
			flags: []string{
				"cmd",
				"-a", "localhost:8085",
				"-i", "30",
				"-f", "metrics-db-flag.json",
				"-r", "true",
			},
			envs: map[string]string{
				"ADDRESS":        "localhost:8086",
				"STORE_INTERVAL": "20",
				"STORE_FILE":     "metrics-db-env.json",
				"RESTORE":        "true",
			},
			want: Config{
				Address:          "localhost:8086",
				StoreInterval:    time.Duration(20) * time.Second,
				StoreFile:        "metrics-db-env.json",
				RestoreSavedData: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.flags

			for env, value := range tt.envs {
				os.Setenv(env, value)

			}

			defer func() {
				for env := range tt.envs {
					os.Unsetenv(env)

				}
			}()

			got, _ := Read()
			assert.Equal(t, tt.want, got)

		})
	}
}
