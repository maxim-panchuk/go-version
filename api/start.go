package api

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/maxim-panchuk/go-version/db"
	"github.com/maxim-panchuk/go-version/kafka"
	"github.com/maxim-panchuk/go-version/service"
)

func Start(pool *pgxpool.Pool) {
	repo := db.GetRepo(pool)
	metadataService := service.GetVersionMetadateService()
	settingService := service.NewVersionSettingServiceImpl(repo, metadataService)
	kafka.ListenForServiceInfo(metadataService, settingService)
}

// func initRepo() *pgxpool.Pool {
// 	config, err := pgxpool.ParseConfig("postgres://postgres:kilibok47Hromos@localhost:5432/diasoft_versioning")
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
// 		os.Exit(1)
// 	}

// 	config.MaxConns = 10
// 	pool, err := pgxpool.ConnectConfig(context.Background(), config)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
// 		os.Exit(1)
// 	}

// 	return pool

// }
