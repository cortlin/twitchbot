package main

/*
	1.[ ] Connect to the twitch IRC server via SSL connection
	2.[ ] Send the PASS and NICK to authentica the bot user
	3.[ ] Join a specific chat room, in this case okayotter
	4.[ ] once connected, read the chat and log the results to the terminal
	5.[ ] Parse the results for commands so we can send a response to a specific command
	6.[x] create a fuction to call the ChuckNorris API
	7.[ ] store the result, and send it back to the twitch channel
	8.[ ] error handling all along the way
*/

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

//Object to store the ChuckNorris Joke
type ChuckNorris struct {
    Value string `json:"value"`
	err string
}

func main() {
	joke := getJoke()
	fmt.Println(joke)
}

//call the Chuck Norris API and return a joke as a string
func getJoke() string {
    response, err := http.Get("https://api.chucknorris.io/jokes/random")
    if err != nil {
        fmt.Println("Something went wrong: unable to retrieve a joke")
        return ("Oops, something went wrong. Could not retrieve joke")
    }
	//TODO: Probably don't ignore the error
    responseBody, _ := ioutil.ReadAll(response.Body)
    chuck1 := ChuckNorris{}
    json.Unmarshal([]byte(responseBody), &chuck1)
    return chuck1.Value
}