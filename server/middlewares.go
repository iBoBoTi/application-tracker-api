package server

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/iBoBoTi/ats/internal/log"
	"github.com/iBoBoTi/ats/internal/models"
)

const (
	AuthorizationHeaderKey  = "Authorization"
	AuthorizationTypeBearer = "Bearer"
)

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Pin")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// CustomLogger is a middleware that logs the request
func CustomLogger(l log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		path := c.Request.URL.Path

		l.Info("[GIN]", map[string]interface{}{
			"start":       start.Format("2006/01/02 - 15:04:05"),
			"latency":     latency.String(),
			"client_ip":   clientIP,
			"method":      method,
			"status_code": fmt.Sprintf("%d", statusCode),
			"path":        path,
			"end":         end.Format("2006/01/02 - 15:04:05"),
		})

	}
}

// ApplyAuthentication is a middleware that checks for the authorization header
// the function does not check if the user is activated or not
func (srv *Server) ApplyAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Vary", AuthorizationHeaderKey)

		authHeader := c.GetHeader(AuthorizationHeaderKey)
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header",
			})
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != AuthorizationTypeBearer {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header",
			})
			return
		}

		token := headerParts[1]
		payload, err := srv.TokenMaker.VerifyToken(token)

		if err != nil {
			if strings.Contains(err.Error(), "expired") {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "token has expired",
				})
				return
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization token",
			})
			return
		}

		user, err := findUserByID(srv.DB, payload.UserID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization token",
			})
			return
		}

		srv.ContextSetUser(c, user)
		srv.ContextSetToken(c, payload)

		c.Next()
	}
}

func findUserByID(db *models.Database, id uuid.UUID) (*models.User, error) {
	user := &models.User{}
	if err := db.GormDB.Model(&models.User{}).Where("id", id).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
