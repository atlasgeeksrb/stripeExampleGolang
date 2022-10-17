package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var routerUrl string = "localhost:8080"

type CheckoutData struct {
	ClientSecret string `json:"client_secret"`
}

func main() {

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

	router.Run(routerUrl)

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
