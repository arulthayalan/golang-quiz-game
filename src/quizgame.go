package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"bufio"
	"strings"
)

const filename = "problems.csv"

type Question struct {
	question string
	answer   string
	response string
	correct bool
}

type Reader struct {
	r io.Reader
}

func ResourceFilePath() string {
	abspath, err := filepath.Abs("../resource/")
	if err != nil {
		log.Fatal("Error finding file path", err)
	}
	return filepath.Join(abspath, filename)
}

func fileexist(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func openfile(path string) (io.Reader, error) {
	f, err := os.Open(path)
	//defer f.Close()
	return f, err
}

func (c Reader) ReadCsv() ([][]string, error) {
	var records = [][]string{}
	reader := csv.NewReader(c.r)

	records, err := reader.ReadAll()
	if err != nil {
		return records, err
	}
	return records, nil
}

func converttoquestion(records [][]string) []Question {
	if records == nil || len(records) <=0 {
		return nil
	}
	var questions = []Question{}
	for _, record := range records {
		q := Question{record[0],record[1], "", false,}
		questions = append(questions, q)
	}
	return questions
}

func askuser(q []Question) []Question {
	answers := []Question{}
	for i, question := range q {
		stdinreader := bufio.NewReader(os.Stdin)
		fmt.Printf("%v Question %v\n", i+1, question.question)
		fmt.Print("Enter answer: ")
		var userresponse = ""
		for {
			userresponse, _ = stdinreader.ReadString('\n')
			userresponse = strings.TrimSuffix(userresponse, "\n")
			fmt.Println(userresponse)
			if userresponse != "" {
				break
			} else {
				fmt.Print("Enter answer: ")
			}
		}
		question.response = userresponse
		if (question.response == question.answer) {
			question.correct = true
		}
		answers = append(answers, question)
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
	filepath := ResourceFilePath()
	fileinfo := fileexist(filepath)
	if !fileinfo {
		fmt.Println("File not exist")
		return
	}

	file, err := openfile(filepath)
	if err != nil {
		log.Printf("unable to open file %v", err)
	}

	reader := Reader{
		file,
	}

	records, err := reader.ReadCsv()

	if err != nil {
		log.Printf("unable to open file %v", err)
	}
    //var questions := make([]&Question,2)
	questions := converttoquestion(records)
	questions = askuser(questions)
	
	correct, incorrect := calculate(questions)
	
	for _, question := range questions {
		fmt.Printf("Question: %v Answer: %v response: %v \n", 
			question.question, question.answer, question.response)
	}

	fmt.Printf("Correct: %v Incorrect: %v\n", correct, incorrect)

}
