package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stripe/stripe-go/v73"
)

func TestGetPayment(t *testing.T) {

	configuration = loadConfiguration()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// add the payment id param
	p := gin.Param{
		Key:   "id",
		Value: "pi_3LtzZ4B88OE6Qmms0OU5FjJO",
	}
	c.Params = append(c.Params, p)
	getPayment(c)

	responseData, readerr := ioutil.ReadAll(w.Body) // []byte
	if readerr != nil {
		t.Fatal(readerr)
	}

	// want := gin.H{
	// 	"result": "get",
	// }
	var got stripe.PaymentIntent
	err := json.Unmarshal(responseData, &got)
	if err != nil {
		t.Fatal(err)
	}

	ref := reflect.ValueOf(got.ID)
	if ref.Kind() != reflect.String {
		assert.Fail(t, "Payment Intent ID not found in response")
	}

	// assert.Equal(t, want, got)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestInitiatePayment(t *testing.T) {
	configuration = loadConfiguration()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	initiatePayment(c)

	responseData, readerr := ioutil.ReadAll(w.Body) // []byte
	if readerr != nil {
		t.Fatal(readerr)
	}

	var got stripe.PaymentIntent
	err := json.Unmarshal(responseData, &got)
	if err != nil {
		t.Fatal(err)
	}

	ref := reflect.ValueOf(got.ID)
	if ref.Kind() != reflect.String {
		assert.Fail(t, "Payment Intent ID not found in response")
	}

	assert.Equal(t, http.StatusOK, w.Code)
}
