package main

import (
	"database/sql"
	  "encoding/json"
	"log"
	"net/http"
	"os"
    

	_ "modernc.org/sqlite"

)
type TextMessage struct {
    Ok     bool `json:"ok"`
    Result []struct {
        UpdateID int `json:"update_id"`
        Message  struct { 
            MessageID int `json:"message_id"`
            From      struct {
                ID           int    `json:"id"`
                IsBot        bool   `json:"is_bot"`
                FirstName    string `json:"first_name"`
                LastName     string `json:"last_name"`
                LanguageCode string `json:"language_code"`
            } `json:"from"`
            Chat struct {
                ID        int    `json:"id"`
                FirstName string `json:"first_name"`
                LastName  string `json:"last_name"`
                Type      string `json:"type"`
            } `json:"chat"`
            Date int    `json:"date"`
            Text string `json:"text"`
        } `json:"message"`
    } `json:"result"`
}


        type message struct {
            MessageID int `json:"message_id"`
            From      struct {
                ID           int    `json:"id"`
                IsBot        bool   `json:"is_bot"`
                FirstName    string `json:"first_name"`
                LastName     string `json:"last_name"`
                LanguageCode string `json:"language_code"`
            } `json:"from"`
            Chat struct {
                ID        int    `json:"id"`
                FirstName string `json:"first_name"`
                LastName  string `json:"last_name"`
                Type      string `json:"type"`
            } `json:"chat"`
            Date int    `json:"date"`
            Text string `json:"text"`
        } 
  

var db *sql.DB

func initDB() error {
	var err error
	db, err = sql.Open("sqlite", "telegram.db")
	if err != nil {
		log.Fatal("Error initializing database: ", err)
	}

	return db.Ping()
}

func insertMessage(message message) {
	db.Ping()


	_, err := db.Exec("INSERT INTO telegram VALUES ($1, $2, $3)",
		message. MessageID, message.Text,message.Chat.ID)

	if err != nil {
		log.Fatal("Error message insert: ", err)
	}
}

func inputMessagesHandler(rw http.ResponseWriter, rq *http.Request) {
	var p TextMessage
	decoder := json.NewDecoder(rq.Body)
	if err := decoder.Decode(&p); err != nil {
		http.Error(rw, "Invalid JSON", http.StatusBadRequest)
		return
	}
	for index, element := range p.Result {

		log.Printf("%d) Inserting data with update_id = %d", index, element.UpdateID)
		insertMessage(element.Message)
	}
}

func main() {
    	log.SetPrefix("DBMS: ")
	f, err := os.OpenFile("logs.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	log.SetOutput(f)

	err = initDB()
	if err != nil {
		log.Fatal("Error while Ping Database: ", err)
	}

	http.HandleFunc("/api/input_messages", inputMessagesHandler)

	http.ListenAndServe(":8080", nil)
}