package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/rs/zerolog/log"
)

const (
	host         = "127.0.0.1"
	port         = "4201"
	gitHubURL    = "https://api.github.com"
	clientID     = "bfae7ef22c970e9166b2"
	clientSecret = "20f42a396563ab6ae4fdbda80918c9d7aa195eed"
	redirectURL  = "http://127.0.0.1:4200"
)

func main() {
 // получаем запрос от сервера
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// получаем данные из запроса
		code := r.RequestURI[2:]
		// формируем ответный запрос на авторизацию и отправляем обратно
		q := "https://github.com/login/oauth/access_token" +
			fmt.Sprintf("?client_id=%s", clientID) +
			fmt.Sprintf("&client_secret=%s", clientSecret) +
			fmt.Sprintf("&%s", code) +
			fmt.Sprintf("&redirect_uri=%s", url.QueryEscape(redirectURL))
		req, err := http.Post(
			q,
			"text/plain",
			nil,
		)
		if err != nil {
			log.Error().Err(err).Msg("Failed to prepare get token request")
			w.WriteHeader(http.StatusInternalServerError)
		}

		switch req.StatusCode {
		case 200:
			// считываем полученные данные от сервера и переправляем клиента
			resp, err := ioutil.ReadAll(req.Body)
			if err != nil {
				// TODO Handle error
				log.Error().Err(err).Msgf("Failed to get token from github: %v", req.Status)
			}
			token := strings.Split(strings.Split(string(resp), "&")[0], "=")[1]
			// сохранили токен для клиента
			http.SetCookie(w, &http.Cookie{Name: "token", Value: token})
			// перенаправили клиента на его страницу
			http.Redirect(w, r, redirectURL, http.StatusMovedPermanently)
		case 422:
			w.WriteHeader(http.StatusUnprocessableEntity)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), nil)
}
