package router

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	simpleUtil "github.com/liserjrqlxue/simple-util"
)

func md5sum(str string) string {
	byteStr := []byte(str)
	sum := md5.Sum(byteStr)
	sumStr := fmt.Sprintf("%x", sum)
	return sumStr
}

type Info struct {
	Src     string
	Token   string
	Message string
}

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
	token := md5sum(strconv.FormatInt(crutime, 10))
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
	token := md5sum(strconv.FormatInt(crutime, 10))
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
		file, handler, err := r.FormFile("uploadfile")
		simpleUtil.CheckErr(err)
		defer simpleUtil.DeferClose(file)
		//Info.Message=fmt.Sprint(handler.Header)
		uploadFile := path.Join("public", "upload", handler.Filename)
		f, err := os.Create(uploadFile)
		simpleUtil.CheckErr(err)
		defer simpleUtil.DeferClose(f)
		_, err = io.Copy(f, file)
		simpleUtil.CheckErr(err)
		src.Message = "upload succeed"
		src.Src = uploadFile
	}
	simpleUtil.CheckErr(t.Execute(w, src))
}
