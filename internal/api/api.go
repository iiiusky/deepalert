package api

import (
	"github.com/deepalert/deepalert/internal/handler"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const (
	contextArgumentKey = "handler.arguments"
	contextRequestID   = "request.id"
)

var logger = handler.Logger

func getArguments(c *gin.Context) *handler.Arguments {
	// In API code, handler.Arguments must be retrieved. If failed, the process must fail
	ptr, ok := c.Get(contextArgumentKey)
	if !ok {
		logger.Fatalf("Config is not set in API as '%s'", contextArgumentKey)
		return nil
	}

	args, ok := ptr.(*handler.Arguments)
	if !ok {
		logger.Fatalf("Config data as '%s' can not be casted", contextArgumentKey)
		return nil
	}

	return args
}

func getRequestID(c *gin.Context) string {
	// In API code, requestID must be retrieved. If failed, the process must fail
	ptr, ok := c.Get(contextRequestID)
	if !ok {
		logger.Fatalf("RequestID is not set in API as '%s'", contextRequestID)
	}

	reqID, ok := ptr.(string)
	if !ok {
		logger.Fatalf("RequestID as '%s' can not be casted", contextRequestID)
	}

	return reqID
}

type errorResponse struct {
	Error interface{} `json:"error"`
}

func resp(c *gin.Context, data interface{}) {
	reqID := getRequestID(c)
	c.Header("DeepAlert-Request-ID", reqID)

	switch v := data.(type) {
	case apiError:
		logger.WithError(v).WithField("message", v.Message()).Error("API Error")
		c.JSON(v.StatusCode(), v.Error())
	case error:
		logger.WithError(v).Error("API Error (not apiError type)")
		c.JSON(500, "SystemError")
	default:
		c.JSON(200, data)
	}
}

const paramReportID = "report_id"

// SetupRoute binds route of gin and API
func SetupRoute(r *gin.RouterGroup, args *handler.Arguments) {
	r.Use(func(c *gin.Context) {
		reqID := uuid.New().String()
		logger.WithFields(logrus.Fields{
			"path":       c.FullPath(),
			"params":     c.Params,
			"request_id": reqID,
			"remote":     c.ClientIP(),
			"ua":         c.Request.UserAgent(),
		}).Info("API request")

		c.Set(contextRequestID, reqID)
		c.Set(contextArgumentKey, args)
		c.Next()
	})

	r.POST("/alert", postAlert)
	r.GET("/report/:"+paramReportID+"/alert", getReportAlerts)
	r.GET("/report/:"+paramReportID+"/section", getReportSections)
	r.GET("/report/:"+paramReportID+"/attribute", getReportAttributes)
}