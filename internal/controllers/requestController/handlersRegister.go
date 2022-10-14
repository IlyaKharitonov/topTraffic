package requestController

import (
	"net/http"
)

func HandlersRegister(c *controller) {
	http.HandleFunc("/placements/request", c.Request)
}
