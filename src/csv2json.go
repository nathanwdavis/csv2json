package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"flag"
	//"fmt"
	"io"
	"os"
	"strconv"

	"csv2json/typeguessing"
)

func main() {

	useLines := flag.Bool("use-lines", false,
		"Json object per line, rather than in an array.")
	numLinesForTypeModel := flag.Int("infer-lines", 20,
		"Number of lines to parse to infer the type model")
	//useStringsOnly := flag.Bool("s")
	flag.Parse()

	remainingArgs := flag.Args()
	if len(remainingArgs) < 1 {
		panic("No input file given")
	}
	inFileName := remainingArgs[0]

	csvf, fIn := getCsvReader(inFileName)
	defer fIn.Close()

	line, err := csvf.Read()
	if err != nil {
		panic(err)
	} else if len(line) < 1 {
		panic("No fields found in header (first) line")
	}
	fields := line

	typeMap, err := inferTypes(csvf, fields, *numLinesForTypeModel)
	if err != nil {
		panic(err)
	}

	// head back to the beginning of the file, then skip the header
	fIn.Seek(0, 0)
	csvf = csv.NewReader(fIn)
	csvf.Read()

	// get a Json Encoder for a file
	jsonw, fOut := getJsonWriter(inFileName + ".json")
	defer fOut.Close()

	if !*useLines {
		fOut.WriteString("[\n")
	}
	for i := 1; i < 9223372036854775806; i++ {
		line, err = csvf.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		data := buildRecord(line, fields, typeMap)
		writeRecord(jsonw, data)
		if !*useLines {
			fOut.WriteString(",")
		}
	}
	if !*useLines {
		fOut.WriteString("]")
	}

}

func getCsvReader(fName string) (*csv.Reader, *os.File) {
	fr, err := os.Open(fName)
	if err != nil {
		panic("Could not open input CSV file: " + err.Error())
	}
	return csv.NewReader(fr), fr
}

func getLine(csv *csv.Reader) ([]string, error) {
	return csv.Read()
}

func getJsonWriter(fName string) (*json.Encoder, *os.File) {
	fw, err := os.Create(fName)
	if err != nil {
		panic("Could not open output JSON file: " + err.Error())
	}
	return json.NewEncoder(bufio.NewWriter(fw)), fw
}

func buildRecord(line, fields []string,
	template map[string]interface{}) map[string]interface{} {

	data := make(map[string]interface{}, len(template))
	for idx, field := range fields {
		var typedv interface{}
		var err error
		switch typ := template[field]; {
		case line[idx] == "" && typ != typeguessing.STRING:
			typedv = nil
		case typ == typeguessing.INT:
			typedv, err = strconv.ParseInt(line[idx], 0, 64)
			if err != nil {
				panic("Field in line did not match inferred type (int): " + err.Error())
			}
		case typ == typeguessing.BOOL:
			typedv, err = strconv.ParseBool(line[idx])
			if err != nil {
				panic("Field in line did not match inferred type (bool): " + err.Error())
			}
		case typ == typeguessing.FLOAT:
			typedv, err = strconv.ParseFloat(line[idx], 64)
			if err != nil {
				panic("Field in line did not match inferred type (float): " + err.Error())
			}
		default:
			typedv = line[idx]
		}
		data[field] = typedv
	}
	return data
}

func writeRecord(jsonw *json.Encoder, data map[string]interface{}) error {
	return jsonw.Encode(data)
}

func inferTypes(csv *csv.Reader, fields []string,
	numLines int) (map[string]interface{}, error) {

	template := make(map[string]interface{})
	learners := make([]*typeguessing.Learner, len(fields))
	for i := 0; i < len(learners); i++ {
		learners[i] = typeguessing.NewLearner()
	}
	for i := 0; i < numLines; i++ {
		line, err := csv.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		for j, _ := range fields {
			learners[j].Feed(line[j])
		}
	}
	for i, f := range fields {
		exampleVal := learners[i].BestGuess()
		template[f] = exampleVal
	}
	return template, nil
}
