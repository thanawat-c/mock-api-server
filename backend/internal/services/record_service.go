package services

import (
	"github.com/kan/mock-api/internal/config"
	"github.com/kan/mock-api/internal/models"
	"gorm.io/datatypes"
)

func UpdateRecordData(recordID string, data datatypes.JSON) error {
	return config.DB.Model(&models.MockRecord{}).Where("id = ?", recordID).Update("data", data).Error
}

func DeleteRecord(recordID string) error {
	return config.DB.Where("id = ?", recordID).Delete(models.MockRecord{}).Error
}

func AddRecordToEndpoint(endpointID string) (*models.MockRecord, error) {
	var endpoint models.Endpoint
	if err := config.DB.First(&endpoint, endpointID).Error; err != nil {
		return nil, err
	}
	data := GenerateMockData(endpoint.Schema)
	record := models.MockRecord{
		EndpointID: endpoint.ID,
		Data:       data,
	}
	if err := config.DB.Create(&record).Error; err != nil {
		return nil, err
	}

	return &record, nil
}
