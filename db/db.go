package db

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
	gouuid "github.com/satori/go.uuid"
)

const (
	lastIdC   = "last_id"
	lastPostC = "last_post"
	usernameC = "username"
	passwordC = "password"
	usersMapC = "users"
)

type DB struct {
	rdb *redis.Client
}

func NewDB(addr string) *DB {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	_, err := rdb.Ping().Result()
	checkErr(err)
	err = rdb.Set(lastIdC, "0", 0).Err()
	checkErr(err)
	err = rdb.Set(lastPostC, "0", 0).Err()
	checkErr(err)
	return &DB{
		rdb: rdb,
	}
}

func (d *DB) RegisterUser(user, pass string) error {
	return registerUser(d.rdb, user, pass)
}

func (d *DB) LogIn(username, password string) (string, error) {
	return logIn(d.rdb, username, password)
}

func (d *DB) LogOut(session string) error {
	return logOut(d.rdb, session)
}

func (d *DB) Follow(follower, followed string) error {
	return follow(d.rdb, follower, followed)
}

func (d *DB) Post(body, owner string) error {
	return post(d.rdb, body, owner)
}

func (d *DB) ShowUserPosts(id string, start, count int) (map[string][]string, error) {
	return showUserPosts(d.rdb, id, int64(start), int64(count))
}

func (d *DB) GetSessionUserID(session string) (string, error) {
	return getSessionUserID(d.rdb, session)
}

func registerUser(rdb *redis.Client, username, password string) error {
	var err error
	exists := checkUsername(rdb, username) != ""
	if exists {
		return fmt.Errorf("user %v already exists", username)
	}
	lastID, err := rdb.Get(lastIdC).Result()
	checkErr(err)
	id := "user:" + lastID
	err = rdb.HSet(id, usernameC, username).Err()
	if err != nil {
		return err
	}
	err = rdb.HSet(id, passwordC, password).Err()
	if err != nil {
		return err
	}
	err = rdb.HSet(usersMapC, username, lastID).Err()
	if err != nil {
		return err
	}
	err = rdb.Incr(lastIdC).Err()
	checkErr(err)
	return nil
}

func logIn(rdb *redis.Client, username, password string) (string, error) {
	var err error
	id := checkUsername(rdb, username)
	if id == "" {
		return "", fmt.Errorf("user %v does not exists", username)
	}
	err = checkPassword(rdb, id, password)
	if err != nil {
		return "", err
	}
	removeSession(rdb, id)
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

func logOut(rdb *redis.Client, session string) error {
	rdb.HDel("auths", session)
	return nil
}

func follow(rdb *redis.Client, follower, followed string) error {
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

func post(rdb *redis.Client, body, owner string) error {
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

func showUserPosts(rdb *redis.Client, id string, start, count int64) (map[string][]string, error) {
	//todo add models
	if checkID(rdb, id) != nil {
		return nil, nil
	}
	posts, err := rdb.LRange("posts:"+id, start, start+count).Result()
	if err != nil {
		return nil, err
	}
	postsMap := make(map[string][]string)
	for _, post := range posts {
		owner, body, _ := showPost(rdb, post)
		postsMap[post] = []string{owner, body}
	}
	return postsMap, nil
}

func showPost(rdb *redis.Client, postId string) (string, string, error) {
	post, err := rdb.HGetAll("post:" + postId).Result()
	if err != nil {
		return "", "", nil
	}
	return post["owner"], post["body"], nil
}

func getSessionUserID(rdb *redis.Client, auth string) (string, error) {
	return rdb.HGet("auths", auth).Result()
}

func removeSession(rdb *redis.Client, id string) error {
	uuid, _ := rdb.HGet("user:"+id, "auth").Result()
	_, _ = rdb.HDel("auths", uuid).Result()
	return nil
}

func checkUsername(rdb *redis.Client, username string) string {
	id, _ := rdb.HGet(usersMapC, username).Result()
	return id
}

func checkID(rdb *redis.Client, id string) error {
	exists, _ := rdb.HExists("user:"+id, "username").Result()
	if !exists {
		return fmt.Errorf("user %v not found", id)
	}
	return nil
}

func checkPassword(rdb *redis.Client, id, password string) error {
	pass, _ := rdb.HGet("user:"+id, passwordC).Result()
	if pass != password {
		log.Printf("id==%v pass %v != %v\n", id, pass, password)
		return errors.New("invalid password")
	}
	return nil
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
