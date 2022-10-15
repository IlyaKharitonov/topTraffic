package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"TopTraffic/internal/controllers/requestController"
	"TopTraffic/internal/services/requestService"
)

func main() {

	var (
		address = flag.String("a", "127.0.0.1","")
		port = flag.String("p", "8080","")
		advertiserAddressesStr = flag.String("d", "", "")
	)
	flag.Parse()

	advertiserAddresses, err := parseAdvertiserAddresses(*advertiserAddressesStr)
	if err != nil {
		log.Fatal(err)
	}

	as := requestService.NewAdvertiserService(advertiserAddresses, &http.Client{})
	rs := requestService.NewRequestService(as)
	rc := requestController.NewController(rs)
	requestController.HandlersRegister(rc)


	fmt.Printf("\nstarted server on %s:%s\n",*address, *port)
	err = http.ListenAndServe(fmt.Sprintf("%s:%s",*address, *port), nil )
	if err != nil {
		log.Fatal(err)
	}
}

func parseAdvertiserAddresses(advertiserAddressesStr string)([]string, error){
	if advertiserAddressesStr == ""{
		return nil, fmt.Errorf("advertiser addresses is empty, enter flag -d")
	}

	advertiserAddresses := strings.Split(advertiserAddressesStr, ",")

	if err := checkFormatAdvertiserAddress(advertiserAddresses); err != nil{
		return nil, err
	}

	return advertiserAddresses, nil
}

func checkFormatAdvertiserAddress(advertiserAddresses []string)error{
	reqexpAddr := `^([01]?\d\d?|2[0-4]\d|25[0-5])\.([01]?\d\d?|2[0-4]\d|25[0-5])\.([01]?\d\d?|2[0-4]\d|25[0-5])\.([01]?\d\d?|2[0-4]\d|25[0-5]):((6553[0-5])|(655[0-2][0-9])|(65[0-4][0-9]{2})|(6[0-4][0-9]{3})|([1-5][0-9]{4})|([0-5]{0,5})|([0-9]{1,4}))$`

	regexp, err := regexp.Compile(reqexpAddr)
	if err != nil{
		return err
	}

	wrongAddresses := make([]string, 0)

	for _, a := range advertiserAddresses{
		if !regexp.Match([]byte(a)){
			wrongAddresses = append(wrongAddresses, a)
		}
	}

	if len(wrongAddresses) != 0{
		return  fmt.Errorf("wrong addresses %s", wrongAddresses)
	}

	return nil
}