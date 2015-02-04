package main

import (
	"bytes"
	"encoding/csv"
	"encoding/gob"
	"fmt"
	"github.com/google/btree"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

type Datum struct {
	Year       string
	Location   string
	Industry   string
	Occupation string
	Wage       float64
}

type lessFunc func(p1, p2 *Datum) bool

func (f Datum) String() string {
	return fmt.Sprintf("[Year: %v Location:%v Industry:%v Occupation:%v Wage:%v]", f.Year, f.Location, f.Industry, f.Occupation, f.Wage)
}

func (f Datum) Index() string {
	return fmt.Sprintf("%v%v%v%v", f.Year, f.Location, f.Industry, f.Occupation)
}

func LoadData() []Datum {
	csvpath := "/Users/jspeiser/code/dataviva-scripts/data/rais/Rais_2002.csv"
	csvfile, err := os.Open(csvpath)

	if err != nil {
		panic(err)
	}

	var rows = []Datum{}
	reader := csv.NewReader(csvfile)
	reader.Comma = 59 // set delim to 59

	// read header
	data, err := reader.Read()
	fmt.Println("Lookup result:", data)

	// read data
	for ; err == nil; data, err = reader.Read() {
		year := data[len(data)-1]

		if year == "Year" {
			continue
		}

		cbo := data[0]
		cnae := data[1]
		bra := data[6]
		wage, _ := strconv.ParseFloat(strings.Replace(data[10], ",", ".", 1), 32)
		d := Datum{year, cbo, cnae, bra, wage}
		rows = append(rows, d)
	}
	// err = reader.ReadStructAll(&rows)
	return rows
}

func SaveData(filepath string, data *[]Datum) {
	b := new(bytes.Buffer)
	e := gob.NewEncoder(b)
	err := e.Encode(*data)
	if err != nil {
		panic(err)
	}

	fmt.Println("Saving to Disk...")
	err = ioutil.WriteFile(filepath, b.Bytes(), 0644)
	fmt.Println("Save complete!")
}

func BuildDb(stuff *[]Datum) map[string][]*Datum {
	var raw = *stuff
	var m map[string][]*Datum
	m = make(map[string][]*Datum)

	for i := 0; i < len(raw); i++ {
		a := raw[i]
		elem, ok := m[a.Index()]
		if !ok {
			m[a.Index()] = []*Datum{&a}
		} else {
			m[a.Index()] = append(elem, &a)
		}
	}
	return m
}

func LoadFromFile(filename string) *[]Datum {
	var m2 []Datum

	b, err := os.Open(filename)
	dec := gob.NewDecoder(b)
	err = dec.Decode(&m2)

	b.Close()

	if err != nil {
		panic("Error loading from file.")
	}

	return &m2
}

func (f *Datum) Add(f2 *Datum) {
	if f.Year == "" {
		f.Year = f2.Year
	}
	if f.Location == "" {
		f.Location = f2.Location
	}
	if f.Industry == "" {
		f.Industry = f2.Industry
	}
	if f.Occupation == "" {
		f.Occupation = f2.Occupation
	}

	f.Wage += f2.Wage
}

func (f *Datum) Less(f2 btree.Item) bool {
	d2 := f2.(*Datum)
	s1 := f.Year + f.Location + f.Industry + f.Occupation
	s2 := d2.Year + d2.Location + d2.Industry + d2.Occupation
	return s1 < s2
}

func Compress(m *map[string][]*Datum) []Datum {
	var dlist []Datum

	for _, value := range *m {
		var moi *Datum = &(Datum{"", "", "", "", 0})
		for d := range value {
			moi.Add(value[d])
		}
		dlist = append(dlist, *moi)
	}
	return dlist
}

func main() {
	// var m map[string][]*Datum
	// var raw = LoadData()
	// m = BuildDb(&raw)
	// var dlist = Compress(&m)
	filepath := "/Users/jspeiser/dbrepo/my.db"
	// SaveData(filepath, &dlist)

	var dlist *[]Datum
	// if _, err := os.Stat(filename); os.IsNotExist(err) {
	// } else {
	// fmt.Println("reading from file!!!!")
	dlist = LoadFromFile(filepath)
	// }

	// datas := *dlist
	// key := "200278241049213120040"

	fmt.Println("Done!")

	datas := *dlist
	// var results []*Datum

	tr := btree.New(4)

	// Deepest level index
	for i := range datas {
		tr.ReplaceOrInsert(&datas[i])
	}

	target_d := Datum{"2002", "141410", "94910", "314930", 0}

	start := time.Now()
	result := tr.Get(&target_d)
	elapsed := time.Since(start)
	fmt.Println("Sort took: ", elapsed)
	fmt.Println("Lookup result:", result)
}
