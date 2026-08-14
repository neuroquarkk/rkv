package server

import (
	"encoding/json"
	"log"
	"net/http"
	"rkv/internal/registry"
	"strconv"
	"time"
)

func (s *Server) Heartbeat(w http.ResponseWriter, r *http.Request) {
	var data registry.MemberReq

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Printf("[REGISTRY] failed to decode body: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	r.Body.Close()

	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.state[data.Address]
	s.state[data.Address] = time.Now()
	if !exists {
		log.Printf("[REGISTRY] new member joined: %v\n", data.Address)
		s.tag++
	}

	w.WriteHeader(http.StatusAccepted)
}

func (s *Server) Members(w http.ResponseWriter, r *http.Request) {
	var tag uint64

	tagRaw := r.Header.Get("x-tag")
	if tagRaw != "" {
		parsedTag, err := strconv.ParseUint(tagRaw, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		tag = parsedTag
	}

	s.mu.RLock()

	currentTag := s.tag
	if tag == currentTag {
		s.mu.RUnlock()
		w.WriteHeader(http.StatusNotModified)
		return
	}

	members := make([]string, 0, len(s.state))
	for member := range s.state {
		members = append(members, member)
	}
	s.mu.RUnlock()

	resp := &registry.MembersResp{Members: members}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("x-tag", strconv.FormatUint(currentTag, 10))
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) Leave(w http.ResponseWriter, r *http.Request) {
	var data registry.MemberReq

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Printf("[REGISTRY] failed to decode body: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	r.Body.Close()

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.state[data.Address]; exists {
		log.Printf("[REGISTRY] member left: %s\n", data.Address)
		delete(s.state, data.Address)
		s.tag++
	}

	w.WriteHeader(http.StatusNoContent)
}
