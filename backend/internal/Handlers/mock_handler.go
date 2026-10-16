package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kan/mock-api/internal/config"
	"github.com/kan/mock-api/internal/models"
	"github.com/kan/mock-api/internal/services"
	"gorm.io/datatypes"
)

func HandleMockRequest(c *gin.Context) {
	rawParam := c.Param("any")
	cleanParam := "/" + strings.TrimPrefix(rawParam, "/")
	noSlashParam := strings.TrimPrefix(rawParam, "/")

	// 1. ลองค้นหา Endpoint ที่ตรงกับ Path และ Method ของคำขอ
	var endpoint models.Endpoint
	var itemID string = ""

	// ค้นหาแบบตรงทั้ง Path และ Method ก่อน
	err := config.DB.Where("(path = ? OR path = ?) AND method = ?", cleanParam, noSlashParam, c.Request.Method).First(&endpoint).Error
	if err != nil {
		// ถ้าไม่พบ ลองหาแบบมี ID ต่อท้าย เช่น /users/14 (ซึ่ง path ของ endpoint คือ /users)
		lastSlashIdx := strings.LastIndex(cleanParam, "/")
		if lastSlashIdx > 0 {
			basePath := cleanParam[:lastSlashIdx]
			noSlashBase := strings.TrimPrefix(basePath, "/")
			candidateID := cleanParam[lastSlashIdx+1:]

			// ตรวจสอบกับ Method ที่ส่งมา
			if err2 := config.DB.Where("(path = ? OR path = ?) AND method = ?", basePath, noSlashBase, c.Request.Method).First(&endpoint).Error; err2 == nil {
				itemID = candidateID
			} else {
				// หรือถ้าเป็น GET/PUT/DELETE โดยอ้างอิงจาก base path ที่อาจเป็น GET
				if err3 := config.DB.Where("path = ? OR path = ?", basePath, noSlashBase).First(&endpoint).Error; err3 == nil {
					itemID = candidateID
				}
			}
		}
	}

	// ถ้ายังไม่พบ endpoint ลองหาตาม Path โดยไม่ฟิกซ์ Method
	if endpoint.ID == 0 {
		err = config.DB.Where("path = ? OR path = ?", cleanParam, noSlashParam).First(&endpoint).Error
	}

	if endpoint.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Mock endpoint not found",
		})
		return
	}

	// --- 🏷️ 1. จำลอง Custom Response Headers ---
	if len(endpoint.Headers) > 0 {
		var headersMap map[string]string
		if err := json.Unmarshal(endpoint.Headers, &headersMap); err == nil {
			for key, val := range headersMap {
				if strings.TrimSpace(key) != "" {
					c.Header(key, val)
				}
			}
		}
	}

	// --- ⏱️ 2. จำลอง Response Delay ---
	if endpoint.Delay > 0 {
		time.Sleep(time.Duration(endpoint.Delay) * time.Millisecond)
	}

	// --- 💥 2. จำลอง Custom Status Code จาก Header เช่น x-mock-status: 500 หรือ Query param ?mock_status=401 ---
	customStatusStr := c.GetHeader("x-mock-status")
	if customStatusStr == "" {
		customStatusStr = c.Query("mock_status")
	}
	if customStatusStr != "" {
		if code, err := strconv.Atoi(customStatusStr); err == nil && code >= 100 && code < 600 {
			c.JSON(code, gin.H{
				"mock_status": code,
				"message":     fmt.Sprintf("Simulated custom status code %d", code),
			})
			return
		}
	}

	// 2. จัดการคำขอตาม Method
	switch c.Request.Method {
	case http.MethodGet:
		// ถ้ามี itemID ให้ค้นหา Record เฉพาะตัวนั้น
		if itemID != "" {
			var record models.MockRecord
			err := config.DB.Where("endpoint_id = ? AND id = ?", endpoint.ID, itemID).First(&record).Error
			if err != nil {
				// ค้นหาใน Data JSON เช่น {"id": itemID}
				var records []models.MockRecord
				config.DB.Where("endpoint_id = ?", endpoint.ID).Find(&records)
				found := false
				for _, rec := range records {
					var dataMap map[string]interface{}
					if jsonErr := json.Unmarshal(rec.Data, &dataMap); jsonErr == nil {
						if fmt.Sprintf("%v", dataMap["id"]) == itemID {
							record = rec
							found = true
							break
						}
					}
				}
				if !found {
					c.JSON(http.StatusNotFound, gin.H{"error": "Record not found with ID: " + itemID})
					return
				}
			}
			c.JSON(http.StatusOK, record.Data)
			return
		}

		// ดึงรายการทั้งหมด
		var records []models.MockRecord
		config.DB.Where("endpoint_id = ?", endpoint.ID).Order("id asc").Find(&records)

		// --- 🔍 3. Filter & Search ตาม Query Parameters ---
		queryParams := c.Request.URL.Query()
		searchQuery := strings.ToLower(queryParams.Get("search"))

		var filteredRecords []datatypes.JSON
		for _, rec := range records {
			var dataMap map[string]interface{}
			if err := json.Unmarshal(rec.Data, &dataMap); err != nil {
				filteredRecords = append(filteredRecords, rec.Data)
				continue
			}

			// ค้นหาด้วย ?search=... (ค้นหาทุกฟิลด์)
			if searchQuery != "" {
				matched := false
				for _, v := range dataMap {
					if strings.Contains(strings.ToLower(fmt.Sprintf("%v", v)), searchQuery) {
						matched = true
						break
					}
				}
				if !matched {
					continue
				}
			}

			// กรองตามฟิลด์เฉพาะ เช่น ?name=John หรือ ?email=...
			fieldMatched := true
			for qKey, qVals := range queryParams {
				if qKey == "page" || qKey == "limit" || qKey == "search" || qKey == "mock_status" {
					continue
				}
				if val, exists := dataMap[qKey]; exists {
					valStr := strings.ToLower(fmt.Sprintf("%v", val))
					targetVal := strings.ToLower(qVals[0])
					if !strings.Contains(valStr, targetVal) {
						fieldMatched = false
						break
					}
				}
			}
			if !fieldMatched {
				continue
			}

			filteredRecords = append(filteredRecords, rec.Data)
		}

		// --- 📄 4. Pagination: ?page=1&limit=10 ---
		pageStr := c.Query("page")
		limitStr := c.Query("limit")
		if pageStr != "" || limitStr != "" {
			page, _ := strconv.Atoi(pageStr)
			limit, _ := strconv.Atoi(limitStr)
			if page <= 0 {
				page = 1
			}
			if limit <= 0 {
				limit = 10
			}

			total := len(filteredRecords)
			startIndex := (page - 1) * limit
			endIndex := startIndex + limit

			if startIndex > total {
				filteredRecords = []datatypes.JSON{}
			} else {
				if endIndex > total {
					endIndex = total
				}
				filteredRecords = filteredRecords[startIndex:endIndex]
			}

			c.JSON(http.StatusOK, gin.H{
				"page":  page,
				"limit": limit,
				"total": total,
				"data":  filteredRecords,
			})
			return
		}

		c.JSON(http.StatusOK, filteredRecords)

	case http.MethodPost:
		var raw datatypes.JSON
		if err := c.ShouldBindJSON(&raw); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
			return
		}

		// ตรวจสอบว่า Body ที่ส่งมา มีฟิลด์ครบตาม Schema หรือไม่
		if missing := services.ValidateRequestBodyWithSchema(raw, endpoint.Schema); len(missing) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":          "Schema validation failed",
				"missing_fields": missing,
			})
			return
		}

		newRecord, err := services.CreateMockRecord(endpoint.ID, raw)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Record created successfully",
			"data":    newRecord,
		})

	case http.MethodPut:
		var raw datatypes.JSON
		if err := c.ShouldBindJSON(&raw); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
			return
		}

		// 1) ถ้าส่ง ID มาต่อท้าย URL ให้ค้นหาและอัปเดต Record นั้น
		if itemID != "" {
			var record models.MockRecord
			err := config.DB.Where("endpoint_id = ? AND id = ?", endpoint.ID, itemID).First(&record).Error
			if err != nil {
				// ค้นหาใน Data JSON เช่น {"id": itemID}
				var records []models.MockRecord
				config.DB.Where("endpoint_id = ?", endpoint.ID).Find(&records)
				found := false
				for _, rec := range records {
					var dataMap map[string]interface{}
					if jsonErr := json.Unmarshal(rec.Data, &dataMap); jsonErr == nil {
						if fmt.Sprintf("%v", dataMap["id"]) == itemID {
							record = rec
							found = true
							break
						}
					}
				}
				if !found {
					c.JSON(http.StatusNotFound, gin.H{"error": "Record not found to update"})
					return
				}
			}

			// อัปเดตข้อมูล
			record.Data = raw
			if err := config.DB.Save(&record).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "Record updated successfully",
				"data":    record,
			})
			return
		}

		// 2) ถ้าไม่ได้ใส่ ID ใน URL แต่ใน Body มีฟิลด์ "id" ให้ค้นหาและอัปเดตแถวที่มี id ตรงกัน
		var bodyMap map[string]interface{}
		if jsonErr := json.Unmarshal(raw, &bodyMap); jsonErr == nil && bodyMap["id"] != nil {
			targetID := fmt.Sprintf("%v", bodyMap["id"])
			var records []models.MockRecord
			config.DB.Where("endpoint_id = ?", endpoint.ID).Find(&records)
			for _, rec := range records {
				// เช็คทั้ง ID ของตารางใน DB และ "id" ใน Data JSON
				if fmt.Sprintf("%d", rec.ID) == targetID {
					rec.Data = raw
					config.DB.Save(&rec)
					c.JSON(http.StatusOK, gin.H{
						"message": "Record updated successfully",
						"data":    rec,
					})
					return
				}
				var recDataMap map[string]interface{}
				if err := json.Unmarshal(rec.Data, &recDataMap); err == nil {
					if fmt.Sprintf("%v", recDataMap["id"]) == targetID {
						rec.Data = raw
						config.DB.Save(&rec)
						c.JSON(http.StatusOK, gin.H{
							"message": "Record updated successfully",
							"data":    rec,
						})
						return
					}
				}
			}
		}

		// 3) ถ้าไม่มี ID ใดๆ เลย ให้บันทึกเป็น Record ใหม่
		newRecord, err := services.CreateMockRecord(endpoint.ID, raw)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Record saved successfully",
			"data":    newRecord,
		})

	case http.MethodDelete:
		if itemID != "" {
			// ลบตาม ID
			res := config.DB.Where("endpoint_id = ? AND id = ?", endpoint.ID, itemID).Delete(&models.MockRecord{})
			if res.RowsAffected == 0 {
				// ลองหาใน Data JSON
				var records []models.MockRecord
				config.DB.Where("endpoint_id = ?", endpoint.ID).Find(&records)
				foundID := uint(0)
				for _, rec := range records {
					var dataMap map[string]interface{}
					if jsonErr := json.Unmarshal(rec.Data, &dataMap); jsonErr == nil {
						if fmt.Sprintf("%v", dataMap["id"]) == itemID {
							foundID = rec.ID
							break
						}
					}
				}
				if foundID > 0 {
					config.DB.Delete(&models.MockRecord{}, foundID)
				} else {
					c.JSON(http.StatusNotFound, gin.H{"error": "Record not found to delete"})
					return
				}
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "Record with ID " + itemID + " deleted successfully",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Delete request received",
		})

	default:
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
	}
}
