package usecase

import (
	"avito-test/internal/domain"
	"context"
	"errors"
	"math/rand"
)

type PullRequest struct {
	pullRequestRepository  PullRequestRepository
	teamRepository         TeamRepository
	userRepository         UserRepository
	requestOwnerRepository RequestOwnerRepository
}

func NewPullRequest(pullRequestRepo PullRequestRepository,
	teamRepository TeamRepository,
	userRepository UserRepository,
	requestOwnerRepository RequestOwnerRepository) PullRequest {
	return PullRequest{
		pullRequestRepository:  pullRequestRepo,
		teamRepository:         teamRepository,
		userRepository:         userRepository,
		requestOwnerRepository: requestOwnerRepository,
	}
}

func (p *PullRequest) CreatePullRequest(ctx context.Context, request *domain.PullRequest) (*domain.PullRequest, error) {
	if request == nil {
		return nil, ErrAuthorNotFound
	}
	if p.teamRepository == nil {
		return nil, ErrTeamRepositoryNotFound
	} else if p.userRepository == nil {
		return nil, ErrUserRepositoryNotFound
	} else if p.requestOwnerRepository == nil {
		return nil, ErrRequestOwnerRepositoryNotFound
	} else if p.pullRequestRepository == nil {
		return nil, ErrPullRequestRepositoryNotFound
	}

	author, err := p.userRepository.GetUserByID(ctx, request.AuthorID)
	if errors.Is(err, ErrMemberNotFound) {
		return nil, ErrAuthorNotFound
	} else if err != nil {
		return nil, err
	}
	if !author.IsActive {
		return nil, ErrAuthorIsInactive
	}
	_, err = p.pullRequestRepository.GetPullRequestByID(ctx, request.ID)
	if err == nil {
		return nil, ErrPullRequestAlreadyExists
	} else if !errors.Is(err, ErrPullRequestNotFound) {
		return nil, err
	}

	err = p.pullRequestRepository.SavePullRequest(ctx, request)
	if err != nil {
		return nil, err
	}

	err = p.requestOwnerRepository.SaveRequestOwner(ctx, &domain.RequestOwner{RequestID: request.ID, UserID: request.AuthorID, Role: domain.UserRoleAuthor})
	if err != nil {
		return nil, err
	}

	authorTeam, err := p.userRepository.GetTeamsByUserID(ctx, request.AuthorID)
	if errors.Is(err, ErrMemberNotFound) {
		return nil, ErrAuthorNotFound
	} else if err != nil {
		return nil, err
	}
	coworkers := make([]domain.User, 0, 10)
	for _, team := range authorTeam {
		cwrk, err := p.userRepository.GetUsersByTeamName(ctx, team.Name)
		if errors.Is(err, ErrTeamNotFound) {
			return nil, ErrTeamNotFound
		} else if err != nil {
			return nil, err
		}
		for _, user := range cwrk {
			if !user.IsActive {
				continue
			} else if user.ID == request.AuthorID {
				continue
			}
			coworkers = append(coworkers, user)
		}
	}
	rand.Shuffle(len(coworkers), func(i, j int) { coworkers[i], coworkers[j] = coworkers[j], coworkers[i] })

	assignedReviewers := make([]string, 0, 2)
	for i := 2; i > 0 && len(coworkers) > 0; i-- {
		err := p.requestOwnerRepository.SaveRequestOwner(ctx, &domain.RequestOwner{RequestID: request.ID, UserID: coworkers[0].ID, Role: domain.UserRoleReviewer})
		if err != nil {
			return nil, err
		}
		assignedReviewers = append(assignedReviewers, coworkers[0].ID)
		coworkers = coworkers[1:]
	}
	request.AssignedReviewersID = assignedReviewers
	request.Status = domain.RequestStatusOpen
	return request, nil
}

func (p *PullRequest) UpdatePullRequest(ctx context.Context, request *domain.PullRequest) (*domain.PullRequest, error) {
	if request == nil {
		return nil, ErrAuthorNotFound
	}
	if p.teamRepository == nil {
		return nil, ErrTeamRepositoryNotFound
	} else if p.userRepository == nil {
		return nil, ErrUserRepositoryNotFound
	} else if p.requestOwnerRepository == nil {
		return nil, ErrRequestOwnerRepositoryNotFound
	} else if p.pullRequestRepository == nil {
		return nil, ErrPullRequestRepositoryNotFound
	}

	req, err := p.pullRequestRepository.GetPullRequestByID(ctx, request.ID)
	if errors.Is(err, ErrPullRequestNotFound) {
		return nil, ErrPullRequestNotFound
	} else if err != nil {
		return nil, err
	}
	if req.Status == domain.RequestStatusMerged {
		return req, nil
	}
	err = p.pullRequestRepository.UpdatePullRequest(ctx, request)
	if err != nil {
		return nil, err
	}

	req, err = p.pullRequestRepository.GetPullRequestByID(ctx, request.ID)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (p *PullRequest) MergePullRequest(ctx context.Context, id string) (*domain.PullRequest, error) {
	if p.teamRepository == nil {
		return nil, ErrTeamRepositoryNotFound
	} else if p.userRepository == nil {
		return nil, ErrUserRepositoryNotFound
	} else if p.requestOwnerRepository == nil {
		return nil, ErrRequestOwnerRepositoryNotFound
	} else if p.pullRequestRepository == nil {
		return nil, ErrPullRequestRepositoryNotFound
	}
	req, err := p.pullRequestRepository.GetPullRequestByID(ctx, id)
	if errors.Is(err, ErrPullRequestNotFound) {
		return nil, ErrPullRequestNotFound
	} else if err != nil {
		return nil, err
	}
	if req.Status == domain.RequestStatusMerged {
		return req, nil
	}
	req.Status = domain.RequestStatusMerged
	req, err = p.UpdatePullRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (p *PullRequest) ReassignRequest(ctx context.Context, requestID, userID string) (*domain.PullRequest, *domain.User, error) {
	if p.teamRepository == nil {
		return nil, nil, ErrTeamRepositoryNotFound
	} else if p.userRepository == nil {
		return nil, nil, ErrUserRepositoryNotFound
	} else if p.requestOwnerRepository == nil {
		return nil, nil, ErrRequestOwnerRepositoryNotFound
	} else if p.pullRequestRepository == nil {
		return nil, nil, ErrPullRequestRepositoryNotFound
	}

	pr, err := p.pullRequestRepository.GetPullRequestByID(ctx, requestID)
	if errors.Is(err, ErrPullRequestNotFound) {
		return nil, nil, ErrPullRequestNotFound
	} else if err != nil {
		return nil, nil, err
	}

	if pr.Status == domain.RequestStatusMerged {
		return nil, nil, nil
	}

	author, err := p.userRepository.GetUserByID(ctx, pr.AuthorID)
	if errors.Is(err, ErrMemberNotFound) {
		return nil, nil, ErrAuthorNotFound
	} else if err != nil {
		return nil, nil, err
	}
	if !author.IsActive {
		return nil, nil, ErrAuthorIsInactive
	}

	if err := p.requestOwnerRepository.DeleteRequestOwner(ctx, &domain.RequestOwner{
		RequestID: pr.ID,
		UserID:    userID,
		Role:      domain.UserRoleReviewer,
	}); err != nil {
		return nil, nil, err
	}

	owners, err := p.requestOwnerRepository.GetUsersByPullRequestID(ctx, pr.ID)
	if err != nil {
		return nil, nil, err
	}
	assigned := make(map[string]struct{}, len(owners))
	for _, o := range owners {
		assigned[o.UserID] = struct{}{}
	}

	authorTeam, err := p.userRepository.GetTeamsByUserID(ctx, pr.AuthorID)
	if errors.Is(err, ErrMemberNotFound) {
		return nil, nil, ErrAuthorNotFound
	} else if err != nil {
		return nil, nil, err
	}
	coworkers, err := p.userRepository.GetUsersByTeamName(ctx, authorTeam[0].Name)
	if errors.Is(err, ErrTeamNotFound) {
		return nil, nil, ErrTeamNotFound
	} else if err != nil {
		return nil, nil, err
	}

	candidates := make([]domain.User, 0, len(coworkers))
	for _, u := range coworkers {
		if !u.IsActive {
			continue
		}
		if u.ID == author.ID {
			continue
		}
		if _, used := assigned[u.ID]; used {
			continue
		}
		candidates = append(candidates, u)
	}

	if len(candidates) == 0 {
		return nil, nil, ErrCannotFindActiveMembers
	}

	newReviewer := candidates[rand.Intn(len(candidates))]

	if err := p.requestOwnerRepository.SaveRequestOwner(ctx, &domain.RequestOwner{
		RequestID: pr.ID,
		UserID:    newReviewer.ID,
		Role:      domain.UserRoleReviewer,
	}); err != nil {
		return nil, nil, err
	}

	pr, err = p.pullRequestRepository.GetPullRequestByID(ctx, requestID)
	if err != nil {
		return nil, nil, err
	}

	return pr, &newReviewer, nil
}
