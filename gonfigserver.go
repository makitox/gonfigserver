package main

import (
	"github.com/gin-gonic/gin"
	mm "github.com/gonfigserver/mapmonitor"
	cl "github.com/op/go-logging"
	// "github.com/rjeczalik/notify"
)

var logger = cl.MustGetLogger("gonfigserver")
var configManager mm.Monitor

func main() {

	readCommandLIneFlags()

	if !*debug {
		gin.SetMode(gin.ReleaseMode)
	}

	var config = mm.MonitorConfiguration{}
	config.PropertyFileMask = mm.PropertyFileMask
	config.FailOnDuplicates = *pFailOnDup

	var err error
	if configManager, err = New(*propertyRoot, config); err != nil {
		panic(err)
	}

	router := gin.Default()
	v1 := router.Group("/api/v1")
	{
		v1.GET("/namespaces/", fetchAllSpaceNames)
		v1.GET("/namespace/:namespace/key/:paramName", fetchParameter)
		v1.GET("/namespace/:namespace/keylist", fetchParametersNameList)
		v1.GET("/key/:paramName", fetchParameterForDefault)
		v1.GET("/keylist", fetchParametersNameListForDefault)
		//v1.GET("/batch/", processBatchRequest)
		//v1.GET("/help/", fetchHelp)
	}

	router.Run(listenURL)
}
