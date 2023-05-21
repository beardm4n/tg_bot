package main

import (
	"fmt"
	"time"
)

var (
	telegramMethods = map[string]string{
		"GET_UPDATES":  "getUpdates",
		"SEND_MESSAGE": "sendMessage",
	}
	offset   = 0
	botToken string
	baseUrl  string
)

func main() {
	/*
	 Загружаем переменные окружения
	 Внутри метода можно инициализировать переменную, которую берем из env (но сначала надо ее определить в var)
	*/
	err := initEnv()
	if err != nil {
		fmt.Println("Error in loadEnv method: ", err)
		return
	}

	apiUrl := baseUrl + botToken + "/" + telegramMethods["GET_UPDATES"]

	for {
		updates, err := getUpdates(apiUrl, offset)
		if err != nil {
			fmt.Println("Something went wrong in getUpdates: ", err)
		}

		for _, update := range updates {
			checkMainCommand(update.Message)
			
			offset = update.UpdateId + 1
		}

		fmt.Printf("Updates: %v\n", updates)

		// Пауза между запросами на получение обновлений
		time.Sleep(time.Second)
	}
}
