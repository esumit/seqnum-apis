package seqnum

import (
	"encoding/json"
	"net/http"
)

type seqnumApiRqHandler struct {
	data DataManager
}

func NewSeqnumApiRqHandler(data DataManager) *seqnumApiRqHandler {
	return &seqnumApiRqHandler{data}
}

func (h *seqnumApiRqHandler) Get(w http.ResponseWriter, r *http.Request) error {
	// Handle HTTP GET Request
	ctx := r.Context()
	seqNumRs, err := h.data.GenerateSeqNum(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	// Return HTTP Response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(seqNumRs)
	return nil
}
