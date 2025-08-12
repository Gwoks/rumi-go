package infrastructure

import (
	"context"

	"github.com/smira/go-statsd"

	"rumi-go/internal/config"
	database "rumi-go/internal/infrastructure/database"
)

type Infra struct {
	StatsDClient *statsd.Client
	Database     database.Store
}

func NewInfra(_ context.Context, cfg *config.Config) (*Infra, error) {
	infra := &Infra{}

	// init database client using main database config
	databaseStore, err := database.NewStore(context.Background(), cfg.Database)
	if err != nil {
		return nil, err
	}
	infra.Database = databaseStore

	return infra, nil
}
