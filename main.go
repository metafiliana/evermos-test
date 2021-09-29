package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/metafiliana/evermos-test/config"
	"github.com/metafiliana/evermos-test/handler"
	"github.com/metafiliana/evermos-test/order"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	cfg := config.Get()

	// init db dependencies
	dbConn := InitDB(cfg)
	redisPool, err := InitRedis(cfg)
	if err != nil {
		log.Println(`[SERVICE] can't connect to redis'`)
	}

	// init repo
	repo := order.NewRepository(dbConn)
	cacheRepo := order.NewCacheRepository(redisPool)

	// init service
	svc := order.NewService(repo, cacheRepo)

	// init handler
	handler := handler.NewHandler(svc)

	// route start here
	router.POST(`/api/checkout-items`, handler.CheckoutItems)
	router.POST(`/api/order`, handler.CreateOrder)

	if err := http.ListenAndServe(cfg.RestPort, router); err != nil {
		log.Fatal(err)
	}
}

func InitDB(cfg config.Config) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		cfg.DBUserName, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf(err.Error())
		panic(err)
	}

	return db
}

func InitRedis(config config.Config) (*redis.Pool, error) {
	Redis := &redis.Pool{
		MaxIdle:   10,
		MaxActive: 20,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", config.RedisAddress)
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}
	conn := Redis.Get()
	_, err := conn.Do("PING")
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return Redis, nil
}
