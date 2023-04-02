package models

//import "go.mongodb.org/mongo-driver/bson/primitive"

type Repo struct {
	//Id                   primitive.ObjectID `json:"id,omitempty"`
	Name                 string             `json:"name,omitempty" validate:""`
	RampUp               float64            `json:"rampup" validate:""`
	Correctness          float64            `json:"correctness" validate:""`
	ResponsiveMaintainer float64            `json:"responsivemaintainer" validate:""`
	BusFactor            float64            `json:"busfactor" validate:""`
	ReviewCoverage       float64            `json:"reviewcoverage" validate:""`
	DependancyPinning    float64            `json:"dependancypinning" validate:""`
	License              int            	`json:"license" validate:""`
	Net                  float64            `json:"net" validate:""`
}
