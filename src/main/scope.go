package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func ReadCfg(filename string) (*JsonPkg, error) {
	raw_cfg, e := ioutil.ReadFile(filename)
	cfg := JsonPkg{}
	if e != nil {
		return nil, fmt.Errorf("Error reading json file")
	}
	e = json.Unmarshal(raw_cfg, &cfg)
	return &cfg, e
}

func PrepareFilelist(dirnames []string, lang Language) (FileList []string) {
	for _, dirname := range dirnames {

		filepath.Walk(dirname, func(path string, f os.FileInfo, err error) error {
			dirnames = append(dirnames, path)
			return nil
		})

	}

	for _, dirname := range dirnames {
		process_dir := true
		for _, ignore_dir := range lang.IgnoreDirs {
			if strings.Contains(dirname, ignore_dir) {
				process_dir = false
			}
		}

		if !process_dir {
			continue
		}

		files, _ := ioutil.ReadDir(dirname)
		for _, f := range files {
			for _, ext := range lang.Extensions {
				if strings.HasSuffix(f.Name(), ext) == true {
					FileList = append(FileList, dirname+"/"+f.Name())
				}
			}
		}
	}
	return
}

func FillAllSymbolsFromLine(filename, content, parsed_content string, line_num int, symbols_found *Symbols) {
	words := strings.Split(parsed_content, " ")
	var ptr_to_symbs *Symbols = symbols_found
	visited := map[string]bool{}

	for _, w := range words {
		if w == "" {
			continue
		}
		_, ok := visited[w]

		if ok != true {
			visited[w] = true
		}

		o := Occourance{FileName: filename, LineContent: content, LineNumber: line_num, OccouranceCountInThisLine: 1}

		if occourances, ok := (*ptr_to_symbs)[w]; ok == true {
			// occorances exist
			occourances = append(occourances, o)
			(*ptr_to_symbs)[w] = occourances
		} else {
			// empty
			(*ptr_to_symbs)[w] = make([]Occourance, 1)
			(*ptr_to_symbs)[w][0] = o
		}
	}

}

func ReadFileToGetSymbols(filename string, remove []string, symbols_found *Symbols) (err error) {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	line_num := 0
	for scanner.Scan() {
		line_num++
		s := scanner.Text()
		// TODO: reading line by line.

		if len(s) != 0 {
			// strip all keywords
			for _, rem_sym := range remove {
				s = strings.Replace(s, rem_sym, " ", -1)
			}
			if len(strings.Replace(s, " ", "", -1)) != 0 {
				FillAllSymbolsFromLine(filename, scanner.Text(), s, line_num, symbols_found)
			}
		}
	}
	return nil
}

func main() {
	cfg, e := ReadCfg("./package.json")

	if e != nil {
		fmt.Println("failed:", e)
	}
	var dirs []string
	if len(os.Args) > 1 {
		dirs = os.Args[1:]
	} else {
		dirs = []string{"./"}
	}

	all_symbols := Symbols{}

	for _, lang := range cfg.Languages {
		remove_words := []string{}
		remove_words = append(remove_words, lang.Keywords...)
		remove_words = append(remove_words, lang.Operators...)
		remove_words = append(remove_words, lang.Whitespace...)

		FileList := PrepareFilelist(dirs, lang)
		for _, f := range FileList {
			e := ReadFileToGetSymbols(f, remove_words, &all_symbols)
			if e != nil {
				fmt.Println(e.Error())
			}
		}

	} // Done!

	for s, occ := range all_symbols {
		fmt.Println(s, "\t count=(", len(occ), ")")
		for _, o := range occ {
			fmt.Println("\t ", o.LineNumber, o.FileName, o.LineContent)
		}
	}
}
