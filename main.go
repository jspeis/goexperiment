package main

import (
	"bytes"
	// "encoding/csv"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

type Datum struct {
	Year       int
	Location   string
	Industry   string
	Occupation string
}

func (f Datum) String() string {
	return fmt.Sprintf("%v - %v - %v - %v", f.Year, f.Location, f.Industry, f.Occupation)
}

func (f Datum) Index() string {
	return fmt.Sprintf("%v%v%v", f.Year, f.Location, f.Industry)
}

func LoadData() []Datum {
	// check if file exists. If it does, load from file,
	// else build and write file
	var tmp []Datum
	for i := 0; i < 1000000; i++ {
		a := Datum{i % 50000, "4mg", "C", "D"}
		tmp = append(tmp, a)
	}
	return tmp
}

func SaveData(m *map[string][]*Datum) {
	b := new(bytes.Buffer)

	e := gob.NewEncoder(b)

	// Encoding the map
	err := e.Encode(*m)
	if err != nil {
		panic(err)
	}

	fmt.Println("Saving to Disk...")
	// fmt.Println("%v", b.Bytes())
	err = ioutil.WriteFile("/tmp/my.db", b.Bytes(), 0644)
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

func LoadFromFile(filename string) map[string][]*Datum {
	var m2 map[string][]*Datum

	b, err := os.Open(filename)
	dec := gob.NewDecoder(b)
	err = dec.Decode(&m2)

	b.Close()

	if err != nil {
		panic("Error loading from file.")
	}

	return m2
}

func main() {
	filename := "/tmp/my.db"
	var m map[string][]*Datum
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Println("no such file... creating new...")
		var raw = LoadData()
		m = BuildDb(&raw)
		SaveData(&m)
	} else {
		fmt.Println("reading from file!!!!")
		m = LoadFromFile(filename)
	}

	key := "14mgC"
	start := time.Now()
	result := m[key]
	elapsed := time.Since(start)

	fmt.Printf("HELLO. Done!\n")
	fmt.Printf("Lookup took %s\n", elapsed)
	fmt.Println("get100:    ", result)
}
