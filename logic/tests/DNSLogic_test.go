package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sf7293/Drone-Navigation-System/domain"
	"github.com/sf7293/Drone-Navigation-System/logic"
	mocks_domain "github.com/sf7293/Drone-Navigation-System/mocks/domain"
)

type CalculateLocationTestCase struct {
	Input                domain.LogicRequestCalculateLocation
	WantedOutputLocation float64
	WantedOutputErr      interface{}
	WantedError          error
	FailedTestMessage    string
}

func TestCalculateLocation(t *testing.T) {
	tc := CalculateLocationTestCase{
		Input:                *mocks_domain.ReturnCalculateLocationLogicRequest(),
		WantedOutputLocation: 16.8,
		WantedError:          nil,
		FailedTestMessage:    "failed to call CalculateLocation method on DNSLogic",
	}

	l := &logic.DNSLogic{}
	logicReq := mocks_domain.ReturnCalculateLocationLogicRequest()
	resp, err := l.CalculateLocation(TestTracingSpan, context.Background(), *logicReq)
	assert.Nil(t, err, "error is not nil")
	assert.EqualValues(t, tc.WantedOutputLocation, resp.Location, tc.FailedTestMessage)
}
