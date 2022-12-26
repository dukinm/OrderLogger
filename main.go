package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5" //nolint
	"io"
	"net/http" //nolint
	"os"
	"strconv"
)

type Config struct {
	Port                      int    `json:"port"`
	ShowLogInConsole          bool   `json:"show_log_in_console"`
	WriteLogToFolder          bool   `json:"write_log_to_folder"`
	LogFolderPath             string `json:"log_folder_path"`
	CreateLogFolderIfNotExits bool   `json:"create_log_folder_if_not_exits"`
	APIPath                   string `json:"api_path"`
}

func GetConfig() (config Config) {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	configFile := pwd + "/" + "config.json"
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		fmt.Println(err)
		os.Exit(1)
	}

	file, _ := os.ReadFile(configFile)

	err = json.Unmarshal([]byte(file), &config)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if config.WriteLogToFolder && config.CreateLogFolderIfNotExits {
		if _, err := os.Stat(config.LogFolderPath); os.IsNotExist(err) {
			err = os.Mkdir(config.LogFolderPath, 0755)
			if err != nil {
				fmt.Println("Произошла ошибка:")
				fmt.Println(err)
			}
		}
	}
	return config
}

func main() {

	config := GetConfig()
	fmt.Println("Привет, я программа, которая покажет данные о заказах отправляемые на сервер")
	fmt.Println("Не забудь выставить настройки отправки:")
	fmt.Println("Путь к серверу: http://127.0.0.1:" + strconv.Itoa(config.Port) + config.APIPath)
	fmt.Println("Текущая конфигурация:")
	if config.ShowLogInConsole {
		fmt.Println("- Лог будет выводиться в консоль")
	}
	if config.WriteLogToFolder {
		fmt.Println("- Лог будет записываться в папку: " + config.LogFolderPath)
		if config.CreateLogFolderIfNotExits {
			fmt.Println("- Если папка лога не существует, она будет создана")
		}
	}

	r := chi.NewRouter()

	r.Post(config.APIPath, func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, err := io.ReadAll(r.Body)

		if err != nil {
			fmt.Println("Произошла ошибка:")
			fmt.Println(err)
		}

		fmt.Println()
		if config.ShowLogInConsole {
			fmt.Println("Пришел запрос:")
			fmt.Println(string(bodyBytes))
		}

		if config.WriteLogToFolder {
			files, err := os.ReadDir(config.LogFolderPath)
			if err != nil {
				fmt.Println("Произошла ошибка:")
				fmt.Println(err)
			}
			newFile := len(files) + 1
			newFilePath := config.LogFolderPath + "/" + strconv.Itoa(newFile) + ".json"
			fmt.Println("Записываем запрос в файл: " + newFilePath)
			err = os.WriteFile(newFilePath, bodyBytes, 0655)
			if err != nil {
				fmt.Println("Произошла ошибка:")
				fmt.Println(err)
			}
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, err = w.Write(bodyBytes)
		if err != nil {
			fmt.Println("Произошла ошибка:")
			fmt.Println(err)
		}

	})

	err := http.ListenAndServe("127.0.0.1:"+strconv.Itoa(config.Port), r)

	if err != nil {
		fmt.Println("Произошла ошибка:")
		fmt.Println(err)
	}
}
