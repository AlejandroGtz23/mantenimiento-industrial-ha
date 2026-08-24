package router

import (
	"fiber3-app/backend/handler"
	"fiber3-app/backend/middleware"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/etag"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/session"
)

// SetupRoutes exposes the industrial maintenance API.
func SetupRoutes(app *fiber.App, _ *session.Store) {
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}))
	app.Use(logger.New())
	app.Use(compress.New())
	app.Use(etag.New())

	api := app.Group("/api")
	maintenance := handler.NewMaintenanceHandler()

	mobile := api.Group("/movil")
	mobile.Post("/login", maintenance.MobileLogin)
	mobile.Use(middleware.RequireJWT)
	mobile.Get("/maquinas", maintenance.Machines)
	mobile.Get("/maquinas/:id/checklist", maintenance.Checklist)
	mobile.Post("/registrar-mantenimiento", maintenance.Register)
	mobile.Get("/historial", maintenance.History)

	admin := api.Group("/admin")
	admin.Post("/iniciar-sesion", maintenance.AdminLogin)
	admin.Use(middleware.RequireAdminJWT)
	admin.Get("/dashboard", maintenance.Dashboard)
	admin.Get("/mantenimientos", maintenance.Records)
	admin.Get("/mantenimientos/:id", maintenance.Record)
	admin.Get("/maquinas", maintenance.MachinesAdmin)
	admin.Post("/maquinas", maintenance.CreateMachine)
	admin.Put("/maquinas/:id", maintenance.UpdateMachine)
	admin.Delete("/maquinas/:id", maintenance.DeleteMachine)
	admin.Get("/maquinas/:id/checklist", maintenance.ChecklistAdmin)
	admin.Post("/maquinas/:id/checklist", maintenance.CreateChecklistItem)
	admin.Put("/maquinas/:id/checklist/:itemId", maintenance.UpdateChecklistItem)
	admin.Delete("/maquinas/:id/checklist/:itemId", maintenance.DeleteChecklistItem)
	admin.Get("/tecnicos", maintenance.Technicians)
	admin.Post("/tecnicos", maintenance.CreateTechnician)
}
