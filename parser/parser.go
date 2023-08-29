package parser

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type TOMLData struct {
	Grades []GradeData
}

type GradeData struct {
	Grade  string
	Topics []TopicData
}

type TopicData struct {
	Title     string
	Questions []QuestionData
}

type QuestionData struct {
	Question string
	Answer   string
	Feedback string
}

func DecodeTOML() *TOMLData {
	data := &TOMLData{}
	if _, err := toml.DecodeFile("grades.toml", &data); err != nil {
		log.Fatal(err)
	}

	return data
}

func EncodeTOML(data GradeData) {
	f, err := os.Create("./result.toml")
	if err != nil {
		log.Fatal(err)
	}
	if err = toml.NewEncoder(f).Encode(data); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
