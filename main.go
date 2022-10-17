package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v73"
	"github.com/stripe/stripe-go/v73/paymentintent"
)

type CheckoutData struct {
	ClientSecret string `json:"client_secret"`
}

type Configuration struct {
	RouterUrl string
	ReactUrl  string
	ApiMode   string
	StripeKey string
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
	router.GET("/payment", getPayment)
	router.OPTIONS("/payment", preflight)

	router.POST("/initiatePayment", initiatePayment)
	router.OPTIONS("/initiatePayment", preflight)

	router.POST("/retryPayment", retryPayment)
	router.OPTIONS("/retryPayment", preflight)

	router.POST("/completePayment", completePayment)
	router.OPTIONS("/completePayment", preflight)

	router.Run(configuration.RouterUrl)

}

func preflight(c *gin.Context) {
	//@todo add OPTIONS headers here
	c.JSON(http.StatusOK, struct{}{})
}

func addHeaders(c *gin.Context) {
	//@todo add HTTP headers here
}

func getPayment(c *gin.Context) {
	addHeaders(c)
	//@todo get payment status from stripe using paymentIntent ID
	c.IndentedJSON(http.StatusOK, gin.H{"result": "get"})
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
	result, err := paymentintent.New(params)

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

func retryPayment(c *gin.Context) {
	addHeaders(c)
	//@todo retry a payment using paymentIntent ID
	c.IndentedJSON(http.StatusOK, gin.H{"result": "retry"})
}

func completePayment(c *gin.Context) {
	addHeaders(c)
	//@todo final submit for payment
	c.IndentedJSON(http.StatusOK, gin.H{"result": "capture"})
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
