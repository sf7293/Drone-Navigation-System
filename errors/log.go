package errors

import (
	"github.com/sf7293/Drone-Navigation-System/utils/logger"
)

func Log(err error) {
	defer func() { _ = logger.ZSLogger.Sync() }()

	Err, ok := err.(*Error)
	if !ok {
		logger.ZSLogger.Error(err)
		return
	}

	// log whatever details about Err
	ops := Err.Ops()
	data := GetData(Err)

	switch GetSeverity(err) {
	case SeverityInfo:
		logger.ZSLogger.Infow(Err.Error(), "ops", ops, "data", data)
	case SeverityWarning:
		logger.ZSLogger.Warnw(Err.Error(), "ops", ops, "data", data)
	default:
		logger.ZSLogger.Errorw(Err.Error(), "ops", ops, "data", data)
	}
}
