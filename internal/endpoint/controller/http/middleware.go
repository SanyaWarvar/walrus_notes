package http

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	apperrors "wn/internal/errors"
	"wn/pkg/apperror"
	"wn/pkg/applogger"
	"wn/pkg/token"
	"wn/pkg/util"

	"github.com/gin-contrib/cors"
	"github.com/google/uuid"

	"wn/pkg/constants"
	"wn/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func HeaderCtxHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), constants.RequestIdCtx, c.Request.Header.Get(constants.RequestIdHeader))
		ctx = context.WithValue(ctx, constants.ApiNameCtx, c.Request.URL.Path)
		c.Request = c.Request.WithContext(ctx)

	}
}

func RequestIdValidationHandler(c *gin.Context) {
	_, err := uuid.Parse(c.Request.Header.Get(constants.RequestIdHeader))
	if err != nil {
		_ = c.Error(apperror.NewBadRequestError(err.Error(), "invalid_"+constants.RequestIdHeader))
		c.Abort()

	}

}

func AuthorizationHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader(constants.AuthorizationHeader)

		headerParts := strings.Fields(header)
		if len(headerParts) != 2 || len(headerParts[1]) == 0 {
			_ = c.Error(apperrors.InvalidAuthorizationHeader)
			c.Abort()
			return
		}

		claims, err := token.ParseTokenWithoutKeyCheck(headerParts[1])
		if err != nil {
			_ = c.Error(apperrors.InvalidTokenError)
			c.Abort()
			return
		}
		ctx := context.WithValue(c.Request.Context(), constants.UserRoleCtx, token.GetUserRole(claims))
		ctx = context.WithValue(ctx, constants.UserIdCtx, token.GetUserId(claims).String())
		ctx = context.WithValue(ctx, "mainLayoutId", token.GetMainLayoutId(claims).String())

		c.Request = c.Request.WithContext(ctx)
	}
}

func LoggerHandler(logger applogger.Logger, logInputParamOnErr bool) gin.HandlerFunc {

	return func(c *gin.Context) {
		var (
			requestBody []byte
			clientIP    string
			userAgent   string
			method      string
			headers     http.Header
		)

		ctx := context.WithValue(c.Request.Context(), constants.RequestIdCtx, c.Request.Header.Get(constants.RequestIdHeader))
		ctx = context.WithValue(ctx, constants.ApiNameCtx, c.Request.URL.Path)
		c.Request = c.Request.WithContext(ctx)

		if logInputParamOnErr {
			method = c.Request.Method
			clientIP = c.ClientIP()
			userAgent = c.Request.UserAgent()
			headers = c.Request.Header
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		c.Next()

		if len(c.Errors) > 0 {

			logger.WithCtx(c.Request.Context()).Errorf(c.Errors.Last().Error())
			if len(string(requestBody)) > 2000 {
				return
			}
			if logInputParamOnErr {
				msg := fmt.Sprintf("ERR: %s, IP: %s, AGENT: %s, METHOD: %s, STATUS: %d, QUERY_PARAM: %s, HEADERS: %v, REQUEST_BODY: %v",
					c.Errors.Last().Err.Error(), clientIP, userAgent, method, c.Writer.Status(), c.Request.URL.RawQuery, util.MaskHeaders(headers), string(requestBody))
				logger.WithCtx(c.Request.Context()).Errorf(msg)
			}
		}
	}
}

func ErrorHandler(builder *response.Builder) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) <= 0 {
			return
		}

		err := c.Errors.Last().Err
		status := http.StatusInternalServerError
		appError, ok := errors.Cause(err).(*apperror.AppError)
		if ok {
			status = apperror.GetHttpStatusByErrorType(appError.Type)
		} else {
			appError = apperror.NewInternalError(err)
		}
		c.AbortWithStatusJSON(status, builder.BuildErrorResponse(c.Request.Context(), appError.Message, appError.Code, err))
	}
}

func CORSHandler() gin.HandlerFunc {
	return cors.New(cors.Config{

		//Адреса которые могут обращаться к нам
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		//Методы которые могут кидать к нам
		AllowMethods: []string{"DELETE", "GET", "POST", "PUT", "OPTIONS"},
		//Хедеры которые могут кидать к нам
		AllowHeaders: []string{"Origin", "Content-Type", constants.AuthorizationHeader, constants.RequestIdHeader, constants.RefreshHeader},
		//Допустимы параметры авторизации в куках
		AllowCredentials: true,
		//хедеры которые я могу прокинуть клиенту
		ExposeHeaders: []string{"Content-Length", "Content-Type", constants.AuthorizationHeader, constants.RequestIdHeader, constants.RefreshHeader},
		//время хранения префлайт запрсоов
		MaxAge: 12 * time.Hour,
	})
}
