package interceptors

import (
	"context"

	"github.com/kotche/url-shortening-service/internal/app/config"
	"github.com/kotche/url-shortening-service/internal/app/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// UnaryCookieInterceptor checks for the presence of the user ID in the cookie file. If not, then a new one is issued
func UnaryCookieInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	userID := utils.GetUserIDFromMD(ctx)
	if userID != "" {
		return handler(ctx, req)
	}
	_, userIDCookie := utils.MakeUserIDCookie()
	md := metadata.New(map[string]string{string(config.UserIDCookieName): userIDCookie})
	newCtx := metadata.NewIncomingContext(ctx, md)
	return handler(newCtx, req)
}
