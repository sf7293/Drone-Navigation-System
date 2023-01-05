package controller

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sf7293/Drone-Navigation-System/config"
	core_config "github.com/sf7293/Drone-Navigation-System/config/core"
	"github.com/sf7293/Drone-Navigation-System/constants"
	"github.com/sf7293/Drone-Navigation-System/domain"
	"github.com/sf7293/Drone-Navigation-System/errors"
)

func returnErr(ctx *gin.Context, responseHTTPStatusCode int, err error, langCode string) {
	env := os.Getenv("env")
	if len(env) == 0 {
		env = core_config.EnvDev
	}

	errors.Log(err)

	errorCause := errors.GetCause(err)
	if strings.Contains(langCode, "-") {
		langCode = strings.Split(langCode, "-")[0]
	}

	_, langIsDefinied := config.TranslationMap[langCode]
	if !langIsDefinied {
		langCode = constants.LangCodeEnglish
	}

	errorTitleTranslation, ok := config.TranslationMap[langCode][errorCause.TitleLabel]
	if !ok {
		errorTitleTranslation, ok = config.TranslationMap[langCode][errors.ErrInternal.TitleLabel]
		if !ok {
			errorTitleTranslation = config.TranslationMap[constants.LangCodeEnglish][errors.ErrInternal.TitleLabel]
		}
	}

	errorTitle := ""
	if errorTitleTranslation != nil {
		errorTitle = errorTitleTranslation.Translation
	}

	errorMessageTranslation, ok := config.TranslationMap[langCode][errorCause.MessageLabel]
	if !ok {
		errorMessageTranslation, ok = config.TranslationMap[langCode][errors.ErrInternal.MessageLabel]
		if !ok {
			errorMessageTranslation = config.TranslationMap[constants.LangCodeEnglish][errors.ErrInternal.MessageLabel]
		}
	}

	errorMessage := ""
	if errorMessageTranslation != nil {
		errorMessage = errorMessageTranslation.Translation
	}

	serverError := ""
	var serverErrorData map[string]interface{}
	serverErrorData = nil

	if env == core_config.EnvDev {
		serverError = err.Error()

		childErr, ok := err.(*errors.Error)
		if ok {
			serverError = childErr.Error()
			serverErrorData = errors.GetData(childErr)
		}
	}

	var publicErrorData map[string]interface{}
	publicErrorData = nil
	errorData := errors.GetData(err)
	if errorData != nil {
		publicErrorDataMap, ok := errorData["public_error_data"]
		if ok {
			publicErrorData = publicErrorDataMap.(map[string]interface{})
		}
	}

	//logger.ZSLogger.Errorw("error returned to user", "lang_code", langCode, "title", errorTitle, "message", errorMessage, "button_text", errorButtonText, "button_type", errorCause.ButtonType, "button_action", errorCause.ButtonAction)
	//_ = logger.ZSLogger.Sync()

	responseError := domain.RouterResponseError{
		Success: false,
		Data: domain.ErrorResponseData{
			Label:           errorCause.TitleLabel,
			Title:           errorTitle,
			Message:         errorMessage,
			ServerError:     serverError,
			ServerErrorData: serverErrorData,
			PublicErrorData: publicErrorData,
		},
	}

	ctx.JSON(responseHTTPStatusCode, responseError)
}
