package main

import (
	"net/http"
	"fmt"
	"strings"
)

func main(){
	reader := strings.NewReader("行业交流群")
//	csv.Reader{}
	fmt.Println(strings.Index("行业交流群","交流"))
	fmt.Println(strings.IndexRune("行业交流群",'交'))
	a, s, e := reader.ReadRune()
	fmt.Println(string(a),s,e)
//	http.HandleFunc("/", handler1)
//	http.ListenAndServe(":8080", nil)
}

func handler1(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, " Request url is %s", r.RequestURI)
}

