package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"nhooyr.io/websocket"
	"os"
	"os/signal"
	"time"
)

func main() {
	listen, err := net.Listen("tcp", ":8000")
	if err != nil {
		panic(err)
	}

	s := &http.Server{
		Handler:      websocketServer{},
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Print("Server is starting", listen.Addr())

	errc := make(chan error, 1)
	go func() {
		errc <- s.Serve(listen)
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)

	select {
	case err := <-errc:
		log.Println("Failed to serve", err)
	case sig := <-sigs:
		log.Println("Terminating:", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = s.Shutdown(ctx)
	if err != nil {
		log.Println("Error shutting down server:", err)
	}
}

type websocketServer struct {
}

func (s websocketServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Println("Error accepting websocket:", err)
		return
	}

	ctx := context.Background()
	conn := websocket.NetConn(ctx, c, websocket.MessageBinary)
	go ServeNetConn(conn)
}

func ServeNetConn(conn net.Conn) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Println("Error closing net.Conn", err)
		}
	}()

	timeout := 50 * time.Second
	timeoutChan := make(chan uint8, 1)

	const StopTimeout uint8 = 0
	const ContTimeout uint8 = 1
	const MaxMsgSize int = 4 * 1024

	// ReadData
	go func() {
		msg := make([]byte, MaxMsgSize)
		for {
			n, err := conn.Read(msg)

			if err != nil {
				log.Println("Websocket read error:", err)
				timeoutChan <- StopTimeout // Stop timeout because of a read error
				return
			}

			// tick the timeout watcher so we don't timeout!
			timeoutChan <- ContTimeout

			log.Println("Message:", msg[:n])
		}
	}()

	// Manage timeout
ExitTimeout:
	for {
		select {
		case res := <-timeoutChan:
			if res == StopTimeout {
				log.Println("Manually stopping timeout manager")
				break ExitTimeout
			}
		case <-time.After(timeout):
			log.Println("User timeout!")
			break ExitTimeout

		}
	}

}
