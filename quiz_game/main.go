package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	flagFilePath string
	flagTime     int
	flagShuffle  bool
)

// Question 問題と回答の構造体
type Question struct {
	Quiz string
	Ans  string
}

func init() {
	flag.StringVar(&flagFilePath, "csv", "problems.csv", "a csv file in the format of \"question\"")
	flag.IntVar(&flagTime, "limit", 30, "the time limit for the quiz in seconds")
	flag.BoolVar(&flagShuffle, "shuffle", true, "shuffle quiz")
	flag.Parse()
}

func main() {
	csvPath, err := filepath.Abs(flagFilePath)
	if err != nil {
		log.Fatalln(err)
	}
	file, err := os.Open(csvPath)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	csvData, err := csvReader.ReadAll()
	if err != nil {
		log.Fatalln(err)
	}

	numQuestions := len(csvData)
	questions := make([]Question, numQuestions)
	for i, data := range csvData {
		var q Question
		q.Quiz = data[0]
		q.Ans = data[1]
		questions[i] = q
	}

	fmt.Println("Press Any to start Game!")
	bufio.NewScanner(os.Stdout).Scan()

	if flagShuffle {
		rand.Seed(time.Now().UTC().UnixNano())
	}
	quizIndex := rand.Perm(numQuestions)

	correct, wrong := 0, 0
	respondTo := make(chan bool)
	timeUp := time.After(time.Second * time.Duration(flagTime))

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
	LOOP:
		for _, i := range quizIndex {
			go askQuestion(questions[i], respondTo)
			select {
			case <-timeUp:
				break LOOP
			case ans := <-respondTo:
				if ans {
					correct++
				} else {
					wrong++
				}
			}
		}
		wg.Done()
	}()
	wg.Wait()

	fmt.Printf("You correctly/answered %v/%v, the number of questions %v)\n",
		correct, correct+wrong, numQuestions)

}

func askQuestion(q Question, replyTo chan bool) {
	fmt.Printf("Question: %v\n", q.Quiz)
	fmt.Printf("Answer: ")
	var ans string
	fmt.Scan(&ans)
	if ans == q.Ans {
		replyTo <- true
	} else {
		replyTo <- false
	}
}
