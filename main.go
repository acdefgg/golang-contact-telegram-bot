package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Msg_struct struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Text  string `json:"text"`
}

func splitMessage(message string) []string {
	maxLen := 4000
	var parts []string
	for len(message) > maxLen {
		parts = append(parts, message[:maxLen])
		message = message[maxLen:]
	}
	parts = append(parts, message)
	return parts
}

func main() {
	bot, err := tgbotapi.NewBotAPI("YOUR_BOT_TOKEN")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.Error(w, "404 not found.", http.StatusNotFound)
			return
		}
		//fmt.Print(r.FormValue("name"))
		if r.Method == "POST" {

			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Ошибка чтения тела запроса", http.StatusBadRequest)
				return
			}
			defer r.Body.Close()

			var data Msg_struct
			if err := json.Unmarshal(body, &data); err != nil {
				http.Error(w, "Ошибка разбора JSON", http.StatusBadRequest)
				return
			}

			messageParts := splitMessage(data.Name + "\n\n" + data.Phone + "\n\n" + data.Text)

			for _, part := range messageParts {
				msg := tgbotapi.NewMessage(347468059, part)
				bot.Send(msg)
			}
			w.WriteHeader(http.StatusOK)
		}
	})

	log.Fatal(http.ListenAndServe(":3003", nil))
}
