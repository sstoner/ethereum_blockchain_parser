package parser

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// Ethereum ...
type Ethereum struct {
	endpoint string
	client   *http.Client
}

// NewEthereum ...
func NewEthereum(endpoint string, timeout time.Duration) *Ethereum {
	return &Ethereum{
		endpoint: endpoint,
		// TODO Client build from Transport
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// GetCurrentBlock ...
func (e *Ethereum) GetCurrentBlock() (int, error) {
	reqBody := GetCurrentBlockRequest{
		Request: Request{
			Jsonrpc: "2.0",
			Method:  "eth_blockNumber",
			// TODO build ID
			ID: rand.Intn(math.MaxInt32),
		},
	}
	jsonReq, err := json.Marshal(reqBody)
	if err != nil {
		return -1, err
	}
	rd := bytes.NewReader(jsonReq)
	req, err := http.NewRequest(http.MethodPost, e.endpoint, rd)
	resp, err := e.client.Do(req)
	if err != nil {
		return -1, err
	}
	// decode resp
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	result := GetCurrentBlockResponse{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return -1, err
	}

	blockNum, err := strconv.ParseInt(result.Result, 0, 64)
	if err != nil {
		return -1, err
	}

	return int(blockNum), nil
}

// GetTransactions ...
func (e *Ethereum) GetTransactions(address string) ([]Transaction, error) {
	// get logs
	reqBody := GetLogsRequest{
		Request: Request{
			Jsonrpc: "2.0",
			Method:  "eth_getLogs",
			ID:      rand.Intn(math.MaxInt32),
		},
		Params: []struct {
			Address []string `json:"addresses"`
		}{
			{
				Address: []string{address},
			},
		},
	}

	jsonReq, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	rd := bytes.NewReader(jsonReq)
	req, err := http.NewRequest(http.MethodPost, e.endpoint, rd)
	if err != nil {
		return nil, err
	}

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, err
	}

	// decode resp
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	result := GetLogsResponse{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	// get transactions
	var transactions []Transaction
	for _, log := range result.Result {
		transaction, err := e.getTransactionByHash(log.TransactionHash)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, *transaction)
	}
	return transactions, nil
}

// GetTransactionByHash ...
func (e *Ethereum) getTransactionByHash(txHash string) (*Transaction, error) {
	reqBody := GetTransactionByHashRequest{
		Request: Request{
			Jsonrpc: "2.0",
			Method:  "eth_getTransactionByHash",
			ID:      rand.Intn(math.MaxInt32),
		},
		Params: []string{
			txHash,
		},
	}

	jsonReq, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	rd := bytes.NewReader(jsonReq)
	req, err := http.NewRequest(http.MethodPost, e.endpoint, rd)
	if err != nil {
		return nil, err
	}
	resp, err := e.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	result := GetTransactionByHashResponse{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result.Result, nil
}

// Request ...
type Request struct {
	Jsonrpc string `json:"jsonrpc"`
	// eth_blockNumber, eth_getLogs, eth_getTransactionByHash
	Method string `json:"method"`
	ID     int    `json:"id"`
}

// Response ...
type Response struct {
	ID      int    `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
}

// GetCurrentBlockRequest ...
type GetCurrentBlockRequest struct {
	Request
}

// GetCurrentBlockResponse ...
type GetCurrentBlockResponse struct {
	Response
	// blockNum is a hex string
	Result string `json:"result"`
}

// GetLogsRequest ...
type GetLogsRequest struct {
	Request
	Params []struct {
		Address []string `json:"addresses"`
	} `json:"params"`
}

// GetLogsResponse ...
type GetLogsResponse struct {
	Response
	Result []struct {
		TransactionHash string `json:"transactionHash"`
	} `json:"result"`
}

// GetTransactionByHashRequest ...
type GetTransactionByHashRequest struct {
	Request
	Params []string `json:"params"`
}

// GetTransactionByHashResponse ...
type GetTransactionByHashResponse struct {
	Response
	Result Transaction `json:"result"`
}

// Transaction ...
type Transaction struct {
	BlockHash        string `json:"blockHash"`
	BlockNumber      string `json:"blockNumber"`
	From             string `json:"from"`
	Gas              string `json:"gas"`
	GasPrice         string `json:"gasPrice"`
	Hash             string `json:"hash"`
	Input            string `json:"input"`
	Nonce            string `json:"nonce"`
	To               string `json:"to"`
	TransactionIndex string `json:"transactionIndex"`
	Value            string `json:"value"`
	V                string `json:"v"`
	R                string `json:"r"`
	S                string `json:"s"`
}
