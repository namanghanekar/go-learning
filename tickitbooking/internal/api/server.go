package api

import (
	"context"
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"
	"worldtour-tickets/internal/contracts"
	"worldtour-tickets/internal/pubsub"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/graphql-go/graphql"
)

type Server struct {
	inventory contracts.InventoryServiceClient
	payment   contracts.PaymentServiceClient
	events    *pubsub.Hub[contracts.SeatEvent]
	schema    graphql.Schema
	roomID    string
}

func NewServer(inventory contracts.InventoryServiceClient, payment contracts.PaymentServiceClient, roomID string) (*Server, error) {
	s := &Server{
		inventory: inventory,
		payment:   payment,
		events:    pubsub.NewHub[contracts.SeatEvent](),
		roomID:    roomID,
	}
	schema, err := s.buildSchema()
	if err != nil {
		return nil, err
	}
	s.schema = schema
	return s, nil
}

func (s *Server) StartInventoryListener(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			stream, err := s.inventory.OpenSeatStream(ctx)
			if err != nil {
				log.Printf("api inventory stream connect failed: %v", err)
				time.Sleep(2 * time.Second)
				continue
			}
			for {
				msg, recvErr := stream.Recv()
				if recvErr != nil {
					log.Printf("api inventory stream disconnected: %v", recvErr)
					break
				}
				if msg.Event != nil {
					s.events.Publish(*msg.Event)
				}
			}
			time.Sleep(2 * time.Second)
		}
	}()
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/graphql", s.handleGraphQL)
	mux.HandleFunc("/graphql/subscribe", s.handleSubscription)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	if staticFS := frontendFS(); staticFS != nil {
		mux.Handle("/", http.FileServer(http.FS(staticFS)))
	}
	return mux
}

func frontendFS() fs.FS {
	candidates := []string{
		"web",
		filepath.Join("..", "..", "web"),
	}
	for _, candidate := range candidates {
		if _, err := os.Stat(filepath.Join(candidate, "index.html")); err == nil {
			return os.DirFS(candidate)
		}
	}
	return nil
}

func (s *Server) buildSchema() (graphql.Schema, error) {
	seatType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Seat",
		Fields: graphql.Fields{
			"id":        &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"roomId":    &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"number":    &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"status":    &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"heldBy":    &graphql.Field{Type: graphql.String},
			"holdToken": &graphql.Field{Type: graphql.String},
			"expiresAt": &graphql.Field{Type: graphql.String},
			"updatedAt": &graphql.Field{Type: graphql.String},
		},
	})

	query := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"seats": &graphql.Field{
				Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(seatType))),
				Args: graphql.FieldConfigArgument{
					"roomId": &graphql.ArgumentConfig{Type: graphql.String},
				},
				Resolve: func(p graphql.ResolveParams) (any, error) {
					roomID, _ := p.Args["roomId"].(string)
					if roomID == "" {
						roomID = s.roomID
					}
					ctx, cancel := context.WithTimeout(p.Context, 3*time.Second)
					defer cancel()
					resp, err := s.inventory.ListSeats(ctx, &contracts.ListSeatsRequest{RoomID: roomID})
					if err != nil {
						return nil, err
					}
					return toGraphQLSeats(resp.Seats), nil
				},
			},
		},
	})

	mutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"holdSeat": &graphql.Field{
				Type: seatType,
				Args: graphql.FieldConfigArgument{
					"roomId":      &graphql.ArgumentConfig{Type: graphql.String},
					"seatNumber":  &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"userId":      &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"amountCents": &graphql.ArgumentConfig{Type: graphql.Int},
				},
				Resolve: func(p graphql.ResolveParams) (any, error) {
					roomID, _ := p.Args["roomId"].(string)
					if roomID == "" {
						roomID = s.roomID
					}
					seatNumber := p.Args["seatNumber"].(string)
					userID := p.Args["userId"].(string)
					amountCents, _ := p.Args["amountCents"].(int)
					holdToken := uuid.NewString()

					holdCtx, cancel := context.WithTimeout(p.Context, 3*time.Second)
					defer cancel()
					holdResp, err := s.inventory.HoldSeat(holdCtx, &contracts.HoldSeatRequest{
						RoomID:     roomID,
						SeatNumber: seatNumber,
						UserID:     userID,
						HoldToken:  holdToken,
						TTLSeconds: 300,
					})
					if err != nil {
						return nil, err
					}

					registerCtx, registerCancel := context.WithTimeout(p.Context, 3*time.Second)
					defer registerCancel()
					if _, err := s.payment.RegisterHold(registerCtx, &contracts.RegisterHoldRequest{
						SeatID:      holdResp.Seat.ID,
						SeatNumber:  holdResp.Seat.Number,
						RoomID:      holdResp.Seat.RoomID,
						HoldToken:   holdResp.Seat.HoldToken,
						UserID:      userID,
						ExpiresAt:   holdResp.Seat.ExpiresAt,
						AmountCents: int64(amountCents),
					}); err != nil {
						return nil, err
					}

					return toGraphQLSeat(holdResp.Seat), nil
				},
			},
			"payForSeat": &graphql.Field{
				Type: graphql.NewObject(graphql.ObjectConfig{
					Name: "PaymentResult",
					Fields: graphql.Fields{
						"accepted":   &graphql.Field{Type: graphql.NewNonNull(graphql.Boolean)},
						"message":    &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
						"paymentRef": &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
					},
				}),
				Args: graphql.FieldConfigArgument{
					"seatId":      &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"holdToken":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"userId":      &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"amountCents": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				},
				Resolve: func(p graphql.ResolveParams) (any, error) {
					ctx, cancel := context.WithTimeout(p.Context, 3*time.Second)
					defer cancel()
					resp, err := s.payment.SubmitPayment(ctx, &contracts.SubmitPaymentRequest{
						SeatID:      p.Args["seatId"].(string),
						HoldToken:   p.Args["holdToken"].(string),
						UserID:      p.Args["userId"].(string),
						AmountCents: int64(p.Args["amountCents"].(int)),
					})
					if err != nil {
						return nil, err
					}
					return map[string]any{
						"accepted":   resp.Accepted,
						"message":    resp.Message,
						"paymentRef": resp.PaymentRef,
					}, nil
				},
			},
		},
	})

	return graphql.NewSchema(graphql.SchemaConfig{
		Query:    query,
		Mutation: mutation,
	})
}

func (s *Server) handleGraphQL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Query         string         `json:"query"`
		OperationName string         `json:"operationName"`
		Variables     map[string]any `json:"variables"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result := graphql.Do(graphql.Params{
		Schema:         s.schema,
		RequestString:  req.Query,
		OperationName:  req.OperationName,
		VariableValues: req.Variables,
		Context:        r.Context(),
	})
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}

func (s *Server) handleSubscription(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	var subMu sync.Mutex
	var sub chan contracts.SeatEvent

	for {
		var msg map[string]any
		if err := conn.ReadJSON(&msg); err != nil {
			return
		}
		msgType, _ := msg["type"].(string)
		switch msgType {
		case "connection_init":
			_ = conn.WriteJSON(map[string]any{"type": "connection_ack"})
		case "subscribe":
			payload, _ := msg["payload"].(map[string]any)
			query, _ := payload["query"].(string)
			roomID := parseSubscriptionRoomID(query, s.roomID)
			subMu.Lock()
			if sub != nil {
				s.events.Unsubscribe(sub)
			}
			sub = s.events.Subscribe(64)
			subMu.Unlock()

			subscriptionID, _ := msg["id"].(string)
			go func(room string, id string, subCh chan contracts.SeatEvent) {
				for event := range subCh {
					if room != "" && event.Seat.RoomID != room {
						continue
					}
					_ = conn.WriteJSON(map[string]any{
						"id":   id,
						"type": "next",
						"payload": map[string]any{
							"data": map[string]any{
								"seatUpdates": toGraphQLSeat(event.Seat),
							},
						},
					})
				}
			}(roomID, subscriptionID, sub)
		case "complete":
			subMu.Lock()
			if sub != nil {
				s.events.Unsubscribe(sub)
				sub = nil
			}
			subMu.Unlock()
		}
	}
}

func parseSubscriptionRoomID(query, fallback string) string {
	re := regexp.MustCompile(`roomId\s*:\s*"([^"]+)"`)
	matches := re.FindStringSubmatch(query)
	if len(matches) == 2 {
		return matches[1]
	}
	return fallback
}

func toGraphQLSeats(seats []contracts.SeatDTO) []map[string]any {
	out := make([]map[string]any, 0, len(seats))
	for _, seat := range seats {
		out = append(out, toGraphQLSeat(seat))
	}
	return out
}

func toGraphQLSeat(seat contracts.SeatDTO) map[string]any {
	return map[string]any{
		"id":        seat.ID,
		"roomId":    seat.RoomID,
		"number":    seat.Number,
		"status":    seat.Status,
		"heldBy":    seat.HeldBy,
		"holdToken": seat.HoldToken,
		"expiresAt": formatTime(seat.ExpiresAt),
		"updatedAt": formatTime(seat.UpdatedAt),
	}
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format(time.RFC3339)
}
