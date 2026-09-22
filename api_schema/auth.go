package api_schema




type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}


type VerifyEmailRequest struct {
	Code string `json:"code"`
}



type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

