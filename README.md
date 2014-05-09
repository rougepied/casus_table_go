# casus_table_go

A Tool to ease the creation of pure plain text table.

There are two different modes:

- command-line
- web app


## License

casus_table_go is distributed under the terms of the [WTFPL](http://www.wtfpl.net/).
See the LICENSE file for more details. 

## Installation


Make sure you have a working Go environment. [See the install instructions](http://golang.org/doc/install.html).

To install casus_table_go, simply run:
```
go get github.com/rougepied/casus_table_go
go build ...casus_table_go
```

## Usage exemple

### Command line

Imagine you have a `input.txt` file containing 

```
XXX;YYY;ZZZ
1;2;3
4;5;6_666
```

Just run: `casus_table_go -file input.txt` and it will display 
```
=== === =====
XXX YYY   ZZZ
=== === =====
  1   2     3
  4   5 6_666
=== === =====
```

### Web app

Just call `casus_table_go` and visit `http://localhost:8080`.

Note: you can specify another port. Ex: `casus_table_go -port=":6666"`

