package server

import (
	"encoding/json"
	"log"
	"net/http"
	"rkv/internal/registry"
	"time"
)

func (s *Server) Heartbeat(w http.ResponseWriter, r *http.Request) {
	var data registry.HeartbeatReq

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Printf("[REGISTRY] failed to decode body: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	r.Body.Close()

	s.mu.Lock()
	_, exists := s.state[data.Address]
	s.state[data.Address] = time.Now()
	s.mu.Unlock()
	if !exists {
		log.Printf("[REGISTRY] new member joined: %v\n", data.Address)
	}

	w.WriteHeader(http.StatusAccepted)
}

func (s *Server) Members(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	members := make([]string, 0, len(s.state))
	for member := range s.state {
		members = append(members, member)
	}
	s.mu.RUnlock()

	resp := &registry.MembersResp{Members: members}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
