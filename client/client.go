package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type CotacaoJSON struct {
	USDBRL struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var cotacaoAPI CotacaoJSON
	json.Unmarshal(body, &cotacaoAPI)

	var f *os.File
	if _, err := os.Stat("cotacao.txt"); errors.Is(err, os.ErrNotExist) {
		f, _ = os.Create("cotacao.txt")
	} else {
		f, _ = os.OpenFile("cotacao.txt", os.O_RDWR|os.O_APPEND, 0660)
	}

	defer f.Close()
	_, err = f.Write([]byte(fmt.Sprintf("\nDólar: %v", cotacaoAPI.USDBRL.Bid)))
	if err != nil {
		panic(err)
	}

	fmt.Println("Created")
}
