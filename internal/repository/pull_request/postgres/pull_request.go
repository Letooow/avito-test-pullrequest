package postgres

import (
	"avito-test/internal/db"
	"avito-test/internal/domain"
	"avito-test/internal/usecase"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type PullRequestRepository struct {
	db *db.Queries
}

func NewPullRequestRepository(db *db.Queries) *PullRequestRepository {
	return &PullRequestRepository{db: db}
}

func (p *PullRequestRepository) SavePullRequest(ctx context.Context, pull *domain.PullRequest) error {
	err := p.db.CreatePullRequest(ctx, db.CreatePullRequestParams{Pullrequestid: pull.ID, Name: sql.NullString{String: pull.Name, Valid: true}, Status: string(pull.Status)})
	if err != nil {
		return fmt.Errorf("save pull request: %w", err)
	}
	return nil
}

func (p *PullRequestRepository) UpdatePullRequest(ctx context.Context, pull *domain.PullRequest) error {
	err := p.db.UpdatePullRequestStatus(ctx, db.UpdatePullRequestStatusParams{Pullrequestid: pull.ID, Status: string(pull.Status), Mergedat: sql.NullTime{Time: pull.MergedAt, Valid: true}})
	if err != nil {
		return fmt.Errorf("update pull request: %w", err)
	}
	return nil
}

func (p *PullRequestRepository) GetPullRequestByID(ctx context.Context, id string) (*domain.PullRequest, error) {
	pr, err := p.db.GetPullRequestByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, usecase.ErrPullRequestAlreadyExists
	} else if err != nil {
		return nil, fmt.Errorf("can't get pull request by id: %w", err)
	}
	return &domain.PullRequest{ID: pr.Pullrequestid, Name: pr.Name.String, Status: domain.RequestStatus(pr.Status), CreatedAt: pr.Createdat}, nil
}

func (p *PullRequestRepository) GetPullRequests(ctx context.Context) ([]domain.PullRequest, error) {
	prs, err := p.db.GetPullRequests(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return []domain.PullRequest{}, nil
	} else if err != nil {
		return nil, fmt.Errorf("can't get pull requests: %w", err)
	}
	result := make([]domain.PullRequest, len(prs))
	for i, pr := range prs {
		result[i] = domain.PullRequest{ID: pr.Pullrequestid, Name: pr.Name.String, Status: domain.RequestStatus(pr.Status), CreatedAt: pr.Createdat, MergedAt: pr.Mergedat.Time}
	}
	return result, nil
}
