// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package models

type NewOrganisation struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

type NewUser struct {
	FirstName      string  `json:"firstName"`
	LastName       string  `json:"lastName"`
	Email          string  `json:"email"`
	Password       string  `json:"password"`
	OrganisationID *string `json:"organisationId"`
	Role           int     `json:"Role"`
}

type Organisation struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	UserPoolID  string  `json:"userPoolId"`
	WorkerLimit int     `json:"workerLimit"`
	UserLimit   int     `json:"userLimit"`
	Disabled    *bool   `json:"disabled"`
	CreatedAt   int     `json:"created_at"`
	UpdatedAt   *int    `json:"updated_at"`
}

type OrganisationUsers struct {
	OrganisationID *string `json:"organisationId"`
	UserID         *string `json:"userId"`
}

type Token struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiredAt    int    `json:"expiredAt"`
}

type User struct {
	ID             string  `json:"id"`
	FirstName      string  `json:"firstName"`
	LastName       string  `json:"lastName"`
	Email          string  `json:"email"`
	Password       string  `json:"password"`
	OrganisationID *string `json:"organisationId"`
	CreatedAt      int     `json:"created_at"`
	UpdatedAt      *int    `json:"updatedAt"`
	Disabled       *bool   `json:"disabled"`
	Role           *int    `json:"role"`
}
