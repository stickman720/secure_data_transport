package api_schema

import (
	sqlcpkg "secure_data_transport/sqlcpkg"
)



//====================================================================
//requests
//====================================================================

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


//====================================================================
//responses
//====================================================================

type ErrorResponse struct {
	Code string  `json:"code"`
	Error string `json:"error"`
}


type RegisterResponse struct {
	Message string `json:"message"`
}


type VerifyEmailResponse struct {
	User sqlcpkg.User `json:"user"`
}


type LoginResponse struct {
	Token string `json:"token"`
	RefreshToken string `json:"refreshtoken"`
}