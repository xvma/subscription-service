package main

import (
    "fmt"
    "net/http"
    "subscription-service/internal/config"
    "subscription-service/internal/database"
    "subscription-service/internal/handlers"
    "subscription-service/internal/middleware"
    "subscription-service/internal/repository"

    "github.com/gin-gonic/gin"
    "github.com/sirupsen/logrus"
)

func main() {
    // Load configuration
    cfg, err := config.LoadConfig()
    if err != nil {
        logrus.Fatal("Failed to load configuration: ", err)
    }

    // Initialize database
    db, err := database.NewPostgresDB(cfg.GetDSN())
    if err != nil {
        logrus.Fatal("Failed to connect to database: ", err)
    }
    defer db.Close()

    // Initialize repository and handler
    repo := repository.NewSubscriptionRepository(db.DB)
    handler := handlers.NewSubscriptionHandler(repo)

    // Setup router
    router := gin.Default()
    router.Use(middleware.LoggingMiddleware())

    // Swagger documentation - serve static files
    router.StaticFile("/swagger.json", "./docs/swagger.json")
    router.StaticFile("/swagger.yaml", "./docs/swagger.yaml")
    
    // Swagger UI handler
    router.GET("/swagger/", func(c *gin.Context) {
        c.Header("Content-Type", "text/html")
        c.String(http.StatusOK, `
        <!DOCTYPE html>
        <html lang="en">
        <head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <title>Subscription Service API Documentation</title>
            <link rel="stylesheet" type="text/css" href="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/5.0.0/swagger-ui.min.css">
            <style>
                html { box-sizing: border-box; overflow: -moz-scrollbars-vertical; overflow-y: scroll; }
                *, *:before, *:after { box-sizing: inherit; }
                body { margin: 0; background: #fafafa; }
            </style>
        </head>
        <body>
            <div id="swagger-ui"></div>
            <script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/5.0.0/swagger-ui-bundle.min.js"></script>
            <script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/5.0.0/swagger-ui-standalone-preset.min.js"></script>
            <script>
                window.onload = function() {
                    const ui = SwaggerUIBundle({
                        url: "/swagger.json",
                        dom_id: '#swagger-ui',
                        deepLinking: true,
                        presets: [
                            SwaggerUIBundle.presets.apis,
                            SwaggerUIStandalonePreset
                        ],
                        plugins: [
                            SwaggerUIBundle.plugins.DownloadUrl
                        ],
                        layout: "StandaloneLayout"
                    });
                    window.ui = ui;
                }
            </script>
        </body>
        </html>
        `)
    })

    // Also support /swagger/index.html for compatibility
    router.GET("/swagger/index.html", func(c *gin.Context) {
        c.Redirect(http.StatusMovedPermanently, "/swagger/")
    })

    // API routes
    api := router.Group("/api/v1")
    {
        api.GET("/", func(c *gin.Context) {
            c.JSON(200, gin.H{
                "message": "Subscription Service API",
                "version": "1.0",
                "documentation": "http://localhost:8080/swagger/",
            })
        })
        
        subscriptions := api.Group("/subscriptions")
        {
            subscriptions.POST("/", handler.CreateSubscription)
            subscriptions.GET("/:id", handler.GetSubscription)
            subscriptions.PUT("/:id", handler.UpdateSubscription)
            subscriptions.DELETE("/:id", handler.DeleteSubscription)
            subscriptions.POST("/total-cost", handler.GetTotalCost)
        }
    }

    // Health check
    router.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "healthy"})
    })

    // Start server
    addr := fmt.Sprintf(":%s", cfg.ServerPort)
    logrus.WithField("port", cfg.ServerPort).Info("Starting server")
    logrus.Info("Swagger documentation available at: http://localhost:8080/swagger/")
    
    if err := router.Run(addr); err != nil {
        logrus.Fatal("Failed to start server: ", err)
    }
}