package engine

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthErrorCode int32

const (
	AuthErrorCodeMissingAuthClaims AuthErrorCode = 2001
	AuthErrorCodeInvalidClaims     AuthErrorCode = 2002
)

var (
	ErrMissingAuthClaims = status.Error(codes.Code(AuthErrorCodeMissingAuthClaims), "context missing authz claims")
	ErrInvalidClaims     = status.Error(codes.Code(AuthErrorCodeInvalidClaims), "invalid claims")
)
