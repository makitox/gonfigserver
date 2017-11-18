package main

import (
	"github.com/gin-gonic/gin"
)

var parameters = make(map[string]map[string]string)
var namespacesParamsNames = make(map[string][]string)
var namespaces []string

func main() {
	readParameters()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	v1 := router.Group("/api/v1")
	{
		v1.GET("/single/", fetchAllSpaceNames)
		v1.GET("/single/:namespace/:paramName", fetchParameter)
		v1.GET("/single/:namespace/", fetchParametersNameList)
		v1.GET("/batch/", processBatchRequest)
		v1.GET("/help/", fetchHelp)
	}

	router.Run("0.0.0.0:8080")
}

func readParameters() {
	parameters["default"] = make(map[string]string)
	(parameters["default"])["param"] = "paramvalue"
	(parameters["default"])["antiparam"] = "antiparam_value"
	(parameters["default"])["zero"] = "null"

	prepareIndexes()
}

func prepareIndexes() {
	for namespace := range parameters {
		namespaces = append(namespaces, namespace)
		for paramName := range parameters[namespace] {
			namespacesParamsNames[namespace] = append(namespacesParamsNames[namespace], paramName)
		}
	}
}
