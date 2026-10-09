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
	"log"
	"secure_data_transport/core"
	"secure_data_transport/config"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"errors"
)




 



type AuthHandler struct{
	quesries *sqlcpkg.Queries
	cach widget.Cach
	mailer widget.Mailer
	security core.Security
	conf *config.Config
}

func NewAuthHandler(queries *sqlcpkg.Queries, cach widget.Cach , mailer widget.Mailer , security core.Security ,conf *config.Config) *AuthHandler{
	return &AuthHandler{quesries: queries, cach: cach , mailer: mailer,security: security , conf: conf}
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




	
	passwordhash := ah.security.HashPassword(input.Password)

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
	if err != nil {
		c.JSON(http.StatusInternalServerError , api_schema.ErrorResponse{
			Code: "INTERNAL_SERVER_ERROR",
			Error: "unable to cleane refresh token",
		})
		return
	}
	
	crt , ccr := context.WithTimeout(c , time.Second)
	defer ccr()
	_ , err =ah.quesries.CreateRefreshToken(crt,sqlcpkg.CreateRefreshTokenParams{
		Token: refresh_token,
		Userid: user.ID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError , api_schema.ErrorResponse{
			Code: "INTERNAL_SERVER_ERROR",
			Error: "unable to create refresh token",
		})
		return
	}
	

	c.JSON(http.StatusOK , api_schema.LoginResponse{
		Token: token,
		RefreshToken: refresh_token,
	})
	
} 



func (ah *AuthHandler)ForgotPass(c *gin.Context){
	input := &api_schema.ForgetPassRequest{}
	c.ShouldBindJSON(input)

	user , err := ah.quesries.GetUserByEmail(c , input.Email)
	if err != nil {
		if err.Error() == "no rows in result set"{
		c.JSON(http.StatusUnauthorized , api_schema.ErrorResponse{
			Code: "UNAUTHORIZED",
			Error: "emali is not exists",
		})
		return
		}
		c.JSON(http.StatusInternalServerError , api_schema.ErrorResponse{
			Code: "INTERNAL_SERVER_ERROR",
			Error: "unable to connect to database",
		})
		return
	}
	

	pass_change_tk := "change-pass|"+uuid.New().String()

	err  = ah.cach.Set(c , pass_change_tk , cache_schema.PassChangeCache{UserID: user.ID} , time.Minute * 15)
	if err != nil {
		c.JSON(http.StatusInternalServerError , api_schema.ErrorResponse{
			Code: "INTERNAL_SERVER_ERROR",
			Error: "unable to caching",
		})
		return
	}
	err = ah.mailer.SendMail([]string{input.Email} , "change pass token :" , pass_change_tk)
	if err != nil {
		_ =ah.cach.DropKey(c, pass_change_tk)

		c.JSON(http.StatusInternalServerError , api_schema.ErrorResponse{
			Code: "INTERNAL_SERVER_ERROR",
			Error: "unable to send email",
		})
		return
	}
	
	c.JSON(http.StatusOK , api_schema.ForgotPassResponse{
		Message: "change pass token sent for your email",
	})
}



func (h *AuthHandler) ChangePassword(c *gin.Context) {
	ctx := c.Request.Context()

	// ---- 1. Bind & validate ----
	var req api_schema.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}



	// ---- 3. Look up token in cache ----
	var payload cache_schema.PassChangeCache
	if err := h.cach.Get(ctx, req.Token, &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or expired token"})
		return
	}

	// ---- 4. Hash new password ----
	hashed := h.security.HashPassword(req.NewPassword)
	

	// ---- 5. Update DB ----
	if err := h.quesries.UpdateUserPassword(ctx, sqlcpkg.UpdateUserPasswordParams{
		PasswordHash: hashed,
		ID:           payload.UserID,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update password"})
		return
	}

	// ---- 6. Invalidate token (single-use) ----
	if err := h.cach.DropKey(ctx, req.Token); err != nil {
		log.Printf("failed to drop pass change token: %v", err)
	}

	// ---- 7. Respond ----
	c.JSON(http.StatusOK, api_schema.ChangePasswordResponse{
		Message: "password updated successfully",
	})
}



func (h *AuthHandler) RefreshLogin(c *gin.Context) {
	// 1. Parse & validate body
	var req api_schema.RefreshLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api_schema.ErrorResponse{
			Code:  "BAD_REQUEST",
			Error: "bad request",
		})
		return
	}

	if req.Token == "" {
		c.JSON(http.StatusBadRequest, api_schema.ErrorResponse{
			Code:  "BAD_REQUEST",
			Error: "bad request",
		})
		return
	}

	// 2. Look up the refresh token and make sure the user is active
	user, err := h.quesries.ValidRefreshToken(c.Request.Context(), sqlcpkg.ValidRefreshTokenParams{
		Token: req.Token,
		Column2: int32(h.conf.RefreshTkExpireTime),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusUnauthorized, api_schema.ErrorResponse{
				Code:  "UNAUTHORIZED",
				Error: "invalid token",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, api_schema.ErrorResponse{
			Code:  "INTERNAL_SERVER_ERROR",
			Error: "internal server error",
		})
		return
	}

	// 3. Issue a new JWT
	jwtToken, err := h.security.CreateJWT(user.Username, user.Role , time.Minute*15)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api_schema.ErrorResponse{
			Code:  "INTERNAL_SERVER_ERROR",
			Error: "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, api_schema.RefrshLoginResponse{
		JWT_Token: jwtToken,
	})
}