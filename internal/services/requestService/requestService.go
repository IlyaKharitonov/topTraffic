package requestService

import (
	"context"
	"fmt"

	"TopTraffic/internal/entities"
)

type advertiser interface {
	Request(ctx context.Context, request *entities.Request) (*entities.Response,error)
}

type requestService struct {
	advertiser advertiser
}

func NewRequestService(advertiser advertiser) *requestService {
	return &requestService{advertiser: advertiser}
}

func (rs *requestService)Request(ctx context.Context, request *entities.Request)(*entities.Response, error){
	resp, err := rs.advertiser.Request(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("(rs *requestService)Request #1\n Error: %s\n", err.Error())
	}

	return resp, nil
}


