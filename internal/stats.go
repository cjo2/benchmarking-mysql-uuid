package internal

import "sync"

type Stats struct {
	successfulInserts int
	failedInserts     int

	lock sync.Mutex
}

func (s *Stats) IncrementSuccessfulInserts() {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.successfulInserts++
}

func (s *Stats) IncrementFailedInserts() {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.failedInserts++
}

func (s *Stats) GetSuccessfulInserts() int {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.successfulInserts
}

func (s *Stats) GetFailedInserts() int {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.failedInserts
}
