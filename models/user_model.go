package models

//import "go.mongodb.org/mongo-driver/bson/primitive"

type Repo struct {
	//Id                   primitive.ObjectID `json:"id,omitempty"`
	Name                 string             `json:"name,omitempty" validate:"required"`
	RampUp               float64            `json:"rampup,omitempty" validate:"required"`
	Correctness          float64            `json:"correctness,omitempty" validate:"required"`
	ResponsiveMaintainer float64            `json:"responsivemaintainer,omitempty" validate:"required"`
	BusFactor            float64            `json:"busfactor,omitempty" validate:"required"`
	ReviewCoverage       float64            `json:"reviewcoverage,omitempty" validate:"required"`
	DependancyPinning    float64            `json:"dependancypinning,omitempty" validate:"required"`
	License              float64            `json:"license,omitempty" validate:"required"`
	Net                  float64            `json:"net,omitempty" validate:"required"`
}
