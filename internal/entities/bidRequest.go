package entities

type BidRequest struct {
	ID 		string 		`json:"id"`
	Imp 	[]BidImp 	`json:"imp"`
	Context Context 	`json:"context"`
}

type BidImp struct {
	ID 			uint 	`json:"id"`
	MinWidth 	uint 	`json:"min_width"`
	MinHeight 	uint 	`json:"min_height"`
}