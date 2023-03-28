package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type Cotacao struct {
	Code       string
	Codein     string
	Name       string
	High       string
	Low        string
	VarBid     string
	PctChange  string
	Bid        string
	Ask        string
	Timestamp  string
	CreateDate string
}

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
	db, err := gorm.Open(sqlite.Open("cotacao.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&Cotacao{})

	http.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		GetCotacaoHandler(w, r, db)
	})
	http.ListenAndServe(":8080", nil)
}

func GetCotacaoHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	ctxAPI := context.Background()
	ctxAPI, cancelAPI := context.WithTimeout(ctxAPI, time.Millisecond*200)
	defer cancelAPI()

	req, err := http.NewRequestWithContext(ctxAPI, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
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

	var cotacoes []Cotacao
	db.Find(&cotacoes)
	fmt.Println(cotacoes)

	var cotacaoAPI CotacaoJSON
	json.Unmarshal(body, &cotacaoAPI)

	SaveCotacao(cotacaoAPI, db)

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func SaveCotacao(body CotacaoJSON, db *gorm.DB) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*10)
	defer cancel()

	db.WithContext(ctx).Create(&Cotacao{
		Code:       body.USDBRL.Code,
		Codein:     body.USDBRL.Codein,
		Name:       body.USDBRL.Name,
		High:       body.USDBRL.High,
		Low:        body.USDBRL.Low,
		VarBid:     body.USDBRL.VarBid,
		PctChange:  body.USDBRL.PctChange,
		Bid:        body.USDBRL.Bid,
		Ask:        body.USDBRL.Ask,
		Timestamp:  body.USDBRL.Timestamp,
		CreateDate: body.USDBRL.CreateDate,
	})
}
