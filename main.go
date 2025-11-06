package main

import (
	"fmt"
	datastructures "go/rdtbot/dataStructures"
	"go/rdtbot/database"
	discordclient "go/rdtbot/discordClient"
	redditclient "go/rdtbot/redditClient"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading dotenv:", err)
	}

	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	clientId := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("SECRET")
	discordAuth := os.Getenv("AUTHORIZATION")

	if username == "" || password == "" || clientId == "" || clientSecret == "" {
		log.Fatal("You must set all env vars before running the program")
	}

	history, err := database.NewHistory()

	if err != nil {
		log.Fatal(err)
	}

	defer history.Close()

	client := redditclient.NewRedditClient(username, password, clientId, clientSecret)
	dclient := discordclient.New(discordAuth, "895480798709428255")
	q := datastructures.Queue[string]{}
	arcQ := datastructures.NewArc(q)

	go func() {
		subreddits := []string{"eu_nvr", "animebrasil"}

		for {
			for _, i := range subreddits {
				resp, err := client.GetSubredditPosts(i, "10")
				if err != nil {
					log.Printf("Error getting posts: %v\n", err)
					continue
				}

				itens, err := redditclient.NewListing(resp)
				if err != nil {
					log.Printf("Error creating listing: %v\n", err)
					continue
				}

				h, err := history.Get()
				if err != nil {
					log.Printf("Error querying the database: %v\n", err)
					continue
				}

				for _, i := range itens.GetLinks() {
					if !h.Contains(i) {
						copy := arcQ.Get()
						copy.Push(i) // updating the queue copy to avoid data races
						arcQ.Set(copy)
						if err = history.Insert(i); err != nil {
							log.Println("History error:", err)
						}
					}
				}
			}

			time.Sleep(10 * time.Minute)
		}
	}()

	go func() {
		time.Sleep(15 * time.Second)

		for {
			copy := arcQ.Get()
			next := copy.Pop()

			if next != "" {
				log.Println("Sending link:", next)
				if err := dclient.SendMsg(next); err != nil {
					log.Printf("Error sending content: %v\n", err)
				}
			}

			arcQ.Set(copy)
			time.Sleep(time.Minute * 3)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop

	fmt.Println("\nQuiting...")
}
