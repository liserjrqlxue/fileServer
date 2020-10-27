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

var (
	PublicPath   = "public"
	UploadPath   = "upload"
	TemplatePath = "template"
)

func Mp4play(w http.ResponseWriter, r *http.Request) {
	var (
		t   *template.Template
		src Info
		err error
	)
	err = r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Println("r.Form:\t", r.Form)
	fmt.Println("path:\t", r.URL.Scheme)
	fmt.Println("url_long:\t", r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("\tkey:\t", k)
		fmt.Println("\tval:\t", strings.Join(v, " "))
	}

	// token
	var crutime = time.Now().Unix()
	token := cryptoUtil.Md5sum(strconv.FormatInt(crutime, 10))
	fmt.Println("token:\t", token)

	//{Src:r.Form["file"][0],Token:token}
	src.Token = token
	if len(r.Form["file"]) > 0 {
		src.Src = r.Form["file"][0]
	}

	t, err = template.ParseFiles(filepath.Join(TemplatePath, "mp4play.html"))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if t != nil {
		err = t.Execute(w, src)
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
	} else {
		http.Error(w, "mp4play parse failed!", 500)
	}
}

func Upload(w http.ResponseWriter, r *http.Request) {
	var (
		t   *template.Template
		src Info
		err error
	)
	log.Println(r.URL.Path, " method:", r.Method)

	// token
	var crutime = time.Now().Unix()
	var token = cryptoUtil.Md5sum(strconv.FormatInt(crutime, 10))
	log.Println("token:\t", token)
	src.Token = token

	t, err = template.ParseFiles(filepath.Join(TemplatePath, "upload.html"))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if t != nil {
		if r.Method == "POST" {
			err = r.ParseMultipartForm(32 << 20)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			var dest, destUrl string
			if len(r.Form["dest"]) > 0 {
				destUrl = r.Form["dest"][0]
				if destUrl == "" {
					destUrl = UploadPath
				}
				dest = filepath.Join(PublicPath, destUrl)
				err = os.MkdirAll(dest, 0755)
				if err != nil {
					http.Error(w, err.Error(), 500)
					return
				}
			}
			file, handler, err := r.FormFile("uploadfile")
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			defer simpleUtil.DeferClose(file)
			//Info.Message=fmt.Sprint(handler.Header)
			var uploadFile = path.Join(dest, handler.Filename)
			simpleUtil.CheckErr(err, "create error")
			f, err := os.Create(uploadFile)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			defer simpleUtil.DeferClose(f)
			_, err = io.Copy(f, file)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			src.Message = "upload succeed"
			src.Src = path.Join(destUrl, handler.Filename)
		}
		err = t.Execute(w, src)
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
	} else {
		http.Error(w, "upload parse failed!", 500)
	}
}

func Download(w http.ResponseWriter, r *http.Request) {
	var relPath = filepath.Join(PublicPath, r.URL.Path)
	log.Printf("[%s]\t->\t[%s]\n", r.URL.Path, relPath)
	http.ServeFile(w, r, relPath)
}
