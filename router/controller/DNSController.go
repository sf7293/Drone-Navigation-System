package controller

import (
	"context"
	"net/http"
	"strconv"

	"github.com/sf7293/Drone-Navigation-System/constants"
	"github.com/sf7293/Drone-Navigation-System/service/metrix"
	"github.com/sf7293/Drone-Navigation-System/utils/logger"

	"github.com/sf7293/Drone-Navigation-System/domain"
	"github.com/sf7293/Drone-Navigation-System/logic"

	"github.com/sf7293/Drone-Navigation-System/errors"

	"github.com/gin-gonic/gin"
)

type DNSControllerInterface interface {
	CalculateLocation(inputContext *gin.Context)
}

type DNSController struct {
	DNSLogic logic.DNSLogicInterface
}

func NewDNSController() DNSControllerInterface {
	dl := &logic.DNSLogic{}

	return &DNSController{
		DNSLogic: dl,
	}
}

func (r *DNSController) CalculateLocation(inputContext *gin.Context) {
	const op errors.Op = "Router.Controller.DNS.CalculateLocation"
	span := metrix.CreateSpan("Router.Controller.DNS.CalculateLocation", context.Background())
	defer func() {
		span.Finish()
		_ = logger.ZSLogger.Sync()
	}()

	req := &domain.RouterRequestCalculateLocation{}
	if err := inputContext.ShouldBind(req); err != nil {
		err = errors.E(op, errors.ErrRequestBinding, errors.Data{"error": err.Error()})
		returnErr(inputContext, 400, err, constants.LangCodeEnglish)
		return
	}

	x, err := strconv.ParseFloat(req.XCord, 64)
	if err != nil {
		err = errors.E(op, errors.ErrRequestBinding, errors.Data{"error": "X cord is not float!", "error_2": err.Error()})
		returnErr(inputContext, 400, err, constants.LangCodeEnglish)
		return
	}

	y, err := strconv.ParseFloat(req.YCord, 64)
	if err != nil {
		err = errors.E(op, errors.ErrRequestBinding, errors.Data{"error": "Y cord is not float!", "error_2": err.Error()})
		returnErr(inputContext, 400, err, constants.LangCodeEnglish)
		return
	}

	z, err := strconv.ParseFloat(req.ZCord, 64)
	if err != nil {
		err = errors.E(op, errors.ErrRequestBinding, errors.Data{"error": "Z cord is not float!", "error_2": err.Error()})
		returnErr(inputContext, 400, err, constants.LangCodeEnglish)
		return
	}

	vel, err := strconv.ParseFloat(req.Velocity, 64)
	if err != nil {
		err = errors.E(op, errors.ErrRequestBinding, errors.Data{"error": "Velocity cord is not float!", "error_2": err.Error()})
		returnErr(inputContext, 400, err, constants.LangCodeEnglish)
		return
	}

	DNSLogicReq := domain.LogicRequestCalculateLocation{
		XCord:    x,
		YCord:    y,
		ZCord:    z,
		Velocity: vel,
	}

	DNSLogicResp, err := r.DNSLogic.CalculateLocation(span, inputContext, DNSLogicReq)
	if err != nil {
		err = errors.E(op, err)
		returnErr(inputContext, 400, err, constants.LangCodeEnglish)
		return
	}

	resp := domain.RouterResponseCalculateLocation{
		Location: DNSLogicResp.Location,
	}

	inputContext.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    resp,
	})
}
