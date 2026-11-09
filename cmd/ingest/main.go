package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/tolaniverse/streamforge/internal/producer"
	"github.com/tolaniverse/streamforge/internal/router"
	"github.com/twmb/franz-go/pkg/kgo"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

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

	srv := &http.Server{
		Addr:    ":3000",
		Handler: r,
	}

	go func() {
		log.Println("Server listening on :3000")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutdown signal received, draining connections...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Forced shutdown due to error: %v", err)
	}

	log.Println("Server exited cleanly")
}
