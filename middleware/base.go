package middleware

import (
	"github.com/gin-gonic/gin"
)

const (
	TimeFormat = "2006-01-02T15:04:05.000Z0700"
	UTC        = false
	Stack      = true
)

type MiddleConfig struct {
	//TimeFormat 时间格式
	TimeFormat string `yaml:"timeFormat" json:"timeFormat"`

	//UTC 是否使用UTC时间, 否则使用本地时间
	UTC bool `yaml:"utc" json:"utc"`

	//SkipPaths  跳过url
	SkipPaths []string `yaml:"skipPaths" json:"skipPaths"`

	//Stack 是否打印堆栈信息
	Stack bool `yaml:"stack" json:"stack"`
}

func NewMiddleConfig(args ...MiddleConfigOption) (config *MiddleConfig) {
	config = &MiddleConfig{
		TimeFormat: TimeFormat,
		UTC:        UTC,
		SkipPaths:  nil,
		Stack:      Stack,
	}

	for _, arg := range args {
		arg(config)
	}

	return config
}

type MiddleConfigOption func(middleConfig *MiddleConfig)

func MiddleConfigWithTimeFormat(timeFormat string) MiddleConfigOption {
	return func(middleConfig *MiddleConfig) {
		middleConfig.TimeFormat = timeFormat
	}
}

func MiddleConfigWithUTC(utc bool) MiddleConfigOption {
	return func(middleConfig *MiddleConfig) {
		middleConfig.UTC = utc
	}
}

func MiddleConfigWithSkipPath(skipPath []string) MiddleConfigOption {
	return func(middleConfig *MiddleConfig) {
		middleConfig.SkipPaths = skipPath
	}
}

func MiddleConfigWithStack(stack bool) MiddleConfigOption {
	return func(middleConfig *MiddleConfig) {
		middleConfig.Stack = stack
	}
}

type MidFunc func(config *MiddleConfig) gin.HandlerFunc

var middlewares = []MidFunc{
	GinZap,
	RecoveryWithZap,
}

func UseMiddles(handler *gin.Engine, config *MiddleConfig) {
	if len(middlewares) == 0 {
		return
	}

	for _, middleware := range middlewares {
		handler.Use(middleware(config))
	}
}
