package model

import "sync"

type Subscriptions struct {
	mutex *sync.Mutex
	subs  map[string][]*Filters
}

func NewSubscriptions() *Subscriptions {
	return &Subscriptions{
		subs: make(map[string][]*Filters),
	}
}

func (subs Subscriptions) AddSubscription(subscriptionID string, filters []*Filters) {
	subs.mutex.Lock()
	defer subs.mutex.Unlock()
	subs.subs[subscriptionID] = filters
}

func (subs Subscriptions) CloseSubscription(subscriptionID string) {
	subs.mutex.Lock()
	defer subs.mutex.Unlock()
	delete(subs.subs, subscriptionID)
}

func (subs Subscriptions) MatchEvent(event Event) string {
	for subId, filters := range subs.subs {
		for _, filter := range filters {
			if filter.Match(event) {
				return subId
			}
		}
	}
	return ""
}
