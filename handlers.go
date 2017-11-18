package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

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
