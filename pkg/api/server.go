package api

import "net/http"

func StartServer(handler *Handler) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/tickers", handler.GetTickers)
	mux.HandleFunc("/api/quit", handler.BasicAuth(handler.QuitApp))

	server := &http.Server{
		Addr:    handler.c.ServerPort,
		Handler: mux,
	}

	return server
}
