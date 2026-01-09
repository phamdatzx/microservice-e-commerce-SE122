package controller

import (
	"encoding/json"
	"io"
	"net/http"
	"order-service/client"
	"order-service/client/payment"
	"order-service/dto"
	appError "order-service/error"
	"order-service/service"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v84"
)

type OrderController struct {
	service       service.OrderService
	paymentClient payment.PaymentClient
	GHNClient     *client.GHNClient
}

func NewOrderController(service service.OrderService, paymentClient payment.PaymentClient) *OrderController {
	return &OrderController{
		service:       service,
		paymentClient: paymentClient,
		GHNClient:     client.NewGHNClient(),
	}
}

func (c *OrderController) Checkout(ctx *gin.Context) {
	userID := ctx.GetHeader("X-User-Id")
	if userID == "" {
		ctx.Error(appError.NewAppError(401, "User ID not found in header"))
		return
	}

	var request dto.CheckoutRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.Error(appError.NewAppErrorWithErr(400, "Invalid request body", err))
		return
	}

	response, err := c.service.Checkout(userID, request)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, response)
}

func (c *OrderController) StripeWebhook(ctx *gin.Context) {
	const MaxBodyBytes = int64(65536)

	// Stripe yêu cầu raw body
	ctx.Request.Body = http.MaxBytesReader(
		ctx.Writer,
		ctx.Request.Body,
		MaxBodyBytes,
	)

	payload, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.Error(appError.NewAppErrorWithErr(400, "Cannot read webhook body", err))
		return
	}

	signature := ctx.GetHeader("Stripe-Signature")

	// Call stripeClient to construct event
	event, err := c.paymentClient.ConstructEvent(payload, signature)
	if err != nil {
		ctx.Error(err)
		return
	}

	switch event.Type {

	case "checkout.session.completed":
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			ctx.Error(appError.NewAppErrorWithErr(400, "Invalid session payload", err))
			return
		}

		orderID := session.Metadata["order_id"]
		if orderID == "" {
			ctx.Error(appError.NewAppError(400, "Missing order_id in metadata"))
			return
		}

		err := c.service.HandleCheckoutSessionCompleted(
			ctx,
			orderID,
			session.ID,
			session.PaymentIntent.ID,
		)
		if err != nil {
			ctx.Error(err)
			return
		}

	case "payment_intent.payment_failed":
		var pi stripe.PaymentIntent
		if err := json.Unmarshal(event.Data.Raw, &pi); err != nil {
			ctx.Error(appError.NewAppErrorWithErr(400, "Invalid payment_intent payload", err))
			return
		}

		err := c.service.HandlePaymentFailed(ctx, pi.ID)
		if err != nil {
			ctx.Error(err)
			return
		}

	default:
		// Ignore other events
	}

	// Stripe cần 200 OK
	ctx.JSON(200, gin.H{"status": "ok"})
}

func (c *OrderController) CreatePayment(ctx *gin.Context) {
	orderID := ctx.Param("orderId")
	if orderID == "" {
		ctx.Error(appError.NewAppError(400, "Order ID is required"))
		return
	}

	var request dto.CreatePaymentRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.Error(appError.NewAppErrorWithErr(400, "Invalid request body", err))
		return
	}

	response, err := c.service.CreatePaymentForOrder(ctx, orderID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, response)
}

func (c *OrderController) GetOrders(ctx *gin.Context) {
	userID := ctx.GetHeader("X-User-Id")
	if userID == "" {
		ctx.Error(appError.NewAppError(401, "User ID not found in header"))
		return
	}

	var request dto.GetOrdersRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.Error(appError.NewAppErrorWithErr(400, "Invalid query parameters", err))
		return
	}

	response, err := c.service.GetUserOrders(ctx, userID, request)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, response)
}

func (c *OrderController) Test(ctx *gin.Context) {
	district := ctx.Query("district")
	province := ctx.Query("province")
	ward := ctx.Query("forward")
	result, _ := c.GHNClient.ResolveGHNAddress(province, district, ward)
	ctx.JSON(200, result)
}
