package main

/*
	1.[-] Connect to the twitch IRC server via SSL connection
	2.[x] Send the PASS and NICK to authentica the bot user
	3.[x] Join a specific chat room, in this case okayotter
	4.[x] once connected, read the chat and log the results to the terminal
	5.[ ] Parse the results for commands so we can send a response to a specific command
	6.[ ] Reply to PING from the server with PONG to keep the connection alive
	6.[x] create a fuction to call the ChuckNorris API
	7.[ ] store the result, and send it back to the twitch channel
	8.[ ] error handling all along the way
*/

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/textproto"
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
	channel := "okayotter"
    username := "okayotterbot"
    //The auth key ideally should be stored in an env or config file
    oauth := "oauth: removed"
	chat := connectToTwitch(username, oauth)
	chat.joinChannel(channel)

	reader := bufio.NewReader(chat.connection)
    tp := textproto.NewReader(reader)
    //Infinite loop to continuously check for new messages
    for {
        response, err := tp.ReadLine()
        if err != nil {
            log.Fatal("disconnected from twitch channel. Reason: ", err)
        }
        //Print all of the messages received from twitch to the console
        fmt.Printf("[%s] %s\r\n", time.Now().Format(time.UnixDate), response)
	}
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
        log.Fatal("oops. connection to twitch irc server failed: ", err)
    }

    //store the connection data into an object
    irc := &chatConnection{
        connection: response,
    }

	//must send in this order according to twitch
    irc.sendMessage("PASS " + oauth)
    irc.sendMessage("NICK " + username)
    fmt.Printf("[%s] connected to twitch irc service\n", time.Now().Format(time.UnixDate))
    return irc
}

//Send a message to the twitch server
func (irc *chatConnection) sendMessage(msg string) {
    _, err := irc.connection.Write([]byte(msg + "\r\n"))
    if err != nil {
        log.Fatal("Couldn't send a message to the IRC service. Reason: ", err)
    }
}

//Join the selected twitch channel
func (irc *chatConnection) joinChannel(channel string) {

	//Send the JOIN message to the Twitch IRC service to join the server
    irc.sendMessage("JOIN #" + channel)
    fmt.Printf("[%s] joined channel %s\n", time.Now().Format(time.UnixDate), channel)
}