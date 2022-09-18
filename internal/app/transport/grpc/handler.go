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
func (h *Handler) Ping(ctx context.Context, r *pb.PingRequest) (*pb.PingResponse, error) {
	if err := h.Service.Ping(ctx); err != nil {
		return nil, status.Errorf(codes.Internal, "ping error: %s", err.Error())
	}
	response := pb.PingResponse{Status: int32(http.StatusOK)}
	return &response, nil
}

// HandlePost accepts the URL string in the request body and returns its shorten version.
func (h *Handler) HandlePost(ctx context.Context, r *pb.HandlePostRequest) (*pb.HandlePostResponse, error) {
	userID := h.Cm.GetUserID(ctx)
	if userID == "" {
		return nil, status.Errorf(codes.Internal, "HandlePost error: %s", "user ID is empty")
	}

	originURL := r.OriginURL
	urlModel, err := h.Service.GetURLModel(ctx, userID, originURL)

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
	url, err := h.Service.GetURLModelByID(ctx, shortURL)

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

// HandleGetUserURLs gets all shortened links by the user
func (h *Handler) HandleGetUserURLs(ctx context.Context, r *pb.HandleGetUserURLsRequest) (*pb.HandleGetUserURLsResponse, error) {
	userID := h.Cm.GetUserID(ctx)
	if userID == "" {
		return nil, status.Errorf(codes.Internal, "HandleGetUserURLs error: %s", "user ID is empty")
	}

	userUrls, err := h.Service.GetUserURLs(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "HandleGetUserURLs error: %s", err.Error())
	}

	if len(userUrls) == 0 {
		return nil, status.Errorf(codes.NotFound, "HandleGetUserURLs error: %s", "no shortened URLs")
	}

	response := pb.HandleGetUserURLsResponse{}
	for _, v := range userUrls {
		responseURL := pb.SetURLsResponse{
			ShortURL:  h.Conf.BaseURL + "/" + v.Short,
			OriginURL: v.Origin,
		}
		response.SetURLs = append(response.SetURLs, &responseURL)
	}

	return &response, nil
}

// HandlePostShortenBatch accepts in the request body a set of URLs for shorten
func (h *Handler) HandlePostShortenBatch(ctx context.Context, r *pb.HandlePostShortenBatchRequest) (*pb.HandlePostShortenBatchResponse, error) {
	userID := h.Cm.GetUserID(ctx)
	if userID == "" {
		return nil, status.Errorf(codes.Internal, "HandlePostShortenBatch error: %s", "user ID is empty")
	}

	if len(r.CorrelationURL) == 0 {
		return nil, status.Errorf(codes.Internal, "HandlePostShortenBatch error: %s", "correlation URL empty")
	}

	inputDataList := make([]model.InputCorrelationURL, 0, len(r.CorrelationURL))

	for _, corURLInput := range r.CorrelationURL {
		inputDataList = append(inputDataList, model.InputCorrelationURL{
			CorrelationID: corURLInput.Id,
			Origin:        corURLInput.OriginalURL,
		})
	}

	outputDataList, err := h.Service.ShortenBatch(ctx, userID, inputDataList)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "HandlePostShortenBatch error: %s", err.Error())
	}

	response := pb.HandlePostShortenBatchResponse{}
	for ind := range outputDataList {
		response.CorrelationURL = append(response.CorrelationURL, &pb.CorrelationURLResponse{
			Id:       outputDataList[ind].CorrelationID,
			ShortURL: h.Conf.BaseURL + "/" + outputDataList[ind].Short,
		})
	}
	response.Status = int32(http.StatusCreated)

	return &response, nil
}

// HandleDeleteURLs accepts a list of shortened URL to delete
func (h *Handler) HandleDeleteURLs(ctx context.Context, r *pb.HandleDeleteURLsRequest) (*pb.HandleDeleteURLsResponse, error) {
	userID := h.Cm.GetUserID(ctx)
	if userID == "" {
		return nil, status.Errorf(codes.Internal, "HandleDeleteURLs error: %s", "user ID is empty")
	}

	toDelete := make([]string, 0, len(r.DeleteURLs))
	for _, delURL := range r.DeleteURLs {
		toDelete = append(toDelete, delURL)
	}

	go func() {
		h.Service.DeleteURLs(userID, toDelete)
	}()

	response := pb.HandleDeleteURLsResponse{Status: int32(http.StatusAccepted)}
	return &response, nil
}

// HandleGetStats returns the number of shortened urls and the number of users in the service
func (h *Handler) HandleGetStats(ctx context.Context, r *pb.HandleGetStatsRequest) (*pb.HandleGetStatsResponse, error) {
	stats, err := h.Service.GetStats(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "HandleGetStats error: %s", err.Error())
	}

	response := pb.HandleGetStatsResponse{Urls: int64(stats.NumberOfURLs), Users: int64(stats.NumberOfUsers)}
	return &response, nil
}
