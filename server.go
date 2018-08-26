package main

import (
	//"net/http"

	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name" binding:"required"`
	Age  int    `json:"age" binding:"exists"`
}

// Server implements a http server
type Server struct {
	Addr string

	router      *gin.Engine
	batchRouter *gin.Engine
	usersIndex  map[int]*User
	users       []*User
}

func (s *Server) Start() error {
	return s.router.Run(s.Addr)
}

func NewServer() *Server {
	router := gin.Default()
	batchRouter := gin.Default()

	s := &Server{
		router:      router,
		batchRouter: batchRouter,
		usersIndex:  make(map[int]*User),
		users:       make([]*User, 0),
	}
	router.GET("/users/:userId", s.getUser)
	router.DELETE("/users/:userId", s.deleteUser)

	router.POST("/users", s.createUser)
	router.GET("/users", s.getUserList)

	router.POST("/batch", s.batchProcess)
	batchRouter.DELETE("/users/:userId", s.deleteUser)

	return s
}

func (s *Server) getUserList(c *gin.Context) {
	c.JSON(200, s.users)
}

func (s *Server) createUser(c *gin.Context) {
	var p User
	err := c.Bind(&p)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	id := len(s.users) + 1
	p.ID = id
	s.users = append(s.users, &p)
	s.usersIndex[id] = &p
	c.JSON(200, struct {
		ID int
	}{
		id,
	})
}

func (s *Server) getUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.AbortWithError(400, err)
	}

	c.JSON(200, s.usersIndex[userID])
}

func (s *Server) deleteUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.AbortWithError(400, err)
	}

	delete(s.usersIndex, userID)
	users := make([]*User, 0, len(s.users)-1)
	for _, u := range s.users {
		if u.ID != userID {
			users = append(users, u)
		}
	}

	/*
		if userID < len(s.users){
			c.AbortWithError(400, errors.New(invalidUserID))
			return
		}
	*/
	s.users = users
	c.JSON(200, struct {
		DeletedUserID int `json:"deletedUserId"`
	}{
		userID,
	})

}

func (s *Server) batchProcess(c *gin.Context) {
	mr, err := c.Request.MultipartReader()
	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	for i := 0; ; i++ {
		log.Printf("++++++++++++++get multipart, %d", i)
		part, errPart := mr.NextPart()
		if errPart == io.EOF {
			log.Print("multipart EOF found\n")
			break
		}

		if errPart != nil {
			log.Printf("multipart err, %T, %s", errPart, errPart)
			c.AbortWithError(400, errPart)
			return
		}

		log.Printf("header: %+v\n", part.Header)
		data, _ := ioutil.ReadAll(part)
		log.Printf("data: %s\n", data)

		req, err := http.ReadRequest(bufio.NewReader(strings.NewReader(string(data) + "\r\n")))
		if err != nil {
			log.Printf("failed to read request, %s", err)
			return
		}
		log.Printf("request: %+v\n", req)
		/*
			requestBody, err := ioutil.ReadAll(req.Body)
			if err != nil {
				log.Printf("failed to read request body, %s", err)
				return
			}
			log.Printf("request body: %s\n", string(requestBody))
		*/
		w := newImMemoryResponseWriter()
		s.batchRouter.ServeHTTP(w, req)
		log.Printf("response, statusCode=%d, %+v, %s", w.statusCode, w.header, string(w.data.Bytes()))

	}

	c.JSON(200, struct{}{})
}

type InMemoryResponseWriter struct {
	header     http.Header
	data       bytes.Buffer
	statusCode int
}

func newImMemoryResponseWriter() *InMemoryResponseWriter {
	return &InMemoryResponseWriter{
		header: make(map[string][]string),
	}
}

func (rw *InMemoryResponseWriter) Header() http.Header {
	return rw.header
}

func (rw *InMemoryResponseWriter) Write(d []byte) (int, error) {
	return rw.data.Write(d)
}

func (rw *InMemoryResponseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
}
