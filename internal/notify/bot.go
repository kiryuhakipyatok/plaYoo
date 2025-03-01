package notify

import (
	"avantura/backend/internal/db/postgres"
	"avantura/backend/internal/models"
	"avantura/backend/pkg/constants"
	"log"
	//"os"
	//"os/signal"
	"sync"
	//"syscall"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	mu sync.Mutex
	Bot *tgbotapi.BotAPI
)

func CreateBot(stop chan struct{}) {
	var err error
	Bot, err = tgbotapi.NewBotAPI(constants.TelegramBotToken)
	if err != nil {
		log.Printf("Error creating bot: %v", err)
	}

	log.Printf("Authorized on bot %s", Bot.Self.UserName)
	listenForUpdates(stop)
	// var wg sync.WaitGroup
	// c:=make(chan struct{})
	// wg.Add(1)
	// go func ()  {
	// go	listenForUpdates()
	// 	defer wg.Done()
	// }()
	// quit:=make(chan os.Signal,1)
	// signal.Notify(quit,syscall.SIGINT,syscall.SIGTERM)
	// <-quit
	// close(c)
	// wg.Wait()
}

func listenForUpdates(stop chan struct{}) {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates := Bot.GetUpdatesChan(updateConfig)
	for {
		select{
		case <-stop:
			log.Println("Stopping Bot...")
			return
		case update:=<-updates:
			if update.Message != nil {
				handleMessage(update)
			}
		}
	}

}

func handleMessage(update tgbotapi.Update) {
	username := update.Message.From.UserName
	chatID := update.Message.Chat.ID
	text := update.Message.Text

	user := models.User{}
	if err := postgres.Database.Find(&user, "tg=?", username).Error; err != nil {
		log.Printf("Error finding user: %v", err)
		return
	}

	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("✅ Да, хочу"),
			tgbotapi.NewKeyboardButton("❌ Нет, не хочу"),
		),
	)

	switch text {
	case "✅ Да, хочу":
		if user.ChatId == "" {
			storeChatID(&user, chatID)
			msg := tgbotapi.NewMessage(chatID, "Теперь вы будете получать уведомления от plaYoo о начале ивентов!")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			if _, err := Bot.Send(msg); err != nil {
				log.Printf("Failed to send message: %v", err)
			}
		} else {
			msg := tgbotapi.NewMessage(chatID, "Вы уже подписаны на уведомления.")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			if _, err := Bot.Send(msg); err != nil {
				log.Printf("Failed to send message: %v", err)
			}
		}

	case "❌ Нет, не хочу":
		if user.ChatId != "" {
			removeChatID(&user)
			msg := tgbotapi.NewMessage(chatID, "Вы отписались от уведомлений.")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			if _, err := Bot.Send(msg); err != nil {
				log.Printf("Failed to send message: %v", err)
			}
		} else {
			msg := tgbotapi.NewMessage(chatID, "Вы не подписаны на уведомления.")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			if _, err := Bot.Send(msg); err != nil {
				log.Printf("Failed to send message: %v", err)
			}
		}

	default:
		msg := tgbotapi.NewMessage(chatID, "Хотите ли вы получать уведомления о начале ивентов, к которым вы присоединились?")
		msg.ReplyMarkup = keyboard
		if _, err := Bot.Send(msg); err != nil {
			log.Printf("Failed to send message: %v", err)
		}
	}
}

