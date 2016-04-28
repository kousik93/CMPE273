package main

import (
	"log"
    "encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"github.com/gorilla/mux"
)

type Obj struct{
	Key string `json:"key"`
	Value string `json:"value"`
}

type MainStruct struct{	
	obj Obj
	Data map[string]Obj
}

func NewMainStruct() *MainStruct {
    var m MainStruct
    m.Data = make(map[string]Obj)
 
    return &m
}

func (m MainStruct) PutData(w http.ResponseWriter, r *http.Request) {
	vars:=mux.Vars(r)
	key := vars["key"]
	value := vars["value"]

	m.obj.Key=key
	m.obj.Value=value
	m.Data[key] = m.obj

	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte(""))
	log.Println("Put SuccessFull.")

}

func (m MainStruct) GetData(w http.ResponseWriter, r *http.Request) {
	vars:=mux.Vars(r)
	key := vars["key"]

	a, err := json.Marshal(m.Data[key])
	if err != nil {
		log.Println("error:", err)
	}
	w.Write([]byte(a))	
	log.Println("Get SuccessFull.")
}

func (m MainStruct) GetDataAll(w http.ResponseWriter, r *http.Request) {
	var ooo []Obj
	for k, _ := range m.Data { 
		ooo=append(ooo,m.Data[k])
	 }	
	a, _ := json.Marshal(ooo)
	log.Println(string(a))
	w.Write([]byte(a))
 }

func runServers(hport int){
	m:=NewMainStruct()

	router:=mux.NewRouter()
	router.HandleFunc("/{key}/{value}",m.PutData).Methods("PUT")
	router.HandleFunc("/",m.GetDataAll).Methods("GET")
	router.HandleFunc("/{key}",m.GetData).Methods("GET")
	http.ListenAndServe(":"+strconv.Itoa(hport), router)
}

func main(){
	arg := os.Args[1]
    s:=strings.Split(arg,"-")
	startp,endp:=s[0],s[1]
	startport,_ := strconv.Atoi(startp)
	endport,_ :=strconv.Atoi(endp)
	for i:=startport;i<=endport;i++ {
		go runServers(i)
	}
		log.Println("Running server on:"+strconv.Itoa(startport)+"-"+strconv.Itoa(endport))
	select{}
}