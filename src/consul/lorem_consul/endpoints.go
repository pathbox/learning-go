package lorem_consul

import (
	"context"
	"errors"
	"strings"

	"github.com/go-kit/kit/endpoint"
)

var (
	ErrRequestTypeNotFound = errors.New("Request type only valid for word, sentence and paragraph")
)

//Lorem Request
type LoremRequest struct {
	RequestType string `json:"requestType"`
	Min         int    `json:"min"`
	Max         int    `json:"max"`
}

//Lorem Response
type LoremResponse struct {
	Message string `json:"message"`
	Err     error  `json:"err,omitempty"`
}

//Health Request
type HealthRequest struct {
}

//Health Response
type HealthResponse struct {
	Status bool `json:"status"`
}

// endpoints wrapper
type Endpoints struct {
	LoremEndpoint  endpoint.Endpoint
	HealthEndpoint endpoint.Endpoint
}

// creating Lorem Ipsum Endpoint
func MakeLoremLoggingEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(LoremRequest)

		var (
			txt      string
			min, max int
		)

		min = req.Min
		max = req.Max

		if strings.EqualFold(req.RequestType, "Word") {
			txt = svc.Word(min, max)
		} else if strings.EqualFold(req.RequestType, "Sentence") {
			txt = svc.Sentence(min, max)
		} else if strings.EqualFold(req.RequestType, "Paragraph") {
			txt = svc.Paragraph(min, max)
		} else {
			return nil, ErrRequestTypeNotFound
		}
		return LoremResponse{Message: txt}, nil
	}

}

// creating health endpoint
func MakeHealthEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		status := svc.HealthCheck()
		return HealthResponse{Status: status}, nil
	}
}
