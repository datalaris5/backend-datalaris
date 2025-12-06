package controllers

import (
	"fmt"
	"go-datalaris/config"
	"go-datalaris/constant"
	"go-datalaris/models"
	"go-datalaris/services"
	"go-datalaris/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateStore(c *gin.Context) {
	input, errBind := utils.BindJSON[models.Store](c)
	if errBind != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid input", errBind.Error())
		return
	}

	db := config.DB
	err := services.WithTransaction(db.WithContext(c.Request.Context()), func(tx *gorm.DB) error {
		_, err := services.GetWhereFirst[models.Store]("name = ?", input.Name)
		if err == nil {
			return fmt.Errorf(constant.StoreConst + constant.ErrorAlreadyExist)
		}
		savedInput, err := services.Save[models.Store](input, tx)
		if err != nil {
			return err
		}

		input = savedInput
		return nil
	})

	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Transaction failed", utils.ParseDBError(err))
		return
	}

	utils.Success(c, constant.StoreConst+constant.SuccessSaved, input)
}

func UpdateStore(c *gin.Context) {
	input, errBind := utils.BindJSON[models.Store](c)
	if errBind != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid input", errBind.Error())
		return
	}

	existing, err := services.GetById[models.Store](input.ID, c)
	if err == gorm.ErrRecordNotFound {
		utils.Error(c, http.StatusNotFound, constant.StoreConst+constant.ErrorNotFound, nil)
		return
	}

	db := config.DB
	err = services.WithTransaction(db.WithContext(c.Request.Context()), func(tx *gorm.DB) error {
		existing, err = services.Update[models.Store](input, existing, tx)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Transaction failed", utils.ParseDBError(err))
		return
	}

	utils.Success(c, constant.StoreConst+constant.SuccessUpdate, existing)
}

func GetStoreById(c *gin.Context) {
	id := c.Param("id")
	store, err := services.GetWhereFirst[models.Store]("id = ?", id)
	if err == gorm.ErrRecordNotFound {
		utils.Error(c, http.StatusNotFound, constant.StoreConst+constant.ErrorNotFound, nil)
		return
	}

	utils.Success(c, constant.StoreConst+constant.SuccessFetch, store)
}
