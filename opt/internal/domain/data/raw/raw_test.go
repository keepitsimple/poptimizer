package raw

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/WLM1ke/poptimizer/opt/internal/domain"
	"github.com/WLM1ke/poptimizer/opt/internal/domain/data"
	"github.com/stretchr/testify/assert"
)

type testRepo struct {
	qid    domain.QID
	entity Table
}

func (r testRepo) Get(_ context.Context, qid domain.QID) (domain.Aggregate[Table], error) {
	agg := domain.Aggregate[Table]{}

	if qid != r.qid {
		return agg, fmt.Errorf("repo error")
	}

	agg.Entity = r.entity

	return agg, nil
}

type testPublisher struct {
	events []domain.Event
}

func (t *testPublisher) Publish(event domain.Event) {
	t.events = append(t.events, event)
}

func TestCheckRawHandler_Handle(t *testing.T) {
	t.Parallel()

	repo := testRepo{
		qid: domain.QID{
			Sub:   data.Subdomain,
			Group: _rawGroup,
			ID:    "ABRD",
		},
		entity: Table{
			{
				Date:     time.Date(2021, time.July, 12, 0, 0, 0, 0, time.UTC),
				Value:    2.86,
				Currency: "RUR",
			},
			{
				Date:     time.Date(2022, time.July, 8, 0, 0, 0, 0, time.UTC),
				Value:    3.44,
				Currency: "RUR",
			},
		},
	}

	pub := testPublisher{}

	handler := NewCheckRawHandler(&pub, repo)

	testTable := []struct {
		divDate time.Time
		event   int
	}{
		{time.Date(2022, time.July, 7, 0, 0, 0, 0, time.UTC), 1},
		{time.Date(2022, time.July, 8, 0, 0, 0, 0, time.UTC), 0},
		{time.Date(2022, time.July, 9, 0, 0, 0, 0, time.UTC), 1},
	}

	for _, row := range testTable {
		event := domain.Event{
			QID: domain.QID{
				Sub:   data.Subdomain,
				Group: _statusGroup,
				ID:    "ABRD",
			},
			Timestamp: time.Date(2022, time.July, 4, 0, 0, 0, 0, time.UTC),
			Data: Status{
				Ticker: "ABRD",
				Date:   row.divDate,
			},
		}

		handler.Handle(context.Background(), event)

		assert.Equal(t, row.event, len(pub.events), "must be 1 event %d", len(pub.events))

		if row.event > 0 {
			assert.Equal(
				t,
				domain.QID{Sub: data.Subdomain, Group: _rawGroup, ID: "ABRD"},
				pub.events[0].QID,
				"incorrect event id",
			)
			assert.Equal(
				t,
				time.Date(2022, time.July, 4, 0, 0, 0, 0, time.UTC),
				pub.events[0].Timestamp,
				"incorrect event timestamp",
			)
			assert.ErrorContains(
				t,
				pub.events[0].Data.(error),
				"missed dividend at",
			)
		}

		pub.events = nil
	}
}

func TestCheckRawHandler_Handle_DataErr(t *testing.T) {
	t.Parallel()

	repo := testRepo{}
	pub := testPublisher{}

	handler := NewCheckRawHandler(&pub, repo)
	handler.Handle(context.Background(), domain.Event{})

	assert.Equal(t, 1, len(pub.events), "no error on incorrect event data")
	assert.ErrorContains(t, pub.events[0].Data.(error), "can't parse Event", "incorrect error")
}

func TestCheckRawHandler_Handle_RepoErr(t *testing.T) {
	t.Parallel()

	repo := testRepo{}
	pub := testPublisher{}

	handler := NewCheckRawHandler(&pub, repo)
	handler.Handle(context.Background(), domain.Event{Data: Status{}})

	assert.Equal(t, 1, len(pub.events), "no error with bad repo")
	assert.ErrorContains(t, pub.events[0].Data.(error), "repo error", "incorrect error")
}
