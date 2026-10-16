package services

import (
	"github.com/kan/mock-api/internal/config"
	"github.com/kan/mock-api/internal/models"
)

func GetAllEndpoints() ([]models.Endpoint, error) {
	var endpoints []models.Endpoint

	if err := config.DB.Find(&endpoints).Error; err != nil {
		return nil, err
	}

	return endpoints, nil
}

type EndpointDetail struct {
	models.Endpoint
	Record []models.MockRecord `json:"records"`
}

func GetEndpointById(id string) (EndpointDetail, error) {
	var endpoint EndpointDetail

	if err := config.DB.Where("id = ?", id).First(&endpoint.Endpoint).Error; err != nil {
		return endpoint, err
	}

	if err := config.DB.Where("endpoint_id = ?", id).Order("id asc").Find(&endpoint.Record).Error; err != nil {
		return endpoint, err
	}

	return endpoint, nil
}

func CreateEndpoint(endpoint *models.Endpoint, count int) error {
	if err := config.DB.Create(endpoint).Error; err != nil {
		return err
	}

	// ถ้าเป็น GET ให้สุ่ม Mock Data เริ่มต้นตาม count (default 5)
	// ถ้าเป็น POST/PUT/DELETE ให้เริ่มต้นเป็น 0 แถว รอให้ผู้ใช้ยิงข้อมูลเข้ามาเอง
	if endpoint.Method == "GET" {
		if count <= 0 {
			count = 5
		}

		records := make([]models.MockRecord, 0, count)
		for i := 0; i < count; i++ {
			mock := GenerateMockData(endpoint.Schema)
			records = append(records, models.MockRecord{
				EndpointID: endpoint.ID,
				Data:       mock,
			})
		}

		// ⚡ Batch Insert: ยิง SQL บันทึกทั้งก้อนในคำสั่งเดียว ไม่ต้องต่อคิววนลูป
		if len(records) > 0 {
			if err := config.DB.Create(&records).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

func DeleteEndpoint(id string) error {
	if err := config.DB.Where("endpoint_id = ?", id).Delete(&models.MockRecord{}).Error; err != nil {
		return err
	}

	if err := config.DB.Where("id = ?", id).Delete(&models.Endpoint{}).Error; err != nil {
		return err
	}

	return nil
}

func UpdateEndpoint(id string, updatedata *models.Endpoint, count int) error {
	var existingendpoint models.Endpoint
	if err := config.DB.First(&existingendpoint, id).Error; err != nil {
		return err
	}

	existingendpoint.Path = updatedata.Path
	existingendpoint.Method = updatedata.Method
	existingendpoint.Delay = updatedata.Delay
	existingendpoint.Headers = updatedata.Headers
	existingendpoint.Schema = updatedata.Schema

	if err := config.DB.Save(&existingendpoint).Error; err != nil {
		return err
	}

	var records []models.MockRecord
	if err := config.DB.Where("endpoint_id = ?", existingendpoint.ID).Find(&records).Error; err != nil {
		return err
	}

	for _, rec := range records {
		rec.Data = SyncRecordWithSchema(rec.Data, existingendpoint.Schema)

		if err := config.DB.Save(&rec).Error; err != nil {
			return err
		}
	}

	return nil
}

// ResetEndpointRecords ลบข้อมูลเดิมทั้งหมดแล้วสุ่ม Mock Data เริ่มต้นขึ้นมาใหม่ตาม Schema
func ResetEndpointRecords(endpointID string) ([]models.MockRecord, error) {
	var endpoint models.Endpoint
	if err := config.DB.First(&endpoint, endpointID).Error; err != nil {
		return nil, err
	}

	// ลบข้อมูล Record เดิมทั้งหมดของ Endpoint นี้
	if err := config.DB.Where("endpoint_id = ?", endpoint.ID).Delete(&models.MockRecord{}).Error; err != nil {
		return nil, err
	}

	count := 5
	newRecords := make([]models.MockRecord, 0, count)
	for i := 0; i < count; i++ {
		mock := GenerateMockData(endpoint.Schema)
		newRecords = append(newRecords, models.MockRecord{
			EndpointID: endpoint.ID,
			Data:       mock,
		})
	}

	// ⚡ Batch Insert ทีเดียว
	if len(newRecords) > 0 {
		if err := config.DB.Create(&newRecords).Error; err != nil {
			return nil, err
		}
	}

	return newRecords, nil
}
