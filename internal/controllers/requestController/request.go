package requestController

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/asaskevich/govalidator"

	"TopTraffic/internal/entities"
	"TopTraffic/pkg"
)


type RequestService interface {
	Request(ctx context.Context, request *entities.Request) (*entities.Response,error)
}

type controller struct {
	rs RequestService
}

func NewController(rs RequestService) *controller {
	return &controller{rs: rs}
}

func (c *controller) Request(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var request = &entities.Request{}

	err := pkg.Unmarshal(request, req.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("(c *controller) Request#1 \n Error: %s\n", err.Error())
		return
	}

	res, err := govalidator.ValidateStruct(request)
	if !res{
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("(c *controller) Request#2 \n Error: %s\n", err.Error())
		return
	}
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("(c *controller) Request#3 \n Error: %s\n", err.Error())
		return
	}

	resp, err := c.rs.Request(context.Background(), request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("(c *controller) Request#4 \n Error: %s\n", err.Error())
		return
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("(c *controller) Request#5 \n Error: %s\n", err.Error())
		return
	}

}


