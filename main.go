package main

import (
	"encoding/json"
	"errors"
	"github.com/freeeve/uci"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type EngineWrapper struct {
	Engine *uci.Engine
	LastAccessed time.Time
}

func main() {
	http.HandleFunc("/move", ChessServer)
	http.ListenAndServe(":8081", nil)
}

func ChessServer(w http.ResponseWriter, r *http.Request) {
	var responseObj interface{}

	w.Header().Set("Content-Type", "application/json")

	if game, ok := r.URL.Query()["game"];ok {
		if fenString, ok := r.URL.Query()["fen"]; ok {
			result, err := GetStockfishResults(strings.Join(game, " "), fenString[0])

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

var engines map[string]EngineWrapper = map[string]EngineWrapper{}

func GetEngine(gameID string) (engine *uci.Engine, err error) {
	if wrapper, ok := engines[gameID]; ok {
		wrapper.LastAccessed = time.Now()
		engine = wrapper.Engine
	} else {
		engine, err = uci.NewEngine(os.Getenv("STOCKFISH_PATH"))
		if err == nil {
			// set some engine options
			engine.SetOptions(uci.Options{
				Hash:    128,
				Ponder:  true,
				OwnBook: true,
				MultiPV: 4,
				Threads: 10,
			})

			engine.SendOption("UCI_LimitStrength", 800)

			wrapper := EngineWrapper{
				Engine:       engine,
				LastAccessed: time.Now(),
			}
			engines[gameID] = wrapper
		}
	}

	go func(gameID string) {
		time.Sleep(10 * time.Minute)
		if wrapper, ok := engines[gameID]; ok {
			if wrapper.LastAccessed.Add(5 * time.Minute).Before(time.Now()) {
				delete(engines, gameID)
				wrapper.Engine.Close()
			}
		}
	}(gameID)

	return engine, err
}

func GetStockfishResults(gameID string, fenString string) (result *uci.Results, err error) {
	eng, err := GetEngine(gameID)
	if err == nil {
		// set the starting position
		eng.SetFEN(fenString)

		// set some result filter options
		resultOpts := uci.HighestDepthOnly | uci.IncludeUpperbounds | uci.IncludeLowerbounds
		result, err = eng.GoDepth(10, resultOpts)
	}

	return result, err
}