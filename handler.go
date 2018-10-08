package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

const (
	GET  string = "GET"
	POST string = "POST"
)

func Route(r *mux.Router, path string, handler func(http.ResponseWriter, *http.Request), method string) {
	r.HandleFunc(path, handler).Methods(method)
}

func Run(r *mux.Router, address string, port int) {
	server := http.Server{
		Addr:    fmt.Sprintf("%s:%d", address, port),
		Handler: r,
	}

	log.Println("Server started:" + server.Addr)
	log.Fatal(server.ListenAndServe())
}

func MainHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Who am I ?"))
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	var templatefile = template.Must(template.ParseFiles("upload.html"))
	templatefile.Execute(w, "index.html")
}

func redirectToErrorPage(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/error/", http.StatusFound)
}

func errorPageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", "<p>Internal Server Error</p>")
}

func RemoveContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

func UploadFile(w http.ResponseWriter, r *http.Request) {

	data1, err := ioutil.ReadFile("/home/g23516033/anonimserve/anonim1/tor.url")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	url1 := string(data1)

	data2, err := ioutil.ReadFile("/home/g23516033/anonimserve/anonim2/tor.url")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	url2 := string(data2)

	data3, err := ioutil.ReadFile("/home/g23516033/anonimserve/anonim3/tor.url")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	url3 := string(data3)

	var m1 = [3][2]string{{"anonim1/upload", url1}, {"anonim2/upload", url2}, {"anonim3/upload", url3}}
	i := rand.Intn(3)
	dirName := fmt.Sprintf("/home/g23516033/anonimserve/%s/", m1[i][0])
	urlName := fmt.Sprintf("%s.onion", m1[i][1])

	err = RemoveContents(dirName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	reader, err := r.MultipartReader()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}

		if part.FileName() == "" {
			continue
		}

		uploadedFile, err := os.Create(dirName + part.FileName())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			uploadedFile.Close()
			redirectToErrorPage(w, r)
			return
		}

		_, err = io.Copy(uploadedFile, part)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			uploadedFile.Close()
			redirectToErrorPage(w, r)
			return
		}

	}

	fmt.Fprintf(w, "<p>Your Url: %s</p><p><a href=\"/upload/\">Click to Upload</a></p>", urlName)
}
