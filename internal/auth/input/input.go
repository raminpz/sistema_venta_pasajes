package input

// LoginInput contiene las credenciales de acceso.
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RefreshInput contiene el refresh token para renovar el access token.
type RefreshInput struct {
	RefreshToken string `json:"refresh_token"`
}

// TokenPairOutput es la respuesta al hacer login o refresh.
type TokenPairOutput struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"` // segundos hasta que expira el access token
}
