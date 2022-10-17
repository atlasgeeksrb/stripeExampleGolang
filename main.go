package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type CheckoutData struct {
	ClientSecret string `json:"client_secret"`
}

type Configuration struct {
	RouterUrl string
	ReactUrl  string
	ApiMode   string
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
	//@todo get paymentIntent from stripe
	c.IndentedJSON(http.StatusOK, gin.H{"result": "initiate"})
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
