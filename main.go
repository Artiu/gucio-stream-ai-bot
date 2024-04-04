package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc/v4"
	"github.com/joho/godotenv"
	"github.com/nicklaw5/helix"
)

const gucioUserId = "36954803"

func isGucioLive(client *helix.Client) bool {
	res, err := client.GetStreams(&helix.StreamsParams{UserIDs: []string{gucioUserId}, Type: "live"})
	if err != nil {
		log.Println("error getting streams: ", err)
		return false
	}
	if len(res.Data.Streams) == 0 {
		return false
	}
	return true
}

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
	twitchClientId := os.Getenv("TWITCH_CLIENT_ID")
	if twitchAccessToken == "" || twitchClientId == "" {
		log.Fatal("TWITCH_ACCESS_TOKEN and TWITCH_CLIENT_ID must be set")
	}

	helixClient, err := helix.NewClient(&helix.Options{
		ClientID:        twitchClientId,
		UserAccessToken: twitchAccessToken,
	})
	if err != nil {
		log.Fatal("error creating helix client: ", err)
	}

	cooldown := time.NewTimer(0)

	client := twitch.NewClient("guciostreamai", "oauth:"+twitchAccessToken)
	client.OnConnect(func() {
		log.Println("connected to twitch")
	})
	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		if strings.TrimSpace(message.Message) != "!czydzisstream" {
			return
		}
		select {
		case <-cooldown.C:
			cooldown.Reset(5 * time.Second)
		default:
			return
		}
		if isGucioLive(helixClient) {
			client.Say(message.Channel, fmt.Sprintf("@%v gucci jest live!", message.User.DisplayName))
			return
		}
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
	})
	client.Join("h2p_gucio")
	err = client.Connect()
	if err != nil {
		log.Fatal("error connecting to twitch: ", err)
	}
}
