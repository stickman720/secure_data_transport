package api

import (
	"context"
	"errors"
	"fmt"

	"net/http"
	"secure_data_transport/api_schema"
	"secure_data_transport/config"
	"secure_data_transport/core"
	sqlcpkg "secure_data_transport/sqlcpkg"
	"secure_data_transport/widget"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)



type TeamHandler struct{
	quesries *sqlcpkg.Queries
	cach widget.Cach
	mailer widget.Mailer
	security core.Security
	conf *config.Config
}

func NewTeamHandler(queries *sqlcpkg.Queries, cach widget.Cach , mailer widget.Mailer , security core.Security ,conf *config.Config) *TeamHandler{
	return &TeamHandler{quesries: queries, cach: cach , mailer: mailer,security: security , conf: conf}
} 



func (h *TeamHandler) CreateTeam(c *gin.Context) {
	// --- 1. Parse body ---
	var req api_schema.CreateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api_schema.ErrorResponse{
			Code:  "BAD_REQUEST",
			Error: "bad request",
		})
		return
	}

	// --- 2. Validate fields ---
	req.Team.Title = strings.TrimSpace(req.Team.Title)
	if req.JWTToken == "" || req.Team.Title == "" {
		c.JSON(http.StatusBadRequest, api_schema.ErrorResponse{
			Code:  "BAD_REQUEST",
			Error: "bad request",
		})
		return
	}
	if req.Team.UserCount <= 0 ||
		req.Team.RateLimitPerUser <= 0 ||
		req.Team.StorageSpace <= 0 ||
		req.Team.Transfer <= 0 {
		c.JSON(http.StatusBadRequest, api_schema.ErrorResponse{
			Code:  "BAD_REQUEST",
			Error: "bad request",
		})
		return
	}

	// --- 3. Authenticate ---
	claims, err := h.security.ParseToken(req.JWTToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, api_schema.ErrorResponse{
			Code:  "UNAUTHORIZED",
			Error: "invalid token",
		})
		return
	}


	// --- 5. Insert ---
	team, err := h.quesries.CreateTeam(c.Request.Context(), sqlcpkg.CreateTeamParams{
		Title:            req.Team.Title,
		UserCount:        int32(req.Team.UserCount),
		RateLimitPerUser: int32(req.Team.RateLimitPerUser),
		StorageSpace:     int32(req.Team.StorageSpace),
		Transfer:       int32(req.Team.Transfer),
		Userid:           claims.UserID,
	})
	if err != nil {
		// Duplicate title -> 409
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			c.JSON(http.StatusConflict, api_schema.ErrorResponse{
				Code:  "TEAM_ALREDY_EXIST",
				Error: "team with this title is alredy exist",
			})
			return
		}

		// Anything else -> 500 (never leak DB details)
		c.JSON(http.StatusInternalServerError, api_schema.ErrorResponse{
			Code:  "INTERNAL_SERVER_ERROR",
			Error: "internal server error",
		})
		return
	}

	// --- 6. Success ---
	c.JSON(http.StatusOK, api_schema.CreateTeamResponse{
		Message: "team created",
		Team:    team,
	})
}



func (h *TeamHandler) GetTeams(c *gin.Context){
	token := c.Param("token")

	claim , err := h.security.ParseToken(token)
	if err != nil{
		c.JSON(http.StatusUnauthorized , api_schema.ErrorResponse{
			Code: "UNAUTHORIZED",
			Error: "invalid tokan",
		})
		return
	}

	if claim.Role != "owner"{
		c.JSON(http.StatusForbidden,api_schema.ErrorResponse{
			Code: "FORBIDEN",
			Error: "you haven`t access",
		})
		return
	}

	togt , cgt := context.WithTimeout(c , time.Second)
	defer cgt()
	teams , err := h.quesries.GetTeams(togt , claim.UserID)

	if err != nil{
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError , api_schema.ErrorResponse{
			Code: "INTERNAL_SERVER_ERROR",
			Error: "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK , api_schema.GetTeamsResponse{
		Teams: teams,
	})

}



func (h *TeamHandler) DeleteTeam(c *gin.Context) {
    ctx := c.Request.Context()

    // -------- 1. Extract & validate JWT from path --------
    jwtToken := c.Param("token")
    if jwtToken == "" {
        c.JSON(http.StatusUnauthorized, api_schema.ErrorResponse{
            Code:  "UNAUTHORIZED",
            Error: "invalid token",
        })
        return
    }

    claims, err := h.security.ParseToken(jwtToken)
    if err != nil {
        c.JSON(http.StatusUnauthorized, api_schema.ErrorResponse{
            Code:  "UNAUTHORIZED",
            Error: "invalid token",
        })
        return
    }

	if claims.Role != "owner"{
		c.JSON(http.StatusForbidden , api_schema.ErrorResponse{
			Code:"FORBIDEN",
			Error: "you haven`t access",
		})
		return
	}


    // -------- 2. Bind and validate body --------
    var req api_schema.DeleteTeamRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, api_schema.ErrorResponse{
            Code:  "BAD_REQUEST",
            Error: "invalid request body",
        })
        return
    }

    req.Title = strings.TrimSpace(req.Title)
    if req.Title == "" {
        c.JSON(http.StatusBadRequest, api_schema.ErrorResponse{
            Code:  "BAD_REQUEST",
            Error: "title is required",
        })
        return
    }

    // -------- 3. Try to delete as owner --------
    team, err := h.quesries.DeleteTeam(ctx, sqlcpkg.DeleteTeamParams{
		Userid: claims.UserID,
		Title: req.Title,
	})

    if err == nil {
        // Success — return the deleted team
        c.JSON(http.StatusOK, api_schema.DeleteTeamResponse{Team: team})
        return
    }

    // Only pgx.ErrNoRows means "no rows matched" — anything else is a real failure.
    if !errors.Is(err, pgx.ErrNoRows) {
        var pgErr *pgconn.PgError
        if errors.As(err, &pgErr) {
            // optionally log pgErr
        }
        c.JSON(http.StatusInternalServerError, api_schema.ErrorResponse{
            Code:  "INTERNAL_SERVER_ERROR",
            Error: "internal server error",
        })
        return
    }

    // -------- 4. Nothing deleted — figure out 404 vs 403 --------
    exists, exErr := h.quesries.TeamExistsByTitle(ctx, req.Title)
    if exErr != nil {
        c.JSON(http.StatusInternalServerError, api_schema.ErrorResponse{
            Code:  "INTERNAL_SERVER_ERROR",
            Error: "internal server error",
        })
        return
    }

    if exists {
        // Team exists but belongs to someone else
        c.JSON(http.StatusForbidden, api_schema.ErrorResponse{
            Code:  "FORBIDEN",
            Error: "you haven`t access",
        })
        return
    }

    // No such team at all
    c.JSON(http.StatusNotFound, api_schema.ErrorResponse{
        Code:  "TEAM_NOT_FOUND",
        Error: "there is not any team with this title for you",
    })
}



func (h *TeamHandler) CreateUser(c *gin.Context) {
	// ---------------------------------------------------------
	// 1. read token from path param
	// ---------------------------------------------------------
	tokenStr := c.Param("token")
	if tokenStr == "" {
		c.JSON(http.StatusUnauthorized, api_schema.ErrorResponse{
			Code:  "UNAUTHORIZED",
			Error: "invalid token",
		})
		return
	}

	// ---------------------------------------------------------
	// 2. parse / validate JWT
	// ---------------------------------------------------------
	claims, err := h.security.ParseToken(tokenStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, api_schema.ErrorResponse{
			Code:  "UNAUTHORIZED",
			Error: "invalid token",
		})
		return
	}

	// ---------------------------------------------------------
	// 3. authorization: only admin/owner may add team users
	// ---------------------------------------------------------
	if claims.Role != "admin" && claims.Role != "owner" {
		c.JSON(http.StatusForbidden, api_schema.ErrorResponse{
			Code:  "FORBIDEN",
			Error: "you haven`t access",
		})
		return
	}

	// ---------------------------------------------------------
	// 4. bind request body
	// ---------------------------------------------------------
	var req api_schema.CreateTeamUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api_schema.ErrorResponse{
			Code:  "BAD_REQUEST",
			Error: "invalid request body",
		})
		return
	}
	if strings.TrimSpace(req.Username) == "" || strings.TrimSpace(req.Title) == "" {
		c.JSON(http.StatusBadRequest, api_schema.ErrorResponse{
			Code:  "BAD_REQUEST",
			Error: "username or title are required",
		})
		return
	}

	// ---------------------------------------------------------
	// 5. generate a random token for the new teamuser
	// ---------------------------------------------------------
	userToken := uuid.New().String()


	// ---------------------------------------------------------
	// 6. insert into teamuser
	// ---------------------------------------------------------
	user, err := h.quesries.CreateTeamUser(c.Request.Context(), sqlcpkg.CreateTeamUserParams{
		Username: req.Username,
		Title:    req.Title,   // <-- maps to $2 (team title)
		Token:    userToken,   // <-- maps to $3
	})
	if err != nil {
		// team title lookup returned no rows -> subquery is NULL
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, api_schema.ErrorResponse{
				Code:  "NOT_FOUND",
				Error: "team not found",
			})
			return
		}

		// unique violation (e.g. username already on this team)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			c.JSON(http.StatusConflict, api_schema.ErrorResponse{
				Code:  "CONFLICT",
				Error: "user already exists in team",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, api_schema.ErrorResponse{
			Code:  "INTERNAL_SERVER_ERROR",
			Error: "internal server error",
		})
		return
	}

	// ---------------------------------------------------------
	// 7. respond
	// ---------------------------------------------------------
	c.JSON(http.StatusOK, api_schema.CreateTeamUserResponse{
		User:  user,
		Token: userToken,
	})
}

