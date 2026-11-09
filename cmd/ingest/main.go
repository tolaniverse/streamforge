package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/tolaniverse/streamforge/internal/producer"
	"github.com/tolaniverse/streamforge/internal/router"
	"github.com/twmb/franz-go/pkg/kgo"
)

func main() {
	seeds := []string{"localhost:9092"}
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(seeds...),
		kgo.ConsumeTopics("streamforge.events"),
	)
	if err != nil {
		log.Fatalf("error on kafka %v\n", err)
	}
	defer cl.Close()

	r := chi.NewRouter()
	ep := &producer.EventProducer{
		KafkaClient: cl,
	}
	router.Routes(r, ep)

	port := ":3000"
	fmt.Printf("Listening at port %s", port)
	http.ListenAndServe(":3000", r)
}
