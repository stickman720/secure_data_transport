package api

import (
	"secure_data_transport/core"
	sqlcpkg "secure_data_transport/sqlcpkg"

	"github.com/gin-gonic/gin"

	"secure_data_transport/widget"
)


func InitAuthRoutes(r *gin.RouterGroup , quesries *sqlcpkg.Queries, cach *widget.Cach , mailer widget.Mailer , security core.Security) {
	ah := NewAuthHandler(quesries, *cach, mailer , security)

	r.POST("/register/",ah.registerhandler)
	r.POST("/verifyemail/",ah.VerifyEmail)
	r.POST("/login" , ah.Login)
}