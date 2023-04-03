package gamecore

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// var sugarLogger *zap.SugaredLogger

var M = map[int]int{
	1: 0,
	2: 5,
	3: 3,
	4: 8,
	5: 0,
	6: 10,
}

func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func Play(w http.ResponseWriter, r *http.Request) {
	// http://127.0.0.1:<port>/play?Ante=1
	// bet_str := strings.TrimPrefix(r.URL.Path, "/play/")
	q := r.URL.Query()
	anteArr, ok := q["ante"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("<h1>no this field</h1>"))
		return
	}
	ante := strings.Join(anteArr, "")
	rand.Seed(time.Now().UnixNano())
	ran := randomInt(1, 7)
	// logrus.Debug("bet: %s, multiple: %v, random number: : %v", bet_str, M[ran], ran)
	fmt.Printf("bet: %v, random number: %v, multiple: %v \n", ante, ran, M[ran])
	bet, err := strconv.Atoi(ante)
	if err != nil {
		log.Fatal(err)
	}

	res := strconv.Itoa(bet * M[ran])

	fmt.Fprintf(w, res)
}
