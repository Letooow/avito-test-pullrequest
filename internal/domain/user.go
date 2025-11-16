package domain

// User - участник команды с уникальным идентификатором, именем и флагом активности `isActive`.
type User struct {
	// ID - id участника
	ID string `json:"id"`
	// UserName - Никнейм участника
	Username string `json:"username"`
	//  IsActive - Активен ли участник
	IsActive bool `json:"is_active"`
}

type Role string

const (
	UserRoleReviewer Role = "reviewer"
	UserRoleAuthor   Role = "author"
)

type RequestOwner struct {
	UserID    string `json:"user_id"`
	RequestID string `json:"request_id"`
	Role      Role   `json:"role"`
}
