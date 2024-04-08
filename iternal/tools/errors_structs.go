package tools

import "Edos_Docer/iternal/domain"

type BaseError struct {
	Code   int
	String string
	Err    error
	Result domain.User
}
