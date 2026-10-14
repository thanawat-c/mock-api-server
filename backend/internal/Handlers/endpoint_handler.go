package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kan/mock-api/internal/models"
	"github.com/kan/mock-api/internal/services"
	"gorm.io/datatypes"
)

type UpdateEndpointInput struct {
	Path    string         `json:"path" binding:"required"`
	Method  string         `json:"method" binding:"required"`
	Delay   int            `json:"delay"`
	Headers datatypes.JSON `json:"headers"`
	Schema  datatypes.JSON `json:"schema" binding:"required"`
	Count   int            `json:"count"`
}

type CreateEndpointInput struct {
	Path    string         `json:"path" binding:"required"`
	Method  string         `json:"method" binding:"required"`
	Delay   int            `json:"delay"`
	Headers datatypes.JSON `json:"headers"`
	Schema  datatypes.JSON `json:"schema" binding:"required"`
	Count   int            `json:"count"`
}

func GetAllEndpoints(c * gin.Context) {
	res, err := services.GetAllEndpoints()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Get all endpoints successfully",
		"data": res,
	})
}

func GetEndpointById(c * gin.Context){
	id := c.Param("id")

	rec, err := services.GetEndpointById(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Get endpoint by id successfully",
		"data":    rec,
	})
}

func CreateEndpoint( c *gin.Context) {
	var endpoint CreateEndpointInput

	if err := c.ShouldBindJSON(&endpoint); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ep := models.Endpoint{
		Path:    endpoint.Path,
		Method:  endpoint.Method,
		Delay:   endpoint.Delay,
		Headers: endpoint.Headers,
		Schema:  endpoint.Schema,
	}
	
	if err := services.CreateEndpoint(&ep, endpoint.Count); err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			c.JSON(http.StatusConflict, gin.H{
				"error": "Endpoint with this path and method already exists",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Endpoint created successfully",
		"data":    ep,
	})
}

func DeleteEndpoint(c *gin.Context) {
	id := c.Param("id")

	if err := services.DeleteEndpoint(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Endpoint deleted successfully",
	})
}

func UpdateEndpoint(c *gin.Context) {
	id := c.Param("id")
	var input UpdateEndpointInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	endpointData := models.Endpoint{
		Path:    input.Path,
		Method:  input.Method,
		Delay:   input.Delay,
		Headers: input.Headers,
		Schema:  input.Schema,
	}

	if err := services.UpdateEndpoint(id, &endpointData, input.Count); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Endpoint Update successfully",
	})
}

func ResetEndpointRecords(c *gin.Context) {
	id := c.Param("id")
	records, err := services.ResetEndpointRecords(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Records reset to default successfully",
		"data":    records,
	})
}