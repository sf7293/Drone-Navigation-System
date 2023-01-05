package domain

type LogicRequestCalculateLocation struct {
	XCord    float64
	YCord    float64
	ZCord    float64
	Velocity float64
}

type LogicResponseCalculateLocation struct {
	Location float64
}
