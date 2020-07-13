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
	// "time"
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

	numQuestions, numCorrect := runGame(quizReader, inputReader)
	fmt.Printf("Out of %d questions , you got %d correct!\n", numQuestions, numCorrect)
}

func runGame(quizReader *csv.Reader, inputReader *bufio.Reader) (int, int) {
	numberOfQuestions := 0
	numberCorrect := 0

	timer
	for {
		data, err := quizReader.Read()
		numberOfQuestions += 1

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

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

	return numberOfQuestions, numberCorrect
}