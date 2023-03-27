package seqnum

import (
	"context"
)

type DataManager interface {
	GenerateSeqNum(ctx context.Context) (*SeqnumRs, error)
}

type SeqnumRq struct {
}

type SeqnumRs struct {
	SeqNum int64 `json:"seq_num"`
}

type SeqNumAPIs interface {
	Get(ctx context.Context) (*SeqnumRs, error)
}
