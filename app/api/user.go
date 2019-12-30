package api


type User struct {
	ID string
	OrganisationId string
	CreatedAt int64
	UpdatedAt int64
	Disabled bool
	Role int
}