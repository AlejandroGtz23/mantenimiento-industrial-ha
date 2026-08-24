package main

import (
	"fiber3-app/backend/config"
	"fiber3-app/backend/database"
	"fiber3-app/backend/middleware"
	"fiber3-app/backend/router"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/gofiber/fiber/v3/middleware/static"
)

func main() {
	// Load config Env variables
	config.LoadConfig()

	// DB Connection
	database.ConnectDB()
	database.ConnectRedis()
	if err := database.EnsureMaintenanceSchema(); err != nil {
		log.Fatal("Error al preparar tablas de mantenimiento: ", err)
	}
	if err := os.MkdirAll(config.Env.EvidenceDirectory, 0750); err != nil {
		log.Fatal("Error al preparar directorio de evidencias: ", err)
	}

	// Session in memory by default
	store := session.NewStore(session.Config{
		IdleTimeout:     30 * time.Minute,
		AbsoluteTimeout: 24 * time.Hour,
	})
	middleware.InitSessionStore(store)
	// Fiber Web server
	app := fiber.New()
	router.SetupRoutes(app, store)

	//
	app.Use(devSecurity())

	// Servir archivos de dist generado por Vite
	app.Use("/uploads", static.New("./uploads"))
	app.Use("/", static.New("./internal/web/dist"))

	// SPA fallback (para Vue/React/etc)
	app.Get("*", func(c fiber.Ctx) error {
		return c.SendFile("./internal/web/dist/index.html")
	})

	// Start the server with Prefork enabled via ListenConfig
	// fiber.ListenConfig{EnablePrefork: true}
	if err := app.Listen(config.Env.Port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

// En desarrollo Vite usa WebSockets y eval().
func devSecurity() fiber.Handler {
	return helmet.New(helmet.Config{
		CrossOriginResourcePolicy: "cross-origin",
		CrossOriginEmbedderPolicy: "unsafe-none",
		CrossOriginOpenerPolicy:   "same-origin",
		XFrameOptions:             "SAMEORIGIN",
		//XContentTypeOptions: "nosniff",
		ReferrerPolicy: "strict-origin-when-cross-origin",

		ContentSecurityPolicy: `
			default-src 'self' http: https: ws: wss: data: blob:;
			script-src 'self' 'unsafe-inline' 'unsafe-eval' http: https:;
			style-src 'self' 'unsafe-inline' https:;
			img-src 'self' data: blob: https:;
			font-src 'self' https: data:;
			connect-src 'self' http: https: ws: wss:;
			frame-ancestors 'self';
		`,
	})
}

// En produccion el acceso debe ser estricto
func productionSecurity() fiber.Handler {
	return helmet.New(helmet.Config{
		CrossOriginResourcePolicy: "same-origin",
		CrossOriginEmbedderPolicy: "unsafe-none",
		CrossOriginOpenerPolicy:   "same-origin",
		XFrameOptions:             "SAMEORIGIN",
		//XContentTypeOptions: "nosniff",
		XSSProtection:  "1; mode=block",
		ReferrerPolicy: "strict-origin-when-cross-origin",

		ContentSecurityPolicy: `
			default-src 'self';
			script-src 'self';
			style-src 'self' 'unsafe-inline' https:;
			img-src 'self' data: https:;
			font-src 'self' https: data:;
			connect-src 'self' https: wss:;
			frame-ancestors 'self';
			base-uri 'self';
			form-action 'self';
		`,
	})
}
