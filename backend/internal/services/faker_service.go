package services

import (
	"encoding/json"

	"github.com/brianvoe/gofakeit/v6"
	"gorm.io/datatypes"
)

type SchemaField struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func GenerateMockData(schemaJSON datatypes.JSON) datatypes.JSON {
	var schema []SchemaField
	if err := json.Unmarshal(schemaJSON, &schema); err != nil {
		return schemaJSON
	}

	result := make(map[string]interface{})

	for _, field := range schema {

		switch field.Type {
		case "faker.name":
			result[field.Name] = gofakeit.Name()
		case "faker.first_name":
			result[field.Name] = gofakeit.FirstName()
		case "faker.last_name":
			result[field.Name] = gofakeit.LastName()
		case "faker.email":
			result[field.Name] = gofakeit.Email()
		case "faker.phone":
			result[field.Name] = gofakeit.Phone()
		case "faker.gender":
			result[field.Name] = gofakeit.Gender()

		case "faker.int":
			result[field.Name] = gofakeit.Number(1, 1000)
		case "faker.price":
			result[field.Name] = gofakeit.Price(10, 5000)
		case "faker.boolean":
			result[field.Name] = gofakeit.Bool()

		case "faker.image", "faker.avatar":
			result[field.Name] = "https://i.pravatar.cc/300?u=" + gofakeit.UUID()

		case "faker.sentence":
			result[field.Name] = gofakeit.Sentence(6)
		case "faker.paragraph":
			result[field.Name] = gofakeit.Paragraph(2, 3, 5, " ")

		case "faker.city":
			result[field.Name] = gofakeit.City()
		case "faker.country":
			result[field.Name] = gofakeit.Country()
		case "faker.street":
			result[field.Name] = gofakeit.Street()
		case "faker.zip":
			result[field.Name] = gofakeit.Zip()

		case "faker.product_name":
			result[field.Name] = gofakeit.ProductName()
		case "faker.company":
			result[field.Name] = gofakeit.Company()

		case "faker.date":
			result[field.Name] = gofakeit.Date().Format("2006-01-02")
		case "faker.uuid":
			result[field.Name] = gofakeit.UUID()

		default:
			result[field.Name] = field.Type
		}
	}

	resultBytes, _ := json.Marshal(result)
	return datatypes.JSON(resultBytes)
}

func GenerateSingleFakeValue(fakerType string) interface{} {
	switch fakerType {
	case "faker.name":
		return gofakeit.Name()
	case "faker.first_name":
		return gofakeit.FirstName()
	case "faker.last_name":
		return gofakeit.LastName()
	case "faker.email":
		return gofakeit.Email()
	case "faker.phone":
		return gofakeit.Phone()
	case "faker.gender":
		return gofakeit.Gender()

	case "faker.int":
		return gofakeit.Number(1, 1000)
	case "faker.price":
		return gofakeit.Price(10, 5000)
	case "faker.boolean":
		return gofakeit.Bool()

	case "faker.image", "faker.avatar":
		return "https://i.pravatar.cc/300?u=" + gofakeit.UUID()

	case "faker.sentence":
		return gofakeit.Sentence(6)
	case "faker.paragraph":
		return gofakeit.Paragraph(2, 3, 5, " ")

	case "faker.city":
		return gofakeit.City()
	case "faker.country":
		return gofakeit.Country()
	case "faker.street":
		return gofakeit.Street()
	case "faker.zip":
		return gofakeit.Zip()

	case "faker.product_name":
		return gofakeit.ProductName()
	case "faker.company":
		return gofakeit.Company()

	case "faker.date":
		return gofakeit.Date().Format("2006-01-02")
	case "faker.uuid":
		return gofakeit.UUID()

	default:
		return fakerType
	}
}

func SyncRecordWithSchema(currentDataJSON datatypes.JSON, newSchemaJSON datatypes.JSON) datatypes.JSON {
	var currentData map[string]interface{}
	var newSchema []SchemaField

	json.Unmarshal(currentDataJSON, &currentData)
	json.Unmarshal(newSchemaJSON, &newSchema)

	if currentData == nil {
		currentData = make(map[string]interface{})
	}

	validKeys := make(map[string]bool)

	// วนลูปดู Schema ใหม่ ถ้าใน Record ปัจจุบันยังไม่มีคีย์นั้น ให้สุ่มค่าเติมเข้าไป
	for _, field := range newSchema {
		validKeys[field.Name] = true
		if _, exists := currentData[field.Name]; !exists {
			currentData[field.Name] = GenerateSingleFakeValue(field.Type)
		}
	}

	for key := range currentData {
		if !validKeys[key] {
			delete(currentData, key)
		}
	}

	resultBytes, _ := json.Marshal(currentData)
	return datatypes.JSON(resultBytes)
}

// ValidateRequestBodyWithSchema ตรวจสอบว่า JSON Body ที่ส่งเข้ามามีฟิลด์ครบตาม Schema หรือไม่
func ValidateRequestBodyWithSchema(bodyJSON datatypes.JSON, schemaJSON datatypes.JSON) []string {
	var schema []SchemaField
	if err := json.Unmarshal(schemaJSON, &schema); err != nil || len(schema) == 0 {
		return nil
	}

	var body map[string]interface{}
	if err := json.Unmarshal(bodyJSON, &body); err != nil {
		return []string{"invalid JSON format"}
	}

	var missingFields []string
	for _, field := range schema {
		if _, exists := body[field.Name]; !exists {
			missingFields = append(missingFields, field.Name)
		}
	}

	return missingFields
}
