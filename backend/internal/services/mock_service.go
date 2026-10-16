package services

import (
	"github.com/kan/mock-api/internal/config"
	"github.com/kan/mock-api/internal/models"
	"gorm.io/datatypes"
)

func GetMockRecordByPath(path string) ([]models.MockRecord, error) {
	var endpoint models.Endpoint

	err := config.DB.Where("path = ?", path).Order("id asc").First(&endpoint).Error
	if err != nil {
		return nil, err
	}

	var records []models.MockRecord
	err = config.DB.Where("endpoint_id = ?", endpoint.ID).Order("id asc").Find(&records).Error

	if err != nil {
		return nil, err
	}

	return records, nil

}

func CreateMockRecord(endpointID uint, data datatypes.JSON) (*models.MockRecord, error) {
	record := models.MockRecord{
		EndpointID: endpointID,
		Data:       data,
	}

	err := config.DB.Create(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}
