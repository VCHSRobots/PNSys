// --------------------------------------------------------------------
// ginlogger.go -- A logger built for GIN
//
// Created 2017-12-18 DLB
// --------------------------------------------------------------------

package main

import (
	"epic/lib/log"
	"epic/lib/util"
	"github.com/gin-gonic/gin"
)

type GinLoggerType string

const (
	GinLoggerType_Error GinLoggerType = "Error"
	GinLoggerType_Info  GinLoggerType = "Info"
)

type GinLogger struct {
	logtype GinLoggerType
}

func init() {
	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()
	gin.DefaultWriter = NewGinLogger(GinLoggerType_Info)
	gin.DefaultErrorWriter = NewGinLogger(GinLoggerType_Error)
}

// NewGinLogger inializes the logger for use with GIN. Use the returned
// object for gin.DefaultWriter and gin.DefaultErrorWriter.
func NewGinLogger(logtype GinLoggerType) *GinLogger {
	logger := new(GinLogger)
	logger.logtype = logtype
	return logger
}

func (x *GinLogger) Write(p []byte) (n int, err error) {
	s := util.CleanStr(string(p), "")
	if x.logtype == GinLoggerType_Error {
		log.Errorf("Gin issued error!")
	}
	log.Passf(s)
	return len(p), nil
}
