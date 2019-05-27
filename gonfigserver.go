package main

import (
	"github.com/gin-gonic/gin"
	mm "github.com/makitox/gonfigserver/mapmonitor"
	cl "github.com/op/go-logging"
)

var logger = cl.MustGetLogger("gonfigserver")
var configManager mm.Monitor

func main() {

	readCommandLIneFlags()

	if !*debug {
		gin.SetMode(gin.ReleaseMode)
	}

	var err error
	if configManager, err = NewMonitor(); err != nil {
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
		v1.GET("/pattern/:pattern/keylist/", fetchPatternKeyListForDefault)
	}

	router.Run(listenURL)
}
