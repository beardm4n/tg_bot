package file

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	types "tg_bot/helpers"
)

func getFile(baseUrl, botToken, telegramMethod, fileId string) types.File {
	urlGetFile := baseUrl + botToken + "/" + telegramMethod + "?file_id=" + fileId

	resp, err := http.Get(urlGetFile)
	if err != nil {
		fmt.Println("Something went wrong in error handling in getFile: ", err)
	}

	// ответ от сервера получаем в байтах, необходимо обработать его
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		if err != nil {
			fmt.Println("Something went wrong in error handling in getUpdates: ", err)
		}
	}

	defer resp.Body.Close()

	var restResponse types.VoiceResponse

	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		fmt.Println("Something went wrong in parse json in getFile: ", err)
	}

	return restResponse.Result
}

func DownloadFile(baseUrl, botToken, telegramMethod, fileUrl, basePathToSaveFile string, message types.Message) {
	file := getFile(baseUrl, botToken, telegramMethod, message.Voice.FileId)

	downloadFileUrl := filepath.FromSlash(fmt.Sprintf(fileUrl + botToken + "/" + file.FilePath))

	// Выполняем запрос для получения содержимого файла
	fileResp, err := http.Get(downloadFileUrl)
	if err != nil {
		fmt.Println("Something went wrong in error handling in download file")
		return
	}

	defer fileResp.Body.Close()

	filePathOnDisk := filepath.FromSlash(fmt.Sprintf("%s/voice_%s.oga", basePathToSaveFile, message.Voice.FileId))

	// Создаем файл на диске для сохранения содержимого
	fileOnDisk, err := os.Create(filePathOnDisk)
	if err != nil {
		fmt.Println("Something went wrong in create file into disk")
		return
	}

	defer fileOnDisk.Close()

	// Копируем полученный файл в файл на диске
	_, err = io.Copy(fileOnDisk, fileResp.Body)
	if err != nil {
		fmt.Println("Something went wrong when save file into disk")
		return
	}
}
