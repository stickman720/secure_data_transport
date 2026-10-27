package api_schema


import (
	sqlcpkg"secure_data_transport/sqlcpkg"
)



//========================================================
//requests
//========================================================
type CreateTeamRequest struct {
	Team TeamInput `json:"team"`
}



type GetTeamsRequest struct {
	Title string `json:"title"`
	Username string `json:"username"`
}



type DeleteTeamRequest struct {
	Title string `json:"title"`
}



type CreateTeamUserRequest struct{
	Username string `json:"username"`
	Title string `json:"title"`
}



type DeleteTeamUserRequest struct {
	Username string `json:"username"`
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



type CreateTeamUserResponse struct {
	User sqlcpkg.Teamuser `json:"user"`
	Token string `json:"user-token"`
}



type DeleteTeamUserResponse struct {
	Message string `json:"message"`
}



type GetTeamUsersResponse struct {
	Users []sqlcpkg.GetAllTeamUsersRow `json:"users"`
}