package api

import (
	"secure_data_transport/core"
	sqlcpkg "secure_data_transport/sqlcpkg"

	"github.com/gin-gonic/gin"

	"secure_data_transport/widget"
	"secure_data_transport/config"
)


func InitAuthRoutes(r *gin.RouterGroup , quesries *sqlcpkg.Queries, cach *widget.Cach , mailer widget.Mailer , security core.Security , conf *config.Config) {
	ah := NewAuthHandler(quesries, *cach, mailer , security , conf)

	r.POST("/register/"        , ah.registerhandler)
	r.POST("/verifyemail/"     , ah.VerifyEmail)
	r.POST("/login"            , ah.Login)
	r.POST("/forgot-pass"      , ah.ForgotPass)
	r.POST("/change-password"  , ah.ChangePassword)
	r.POST("/refresh-login"    , ah.RefreshLogin)

	r.GET("/change-password/:token" , ah.GetChangePasswordTk)
}



func InitTeamRoutes(r *gin.RouterGroup , quesries *sqlcpkg.Queries, cach *widget.Cach , mailer widget.Mailer , security core.Security , conf *config.Config) {
	th := NewTeamHandler(quesries, *cach, mailer , security , conf)

	r.POST("/" , th.CreateTeam)
	r.POST("/user/owner/" , th.CreateUser)
	r.POST("/user/me/login" , th.LoginMeTeamUser)

	r.GET("/" , th.GetTeams)
	r.GET("/user/owner" , th.GetOwnerTeamUsers)

	r.DELETE("/" , th.DeleteTeam)
	r.DELETE("/user/owner/" , th.DeleteUser)
}