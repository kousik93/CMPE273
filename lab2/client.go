package main

import (
	"sort"
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
	Mapto string `json:mapto`
}

type Ring struct {
	Nodes Nodes
}

type Node struct {
	Id     string
	HashId uint32
}

func (n Nodes) Len() int {
 return len(n) 
}

func (n Nodes) Swap(i, j int) {
 n[i], n[j] = n[j], n[i] 
}

func (n Nodes) Less(i, j int) bool {
 return n[i].HashId < n[j].HashId 
}

func NewNode(id string) *Node {
	return &Node{
		Id:     id,
		HashId: hashId(id),
	}
}

func NewRing() *Ring {
	return &Ring{Nodes: Nodes{}}
}

func (r *Ring) AddNode(id string) {
	node := NewNode(id)
	r.Nodes = append(r.Nodes, node)
	sort.Sort(r.Nodes)
}

func (r *Ring) Get(id string) string {
	searchfn := func(i int) bool {
		return r.Nodes[i].HashId >= hashId(id)
	}

	i := sort.Search(r.Nodes.Len(), searchfn)
	if i >= r.Nodes.Len() {
		i = 0
	}

	return r.Nodes[i].Id
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
		r.AddNode(strconv.Itoa(i))
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
		baseUrl=baseUrl+DataAll[i].Mapto+"/"+DataAll[i].Key+"/"+DataAll[i].Value
		fmt.Println(baseUrl)
		doPut(baseUrl)
	}
}