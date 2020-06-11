package main

import (
	"encoding/json"
	"errors"
	"github.com/freeeve/uci"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/move", ChessServer)
	http.ListenAndServe(":8081", nil)
}

func ChessServer(w http.ResponseWriter, r *http.Request) {
	var responseObj interface{}

	w.Header().Set("Content-Type", "application/json")


	if _, ok := r.URL.Query()["game"];ok {
		if fenString, ok := r.URL.Query()["fen"]; ok {
			result, err := GetStockfishResults(fenString[0])

			if err == nil {
				responseObj = result
			} else {
				responseObj = err
			}
		} else {
			responseObj = errors.New("fen parameter is missing")
		}
	} else {
		responseObj = errors.New("game parameter is missing")
	}

	if errObj, ok := responseObj.(error); ok {
		responseObj = map[string]string {
			"errorMessage": errObj.Error(),
		}
	}

	bytes, err := json.Marshal(responseObj)

	if err != nil {
		log.Fatal(err)

		w.WriteHeader(500)
	} else {
		w.WriteHeader(200)
		w.Write(bytes)
	}
}
func GetStockfishResults(fenString string) (result *uci.Results, err error) {
	eng, err := uci.NewEngine(os.Getenv("STOCKFISH_PATH"))
	if err == nil {
		// set some engine options
		eng.SetOptions(uci.Options{
			Hash:    128,
			Ponder:  false,
			OwnBook: true,
			MultiPV: 4,
		})

		// set the starting position
		eng.SetFEN(fenString)

		// set some result filter options
		resultOpts := uci.HighestDepthOnly | uci.IncludeUpperbounds | uci.IncludeLowerbounds
		result, err = eng.GoDepth(10, resultOpts)

		eng.Close()
	}

	return result, err
}