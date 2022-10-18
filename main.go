package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v73"
	"github.com/stripe/stripe-go/v73/paymentintent"
)

type CheckoutData struct {
	ClientSecret string `json:"client_secret"`
}

type Configuration struct {
	RouterUrl      string
	AcceptedOrigin string
	ApiMode        string
	StripeKey      string
}

type StripeApiStatus struct {
	Status  int64  `json:"status"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

var configuration Configuration

func main() {

	configuration = loadConfiguration()
	if configuration.ApiMode == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	// add OPTIONS for each route that is publicly accessible
	router.GET("/payment/:id", getPayment)
	router.OPTIONS("/payment/:id", preflight)

	router.POST("/initiatePayment", initiatePayment)
	router.OPTIONS("/initiatePayment", preflight)

	router.Run(configuration.RouterUrl)

}

func preflight(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", configuration.AcceptedOrigin) //"*"
	c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers,x-auth-token,x-auth-user,content-type")
	c.Header("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST")
	c.JSON(http.StatusOK, struct{}{})
}

func addHeaders(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", configuration.AcceptedOrigin) //"*"
	c.Header("Access-Control-Expose-Headers", "x-auth-token, x-auth-user")
}

func getPayment(c *gin.Context) {
	addHeaders(c)
	paymentId := c.Param("id")
	if len(strings.TrimSpace(paymentId)) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "payment intent id not found"})
		return
	}

	stripe.Key = configuration.StripeKey
	pi, err := paymentintent.Get(
		paymentId,
		nil,
	)
	if nil != err {
		errmsg := fmt.Sprint(err)
		c.JSON(http.StatusNotFound, gin.H{"error": errmsg})
		return
	}
	c.JSON(http.StatusOK, pi)
}

func initiatePayment(c *gin.Context) {
	addHeaders(c)

	// Keys- https://dashboard.stripe.com/apikeys
	stripe.Key = configuration.StripeKey
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(10100), // this is $101 = 10,100 cents
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		// can also set payment method types:
		// PaymentMethodTypes: []*string{
		// 	stripe.String("card"),
		// },
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}
	params.SetIdempotencyKey("RZzYwj0S5rssCYTo") //@todo dynamic
	result, err := paymentintent.New(params)
	//@todo inspect error, determine whether to retry
	if nil != err {
		var stripeStatus StripeApiStatus
		errmsg := fmt.Sprint(err)
		marshalErr := json.Unmarshal([]byte(errmsg), &stripeStatus)
		if nil != marshalErr {
			c.IndentedJSON(http.StatusBadRequest, errmsg)
		} else {
			c.IndentedJSON(int(stripeStatus.Status), stripeStatus)
		}
		return
	}

	c.IndentedJSON(http.StatusOK, result)
}

func loadConfiguration() Configuration {
	var configuration Configuration

	configpath, configerr := os.Getwd()
	if nil != configerr {
		panic(configerr)
	}

	file, _ := os.Open(configpath + string(os.PathSeparator) + "config.json")
	defer file.Close()

	decoder := json.NewDecoder(file)
	err := decoder.Decode(&configuration)
	if err != nil {
		panic(err)
	}
	return configuration
}
