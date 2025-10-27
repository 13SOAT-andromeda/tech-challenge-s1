package session

type UserOutput struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Contact string `json:"contact"`
	Role    string `json:"role"`
	Active  bool   `json:"active"`
}

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	AccessToken  string     `json:"access_token"`
	RefreshToken string     `json:"refresh_token"`
	ExpiresIn    int64      `json:"expires_in"`
	TokenType    string     `json:"token_type"`
	User         UserOutput `json:"user"`
}

type ValidateInput struct {
	Token string
}

type ValidateOutput struct {
	Valid bool        `json:"valid"`
	User  *UserOutput `json:"user,omitempty"`
}

type RefreshInput struct {
	RefreshToken string
}

type RefreshOutput struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type LogoutInput struct {
	RefreshToken string
}

type LogoutOutput struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
