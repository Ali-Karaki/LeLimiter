package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"rate-limiter/api/controllers"
	"rate-limiter/api/middleware"
	"rate-limiter/utils/logger"

	// "rate-limiter/api/middleware"
	"rate-limiter/pkg/constants"
	client "rate-limiter/pkg/redis/client"
	"rate-limiter/utils"
	"time"

	"github.com/redis/go-redis/v9"
)

func Simulate(endpoint string, rdb *redis.Client, logPath string, users []string) {
	logger.SwitchLogFile(logPath)

	ctx, cancel := context.WithTimeout(context.Background(), constants.CallsTimeout)
	defer cancel()

	log.Printf("Sending requests to %s", endpoint)
	utils.CallConcurrentReqs(ctx, endpoint, users)
	
	time.Sleep(1 * time.Second)
	log.Printf("Switching...")
}


func main() {

	logger.InitLogger()

	rdb, err := client.InitRedis()
	if err != nil {
		log.Fatalf("Error creating Redis client: %v", err)
	}
	log.Println("Redis client created")

	http.HandleFunc("/nolua", middleware.RateLimit(controllers.FancyController))
	http.HandleFunc("/lua", middleware.RateLimitLua(controllers.FancyController))
	http.HandleFunc("/luasha", middleware.RateLimitLuaSha(controllers.FancyController))

	go func() {
		if err := http.ListenAndServe("localhost:8080", nil); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
		log.Print("Server started")
	}()

	time.Sleep(1 * time.Second)


	noluaUsers := []string{"nolua_1", "nolua_2", "nolua_3", "nolua_4", "nolua_5"}
	luaUsers := []string{"lua_1", "lua_2", "lua_3", "lua_4", "lua_5"}
	luaShaUsers := []string{"luasha_1", "luasha_2", "luasha_3", "luasha_4", "luasha_5"}

	Simulate("/nolua", rdb, logger.NoLuaPath, noluaUsers)
	Simulate("/lua", rdb, logger.LuaPath, luaUsers)
	Simulate("/luasha", rdb, logger.LuaShaPath, luaShaUsers)

	log.Print("Check Stats in ", logger.StatsPath)
	logger.RunAwk()

	os.Exit(0)

	select {}

}
