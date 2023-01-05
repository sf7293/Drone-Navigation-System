package tests

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/sf7293/Drone-Navigation-System/domain"
	mocks_domain "github.com/sf7293/Drone-Navigation-System/mocks/domain"
	mocks_logic "github.com/sf7293/Drone-Navigation-System/mocks/logic"
	"github.com/sf7293/Drone-Navigation-System/router/controller"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type CalculateLocationTestCase struct {
	Input                domain.RouterRequestCalculateLocation
	WantedOutputLocation float64
	WantedOutputErr      interface{}
	WantedStatusCode     int64
	WantedError          error
	FailedTestMessage    string
}

func TestCalculateLocation(t *testing.T) {
	mockedRequest := mocks_domain.ReturnRouterRequestCalculateLocation()
	mockedResponse := mocks_domain.ReturnRouterResponseCalculateLocation()
	tc := CalculateLocationTestCase{
		Input:                *mockedRequest,
		WantedOutputLocation: mockedResponse.Location,
		WantedStatusCode:     200,
		WantedError:          nil,
		FailedTestMessage:    "failed to call CalculateLocation method on DNSController",
	}

	mockedLogicRequest := *mocks_domain.ReturnCalculateLocationLogicRequest()
	mockedLogicResponse := mocks_domain.ReturnCalculateLocationLogicResponse()
	respRecorder := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	marshalledRequest, err := json.Marshal(mockedRequest)
	assert.Nil(t, err, tc.FailedTestMessage)
	ginContext, _ := gin.CreateTestContext(respRecorder)
	ginContext.Request, err = http.NewRequest("POST", "localhost:8080/v1/dns/location", bytes.NewBuffer(marshalledRequest))
	assert.Nil(t, err, tc.FailedTestMessage)
	ginContext.Request.Header.Add("Content-Type", "application/json")

	ctrl := gomock.NewController(t)
	mockedLogic := mocks_logic.NewMockDNSLogicInterface(ctrl)
	mockedLogic.EXPECT().CalculateLocation(TestTracingSpan, ginContext, mockedLogicRequest).Return(*mockedLogicResponse, nil)

	r := controller.DNSController{
		DNSLogic: mockedLogic,
	}

	r.CalculateLocation(ginContext)

	type ResponseModel struct {
		Success bool                                   `json:"status"`
		Data    domain.RouterResponseCalculateLocation `json:"data"`
	}

	var response ResponseModel
	err = json.Unmarshal(respRecorder.Body.Bytes(), &response)
	assert.Nil(t, err, tc.FailedTestMessage)

	assert.EqualValues(t, tc.WantedStatusCode, respRecorder.Code, tc.FailedTestMessage)
	assert.EqualValues(t, tc.WantedOutputLocation, response.Data.Location, tc.FailedTestMessage)
}
