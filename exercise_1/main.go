package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

func readCSV(filepath string) *csv.Reader {

	file, err := os.ReadFile(filepath) // For read access.
	if err != nil {
		log.Fatal(err)
	}
	return csv.NewReader(strings.NewReader(string(file)))
}
func checkForWin(err error) bool {
	if err == io.EOF {
		return true
	}
	if err != nil {
		log.Fatal(err)
	}
	return false
}

func getInput(question string, index int, answerCh chan string) {
	var input string
	fmt.Printf("Problem #%d: %s = ", index, question)
	fmt.Scanln(&input)
	answerCh <- input
}

func checkForAnswer(input, answer string) bool {
	if input == answer {
		return true
	}
	return false

}

func shuffleArr(arr [][]string) [][]string {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(arr), func(i, j int) { arr[i], arr[j] = arr[j], arr[i] })
	return arr
}

func main() {

	inputFileFlag := flag.String("i", "./problems.csv", "path to questions file")
	shuffleFlag := flag.Bool("s", false, "shuffles the questions for an extra challenge")
	timeFlag := flag.Int("l", 30, "the limit for answering the questions in seconds")
	flag.Parse()

	records, err := readCSV(*inputFileFlag).ReadAll()
	if err != nil {
		fmt.Println("something went wrong")
	}

	var r [][]string
	if *shuffleFlag {
		fmt.Println("shuffle is on")
		r = shuffleArr(records)
	} else {
		r = records
	}

	timer := time.NewTimer(time.Duration(*timeFlag) * time.Second)
	correct := 0
	for index, record := range r {
		answerCh := make(chan string)
		go func() {
			getInput(record[0], index+1, answerCh)
		}()
		select {
		case <-timer.C:
			fmt.Printf("\ntime is over")
			fmt.Printf("\nyou solved %d out of %d \n", index, len(r))
			os.Exit(1)
		case answer := <-answerCh:
			if checkForAnswer(record[1], answer) {
				correct++
			}
		}
	}
	fmt.Printf("you solved %d out of %d \n", correct, len(r))
	fmt.Println("congrats, you win")
}
