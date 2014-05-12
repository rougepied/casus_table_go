package main

import (
	"bufio"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

//var templates = template.Must(template.ParseFiles("casus_table.html"))
var html = `<html>
<style type="text/css">
.err { color: red; }
</style>
<body>
	<h1>Casus Table</h1>
	{{if .Err}}
	<p class="err">{{.Err}}</p>
	{{end}}
	<form action="." method="POST">
		<ul>
		<li>Separator:
		<input type="radio" name="sep" value="semico" {{if eq .Sep "semico"}}checked{{end}} ><label for="sep">semicolon</label>
		<input type="radio" name="sep" value="tabula" {{if eq .Sep "tabula"}}checked{{end}} ><label for="sep">tab</label>
		</li>
		<li>Style:
		<input type="radio" name="style" value="grid" {{if eq .Style "grid"}}checked{{end}} ><label for="sep">Grid table</label>
		<input type="radio" name="style" value="simp" {{if eq .Style "simp"}}checked{{end}} ><label for="sep">Simple table</label>
		</li>
		</ul>
		<textarea name="input" rows="20" cols="80">{{printf "%s" .Input}}</textarea>
		<input type="submit" value="Convert">
		<textarea name="output" rows="20" cols="80">{{printf "%s" .Output}}</textarea>
	</form>
</body>
</html>`

func main() {
	// Parsing command line args
	path := flag.String("file", "", "path to file")
	sep := flag.String("sep", ";", "column separator")
	port := flag.String("port", ":8080", "port (if running as a web app)")
	flag.Parse()

	if *path != "" {
		// Input file
		file, err := os.Open("input.csv")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		table, size, _ := measure(file, *sep)
		table2 := formatTable(table, size)
		fmt.Println(rst2TableSimple(table2, size))
	} else {
		// run as a server
		log.Println("Running as server")
		http.HandleFunc("/", handler)
		http.ListenAndServe(*port, nil)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	p := make(map[string]interface{})

	// processing
	input := r.FormValue("input")
	var sep string
	switch r.FormValue("sep") {
	default:
		sep = ";"
		p["Sep"] = "semico"
	case "semico":
		sep = ";"
		p["Sep"] = "semico"
	case "tabula":
		sep = "\t"
		p["Sep"] = "tabula"
	}
	var style string
	switch r.FormValue("style") {
	default:
		style = "simp"
		p["Style"] = "simp"
	case "simp":
		style = "simp"
		p["Style"] = "simp"
	case "grid":
		style = "grid"
		p["Style"] = "grid"
	}
	table, size, err := measure(strings.NewReader(input), sep)
	table2 := formatTable(table, size)
	var output string
	switch style {
	case "simp":
		output = rst2TableSimple(table2, size)
	case "grid":
		output = rst2TableGrid(table2, size)
	}
	// responsing
	p["Input"] = input
	p["Output"] = output

	if err != nil {
		p["Err"] = err
	}

	t, err := template.New("casus_table").Parse(html)
	err = t.Execute(w, p)
	// err := templates.Execute(w, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func rst2TableGrid(table [][]string, size []int) (r string) {
	separator := rstSeparator(size)
	for idx, line := range table {
		if idx == 0 {
			r += separator + "\n"
		}
		r += "|" + strings.Join(line, "|") + "|\n"
		r += separator + "\n"
	}
	return
}

func rst2TableSimple(table [][]string, size []int) (r string) {
	nbline := len(table)
	separator := rstSeparator2(size)
	for idx, line := range table {
		if idx == 0 || idx == 1 {
			r += separator + "\n"
		}
		r += strings.Join(line, " ") + "\n"
		if idx == nbline-1 {
			r += separator + "\n"
		}
	}
	return
}

func formatTable(table [][]string, size []int) (r [][]string) {
	for _, line := range table {
		var l []string
		for idx, cell := range line {
			l = append(l, formatCol(cell, size[idx]))
		}
		r = append(r, l)
	}
	return
}

func rstSeparator(cols []int) string {
	t := make([]string, len(cols))
	for idx, val := range cols {
		for j := 0; j < val; j++ {
			t[idx] += "-"
		}
	}

	return "+" + strings.Join(t, "+") + "+"
}

func rstSeparator2(cols []int) string {
	t := make([]string, len(cols))
	for idx, val := range cols {
		for j := 0; j < val; j++ {
			t[idx] += "="
		}
	}

	return strings.Join(t, " ")
}

// measure measures the width of columns.
func measure(input io.Reader, sep string) ([][]string, []int, error) {
	var r []int
	var t [][]string

	scanr := bufio.NewScanner(input)
	for scanr.Scan() {
		line := scanr.Text()
		tab := strings.Split(line, sep)
		t = append(t, tab)
		if len(r) == 0 {
			r = make([]int, len(tab))
		}

		rline := make([]int, len(tab))

		if len(rline) > len(r) {
			return nil, nil, fmt.Errorf("Error: the line “%s” is too long", line)
		}

		for idx, val := range tab {
			rline[idx] = len(val)
			r[idx] = max(r[idx], len(val))
		}
	}

	if err := scanr.Err(); err != nil {
		log.Fatal(err)
	}
	return t, r, nil
}

// max return the maximum value between 2 int
func max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

// formatCol formats input string to add heading space characters.
func formatCol(text string, width int) string {
	return fmt.Sprintf("%"+strconv.Itoa(width)+"s", text)
}
