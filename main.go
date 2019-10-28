package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/nlopes/slack"
)

var atReplies = []string{
	"stop using my name to get attention",
	"hey there!",
	"howdy friends",
}

type MessageResponder interface {
	prepMessageResponse(incoming string, channel string) slack.OutgoingMessage
}

type messageResponder struct {
	RTM slack.RTM
}

func (mr messageResponder) prepMessageResponse(incoming string, channel string) slack.OutgoingMessage {
	processedString := cleanString(incoming)
	words := strings.Fields(strings.ToLower(processedString))

	switch {
	case containsWord(words, "hello") || containsWord(words, "hi"):
		return *mr.RTM.NewOutgoingMessage("hola!", channel)
	case containsWord(words, "go"):
		return *mr.RTM.NewOutgoingMessage("Go is the best programming language", channel)
	}

	return slack.OutgoingMessage{}
}

func main() {
	api := slack.New(os.Getenv("SLACK_KEY"))

	// Here we make a new real-time messanger
	rtm := api.NewRTM()
	mr := messageResponder{RTM: *rtm}

	// Fire a goroutine to manage the connection
	go rtm.ManageConnection()

	// Range over incoming messages
	for msg := range rtm.IncomingEvents {
		fmt.Println("Event Received: ", msg.Type)

		switch ev := msg.Data.(type) {

		case *slack.MessageEvent:

			botTagString := fmt.Sprintf("<@%s>", rtm.GetInfo().User.ID)

			processedString := cleanString(ev.Msg.Text)
			words := strings.Fields(strings.ToLower(processedString))

			if strings.Contains(ev.Msg.Text, botTagString) {
				greetDex := randomFromZeroTo(len(atReplies))
				rtm.SendMessage(rtm.NewOutgoingMessage(atReplies[greetDex], ev.Channel))
				continue
			}

			if containsWord(words, "rip") {
				msgRef := slack.NewRefToMessage(ev.Channel, ev.Timestamp)
				if err := api.AddReaction("rip", msgRef); err != nil {
					fmt.Printf("Error adding reaction: %s\n", err)
					return
				}
			}

			reply := mr.prepMessageResponse(ev.Msg.Text, ev.Channel)
			if len(reply.Text) > 0 {
				rtm.SendMessage(&reply)
			}
		// case *slack.UserTypingEvent:
		// 	rtm.SendMessage(rtm.NewOutgoingMessage("type faster you worm!", ev.Channel))
		default:

		}
	}
}
