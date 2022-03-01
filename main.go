package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	healthcheck "github.com/RaMin0/gin-health-check"
	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/unrolled/secure"

	"restaurant/config"

	cMenu "restaurant/controller/menu"
	cOrder "restaurant/controller/order"
	cUser "restaurant/controller/user"

	_ "restaurant/docs"

	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

/* init json configuration */
func init() {
	viper.SetConfigType("json")

	viper.SetConfigFile("config/Config.json")
	err := viper.ReadInConfig()
	if err != nil {
		//Fatal logs a message at level Fatal on the standard logger
		//then the process will exit with status set to 1.
		log.Fatal(err)
	}
}

func main() {

	setupEnvironment()
}

// @BasePath /api/v1/restaurant
// @title Web Order API
// @version 1.0
// @description This page is API documentation for simple order management system.
// @contact.name Developer
// @host localhost:8081
func setupEnvironment() {
	db := config.ConnectDB()
	defer db.Close()

	port := config.ListenAndServeServerPort()
	host := config.Hostname()
	env := config.Environment()
	if env == "development" {
		gin.SetMode(gin.DebugMode)
	} else if env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	//Prevent clickjacking
	//[FrameDeny] X-Frame-Options allows content publishers to prevent their own content from being used in an invisible frame by attackers
	secureMiddleware := secure.New(secure.Options{
		FrameDeny:             true,
		ContentSecurityPolicy: "frame-ancestors 'none'",
	})

	secureFunc := func() gin.HandlerFunc {
		return func(c *gin.Context) {
			err := secureMiddleware.Process(c.Writer, c.Request)

			if err != nil {
				c.Abort()
				return
			}

			if status := c.Writer.Status(); status > 300 && status < 399 {
				c.Abort()
			}
		}
	}()

	router := gin.New()
	router.Use(secureFunc)
	router.GET("/db/healthcheck", config.GetDBHealthCheck)
	// docs.SwaggerInfo.BasePath = "/api/v1"
	url := ginSwagger.URL(host + port + "/api/v1/restaurant/swagger/doc.json") // pointing to API definition
	middleware := router.Group("/api/v1/restaurant")
	{
		middleware.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
		middleware.POST("/menu", cMenu.CreateMenu)
		middleware.POST("/menucategory", cMenu.CreateMenuCategory)
		middleware.POST("/user/customer", cUser.CreateCustomer)
		middleware.POST("/order", cOrder.CreateOrder)
		middleware.GET("/menu/:id", cMenu.FindByID)
		middleware.GET("/menu/list", cMenu.GetMenus)
		middleware.GET("/order/list", cOrder.GetOrders)
	}

	c := cors.AllowAll() //allows all method
	handler := c.Handler(router)

	//Before it gets started
	//Important to do overall check-up to make sure the microservice is fine
	router.Use(healthcheck.Default())
	server := &http.Server{
		Addr:         port,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1) //len=1
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit //read
		log.Println("Server is shutting down")
		ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
		defer cancel()

		//control whether the http.setKeepAlives are enabled
		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			log.Fatal("Error occurred when shutdown the server")
		}
		close(done) //close a channel
	}()

	log.Printf("Server is ready to handle request on %s:%s", host, port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Couldn't listen on %s:%v\n", port, err)
	}

	<-done
	log.Println("Server stopped")
}
