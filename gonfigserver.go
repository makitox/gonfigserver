package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	mm "github.com/gonfigserver/mapmonitor"
	cl "github.com/op/go-logging"
	"strconv"
	// "github.com/rjeczalik/notify"
)

var logger = cl.MustGetLogger("gonfigserver")
var configManager mm.Monitor

func main() {
	flag.Parse()

	sPort := strconv.Itoa(*iPort)
	sAddress := *pAddress
	logger.Info("Starting Gonfig Server with configuration:")
	logger.Info("\tport: " + sPort)
	logger.Info("\taddress: " + sAddress)
	url := sAddress + ":" + sPort
	logger.Info("\tlistening url: " + url)
	logger.Info("\tproperty root: " + *propertyRoot)
	//logger.Info("\tlogger root: " + *loggerRoot)
	if *debug {
		logger.Info("\tmode: debug")
	} else {
		logger.Info("\tmode: release")
		gin.SetMode(gin.ReleaseMode)
	}
	logger.Info("\tfail on paramkey duplicates: " + strconv.FormatBool(*pFailOnDup))
	logger.Info("\treplace character: '" + *pReplaceChar + "'")
	logger.Info("\n")

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

	router.Run(url)
}
