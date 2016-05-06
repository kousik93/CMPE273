package main

import (
    "fmt"
    "os"
    "strconv"
    "strings"
    // /"time"
    "encoding/json"
    "io/ioutil"
    "net/http"
    //"gopkg.in/mgo.v2/bson"
    "bytes"
    "github.com/jasonlvhit/gocron"
)

type InputBSON struct{
  Id int64 	`json:"id"`
  Success_http_response_code int 	`json:"success_http_response_code"`
  Max_retries int `josn:"max_retries"`
  Callback_webhook_url string `josn:"callback_webhook_url"`
	  Request struct{
		Url string 	`json:"url"`
		Method string 	`json:"method"`
		Http_headers struct{
				Content_Type string `json:"Content-Type"`
				Accept string `json:"Accept"`
			}`json:"http_headers"`
		Body struct{
			Foo string `json:"foo"`
		}`json:"body"`
	} `json:"request"`		

}

type OutputBSON struct{
    Job struct{
    	Status string `json:"status"`
    	Num_retries int `json:"num_retries"`
    }`json:"job"`
    Input InputBSON `json:"input"`
    Output HttpResponse `json:"output"`
    Callback_response_code int `json:"callback_response_code"`
}


type HttpResponse struct {
	Response struct{
		Http_response_code int `json:"http_response_code"`
		Http_headers struct{
			Date string `json:"Date"`
			Content_Type string `json:"Content-Type"`
			Content_Length int `json:"Content-Length"`
		}`json:"response"`
		Body struct{
			Hello string `json:"hello"`
		}`json:"body"`
	}
}

var input InputBSON
var httpResponse HttpResponse
var outputBSON OutputBSON

func doRequest(){
	//Start Make initial Request
		url :=input.Request.Url
	    fmt.Println("Making request to: ", url)
	    jsonStr,_:=json.Marshal(input.Request.Body)
	    jsonBody := []byte(jsonStr)
	    method:=input.Request.Method
	    req, _ := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	    content_type:=input.Request.Http_headers.Content_Type
	    accept:=input.Request.Http_headers.Accept
	    req.Header.Set("Accept", accept)
	    req.Header.Set("Content-Type", content_type)
	    client := &http.Client{}
	    resp, _ := client.Do(req)
	    defer resp.Body.Close()
    //End Make initial Request

    //Start filling HttpResponse
	    	//Status
	    s:=strings.Split(resp.Status," ")
		Status,_:=strconv.Atoi(s[0])
		httpResponse.Response.Http_response_code=Status
			//Date
	    Date:=resp.Header["Date"][0]
	    httpResponse.Response.Http_headers.Date=Date
	    	//Content-Type
	    Content_type:=resp.Header["Content-Type"][0]
	    httpResponse.Response.Http_headers.Content_Type=Content_type
	    	//Content-Length
	    Content_length:=resp.Header["Content-Length"][0]
	    httpResponse.Response.Http_headers.Content_Length,_=strconv.Atoi(Content_length)
	    	//Body
	    body, _ := ioutil.ReadAll(resp.Body)
	    httpResponse.Response.Body.Hello=  string(body)
	//End filling HttpResponse
	
	//Start fill Input and Output in OutputBSOn
		outputBSON.Job.Num_retries=0
		outputBSON.Input=input
		outputBSON.Output=httpResponse
		outputBSON.Job.Num_retries=outputBSON.Job.Num_retries+1
	//End fill Input and Output in OutputBSOn

	//If success do callback webhook
	    if resp.Status == "200 OK"{
	    	outputBSON.Job.Status="COMPLETED"
	    	//Do callback
	    	req,_ := http.NewRequest("POST", input.Callback_webhook_url, bytes.NewBuffer([]byte("")))
	    	client := &http.Client{}
	    	resp,_ := client.Do(req)
	    	s:=strings.Split(resp.Status," ")
			Status,_:=strconv.Atoi(s[0])
	    	outputBSON.Callback_response_code=Status
	    }
	//End Callback

	//If failure, put Still Trying and write to output. Leave callback response code empty
	    if resp.Status != "200 OK"{
	    	outputBSON.Job.Status="STILL_TRYING"
	    }
	//End if failure

	//Write to output
	    a,_:=json.Marshal(outputBSON)
	    writeFile([]byte(a),"output.bson")
}




func writeFile(data []byte, fileName string)  {
    file, _ := os.Create(fileName)
    defer file.Close()   
    numBytes, _ := file.Write(data)
    fmt.Printf("\nwrote %d bytes to %s\n", numBytes, fileName)
    file.Sync()
}
func task() {
	//Get input.bson
	arg := os.Args[1]
	f, _ := os.Open(arg)
	defer f.Close()
	buf, _ := ioutil.ReadAll(f)
    json.Unmarshal(buf, &input); 

    //Do the required Http request
    doRequest()





    // t := makeTimestamp()
    // fmt.Printf("Running task...@%d\n", t)
    // data := createBSON()
    // writeFile(data, "data.bson")
    // readFile("data.bson")
}

func main() {

	

    s := gocron.NewScheduler()
    s.Every(5).Seconds().Do(task)
    <- s.Start()
}
