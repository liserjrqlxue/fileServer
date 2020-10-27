package router

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/liserjrqlxue/goUtil/cryptoUtil"
	"github.com/liserjrqlxue/goUtil/simpleUtil"
)

type Info struct {
	Src     string
	Token   string
	Message string
}

var Public = "."

func Mp4play(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	fmt.Println("r.Form:\t", r.Form)
	fmt.Println("path:\t", r.URL.Scheme)
	fmt.Println("url_long:\t", r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("\tkey:\t", k)
		fmt.Println("\tval:\t", strings.Join(v, " "))
	}

	// token
	crutime := time.Now().Unix()
	token := cryptoUtil.Md5sum(strconv.FormatInt(crutime, 10))
	fmt.Println("token:\t", token)
	t, err := template.ParseFiles("template/mp4play.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	var src Info
	//{Src:r.Form["file"][0],Token:token}
	src.Token = token
	if len(r.Form["file"]) > 0 {
		src.Src = r.Form["file"][0]
	}
	if t != nil {
		err = t.Execute(w, src)
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
	} else {
		http.Error(w, "template/mp4play parse failed!", 500)
	}
}

func Upload(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path, " method:", r.Method)

	// token
	crutime := time.Now().Unix()
	token := cryptoUtil.Md5sum(strconv.FormatInt(crutime, 10))
	fmt.Println("token:\t", token)
	var src Info
	src.Token = token
	t, err := template.ParseFiles("template/upload.html")
	simpleUtil.CheckErr(err)

	if r.Method == "POST" {
		err = r.ParseMultipartForm(32 << 20)
		if err != nil {
			log.Print(err)
			return
		}
		var dest string
		if len(r.Form["dest"]) > 0 {
			dest = r.Form["dest"][0]
			if dest == "" {
				dest = "."
			}
			simpleUtil.CheckErr(os.MkdirAll(dest, 0755))
		}
		var file, handler, err = r.FormFile("uploadfile")
		simpleUtil.CheckErr(err)
		defer simpleUtil.DeferClose(file)
		//Info.Message=fmt.Sprint(handler.Header)
		var uploadFile = path.Join(dest, handler.Filename)
		simpleUtil.CheckErr(err, "create error")
		f, err := os.Create(uploadFile)
		defer simpleUtil.DeferClose(f)
		_, err = io.Copy(f, file)
		simpleUtil.CheckErr(err)
		src.Message = "upload succeed"
		src.Src = uploadFile
	}
	simpleUtil.CheckErr(t.Execute(w, src))
}

func Download(w http.ResponseWriter, r *http.Request) {
	var relPath = filepath.Join(Public, r.URL.Path)
	log.Printf("[%s]\t->\t[%s]\n", r.URL.Path, relPath)
	http.ServeFile(w, r, relPath)
}
