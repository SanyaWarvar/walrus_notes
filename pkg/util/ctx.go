package util

import (
	"context"
	"fmt"
	"wn/pkg/constants"

	"github.com/google/uuid"
)

func GetUserRole(c context.Context) (string, error) {
	return getStringFromContext(c, constants.UserRoleCtx)
}

func GetUserId(c context.Context) (uuid.UUID, error) {
	id, err := GetUUIDFromContext(c, constants.UserIdCtx)
	if err != nil {
		return [16]byte{}, fmt.Errorf("GetUserId: %w", err)
	}
	return *id, nil
}

func CopyContextValues(parentCtx context.Context, keys ...any) context.Context {
	newCtx := context.Background()
	for _, key := range keys {
		if value := parentCtx.Value(key); value != nil {
			newCtx = context.WithValue(newCtx, key, value)
		}
	}
	return newCtx
}

func GetRequestId(c context.Context) (string, error) {
	id, err := getStringFromContext(c, constants.RequestIdCtx)
	if err != nil {
		return "", fmt.Errorf("getStringFromContext: %w", err)
	}
	return id, nil
}

func GetTrace(c context.Context) (string, error) {
	id, err := getStringFromContext(c, constants.TraceIdCtx)
	if err != nil {
		return "", fmt.Errorf("GetTrace: %w", err)
	}
	return id, nil
}

func GetSpan(c context.Context) (string, error) {
	id, err := getStringFromContext(c, constants.SpanIdCtx)
	if err != nil {
		return "", fmt.Errorf("GetSpan: %w", err)
	}
	return id, nil
}

func getStringFromContext(c context.Context, context string) (string, error) {
	val, ok := c.Value(context).(string)
	if !ok {
		return "", fmt.Errorf("cant convert %s to string", context)
	}
	return val, nil
}

func GetUUIDFromContext(c context.Context, context string) (*uuid.UUID, error) {
	idFromCtx, err := getStringFromContext(c, context)
	if err != nil {
		return nil, err
	}
	id, err := UUIDFromString(idFromCtx)
	if err != nil {
		return nil, fmt.Errorf("cant deserialise id %s", context)
	}
	return &id, nil
}
