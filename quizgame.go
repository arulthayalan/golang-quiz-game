package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const defaultFilename = "problems.csv"

type Question struct {
	question string
	answer   string
	response string
	correct  bool
}

type Reader struct {
	r io.Reader
}

//For given relative path string,
//returns absolute path of problems.csv
func resourceFilePath(path string) string {
	abspath, err := filepath.Abs(path)
	if err != nil {
		log.Fatal("Error finding file path", err)
	}
	return abspath
}

func fileExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func openFile(path string) (io.Reader, error) {
	f, err := os.Open(path)
	//defer f.Close()
	return f, err
}

//return array of records
func (c Reader) readCsv() ([][]string, error) {
	var records = [][]string{}
	reader := csv.NewReader(c.r)

	records, err := reader.ReadAll()
	if err != nil {
		return records, err
	}
	return records, nil
}

func questions(records [][]string) []Question {
	if records == nil || len(records) <= 0 {
		return nil
	}
	var questions = []Question{}
	for _, record := range records {
		q := Question{
			question: record[0],
			answer:   record[1],
		}
		questions = append(questions, q)
	}
	return questions
}

func promptUser(q []Question, limit int) []Question {
	answers := []Question{}
	timer := time.NewTimer(time.Duration(limit) * time.Second)
	for i, question := range q {

		fmt.Printf("%v Question %v\n", i+1, question.question)
		fmt.Print("Enter answer: ")

		answerCh := make(chan string)

		go func() {
			stdinreader := bufio.NewReader(os.Stdin)
			
			for {
				response, _ := stdinreader.ReadString('\n')
				response = strings.TrimSuffix(response, "\n")
				fmt.Println(response)
				if response != "" {
					answerCh <- response
					break
				} else {
					fmt.Print("Enter answer: ")
				}
			}
			
		}()

		select {
		case <-timer.C:
			return answers

		case answer := <-answerCh:
			question.response = answer
			if question.response == question.answer {
				question.correct = true
			}
			answers = append(answers, question)
		}

	}
	return answers
}

func calculate(q []Question) (int, int) {
	var correct int
	for _, question := range q {
		if question.correct {
			correct++
		}
	}
	return correct, len(q) - correct
}

func main() {

	fmt.Println("Welcome to Quiz game")

	csvFileName := flag.String("csvFileName", defaultFilename, "file name under resources directory")
	timeLimit := flag.Int("timeLimit", 30, "quiz time limit")
	flag.Parse()

	relativepath := resourceFilePath("./resource/")
	abspath := filepath.Join(relativepath, *csvFileName)
	fileinfo := fileExist(abspath)
	if !fileinfo {
		fmt.Println("File not exist")
		return
	}

	file, err := openFile(abspath)
	if err != nil {
		log.Printf("unable to open file %v", err)
	}

	reader := Reader{
		file,
	}

	records, err := reader.readCsv()

	if err != nil {
		log.Printf("unable to open file %v", err)
	}
	//var questions := make([]&Question,2)
	questions := questions(records)
	questions = promptUser(questions, *timeLimit)

	correct, incorrect := calculate(questions)

	for _, question := range questions {
		fmt.Printf("Question: %v Answer: %v response: %v \n",
			question.question, question.answer, question.response)
	}

	fmt.Printf("Correct: %v Incorrect: %v\n", correct, incorrect)

}
