package domain

import "github.com/sf7293/Drone-Navigation-System/domain"

func ReturnCalculateLocationLogicRequest() *domain.LogicRequestCalculateLocation {
	return &domain.LogicRequestCalculateLocation{
		XCord:    0.2,
		YCord:    2.3,
		ZCord:    3.4,
		Velocity: 5,
	}
}

func ReturnCalculateLocationLogicResponse() *domain.LogicResponseCalculateLocation {
	return &domain.LogicResponseCalculateLocation{
		Location: 16.8,
	}
}
