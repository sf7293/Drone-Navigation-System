package domain

import (
	"github.com/sf7293/Drone-Navigation-System/domain"
	"github.com/sf7293/Drone-Navigation-System/errors"
)

func ReturnRouterResponseError(err errors.Cause) *domain.RouterResponseError {
	return &domain.RouterResponseError{
		Success: false,
		Data: domain.ErrorResponseData{
			Label:   err.TitleLabel,
			Message: err.MessageLabel,
		},
	}
}

func ReturnRouterRequestCalculateLocation() *domain.RouterRequestCalculateLocation {
	return &domain.RouterRequestCalculateLocation{
		XCord:    "0.2",
		YCord:    "2.3",
		ZCord:    "3.4",
		Velocity: "5",
	}
}

func ReturnRouterResponseCalculateLocation() *domain.RouterResponseCalculateLocation {
	return &domain.RouterResponseCalculateLocation{
		Location: 16.8,
	}
}
