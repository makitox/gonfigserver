package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func fetchAllSpaceNames(c *gin.Context) {
	names := configManager.Namespaces()
	if names == nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, names)
}

func fetchParametersNameList(c *gin.Context) {
	pNamespace := c.Param("namespace")
	names := configManager.ConfigKeys(pNamespace)
	if names == nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, names)
}

func fetchParameter(c *gin.Context) {
	pNamespace := c.Param("namespace")
	pParamName := c.Param("paramName")

	value := configManager.ParameterValue(pNamespace, pParamName)
	if value == nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, *value)
}

func fetchParameterForDefault(c *gin.Context) {
	pParamName := c.Param("paramName")

	value := configManager.ParameterValue(DefaultNamespace, pParamName)
	if value == nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, *value)
}

func fetchParametersNameListForDefault(c *gin.Context) {
	names := configManager.ConfigKeys(DefaultNamespace)
	if names == nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, names)
}

func processBatchRequest(c *gin.Context) {
	// q := c.Query("q")
	c.Status(http.StatusNotImplemented)
}

func fetchHelp(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}
