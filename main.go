package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

var (
	reactionsListUrl string = "https://slack.com/api/reactions.list"
	reactionList     []Reaction
	cursor           string = "first cursor"
	token            string = os.Getenv("SLACK_TOKEN")
)

type Response struct {
	ResponseMetadata ResponseMetadata `json:"response_metadata"`
	Items            []Item           `json:"items"`
}

type ResponseMetadata struct {
	NextCursor string `json:"next_cursor"`
}

type Item struct {
	Type    string  `json:"type"`
	Channel string  `json:"channel"`
	Message Message `json:"message"`
}

type Message struct {
	Type      string     `json:"type"`
	Reactions []Reaction `json:"reactions"`
}

type Reaction struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func main() {
	if token == "" {
		log.Fatal("SLACK_TOKEN environment variable should be set")
	}

	for {
		if result := getReactions(); result {
			break
		}
	}

	reactions := map[string]int{}
	for _, reaction := range reactionList {
		count, ok := reactions[reaction.Name]
		if ok == false {
			reactions[reaction.Name] = reaction.Count
		} else {
			reactions[reaction.Name] = count + reaction.Count
		}
	}

	fmt.Println(len(reactions))
	for key, value := range reactions {
		fmt.Println(key + " : " + strconv.Itoa(value))
	}
}

func getReactions() bool {
	req, err := http.NewRequest("GET", reactionsListUrl, nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	q := req.URL.Query()
	q.Add("token", token)
	if cursor != "first cursor" {
		q.Add("cursor", cursor)
	}
	if cursor == "" {
		return true
	}
	req.URL.RawQuery = q.Encode()
	fmt.Println(req.URL.String())

	//resp, err := http.Get(reactionsListUrl + "?" + values.Encode())
	resp, err := http.Get(req.URL.String())
	if err != nil {
		fmt.Println(err)
		return true
	}

	defer resp.Body.Close()

	fmt.Println(resp.Body)

	response := &Response{}
	err = json.NewDecoder(resp.Body).Decode(response)

	for _, item := range response.Items {
		for _, reaction := range item.Message.Reactions {
			fmt.Println(reaction.Name)
			fmt.Println(reaction.Count)
			reactionList = append(reactionList, reaction)
		}
	}

	fmt.Println(len(response.Items))
	fmt.Println(len(reactionList))
	fmt.Println(response.ResponseMetadata.NextCursor)
	cursor = response.ResponseMetadata.NextCursor
	fmt.Println(cursor)
	return false
}
