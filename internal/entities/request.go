package entities

type Request struct {
	ID          string    `json:"id,omitempty" valid:"required"`
	Tiles  		[]Tile    `json:"tiles,omitempty" valid:"required"`
	Context 	Context   `json:"context,omitempty" valid:"required"`
}

type Tile struct {
	ID    uint 		`json:"id,omitempty" valid:"required"`
	Width uint 		`json:"width,omitempty" valid:"required"`
	Ratio float64 	`json:"ratio,omitempty" valid:"required"`
}

//