package dto

type AddressModel struct {
	StreetName string ` json:"streetName"`
	LandMark   string ` json:"landMark"`
	ZipCode    string ` json:"zipCode"`
	City       string ` json:"city"`
	State      string ` json:"state"`
}
