package requestService

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"sync"
	"time"

	"TopTraffic/internal/entities"
	"TopTraffic/pkg"
)

type advertiserService struct {
	addresses []string
	cli *http.Client
}

func NewAdvertiserService(addresses []string, cli *http.Client)*advertiserService{
	return &advertiserService{addresses: addresses, cli: cli }
}

func (as *advertiserService)Request(ctx context.Context, request *entities.Request)(*entities.Response, error){
	var (
		impRespCh = make(chan []entities.Imp)
		respCh = make(chan *entities.Response, len(as.addresses))
		//errCh = make(chan error, len(as.addresses))
		response = &entities.Response{ID: request.ID, Imp: []entities.Imp{}}
		wg = &sync.WaitGroup{}
		ctxWithTimeout, cansel = context.WithTimeout(ctx, time.Millisecond*200)
	)
	defer cansel()

	reqBody, err := getReqBody(request)
	if err != nil {
		return nil, fmt.Errorf("(as *advertiserService)Request error:%s", err.Error())
	}

	go prepareResponse(request, respCh, impRespCh)

	for _, a := range as.addresses{
		wg.Add(1)
		go execBidRequest(ctxWithTimeout, as.cli, a, reqBody, respCh, wg)
	}
	//после исполнения всех горутин закрываем канал
	//чтобы спровоцировать выход из цикла в prepareResponse
	wg.Wait()
	close(respCh)

	response.Imp = <- impRespCh

	return response, nil
}

func getReqBody(request *entities.Request)([]byte, error){
	reqImps := make([]entities.BidImp, 0)
	for _,elem := range request.Tiles{
		var bidImp = entities.BidImp{
			ID:        elem.ID,
			MinWidth:  elem.Width,
			MinHeight: uint(math.Floor(float64(elem.Width) * elem.Ratio)),
		}
		reqImps = append(reqImps, bidImp)
	}

	reqBodyStruct := &entities.BidRequest{
		ID: request.ID,
		Imp: reqImps,
		Context: request.Context,
	}

	reqBody, err :=  json.Marshal(reqBodyStruct)
	if err != nil {
		return nil, err
	}

	return reqBody, nil
}

func execBidRequest(ctx context.Context, cli *http.Client, addr string, reqBody []byte, respCh chan *entities.Response, wg *sync.WaitGroup){
	var respBody = &entities.Response{}

	defer wg.Done()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, addr, bytes.NewBuffer(reqBody))
	if err != nil {
		log.Printf("requestService.execBidRequest() #1, request error:%s ,addr: %s", err.Error(), addr)
		return
	}

	resp,err := cli.Do(req)
	if err != nil {
		log.Printf("requestService.execBidRequest() #2, request error:%s ,addr: %s", err.Error(), addr)
		//errCh <- fmt.Errorf("requestService.execBidRequest(), request error:%s ,addr: %s", err.Error(), addr)
		return
	}

	if resp.StatusCode == http.StatusNoContent{
		return
	}

	err = pkg.Unmarshal(respBody, resp.Body)
	//полученные данные отправляем в prepareResponse
	respCh <- respBody

}

func prepareResponse(request *entities.Request, respCh chan *entities.Response, impRespCh chan []entities.Imp){
	//тк не известко как упорядочен слайс tile в реквесте,
	//копируем порядок для imp слайcа в респонзе.
	impResponse := make([]entities.Imp, 0)
	for _,tile := range request.Tiles{
		imp := entities.Imp{ID:tile.ID}
		impResponse = append(impResponse, imp)
	}

	//читаем ответы рекламодателей и обрабатываем их
	for r := range respCh{
		writeInImpResponse(impResponse, r)
	}

	impRespCh <- impResponse

}

// writeInImpResponse обрабатывает r - то что пришло от одного рекламодателя
func writeInImpResponse(impResponse []entities.Imp, r *entities.Response){
	for key := range r.Imp{
		for i := range impResponse{
			if impResponse[i].ID == r.Imp[key].ID && impResponse[i].Price > r.Imp[key].Price{
				//перезаписываем каждое поле отдельно, чтобы исключить Price в ответе
				impResponse[i].Width = r.Imp[key].Width
				impResponse[i].Height = r.Imp[key].Height
				impResponse[i].Title = r.Imp[key].Title
				impResponse[i].URL = r.Imp[key].URL
				break
			}
		}
	}
}

