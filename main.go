package main

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/gommon/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"math/rand/v2"
	"os"
	"path/filepath"
)

type JSON json.RawMessage

// Scan scan value into JSON, implements sql.Scanner interface
func (j *JSON) Scan(value interface{}) error {
	bytes, ok := value.(string)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSON value:", value))
	}

	result := json.RawMessage{}
	err := json.Unmarshal([]byte(bytes), &result)
	*j = JSON(result)
	return err
}

// Value return json value, implement driver.Valuer interface
func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.RawMessage(j).MarshalJSON()
}

// GormDataType is the fallback of gorm if there's no annotation tag `gorm:"type:jsonb"`
func (JSON) GormDataType() string {
	return "json"
}

// GormDBDataType is for Migration
func (JSON) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return "JSON"
}

type Message struct {
	ID   uint `gorm:"primaryKey"`
	Data JSON `gorm:"type:json"`
}

func (Message) TableName() string {
	return "messages_table"
}

type MessageData struct {
	Worker            int    `json:"worker"`
	Message           string `json:"message"`
	Interval          int    `json:"interval"`
	DestinationWorker int    `json:"destination_worker"`
}

func main() {
	var err error
	var dir string

	// Development-specific logic
	dir = filepath.Join(".", "db")

	err = os.MkdirAll(dir, 0755)
	if err != nil {
		log.Fatal(err)
		return
	}

	db, err := gorm.Open(sqlite.Open(filepath.Join(dir, "workers.db")), &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		log.Fatal("Error starting GORM: ", err)
		return
	}

	defer func() {
		sqlDbCloser, err := db.DB()
		if err != nil {
			log.Fatal("Error preparing the database to close: ", err)
			return
		}
		err = sqlDbCloser.Close()
		if err != nil {
			log.Fatal("Error closing the database: ", err)
			return
		}
	}()

	err = db.AutoMigrate(&Message{})
	if err != nil {
		log.Fatal("Error creating tables with AutoMigrate: ", err)
		return
	}

	//rand.Seed(uint64(time.Now().UnixNano()))
	seed := 1729727730264126872

	var seedAs32bytes [32]byte
	for i := 0; i < 8; i++ {
		seedAs32bytes[i] = byte(seed >> (8 * i))
	}
	rand.New(rand.NewChaCha8(seedAs32bytes))

	for i := 0; i < 10; i++ {
		baseInterval := (rand.IntN(3) + 1) * i

		// Select 3 distinct workers
		workers := rand.Perm(5)[:3]
		workerA := workers[0] + 1 // +1 because rand.Perm starts at 0
		workerB := workers[1] + 1
		workerC := workers[2] + 1

		// Generate random time intervals (to be displayed in milliseconds)
		intervals := make([]int, 6)
		for j := 0; j < 6; j++ {
			intervals[j] = rand.IntN(2) + 1 // Interval between 1 and 2 milliseconds
		}

		// Messages to be stored
		messages := []MessageData{
			{
				Worker:            workerA,
				Message:           "Olá, como vai?",
				Interval:          (intervals[0]) + baseInterval,
				DestinationWorker: workerB,
			},
			{
				Worker:            workerB,
				Message:           fmt.Sprintf("Vou bem, e você? Pode falar com o worker %d?", workerC),
				Interval:          (intervals[1] + intervals[0]) + baseInterval,
				DestinationWorker: workerA,
			},
			{
				Worker:            workerA,
				Message:           fmt.Sprintf("Olá, como vai? A worker %d pediu para falar contigo.", workerB),
				Interval:          (intervals[2] + intervals[1] + intervals[0]) + baseInterval,
				DestinationWorker: workerC,
			},
			{
				Worker:            workerA,
				Message:           "De novo! Vou bem, obrigada!",
				Interval:          (intervals[3] + intervals[1] + intervals[0]) + baseInterval,
				DestinationWorker: workerB,
			},
			{
				Worker:            workerC,
				Message:           "Vou bem, obrigado! O que ela quer?",
				Interval:          (intervals[4] + intervals[2] + intervals[1] + intervals[0]) + baseInterval,
				DestinationWorker: workerA,
			},
			{
				Worker:            workerA,
				Message:           "Não sei. Pergunte para ela, meu chapa.",
				Interval:          (intervals[5] + intervals[4] + intervals[2] + intervals[1] + intervals[0]) + baseInterval,
				DestinationWorker: workerC,
			},
		}

		// Stores messages in JSON1 column
		for j, msg := range messages {
			jsonData, err := json.Marshal(msg)
			if err != nil {
				log.Errorf("Error marshalling a message (%d): %v", i*6+j, err)
				continue
			}

			// Create messages
			err = db.Exec(`
				INSERT INTO messages_table (data) VALUES (json(?));`, string(jsonData)).
				Error
			if err != nil {
				log.Errorf("Error inserting a message (%d): %v", i*6+j, err)
				continue
			}
		}
	}

	//// Extra tip if you are going to use Go
	//// Implement the Scanner and Valuer functions shown above the main,
	////  as well as the GormDataType and the GormDBDataType
	//
	//// Below it is not all of what you are going to do, but here is a tip for how to query json columns
	//// Then, inside the for loop, you should use a json deserialization library to parse the JSON column,
	////  such as using json.Unmarshal
	//
	//// Example of Raw Querying but with the correct table:
	//var unparsedJSON []JSON
	//result := db.Raw(`
	//	SELECT data FROM messages_table;`)
	//err = result.Error
	//if err != nil {
	//	log.Error("Error getting messages in main: ", err)
	//	return
	//}
	//
	//result = result.Scan(&unparsedJSON)
	//err = result.Error
	//if err != nil {
	//	log.Error("Error scanning messages in main: ", err)
	//	return
	//}
	//
	//if len(unparsedJSON) == 0 {
	//	log.Error("No messages found.")
	//	return
	//}
	//
	//for i, msg := range unparsedJSON {
	//	...
	//}
}
