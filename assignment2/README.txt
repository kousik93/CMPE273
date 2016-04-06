I have made 2 copies of the server to showcase their working. They are in folders server1, server2.

I have also attached a screenshot of my test run. It is `testrun.png`

* sever1/server1.go and server2/server2.go are the same exact files. Only renamed for readability.

* server1/app1.toml and server2/app2.toml are different files with different port numbers

* To run my application (Run server1.go, server2.go on different terminals) (You HAVE to be in the respective folders to run. Dont do `go run server1/server1.go`. It wont work):

~/$ go run server1.go app1.toml

~/$ go run server2.go app2.toml

* Now both the servers will be running. You will notice 2 more files created in each server folder. They are the respective local databases for that server.

* The database file is userdetails and the collection file is userdetails_details

* Now that both the servers are running, you can send HTTP GET, POST, PUT, DELETE to either of the REST API ports and you will find that both the databases are always consistant after every request. 

* You can diff userdetails_details in both the folders after every request to check consistancy. (I used the test file fromt the previous assignment to make the requests and afterwards did a diff and checked)

* That test file I used for making HTTP requests is the grader.py file
