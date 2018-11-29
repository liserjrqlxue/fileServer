package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"crypto/md5"
	"time"
)

func md5sum(str string) string {
	byteStr := []byte(str)
	sum := md5.Sum(byteStr)
	sumStr := fmt.Sprintf("%x", sum)
	return sumStr
}

type MP4 struct{
    Src string
    Token string
}

func mp4play(w http.ResponseWriter, r *http.Request){
    r.ParseForm()
    fmt.Println("r.Form:\t",r.Form)
    fmt.Println("path:\t",r.URL.Scheme)
    fmt.Println("url_long:\t",r.Form["url_long"])
    for k,v:=range r.Form{
	fmt.Println("\tkey:\t",k)
	fmt.Println("\tval:\t",strings.Join(v," "))
    }

    // token
    crutime:=time.Now().Unix()
    token:=md5sum(strconv.FormatInt(crutime,10))
    fmt.Println("token:\t",token)
    t,_:=template.ParseFiles("template/mp4play.gtpl")
    var src MP4
    //{Src:r.Form["file"][0],Token:token}
    src.Token=token
    if len(r.Form["file"])>0{
	src.Src=r.Form["file"][0]
    }
    t.Execute(w,src)
}

func main() {
	http.HandleFunc("/mp4", mp4play)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Path := r.URL.Path
		path := fmt.Sprintf("%s", Path)
		fmt.Println(Path)
		http.ServeFile(w, r, "./"+path)

	}) //设置访问的路由
	fmt.Println("start")
	err := http.ListenAndServe(":9093", nil) //设置监听的端口
	if err != nil {
		log.Fatal(err)
	}
}
