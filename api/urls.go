package api

import (
	"github.com/gin-gonic/gin"
	sqlcpkg"secure_data_transport/sqlcpkg"
)


func InitAuthRoutes(r *gin.RouterGroup , quesries *sqlcpkg.Queries) {
	ah := NewAuthHandler(quesries)

	r.POST("/register/",ah.registerhandler)


}