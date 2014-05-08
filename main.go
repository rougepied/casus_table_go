package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	// Parsing command line args
	path := flag.String("file", "", "path to file")
	sep := flag.String("sep", ";", "column separator")
	flag.Parse()

	// Input file
	if *path != "" {

		file, err := os.Open("input.csv")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		table, size := measure(file, sep)
		table2 := formatTable(table, size)
		fmt.Println(rst2Table(table2, size))
	}
}

func rst2Table(table [][]string, size []int) (r string) {
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
func measure(input io.Reader, sep *string) ([][]string, []int) {
	var r []int
	var t [][]string

	scanr := bufio.NewScanner(input)
	for scanr.Scan() {
		line := scanr.Text()
		tab := strings.Split(line, *sep)
		t = append(t, tab)
		if len(r) == 0 {
			r = make([]int, len(tab))
		}

		rline := make([]int, len(tab))
		for idx, val := range tab {
			rline[idx] = len(val)
			r[idx] = max(r[idx], len(val))
		}
	}

	if err := scanr.Err(); err != nil {
		log.Fatal(err)
	}
	return t, r
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
