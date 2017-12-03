package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"text/template"
)

type Ability struct {
	Name string
	Val  int
}

type Data struct {
	Name      string
	Number    string
	ImagePath string
	Abilities []Ability
	Bottom    string
}

func (data *Data) create_html_card(path string) {

	tmpl := `
	<!DOCTYPE html>
	<html>
	<link rel="stylesheet" href="style.css">

	<div class="card">
		<div class="row">
			<h1 class="title_left">{{.Name}}</h1>
			<h1 class="title_right">{{.Number}}</h1>
		</div>
		
		<div class="image_container"><img src="{{.ImagePath}}" alt="Image"></div>

		<div class="row">
		 	{{  range $ability := .Abilities }}
	 		<div class="left">
	 			<div>{{ $ability.Name }}</div>
	 		</div>
	 		<div class="right">
	 			<div>{{ $ability.Val }}</div>
	 		</div>
			{{end}}

			<div class="bottom">{{.Bottom}}</div>

		</div>
	</div>
	</html>
	`

	t := template.Must(template.New("webpage").Parse(tmpl))
	f, _ := os.OpenFile(path+"/output_html/card_"+data.Name+".html", os.O_CREATE|os.O_WRONLY, 0777)
	defer f.Close()

	err := t.Execute(f, data)

	if err != nil {
		panic(err)
	}

	fmt.Println(data.Name+"'s card was saved in", path)
}

func main() {

	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Please provide path to csv.")
		os.Exit(1)
	}

	path_to_csv := args[0]
	number_of_cards := number_rows_csv(path_to_csv)
	slice_path := strings.Split(path_to_csv, "/")
	slice_path = []string(slice_path)
	path_to_dir := strings.Join(slice_path[0:(len(slice_path)-1)], "/")

	for card := 1; card < number_of_cards; card++ {
		card_data := read_csv(path_to_csv, card)
		card_data.create_html_card(path_to_dir)
	}

}

func read_csv(path string, row int) *Data {
	source, _ := os.Open(path)

	// Create a new reader.
	var reader = csv.NewReader(bufio.NewReader(source))
	reader.Comma = ','

	d := Data{}

	var index_abilities []int
	counter_abilities := 0
	array_abilities := []string{}

	header, _ := reader.Read()
	for i := range header {
		res, _ := regexp.MatchString("Ability", header[i])
		if res {
			index_abilities = append(index_abilities, i)
			re := regexp.MustCompile("Ability ")
			array_abilities = append(array_abilities, re.ReplaceAllString(header[i], ""))
			counter_abilities += 1
		}
	}

	counter := 1
	for {
		record, err := reader.Read()
		// Stop at EOF.
		if (err == io.EOF) || (counter > row) {
			break
		}

		d.Name = record[0]
		d.Number = record[1]
		d.ImagePath = record[2]

		var scores []int
		for _, element := range index_abilities {

			s, _ := strconv.Atoi(record[element])
			scores = append(scores, s)
		}
		var abilities []Ability
		for index, value := range array_abilities {
			abilities = append(abilities, Ability{value, scores[index]})
		}
		d.Abilities = abilities
		d.Bottom = record[len(record)-1]

		counter += 1
	}
	return (&d)
}

func number_rows_csv(path string) int {
	source, _ := os.Open(path)

	// Create a new reader.
	var reader = csv.NewReader(bufio.NewReader(source))
	reader.Comma = ','

	var counter int
	counter = 0

	for {
		_, err := reader.Read()
		// Stop at EOF.
		if err == io.EOF {
			break
		}

		counter += 1

	}
	return counter
}
