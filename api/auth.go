package api

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/http"
	"secure_data_transport/api_schema"
	"secure_data_transport/cache_schema"
	sqlcpkg "secure_data_transport/sqlcpkg"
	"secure_data_transport/widget"
	"time"

	"secure_data_transport/core"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)




 



type AuthHandler struct{
	quesries *sqlcpkg.Queries
	cach widget.Cach
	mailer widget.Mailer
	security core.Security
}

func NewAuthHandler(queries *sqlcpkg.Queries, cach widget.Cach , mailer widget.Mailer , security core.Security) *AuthHandler{
	return &AuthHandler{quesries: queries, cach: cach , mailer: mailer,security: security}
}




func (ah *AuthHandler)registerhandler(c *gin.Context) {

	input := &api_schema.RegisterRequest{}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code" : "BAD_REQUEST", 
			"error": "bad request",
		})
		return
	}

	toutue , cancelue:= context.WithTimeout(c, 3 * time.Second)
	defer cancelue()
	user_exists , err := ah.quesries.UserExists(toutue, sqlcpkg.UserExistsParams{
		Username: input.Username,
		Email: input.Email,
	})

	if err != nil{
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError , gin.H{
			"error" : "unable to complete request, please try again later",
		})
		return
	}

	if user_exists {
		c.JSON(http.StatusConflict , gin.H{
			"error" : "user with this username or email already exists",
		})
		return
	}




	h := sha256.Sum256([]byte(input.Password))
	passwordhash := hex.EncodeToString(h[:])

	userid , err := ah.quesries.CreateUser(c , sqlcpkg.CreateUserParams{
		Username: input.Username,
		Email: input.Email,
		PasswordHash: string(passwordhash),
		Role: "owner",
	})

	if err != nil{
		c.JSON(http.StatusInternalServerError , gin.H{
			"error" : "unable to insert user into database, please try again later",
		})
		return
	}


	vcode := uuid.New().String()
	ch := cache_schema.EmailVerificationCache{
		UserId: userid,
		VerificationCode: vcode,
	}

	err = ah.cach.Set(c, vcode, ch, 15 * time.Minute)
	if err != nil{
		c.JSON(http.StatusInternalServerError , gin.H{
			"error" : "unable to cache vrifying code, please try again later",
		})
		ah.quesries.DeleteUser(c, userid)
		return
	}


	

	msg := fmt.Sprintf("Your verification code is: %s" , vcode)
	err = ah.mailer.SendMail([]string{input.Email}, "Email Verification", msg)
	if err != nil{

	}
	c.JSON(http.StatusOK , gin.H{
		"message" : "verification code sent to your email",
	})
	
}



func (ah *AuthHandler)VerifyEmail(c *gin.Context){

	input := &api_schema.VerifyEmailRequest{}

	if err := c.ShouldBindJSON(input); err != nil{
		c.JSON(http.StatusBadRequest , api_schema.ErrorResponse{
			Code : "BAD_REQUEST", 
			Error: "bad request",
		})
		return
	}

	ch := cache_schema.EmailVerificationCache{}
	err := ah.cach.Get(c , input.Code , &ch)

	
	
	if err != nil{
		if err.Error() == "key dosnt exist or expired"{
			c.JSON(http.StatusNonAuthoritativeInfo , api_schema.ErrorResponse{
				Code:"UNAUTHORISE",
				Error:"invalid code or expired code",
			})
			return
		}	
		c.JSON(http.StatusInternalServerError , api_schema.ErrorResponse{
			Code:"INTERNAL_SERVER_ERROR",
			Error:"internal server error 1",
		})
		return
	}
	
	user , err := ah.quesries.GetUserByID(c , ch.UserId)
	if err != nil{
		c.JSON(http.StatusInternalServerError , api_schema.ErrorResponse{
			Code:"INTERNAL_SERVER_ERROR",
			Error:"internal server error 2",
		})
		return
	}

	if user.Isactive.Bool {
		c.JSON(http.StatusOK , api_schema.RegisterResponse{
			Message:"user already verified",
		})
		return
	}

	user , err = ah.quesries.VerifyUser(c , ch.UserId)
	if err != nil{
		c.JSON(http.StatusInternalServerError , api_schema.ErrorResponse{
			Code:"INTERNAL_SERVER_ERROR",
			Error:"internal server error 3",
		})
		return
	}
	c.JSON(http.StatusOK , api_schema.VerifyEmailResponse{
		User:user,
	})

}



func (ah * AuthHandler) Login (c *gin.Context ){
	input := &api_schema.LoginRequest{}
	err := c.ShouldBindJSON(input)

	if err != nil {
		c.JSON(http.StatusBadRequest , api_schema.ErrorResponse{
			Code : "BAD_REQUEST",
			Error: "bad request",
		})
		return
	}


	h := sha256.Sum256([]byte(input.Password))
	passwordhash := hex.EncodeToString(h[:])

	lt , ltc := context.WithTimeout(c , time.Second)
	defer ltc()
	user , err := ah.quesries.Login(lt , sqlcpkg.LoginParams{
		Username: input.Username,
		PasswordHash: passwordhash,
	})

	if err != nil {
		if err == sql.ErrNoRows{
			c.JSON(http.StatusUnauthorized , api_schema.ErrorResponse{
				Code : "UNAUTHORIZED",
				Error: "username or password is not corect",
			})
			return
		}
		c.JSON(http.StatusInternalServerError , api_schema.ErrorResponse{
			Code : "INTERNAL_SERVER_ERROR",
			Error: "internal server error",
		})
		return
	}

	token , err   := ah.security.CreateJWT(user.Username , user.Role , time.Minute * 15)
	if err != nil {
		c.JSON(http.StatusInternalServerError , api_schema.ErrorResponse{
			Code: "INTERNAL_SERVER_ERROR",
			Error: "internal server error",
		})
		return
	}

	refresh_token := "refresh-tk->"+uuid.New().String()

	err = ah.quesries.CleanRefreshToken(c , user.Username)
	
	crt , ccr := context.WithTimeout(c , time.Second)
	defer ccr()
	ah.quesries.CreateRefreshToken(crt,sqlcpkg.CreateRefreshTokenParams{
		Token: refresh_token,
		Userid: user.ID,
	})

	

	c.JSON(http.StatusOK , api_schema.LoginResponse{
		Token: token,
		RefreshToken: refresh_token,
	})
	
} 





