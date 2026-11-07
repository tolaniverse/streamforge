// Package producer push events to consumers
package producer

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

func PushEvent(w http.ResponseWriter, r *http.Request) {
	var event Events
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	eventBytes, err := json.Marshal(event)
	if err != nil {
		log.Printf("%v\n", err)
		http.Error(w, "Could not parse event payload", http.StatusBadRequest)
		return
	}

	seeds := []string{"localhost:9092"}
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(seeds...),
		kgo.ConsumeTopics("streamforge.events"),
	)
	if err != nil {
		log.Printf("error on kafka %v\n", err)
		http.Error(w, "Unable to initialize client", http.StatusInternalServerError)
		return
	}
	defer cl.Close()

	ctx, cancel := context.WithTimeout(
		r.Context(),
		5*time.Second,
	)
	defer cancel()

	record := &kgo.Record{
		Topic: "streamforge.events",
		Value: eventBytes,
	}
	if err := cl.ProduceSync(ctx, record).FirstErr(); err != nil {
		log.Printf("err :: could not produce message |> %v\n", err)
		http.Error(
			w,
			"Unable to publish event",
			http.StatusServiceUnavailable,
		)
		return
	}

	err = json.NewEncoder(w).Encode(APIResponse{
		Status:  http.StatusAccepted,
		Message: "Event Pushed Successfully",
		Data:    map[string]string{"event_id": event.EventID},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
