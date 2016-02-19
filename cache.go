package main

import (
	"container/list"
	"fmt"
)


// DO NOT CHANGE THIS CACHE SIZE VALUE
const CACHE_SIZE int = 3

//Global Vars
var l = list.New()
var m = make(map[int]int)
var size int



func push(val int){
	l.PushFront(val)
}


func pop(){
	l.Remove(l.Back())
}

func mostRecent(val int){

	for e:=l.Front(); e != nil; e = e.Next() {
		if e.Value == val{
			l.MoveToFront(e)
			break
		}
	}
}




func Set(key int, value int) {
	// TODO: add your code here!
	var val int 
	val=value
	if l.Len()>=CACHE_SIZE{
		e:=l.Back()
		for key, value := range m {
    		if value==e.Value{
    			delete(m, key)
    			break
    		}
		}

		pop()
	}
	m[key]=val
	l.PushFront(val)
	fmt.Println("Key and Value added")
}

func Get(key int) int {
	// TODO: add your code here!
	
	i:=m[key]
	if i==0 {
		fmt.Println("Key value not found. ")
		return -1
	}

	if i!=0{

		fmt.Println("Value is found. It is:")
		fmt.Println(i)
		mostRecent(i)
		return i
	}
	return -1

}
