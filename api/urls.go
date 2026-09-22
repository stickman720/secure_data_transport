package api

import (
	"github.com/gin-gonic/gin"
	sqlcpkg"secure_data_transport/sqlcpkg"

	"secure_data_transport/widget"
)


func InitAuthRoutes(r *gin.RouterGroup , quesries *sqlcpkg.Queries, cach *widget.Cach , mailer widget.Mailer) {
	ah := NewAuthHandler(quesries, *cach, mailer)

	r.POST("/register/",ah.registerhandler)
	r.POST("/verifyemail/",ah.VerifyEmail)

}