package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", handleRequest)

	// Укажите порт, на котором будет запущен сервер
	port := "8080"

	log.Printf("Сервер запущен на http://localhost:%s\n", port)

	// Запустите сервер на указанном порту
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("Ошибка при запуске сервера:", err)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Временно наживать F5 для получения обновлений")
	// Отправьте GET-запрос на метод getUpdates Телеграм API
	response, err := http.Get("https://api.telegram.org/bot6190524855:AAFug058VCgIcqfmgyKiXq8j5GhQz3yEL-M/getUpdates")
	if err != nil {
		log.Println("Ошибка при получении обновлений:", err)
		return
	}
	defer response.Body.Close()

	// Прочитайте ответ в виде JSON-данных
	var result struct {
		Ok     bool     `json:"ok"`
		Result []Update `json:"result"`
	}
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		log.Println("Ошибка при чтении ответа:", err)
		return
	}

	// Обработайте полученные обновления
	for _, update := range result.Result {
		// Выполните необходимые действия с каждым обновлением
		log.Printf("Получено обновление: %+v\n", update)
	}

	// Отправьте ответ клиенту
	w.WriteHeader(http.StatusOK)
}
