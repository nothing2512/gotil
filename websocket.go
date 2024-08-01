package gotil

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
)

// websocket base model
type WebSocket struct {
	// Server
	clients map[string]*websocket.Conn
	handle  func(m WebSocketMessage)
	baseUri string

	// Util
	encryption *Encryption

	// Client
	id         string
	token      string
	connected  bool
	connection *websocket.Conn
}

// web socket base message model
type WebSocketMessage struct {
	Command string `json:"command"`
	Message string `json:"message"`
	Token   string `json:"token"`
	To      string `json:"to"`
}

// Create New Websocket
func NewWebSocket(baseUri string) *WebSocket {
	return &WebSocket{baseUri: baseUri}
}

// Start Server
func (s *WebSocket) Server(secret, iv string) error {
	s.encryption = NewEncryption(secret, iv)
	s.clients = make(map[string]*websocket.Conn)

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	http.HandleFunc("/ws/connect", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		for {
			t, message, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				break
			}
			switch t {
			case websocket.TextMessage:

				var m WebSocketMessage
				_ = json.Unmarshal(message, &m)

				if m.Command == "connect" {
					uid := UUID()
					b, _ := json.Marshal(JSON{
						"id":      uid,
						"token":   s.encryption.Encrypt(uid),
						"command": "connect",
					})
					err := conn.WriteMessage(t, b)
					s.clients[uid] = conn
					if err != nil {
						log.Println(err)
						return
					}
				} else if m.Command == "send" {
					s.send(m)
				} else if m.Command == "disconnect" {
					token := m.Token
					uid := s.encryption.Decrypt(token)
					c := s.clients[uid]
					c.WriteMessage(websocket.TextMessage, []byte("Disconnected"))
					c.Close()
					delete(s.clients, uid)
				} else if m.Command == "command" {
					if s.handle != nil {
						var cmd WebSocketMessage
						_ = json.Unmarshal([]byte(m.Message), &cmd)
						cmd.Token = m.Token
						cmd.To = s.encryption.Decrypt(m.Token)
						s.handle(cmd)
					}
				}
			}
		}
	})
	http.HandleFunc("/ws/send", func(w http.ResponseWriter, r *http.Request) {
		var m WebSocketMessage
		_ = json.NewDecoder(r.Body).Decode(&m)
		if m.Command == "send" {
			s.send(m)
			w.Write([]byte("Sent"))
		} else if m.Command == "command" {
			token := m.Token
			uid := s.encryption.Decrypt(token)
			if _, exists := s.clients[uid]; !exists {
				w.Write([]byte("Unsent"))
				return
			}
			if s.handle != nil {
				var cmd WebSocketMessage
				_ = json.Unmarshal([]byte(m.Message), &cmd)
				s.handle(cmd)
			}
			w.Write([]byte("Sent"))
		} else {
			w.Write([]byte("Unsent"))
		}
	})

	ln, err := net.Listen("tcp", s.baseUri)
	if err != nil {
		return err
	}

	fmt.Println("Listening WebSocket On " + s.baseUri)
	if err = http.Serve(ln, nil); err != nil {
		return err
	}

	return nil
}

// Handle Command On Server
func (s *WebSocket) OnCommand(handle func(m WebSocketMessage)) {
	s.handle = handle
}

// Server Send Message
func (s *WebSocket) send(m WebSocketMessage) {
	token := m.Token
	uid := s.encryption.Decrypt(token)
	if _, exists := s.clients[uid]; exists {
		if to, exists := s.clients[m.To]; exists {
			to.WriteMessage(websocket.TextMessage, []byte(m.Message))
		}
	}
}

// Server Send Reply Message
func (s *WebSocket) Reply(m WebSocketMessage, message string) error {
	m.Command = "reply"
	m.Message = message
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	m.Message = string(b)
	s.send(m)
	return nil
}

// Server Blast Message
func (s *WebSocket) Blast(m string) error {
	b, err := json.Marshal(WebSocketMessage{
		Command: "blast",
		Message: m,
	})
	if err != nil {
		return err
	}
	go func() {
		for _, c := range s.clients {
			c.WriteMessage(websocket.TextMessage, b)
		}
	}()
	return nil
}

// Start Client
func (s *WebSocket) Client() error {
	s.baseUri = fmt.Sprintf("ws://%v", s.baseUri)

	if s.connected {
		return nil
	}
	uri := fmt.Sprintf("%v/connect", s.baseUri)
	conn, _, err := websocket.DefaultDialer.Dial(uri, nil)
	if err != nil {
		log.Fatal(err)
	}

	message, _ := json.Marshal(WebSocketMessage{Command: "connect"})
	err = conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		return err
	}

	_, reply, err := conn.ReadMessage()
	if err != nil {
		return err
	}
	var data struct {
		Id    string `json:"id"`
		Token string `json:"token"`
	}
	_ = json.Unmarshal(reply, &data)

	s.id = data.Id
	s.token = data.Token
	s.connected = true
	s.connection = conn

	return nil
}

// Disconnect Client
func (s *WebSocket) Disconnect() error {
	if !s.connected {
		return nil
	}

	message, _ := json.Marshal(WebSocketMessage{
		Token:   s.token,
		Command: "disconnect",
	})
	err := s.connection.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		return err
	}

	s.connection.Close()
	s.connected = false

	return nil
}

// Client Send Command
func (s *WebSocket) Command(data WebSocketMessage) error {
	if !s.connected {
		return nil
	}

	messageData, _ := json.Marshal(data)

	message, _ := json.Marshal(WebSocketMessage{
		Token:   s.token,
		Command: "command",
		Message: string(messageData),
	})
	err := s.connection.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		return err
	}

	return nil
}

// Client On Message Received
func (s *WebSocket) OnMessage(handle func(m WebSocketMessage)) {
	for {
		_, message, err := s.connection.ReadMessage()
		if err != nil {
			s.connected = false
			return
		}
		var m WebSocketMessage
		_ = json.Unmarshal(message, &m)
		handle(m)
	}
}

// Client Send Message To Other
func (s *WebSocket) Send(to string, msg string) error {
	if !s.connected {
		return nil
	}

	messageData, _ := json.Marshal(WebSocketMessage{
		Message: msg,
	})

	message, _ := json.Marshal(WebSocketMessage{
		Token:   s.token,
		To:      to,
		Command: "send",
		Message: string(messageData),
	})
	err := s.connection.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		return err
	}

	return nil
}
