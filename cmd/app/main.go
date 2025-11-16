package main

import (
	"avito-test/internal/db"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	gateway "avito-test/internal/gateway/http"
	pr "avito-test/internal/repository/pull_request/postgres"
	tr "avito-test/internal/repository/team/postgres"
	ur "avito-test/internal/repository/user/postgres"
	"avito-test/internal/usecase"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func setupDB(ctx context.Context) (*sql.DB, error) {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	pass := getEnv("DB_PASSWORD", "postgres")
	name := getEnv("DB_NAME", "avito")
	ssl := getEnv("DB_SSLMODE", "disable")

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, pass, host, port, name, ssl,
	)

	dbConn, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	dbConn.SetMaxOpenConns(10)
	dbConn.SetMaxIdleConns(5)
	dbConn.SetConnMaxLifetime(30 * time.Minute)

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := dbConn.PingContext(pingCtx); err != nil {
		_ = dbConn.Close()
		return nil, fmt.Errorf("database.PingContext: %w", err)
	}

	log.Printf("data base connected: %s", dsn)
	return dbConn, nil
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	ctx := context.Background()

	conn, err := setupDB(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	database := db.New(conn)

	userRepo := ur.NewUserRepository(database)
	teamRepo := tr.NewTeamRepository(database)
	prRepo := pr.NewPullRequestRepository(database)
	reqOwnerRepo := ur.NewRequestOwnerRepository(database)

	prUC := usecase.NewPullRequest(prRepo, teamRepo, userRepo, reqOwnerRepo)
	teamUC := usecase.NewTeam(teamRepo, userRepo)
	userUC := usecase.NewUser(userRepo, reqOwnerRepo, prRepo)

	usecases := gateway.UseCases{
		User:        userUC,
		Team:        teamUC,
		PullRequest: prUC,
	}

	server := gateway.NewServer(usecases,
		gateway.WithHost("localhost"),
		gateway.WithPort(8080),
	)

	if err := server.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
