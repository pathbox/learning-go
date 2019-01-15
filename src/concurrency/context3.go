package main

import (
	"context"
	"fmt"
)

func main() {
	ProcessRequest("jane", "abc123")
}

type ctxKey int

const (
	ctxUserID ctxKey = iota
	ctxAuthToken
)

func UserID(c context.Context) string {
	return c.Value(ctxUserID).(string)
}
func AuthToken(c context.Context) string {
	return c.Value(ctxAuthToken).(string)
}

func ProcessRequest(userID, authToken string) {
	ctx := context.WithValue(context.Background(), ctxUserID, userID)
	ctx = context.WithValue(ctx, ctxAuthToken, authToken)
	HandleResponse(ctx)
}

func HandleResponse(ctx context.Context) {
	fmt.Printf(
		"handling response for %v (auth: %v)",
		UserID(ctx),
		AuthToken(ctx),
	)
}
