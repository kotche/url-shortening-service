package grpc

import (
	"context"
	"net/http"

	"github.com/kotche/url-shortening-service/internal/app/config"
	"github.com/kotche/url-shortening-service/internal/app/model"
	"github.com/kotche/url-shortening-service/internal/app/service"
	pb "github.com/kotche/url-shortening-service/internal/app/transport/grpc/proto"
)

// ICookieManager retrieves the user id from cookies
type ICookieManager interface {
	GetUserID(r *http.Request) string
}

type Handler struct {
	Service *service.Service
	pb.UnimplementedShortenerServer
	Conf *config.Config
	Cm   ICookieManager
}

func NewHandler(service *service.Service, conf *config.Config) *Handler {
	handler := &Handler{
		Service: service,
		Conf:    conf,
		Cm:      model.CookieManager{},
	}
	return handler
}

// Ping implement method for gRPC for check DB connection
func (h *Handler) Ping(ctx context.Context, _ *pb.EmptyRequest) (*pb.PingResponse, error) {
	var response pb.PingResponse
	var err error
	status := http.StatusOK
	if err = h.Service.Ping(ctx); err != nil {
		status = http.StatusInternalServerError
	}
	response.Code = int32(status)

	return &response, err
}
