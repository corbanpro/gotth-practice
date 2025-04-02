package store

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDb() *gorm.DB {
	dbName := "./todo.db"

	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&TodoItem{}, &User{})

	if err != nil {
		panic(err)
	}

	// test connection
	err = db.Exec("SELECT 1").Error
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to Gorm")

	return db
}

func InitRedis(c *gin.Context) *redis.Client {
	db := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	_, err := db.Ping(c).Result()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to Redis")

	return db
}
