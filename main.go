package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc/v4"
	"github.com/joho/godotenv"
)

func getStreamProbability(day, month, weekday int) (float64, error) {
	if day < 1 || day > 31 || month < 1 || month > 12 || weekday < 0 || weekday > 6 {
		return 0.0, fmt.Errorf("invalid input")
	}
	predictions := score([]float64{float64(day), float64(month), float64(weekday)})
	return predictions[1], nil
}

func main() {
	godotenv.Load()
	twitchAccessToken := os.Getenv("TWITCH_ACCESS_TOKEN")
	if twitchAccessToken == "" {
		log.Fatal("TWITCH_ACCESS_TOKEN must be set")
	}
	client := twitch.NewClient("guciostreamai", "oauth:"+twitchAccessToken)
	client.OnConnect(func() {
		log.Println("connected to twitch")
	})
	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		if strings.TrimSpace(message.Message) == "!czydzisstream" {
			now := time.Now()
			day := now.Day()
			month := int(now.Month())
			weekday := int(now.Weekday()) - 1
			if weekday < 0 {
				weekday = 6
			}
			prediction, err := getStreamProbability(day, month, weekday)
			if err != nil {
				log.Printf("%v: day: %v, month: %v, weekday %v", err, day, month, weekday)
				return
			}
			msg := fmt.Sprintf("@%v Prawdopodobieństwo streama dzisiaj wynosi %v%% gucci", message.User.DisplayName, int(prediction*100))
			client.Say(message.Channel, msg)
		}
	})
	client.Join("h2p_gucio")
	err := client.Connect()
	if err != nil {
		log.Fatal("error connecting to twitch: ", err)
	}
}
