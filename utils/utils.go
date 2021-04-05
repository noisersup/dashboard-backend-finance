package utils

import (
	"errors"
)

func Err(msg string, err error) error{
	return errors.New("["+msg+"] "+err.Error())
}