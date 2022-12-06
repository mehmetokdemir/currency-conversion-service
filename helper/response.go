package helper

import (
	// Go imports
	"net/http"
	"strings"

	// External imports
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Success    bool                 `json:"success" extensions:"x-order=1" example:"true"`
	StatusCode int                  `json:"status_code" extensions:"x-order=2" example:"200"`
	Warnings   ResponseWarningArray `json:"warnings,omitempty" extensions:"x-order=3"`
	Error      *ResponseError       `json:"error,omitempty" extensions:"x-order=4"`
	Data       interface{}          `json:"data,omitempty" extensions:"x-order=5"`
}

type ResponseError struct {
	Message string `json:"message" extensions:"x-order=1" example:"NOT_FOUND"`
	Detail  string `json:"detail" extensions:"x-order=2" example:"user not found"`
}

type ResponseWarning struct {
	Field string `json:"field" extensions:"x-order=1" example:"password"`
	Code  string `json:"code" extensions:"x-order=2" example:"invalid"`
}

// Success response
func Success(ctx *gin.Context, data interface{}) {
	res := new(Response)
	res.Success = true
	res.StatusCode = http.StatusOK
	res.Warnings = nil
	res.Error = nil
	res.Data = data
	ctx.JSON(http.StatusOK, res)
	ctx.Next()
}

// Error response
func Error(ctx *gin.Context, statusCode int, message, detail string) {
	res := new(Response)
	res.Success = false
	res.StatusCode = statusCode
	res.Warnings = nil
	res.Error = &ResponseError{
		Message: message,
		Detail:  detail,
	}
	res.Data = nil
	ctx.JSON(statusCode, res)
	ctx.Next()
}

type ResponseWarningArray []ResponseWarning

func (w ResponseWarningArray) Add(field, code string) ResponseWarningArray {
	return append(w, ResponseWarning{
		Field: field,
		Code:  code,
	})
}

func WarningsFromValidationError(err error) (warnings ResponseWarningArray) {
	const errorSeparator = "|"
	switch errs := err.(type) {
	case govalidator.Errors:
		for _, e := range errs {
			parts := strings.Split(e.Error(), errorSeparator)
			if len(parts) != 2 {
				continue
			}
			warnings = append(warnings, ResponseWarning{
				Field: parts[0],
				Code:  parts[1],
			})
		}
	}
	return
}

var WarningMap map[string]string

// Warning returns required empty values
func Warning(ctx *gin.Context, warnings ResponseWarningArray) {
	res := new(Response)
	res.Success = false
	res.StatusCode = http.StatusBadRequest
	for _, warning := range warnings {
		if v, ok := WarningMap[warning.Code]; ok && len(v) > 0 {
			warning.Code = v
		}
		res.Warnings = append(res.Warnings, warning)
	}
	res.Error = nil
	res.Data = nil
	ctx.JSON(http.StatusBadRequest, res)
	ctx.Next()
}
