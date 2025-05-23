package api

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"net/http"
	"testgo/pkg/config"
	"testgo/pkg/service"
)

type Handler struct {
	AppService *service.AppService
	c          *config.Config
}

func NewHandler(appService *service.AppService, c *config.Config) *Handler {
	return &Handler{
		AppService: appService,
		c:          c,
	}
}

func (h *Handler) GetTickers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "указан неверный метод запроса", http.StatusMethodNotAllowed)
		return
	}
	result := h.AppService.GiveCurrencyMap()

	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}
}

func (h *Handler) BasicAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok {
			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(password))
			expectedUsernameHash := sha256.Sum256([]byte(h.c.Auth.Username))
			expectedPasswordHash := sha256.Sum256([]byte(h.c.Auth.Password))

			usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
			passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

			if usernameMatch && passwordMatch {
				next.ServeHTTP(w, r)
				return
			}
		}
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}

func (h *Handler) QuitApp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "неверный метод", http.StatusMethodNotAllowed)
		return
	}

	h.AppService.StopWork()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Сервер выключен"))
}
