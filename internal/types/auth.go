package types

type AuthParams struct {
	Email    string `json: "email"`
	Password string `json: "password`
}

type AuthResponse struct {
	User  *User  `json:"user,omitempty"`
	Token string `json:"token,omitempty"`
}
