package webhooktomqtt

import (
	"io"
	"log"
	"net/http"
	"strings"
)

type mqttClient interface {
	Publish(topic string, payload []byte)
}

type webhooktomqttHandler struct {
	client mqttClient
}

func NewHandler(client mqttClient) http.Handler {
	return &webhooktomqttHandler{client: client}
}

func (h *webhooktomqttHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("received request: %s %s", r.Method, r.URL.Path)

	urlpath := strings.TrimPrefix(r.URL.Path, "/api/")

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("error reading request body: %s", err)
		http.Error(w, `{"status":"error"}`, http.StatusInternalServerError)
		return
	}

	h.client.Publish(urlpath, payload)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"status":"success"}`))
}
