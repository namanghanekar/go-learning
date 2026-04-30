package graph

import (
	"context"
	"graphql-service/graph/model"
	"ticket-system/shared/logger"
)

func (r *mutationResolver) LockSeat(ctx context.Context, seatID int32, userID string) (*string, error) {
	if err := r.Client.Lock(int(seatID), userID); err != nil {
		return nil, err
	}

	msg := "locked"
	return &msg, nil
}
func (r *mutationResolver) ConfirmSeat(ctx context.Context, seatID int32, userID string) (*string, error) {
	if err := r.Client.Pay(int(seatID), userID); err != nil {
		return nil, err
	}

	msg := "booked"
	return &msg, nil
}
func (r *queryResolver) Seats(ctx context.Context) ([]*model.Seat, error) {
	res, err := r.Client.GetSeats()
	if err != nil {
		return nil, err
	}

	var seats []*model.Seat

	for _, s := range res.Seats {

		var userID *string
		if s.UserId != "" {
			userID = &s.UserId
		}

		seats = append(seats, &model.Seat{
			ID:     s.Id,
			Status: s.Status,
			UserID: userID,
		})
	}
	return seats, nil
}
func (r *queryResolver) Logs(ctx context.Context) ([]string, error) {
	return logger.GetLogs(), nil
}
func (r *subscriptionResolver) SeatUpdates(ctx context.Context) (<-chan *model.Seat, error) {
	id, ch := r.AddSubscriber()
	go func() {
		<-ctx.Done()
		r.RemoveSubscriber(id)
	}()
	return ch, nil
}
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }
func (r *Resolver) Query() QueryResolver       { return &queryResolver{r} }

func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
