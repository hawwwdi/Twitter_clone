package db

import (
	"github.com/go-redis/redis"
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
	if rdb.Get(lastIdC).Err() == redis.Nil {
		err = rdb.Set(lastIdC, "0", 0).Err()
		checkErr(err)
	}
	if rdb.Get(lastPostC).Err() == redis.Nil {
		err = rdb.Set(lastPostC, "0", 0).Err()
		checkErr(err)
	}
	_ = rdb.Del("auths")
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

func (d *DB) ShowTimeLinePosts(count int) (map[string][]string, error) {
	return showTimeLinePosts(d.rdb, int64(count))
}

func (d *DB) ShowUserPosts(id string, start, count int) (map[string][]string, error) {
	return showUserPosts(d.rdb, id, int64(start), int64(count))
}

func (d *DB) GetSessionUserID(session string) (string, error) {
	return getSessionUserID(d.rdb, session)
}
