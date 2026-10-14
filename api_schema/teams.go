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



type GetTeamsRequest struct {
	Title string `json:"title"`
	Username string `json:"username"`
}



type DeleteTeamRequest struct {
	Title string `json:"title"`
}

//========================================================
//responses
//========================================================

type CreateTeamResponse struct {
	Message string `json:"message"`
	Team sqlcpkg.Team `json:"team"`
}



type GetTeamsResponse struct {
	Teams []sqlcpkg.Team `json:"teams"`
}



type DeleteTeamResponse struct {
    Team sqlcpkg.Team `json:"team"`
}