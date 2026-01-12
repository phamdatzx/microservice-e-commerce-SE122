package dto

type GHNItem struct {
	Name     string `json:"name"`
	Code     string `json:"code"`
	Quantity int    `json:"quantity"`
	Price    int    `json:"price"`
	Weight   int    `json:"weight"`
}

type GHNCreateOrderRequest struct {
	PaymentTypeID      int     `json:"payment_type_id"`
	RequiredNote       string  `json:"required_note"`
	FromName           string  `json:"from_name"`
	FromPhone          string  `json:"from_phone"`
	FromAddress        string  `json:"from_address"`
	FromWardName       string  `json:"from_ward_name"`
	FromDistrictName   string  `json:"from_district_name"`
	FromProvinceName   string  `json:"from_province_name"`
	ToName             string  `json:"to_name"`
	ToPhone            string  `json:"to_phone"`
	ToAddress          string  `json:"to_address"`
	ToWardCode         string  `json:"to_ward_code"`
	ToDistrictID       int     `json:"to_district_id"`
	CodAmount          int     `json:"cod_amount"`
	Content            string  `json:"content"`
	Weight             int     `json:"weight"`
	Length             int     `json:"length"`
	Width              int     `json:"width"`
	Height             int     `json:"height"`
	ServiceID          int     `json:"service_id"`
	ServiceTypeID      int     `json:"service_type_id"`
	Items              []GHNItem  `json:"items"`
}

func NewGHNCreateOrderRequest() GHNCreateOrderRequest {
	return GHNCreateOrderRequest{
		PaymentTypeID:    2,
		RequiredNote:     "CHOXEMHANGKHONGTHU",
		FromName:         "TinTest124",
		FromPhone:        "0987654321",
		FromAddress:      "72 Thành Thái, Phường 14, Quận 10, Hồ Chí Minh, Vietnam",
		FromWardName:     "Phường 14",
		FromDistrictName: "Quận 10",
		FromProvinceName: "HCM",
		ToName:           "TinTest124",
		ToPhone:          "0987654321",
		ToAddress:        "72 Thành Thái, Phường 14, Quận 10, Hồ Chí Minh, Vietnam",
		ToWardCode:       "20308",
		ToDistrictID:     1444,
		CodAmount:        0,
		Content:          "Theo New York Times",
		Weight:           200,
		Length:           1,
		Width:            19,
		Height:           10,
		ServiceID:        0,
		ServiceTypeID:    2,
		Items: []GHNItem{
			{
				Name:     "Áo Polo",
				Code:     "Polo123",
				Quantity: 1,
				Price:    200000,
				Weight:   1200,
			},
		},
	}
}
