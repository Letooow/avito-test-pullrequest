package usecase

import (
	"avito-test/internal/domain"
	"context"
	"errors"
)

type User struct {
	userRepository         UserRepository
	requestOwnerRepository RequestOwnerRepository
	pullRequestRepository  PullRequestRepository
}

func NewUser(userRepository UserRepository, requestOwnerRepository RequestOwnerRepository, pullRequestRepository PullRequestRepository) User {
	return User{
		userRepository:         userRepository,
		requestOwnerRepository: requestOwnerRepository,
		pullRequestRepository:  pullRequestRepository,
	}
}

func (u *User) SetActive(ctx context.Context, userID string, active bool) (*domain.User, error) {
	if u.userRepository == nil {
		return nil, ErrUserRepositoryNotFound
	} else if u.requestOwnerRepository == nil {
		return nil, ErrRequestOwnerRepositoryNotFound
	} else if u.pullRequestRepository == nil {
		return nil, ErrPullRequestRepositoryNotFound
	}

	if u.userRepository == nil {
		return nil, ErrMemberNotFound
	}
	user, err := u.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		return nil, ErrMemberNotFound
	}
	user.IsActive = active
	err = u.userRepository.UpdateUser(ctx, user)
	if err != nil {
		return nil, ErrMemberNotFound
	}
	return user, nil
}

func (u *User) GetUserPullRequests(ctx context.Context, userID string) ([]domain.PullRequest, error) {
	if u.userRepository == nil {
		return nil, ErrUserRepositoryNotFound
	} else if u.requestOwnerRepository == nil {
		return nil, ErrRequestOwnerRepositoryNotFound
	} else if u.pullRequestRepository == nil {
		return nil, ErrPullRequestRepositoryNotFound
	}
	_, err := u.userRepository.GetUserByID(ctx, userID)
	if errors.Is(err, ErrMemberNotFound) {
		return nil, ErrMemberNotFound
	} else if err != nil {
		return nil, err
	}
	pr, err := u.requestOwnerRepository.GetRequestsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	result := make([]domain.PullRequest, 0, len(pr))
	for _, reqOwner := range pr {
		if reqOwner.Role == domain.UserRoleReviewer {
			gottenPr, err := u.pullRequestRepository.GetPullRequestByID(ctx, reqOwner.RequestID)
			if err != nil {
				return nil, errors.Join(ErrPullRequestNotFound, err)
			}
			result = append(result, *gottenPr)
		}
	}
	return result, nil
}
