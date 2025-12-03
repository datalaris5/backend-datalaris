package constant

type RoleTypeStruct struct {
	Global      string
	TenantSuper string
	TenantAdmin string
}

var RoleType = RoleTypeStruct{
	Global:      "GLOBAL_SUPERADMIN",
	TenantSuper: "TENANT_SUPERADMIN",
	TenantAdmin: "TENANT_ADMIN",
}

var AllowedRoleTypes = []string{
	RoleType.Global,
	RoleType.TenantSuper,
	RoleType.TenantAdmin,
}

func IsValidRoleType(roleType string) bool {
	for _, r := range AllowedRoleTypes {
		if r == roleType {
			return true
		}
	}
	return false
}

// Success
const SuccessFetch = " Fetched Successfully"
const SuccessSaved = " Saved Successfully"
const SuccessDelete = " Delete Successfully"
const SuccessUpdate = " Update Successfully"
const SuccessUploadFile = "Upload File Successfully"

// Error
const ErrorNotFound = " Not Found"
const ErrorAlreadyExist = " Already Exist"
const ErrorFailedStartTransaction = "Failed To Start Transaction"
const ErrorTransactionPanic = "Transaction Panic"
const ErrorTransactionFailed = "Transaction Failed"
const ErrorFailedCommitTransaction = "Failed To Commit Transaction"
const ErrorUploadFile = "Failed To Upload File"
const ErrorCopyFile = "Failed To Copy File"
const ErrorSaveFile = "Failed To Save File"
const ErrorCreateDir = "Failed To Create Directory"
const ErrorOpenFile = "Failed To Open File"
const ErrorInvalidFileType = "File Type Invalid"
const OnlyGlobalSuperadminCanAccess = "Only global superadmin can access"
const OnlyGlobalSuperadminCanCreate = "Only global superadmin can create"

// Modul
const TenantConst = "Tenant"
const RoleConst = "Role"
const UserConst = "User"
const LovConst = "Lov"

// Status
const ActiveConst = "active"
const InactiveConst = "inactive"
