package api


type User struct {
	ID string
	FirstName string
	LastName string
	Email string
	Password string
	OrganisationId string
	CreatedAt int64
	UpdatedAt int64
	Disabled bool
	Role int
}