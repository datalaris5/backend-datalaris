package controllers

import (
	"go-datalaris/constant"
	"go-datalaris/models"
	"go-datalaris/services"
	"go-datalaris/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

type loginReq struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RoleMenuResponse struct {
	MenuID    uint   `json:"menu_id"`
	Name      string `json:"name"`
	Path      string `json:"path"`
	CanCreate bool   `json:"can_create"`
	CanRead   bool   `json:"can_read"`
	CanUpdate bool   `json:"can_update"`
	CanDelete bool   `json:"can_delete"`
}

// ------------------ LOGIN ------------------
func Login(c *gin.Context) {
	input, errBind := utils.BindJSON[loginReq](c)
	if errBind != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid input", errBind.Error())
		return
	}

	user, err := services.GetWhereFirst[models.User]("email = ?", input.Email)
	if err == gorm.ErrRecordNotFound {
		utils.Error(c, http.StatusNotFound, constant.TenantConst+constant.ErrorNotFound, nil)
		return
	}

	if !utils.CheckPasswordHash(input.Password, user.Password) {
		utils.Error(c, http.StatusUnauthorized, "Invalid credentials", "Incorrect password")
		return
	}

	if !user.IsActive {
		utils.Error(c, http.StatusForbidden, "Account inactive", "Your account is currently inactive. Please contact the administrator.")
		return
	}

	role, err := services.GetById[models.Role](*user.RoleID, c)
	if err == gorm.ErrRecordNotFound {
		utils.Error(c, http.StatusNotFound, constant.TenantConst+constant.ErrorNotFound, nil)
		return
	}

	var tenantName, tenantKey string
	if user.TenantID != nil && *user.TenantID != 0 {
		tenant, _ := services.GetById[models.Tenant](*user.TenantID, c)
		tenantName = tenant.Name
		tenantKey = tenant.TenantKey
	}

	claims := jwt.MapClaims{
		"user_id":    user.ID,
		"tenant_id":  user.TenantID,
		"tenant_key": tenantKey,
		"role_id":    role.ID,
		"role":       role.Name,
		"exp":        time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString([]byte(utils.GetKey()))

	utils.Success(c, "Login successful", gin.H{
		"token": tokenStr,
		"user": gin.H{
			"id":          user.ID,
			"email":       user.Email,
			"tenant_id":   user.TenantID,
			"tenant_name": tenantName,
			"tenant_key":  tenantKey,
			"role": gin.H{
				"id":        role.ID,
				"name":      role.Name,
				"role_type": role.RoleType,
			},
		},
	})
}

// // ------------------ FORGOT PASSWORD ------------------
// func ForgotPassword(c *gin.Context, db *gorm.DB) {
// 	var input struct {
// 		Email string `json:"email"`
// 	}
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		utils.Error(c, http.StatusBadRequest, "Invalid input", err.Error())
// 		return
// 	}

// 	var user models.User
// 	if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
// 		utils.Error(c, http.StatusNotFound, "User not found", err.Error())
// 		return
// 	}

// 	err := services.WithTransaction(db.WithContext(c.Request.Context()), func(tx *gorm.DB) error {
// 		newPassword := generateRandomPassword(10)
// 		hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
// 		if err != nil {
// 			return err
// 		}

// 		user.Password = string(hashed)
// 		if err := tx.Save(&user).Error; err != nil {
// 			return err
// 		}

// 		data := dto.EmailDto{
// 			Title:       "Forgot your password ?",
// 			Email:       user.Email,
// 			NewPassword: newPassword,
// 		}

// 		subject := "Reset Password"
// 		body := utils.RenderEmailTemplate(data)
// 		// body := "Hello,\n\nYour new password is: " + newPassword + "\nPlease login and change your password immediately."
// 		if err := utils.SendEmail(user.Email, subject, body); err != nil {
// 			return errors.New("password reset successful but failed to send email")
// 		}

// 		return nil
// 	})

// 	if err != nil {
// 		utils.Error(c, http.StatusInternalServerError, constant.ErrorTransactionFailed, utils.ParseDBError(err))
// 		return
// 	}

// 	utils.Success(c, "Password reset successful and email sent", gin.H{"email": user.Email})
// }

// // ------------------ RESET PASSWORD (ADMIN) ------------------
// func ResetPassword(c *gin.Context, db *gorm.DB) {
// 	isGlobal, _ := c.Get("is_global_superadmin")
// 	isTenantSuper, _ := c.Get("is_tenant_superadmin")

// 	if !(isGlobal == true || isTenantSuper == true) {
// 		utils.Error(c, http.StatusForbidden, "Access denied", "Not allowed")
// 		return
// 	}

// 	var input struct {
// 		UserID uint `json:"user_id"`
// 	}
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		utils.Error(c, http.StatusBadRequest, "Invalid input", err.Error())
// 		return
// 	}

// 	err := services.WithTransaction(db.WithContext(c.Request.Context()), func(tx *gorm.DB) error {
// 		var user models.User
// 		if err := tx.First(&user, input.UserID).Error; err != nil {
// 			return err
// 		}

// 		newPassword := generateRandomPassword(10)
// 		hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
// 		if err != nil {
// 			return err
// 		}

// 		user.Password = string(hashed)
// 		if err := tx.Save(&user).Error; err != nil {
// 			return err
// 		}

// 		data := dto.EmailDto{
// 			Title:       "Password Reset",
// 			Email:       user.Email,
// 			NewPassword: newPassword,
// 		}

// 		subject := "Reset Password"
// 		body := utils.RenderEmailTemplate(data)
// 		// body := "Hello,\n\nYour new password is: " + newPassword + "\nPlease login and change your password immediately."
// 		if err := utils.SendEmail(user.Email, subject, body); err != nil {
// 			return errors.New("password reset successful but failed to send email")
// 		}

// 		utils.Success(c, "Password reset successfully", gin.H{
// 			"user_id":  user.ID,
// 			"email":    user.Email,
// 			"new_pass": newPassword,
// 		})
// 		return nil
// 	})

// 	if err != nil {
// 		utils.Error(c, http.StatusInternalServerError, constant.ErrorTransactionFailed, utils.ParseDBError(err))
// 	}
// }

// // ------------------ CHANGE PASSWORD ------------------
// func ChangePassword(c *gin.Context, db *gorm.DB) {
// 	uid, _ := c.Get("user_id")

// 	input, errBind := utils.BindJSON[struct {
// 		OldPassword string `json:"old_password"`
// 		NewPassword string `json:"new_password"`
// 	}](c)
// 	if errBind != nil {
// 		utils.Error(c, http.StatusBadRequest, "Invalid input", errBind.Error())
// 		return
// 	}

// 	err := services.WithTransaction(db.WithContext(c.Request.Context()), func(tx *gorm.DB) error {
// 		var user models.User
// 		if err := tx.First(&user, uid).Error; err != nil {
// 			return err
// 		}

// 		if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.OldPassword)) != nil {
// 			return errors.New("wrong old password")
// 		}

// 		hashed, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
// 		if err != nil {
// 			return err
// 		}

// 		user.Password = string(hashed)
// 		return tx.Save(&user).Error
// 	})

// 	if err != nil {
// 		utils.Error(c, http.StatusInternalServerError, constant.ErrorTransactionFailed, utils.ParseDBError(err))
// 		return
// 	}

// 	utils.Success(c, "Password changed successfully", nil)
// }

// // ------------------ ACTIVATE / DEACTIVATE USER ------------------
// func ToggleUserActive(c *gin.Context, db *gorm.DB) {
// 	isGlobal, _ := c.Get("is_global_superadmin")
// 	isTenantSuper, _ := c.Get("is_tenant_superadmin")

// 	if !(isGlobal == true || isTenantSuper == true) {
// 		utils.Error(c, http.StatusForbidden, "Access denied", "Not allowed")
// 		return
// 	}

// 	input, errBind := utils.BindJSON[struct {
// 		UserID   uint `json:"user_id"`
// 		IsActive bool `json:"is_active"`
// 	}](c)
// 	if errBind != nil {
// 		utils.Error(c, http.StatusBadRequest, "Invalid input", errBind.Error())
// 		return
// 	}

// 	err := services.WithTransaction(db.WithContext(c.Request.Context()), func(tx *gorm.DB) error {
// 		var user models.User
// 		if err := tx.First(&user, input.UserID).Error; err != nil {
// 			return err
// 		}

// 		user.IsActive = input.IsActive
// 		if err := tx.Save(&user).Error; err != nil {
// 			return err
// 		}

// 		status := "User deactivated"
// 		if user.IsActive {
// 			status = "User activated"
// 		}

// 		utils.Success(c, status, gin.H{
// 			"user_id":   user.ID,
// 			"email":     user.Email,
// 			"is_active": user.IsActive,
// 		})
// 		return nil
// 	})

// 	if err != nil {
// 		utils.Error(c, http.StatusInternalServerError, constant.ErrorTransactionFailed, utils.ParseDBError(err))
// 	}
// }

// // ------------------ HELPER ------------------
// func generateRandomPassword(length int) string {
// 	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
// 	rand.Seed(time.Now().UnixNano())
// 	b := make([]byte, length)
// 	for i := range b {
// 		b[i] = charset[rand.Intn(len(charset))]
// 	}
// 	return string(b)
// }
