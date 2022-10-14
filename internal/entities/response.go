package entities

type Response struct {
	ID 		string `json:"id"`
	Imp 	[]Imp `json:"imp"`
}


type Imp struct {
	ID 		uint 	`json:"id,omitempty"`
	Width 	uint 	`json:"width,omitempty"`
	Height 	uint 	`json:"height,omitempty"`
	Title 	string 	`json:"title,omitempty"`
	URL 	string  `json:"url,omitempty"`
	Price 	float64 `json:"price,omitempty"`
}
