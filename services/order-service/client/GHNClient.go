package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"order-service/dto"
	"order-service/model"
	"os"
	"strconv"
	"strings"
	"time"
)

type GHNClient struct {
	baseURL string
	token   string
}

type NoData struct {
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

type CreateOrderData struct {
	OrderCode string `json:"order_code"`
	SortCode  string `json:"sort_code"`
}

type CreateOrderResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    CreateOrderData `json:"data"`
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
		token:   token,
	}
}

func callGHN[T any](url string, token string) ([]T, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("token", token)

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

func callGHNWithBody[T any](
	method string,
	url string,
	token string,
	body any,
) (T, error) {
	var zero T

	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return zero, err
		}
		reqBody = bytes.NewBuffer(b)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return zero, err
	}

	req.Header.Set("token", token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return zero, err
	}
	defer resp.Body.Close()

	// Wrapper response từ GHN
	var ghnResp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    T      `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&ghnResp); err != nil {
		return zero, err
	}

	if ghnResp.Code != 200 {
		return zero, fmt.Errorf("GHN error: %s", ghnResp.Message)
	}

	return ghnResp.Data, nil
}

func (c *GHNClient) ResolveGHNAddress(
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

func (c *GHNClient) CreateOrder(request dto.GHNCreateOrderRequest) (string, error) {
	data, err := callGHNWithBody[CreateOrderData]("POST", fmt.Sprintf("%s/shiip/public-api/v2/shipping-order/create", c.baseURL), c.token, request)
	if err != nil {
		return "", err
	}
	return data.OrderCode, nil
}

func (c *GHNClient) CreateRequest(order model.Order, seller dto.UserResponse) (dto.GHNCreateOrderRequest, error) {

	shippingAddress := order.ShippingAddress

	request := dto.NewGHNCreateOrderRequest()

	//map address, sender info
	request.FromName = seller.Name
	request.FromPhone = seller.Phone
	request.FromAddress = seller.Address.AddressLine
	request.FromWardName = seller.Address.Ward
	request.FromDistrictName = seller.Address.District
	request.FromProvinceName = seller.Address.Province

	//map receiver info
	request.ToName = shippingAddress.FullName
	request.ToPhone = shippingAddress.Phone
	request.ToAddress = shippingAddress.AddressLine
	request.ToWardCode = shippingAddress.WardCode
	toDistrictId, _ := strconv.Atoi(shippingAddress.DistrictID)
	request.ToDistrictID = toDistrictId
	//map items
	for _, item := range order.Items {
		request.Items = append(request.Items, dto.GHNItem{
			Name:     item.ProductName,
			Code:     item.SKU,
			Quantity: item.Quantity,
			Price:    item.Price,
		})
	}

	//map cod if needed
	if order.PaymentMethod == "COD" {
		request.CodAmount = int(order.Total) + order.DeliveryFee
	}

	return request, nil
}

// CalculateFee calls GHN API to calculate shipping fee
func (c *GHNClient) CalculateFee(request dto.GHNCalculateFeeRequest) (*dto.GHNCalculateFeeResponse, error) {
	url := fmt.Sprintf("%s/shiip/public-api/v2/shipping-order/fee", c.baseURL)
	
	data, err := callGHNWithBody[dto.GHNCalculateFeeResponse]("POST", url, c.token, request)
	if err != nil {
		return nil, err
	}
	
	return &data, nil
}

func (c *GHNClient) CreateCalculateFeeRequest(order model.Order, seller dto.UserResponse) (*dto.GHNCalculateFeeRequest, error) {
	shippingAddress := order.ShippingAddress

	request := dto.NewGHNCalculateFeeRequest()

	//map address, sender info
	fromDistrictId, _ := strconv.Atoi(seller.Address.DistrictID)
	request.FromDistrictID = fromDistrictId
	request.FromWardCode = seller.Address.WardCode

	//map receiver info
	toDistrictId, _ := strconv.Atoi(shippingAddress.DistrictID)
	request.ToDistrictID = toDistrictId
	request.ToWardCode = shippingAddress.WardCode

	//
	request.ServiceID = order.DeliveryServiceID
	return &request, nil
}
	