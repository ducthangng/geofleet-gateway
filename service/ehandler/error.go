package ehandler

import "errors"

var (
	ERROR_MISSING_USER = errors.New("service does not receive any params related to User")
	ERROR_USER_SERVICE = errors.New("user service temporary close, please try again")
)
