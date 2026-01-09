package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)


type GHNClient struct {
	baseURL    string
	token      string
}

type GHNAddressCode struct {
	ProvinceID int
	DistrictID int
	WardCode   string
}

type Province struct {
	ProvinceID   int    `json:"ProvinceID"`
	ProvinceName string `json:"ProvinceName"`
	Code         string `json:"Code"`
}

type District struct {
	DistrictID   int    `json:"DistrictID"`
	ProvinceID   int    `json:"ProvinceID"`
	DistrictName string `json:"DistrictName"`
	Code         string `json:"Code"`
}

type Ward struct {
	WardCode   string `json:"WardCode"`
	DistrictID int    `json:"DistrictID"`
	WardName   string `json:"WardName"`
}

type GHNResponse[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    []T    `json:"data"`
}


func NewGHNClient() *GHNClient {
	baseURL := os.Getenv("GHN_URL")
	token := os.Getenv("GHN_TOKEN")

	if baseURL == "" || token == "" {
		panic("GHN_URL and GHN_TOKEN must be set")
	}

	return &GHNClient{
		baseURL: baseURL,
		token: token,
	}
}

func callGHN[T any](url string, token string) ([]T, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Token", token)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res GHNResponse[T]
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	if res.Code != 200 {
		return nil, fmt.Errorf("GHN error: %s", res.Message)
	}

	return res.Data, nil
}

func (c *GHNClient)ResolveGHNAddress(
	provinceName string,
	districtName string,
	wardName string,
) (*GHNAddressCode, error) {

	// 1️⃣ Province
	provinces, err := callGHN[Province](
		fmt.Sprintf("%s/shiip/public-api/master-data/province", c.baseURL),
		c.token,
	)
	if err != nil {
		return nil, err
	}

	var provinceID int
	for _, p := range provinces {
		if strings.EqualFold(p.ProvinceName, provinceName) {
			provinceID = p.ProvinceID
			break
		}
	}
	if provinceID == 0 {
		return nil, fmt.Errorf("province not found")
	}

	// 2️⃣ District
	districtURL := fmt.Sprintf(
		"%s/shiip/public-api/master-data/district?province_id=%d",
		c.baseURL,
		provinceID,
	)

	districts, err := callGHN[District](districtURL, c.token)
	if err != nil {
		return nil, err
	}

	var districtID int
	for _, d := range districts {
		if strings.EqualFold(d.DistrictName, districtName) {
			districtID = d.DistrictID
			break
		}
	}
	if districtID == 0 {
		return nil, fmt.Errorf("district not found")
	}

	// 3️⃣ Ward
	wardURL := fmt.Sprintf(
		"%s/shiip/public-api/master-data/ward?district_id=%d",
		c.baseURL,
		districtID,
	)

	wards, err := callGHN[Ward](wardURL, c.token)
	if err != nil {
		return nil, err
	}

	var wardCode string
	for _, w := range wards {
		if strings.EqualFold(w.WardName, wardName) {
			wardCode = w.WardCode
			break
		}
	}
	if wardCode == "" {
		return nil, fmt.Errorf("ward not found")
	}

	return &GHNAddressCode{
		ProvinceID: provinceID,
		DistrictID: districtID,
		WardCode:   wardCode,
	}, nil
}