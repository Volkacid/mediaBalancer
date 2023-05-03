package handlers

import (
	"encoding/json"
	"github.com/Volkacid/mediaBalancer/internal/app"
	"github.com/Volkacid/mediaBalancer/internal/service"
	"github.com/Volkacid/mediaBalancer/internal/storage"
	"io"
	"net/http"
	"strconv"
)

type Room struct {
	Address string
}

func CreateNewRoom(servPool *service.AmazonServices, rdb *storage.RedisStorage) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer request.Body.Close()

		data, err := io.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, "invalid data", http.StatusBadRequest)
			return
		}

		if len(data) == 0 {
			http.Error(writer, "request body must not be empty", http.StatusBadRequest)
			return
		}

		duration, err := strconv.Atoi(string(data))
		if err != nil {
			http.Error(writer, "request body must contain duration", http.StatusBadRequest)
			return
		}

		roomAddr, err := app.NewRoom(request.Context(), duration, servPool, rdb)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		result, err := json.Marshal(Room{Address: roomAddr})
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusCreated)
		_, err = writer.Write(result)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
