package topic_service

import (
	"CandidateTestSystemBackend3/types"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func GetTopicById(topicID int) (*types.Topic, error) {
	//db, err := sqlx.Open("postgres", "user:12345@(127.0.0.1:5432)/candidatetestsystem?parseTime=true?sslmode=disable")
	db, err := sqlx.Connect("postgres", "user=user dbname=candidatetestsystem sslmode=disable password=12345")
	if err != nil {
		log.Fatalln(err)
	}
	topic := types.Topic{}
	err = db.Get(&topic, "SELECT * FROM topics WHERE id=$1", topicID)
	return &topic, err
}

func GetAllTopics()
