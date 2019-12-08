package gorest

import (
	"net/http"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

// MessageResponse use to send a json response
type MessageResponse struct {
	// Message
	// required: true
	Message string `json:"message, required"`
}

// ErrorResponse use to send a json error response
type ErrorResponse struct {
	// in: body
	Body struct {
		// Error message
		// required: true
		Message string `json:"message,required"`
	}
}

func Write200(c echo.Context, message string) error {
	return writeJSON(c, http.StatusOK, message)
}

func Write201(c echo.Context, message string) error {
	return writeJSON(c, http.StatusCreated, message)
}

func writeJSON(c echo.Context, httpStatus int, message string) error {
	var mr MessageResponse
	mr.Message = message
	return c.JSON(httpStatus, mr)
}

// WriteJSONError use to write a json error
func WriteJSONError(c echo.Context, httpStatus int, userMessage interface{}, errorMessage string) error {
	if len(errorMessage) == 0 {
		log.Errorf("%s - error %d: %s", c.Path(), httpStatus, userMessage)
	} else {
		log.Errorf("%s - error %d: %s", c.Path(), httpStatus, errorMessage)
	}
	return c.JSON(httpStatus, userMessage)
}
