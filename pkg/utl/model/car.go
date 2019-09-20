package rest

// Car represents company model
type Car struct {
	Base
	Name         string `json:"name"`
	Merk         string `json:"merk"`
	Model        string `json:"model"`
	PoliceNumber string `json:"police_number"`
	Colour       string `json:"colour"`
	CarType      string `json:"car_type"`
	Owner        User   `json:"owner"`
}
