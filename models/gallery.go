package models

type Gallery struct {
	BaseModel
	GalleryableID   uint   `json:"galleryable_id" `
	GalleryableType string `json:"galleryable_type"`
	FilePath        string `json:"file_path"`
	FileType        string `json:"file_type"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	IsUsed          bool   `json:"is_used"`
	CtaText         string `json:"cta_text"`
	CtaUrl          string `json:"cta_url"`
	Order           int    `json:"order" gorm:"column:order_"`
	IsThumbnail     bool   `json:"is_thumbnail"`
	IsLogo          bool   `json:"is_logo"`
	IsActive        bool   `gorm:"type:bool;default:true" json:"is_active"`
	IsDeleted       bool   `gorm:"type:bool;default:false" json:"is_deleted"`
}
