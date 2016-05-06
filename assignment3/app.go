package main

import (
    "fmt"
    "os"
    "strconv"
    "strings"
    "encoding/json"
    "io/ioutil"
    "time"
    "net/http"
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
var requestsYetToSucceed int

func doInitialRequest(){
	//Start Make initial Request
		url :=input.Request.Url
	    fmt.Println("\nMaking initial request for: ", url)
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
	
	//Start fill Input, Output, Num_retries in OutputBSON
		outputBSON.Job.Num_retries=0
		outputBSON.Input=input
		outputBSON.Output=httpResponse
		outputBSON.Job.Num_retries=outputBSON.Job.Num_retries+1
	//End fill Input, Output, Num_retries in OutputBSON

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
	    	fmt.Println("Request was Completed Successfully. Status is COMPLETED")
	    }
	//End Callback webhook

	//If failure, put Still Trying and write to output. Leave callback response code empty
	    if resp.Status != "200 OK"{
	    	outputBSON.Job.Status="STILL_TRYING"
	    	requestsYetToSucceed=requestsYetToSucceed+1
	    	fmt.Println("Request was Unsuccessful. Status is STILL_TRYING. You have "+strconv.Itoa(input.Max_retries)+" retries left.")
	    }
	//End if failure

	//Write to output
	    a,_:=json.Marshal(outputBSON)
	    writeFile([]byte(a),"output.bson")
}

func doRetry(){
	f, _ := os.Open("output.bson")
	defer f.Close()
	buf, _ := ioutil.ReadAll(f)
    json.Unmarshal(buf, &outputBSON);

    if outputBSON.Job.Status=="STILL_TRYING" {
    	//Check if num of retries is expired
    	if outputBSON.Job.Num_retries==outputBSON.Input.Max_retries {
    		outputBSON.Job.Status="FAILED"
    		//Do Callback
    			req,_ := http.NewRequest("POST", outputBSON.Input.Callback_webhook_url, bytes.NewBuffer([]byte("")))
		    	client := &http.Client{}
		    	resp,_ := client.Do(req)
		    	s:=strings.Split(resp.Status," ")
				Status,_:=strconv.Atoi(s[0])
		    	outputBSON.Callback_response_code=Status
    		//End Callback

		    //Write to output.bson
		    	outputBSON.Input.Id=makeTimestamp()
		    	a,_:=json.Marshal(outputBSON)
	    		writeFile([]byte(a),"output.bson")
	    	//End write to output

	    	//Update counter
	    	requestsYetToSucceed=requestsYetToSucceed-1

	    	fmt.Println("\nSorry. Retries Maxed Out. FAILED.")
    	}

    	//Check if num of retries < max retries
    	if outputBSON.Job.Num_retries < outputBSON.Input.Max_retries {
    		//Start Make Retry Request
				url :=outputBSON.Input.Request.Url
			    fmt.Println("\nMaking Retry request to: ", url)
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
    		//End Make Retry Request

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

		    //Start fill Input, Output, Num_retries in OutputBSON
				outputBSON.Output=httpResponse
				outputBSON.Job.Num_retries=outputBSON.Job.Num_retries+1
			//End fill Input, Output, Num_retries in OutputBSON

			//If success do callback webhook
			    if resp.Status == "200 OK"{
			    	outputBSON.Job.Status="COMPLETED"
			    	//Do callback
			    	req,_ := http.NewRequest("POST", outputBSON.Input.Callback_webhook_url, bytes.NewBuffer([]byte("")))
			    	client := &http.Client{}
			    	resp,_ := client.Do(req)
			    	s:=strings.Split(resp.Status," ")
					Status,_:=strconv.Atoi(s[0])
			    	outputBSON.Callback_response_code=Status
			    	requestsYetToSucceed=requestsYetToSucceed-1
			    	fmt.Println("Request was completed Successfully. Status is COMPLETED.")
			    }
			//End Callback webhook

			//If failure, put Still Trying and write to output. Leave callback response code empty
			    if resp.Status != "200 OK"{
			    	outputBSON.Job.Status="STILL_TRYING"
			    	fmt.Println("Request was Unsuccessful. Status is STILL_TRYING. You have "+strconv.Itoa(outputBSON.Input.Max_retries-outputBSON.Job.Num_retries+1)+" retries left.")
			    }
			//End if failure

			//Write to output
			outputBSON.Input.Id=makeTimestamp()
	    	a,_:=json.Marshal(outputBSON)
	    	writeFile([]byte(a),"output.bson")
    	}
    } 	
}

//Helper Function to write to output.json
func writeFile(data []byte, fileName string)  {
    file, _ := os.Create(fileName)
    defer file.Close()   
    file.Write(data)
    file.Sync()
}

//Helper function for Timestamp
func makeTimestamp() int64 {
    return time.Now().UnixNano() / (int64(time.Millisecond)/int64(time.Nanosecond))
}

func task() {
	//Initialize
	requestsYetToSucceed=0
	//Get input.bson
	arg := os.Args[1]
	f, _ := os.Open(arg)
	defer f.Close()
	buf, _ := ioutil.ReadAll(f)
    json.Unmarshal(buf, &input); 
    //Do the required Http request
    doInitialRequest()
    //Do Retries
    for requestsYetToSucceed > 0 {
    	doRetry()
    }
    fmt.Println("\nAll requests are now COMPLETED or FAILED\n")
}

func main() {
    s := gocron.NewScheduler()
    s.Every(5).Seconds().Do(task)
    <- s.Start()
    task()
}
