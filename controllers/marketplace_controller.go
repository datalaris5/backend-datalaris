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

func CreateMarketplace(c *gin.Context) {
	input, errBind := utils.BindJSON[models.Marketplace](c)
	if errBind != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid input", errBind.Error())
		return
	}

	db := config.DB
	err := services.WithTransaction(db.WithContext(c.Request.Context()), func(tx *gorm.DB) error {
		_, err := services.GetWhereFirst[models.Marketplace]("name = ?", input.Name)
		if err == nil {
			return fmt.Errorf(constant.MarketplaceConst + constant.ErrorAlreadyExist)
		}
		savedInput, err := services.Save[models.Marketplace](input, tx)
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

	utils.Success(c, constant.MarketplaceConst+constant.SuccessSaved, input)
}

func UpdateMarketplace(c *gin.Context) {
	input, errBind := utils.BindJSON[models.Marketplace](c)
	if errBind != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid input", errBind.Error())
		return
	}

	existing, err := services.GetById[models.Marketplace](input.ID, c)
	if err == gorm.ErrRecordNotFound {
		utils.Error(c, http.StatusNotFound, constant.MarketplaceConst+constant.ErrorNotFound, nil)
		return
	}

	db := config.DB
	err = services.WithTransaction(db.WithContext(c.Request.Context()), func(tx *gorm.DB) error {
		existing, err = services.Update[models.Marketplace](input, existing, tx)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Transaction failed", utils.ParseDBError(err))
		return
	}

	utils.Success(c, constant.MarketplaceConst+constant.SuccessUpdate, existing)
}

func GetMarketplaceById(c *gin.Context) {
	id := c.Param("id")
	Marketplace, err := services.GetWhereFirst[models.Marketplace]("id = ?", id)
	if err == gorm.ErrRecordNotFound {
		utils.Error(c, http.StatusNotFound, constant.MarketplaceConst+constant.ErrorNotFound, nil)
		return
	}

	utils.Success(c, constant.MarketplaceConst+constant.SuccessFetch, Marketplace)
}
