package main

import (
	"bufio"
	"fmt"
	"github.com/xuri/excelize"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func main() {

	// read all words english
	mAll := make(map[string]int)
	readFileAll, err := os.Open("words.txt")
	if err != nil {
		fmt.Println(err)
	}
	fileAllScanner := bufio.NewScanner(readFileAll)
	fileAllScanner.Split(bufio.ScanLines)
	for fileAllScanner.Scan() {
		text := strings.ToLower(fileAllScanner.Text())
		index := strings.Index(text, " ")
		if index == -1 {
			mAll[text] = 1
		} else {
			mAll[text[0:index]] = 1
		}
	}
	readFileAll.Close()

	// read file ignore text
	mIgnore := make(map[string]int)
	readFile, err := os.Open("ignore_words.txt")
	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		//fmt.Println(fileScanner.Text())
		text := strings.ToLower(fileScanner.Text())
		//s := [0]string{text}
		index := strings.Index(text, " ")
		if index == -1 {
			mIgnore[text] = 1
		} else {
			mIgnore[text[0:index]] = 1
		}
	}
	readFile.Close()

	// read file subtitle
	m := make(map[string]int)
	//argFileName := os.Args[1]
	argFileName := "tt01.IGN"
	argFileNameArray := strings.Split(argFileName, ",")
	for _, fileName := range argFileNameArray {
		f, err := os.Open(fileName)
		//f, err := os.Open(argFileName)
		if err != nil {
			fmt.Println(err)
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		scanner.Split(bufio.ScanWords)

		for scanner.Scan() {
			//fmt.Println(scanner.Text())
			text := strings.ToLower(scanner.Text())
			text = strings.ReplaceAll(text, "?", "")
			text = strings.ReplaceAll(text, "!", "")
			text = strings.ReplaceAll(text, ",", "")
			text = strings.ReplaceAll(text, ".", "")
			text = strings.ReplaceAll(text, "-", "")
			text = strings.ReplaceAll(text, "\"", "")
			text = strings.ReplaceAll(text, ":", "")
			text = strings.ReplaceAll(text, "”", "")
			text = strings.ReplaceAll(text, "~", "")
			text = strings.ReplaceAll(text, "—", "")
			//text = regexp.MustCompile(`[^a-zA-Z]+`)
			//.ReplaceAllString(text, "")
			numberRegexp := regexp.MustCompile(`\d`)
			matchNumber := numberRegexp.MatchString(text)
			if matchNumber {
				continue
			}
			if strings.TrimSpace(text) == "" || strings.Contains(text, "_") || strings.Contains(text, "$") {
				continue
			}
			if _, ok := mIgnore[text]; ok {
				continue
			}
			if _, ok := mAll[text]; !ok {
				continue
			}
			count := 1
			if val, ok := m[text]; ok {
				count += val
			}
			m[text] = count
		}

		if err := scanner.Err(); err != nil {
			fmt.Println(err)
		}
	}

	// sort text
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.SliceStable(keys, func(i, j int) bool {
		return m[keys[i]] > m[keys[j]]
	})

	//for _, k := range keys {
	//	fmt.Println(k, m[k])
	//}
	newFile := excelize.NewFile()
	i := 1
	for _, k := range keys {
		newFile.SetCellValue("Sheet1", "A"+strconv.Itoa(i), k)
		newFile.SetCellValue("Sheet1", "B"+strconv.Itoa(i), m[k])
		//newFile.SetCellValue("Sheet1", "C"+strconv.Itoa(i), translate(k))
		i++
	}

	if err := newFile.SaveAs("simple.xlsx"); err != nil {
		log.Fatal(err)
	}
}

//func translate(text string) string {
//	url := "https://google-translate1.p.rapidapi.com/language/translate/v2"
//
//	payload := strings.NewReader("q=" + text + "&target=vi&source=en")
//
//	req, _ := http.NewRequest("POST", url, payload)
//
//	req.Header.Add("content-type", "application/x-www-form-urlencoded")
//	req.Header.Add("Accept-Encoding", "application/gzip")
//	req.Header.Add("X-RapidAPI-Key", "95d772dfb5msh80853654d278817p14d19ejsnc15d9601610a")
//	req.Header.Add("X-RapidAPI-Host", "google-translate1.p.rapidapi.com")
//
//	res, _ := http.DefaultClient.Do(req)
//
//	defer res.Body.Close()
//	body, _ := ioutil.ReadAll(res.Body)
//
//	//fmt.Println(res)
//	fmt.Println(string(body))
//	return body.data.translations
//}

type BodyGoogleTranslateAPI struct {
	translations string
}
