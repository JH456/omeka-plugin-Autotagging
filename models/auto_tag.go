package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type omeka_search_text struct {
	id          int64
	record_type []byte
	record_id   []byte
	public      bool
	title       string
	text        string
}

type omeka_file struct {
	id                   int
	item_id              int
	order                sql.NullInt64
	size                 int
	has_derivative_image bool
	authentication       []byte
	mime_type            []byte
	type_os              []byte
	filename             string
	original_filename    string
	modified             []byte
	added                []byte
	stored               bool
	metadata             string
}

func cleanDirectory(baseDir string, suffixes []string) {
	files, err := ioutil.ReadDir(baseDir)
	for _, f := range files {
		for _, s := range suffixes {
			if strings.HasSuffix(f.Name(), s) {
				err = os.Remove(baseDir + f.Name())
				if err != nil {
					panic(err)
				}
			}
		}
	}
}

func main() {
	args := os.Args[1:]
	username := args[0]
	password := args[1]
	api_key := args[2]
	database := args[3]
	db, err := sql.Open("mysql", username+":"+password+"@/"+database)
	defer db.Close()
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Connection established.")
	}

	results, err := db.Query("SELECT * FROM omeka_files")
	if err != nil {
		panic(err)
	}
	defer results.Close()

	for results.Next() {
		var o omeka_file
		if err := results.Scan(
			&o.id,
			&o.item_id,
			&o.order,
			&o.size,
			&o.has_derivative_image,
			&o.authentication,
			&o.mime_type,
			&o.type_os,
			&o.filename,
			&o.original_filename,
			&o.modified,
			&o.added,
			&o.stored,
			&o.metadata,
		); err != nil {
			panic(err)
		}

		tagAndUpdateDocument(o, api_key)
	}
}

func tagAndUpdateDocument(o omeka_file, api_key string) {
	id := strconv.Itoa(o.item_id)
	item_url := "http://localhost/api/items/" + id
	response, err := http.Get(item_url)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	buf, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	var result interface{}
	json.Unmarshal(buf, &result)
	m := result.(map[string]interface{})

	var text string
	element_texts := m["element_texts"].([]interface{})
	for i := range element_texts {
		element_text := element_texts[i].(map[string]interface{})
		element := element_text["element"].(map[string]interface{})
		element_set := element_text["element_set"].(map[string]interface{})
		if element["name"] == "Text" &&
			element_set["name"] == "Item Type Metadata" {
			text = element_text["text"].(string)
		}
	}

	extractedTags := extractTags(text)

	tags := m["tags"].([]interface{})
	for _, t := range extractedTags {
		new := map[string]interface{}{
			"resource": "tags",
			"name":     t,
		}

		tags = append(tags, new)
	}
	m["tags"] = tags

	json_bytes, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	updated := strings.NewReader(string(json_bytes))
	request, err := http.NewRequest("PUT", item_url+"?key="+api_key, updated)
	request.Header.Set("Content-Type", "application/json")

	response, err = http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	buf, err = ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(response)
	fmt.Println("\tResponse: ", response.Status)
}

func extractTags(text string) []string {
	return append(GetKeywords(text, 1.5), GetNamedEntities(text)...)
}
