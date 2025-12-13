package seed

import (
	"go-datalaris/models"
	"go-datalaris/utils"
	"log"

	"gorm.io/gorm"
)

type MenuItem struct {
	Path string
	Name string
	Crud bool
}

func Seed(db *gorm.DB) {
	// --- TENANTS ---
	var tenantTokoA, tenantTokoB models.Tenant

	if err := db.Where("name = ?", "Toko A").First(&tenantTokoA).Error; err == gorm.ErrRecordNotFound {
		tenantTokoA = models.Tenant{
			Name:        "Toko A",
			Description: "Toko A",
			TenantKey:   "toko-a-123",
		}
		if err := db.Create(&tenantTokoA).Error; err != nil {
			log.Fatalf("❌ Failed to create tenant Toko A: %v", err)
		}
		log.Println("✅ Tenant Toko A created")
	}

	if err := db.Where("name = ?", "Toko B").First(&tenantTokoB).Error; err == gorm.ErrRecordNotFound {
		tenantTokoB = models.Tenant{
			Name:        "Toko B",
			Description: "TOko B",
			TenantKey:   "toko-b-123",
		}
		if err := db.Create(&tenantTokoB).Error; err != nil {
			log.Fatalf("❌ Failed to create tenant Toko B: %v", err)
		}
		log.Println("✅ Tenant Toko B created")
	}

	var tenantSuperTokoA models.Role
	if err := db.Where("name = ?", "superadmin_toko_a").
		First(&tenantSuperTokoA).Error; err == gorm.ErrRecordNotFound {

		tenantSuperTokoA = models.Role{
			Name: "superadmin_toko_a",
		}
		if err := db.Create(&tenantSuperTokoA).Error; err != nil {
			log.Fatalf("❌ Failed to create role superadmin_toko_a: %v", err)
		}
		log.Println("✅ Role superadmin_toko_a created")
	}

	var tenantSuperTokoB models.Role
	if err := db.Where("name = ?", "superadmin_toko_b").
		First(&tenantSuperTokoB).Error; err == gorm.ErrRecordNotFound {

		tenantSuperTokoB = models.Role{
			Name: "superadmin_toko_b",
		}
		if err := db.Create(&tenantSuperTokoB).Error; err != nil {
			log.Fatalf("❌ Failed to create role superadmin_toko_b: %v", err)
		}
		log.Println("✅ Role superadmin_toko_b created")
	}

	var roleAdminTokoA models.Role
	if err := db.Where("name = ?", "admin_toko_a").
		First(&roleAdminTokoA).Error; err == gorm.ErrRecordNotFound {

		roleAdminTokoA = models.Role{
			Name: "admin_toko_a",
		}
		if err := db.Create(&roleAdminTokoA).Error; err != nil {
			log.Fatalf("❌ Failed to create role admin_toko_a: %v", err)
		}
		log.Println("✅ Role admin_sanur created")
	}

	var roleAdminTokoB models.Role
	if err := db.Where("name = ?", "admin_toko_b").
		First(&roleAdminTokoB).Error; err == gorm.ErrRecordNotFound {

		roleAdminTokoB = models.Role{
			Name: "admin_toko_b",
		}
		if err := db.Create(&roleAdminTokoB).Error; err != nil {
			log.Fatalf("❌ Failed to create role admin_toko_b: %v", err)
		}
		log.Println("✅ Role admin_toko_b created")
	}

	// --- USERS ---
	createUser := func(email, name, password string, tenantID *uint, roleID uint) {
		var user models.User
		if err := db.Where("email = ?", email).First(&user).Error; err == gorm.ErrRecordNotFound {
			hash, _ := utils.HashPassword(password)
			user = models.User{
				Name:     name,
				Email:    email,
				Password: hash,
				TenantID: tenantID,
				RoleID:   &roleID,
			}
			if err := db.Create(&user).Error; err != nil {
				log.Fatalf("❌ Failed to create user %s: %v", email, err)
			}
			log.Printf("✅ User %s created\n", email)
		}
	}

	createUser("superadmin@tokoa.com", "Super Admin Toko A", "P@ssw0rd", &tenantTokoB.ID, tenantSuperTokoB.ID)
	createUser("superadmin@tokob.com", "Super Admin Toko B", "P@ssw0rd", &tenantTokoA.ID, tenantSuperTokoA.ID)
	createUser("admin@tokoa.com", "Admin Toko A", "P@ssw0rd", &tenantTokoA.ID, roleAdminTokoA.ID)
	createUser("admin@tokob.com", "Admin Toko B", "P@ssw0rd", &tenantTokoB.ID, roleAdminTokoB.ID)

}
