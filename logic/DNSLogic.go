package logic

import (
	"context"

	"github.com/sf7293/Drone-Navigation-System/config"
	"github.com/sf7293/Drone-Navigation-System/service/metrix"
	"github.com/sf7293/Drone-Navigation-System/utils/logger"

	"github.com/opentracing/opentracing-go"

	"github.com/sf7293/Drone-Navigation-System/errors"

	"github.com/sf7293/Drone-Navigation-System/domain"
)

type DNSLogicInterface interface {
	CalculateLocation(inputSpan opentracing.Span, inputContext context.Context, req domain.LogicRequestCalculateLocation) (resp domain.LogicResponseCalculateLocation, err error)
}

// DNSLogic handles the top-level management of stuffs that are related to Drone Navigation System.
type DNSLogic struct{}

func (l DNSLogic) CalculateLocation(inputSpan opentracing.Span, inputContext context.Context, req domain.LogicRequestCalculateLocation) (resp domain.LogicResponseCalculateLocation, err error) {
	const op errors.Op = "Logic.DNS.CalculateLocation"
	span := metrix.CreateChildSpan("Logic.DNS.CalculateLocation", inputSpan)
	defer func() {
		span.Finish()
		_ = logger.ZSLogger.Sync()
	}()

	location := req.XCord*float64(config.App.SectorID) + req.YCord*float64(config.App.SectorID) + req.ZCord*float64(config.App.SectorID) + req.Velocity

	resp = domain.LogicResponseCalculateLocation{
		Location: location,
	}

	return resp, nil
}
