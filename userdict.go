package gorimer

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

type KeyValue struct {
	Key   string
	Value string
}

type ByValue []KeyValue

// Len Less and Swap implement interface sort.Interface to leverage sort functionality.
// Len is the number of elements in the collection.
func (bv ByValue) Len() int {
	return len(bv)
}

// Less reports whether the element with index i
// must sort before the element with index j.
func (bv ByValue) Less(i, j int) bool {
	return bv[i].Value < bv[j].Value
}

// Swap swaps the elements with indexes i and j.
func (bv ByValue) Swap(i, j int) {
	bv[i], bv[j] = bv[j], bv[i]
}

// check checks err, if not nil, the program will exit.
func check(err error) {
	if err != nil {
		panic(err)
	}
}

// Alphabet sorts the items in user dict of rime
func Alphabet(f string) {
	// open file
	file, err := os.OpenFile(f, os.O_RDWR, 0666)
	check(err)
	// close the file after leaving the func
	defer file.Close()
	// read lines of the file
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	//find the target lines and put them in ByValue collection as KeyValue struct
	var items ByValue
	index := 0 //find the start of the target items
	isStarted := false
	item := regexp.MustCompile(`^[\p{Han}]+\t[a-z]+$`)
	for _, line := range lines {
		found := item.FindString(line)
		if found != "" {
			if !isStarted {
				isStarted = true
			}
			words := strings.Split(found, "\t")
			items = append(items, KeyValue{words[0], words[1]})
		}
		if !isStarted {
			index++
		}
	}
	err = scanner.Err()
	check(err)
	// sort them
	sort.Sort(items)
	// replace items with sorted ones
	lines = lines[:index]
	for _, v := range items {
		lines = append(lines, fmt.Sprintf("%s\t%s", v.Key, v.Value))
	}
	// write them back to the file
	_, err = file.Seek(0, 0)
	check(err)
	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err = writer.WriteString(line + "\n")
		check(err)
	}
	err = writer.Flush()
	check(err)
}
