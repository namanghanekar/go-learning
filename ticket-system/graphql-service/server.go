package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/cors"

	"graphql-service/client"
	"graphql-service/graph"
	"graphql-service/graph/model"
	pb "ticket-system/proto"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	cl := client.NewClient()
	resolver := &graph.Resolver{
		Client: cl,
	}
	startSeatSubscriptionBridge(cl, resolver)

	srv := handler.New(
		graph.NewExecutableSchema(
			graph.Config{
				Resolvers: resolver,
			},
		),
	)

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
	})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)
	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir(findWebDir()))))

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	handler := c.Handler(http.DefaultServeMux)

	log.Printf("graphql service running on http://localhost:%s/", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

func findWebDir() string {
	candidates := []string{
		filepath.Join("..", "web"),
		"web",
	}

	for _, candidate := range candidates {
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return candidate
		}
	}

	return filepath.Join("..", "web")
}

func startSeatSubscriptionBridge(cl *client.Client, resolver *graph.Resolver) {
	go func() {
		for {
			stream, err := cl.Inventory.SeatUpdates(context.Background(), &pb.Empty{})
			if err != nil {
				log.Println("seat stream connect error:", err)
				time.Sleep(2 * time.Second)
				continue
			}

			for {
				seat, recvErr := stream.Recv()
				if recvErr != nil {
					log.Println("seat stream recv error:", recvErr)
					break
				}

				var userID *string
				if seat.UserId != "" {
					id := seat.UserId
					userID = &id
				}
				resolver.PublishSeatUpdate(&model.Seat{
					ID:     seat.Id,
					Status: seat.Status,
					UserID: userID,
				})
			}

			time.Sleep(1 * time.Second)
		}
	}()
}
