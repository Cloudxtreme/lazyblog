package model

import (
	"encoding/json"
	"sync"
)

// MemoryStore is the in-memory database we use.
type MemoryStore struct {
	lock    sync.RWMutex
	htmlMap map[string][]byte
	jsonMap map[string][]byte
}

// NewMemoryStore returns a new in-memory data store.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		htmlMap: make(map[string][]byte),
		jsonMap: make(map[string][]byte),
	}
}

// SetPost sets a post in the in-memory data store.
func (s *MemoryStore) SetPost(p *Post) ([]byte, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	post, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return nil, err
	}
	s.jsonMap[string(p.ID)] = post
	s.htmlMap[string(p.ID)] = []byte(p.Body)

	return p.ID, nil
}

// GetPostHTML returns the post HTML for the given post ID from the in-memory
// data store.
func (s *MemoryStore) GetPostHTML(id []byte) ([]byte, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.htmlMap[string(id)], nil
}

// GetPostJSON returns the post JSON for the given post ID from the in-memory
// data store.
func (s *MemoryStore) GetPostJSON(id []byte) ([]byte, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.jsonMap[string(id)], nil
}

// GetPosts returns all Posts from the in-memory data store.
func (s *MemoryStore) GetPosts(num, offset int) ([]*Post, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	// warning: slow implementation
	var posts []*Post
	for _, p := range s.jsonMap {
		var post *Post
		err := json.Unmarshal(p, &post)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

// SetUser is not implemented.
func (s *MemoryStore) SetUser(username, password []byte) error {
	return nil
}

// GetUser is not implemented.
func (s *MemoryStore) GetUser(username, password []byte) error {
	return nil
}
