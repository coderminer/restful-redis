### 使用`httprouter`和`Redis`构建`RESTful API`服务

#### 安装依赖

```
go get github.com/julienschmidt/httprouter
go get github.com/garyburd/redigo/redis
```

编译运行 

```
go build
./restful-redis
```

在浏览器中访问  

```
http://localhost:8080           //欢迎信息
http://localhost:8080/posts     //获取所有的提交信息
http://localhost:8080/posts/1   //查看提交id=1的信息
```

#### 数据model `models.go`

```
package main

import "time"

type User struct {
	Id       int    `json:"id"`
	UserName string `json:"username"`
	Email    string `json:"email"`
}

type Comment struct {
	Id        int       `json:"id"`
	User      User      `json:"user"`
	Text      string    `json:"text"`
	Timestamp time.Time `json:"timestamp"`
}

type Post struct {
	Id        int       `json:"id"`
	User      User      `json:"user"`
	Topic     string    `json:"topic"`
	Text      string    `json:"text"`
	Comment   Comment   `json:"comment"`
	Timestamp time.Time `json:"timestamp"`
}

type Posts []Post
type Comments []Comment
type Users []User
```

#### 路由 `routes.go`

```
package main

import (
	mux "github.com/julienschmidt/httprouter"
)

type Route struct {
	Method string
	Path   string
	Handle mux.Handle
}

type Routes []Route

var routes = Routes{
	Route{
		"GET",
		"/",
		Index,
	},
	Route{
		"GET",
		"/posts",
		PostIndex,
	},
	Route{
		"GET",
		"/posts/:id",
		PostShow,
	},
	Route{
		"POST",
		"/posts",
		PostCreate,
	},
}

func NewRouter() *mux.Router {
	router := mux.New()
	for _, route := range routes {
		router.Handle(route.Method, route.Path, route.Handle)
	}
	return router
}

```


#### 路由方法 `handlers.go`

```
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	mux "github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, _ mux.Params) {
	fmt.Fprintf(w, "<h1>Hello,welcome to my blog!</h1>")
}

func PostIndex(w http.ResponseWriter, r *http.Request, _ mux.Params) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	posts := FindAll()
	if err := json.NewEncoder(w).Encode(posts); err != nil {
		panic(err)
	}
}

func PostShow(w http.ResponseWriter, r *http.Request, ps mux.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	HandleError(err)

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	post := FindPost(id)

	if err := json.NewEncoder(w).Encode(post); err != nil {
		panic(err)
	}

}

func PostCreate(w http.ResponseWriter, r *http.Request, _ mux.Params) {
	var post Post

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	HandleError(err)

	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &post); err != nil {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.WriteHeader(422)

		if err := json.NewEncoder(w).Encode(post); err != nil {
			panic(err)
		}
	}
	CreatePost(post)
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
}

```

#### Redis 数据库 `db.go`

```
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

```

[更多精彩内容](http://coderminer.com)  