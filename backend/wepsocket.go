package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var clients = make(map[*websocket.Conn]bool)     //stroe the connected user
var brodcast = make(chan map[string]interface{}) //channel for sending the new data to js

func ws(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("err in upgrader", err)
		return
	}
	defer conn.Close()

	clients[conn] = true
	fmt.Println("new user")

	for {
		var msg map[string]interface{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println("err reading msg", err)
			delete(clients, conn)
			break
		}
	}
}

func brodcastUpdate() {
	for {
		message := <-brodcast
		for client := range clients {
			err := client.WriteJSON(message)
			if err != nil {
				fmt.Println("Err writing message", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
