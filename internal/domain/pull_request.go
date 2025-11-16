package domain

import "time"

type RequestStatus string

const (
	RequestStatusOpen   RequestStatus = "OPEN"
	RequestStatusMerged RequestStatus = "MERGED"
)

// PullRequest - сущность с идентификатором, названием, автором, статусом `OPEN|MERGED`и списком назначенных ревьюверов (до 2).
type PullRequest struct {
	// ID - id реквеста
	ID string `json:"id" db:"PullRequestID"`
	// Name - название реквеста
	Name string `json:"name" db:"Name"`
	// AuthorID - создатель реквеста
	AuthorID string `json:"author" db:"AuthorID"`
	// Status - текущий статус реквеста
	Status RequestStatus `json:"status" db:"Status"`
	// AssignedReviewersID - прикрепленные проверяющие
	AssignedReviewersID []string `json:"assigned_reviewers"`
	// CreatedAt - время создания
	CreatedAt time.Time `json:"created_at"`
	// MergedAt - время слияние
	MergedAt time.Time `json:"merged_at"`
}
