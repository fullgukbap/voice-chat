package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	log.Println("[ENTRYPOINT-serveWs] 새로운 웹소켓 연결 요청이 들어왔습니다.")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("http 요청을 websocket 요청으로 수정하는데 오류가 발생했습니다: %v\n", err)
		return
	}

	log.Printf("[INIT] %s의 클라이언트 인스턴스를 생성합니다. \n", r.RemoteAddr)
	// create a new client
	client := &Client{
		conn: conn,
		send: make(chan []byte, 256),
	}
	hub.register <- client

	go client.readPump(hub)
	go client.writePump()
	log.Printf("[INIT] %s의 클라이언트의 기본적인 설정을 완료했습니다.\n", r.RemoteAddr)
}

func main() {
	hub := NewHub()
	go hub.Run()

	http.Handle("/", http.FileServer(http.Dir("./static")))

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		// log.Printf("%s가 /ws endpoint에 요청이 들어와 핸들러 진입\n", r.RemoteAddr)
		serveWs(hub, w, r)
	})

	//start server
	log.Println("Server starting at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
