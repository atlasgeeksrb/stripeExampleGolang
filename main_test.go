package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetPayment(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	getPayment(c)

	responseData, readerr := ioutil.ReadAll(w.Body) // []byte
	if readerr != nil {
		t.Fatal(readerr)
	}

	want := gin.H{
		"result": "aok",
	}
	var got gin.H
	err := json.Unmarshal(responseData, &got)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, want, got)
	assert.Equal(t, http.StatusOK, w.Code)
}
