// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type CancelTokenInput struct {
	Token string `json:"token"`
}

type CreateUserInput struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	ReCaptcha string `json:"reCaptcha"`
}

type CreateUserOutput struct {
	User  *User  `json:"user"`
	Token *Token `json:"token"`
}

type Login struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	RememberMe bool   `json:"rememberMe"`
	ReCaptcha  string `json:"reCaptcha"`
}

type RefreshTokenInput struct {
	Token string `json:"token"`
}

type Token struct {
	Value string `json:"value"`
}

type User struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}
