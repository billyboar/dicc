package main

import (
	"fmt"
	"time"

	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const BASE_URL = "http://dictionary.reference.com/browse/%s?s=t"
const TAG_NAME = "def-content"
const TAG_END = "</div>"
const OPEN_TAG = 60
const CLOSE_TAG = 62

//knutt morris pratt algorithm here
func generatePrefix(pat string) (result []int) {
	length := len(pat)
	result = make([]int, length)
	l := 0
	for i := 1; i < length; i++ {
		for l > 0 && pat[l] != pat[i] {
			l = result[l]
		}
		if pat[l] == pat[i] {
			l++
		}
		result[i] = l
	}
	return result
}

func kmpSearch(text []byte, pat string) (results []int) {
	prefixArray := generatePrefix(pat)
	pat_length := len(pat)
	text_length := len(text)
	results = []int{}

	l := 0

	for i := 0; i < text_length; {
		if pat[l] == text[i] {
			l++
			i++
		}

		if l == pat_length {
			results = append(results, i-l+1)
			l = prefixArray[l-1]
		} else if pat[l] != text[i] {
			if l != 0 {
				l = prefixArray[l-1]
			} else {
				i++
			}
		}
	}

	return results
}

func MergeArray(x, y []byte) (merged []byte) {
	merged = []byte{}
	for _, i := range x {
		merged = append(merged, i)
	}

	for _, l := range y {
		merged = append(merged, l)
	}

	return merged
}

func TagRemove(text []byte) []byte {
	var open_tag_pos, close_tag_pos int
	open_tag_pos = SearchFromArray(text, OPEN_TAG, 0)
	for open_tag_pos != -1 {
		close_tag_pos = SearchFromArray(text, CLOSE_TAG, open_tag_pos)
		text = MergeArray(text[:open_tag_pos], text[close_tag_pos+1:])
		open_tag_pos = SearchFromArray(text, OPEN_TAG, 0)
	}
	return text
}

func SearchFromArray(array []byte, element byte, start_pos int) int {
	length := len(array)
	for i := start_pos; i < length; i++ {
		if array[i] == element {
			return i
		}
	}
	return -1
}

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("Please enter the word you want to search")
		return
	}

	word_to_search := os.Args[1]
	search_url := fmt.Sprintf(BASE_URL, word_to_search)
	res, err := http.Get(search_url)

	if err != nil {
		fmt.Println("There is an error: ", err)
		return
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		fmt.Println("There is an error: ", err)
		return
	}

	fmt.Println(time.Now())

	word_res := kmpSearch(body, TAG_NAME)
	if len(word_res) == 0 {
		fmt.Printf("Can't find information about: %s\n", word_to_search)
	}
	for i, x := range word_res {
		tag_end := strings.Index(string(body[x:]), ">")
		div_end := strings.Index(string(body[x:]), TAG_END)
		word_def := TagRemove(body[(x + tag_end + 1):(x + div_end)])
		fmt.Printf("%d) %s\n", i+1, word_def)
	}
}
