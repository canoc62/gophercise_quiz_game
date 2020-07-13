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

func main() {
	var quizFile = flag.String("f", "problems.csv", "usage: filename for quiz")
	var quizTime = flag.Int("t", 30, "usage: time duration for quiz")
	flag.Parse()

	quizData, err := os.Open(*quizFile)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
		os.Exit(1)
	}
	
	quizReader := csv.NewReader(quizData)
	inputReader := bufio.NewReader(os.Stdin)
	gameTimer := time.Duration(*quizTime)

	numQuestions, numCorrect := runGame(quizReader, inputReader, gameTimer)
	fmt.Printf("Out of %d questions , you got %d correct!\n", numQuestions, numCorrect)
}

func runGame(quizReader *csv.Reader, inputReader *bufio.Reader, gameLength time.Duration) (int, int) {
	numberOfQuestions := 0
	numberCorrect := 0

	fmt.Println("Press 'Enter' to start quiz...")
	fmt.Scanln()

	gameTimer := time.NewTimer(gameLength * time.Second)

	G:
		for {
			select {
			case <-gameTimer.C:
				fmt.Println("Times up!")
				break G
			default:
				data, err := quizReader.Read()

				if err == io.EOF {
					break G
				}
				if err != nil {
					log.Fatal(err)
				}
				numberOfQuestions += 1

				fmt.Println(data[0])
				answer, err := inputReader.ReadString('\n')

				if err != nil {
					log.Fatal(err)
				}

				fmt.Printf("Answer bytes are for string from input is: % x\n", answer)
				if strings.TrimRight(answer, "\n") == data[1] {
					numberCorrect += 1
				}
			}
		}

	return numberOfQuestions, numberCorrect
}