# golangchannel

### 下面情况执行【 go test -v ./... --race 】会导致 race condition
main.go:
```
func (s *Server) addUser(user string) {
	s.users[user] = user + "value"
}
```

test.go
```
func TestAddUser(t *testing.T) {
	server := NewServer()

	for i := 0; i < 10; i++ {
		go server.addUser(fmt.Sprintf("user_%d", i))
	}
}
```

### Solution 1 ： 使用互斥锁

main.go
```
type Server struct {
	users map[string]string
	mu    sync.Mutex  
}

func (s *Server) addUserSafe(user string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.users[user] = user + "value"
}
```
test.go
```
func TestAddUserSafe(t *testing.T) {
	server := NewServer()

	for i := 0; i < 10; i++ {
		go server.addUserSafe(fmt.Sprintf("user_%d", i))
	}
}
```

### Solution 2: Channel

main.go

```
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
	}
}
```

test.go

```
func TestAddUser(t *testing.T) {
	server := NewServer()
	server.Start()

	for i := 0; i < 10; i++ {
		go func(i int) {
			server.userch <- fmt.Sprintf("user_%d", i)
		}(i)
	}
}
```
### 下面程序会 fatal error: all goroutines are asleep - deadlock!


是因为在程序的主goroutine中，只有一个select语句等待userch通道的事件，而没有任何goroutine往该通道发送消息，因此主goroutine会永远阻塞在select语句中等待接收userch通道的消息，而不会继续执行下去。

```
userch := make(chan string, 10)
ticker:= time.NewTicker(2 * time.Second)
for {
    select{
        case <- ticker.C：
          print("tick")
        case <- userch:
    }
}
```
为了避免发生死锁，你应该在程序中创建一个新的goroutine来往userch通道发送消息，以确保在主goroutine等待接收userch通道消息时，有其他的goroutine可以往该通道发送消息，从而解除死锁的状态。例如，你可以使用以下代码来避免死锁的发生
```
package main

import "fmt"

func main() {
    userch := make(chan string, 10)
    go func() {
        userch <- "Hello, world!"
    }()
    for {
        select {
        case msg := <-userch:
            fmt.Println(msg)
        }
    }
}
```

### 下面会避免deadlock，但会无限循环

```
func main() {
	userch := make(chan string, 10)
	for {
		select {
		case <-userch:
		default:
		}
	}
}
```