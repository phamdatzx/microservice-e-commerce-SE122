package dto

// GHNCalculateFeeRequest represents the request to calculate shipping fee
type GHNCalculateFeeRequest struct {
	FromDistrictID int    `json:"from_district_id"`
	FromWardCode   string `json:"from_ward_code"`
	ServiceID      int    `json:"service_id"`
	ServiceTypeID  int    `json:"service_type_id"`
	ToDistrictID   int    `json:"to_district_id"`
	ToWardCode     string `json:"to_ward_code"`
	Height         int    `json:"height"`
	Length         int    `json:"length"`
	Weight         int    `json:"weight"`
	Width          int    `json:"width"`
}

// GHNCalculateFeeResponse represents the response from GHN fee calculation API
type GHNCalculateFeeResponse struct {
	Total               int `json:"total"`
	ServiceFee          int `json:"service_fee"`
	InsuranceFee        int `json:"insurance_fee"`
	PickStationFee      int `json:"pick_station_fee"`
	R2SFee              int `json:"r2s_fee"`
	ReturnAgain         int `json:"return_again"`
	DocumentReturn      int `json:"document_return"`
	DoubleCheck         int `json:"double_check"`
	CODFailedFee        int `json:"cod_failed_fee"`
	PickRemoteAreasFee  int `json:"pick_remote_areas_fee"`
	DeliverRemoteAreasFee int `json:"deliver_remote_areas_fee"`
	CODFee              int `json:"cod_fee"`
}

func NewGHNCalculateFeeRequest() GHNCalculateFeeRequest {
	return GHNCalculateFeeRequest{
		ServiceID:      53322,
		ServiceTypeID:  2,
		Height:         50,
		Length:         20,
		Weight:         20,
		Width:          20,
	}
}