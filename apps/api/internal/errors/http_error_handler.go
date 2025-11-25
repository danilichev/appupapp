package errors

import (
	"errors"
	"fmt"
	"net/http"

	z "github.com/Oudwins/zog"
	zconst "github.com/Oudwins/zog/zconst"
	"github.com/labstack/echo/v4"

	"apps/api/internal/api"
)

func HTTPErrorHandler(err error, c echo.Context) {
	var httpErr *echo.HTTPError
	var valErr *ValidationError

	if errors.As(err, &valErr) {
		fmt.Println("Validation error:", valErr)
		c.JSON(
			http.StatusBadRequest,
			api.GeneralError{
				FieldErrors: validationDetails(*valErr.Issues),
				Message:     valErr.Error(),
			},
		)
		return
	}
	if errors.As(err, &httpErr) {
		fmt.Println("Http error:", httpErr.Message)
		c.JSON(
			httpErr.Code,
			api.GeneralError{
				Message: httpErr.Message.(string),
			},
		)
		return
	}

	c.JSON(http.StatusInternalServerError, api.GeneralError{
		Message: "Internal Server Error",
	})
}

func validationDetails(errs z.ZogIssueMap) *map[string]string {
	sanitizedErrs := z.Issues.SanitizeMap(errs)

	details := make(map[string]string)
	for k, v := range sanitizedErrs {
		if k != zconst.ISSUE_KEY_FIRST {
			details[k] = v[0]
		}
	}

	return &details
}
