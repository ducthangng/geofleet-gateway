package gwerr

import (
	"errors"
)

var (
	// --- Authentication & Authorization (401, 403) ---
	ErrUnauthorized          = errors.New("unauthorized: missing or invalid token")
	ErrForbidden             = errors.New("forbidden: you don't have permission to access this resource")
	ErrInvalidAPIKey         = errors.New("invalid api key")
	ErrInvalidInput          = errors.New("invalid api input")
	ErrInvalidAuthentication = errors.New("invalid authentication key")

	// --- Request Validation (400) ---
	ErrBadRequest      = errors.New("bad request: invalid input format")
	ErrPayloadTooLarge = errors.New("request payload too large")

	// --- Routing & Service Discovery (404, 503) ---
	ErrServiceNotFound     = errors.New("requested service is currently unavailable")
	ErrNoAvailableInstance = errors.New("no healthy instances available to handle request")
	ErrInvalidRoute        = errors.New("the requested path does not exist on gateway")

	// --- Flow Control (429, 504) ---
	ErrRateLimitExceeded = errors.New("too many requests: please slow down")
	ErrRequestTimeout    = errors.New("gateway timeout: service took too long to respond")

	// --- Connection & Downstream (500, 502) ---
	ErrInternalGateway    = errors.New("internal gateway error")
	ErrBadGateway         = errors.New("bad gateway: invalid response from upstream service")
	ErrServiceUnavailable = errors.New("service is temporarily down for maintenance")
)
