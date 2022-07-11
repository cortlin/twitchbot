package main

/*
	1.[-] Connect to the twitch IRC server via SSL connection
	2.[x] Send the PASS and NICK to authenticate the bot user
	3.[ ] Join a specific chat room, in this case okayotter
	4.[ ] once connected, read the chat and log the results to the terminal
	5.[ ] parse the results for commands so we can send a response to a specific command
	6.[x] create a function to call the ChuckNorris API
	7.[ ] store the result, and send it back to the twitch channel
	8.[ ] error handling all along the way
*/

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

//Object to store the ChuckNorris Joke
type ChuckNorris struct {
	Value string `json:"value"`
}

type chatConnection struct {
	connection net.Conn
}

func main() {
	//channel := "okayotter"
	username := "okayotterbot"
	//The auth key ideally should be stored in an env or config file
	oauth := "oauth:Removed before pushing"
	connectToTwitch(username, oauth)
	fmt.Println(getJoke().Value)
}

//call the Chuck Norris API and save to chuck1 object/struct
func getJoke() *ChuckNorris {
	response, err := http.Get("https://api.chucknorris.io/jokes/random")
	if err != nil {
		log.Fatal("Something went wrong: unable to retrieve a joke")
	}
	//TODO: Probably don't ignore the error
	responseBody, _ := ioutil.ReadAll(response.Body)
	chuck1 := &ChuckNorris{}
	json.Unmarshal([]byte(responseBody), &chuck1)
	return chuck1
}

//Connect to the twitch IRC service before connecting to a channel
func connectToTwitch(username, oauth string) *chatConnection {
	//connecting over non-SSL for now
	response, err := net.Dial("tcp", "irc.chat.twitch.tv:6667")
	if err != nil {
		log.Fatal("cannot connect to twitch irc server", err)
	}

	//store the connection data into an object
	irc := &chatConnection{
		connection: response,
	}

	//must send in this order according to twitch
	irc.send("PASS " + oauth)
	irc.send("NICK " + username)
	fmt.Printf("[%s] connected to twitch irc service\n", time.Now().Format(time.UnixDate))
	return irc
}

//Send a message to the twitch server
func (irc *chatConnection) send(msg string) {
	_, err := irc.connection.Write([]byte(msg + "\r\n"))
	if err != nil {
		log.Fatal("disconnected from twitch irc server send", err)
	}
}
