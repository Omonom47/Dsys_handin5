## How to start Program

### Servers
To start the program it needs 3 servers being run on precisely port 5000, 5001, and 5002. Do this by writing the following 3 lines into their own seperate terminal
1. go run server/server.go -port 5000
2. go run server/server.go -port 5001
3. go run server/server.go -port 5002

### Clients
When the 3 servers are running write as many clients as you would like in their own seperate terminals with the command:
1. go run client/client.go -id 1

replace the 1 with any number of your choosing

### Functions for client
#### Bid
To write a bid write **"bid"** in the client terminal. A prompt will then come up where you can then write you own bid. Hereafter it will state wether or not the bid was succesful or not
#### Result
To get result write **"result"** in the client terminal and you will be prompted with the highest bid and wether or not the auction is still on going

