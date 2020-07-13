package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type quizData struct {
	question string
	answer string 
}

func main() {
	var quizFileName = flag.String("f", "problems.csv", "usage: filename for quiz")
	var quizTime = flag.Int("t", 30, "usage: time duration for quiz")
	flag.Parse()

	quiz := parseQuizData(*quizFileName)
	
	inputReader := bufio.NewReader(os.Stdin)
	gameTimer := time.Duration(*quizTime)

	numCorrect := runGame(quiz, inputReader, gameTimer)
	fmt.Printf("Out of %d questions , you got %d correct!\n", len(quiz), numCorrect)
}

func newQuizData(question string, answer string) quizData {
	return quizData { question, answer }
}

func parseQuizData(quizFileName string) []quizData {
	quizFile, err := os.Open(quizFileName)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
		os.Exit(1)
	}
	quizReader := csv.NewReader(quizFile)

	quiz := []quizData{}

	for {
		data, err := quizReader.Read()

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		quiz = append(quiz, newQuizData(data[0], data[1]))
	}

	return quiz
}

func runGame(quiz []quizData, inputReader *bufio.Reader, gameLength time.Duration) int {
	numberCorrect := 0

	fmt.Println("Press 'Enter' to start quiz...")
	fmt.Scanln()

	gameTimer := time.NewTimer(gameLength * time.Second)
	correctCounter := make(chan int)
	done := make(chan struct{})
	errC := make(chan error, 0)

	go func(done chan struct{}, errC chan error) {
		for _, data := range(quiz) {
			fmt.Println(data.question)
			answer, err := inputReader.ReadString('\n')

			if err != nil {
				errC <- err
				return
			}

			if strings.TrimRight(answer, "\n") == data.answer {
				correctCounter <- 1
			}
		}

		done <- struct{}{}
	}(done, errC)

	G:
		for {
			select {
			case <-gameTimer.C:
				fmt.Println("Times up!")
				break G
			case <-correctCounter:
				numberCorrect += 1
			case <-done:
				break G
			case err := <-errC:
				log.Fatal(err)
				os.Exit(1)
			}
		}

	return numberCorrect
}