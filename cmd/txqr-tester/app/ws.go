package main

import (
	"encoding/json"
	"log"
	"net"

	"github.com/divan/txqr/cmd/txqr-tester/ws"
	"github.com/gopherjs/websocket"
)

// WSClient implements WebSocket client that will talk to backend.
type WSClient struct {
	address string
	conn    net.Conn

	app *App // TODO(divan): figure out how can we avoid circular dependency
}

func NewWSClient(address string, app *App) *WSClient {
	client := &WSClient{
		address: address,
		app:     app,
	}

	return client
}

// talkToBackend establishes connection with backend and updates
// UI state based on output from backend.
func (w *WSClient) talkToBackend() {
	log.Println("Connecting to", w.address)
	conn, err := websocket.Dial(w.address)
	if err != nil {
		log.Println("[ERROR] Dial:", err)
		return
	}
	w.conn = conn
	defer w.conn.Close()

	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			log.Println("[ERROR] Reading from websocket: %v", err)
			break
		}

		log.Println("[DEBUG] Read:", string(buf[:n]))
		w.processWSCommand(buf[:n])
	}
}

func (w *WSClient) processWSCommand(data []byte) {
	var msg ws.UIRequest
	err := json.Unmarshal(data, &msg)
	if err != nil {
		log.Println("invalid command:", err)
		return
	}

	switch msg.Cmd {
	case ws.CmdConnect:
		log.Println("Got connect")
		w.app.SetConnected(true)
	}
}

func (w *WSClient) sendMsg(msg *ws.UIRequest) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Println("JSON marshal:", err)
		return
	}

	_, err = w.conn.Write(data)
	if err != nil {
		log.Println("write:", err)
		return
	}
}
