package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
)

var (
	reactionsListUrl   string = "https://slack.com/api/reactions.list"
	channelsListUrl    string = "https://slack.com/api/channels.list"
	chatPostMessageUrl string = "https://slack.com/api/chat.postMessage"
	reactionList       []Reaction
	cursor             string = "first cursor"
	token              string = os.Getenv("SLACK_TOKEN")
	slack_channel      string = os.Getenv("SLACK_CHANNEL")
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

type ChannelListResponse struct {
	Channels []Channel `json:"channels"`
}

type Channel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Reaction struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type Emoji struct {
	Key   string
	Value int
}

type EmojiList []Emoji

func (p EmojiList) Len() int           { return len(p) }
func (p EmojiList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p EmojiList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func main() {
	if token == "" {
		log.Fatal("SLACK_TOKEN environment variable should be set")
	}

	if slack_channel == "" {
		slack_channel = "general"
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

	//fmt.Println(len(reactions))
	//for key, value := range reactions {
	//	fmt.Println(key + " : " + strconv.Itoa(value))
	//}

	emojiList := rankByEmojiCount(reactions)
	for _, emoji := range emojiList {
		fmt.Println(emoji.Key + " : " + strconv.Itoa(emoji.Value))
	}

	fmt.Println(getChannelID())
}

func rankByEmojiCount(reactions map[string]int) EmojiList {
	emojiList := make(EmojiList, len(reactions))
	i := 0
	for k, v := range reactions {
		emojiList[i] = Emoji{k, v}
		i++
	}
	sort.Sort(sort.Reverse(emojiList))
	return emojiList
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

func getChannelID() string {
	req, err := http.NewRequest("GET", channelsListUrl, nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	q := req.URL.Query()
	q.Add("token", token)

	req.URL.RawQuery = q.Encode()
	fmt.Println(req.URL.String())

	//resp, err := http.Get(reactionsListUrl + "?" + values.Encode())
	resp, err := http.Get(req.URL.String())
	if err != nil {
		fmt.Println(err)
		return ""
	}

	defer resp.Body.Close()

	response := &ChannelListResponse{}
	err = json.NewDecoder(resp.Body).Decode(response)

	target_channel_id := ""
	for _, channel := range response.Channels {
		if channel.Name == slack_channel {
			target_channel_id = channel.ID
		}
	}

	return target_channel_id
}
