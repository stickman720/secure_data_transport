package api

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"secure_data_transport/api_schema"
	sqlcpkg "secure_data_transport/sqlcpkg"
	"secure_data_transport/widget"
	"time"
	"github.com/gin-gonic/gin"
	"secure_data_transport/cache_schema"
	"github.com/google/uuid"
	"fmt"
	
)




 



type AuthHandler struct{
	quesries *sqlcpkg.Queries
	cach widget.Cach
	mailer widget.Mailer
}

func NewAuthHandler(queries *sqlcpkg.Queries, cach widget.Cach , mailer widget.Mailer) *AuthHandler{
	return &AuthHandler{quesries: queries, cach: cach , mailer: mailer}
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




	h := sha256.Sum256([]byte(input.Password + input.Email))
	passwordhash := hex.EncodeToString(h[:])

	userid , err := ah.quesries.CreateUser(c , sqlcpkg.CreateUserParams{
		Username: input.Username,
		Email: input.Email,
		PasswordHash: string(passwordhash),
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
		c.JSON(http.StatusBadRequest , gin.H{
			"code" : "BAD_REQUEST", 
			"error": "bad request",
		})
		return
	}

	ch := cache_schema.EmailVerificationCache{}
	err := ah.cach.Get(c , input.Code , &ch)

	
	
	if err != nil{
		if err.Error() == "key dosnt exist or expired"{
			c.JSON(http.StatusNonAuthoritativeInfo , gin.H{
				"code":"UNAUTHORISE",
				"error":"invalid code or expired code",
			})
			return
		}	
		c.JSON(http.StatusInternalServerError , gin.H{
			"code":"INTERNAL_SERVER_ERROR",
			"error":"internal server error 1",
		})
		return
	}
	
	user , err := ah.quesries.GetUserByID(c , ch.UserId)
	if err != nil{
		c.JSON(http.StatusInternalServerError , gin.H{
			"code":"INTERNAL_SERVER_ERROR",
			"error":"internal server error 2",
		})
		return
	}

	if user.Isactive.Bool {
		c.JSON(http.StatusOK , gin.H{
			"message":"user already verified",
		})
		return
	}

	user , err = ah.quesries.VerifyUser(c , ch.UserId)
	if err != nil{
		c.JSON(http.StatusInternalServerError , gin.H{
			"code":"INTERNAL_SERVER_ERROR",
			"error":"internal server error 3",
		})
		return
	}
	c.JSON(http.StatusOK , gin.H{
		"user":user,
	})

}









