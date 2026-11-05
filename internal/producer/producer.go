// Package producer push events to consumers
package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

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

	byteArray, err := json.Marshal(event)
	if err != nil {
		log.Fatal(err)
	}

	seeds := []string{"localhost:9092"}
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(seeds...),
		kgo.ConsumeTopics("streamforge.events"),
	)
	if err != nil {
		panic(err)
	}
	defer cl.Close()

	ctx := context.Background()

	var wg sync.WaitGroup
	wg.Add(1)
	record := &kgo.Record{Topic: "streamforge.events", Value: byteArray}
	cl.Produce(ctx, record, func(_ *kgo.Record, err error) {
		defer wg.Done()
		if err != nil {
			fmt.Printf("record had a produce error: %v\n", err)
		}
	})
	wg.Wait()
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
