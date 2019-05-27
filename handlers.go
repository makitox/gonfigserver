package main

import (
	"github.com/gin-gonic/gin"
	mm "github.com/makitox/gonfigserver/mapmonitor"
	"net/http"
	"regexp"
)

func fetchAllSpaceNames(c *gin.Context) {
	names := configManager.NamespacesList()
	if names == nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, names)
}

func fetchParametersNameList(c *gin.Context) {
	pNamespace := c.Param("namespace")
	names := configManager.KeysList(pNamespace)
	if names == nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, names)
}

func fetchParameter(c *gin.Context) {
	pNamespace := c.Param("namespace")
	pParamName := c.Param("paramName")

	value := configManager.Get(pNamespace, pParamName)
	if value == nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, *value)
}

func fetchParameterForDefault(c *gin.Context) {
	pParamName := c.Param("paramName")

	value := configManager.Get(mm.DefaultNamespace, pParamName)
	if value == nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, *value)
}

func fetchParametersNameListForDefault(c *gin.Context) {
	names := configManager.KeysList(mm.DefaultNamespace)
	if names == nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, names)
}

func fetchPatternKeyListForDefault(c *gin.Context){
	pPattern := c.Param("pattern")
	var compiledPattern = regexp.MustCompile(pPattern)
	names := configManager.KeysList(mm.DefaultNamespace)
	var matchedNames []string
	for _, value := range names {
		if compiledPattern.MatchString(value) {
			matchedNames = append(matchedNames, value)
		}
	}
	c.JSON(http.StatusOK, matchedNames)
}
