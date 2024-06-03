package parser

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var (
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// decode req
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {

		}
		req := Request{}
		err = json.Unmarshal(reqBody, &req)
		if err != nil {

		}
		switch req.Method {
		case "eth_blockNumber":
			m := GetCurrentBlockResponse{
				Result: "0x11",
			}
			json.NewEncoder(w).Encode(m)
		case "eth_getLogs":
			m := GetLogsResponse{
				Result: []struct {
					TransactionHash string `json:"transactionHash"`
				}{
					{TransactionHash: "0x0ce1dd8f41038a9549ea2b63a0f5d671e665d9fb49fdb57b46208dff96e9882b"},
					{TransactionHash: "0xf85073685a92c2bff43af67faa72bbd3cec3a026dd73ecfcdfc841c3e3e37b77"},
				},
			}
			json.NewEncoder(w).Encode(m)
		case "eth_getTransactionByHash":
			m := GetTransactionByHashResponse{
				Result: Transaction{
					From: "0x101",
					To:   "0x102",
				},
			}
			json.NewEncoder(w).Encode(m)
		}
	}))
)

func Test_GetCurrentBlock(t *testing.T) {
	eth := NewEthereum(server.URL, 2*time.Second)
	tt := []struct {
		name   string
		expect int
		err    bool
	}{
		{
			expect: 17,
			err:    false,
		},
	}

	for _, tc := range tt {
		currentBlock, err := eth.GetCurrentBlock()
		if err != nil {
			t.Fatal(err)
		}
		if currentBlock != tc.expect {
			t.Fatalf("expect %d, got: %d", tc.expect, currentBlock)
		}
	}
}

func Test_GetTransactions(t *testing.T) {
	eth := NewEthereum(server.URL, 2*time.Second)
	tt := []struct {
		name string
		txs  int
	}{
		{
			txs: 2,
		},
	}
	for _, tc := range tt {
		txs, _ := eth.GetTransactions("0x11110000")
		if len(txs) != tc.txs {
			t.Fatalf("expect %d, got: %d", tc.txs, len(txs))
		}
	}
}
