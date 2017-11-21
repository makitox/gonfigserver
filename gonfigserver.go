package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"strconv"
	//"io"
	"log"
	"os"
	// "github.com/rjeczalik/notify"
)

var configManager *ConfigurationManager
var logger *log.Logger

var debug = flag.Bool("debug", false, "Add 'debug' option to start in debugging mode")
var iPort = flag.Int("port", 8080, "Specify port for Gonfig Server")
var pAddress = flag.String("address", "0.0.0.0", "Specify address for Gonfig Server")
var propertyRoot = flag.String("root", "", "Specify properties root")
var loggerRoot = flag.String("log", "", "Specify log file full path")

func main() {

	flag.Parse()

	// TODO: orginize right logging here
	//if loggerRoot == nil || *loggerRoot == "" {
	logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	//} else {
	//
	//}
	sPort := strconv.Itoa(*iPort)
	sAddress := *pAddress
	logger.Println("Starting Gonfig Server")
	logger.Println("port: " + sPort)
	logger.Println("address: " + sAddress)
	url := sAddress + ":" + sPort
	logger.Println("listening url: " + url)
	logger.Println("property root: " + *propertyRoot)
	logger.Println("logger root: " + *loggerRoot)
	if *debug {
		logger.Println("mode: debug")
	} else {
		logger.Println("mode: release")
		//gin.SetMode(gin.ReleaseMode)
	}

	configManager = NewConfigurationManager(*propertyRoot)

	router := gin.Default()
	v1 := router.Group("/api/v1")
	{
		v1.GET("/single/", fetchAllSpaceNames)
		v1.GET("/single/:namespace/:paramName", fetchParameter)
		v1.GET("/single/:namespace/", fetchParametersNameList)
		v1.GET("/batch/", processBatchRequest)
		v1.GET("/help/", fetchHelp)
	}

	router.Run(url)
	logger.Println("Gonfig Server finished")
}
