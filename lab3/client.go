package main

import (
	"fmt"
	"net/http"
	"strings"
	"os"
	"strconv"
    "github.com/huichen/murmur"
)

type Obj struct{
	Key string `json:"key"`
	Value string `json:"value"`
	Mapto uint32 `json:mapto`
}

type Ring struct {
	Nodes Nodes
}

type Node struct {
	Id     uint32
}

func (n Nodes) Len() int {
 return len(n) 
}

func NewNode(id uint32) *Node {
	return &Node{
		Id:     id,
	}
}

func NewRing() *Ring {
	return &Ring{Nodes: Nodes{}}
}

func (r *Ring) AddNode(id uint32) {
	node := NewNode(id)
	r.Nodes = append(r.Nodes, node)
}

func weight (node uint32, key string) uint32{
    a := uint32(1103515245)
    b := uint32(1455)
    hash := hashId(key)
    return (a * ((a * node + b) ^ hash) + b) % (2^31)
}

func (r *Ring) Get(key string) uint32{
 
        var weights []uint32
        for i:=0;i<len(r.Nodes);i++ {
            w := weight(r.Nodes[i].Id, key)
            weights=append(weights, w)
        }
		//Calculate the biggest hash and return the port number        
        biggest := weights[0]
        j:=0
		for i:=0;i<len(weights);i++ {
			if weights[i] > biggest {
				biggest = weights[i]
				j=i
			}
		}
        return r.Nodes[j].Id
}

func hashId(text string) uint32 {
 	return murmur.Murmur3([]byte(text))
}

func doPut(url string) {
	client := &http.Client{}
	request, _ := http.NewRequest("PUT", url, strings.NewReader(""))
	request.ContentLength = 0
	client.Do(request) 
}

type Nodes []*Node
var DataAll []Obj
var data Obj

func main(){
	//Start Ring
	r := NewRing()

	ports := os.Args[1]
	datas := os.Args[2]

	//Get Ports
	s:=strings.Split(ports,"-")
	startp,endp:=s[0],s[1]
	startport,_:=strconv.Atoi(startp)
	endport,_:=strconv.Atoi(endp)

	//Insert ports as nodes in Ring
	for i:=startport; i<=endport; i++ {
		r.AddNode(uint32(i))
	}

	//Get Data
	dataset:=strings.Split(datas,",")

	for i:=0; i<len(dataset); i++ {
		tmp:=strings.Split(dataset[i],"->")
		data.Key=tmp[0]
		data.Value=tmp[1]
		data.Mapto=r.Get(data.Key)
		DataAll = append(DataAll,data)
	} 
	
	//Do PUT requests
	for i:=0; i<len(DataAll); i++ {
		var baseUrl="http://localhost:"
		baseUrl=baseUrl+strconv.Itoa(int(DataAll[i].Mapto))+"/"+DataAll[i].Key+"/"+DataAll[i].Value
		fmt.Println(baseUrl)
		doPut(baseUrl)
	}
}