package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
)

func main() {

	html_file, _ := ioutil.ReadFile("go_template.html")

	check := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}
	t, err := template.New("webpage").Parse(string(html_file))
	check(err)

	data := struct {
		Name      string
		Number    string
		ImagePath string
		Ability   []string
		Score     [][]int
		Bottom    string
	}{
		Name:      "Julian",
		Number:    "A1",
		ImagePath: "howtogeek.jpg",
		Ability: []string{
			"technique",
			"strength",
			"balance",
		},
		Score: [][]int{
			[]int{8},
			[]int{7},
			[]int{5},
		},
		Bottom: "bouldern is fun",
	}

	fmt.Println(data.Name)

	f, err := os.OpenFile("card_"+data.Name+".html", os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer f.Close()

	err = t.Execute(f, data)
	if err != nil {
		log.Fatal(err)
	}

}
