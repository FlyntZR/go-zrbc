package main

import (
	"context"
	"go-zrbc/config"
	"strconv"
	"testing"

	"container/list"

	"github.com/go-redis/redis/v8"
)

func TestRedis_RedisList(t *testing.T) {
	config.Init("../cmd/ws-server/config.json")
	redisCli := redis.NewClient(&redis.Options{
		Addr:     config.Global.Redis.Addr,
		Password: config.Global.Redis.Password,
		// TODO: what is it ?
		DB: 0,
	})

	var testStings = []string{}
	for i := 0; i < 350; i++ {
		iStr := strconv.Itoa(i)
		testStings = append(testStings, iStr)
	}

	reply, err := redisCli.RPush(context.TODO(), "bib_test", testStings).Result()
	t.Logf("reply:(%d), err:(%+v)", reply, err)
	result, err := redisCli.LTrim(context.TODO(), "bib_test", -300, -1).Result()
	t.Logf("reply:(%s), err:(%+v)", result, err)
	values, err := redisCli.LRange(context.TODO(), "bib_test", 0, -1).Result()
	if err != nil {
		t.Fatalf("err:(%+v)", err)
	}
	t.Logf("values:(%+v)", values)
}

func TestRedis_ContainerList(t *testing.T) {
	// Create a new list and put some numbers in it.
	l := list.New()
	for i := 0; i < 350; i++ {
		l.PushFront(i)
	}

	for l.Len() > 300 {
		l.Remove(l.Back())
	}

	// Iterate through list and print its contents.
	for e := l.Front(); e != nil; e = e.Next() {
		t.Logf("e.Value:(%v)", e.Value)
	}
}

func TestRedis_AppendList(t *testing.T) {
	var testList = make([]int, 0)
	for i := 0; i < 350; i++ {
		testList = append(testList, i)
	}

	for _, item := range testList {
		t.Logf("item:(%v)", item)
	}
}

func TestRedis_Inrc(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "8.210.127.51:6379",
		Password: "DYzgB5JsV3RZNGrG", // no password set
		DB:       0,                  // use default DB
	})

	// err := rdb.Set(context.TODO(), "key", 0, 0).Err()
	// if err != nil {
	// 	t.Fatal(err)
	// }

	val, err := rdb.Incr(context.TODO(), "tKey").Result()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(val)

	val, err = rdb.Incr(context.TODO(), "tKey").Result()
	if err != nil {
		panic(err)
	}
	t.Log(val)

	val, err = rdb.Get(context.TODO(), "tKey").Int64()
	if err != nil {
		panic(err)
	}
	t.Log(val)
}
