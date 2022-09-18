package grpc

import (
	"context"
	"errors"
	"net/http"

	"github.com/kotche/url-shortening-service/internal/app/config"
	"github.com/kotche/url-shortening-service/internal/app/model"
	"github.com/kotche/url-shortening-service/internal/app/service"
	pb "github.com/kotche/url-shortening-service/internal/app/transport/grpc/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ICookieManager retrieves the user id from cookies
type ICookieManager interface {
	GetUserID(ctx context.Context) string
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
		Cm:      model.CookieManagerMD{},
	}
	return handler
}

// Ping check DB connection
func (h *Handler) Ping(ctx context.Context, r *pb.EmptyRequest) (*pb.PingResponse, error) {
	if err := h.Service.Ping(ctx); err != nil {
		return nil, status.Errorf(codes.Internal, "ping error: %s", err.Error())
	}
	response := pb.PingResponse{Status: int32(http.StatusOK)}
	return &response, nil
}

// HandlePost accepts the URL string in the request body and returns its shorten version.
func (h *Handler) HandlePost(ctx context.Context, r *pb.HandlePostRequest) (*pb.HandlePostResponse, error) {
	originURL := r.OriginURL
	userID := h.Cm.GetUserID(ctx)
	urlModel, err := h.Service.GetURLModel(userID, originURL)

	if errors.As(err, &model.ConflictURLError{}) {
		e := err.(model.ConflictURLError)
		return nil, status.Errorf(codes.AlreadyExists, "handlePost error: %s", e.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "handlePost error: %s", err.Error())
	}

	response := pb.HandlePostResponse{
		Status:     int32(http.StatusCreated),
		ShortenURL: h.Conf.BaseURL + "/" + urlModel.Short,
	}

	return &response, nil
}

// HandleGet gets the original URL from a shortened link
func (h *Handler) HandleGet(ctx context.Context, r *pb.HandleGetRequest) (*pb.HandleGetResponse, error) {
	shortURL := r.ShortURL
	url, err := h.Service.GetURLModelByID(shortURL)

	if errors.As(err, &model.GoneError{}) {
		return nil, status.Errorf(codes.NotFound, "handleGet error: %s", err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "handleGet error: %s", err.Error())
	}

	response := pb.HandleGetResponse{
		Status:    int32(http.StatusTemporaryRedirect),
		OriginURL: url.Origin,
	}

	return &response, err
}

//HandleGetUserURLs, HandlePostShortenBatch, HandleDeleteURLs, HandleGetStats
