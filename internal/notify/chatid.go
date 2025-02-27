package notify

import(
	"log"
	"strconv"
	"avantura/backend/internal/db/postgres"
	"avantura/backend/internal/models"
)


func storeChatID(user *models.User, chatID int64) {
	mu.Lock()
	defer mu.Unlock()
	user.ChatId = strconv.Itoa(int(chatID))
	if err := postgres.Database.Save(&user).Error; err != nil {
		log.Printf("Error saving user chat ID: %v", err)
	}
}

func removeChatID(user *models.User) {
	mu.Lock()
	defer mu.Unlock()
	user.ChatId = ""
	if err := postgres.Database.Save(&user).Error; err != nil {
		log.Printf("Error removing user chat ID: %v", err)
	}
}