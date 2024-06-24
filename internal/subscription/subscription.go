package subscription

import (
	"sync"

	"github.com/dkrasnykh/graphql-app/graph/model"
)

type Subscription struct {
	mu sync.RWMutex

	counter int64
	chs     map[int64]chan *model.Comment
	// for each post save set of subscribers
	// map[postID]map[subsribtionID]true
	postObservers map[int64]map[int64]bool
	// save subscribtion posts
	subscription map[int64][]int64
}

func New() *Subscription {
	return &Subscription{
		mu:            sync.RWMutex{},
		counter:       1,
		chs:           make(map[int64]chan *model.Comment),
		postObservers: make(map[int64]map[int64]bool),
		subscription:  make(map[int64][]int64),
	}
}

func (s *Subscription) Add(posts []int64) (int64, chan *model.Comment) {
	s.mu.Lock()
	defer s.mu.Unlock()

	subscriptionID := s.counter
	s.counter += 1

	ch := make(chan *model.Comment)
	s.chs[subscriptionID] = ch
	s.subscription[subscriptionID] = make([]int64, 0, len(posts))
	for _, post := range posts {
		if _, ok := s.postObservers[post]; !ok {
			s.postObservers[post] = make(map[int64]bool)
		}
		s.postObservers[post][subscriptionID] = true
		s.subscription[subscriptionID] = append(s.subscription[subscriptionID], post)
	}
	// returns subsribtion id and chan with updates from server
	return subscriptionID, ch
}

func (s *Subscription) Delete(subscriptionID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// clear data for every post
	for _, post := range s.subscription[subscriptionID] {
		delete(s.postObservers[post], subscriptionID)
		// If there are no more subscriptions for a post, then delete the post
		if len(s.postObservers[post]) == 0 {
			delete(s.postObservers, post)
		}
	}
	delete(s.chs, subscriptionID)
	delete(s.subscription, subscriptionID)
}

func (s *Subscription) Broadcast(postID int64, comment *model.Comment) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, ok := s.postObservers[postID]; !ok {
		return
	}
	// send update for every subscribtion
	for subscriptionID, _ := range s.postObservers[postID] {
		s.chs[subscriptionID] <- comment
	}
}
