package api

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"secure_data_transport/api_schema"
	sqlcpkg "secure_data_transport/sqlcpkg"

	"github.com/gin-gonic/gin"
)



type AuthHandler struct{
	quesries *sqlcpkg.Queries
}

func NewAuthHandler(queries *sqlcpkg.Queries) *AuthHandler{
	return &AuthHandler{quesries: queries}
}




func (ah *AuthHandler)registerhandler(c *gin.Context) {

	input := &api_schema.RegisterRequest{}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{
			"code" : "BAD_REQUEST", 
			"error": "bad request",
		})
		return
	}


	h := sha256.Sum256([]byte(input.Password + input.Email))
	passwordhash := hex.EncodeToString(h[:])

	userid , err := ah.quesries.CreateUser(c , sqlcpkg.CreateUserParams{
		Username: input.Username,
		Email: input.Email,
		PasswordHash: string(passwordhash),
	})

	if err != nil{
		c.JSON(http.StatusOK , gin.H{
			"error" : err.Error(),
		})
		return
	}

	
	user ,err := ah.quesries.GetUserByID(c , userid)

	c.JSON(200, gin.H{
		"user": user,
	})
}








