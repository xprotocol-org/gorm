package schema_test

import (
	"sync"
	"testing"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func TestVirtualFields(t *testing.T) {
	type User struct {
		ID              uint   `gorm:"primaryKey"`
		Name            string
		Email           string
		ComputedField   string `gorm:"virtual"`
		AnotherHelper   int    `gorm:"virtual"`
	}

	userSchema, err := schema.Parse(&User{}, &sync.Map{}, schema.NamingStrategy{})
	if err != nil {
		t.Fatalf("failed to parse user schema, got error %v", err)
	}

	expectedDBNames := []string{"id", "name", "email"}
	if len(userSchema.DBNames) != len(expectedDBNames) {
		t.Errorf("expected %d DBNames, got %d: %v", len(expectedDBNames), len(userSchema.DBNames), userSchema.DBNames)
	}

	for _, expectedName := range expectedDBNames {
		found := false
		for _, dbName := range userSchema.DBNames {
			if dbName == expectedName {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected DBName %s not found in %v", expectedName, userSchema.DBNames)
		}
	}

	for _, virtualName := range []string{"computed_field", "another_helper"} {
		for _, dbName := range userSchema.DBNames {
			if dbName == virtualName {
				t.Errorf("virtual field %s should not be in DBNames: %v", virtualName, userSchema.DBNames)
			}
		}
	}

	if field := userSchema.LookUpField("ComputedField"); field == nil {
		t.Error("ComputedField should still be accessible via LookUpField")
	}
}

func TestVirtualFieldsWithGormModel(t *testing.T) {
	type Product struct {
		gorm.Model
		Name          string
		Price         float64
		CachedTotal   float64 `gorm:"virtual"`
		DisplayString string  `gorm:"virtual"`
	}

	productSchema, err := schema.Parse(&Product{}, &sync.Map{}, schema.NamingStrategy{})
	if err != nil {
		t.Fatalf("failed to parse product schema, got error %v", err)
	}

	for _, dbName := range productSchema.DBNames {
		if dbName == "cached_total" || dbName == "display_string" {
			t.Errorf("virtual field %s should not be in DBNames: %v", dbName, productSchema.DBNames)
		}
	}

	expectedFields := []string{"id", "created_at", "updated_at", "deleted_at", "name", "price"}
	if len(productSchema.DBNames) != len(expectedFields) {
		t.Errorf("expected %d DBNames, got %d: %v", len(expectedFields), len(productSchema.DBNames), productSchema.DBNames)
	}
}

