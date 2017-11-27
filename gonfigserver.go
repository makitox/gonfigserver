package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	cl "github.com/op/go-logging"
	"strconv"
	// "github.com/rjeczalik/notify"
)

var configManager *ConfigurationManager
var logger = cl.MustGetLogger("gonfigserver")

var debug = flag.Bool("debug", false, "Add 'debug' option to start in debugging mode")

// var fswatch = flag.Bool("fswatch", false, "Watch changes in file system and reload properties map if any change")
// var git = flag.Bool("git", false, "Operate with property root as git repo")
// var loggerRoot = flag.String("log", "", "Specify log file full path")

var pFailOnDup = flag.Bool("failondup", true, "Fail to start if app finds key duplication")
var pReplaceChar = flag.String("replaceChar", ".", "Char to replace not allowed symbols in param key")
var iPort = flag.Int("port", 8080, "Specify port for Gonfig Server")
var pAddress = flag.String("address", "0.0.0.0", "Specify address for Gonfig Server")
var propertyRoot = flag.String("root", "", "Specify properties root")

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

	if cm, err := New(*propertyRoot); err != nil {
		panic(err)
	} else {
		configManager = cm
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
