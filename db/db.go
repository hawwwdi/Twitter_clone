package db

import "github.com/go-redis/redis"

var rdb *redis.Client

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6037",
	})
	_, err := rdb.Ping().Result()
	checkErr(err)
	err = rdb.Set("last_id", "0", 0).Err()
	checkErr(err)
}

func AddUser(usr user) error {
	var err error
	_, username, password := usr.info()
	last_id, err := rdb.Get("last_id").Result()
	id := "user:" + last_id
	err = rdb.HSet(id, "username", username).Err()
	if err != nil {
		return err
	}
	err = rdb.HSet(id, "password", password).Err()
	if err != nil {
		return err
	}
	usr.setId(id)
	err = rdb.HSet("users", username, last_id).Err()
	if err != nil {
		return err
	}
	err = rdb.Incr("last_id").Err()
	checkErr(err)
	return nil
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
