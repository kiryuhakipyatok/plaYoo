package notify

import (
	"avantura/backend/internal/db"
	"avantura/backend/internal/models"
	"log"
	"strconv"
	"time"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
)

func CreateNotice(event models.Event,msg string) models.Notice{
	notice:=models.Notice{
		Id:uuid.New(),
		EventId: event.Id,
		Body: msg,
	}
	if err:=db.Database.Create(&notice).Error;err!=nil{
		log.Printf("Error to create notice")
	}
	return notice
}

func sendStartNotification(event models.Event,message string){
	for _,id:=range event.Members{
		user:=models.User{}
		if err:=db.Database.Find(&user,"id=?",id).Error;err!=nil{
			log.Printf("User not found: %s", user.Id)
		}
		notice:=CreateNotice(event,message)
		user.Notifications = append(user.Notifications, notice.Id.String())
		db.Database.Save(&user)
		if user.ChatId!=""{
			chatID, _ := strconv.ParseInt(user.ChatId, 10, 64)
			msg:=tgbotapi.NewMessage(chatID,message)
			if _,err := Bot.Send(msg); err != nil {
			log.Printf("Failed to send message to user %s: %v", user.Tg, err)
		} else {
			log.Printf("Notification sent to user %s", user.Tg)
		}
		}
	}
}

func sendPreNotification(event models.Event,message string){
	// message:="Событие " + event.Body + " начнется через 10 минут!"
	for _,id:=range event.Members{
		user:=models.User{}
		if err:=db.Database.Find(&user,"id=?",id).Error;err!=nil{
			log.Printf("User not found: %s", user.Id)
		}
		notice:=CreateNotice(event,message)
		user.Notifications = append(user.Notifications, notice.Id.String())
		db.Database.Save(&user)
		log.Printf("Notification save to user %s", user.Tg)
		if user.ChatId!=""{
			chatID, _ := strconv.ParseInt(user.ChatId, 10, 64)
			msg:=tgbotapi.NewMessage(chatID,message)
			if _,err := Bot.Send(msg); err != nil {
			log.Printf("Failed to send message to user %s: %v", user.Tg, err)
			} else {
			log.Printf("Notification sent to user %s", user.Tg)
			}
		}
	}
}

func ScheduleNotify() {
    c := cron.New()
    c.AddFunc("@every 1m", func() {
        now := time.Now()
        tenMin := now.Add(10 * time.Minute).Add(30 * time.Second)
        var upcomingEvents []models.Event
        if err := db.Database.Where("time <= ?", tenMin).Find(&upcomingEvents).Error; err != nil {
            log.Printf("Ошибка при получении предстоящих событий: %v", err)
            return
        }
        for _, event := range upcomingEvents {
			if !event.NotifiedPre{
				premsg:="Событие " + event.Body + " начнется через 10 минут!"
				sendPreNotification(event,premsg)
				log.Printf("Уведомление о предстоящем событии %v отправлено в %v", event.Body, time.Now())
				event.NotifiedPre = true
				if err:=db.Database.Save(&event).Error;err!=nil{
					log.Printf("Ошибка при получении предстоящих событий: %v", err)
					return
				}
			}
        }

        var events []models.Event
        if err := db.Database.Where("time <= ?", now.Add(1 * time.Minute).Add(30*time.Second)).Find(&events).Error; err != nil {
            log.Printf("Ошибка при получении текущих событий: %v", err)
            return
        }
		
        for _, event := range events {
			startmsg:="Событие " + event.Body + " началось!"
			sendStartNotification(event,startmsg)
				log.Printf("Уведомление о начале события %v отправлено в %v", event.Body, time.Now())
				
            	if err := db.Database.Delete(&event).Error; err != nil {
                log.Printf("Ошибка при удалении события: %v", err)
            	}
				for _,id:= range event.Members{
					user:=models.User{}
					// uuId,_:=uuid.Parse(id)
					if err:=db.Database.First(&user,"id=?",id).Error;err!=nil{
						log.Printf("Ошибка при поиска пользователя: %v", err)
					}
					updateEvents:=make([]string,0,len(user.Events))
					for _, e := range user.Events {
						if e != event.Id.String() {
							updateEvents = append(updateEvents, e)
						}
					}
					user.Events = updateEvents
					db.Database.Save(&user)
				}
        }
    })
    c.Start()
}