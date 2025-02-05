package db

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	database = "database"
	password = "password"
	username = "user"
	port     = "5432"
	host     = "127.0.0.1"
)

func mustStartPostgresContainer() (func(context.Context) error, error) {
	var (
		dbName = "database"
		dbPwd  = "password"
		dbUser = "user"
	)

	dbContainer, err := postgres.Run(
		context.Background(),
		"postgres:latest",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPwd),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, err
	}

	database = dbName
	password = dbPwd
	username = dbUser

	dbHost, err := dbContainer.Host(context.Background())
	fmt.Println("Created database host", dbHost)
	if err != nil {
		return dbContainer.Terminate, err
	}

	dbPort, err := dbContainer.MappedPort(context.Background(), "5432/tcp")
	if err != nil {
		return dbContainer.Terminate, err
	}

	host = dbHost
	port = dbPort.Port()

	return dbContainer.Terminate, err
}

func TestMain(m *testing.M) {
	teardown, err := mustStartPostgresContainer()
	if err != nil {
		log.Fatalf("could not start postgres container: %v", err)
	}

	m.Run()

	if teardown != nil && teardown(context.Background()) != nil {
		log.Fatalf("could not teardown postgres container: %v", err)
	}
}

func TestNew(t *testing.T) {
	srv := New(DatabaseConfig{
		Host:         host,
		Port:         port,
		DBUserName:   username,
		DBName:       database,
		DBPassword:   password,
		DBSchema:     "public",
		MaxOpenConns: 30,
		MaxIdleConns: 30,
		MaxIdleTime:  "15m",
	})
	if srv == nil {
		t.Fatal("New() returned nil")
	}
}

func TestHealth(t *testing.T) {
	srv := New(DatabaseConfig{
		Host:         host,
		Port:         port,
		DBUserName:   username,
		DBName:       database,
		DBPassword:   password,
		DBSchema:     "public",
		MaxOpenConns: 30,
		MaxIdleConns: 30,
		MaxIdleTime:  "15m",
	})

	stats := srv.Health()

	if stats["status"] != "up" {
		t.Fatalf("expected status to be up, got %s", stats["status"])
	}

	if _, ok := stats["error"]; ok {
		t.Fatalf("expected error not to be present")
	}

	if stats["message"] != "It's healthy" {
		t.Fatalf("expected message to be 'It's healthy', got %s", stats["message"])
	}
}

func TestClose(t *testing.T) {
	srv := New(DatabaseConfig{
		Host:         host,
		Port:         port,
		DBUserName:   username,
		DBName:       database,
		DBPassword:   password,
		DBSchema:     "public",
		MaxOpenConns: 30,
		MaxIdleConns: 30,
		MaxIdleTime:  "15m",
	})

	if srv.Close() != nil {
		t.Fatalf("expected Close() to return nil")
	}
}
