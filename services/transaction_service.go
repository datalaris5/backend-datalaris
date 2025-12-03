package services

import (
	"fmt"

	"gorm.io/gorm"
)

// WithTransaction menjalankan function txFunc di dalam transaksi,
// otomatis commit jika sukses, rollback jika error atau panic,
// dan kirim error response ke gin.Context kalau gagal.
func WithTransaction(db *gorm.DB, txFunc func(tx *gorm.DB) error) (err error) {
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	if err = txFunc(tx); err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
