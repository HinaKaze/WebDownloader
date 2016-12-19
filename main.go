package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	log.Println("Downloader start")
	http.HandleFunc("/", welcome)
	http.HandleFunc("/download", download)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err.Error())
	}
}

func welcome(resp http.ResponseWriter, req *http.Request) {
	tmpl, err := template.ParseFiles("welcome.html")
	if err != nil {
		panic(err.Error())
	}
	tmpl.Execute(resp, nil)
	return
}

func download(resp http.ResponseWriter, req *http.Request) {
	urlBytes := make([]byte, 1024)
	offset, err := req.Body.Read(urlBytes)
	fmt.Println(offset)
	if err != nil && err != io.EOF {
		panic(err.Error())
	}
	fmt.Println(string(urlBytes[:offset]))
	go startDownload(string(urlBytes[:offset]))
	return
}

func startDownload(ustr string) {
	resp, err := http.Get(ustr)
	defer resp.Body.Close()
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println("ContentLength:", resp.ContentLength)
	log.Println(resp.Header)
	contentType := resp.Header["Content-Type"][0]
	if contentTypeCheck(contentType) {
		fileName, err := url.QueryUnescape(ustr[strings.LastIndex(ustr, "/")+1:])
		if err != nil {
			log.Println(err.Error())
			return
		}
		log.Println("FileName:", fileName)
		file, err := os.Create("./files/" + fileName)
		defer file.Close()
		if err != nil {
			log.Println(err.Error())
			return
		}
		bodyBytes := make([]byte, 1024)
		n, err := resp.Body.Read(bodyBytes)
		log.Println("Downloading file size :", n, " Bytes")
		for n > 0 || (err != nil && err != io.EOF) {
			if err != nil && err != io.EOF {
				log.Println(err.Error())
				return
			}
			log.Println("Downloading file :", fileName, "size :", n, " Bytes")
			file.Write(bodyBytes[:n])
			n, err = resp.Body.Read(bodyBytes)
		}
		log.Println("Download file :", fileName, "end!")
	} else {
		log.Println("Resp.Header Content-Type is not  matched :", contentType)
	}
	return
}

func contentTypeCheck(contentType string) bool {
	switch contentType {
	case "application/octet-stream", "text/plain; charset=utf-8", "application/pdf":
		return true
	default:
		return false
	}

	return true
}
