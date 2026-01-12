package dto

type UserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Address  Address `json:"address"`
	Phone    string `json:"phone"`
	Image    string `json:"image"`
}


type Address struct {
	FullName    string    `json:"full_name" bson:"full_name"`
	Phone       string    `json:"phone" bson:"phone"`
	AddressLine string    `json:"address_line" bson:"address_line"`
	Ward        string    `json:"ward" bson:"ward"`
	District    string    `json:"district" bson:"district"`
	Province    string    `json:"province" bson:"province"`
	WardCode    string    `json:"ward_code" bson:"ward_code"`
	ProvinceCode string    `json:"province_code" bson:"province_code"`
	DistrictID  string    `json:"district_id" bson:"district_id"`
	ProvinceID  string    `json:"province_id" bson:"province_id"`
	Country     string    `json:"country" bson:"country"`
	Latitude    float64   `json:"latitude" bson:"latitude"`
	Longitude   float64   `json:"longitude" bson:"longitude"`
	Default     bool      `json:"default" bson:"is_default"`
}
