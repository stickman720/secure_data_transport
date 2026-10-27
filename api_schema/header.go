package api_schema

type AuthHeader struct {
	JWTToken string `header:"jwt-token"`
	RefreshToken string `header:"refresh-token"`
}