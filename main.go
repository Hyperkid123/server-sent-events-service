package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Hyperkid123/server-sent-events-service/src/conditions"
	"github.com/Hyperkid123/server-sent-events-service/src/kafkaconnector"
	"github.com/Hyperkid123/server-sent-events-service/src/mutators"
	"github.com/Hyperkid123/server-sent-events-service/src/sse"
	"github.com/Hyperkid123/server-sent-events-service/src/topics"
	"github.com/gobuffalo/packr"
	"github.com/joho/godotenv"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

func readTopics() map[string]topics.Topics {
	box := packr.NewBox("./static")
	file, foundError := box.FindString("topics.json")
	if foundError != nil {
		fmt.Println("Error while fetching file")
		file = os.Getenv("CONFIG_JSON")
	}
	data := make([]topics.Topics, 0)

	err := json.Unmarshal([]byte(file), &data)
	if err != nil {
		fmt.Println(err)
	}

	dataMap := make(map[string]topics.Topics)
	for i := 0; i < len(data); i++ {
		dataMap[data[i].Topic] = data[i]
	}

	return dataMap
}

func sendToListener(kafkaMessage *kafka.Message, topic topics.Topics) {
	if topic.Event == "" {
		topic.Event = "notification"
	}

	conditions := map[string](func(string, string) bool){
		"inventory": conditions.InventoryCondition,
		"approval":  conditions.ApprovalCondition,
	}

	mutators := map[string](func([]byte) []byte){
		"approval": mutators.ApprovalMutator,
	}

	go func() {
		for messageChannel, connectorInfo := range sse.MessageChannels {
			canSend := true
			responseMessage := kafkaMessage.Value
			for i := 0; i < len(topic.Conditions); i++ {
				canSend = conditions[topic.Conditions[i]](string(kafkaMessage.Value), connectorInfo.AccountNumber)
			}

			for i := 0; i < len(topic.Mutators); i++ {
				responseMessage = mutators[topic.Conditions[i]](responseMessage)
			}

			if canSend {
				msg := sse.FormatSSE(topic.Event, string(responseMessage))
				if topic.Room == "" {
					fmt.Println("No room, broadcasting!")
					messageChannel <- msg
				} else if connectorInfo.Room == topic.Room {
					fmt.Println("Sending to specific room")
					messageChannel <- msg
				} else {
					fmt.Println("Not sending", connectorInfo.Room, topic.Room)
				}
			}
		}
	}()

	fmt.Printf("%% Message on %s:\n%s\n", kafkaMessage.TopicPartition, string(kafkaMessage.Value))
}

func main() {
	godotenv.Load()
	topicsConfig := readTopics()
	apiVersion := os.Getenv("API_VERSION")
	appName := os.Getenv("APP_NAME")

	if apiVersion == "" {
		apiVersion = "v1"
	}

	if appName == "" {
		appName = "notifier"
	}

	go kafkaconnector.ConnectKafka(topicsConfig, sendToListener)

	http.HandleFunc(fmt.Sprintf("/api/%s/%s/connect", appName, apiVersion), sse.ListenHandler)
	http.HandleFunc(fmt.Sprintf("/api/%s/%s/lubdub", appName, apiVersion), func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "lubdub")
	})

	srv := &http.Server{
		Addr:         ":3000",
		Handler:      nil,
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
	}

	log.Fatal(srv.ListenAndServe())
}
