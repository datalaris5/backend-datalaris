package services

import (
	"go-datalaris/config"
	"go-datalaris/dto"
	"go-datalaris/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetById generic function to get a record by id
func GetById[T any](id uint, c *gin.Context, preloads ...string) (item T, err error) {
	db := config.DB
	for _, p := range preloads {
		db = db.Preload(p)
	}
	err = db.First(&item, id).Error
	return
}

func GetByType[T any](typeParam string, c *gin.Context, preloads ...string) (item T, err error) {
	db := config.DB
	for _, p := range preloads {
		db = db.Preload(p)
	}
	err = db.Where("type = ?", typeParam).First(&item).Error
	return
}

func GetBySlug[T any](slug string, c *gin.Context, preloads ...string) (item T, err error) {
	db := config.DB
	for _, p := range preloads {
		db = db.Preload(p)
	}
	err = db.Where("slug = ?", slug).First(&item).Error
	return
}

func GetByKeyName[T any](keyName string, c *gin.Context) (item T, err error) {
	db := config.DB
	err = db.Where("key_name = ?", keyName).First(&item).Error
	return
}

// GetById generic function to get a record by id
func GetFirst[T any](c *gin.Context) (item T, err error) {
	err = config.DB.First(&item).Error
	return
}

func GetList[T any](c *gin.Context, preloads ...string) (item []T, err error) {
	db := config.DB
	for _, p := range preloads {
		db = db.Preload(p)
	}
	err = db.Find(&item).Error
	return
}

// Save generic function to create a record
func Save[T any](input T, tx *gorm.DB) (item T, err error) {
	input = utils.NormalizeStringPointers(input)
	if err = tx.Create(&input).Error; err != nil {
		return item, err
	}
	return input, nil
}

// Update generic function to update a record
func Update[T any](input T, existing T, tx *gorm.DB) (item T, err error) {
	input = utils.NormalizeStringPointers(input)
	err = utils.ApplyNonNilFields(&existing, input)
	if err != nil {
		return item, err
	}
	if err = tx.Omit("created_by", "created_at", "email", "password").Save(&existing).Error; err != nil {
		return item, err
	}

	return existing, nil
}

// Delete generic function to delete a record by id
func Delete[T any](id uint, tx *gorm.DB) error {
	return tx.Where("id = ?", id).Delete(new(T)).Error
}

func DeleteWhere[T any](tx *gorm.DB, query string, args ...interface{}) error {
	return tx.Where(query, args...).Delete(new(T)).Error
}

func GetWhereFind[T any](query string, args ...interface{}) (item []T, err error) {
	db := config.DB
	err = db.Where(query, args...).Find(&item).Error
	return
}

func GetWhereFirst[T any](query string, args ...interface{}) (item T, err error) {
	db := config.DB
	err = db.Where(query, args...).First(&item).Error
	return
}

func GetCount[T any]() (result int64, err error) {
	db := config.DB
	var total int64
	err = db.Count(&total).Error
	return
}

func Paginate[T any](c *gin.Context, db *gorm.DB, input dto.PaginationRequest, orderBy ...string) ([]T, int64, error) {
	// Normalisasi nilai default
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.Limit <= 0 {
		input.Limit = 10
	}
	offset := (input.Page - 1) * input.Limit

	var result []T
	var total int64

	// Hitung total
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	order := "created_at desc"
	if len(orderBy) > 0 && orderBy[0] != "" {
		order = orderBy[0]
	}

	// Fetch data dengan pagination
	if err := db.
		Limit(input.Limit).
		Offset(offset).
		Order(order).
		Find(&result).Error; err != nil {
		return nil, 0, err
	}

	return result, total, nil
}

func PaginateWhere[T any](c *gin.Context, db *gorm.DB, input dto.PaginationRequest, query string, args ...interface{}) ([]T, int64, error) {
	// Normalisasi nilai default
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.Limit <= 0 {
		input.Limit = 10
	}
	offset := (input.Page - 1) * input.Limit

	var result []T
	var total int64

	// Hitung total
	if err := db.Where(query, args...).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch data dengan pagination
	if err := db.
		Where(query, args...).
		Limit(input.Limit).
		Offset(offset).
		Order("updated_at desc").
		Find(&result).Error; err != nil {
		return nil, 0, err
	}

	return result, total, nil
}
