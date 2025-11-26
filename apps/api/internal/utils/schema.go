package utils

import (
	z "github.com/Oudwins/zog"
)

var EmailSchema = z.String().
	Email(z.Message("Invalid email format")).
	Required(z.Message("Email is required"))
