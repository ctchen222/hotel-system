package pgtypes

type PgAuthParams struct {
	Email    string `json: "email"`
	Password string `json: "password`
}

type PgAuthResponse struct {
	User  *PGUser `json:"user,omitempty"`
	Token string  `json:"token,omitempty"`
}
