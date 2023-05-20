package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	BASE_URL  = "https://api.telegram.org/bot"
	BOT_TOKEN = "6253089031:AAHZPs-z7R5TyctRaAnpUWHDo8EoC_2TbxQ"
)

var (
	telegramMethods = map[string]string{
		"GET_UPDATES":  "getUpdates",
		"SEND_MESSAGE": "sendMessage",
	}
	offset = 0
)

func main() {
}

func getUpdates(baseUrl string, botToken string, apiMethod string) ([]Update, error) {
	apiUrl := baseUrl + botToken + "/" + apiMethod

	resp, err := http.Get(apiUrl)
	if err != nil {
		fmt.Println("Ошибка при отправке запроса:", err)
		return nil, err
	}

	// ответ от сервера получаем в байтах, необходимо обработать его
	body, err := io.ReadAll(resp.Body)

	defer resp.Body.Close()

	if err != nil {
		if err != nil {
			fmt.Println("Ошибка при обработки запроса:", err)
			return nil, err
		}
	}

	var restResponse RestResponse

	// необходимо распарсить json, полученный от сервера, который приведем к структуре RestResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		return nil, err
	}

	return restResponse.Result, nil
}
