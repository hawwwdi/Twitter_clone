package db

import "github.com/go-redis/redis"

var rdb *redis.Client

const (
	lastIdC   = "last_id"
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
}

func AddUser(usr user) error {
	var err error
	_, username, password := usr.info()
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
	usr.setId(id)
	err = rdb.HSet(usersMapC, username, lastID).Err()
	if err != nil {
		return err
	}
	err = rdb.Incr(lastIdC).Err()
	checkErr(err)
	return nil
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
