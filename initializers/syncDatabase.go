package initializers

import "my-go/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
}
