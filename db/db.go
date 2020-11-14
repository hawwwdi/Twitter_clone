package db

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
	gouuid "github.com/satori/go.uuid"
)

var rdb *redis.Client

const (
	lastIdC   = "last_id"
	lastPostC = "last_post"
	usernameC = "username"
	passwordC = "password"
	usersMapC = "users"
)

func init() {
	//todo get port from environment variables
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6037",
	})
	_, err := rdb.Ping().Result()
	checkErr(err)
	err = rdb.Set(lastIdC, "0", 0).Err()
	checkErr(err)
	err = rdb.Set(lastPostC, "0", 0).Err()
	checkErr(err)
}

func RegisterUser(usr user) (string, error) {
	var err error
	_, username, password := usr.Info()
	exists := checkUsername(username) != ""
	if exists {
		return "", fmt.Errorf("user %v already exists", username)
	}
	lastID, err := rdb.Get(lastIdC).Result()
	checkErr(err)
	id := "user:" + lastID
	err = rdb.HSet(id, usernameC, username).Err()
	if err != nil {
		return "", err
	}
	err = rdb.HSet(id, passwordC, password).Err()
	if err != nil {
		return "", err
	}
	usr.SetId(lastID)
	err = rdb.HSet(usersMapC, username, lastID).Err()
	if err != nil {
		return "", err
	}
	err = rdb.Incr(lastIdC).Err()
	checkErr(err)
	return lastID, nil
}

func LogIn(username, password string) (string, error) {
	var err error
	id := checkUsername(username)
	if id == "" {
		return "", fmt.Errorf("user %v does not exists", username)
	}
	err = checkPassword(id, password)
	if err != nil {
		return "", err
	}
	removeSession(id)
	uuid := gouuid.NewV4().String()
	err = rdb.HSet("user:"+id, "auth", uuid).Err()
	if err != nil {
		return "", err
	}
	err = rdb.HSet("auths", uuid, id).Err()
	if err != nil {
		return "", err
	}
	return uuid, nil
}

func LogOut(session string) error {
	id, err := IsLoggedIn(session)
	if err != nil {
		return err
	}
	rdb.HDel("auths", session)
	rdb.HDel("user:"+id, "auth")
	return nil
}

func Follow(follower, followed string) error {
	var err error
	currentTime := float64(time.Now().Unix())
	_, err = rdb.ZAdd("followings:"+follower, redis.Z{
		Score:  currentTime,
		Member: followed,
	}).Result()
	if err != nil {
		return err
	}
	_, err = rdb.ZAdd("followers:"+followed, redis.Z{
		Score:  currentTime,
		Member: follower,
	}).Result()
	return err
}

func Post(body, owner string) error {
	var err error
	postID, err := rdb.Get(lastPostC).Result()
	checkErr(err)
	err = rdb.HSet("post:"+postID, "owner", owner).Err()
	if err != nil {
		return err
	}
	err = rdb.HSet("post:"+postID, "body", body).Err()
	if err != nil {
		return err
	}
	err = rdb.Incr(lastPostC).Err()
	checkErr(err)
	_, err = rdb.LPush("posts:"+owner, postID).Result()
	//todo use concurrent pattern
	followers, err := rdb.ZRevRange("followers:"+owner, 0, -1).Result()
	if err != nil {
		return err
	}
	for _, follower := range followers {
		_, err = rdb.LPush("posts:"+follower, postID).Result()
		if err != nil {
			return err
		}
	}
	err = rdb.LPush("timeline", postID).Err()
	if err != nil {
		return err
	}
	err = rdb.LTrim("timeline", 0, 100).Err()
	return err
}

func removeSession(id string) error {
	uuid, _ := rdb.HGet("user:"+id, "auth").Result()
	_, _ = rdb.HDel("auths", uuid).Result()
	return nil
}

func checkUsername(username string) string {
	id, _ := rdb.HGet(usersMapC, username).Result()
	return id
}

func checkPassword(id, password string) error {
	pass, _ := rdb.HGet("user:"+id, passwordC).Result()
	if pass != password {
		log.Printf("id==%v pass %v != %v\n", id, pass, password)
		return errors.New("invalid password")
	}
	return nil
}

func IsLoggedIn(auth string) (string, error) {
	return rdb.HGet("auths", auth).Result()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
