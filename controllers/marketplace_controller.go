package controllers

import (
	"context"
	"fmt"
	"go-datalaris/config"
	"go-datalaris/constant"
	"go-datalaris/dto"
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
		input.IsActive = true
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
	marketplace, err := services.GetWhereFirst[models.Marketplace]("id = ?", id)
	if err == gorm.ErrRecordNotFound {
		utils.Error(c, http.StatusNotFound, constant.MarketplaceConst+constant.ErrorNotFound, nil)
		return
	}

	utils.Success(c, constant.MarketplaceConst+constant.SuccessFetch, marketplace)
}

func GetMarketplaceByPage(c *gin.Context) {
	input, err := utils.BindJSON[dto.PaginationRequest](c)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	db := config.DB.Model(&models.Marketplace{})
	if input.Search != nil && *input.Search != "" {
		searchQuery := "%" + *input.Search + "%"
		db = db.Where("name ILIKE ?", searchQuery)
	}

	if input.Status != nil && *input.Status != "" {
		var status bool
		if *input.Status == constant.ActiveConst {
			status = true
		} else if *input.Status == constant.InactiveConst {
			status = false
		}
		db = db.Where("is_active = ?", status)
	}

	items, total, err := services.PaginateWhere[models.Marketplace](c, db, input, "is_deleted = ?", false)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "", err.Error())
		return
	}

	utils.Success(c, constant.MarketplaceConst+constant.SuccessFetch, utils.PaginationResponse(items, total, input.Page, input.Limit))
}

func SoftDeleteMarketplace(c *gin.Context) {
	id := c.Param("id")
	userID, ok := c.Get("user_id")
	if !ok {
		utils.Error(c, http.StatusUnauthorized, "User ID not found", nil)
		return
	}
	ctx := context.WithValue(c.Request.Context(), utils.UserIDKey, userID)

	db := config.DB
	err := services.WithTransaction(db.WithContext(ctx), func(tx *gorm.DB) error {
		var existing models.Marketplace
		existing, err := services.GetById[models.Marketplace](utils.ParseUintParam(id), c)
		if err != nil {
			return err
		}

		existing.IsDeleted = true
		existing.IsActive = false
		_, err = services.Update[models.Marketplace](existing, existing, tx)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Transaction failed", utils.ParseDBError(err))
		return
	}
	utils.Success(c, constant.MarketplaceConst+constant.SuccessDelete, nil)

}

func ActiveInactiveMarketplace(c *gin.Context) {
	id := c.Param("id")
	status := c.Param("status")
	userID, ok := c.Get("user_id")
	if !ok {
		utils.Error(c, http.StatusUnauthorized, "User ID not found", nil)
		return
	}
	ctx := context.WithValue(c.Request.Context(), utils.UserIDKey, userID)

	db := config.DB
	err := services.WithTransaction(db.WithContext(ctx), func(tx *gorm.DB) error {
		var existing models.Marketplace
		existing, err := services.GetById[models.Marketplace](utils.ParseUintParam(id), c)
		if err != nil {
			return err
		}

		if status == "active" {
			existing.IsActive = true
		} else if status == "inactive" {
			existing.IsActive = false
		}

		_, err = services.Update[models.Marketplace](existing, existing, tx)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Transaction failed", utils.ParseDBError(err))
		return
	}
	utils.Success(c, constant.MarketplaceConst+utils.TernaryString(status == "active", func() string { return constant.SuccessActive }, func() string { return constant.SuccessInactive }), nil)

}

func GetLovMarketplace(c *gin.Context) {
	item, _ := services.GetWhereFind[models.Marketplace]("is_active = ? AND is_deleted = ?", true, false)

	var lovs []dto.ResponseLov
	for i := range item {
		lovs = append(lovs, dto.ResponseLov{
			Id:    item[i].ID,
			Value: item[i].Name,
		})
	}

	utils.Success(c, constant.MarketplaceConst+constant.SuccessFetch, lovs)
}
