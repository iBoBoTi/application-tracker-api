package server

// import (
// 	"github.com/gin-gonic/gin"

// )

// const (
// 	ContextUser           = "context_user"
// 	ContextToken          = "context_token"
// 	ContextSystemSettings = "context_system_settings"
// )

// // ContextSetUser sets the user in the context
// func (srv *Server) ContextSetUser(c *gin.Context, user *db.User) *gin.Context {
// 	c.Set(ContextUser, user)
// 	return c
// }

// // ContextSetToken sets the user in the context
// func (srv *Server) ContextSetToken(c *gin.Context, payload *security.Payload) *gin.Context {
// 	c.Set(ContextToken, payload)
// 	return c
// }

// // ContextGetUser gets the user from the context
// func (srv *Server) ContextGetUser(c *gin.Context) *db.User {
// 	user, ok := c.Get(ContextUser)
// 	if !ok {
// 		panic("missing user value in context")
// 	}
// 	return user.(*db.User)
// }

// // ContextGetToken gets the user from the context
// func (srv *Server) ContextGetToken(c *gin.Context) *security.Payload {
// 	token, ok := c.Get(ContextToken)
// 	if !ok {
// 		panic("missing token value in context")
// 	}
// 	return token.(*security.Payload)
// }
