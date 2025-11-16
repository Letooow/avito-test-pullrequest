package http

import (
	openapi "avito-test/internal/gen/go/go"
	"avito-test/internal/usecase"
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/gin-gonic/gin"
)

type Server struct {
	host   string
	port   uint16
	router *gin.Engine
}

type UseCases struct {
	User        usecase.User
	Team        usecase.Team
	PullRequest usecase.PullRequest
}

func NewServer(useCases UseCases, options ...func(*Server)) *Server {
	r := gin.Default()

	setupRouter(r, useCases)

	s := &Server{router: r, host: "localhost", port: 8080}
	for _, o := range options {
		o(s)
	}

	return s
}

func WithHost(host string) func(*Server) {
	return func(s *Server) {
		s.host = host
	}
}

func WithPort(port uint16) func(*Server) {
	return func(s *Server) {
		s.port = port
	}
}

func (s *Server) Run(ctx context.Context) error {
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.host, s.port),
		Handler: s.router,
	}

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})

	eg.Go(func() error {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	})

	return eg.Wait()
}

func setupRouter(r *gin.Engine, uc UseCases) {
	handlers := openapi.ApiHandleFunctions{
		PullRequestsAPI: openapi.NewPullRequestsAPI(uc.PullRequest),
		TeamsAPI:        openapi.NewTeamsAPI(uc.Team),
		UsersAPI:        openapi.NewUsersAPI(uc.User),
	}

	openapi.NewRouterWithGinEngine(r, handlers)
}
