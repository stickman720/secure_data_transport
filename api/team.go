package api

import (
	"net/http"
	"secure_data_transport/api_schema"
	"secure_data_transport/config"
	"secure_data_transport/core"
	sqlcpkg "secure_data_transport/sqlcpkg"
	"secure_data_transport/widget"
	"errors"
	"strings"
	"github.com/gin-gonic/gin"
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
		req.Team.Transfear <= 0 {
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
		Treansfear:       int32(req.Team.Transfear),
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