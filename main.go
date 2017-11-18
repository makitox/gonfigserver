package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
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

func fetchAllSpaceNames(c *gin.Context) {
	c.JSON(http.StatusOK, namespaces)
}

func fetchParametersNameList(c *gin.Context) {
	pNamespace := c.Param("namespace")
	names, ok := namespacesParamsNames[pNamespace]
	if !ok {
		c.Status(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, names)
}

func fetchParameter(c *gin.Context) {
	pNamespace := c.Param("namespace")
	pParamName := c.Param("paramName")

	namespace, ok := parameters[pNamespace]
	if !ok {
		c.Status(http.StatusNotFound)
		return
	}

	paramValue, ok := namespace[pParamName]
	if !ok {
		c.Status(http.StatusNotFound)
		return
	}
	c.String(http.StatusOK, paramValue)
}

func processBatchRequest(c *gin.Context) {
	// q := c.Query("q")
	c.Status(http.StatusNotImplemented)
}

func fetchHelp(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}
