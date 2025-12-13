package middlewares

import (
	"context"
	"fmt"
	"go-datalaris/config"
	"go-datalaris/models"
	"go-datalaris/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

func AuthMiddlewareWithDB() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("🟢 [AuthMiddlewareWithAuthDB] Start")

		// --- Ambil Header Authorization ---
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			return
		}

		tokenStr := parts[1]
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return utils.GetKey(), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			return
		}

		// Simpan claims mentah ke context
		c.Set("claims", claims)

		// --- Ambil user_id & tenant_id dari token ---
		var userID uint
		var tenantID uint

		if v, ok := claims["user_id"].(float64); ok {
			userID = uint(v)
			c.Set("user_id", userID)
		}
		if v, ok := claims["tenant_id"].(float64); ok {
			tenantID = uint(v)
			c.Set("tenant_id", tenantID)
		}

		// --- Ambil user & role dari AUTH DB ---
		var user models.User
		if err := config.DB.Preload("Role").First(&user, userID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "database error"})
			return
		}

		tenantHeader := c.GetHeader("Tenant-Id")
		if tenantHeader != "" {
			c.Set("tenant_id_header", tenantHeader)
		}

		c.Set("role_id", user.Role.ID)
		c.Set("role_name", user.Role.Name)

		// --- Inject ke context request (untuk downstream service) ---
		ctx := context.WithValue(c.Request.Context(), "claims", claims)
		ctx = context.WithValue(ctx, "user_id", userID)
		ctx = context.WithValue(ctx, "role_id", user.Role.ID)
		ctx = context.WithValue(ctx, "tenant_id", tenantID)
		ctx = context.WithValue(ctx, "tenant_id_header", tenantHeader)
		c.Request = c.Request.WithContext(ctx)

		// fmt.Printf("✅ Authenticated user_id=%d, role=%s, tenant_key=%s\n", user.ID, user.Role.RoleType)
		c.Next()
	}
}

func InjectUserToContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("🟢 [InjectUserToContext] Start")
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			fmt.Println("🔴 No Authorization header in InjectUserToContext")
			c.Next()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ParseJWT(tokenString)
		if err != nil {
			fmt.Println("🔴 ParseJWT error:", err)
			c.Next()
			return
		}

		fmt.Printf("🧩 Inject claims: %+v\n", claims)

		userID, ok := claims["user_id"].(float64)
		if ok {
			uid := uint(userID)
			fmt.Println("✅ Injected user_id:", uid)

			c.Set(string(models.UserIDKey), uid)

			ctx := context.WithValue(c.Request.Context(), models.UserIDKey, uid)
			c.Request = c.Request.WithContext(ctx)
		} else {
			fmt.Println("🔴 user_id not found in InjectUserToContext claims")
		}

		c.Next()
	}
}
