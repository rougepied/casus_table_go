package main

import (
	"fmt"
	"strings"
	"testing"
)

func Test_rstSeparator2(t *testing.T) {
	result := rstSeparator2([]int{3, 5, 6})
	if result != "=== ===== ======" {
		t.Errorf("rstSeparator2 failed: “%s”\n", result)
	}
}

func Test_measure(t *testing.T) {
	var tabTest = []struct {
		Text  string
		Sep   string
		Table [][]string
		Size  []int
	}{
		{`XXX	YYY	ZZZ
1	2	3
4	5	6
`, "\t", [][]string{[]string{"XXX", "YYY", "ZZZ"}, []string{"1", "2", "3"}, []string{"4", "5", "6"}}, []int{3, 3, 3}},
		{`XXX;YYYff;ZZZ
1;2;3
4;5;666666
`, ";", [][]string{[]string{"XXX", "YYYff", "ZZZ"}, []string{"1", "2", "3"}, []string{"4", "5", "666666"}}, []int{3, 5, 6}},
		{`bépo;bepo;bépô
1;2;3
4;5;6
`, ";", [][]string{[]string{"bépo", "bepo", "bépô"}, []string{"1", "2", "3"}, []string{"4", "5", "6"}}, []int{4, 4, 4}},
	}
	for _, tab := range tabTest {
		expected_size := tab.Size
		expected_table := tab.Table
		reader := strings.NewReader(tab.Text)
		sep := tab.Sep
		got_table, got_size, _ := measure(reader, sep)
		if !DirtyCompare(got_size, expected_size) {
			t.Errorf("Got “%v”, expected “%v”\n", got_size, tab.Size)
		}
		if !DirtyCompare(got_table, expected_table) {
			t.Errorf("Got “%v”, expected “%v”\n", got_table, tab.Table)
		}
	}
}

func Test_measure_err(t *testing.T) {
	txt := `auie;
auie;auie;`
	reader := strings.NewReader(txt)
	if _, _, err := measure(reader, ";"); err == nil {
		t.Errorf("Fail: err == nil")
	}
}

func Test_formatCol(t *testing.T) {
	var tab = []struct {
		Text   string
		Size   int
		Result string
	}{
		{"5", 5, "    5"},
	}

	for _, val := range tab {
		result := formatCol(val.Text, val.Size)
		if val.Result != result {
			t.Errorf("Got “%v”, expected “%v”\n", result, val.Result)
		}
	}
}

func Test_DirtyCompare(t *testing.T) {
	var tab = []struct {
		In1 interface{}
		In2 interface{}
	}{
		{[]int{3, 3, 3}, []int{3, 3, 3}},
		{nil, nil},
		{[]int{3, 3, 3}, "[3 3 3]"},
		{[]string{"AA", "BB"}, []string{"AA", "B" + "B"}},
		{[]string{"AA", "BB"}, "[AA BB]"},
	}

	for _, val := range tab {
		if !DirtyCompare(val.In1, val.In2) {
			t.Errorf("Got “%v”, expected “%v”\n", val.In1, val.In2)
		}
	}
}

// DirtyCompare returns true if the input parameter are evaluated as being the same.
// It allows to compare two “arrays” in a dirty way. It is “dirty” because it
// does not realy compare each element of the array but the string produced by
// fmt.Sprintf.
//
// It allows to write quick and dirty tests by comparing the output of a
// function call to a string representing the expected result (not the real
// expected output).
//
// WARNING: DirtyCompare([]string{"AA", "BB"}, "[AA BB]") returns true.
// (See Test_DirtyCompare)
func DirtyCompare(a, b interface{}) bool {
	x := fmt.Sprintf("%v", a)
	y := fmt.Sprintf("%v", b)
	return x == y
}
