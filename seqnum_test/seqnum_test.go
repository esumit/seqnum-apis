package seqnum_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/esumit/seqnum-apis/pkg/seqnum"
)

func handlerAdapter(f func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = f(w, r)
	}
}
func TestGetAPI(t *testing.T) {
	sm := seqnum.NewSeqnumManager()
	handler := seqnum.NewSeqnumApiRqHandler(sm)

	ts := httptest.NewServer(handlerAdapter(handler.Get))
	defer ts.Close()

	t.Run("basic test", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/seqnum")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK, got %v", resp.StatusCode)
		}
		bodyStr := string(body)
		t.Logf("Received response: %s", bodyStr)

		if !strings.Contains(string(body), "seq_num") {
			t.Errorf("Expected response to contain 'seq_num', got %v", string(body))
		}
	})

	t.Run("concurrency test", func(t *testing.T) {
		var wg sync.WaitGroup
		seqNums := make(map[int64]struct{})
		seqNumsMutex := sync.Mutex{}

		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				resp, err := http.Get(ts.URL + "/seqnum")
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
					return
				}

				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					t.Errorf("Failed to read response body: %v", err)
					return
				}

				var seqnumResp seqnum.SeqnumRs
				err = json.Unmarshal(body, &seqnumResp)
				if err != nil {
					t.Errorf("Failed to unmarshal JSON response: %v", err)
					return
				}

				seqNumsMutex.Lock()
				defer seqNumsMutex.Unlock()

				if _, ok := seqNums[seqnumResp.SeqNum]; ok {
					t.Errorf("Duplicate sequence number detected: %d", seqnumResp.SeqNum)
				} else {
					seqNums[seqnumResp.SeqNum] = struct{}{}
				}
			}()
		}

		wg.Wait()
	})
}
