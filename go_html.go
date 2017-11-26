package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
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

func (data *Data) create_html_card() {

	tmpl := `
	<!DOCTYPE html>
	<html>
	<link rel="stylesheet" href="style.css">

	<div class="card">
		<div class="row">
			<h1 class="title_left">{{.Name}}</h1>
			<h1 class="title_right">{{.Number}}</h1>
		</div>

		<img src="{{.ImagePath}}" alt="Image" style="width:100%">

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
	f, _ := os.OpenFile("card_"+data.Name+".html", os.O_CREATE|os.O_WRONLY, 0777)
	defer f.Close()

	err := t.Execute(f, data)

	if err != nil {
		panic(err)
	}

	fmt.Println(data.Name)
}

func main() {

	for climber := 1; climber < 3; climber++ {
		path := "list_climbers.csv"
		card_data := read_csv(path, climber)
		card_data.create_html_card()
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
