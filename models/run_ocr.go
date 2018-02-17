package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"
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

func ToOcr(filepath string) string {
	cmd := exec.Command("curl", "-sI", filepath)
	info, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	fileSize := regexp.MustCompile(`Content-Length.*(\d+)`)
	size, err := strconv.Atoi(strings.Split(string(fileSize.Find(info)), " ")[1])
	if err != nil {
		panic(err)
	}
	if size > 1000000*10 {
		return "File too big for automatic transcription."
	}
	baseDir := "/tmp/pdf_to_ocr_out/"
	err = os.Mkdir(baseDir, 0766)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}
	cleanDirectory(baseDir, []string{".png", "ocr.txt", "png.list"})
	out := baseDir + "out.png"
	cmd = exec.Command("convert", "-adjoin", "-density", "300", filepath,
		"-quality", "100", out)
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
	fmt.Println("\tImage conversion done.")
	files, err := ioutil.ReadDir(baseDir)
	if err != nil {
		panic(err)
	}
	pngList := make([]string, 0)
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".png") {
			pngList = append(pngList, baseDir+f.Name())
		}
	}
	err = ioutil.WriteFile(baseDir+"png.list", []byte(strings.Join(pngList, "\n")), 0666)
	if err != nil {
		panic(err)
	}
	cmd = exec.Command("tesseract", baseDir+"png.list", baseDir+"ocr")
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
	bytes, err := ioutil.ReadFile(baseDir + "ocr.txt")
	if err != nil {
		panic(err)
	}
	cleanDirectory(baseDir, []string{".png", "ocr.txt", "png.list"})
	fmt.Println("\tOCR done.")
	return string(bytes)
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
	db, err := sql.Open("mysql", username+":"+password+"@/omeka")
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

		fmt.Println(o.original_filename)
		filename := "http://localhost/files/original/" + o.filename
		text := ToOcr(filename)

		updateDocument(o, text, api_key)
	}
}

func updateDocument(o omeka_file, ocr, api_key string) {
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
	element_texts := m["element_texts"].([]interface{})
	for i := range element_texts {
		element_text := element_texts[i].(map[string]interface{})
		element := element_text["element"].(map[string]interface{})
		element_set := element_text["element_set"].(map[string]interface{})
		if element["name"] == "Text" &&
			element_set["name"] == "Item Type Metadata" {

			m["element_texts"].([]interface{})[i].(map[string]interface{})["text"] = ocr
		}
	}

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
	fmt.Println("\tResponse: ", response.Status)
}
