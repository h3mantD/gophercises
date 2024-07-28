package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

type answer struct {
	ans       string
	isCorrect bool
}

type question struct {
	que string
	ans string
	answer
}

func (q *question) attempt() {
	fmt.Printf("%s ", q.que)
	fmt.Scanf("%s", &q.answer.ans)

	actualAns := strings.ToLower(strings.TrimSpace(q.ans))
	givenAns := strings.ToLower(strings.TrimSpace(q.answer.ans))

	q.answer.isCorrect = actualAns == givenAns
}

func parseProblemsFile(csvFile *string, shuffle *bool) []*question {
	questions := []*question{}

	f, err := os.ReadFile(*csvFile)
	if err != nil {
		log.Fatal(err.Error())
	}

	csvReader := csv.NewReader(bytes.NewReader(f))

	for {
		line, err := csvReader.Read()
		if err != nil {
			break
		}

		questions = append(questions, &question{
			que: line[0],
			ans: strings.TrimSpace(line[1]),
		})
	}

	if *shuffle {
		r := rand.New(rand.NewSource(time.Now().UnixMicro()))
		r.Shuffle(len(questions), func(i, j int) {
			questions[i], questions[j] = questions[j], questions[i]
		})
	}

	return questions
}

func startTest(questions []*question, ch chan<- bool) {
	for _, question := range questions {
		question.attempt()
	}

	ch <- true
}

func showResult(questions []*question) {
	totalQ := len(questions)
	totalM := 0
	for _, q := range questions {
		if q.answer.isCorrect {
			totalM++
		}
	}

	fmt.Printf("\n\nYou have scored %d out of %d!\n", totalM, totalQ)
}

func main() {
	problemsFile := flag.String("problems", "problems.csv", "A CSV file containing the questions along with answers.")
	timeLimit := flag.Int64("limit", 30, "Time limit for completing the quiz.")
	shuffle := flag.Bool("shuffle", false, "Whether to shuffle the question set or not.")

	flag.Parse()

	questions := parseProblemsFile(problemsFile, shuffle)

	var conf string
	fmt.Print("Are you sure, you want to start the test? [y/n] : ")
	fmt.Scanf("%s", &conf)

	if conf != "y" {
		fmt.Println("TEST ABORTED!")
		return
	}

	ticker := time.NewTicker(time.Second * time.Duration(*timeLimit))

	ch := make(chan bool)
	go startTest(questions, ch)

	select {
	case <-ticker.C:
		showResult(questions)
	case <-ch:
		showResult(questions)
	}
}
