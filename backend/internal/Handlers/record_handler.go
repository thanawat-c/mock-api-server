package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kan/mock-api/internal/services"
	"gorm.io/datatypes"
)

type UpdateRecordInput struct {
	Data datatypes.JSON `json:"data" binding:"required"`
}

func UpdateRecordData( c * gin.Context){

	recordID := c.Param("id")

	var input UpdateRecordInput
	if err := c.ShouldBindJSON(&input); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	
	err := services.UpdateRecordData(recordID, input.Data); 
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Record data updated successfully",
		"data": input.Data,
	})

}

func DeleteRecord( c * gin.Context){
	recordID := c.Param("id")

	err := services.DeleteRecord(recordID)
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Record deleted successfully",
	})
}

func AddRecord( c * gin.Context){
	endpointID := c.Param("id")

	record, err := services.AddRecordToEndpoint(endpointID)
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"message": "Record added successfully",
		"data": record,
	})
}