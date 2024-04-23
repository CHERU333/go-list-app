package main

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type MusicList struct {
	Music []string
}

func New(Music []string) *MusicList {
	return &MusicList{Music: Music}
}

func fileRead(fileName string) []string {
	var MusicList []string
	file, err := os.Open(fileName)
	if os.IsNotExist(err) {
		return nil
	}
	defer file.Close()
	scaner := bufio.NewScanner(file)
	for scaner.Scan() {
		MusicList = append(MusicList, scaner.Text())
	}
	return MusicList
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	MusicList := fileRead("musiclist.txt")
	html, err := template.ParseFiles("view.html")
	if err != nil {
		log.Fatal(err)
	}
	getMusic := New(MusicList)
	if err := html.Execute(w, getMusic); err != nil {
		log.Fatal(err)
	}
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	formValue := r.FormValue("value")
	file, err := os.OpenFile("musiclist.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(0600))
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	_, err = fmt.Fprintln(file, formValue)
	if err != nil {
		log.Fatal(err)
	}
	http.Redirect(w, r, "/view", http.StatusFound)
}

func main() {
	http.HandleFunc("/view", viewHandler)
	http.HandleFunc("/view/create", createHandler)
	fmt.Println("Server Start Up........")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}