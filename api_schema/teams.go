package api_schema


import (
	sqlcpkg"secure_data_transport/sqlcpkg"
)



//========================================================
//requests
//========================================================
type CreateTeamRequest struct {
	JWTToken string `json:"token"`
	Team TeamInput `json:"team"`
}
//========================================================
//responses
//========================================================

type CreateTeamResponse struct {
	Message string `json:"message"`
	Team sqlcpkg.Team `json:"team"`
}
