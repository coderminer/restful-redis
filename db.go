package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
)

var currentPostId int
var currentUserId int

func RedisConnect() redis.Conn {
	c, err := redis.Dial("tcp", ":6379")
	HandleError(err)
	return c
}

func init() {
	CreatePost(Post{
		User: User{
			UserName: "Kevin",
			Email:    "kevin.woo@163.com",
		},
		Topic: "First Post",
		Text:  "Hello everyone! This is awesomes",
	})

	CreatePost(Post{
		User: User{
			UserName: "John",
			Email:    "John@163.com",
		},
		Topic: "Second Post",
		Text:  "This is the second posts",
	})
}

func FindAll() Posts {
	var posts Posts

	c := RedisConnect()
	defer c.Close()

	keys, err := c.Do("KEYS", "posts:*")
	HandleError(err)

	for _, k := range keys.([]interface{}) {
		var post Post
		reply, err := c.Do("GET", k.([]byte))
		HandleError(err)

		if err := json.Unmarshal(reply.([]byte), &post); err != nil {
			panic(err)
		}
		posts = append(posts, post)
	}
	//fmt.Println("FindAll", posts)
	return posts
}

func FindPost(id int) Post {
	var post Post

	c := RedisConnect()
	defer c.Close()

	reply, err := c.Do("GET", "posts:"+strconv.Itoa(id))
	HandleError(err)

	fmt.Println("GET OK")

	if err = json.Unmarshal(reply.([]byte), &post); err != nil {
		panic(err)
	}
	return post
}

func CreatePost(p Post) {
	currentPostId += 1
	currentUserId += 1

	p.Id = currentPostId
	p.User.Id = currentUserId
	p.Timestamp = time.Now()

	c := RedisConnect()
	defer c.Close()

	b, err := json.Marshal(p)
	HandleError(err)

	reply, err := c.Do("SET", "posts:"+strconv.Itoa(p.Id), b)
	HandleError(err)

	fmt.Println("SET", reply)
}
