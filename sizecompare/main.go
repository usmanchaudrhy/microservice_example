package main

import (
	"awesomeProject/gen"
	model "awesomeProject/metadata/pkg"
	"encoding/json"
	"fmt"
)

var metadata = &model.Metadata{
	ID:          "123",
	Title:       "The test movie",
	Description: "Sequel of the legendary the movie",
	Director:    "foo bar",
}

var genmetadata = &gen.Metadata{
	Id:          "123",
	Title:       "The test movie",
	Description: "Sequel of the legendary the movie",
	Director:    "foo bar",
}

func main() {
	jsonBytes, err := serializeToJSON(metadata)
	if err != nil {
		panic(err)
	}

	fmt.Printf("JSON size:\t%dB\n", len(jsonBytes))
}

func serializeToJSON(m *model.Metadata) ([]byte, error) {
	return json.Marshal(m)
}
func serializeToXML()   {}
func serializeToProto() {}
