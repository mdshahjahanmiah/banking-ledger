package main

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/mdshahjahanmiah/banking-ledger/pkg/account"
	"github.com/mdshahjahanmiah/banking-ledger/pkg/config"
	"github.com/mdshahjahanmiah/banking-ledger/pkg/db"
	"github.com/mdshahjahanmiah/explore-go/di"
	eHttp "github.com/mdshahjahanmiah/explore-go/http"
	"github.com/mdshahjahanmiah/explore-go/logging"
	"go.uber.org/dig"
	"log/slog"
)

func main() {
	slog.Info("transaction ledger service is starting...")
	c := di.New()

	c.Provide(func() (config.Config, error) {
		conf, err := config.Load()
		if err != nil {
			slog.Error("failed to load configuration", "err", err)
			return config.Config{}, err
		}
		return conf, nil
	})

	slog.Info("configuration is loaded successfully")

	c.Provide(func(conf config.Config) (*logging.Logger, error) {
		logger, err := logging.NewLogger(conf.LoggerConfig)
		if err != nil {
			slog.Error("initializing logger", "err", err)
			return nil, err
		}

		return logger, nil
	})

	slog.Info("logger is initialized successfully")

	// PostgreSQL connection + migrations (modified section)
	c.Provide(func(conf config.Config, logger *logging.Logger) (*db.DB, error) {
		slog.Info("PostgresDSN", "dsn", conf.PostgresDSN)

		// Initialize DB connection
		database, err := db.NewDB(conf.PostgresDSN, logger)
		if err != nil {
			logger.Error("database initialization", "err", err.Error())
			return nil, err
		}

		// Run migrations with dirty state handling
		m, err := migrate.New(
			"file://migrations",
			conf.PostgresDSN,
		)
		if err != nil {
			logger.Error("migration init failed", "err", err.Error())
			return nil, err
		}

		// Check current version and handle dirty state
		version, dirty, verErr := m.Version()
		if verErr != nil && verErr != migrate.ErrNilVersion {
			logger.Error("failed to check migration version", "err", verErr.Error())
			return nil, verErr
		}

		if dirty {
			logger.Warn("database is dirty, forcing version", "version", version)
			if forceErr := m.Force(int(version)); forceErr != nil {
				logger.Error("failed to force migration version",
					"version", version,
					"error", forceErr.Error(),
				)
				return nil, forceErr
			}
		}

		// Apply migrations
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			logger.Error("migration run failed", "err", err.Error())
			return nil, err
		}

		logger.Info("database migration completed successfully")
		return database, nil
	})

	c.Provide(func(config config.Config) *eHttp.ServerConfig {
		return &eHttp.ServerConfig{
			HttpAddress: config.HttpAddress,
		}
	})

	c.Provide(func(config config.Config, logger *logging.Logger, db *db.DB) (account.Service, error) {
		service := account.NewService(config, logger, db)
		return service, nil
	})

	c.ProvideMonitoringEndpoints("endpoint")
	c.Provide(account.MakeHandler, dig.Group("endpoint"))

	c.Invoke(func(in struct {
		dig.In
		Conf         config.Config
		ServerConfig *eHttp.ServerConfig
		Endpoints    []eHttp.Endpoint `group:"endpoint"`
	}) {
		server := eHttp.NewServer(in.ServerConfig, in.Endpoints, nil)
		c.Provide(func() di.StartCloser { return server }, dig.Group("startclose"))
	})

	err := c.Start()
	if err != nil {
		slog.Error("failed to start server", "err", err)
		return
	}
}
