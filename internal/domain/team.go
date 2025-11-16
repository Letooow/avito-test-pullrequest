package domain

// Team - Команда с уникальным именем.
type Team struct {
	// Name - название команды
	Name string `json:"team_name"`
	// Members - участники команды
	Members []User `json:"members"`
}
