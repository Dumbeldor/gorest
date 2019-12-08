package gorest

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Error500 Trigger a standard error 500 with custom message in json format
func Error500(c echo.Context, err string) error {
	var er ErrorResponse
	er.Body.Message = err
	return c.JSON(http.StatusInternalServerError, er.Body)
}

// Error400 Trigger a standard error 400 with custom message in json format
func Error400(c echo.Context, err string) error {
	var er ErrorResponse
	er.Body.Message = err
	return c.JSON(http.StatusBadRequest, er.Body)
}

// Error4003Trigger a standard error 403 with custom message in json format
func Error403(c echo.Context, err string) error {
	var er ErrorResponse
	er.Body.Message = err
	return c.JSON(http.StatusForbidden, er.Body)
}
