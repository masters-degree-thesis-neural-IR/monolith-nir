package files

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alexflint/go-memdump"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
)

var (
	m      = sync.Mutex{}
	buffer = make(map[string][]string)
)

type Dump struct {
	Index string
}

func fileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if err == nil {
		return !info.IsDir()
	}
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return false
}

func deepCopy(src, dist interface{}) (err error) {
	buf := bytes.Buffer{}
	if err = gob.NewEncoder(&buf).Encode(src); err != nil {
		panic(err)
	}
	return gob.NewDecoder(&buf).Decode(dist)
}

func LoadIndex() map[string][]string {

	index := make(map[string][]string)

	files, err := ioutil.ReadDir("./data/")
	if err != nil {
		log.Fatal(err)
	}

	count := 0
	for _, file := range files {
		count++
		if !file.IsDir() {
			r, err := os.Open("./data/" + file.Name())
			if err != nil {
				log.Fatalln(err.Error())
			}
			var dump *Dump
			memdump.Decode(r, &dump)

			var temp map[string][]string
			json.Unmarshal([]byte(dump.Index), &temp)

			for key, value := range temp {
				index[key] = value
			}

		}
	}

	return index

}

func writer() {
	for {

		time.Sleep(4 * time.Second)
		m.Lock()
		for key, value := range buffer {
			term := map[string][]string{key: value}

			data, err := json.Marshal(term)

			if err != nil {
				panic(err)
			}

			dump := Dump{
				Index: string(data),
			}

			fileName := fmt.Sprintf("./data/%s.memdump", key)

			w, err := os.Create(fileName)
			if err != nil {
				log.Fatalln(err.Error())
			}
			memdump.Encode(w, &dump)
			delete(buffer, key)
		}
		m.Unlock()

	}
}

func DumpIndex(chanIndex chan map[string][]string) {

	go writer()

	for index := range chanIndex {
		m.Lock()
		for key, value := range index {
			buffer[key] = value
		}
		m.Unlock()
	}
}
