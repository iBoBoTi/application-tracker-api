package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iBoBoTi/ats/internal/validator"
)

type Data struct {
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp string      `json:"timestamp,omitempty"`
	Errors    interface{} `json:"errors,omitempty"`
	Status    string      `json:"status,omitempty"`
}

func SuccessJSONResponse(c *gin.Context, status int, message string, data interface{}) {

	responseData := Data{
		Message:   message,
		Data:      data,
		Status:    http.StatusText(status),
		Timestamp: time.Now().Format(time.RFC850),
	}

	c.JSON(status, responseData)
}

func ErrorJSONResponse(c *gin.Context, status int, errs ...error) {

	responseData := Data{
		Message:   "error processing request",
		Status:    http.StatusText(status),
		Timestamp: time.Now().Format(time.RFC850),
	}
	if len(errs) == 1 {
		responseData.Errors = errs[0].Error()
	} else {
		outputErrors := make([]string, 0, len(errs))
		for _, err := range errs {
			outputErrors = append(outputErrors, err.Error())
		}
		responseData.Errors = outputErrors
	}

	c.JSON(status, responseData)
}

func SendValidationError(ctx *gin.Context, errors *validator.ValidationError) {
	responseData := Data{
		Errors:    errors.Fields,
		Message:   errors.Message,
		Status:    http.StatusText(http.StatusUnprocessableEntity),
		Timestamp: time.Now().Format(time.RFC850),
	}

	ctx.JSON(http.StatusUnprocessableEntity, responseData)

}
