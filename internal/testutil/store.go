package testutil

import (
	"context"
	"sync"
	"time"

	"github.com/privacybydesign/irma-meet/internal/meeting"
)

// MemStore is a meeting.Store that keeps everything in memory.
type MemStore struct {
	mu           sync.Mutex
	meetings     map[string]meeting.Meeting
	participants map[string][]string
}

func NewMemStore() *MemStore {
	return &MemStore{meetings: map[string]meeting.Meeting{}, participants: map[string][]string{}}
}

func (s *MemStore) Create(_ context.Context, m meeting.Meeting, emails []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.meetings[m.ID] = m
	s.participants[m.ID] = append([]string(nil), emails...)

	return nil
}

func (s *MemStore) Find(_ context.Context, id string) (meeting.Meeting, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	m, ok := s.meetings[id]
	if !ok {
		return meeting.Meeting{}, meeting.ErrNotFound
	}

	return m, nil
}

func (s *MemStore) SetPasswords(_ context.Context, id, moderator, attendee string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	m := s.meetings[id]
	if m.ModeratorPassword == "" {
		m.ModeratorPassword = moderator
	}

	if m.AttendeePassword == "" {
		m.AttendeePassword = attendee
	}

	s.meetings[id] = m

	return nil
}

func (s *MemStore) DeleteExpired(_ context.Context, cutoff time.Time) (meeting.Removed, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var removed meeting.Removed

	for id, m := range s.meetings {
		if m.StartTime.Before(cutoff) {
			delete(s.meetings, id)
			removed.Meetings++
		}
	}

	for id, emails := range s.participants {
		if _, alive := s.meetings[id]; !alive {
			removed.Participants += int64(len(emails))
			delete(s.participants, id)
		}
	}

	return removed, nil
}

// Insert adds a meeting as is, for tests that need one in a specific state.
func (s *MemStore) Insert(m meeting.Meeting, emails ...string) {
	_ = s.Create(context.Background(), m, emails)
}

// Meetings returns every stored meeting.
func (s *MemStore) Meetings() []meeting.Meeting {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := make([]meeting.Meeting, 0, len(s.meetings))
	for _, m := range s.meetings {
		out = append(out, m)
	}

	return out
}

// Participants returns the invited addresses of a meeting.
func (s *MemStore) Participants(id string) []string {
	s.mu.Lock()
	defer s.mu.Unlock()

	return append([]string(nil), s.participants[id]...)
}
