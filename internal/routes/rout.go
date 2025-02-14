package routes

import (
	"avantura/backend/internal/controllers"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	app.Post("/api/register", controllers.Register)
	app.Post("/api/login", controllers.Login)
	app.Post("/api/logout", controllers.Logout)
	app.Post("/api/follow", controllers.Follow)
	app.Post("/api/unfollow", controllers.Unfollow)
	app.Post("/api/add-game", controllers.AddGame)
	app.Post("/api/add-event/:id", controllers.AddEvent)
	app.Post("/api/join-to-event/:id", controllers.JoinToEvent)
	app.Post("/api/add-news/:id", controllers.AddNews)
	app.Post("/api/add-comment", controllers.AddComment)
	app.Post("/api/record-discord", controllers.RecordDiscord)
	app.Post("/api/upload-avatar/:id", controllers.UploadAvatar)
	app.Post("/api/delete-game/:id", controllers.DeleteGame)
	app.Post("/api/unjoin-from-event/:id", controllers.UnjoinFromEvent)
	app.Post("/api/delete-avatar/:id", controllers.DeleteAvatar)
	app.Post("/api/delete-notification/:id", controllers.DeleteNotification)
	app.Post("/api/delete-all-notifications/:id", controllers.DeleteAllNotifications)
	app.Patch("/api/edit-rating/:id", controllers.EditRating)
	app.Get("/api/get-notifications/:id", controllers.GetNotifications)
	app.Get("/api/user", controllers.User)	
	app.Get("/api/events/:id", controllers.GetConcreteEvent)
	app.Get("/api/events", controllers.GetEvents)	
	app.Get("/api/news", controllers.GetNews)
	app.Get("/api/news/:id", controllers.GetConcreteNews)
	app.Get("/api/users", controllers.GetAllUsers)
	app.Get("/api/users/:id", controllers.GetConcreteUser)
	app.Get("/api/games", controllers.GetAllGames)
	app.Get("/api/games/:name", controllers.GetConcreteGame)
	app.Get("/api/comments/:id", controllers.ShowComments)
}