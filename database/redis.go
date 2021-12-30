package database

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/easilok/mark-notes-server/helpers"
	"github.com/go-redis/redis/v7"
)

var client *redis.Client
var EnableRedis bool = false

func Initialize() {
	//Initializing redis
	useRedis := os.Getenv("REDIS_USAGE")
	fmt.Println("Redis usage: ", useRedis)

	if useRedis == "True" {

		EnableRedis = true
		dsn := os.Getenv("REDIS_DSN")
		if len(dsn) == 0 {
			dsn = "localhost:6379"
		}

		client = redis.NewClient(&redis.Options{
			Addr: dsn, //redis port
		})
		_, err := client.Ping().Result()
		if err != nil {
			panic(err)
		}
	}
}

func CreateAuth(userid uint64, td *helpers.TokenDetails) error {

	if !EnableRedis {
		return nil
	}

	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := client.Set(td.AccessUuid, strconv.Itoa(int(userid)), at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := client.Set(td.RefreshUuid, strconv.Itoa(int(userid)), rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

func FetchAuth(authD *helpers.AccessDetails) (uint64, error) {
	if !EnableRedis {
		return authD.UserId, nil
	}

	userid, err := client.Get(authD.AccessUuid).Result()
	if err != nil {
		return 0, err
	}
	userID, _ := strconv.ParseUint(userid, 10, 64)
	return userID, nil
}

func DeleteAuth(givenUuid string) (int64, error) {
	if !EnableRedis {
		return 0, nil
	}

	deleted, err := client.Del(givenUuid).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}
