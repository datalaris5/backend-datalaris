package utils

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
)

func GenerateSlug(name string) string {
	slug := strings.ToLower(name)

	slug = strings.ReplaceAll(slug, " ", "-")

	reg := regexp.MustCompile(`[^a-z0-9-]+`)
	slug = reg.ReplaceAllString(slug, "")

	regDash := regexp.MustCompile(`-+`)
	slug = regDash.ReplaceAllString(slug, "-")

	slug = strings.Trim(slug, "-")

	return slug
}

func GenerateCode(name string, count int64) string {
	next := count + 1
	slug := GenerateSlug(name)
	code := fmt.Sprintf("%s-%d", slug, next)
	return code
}

func UpdateCode(name string, existingCode string) string {
	slug := GenerateSlug(name)

	parts := strings.Split(existingCode, "-")
	if len(parts) == 0 {
		return slug
	}

	last := parts[len(parts)-1]
	num, err := strconv.Atoi(last)
	if err != nil {
		num = 1
	}

	return fmt.Sprintf("%s-%d", slug, num)
}

func TernaryString(condition bool, a, b func() string) string {
	if condition {
		return a()
	}
	return b()
}

func GenerateUniqueFileName(galleryableType, originalFilename string) string {
	ext := filepath.Ext(originalFilename) // contoh: .jpg
	timestamp := time.Now().Format("20060102_150405")
	randStr := randomString(6)

	return fmt.Sprintf("%s_%s_%s%s", galleryableType, timestamp, randStr, ext)
}

func NormalizeStringPointers[T any](input T) T {
	val := reflect.ValueOf(&input).Elem()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		// handle pointer ke time.Time
		if field.Kind() == reflect.Ptr && field.Type().Elem().String() == "time.Time" {
			if !field.IsNil() {
				t := field.Elem().Interface().(time.Time)
				if t.IsZero() {
					field.Set(reflect.Zero(field.Type())) // set jadi nil
				}
			}
		}

		// string pointer ga usah diubah! biarin "" tetep ""
	}

	return input
}

func BuildUpdateMap[T any](input T) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	v := reflect.ValueOf(input)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()

	if v.Kind() != reflect.Struct {
		return nil, errors.New("input must be a struct or pointer to struct")
	}

	for i := 0; i < v.NumField(); i++ {
		fieldVal := v.Field(i)
		fieldType := t.Field(i)
		fieldName := fieldType.Name

		// Skip unexported fields
		if !fieldVal.CanInterface() {
			continue
		}

		// Check for embedded struct (anonymous)
		if fieldType.Anonymous && fieldVal.Kind() == reflect.Struct {
			embeddedUpdates, err := BuildUpdateMap(fieldVal.Interface())
			if err != nil {
				return nil, err
			}
			for k, v := range embeddedUpdates {
				result[k] = v
			}
			continue
		}

		// Skip certain fields
		if shouldSkipField(fieldName) {
			continue
		}

		// Get DB column name from tag, fallback to snake_case
		tag := fieldType.Tag.Get("gorm")
		column := parseGormColumn(tag)
		if column == "" {
			column = toSnakeCase(fieldName)
		}

		// Skip nil pointer values
		if fieldVal.Kind() == reflect.Ptr && fieldVal.IsNil() {
			continue
		}

		result[column] = fieldVal.Interface()
	}

	return result, nil
}

func shouldSkipField(fieldName string) bool {
	skipped := map[string]bool{
		"CreatedAt": true,
		"CreatedBy": true,
		"UpdatedAt": true, // biar diisi dari BeforeUpdate
		"UpdatedBy": true,
		"DeletedAt": true,
	}
	return skipped[fieldName]
}

func parseGormColumn(tag string) string {
	// parse tag like: gorm:"column:updated_by;otherTag"
	for _, part := range strings.Split(tag, ";") {
		if strings.HasPrefix(part, "column:") {
			return strings.TrimPrefix(part, "column:")
		}
	}
	return ""
}

func toSnakeCase(str string) string {
	// Tangani ID sebagai kasus khusus
	if str == "ID" {
		return "id"
	}

	var result []rune
	for i, r := range str {
		if unicode.IsUpper(r) {
			if i > 0 {
				result = append(result, '_')
			}
			result = append(result, unicode.ToLower(r))
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}

func ApplyNonNilFields[T any](target *T, source T) error {
	targetVal := reflect.ValueOf(target)
	if targetVal.Kind() != reflect.Ptr || targetVal.IsNil() {
		return errors.New("target must be a non-nil pointer to struct")
	}
	targetVal = targetVal.Elem()

	sourceVal := reflect.ValueOf(source)
	if sourceVal.Kind() == reflect.Ptr {
		sourceVal = sourceVal.Elem()
	}

	return applyNonNilFieldsReflect(targetVal, sourceVal)
}

func applyNonNilFieldsReflect(targetVal, sourceVal reflect.Value) error {
	for i := 0; i < sourceVal.NumField(); i++ {
		srcField := sourceVal.Field(i)
		fieldType := sourceVal.Type().Field(i)
		dstField := targetVal.Field(i)

		// Skip unexported or unsettable fields
		if !srcField.CanInterface() || !dstField.CanSet() {
			continue
		}

		// Embedded struct (anonymous) — recursive call
		if fieldType.Anonymous && srcField.Kind() == reflect.Struct {
			err := applyNonNilFieldsReflect(dstField, srcField)
			if err != nil {
				return err
			}
			continue
		}

		// Skip nil pointers
		if srcField.Kind() == reflect.Ptr && srcField.IsNil() {
			continue
		}

		// Set value
		dstField.Set(srcField)
	}

	return nil
}

func GetUserID(ctx context.Context) (uint, bool) {
	val := ctx.Value(UserIDKey)
	if uid, ok := val.(uint); ok {
		return uid, true
	}
	return 0, false
}

func toInt(s string) int {
	i, _ := strconv.Atoi(strings.ReplaceAll(s, ",", ""))
	return i
}

func GetTenantId(c *gin.Context) any {
	tenantID, ok := c.Get("tenant_id")
	if !ok {
		headerTenantID := c.GetHeader("Tenant-Id")
		tenantID = ParseUintParam(headerTenantID)
	}
	return tenantID
}
