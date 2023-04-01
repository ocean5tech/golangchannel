package main

import (
	"fmt"
)

func main() {

}

type Server struct {
	users  map[string]string
	userch chan string
}

func NewServer() *Server {

	return &Server{
		users:  make(map[string]string),
		userch: make(chan string),
	}
}

func (s *Server) Start() {
	go s.loop()
}

func (s *Server) loop() {
	for {
		user := <-s.userch
		s.users[user] = user
		fmt.Printf("adding new user %s\n", user)
	}
}

func (s *Server) addUser(user string) {
	s.users[user] = user + "value"
}

/* func (s *Server) addUserSafe(user string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.users[user] = user + "value"
} */

func sendMessage(msgch chan<- string) {
	msgch <- "hello"
}

func readMessage(msgch <-chan string) {
	msg := <-msgch
	fmt.Println(msg)
}
