package main

import (
	"crypto/md5"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type Params struct {
	input string
	tree  map[string][]string
}

func main() {
	//fmt.Println("findduplicates 0.2")

	params := Params{}
	params.tree = make(map[string][]string)
	flag.StringVar(&params.input, "input", "", "Root directory from where duplicates should be searched.")
	flag.Parse()

	err := searchDuplicates(params)
	if err != nil {
		log.Fatal(err)
	}

	for hash, list := range params.tree {
		if len(list) > 1 {
			fmt.Println(fmt.Sprintf("@rem %s", hash))
			for _, file := range list {
				fmt.Println(fmt.Sprintf("del \"%s\"", file))
			}

		}
	}
}

func searchDuplicates(params Params) error {
	err := filepath.Walk(params.input,
		func(name string, info os.FileInfo, err error) error {
			if err != nil {
				return errors.New(fmt.Sprintf("1 %s\n", err))
			}
			if !info.IsDir() {
				hash := makeFileHash(name)
				if params.tree[hash] == nil {
					var list []string
					list = append(list, name)
					params.tree[hash] = list
				} else {
					//log.Printf("DOUBLE: %s == %s", params.tree[hash], name)
					if !contains(params.tree[hash], name) {
						params.tree[hash] = append(params.tree[hash], name)
					}
				}
			}
			return nil
		})
	if err != nil {
		return errors.New(fmt.Sprintf("2 %s\n", err))
	}
	return nil
}

func makeFileHash(name string) string {
	file, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	h := md5.New()
	if _, err := io.Copy(h, file); err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}

func contains(arr []string, search string) bool {
	for _, value := range arr {
		if value == search {
			return true
		}
	}
	return false
}
