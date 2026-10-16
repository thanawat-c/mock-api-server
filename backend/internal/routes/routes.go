package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kan/mock-api/internal/handlers"
)

func SetupRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		v1.POST("/endpoints", handlers.CreateEndpoint)
		v1.GET("/endpoints", handlers.GetAllEndpoints)
		v1.GET("/endpoints/:id", handlers.GetEndpointById)
		v1.DELETE("/endpoints/:id", handlers.DeleteEndpoint)
		v1.PUT("/endpoints/:id", handlers.UpdateEndpoint)
		v1.PUT("/records/:id", handlers.UpdateRecordData)
		v1.DELETE("/records/:id", handlers.DeleteRecord)
		v1.POST("/endpoints/:id/records", handlers.AddRecord)
		v1.POST("/endpoints/:id/reset", handlers.ResetEndpointRecords)
	}

	r.Any("/mock/*any", handlers.HandleMockRequest)
}
