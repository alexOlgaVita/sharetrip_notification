package apiIntegrationTest_test

import (
	"context"
	"database/sql"
	"github.com/jackc/pgx/v5/pgxpool"
	"job4j.ru/sharetrip-notification/internal/api"
	"job4j.ru/sharetrip-notification/internal/domain"
	_interface "job4j.ru/sharetrip-notification/internal/interface"
	"job4j.ru/sharetrip-notification/internal/repository"
	"job4j.ru/sharetrip-notification/internal/service"
	"log"
	"os"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var (
	testCtx       context.Context
	testDB        *sql.DB
	testPool      _interface.DBQuerier
	testApp       *fiber.App
	testContainer *postgres.PostgresContainer
)

func TestMain(m *testing.M) {
	testCtx = context.Background()

	var err error

	testContainer, err = postgres.Run(
		testCtx,
		"postgres:16",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("password"),
	)
	if err != nil {
		log.Fatalf("start postgres container: %v", err)
	}

	dsn, err := testContainer.ConnectionString(
		testCtx,
		"sslmode=disable",
	)
	if err != nil {
		log.Fatalf("get connection string: %v", err)
	}

	testDB, err = sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("open sql db: %v", err)
	}

	waitReady(testDB)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("unable to parse config: %v", err)
	}

	pgxPool, err := pgxpool.NewWithConfig(testCtx, config)
	if err != nil {
		log.Fatalf("unable to create pgxpool: %v", err)
	}

	// 2. Теперь создаем наш адаптер, который реализует _interface.DBQuerier
	testPool = repository.NewPoolAdapter(pgxPool)

	if err = goose.SetDialect("postgres"); err != nil {
		log.Fatalf("set goose dialect: %v", err)
	}

	if err = goose.Up(testDB, "../../migrations"); err != nil {
		log.Fatalf("run migrations: %v", err)
	}

	if err != nil {
		log.Fatalf("create pgx pool: %v", err)
	}

	repo := repository.NewRepoPg(testPool)
	server := api.NewServer(repo)
	server.NotificationService = &service.NotificationService{
		Pool: testPool,
		NotificationUsecase: &domain.NotificationUsecase{
			NotificationRepo: repo,
		},
	}

	testApp = fiber.New()
	server.Route(testApp.Group(""))

	code := m.Run()

	if testPool != nil {
		defer func() { _ = testPool.Close() }()
	}
	if testDB != nil {
		_ = testDB.Close()
	}
	if testContainer != nil {
		_ = testContainer.Terminate(testCtx)
	}

	os.Exit(code)
}

func waitReady(db *sql.DB) {
	deadline := time.Now().Add(30 * time.Second)

	for time.Now().Before(deadline) {
		ctx, cancel := context.WithTimeout(
			context.Background(),
			2*time.Second,
		)
		err := db.PingContext(ctx)
		cancel()

		if err == nil {
			return
		}

		time.Sleep(500 * time.Millisecond)
	}

	log.Fatalf("database is not ready after timeout")
}
