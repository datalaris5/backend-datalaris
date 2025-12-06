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

func CreateTenant(c *gin.Context) {
	input, errBind := utils.BindJSON[models.Tenant](c)
	if errBind != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid input", errBind.Error())
		return
	}

	db := config.DB
	err := services.WithTransaction(db.WithContext(c.Request.Context()), func(tx *gorm.DB) error {
		_, err := services.GetWhereFirst[models.Tenant]("name = ?", input.Name)
		if err == nil {
			return fmt.Errorf(constant.TenantConst + constant.ErrorAlreadyExist)
		}
		savedInput, err := services.Save[models.Tenant](input, tx)
		if err != nil {
			return err
		}

		input = savedInput

		role := models.Role{
			Name:     "SUPERADMIN " + savedInput.Name,
			RoleType: constant.RoleType.TenantAdmin,
		}

		role, err = services.Save[models.Role](role, tx)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Transaction failed", utils.ParseDBError(err))
		return
	}

	utils.Success(c, constant.TenantConst+constant.SuccessSaved, input)
}

func UpdateTenant(c *gin.Context) {
	input, errBind := utils.BindJSON[models.Tenant](c)
	if errBind != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid input", errBind.Error())
		return
	}

	existing, err := services.GetById[models.Tenant](input.ID, c)
	if err == gorm.ErrRecordNotFound {
		utils.Error(c, http.StatusNotFound, constant.TenantConst+constant.ErrorNotFound, nil)
		return
	}

	db := config.DB
	err = services.WithTransaction(db.WithContext(c.Request.Context()), func(tx *gorm.DB) error {
		existing, err = services.Update[models.Tenant](input, existing, tx)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Transaction failed", utils.ParseDBError(err))
		return
	}

	utils.Success(c, constant.TenantConst+constant.SuccessUpdate, existing)
}

func GetTenantById(c *gin.Context) {
	id := c.Param("id")
	tenant, err := services.GetWhereFirst[models.Tenant]("id = ?", id)
	if err == gorm.ErrRecordNotFound {
		utils.Error(c, http.StatusNotFound, constant.TenantConst+constant.ErrorNotFound, nil)
		return
	}

	utils.Success(c, constant.TenantConst+constant.SuccessFetch, tenant)
}
