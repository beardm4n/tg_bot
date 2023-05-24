package main

import (
	"fmt"
	"time"
)

var (
	telegramMethods = map[string]string{
		"GET_UPDATES":  "getUpdates",
		"SEND_MESSAGE": "sendMessage",
		"GET_FILE":     "getFile",
	}
	commands = map[string]string{
		"/start": "Bot is active now. Send /commands to get a list of commands",
		"/about": "About bot",
	}
	offset   = 0
	botToken string
	baseUrl  string
	fileUrl  string
	basePathToSaveFile string
)

func main() {
	/*
	 Загружаем переменные окружения
	 Внутри метода initEnv можно инициализировать переменную, которую берем из env (но сначала надо ее определить в var)
	*/
	err := initEnv()
	if err != nil {
		fmt.Println("Error in loadEnv method: ", err)
		return
	}

	for {
		updates, err := getUpdates(baseUrl, botToken, telegramMethods["GET_UPDATES"], offset)
		if err != nil {
			fmt.Println("Something went wrong in getUpdates: ", err)
		}

		for _, update := range updates {
			processUpdate(update)

			offset = update.UpdateId + 1
		}

		fmt.Printf("Updates: %v\n", updates)

		// Пауза между запросами на получение обновлений
		time.Sleep(time.Second)
	}
}
