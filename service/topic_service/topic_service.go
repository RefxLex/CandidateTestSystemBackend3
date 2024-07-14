package topic_service

import (
	"CandidateTestSystemBackend3/types"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func GetTopicById(topicID int) (*types.Topic, error) {
	db, err := sqlx.Connect("postgres", "user=user dbname=candidatetestsystem sslmode=disable password=12345")
	if err != nil {
		log.Fatalln(err)
	}

	topic := types.Topic{}
	err = db.Get(&topic, "SELECT * FROM topics WHERE id=$1", topicID)
	return &topic, err
}

func GetAllTopics() ([]types.Topic, error) {
	db, err := sqlx.Connect("postgres", "user=user dbname=candidatetestsystem sslmode=disable password=12345")
	if err != nil {
		log.Fatalln(err)
	}

	topics := []types.Topic{}
	db.Select(&topics, "SELECT * FROM topics ORDER BY name ASC")
	return topics, err
}

func CreateTopic(topic *types.Topic) (*types.Topic, error) {
	db, err := sqlx.Connect("postgres", "user=user dbname=candidatetestsystem sslmode=disable password=12345")
	if err != nil {
		log.Fatalln(err)
	}

	createdTopic := types.Topic{}
	rows, err := db.Queryx(`INSERT INTO topics (name) VALUES ($1) RETURNING *`, topic.Name)
	for rows.Next() {
		err := rows.StructScan(&createdTopic)
		if err != nil {
			log.Fatalln(err)
		}
	}
	return &createdTopic, err
}

func UpdateTopic(topicID int, topic *types.Topic) (*types.Topic, error) {
	db, err := sqlx.Connect("postgres", "user=user dbname=candidatetestsystem sslmode=disable password=12345")
	if err != nil {
		log.Fatalln(err)
	}

	newTopic := types.Topic{}
	rows, err := db.Queryx("UPDATE topics SET name = $1 WHERE id = $2 RETURNING *", topic.Name, topicID)
	for rows.Next() {
		err := rows.StructScan(&newTopic)
		if err != nil {
			log.Fatalln(err)
		}
	}
	return &newTopic, err
}

func DeleteTopic(topicID int) error {
	db, err := sqlx.Connect("postgres", "user=user dbname=candidatetestsystem sslmode=disable password=12345")
	if err != nil {
		log.Fatalln(err)
	}

	_, err = db.Exec("DELETE FROM topics WHERE id=$1", topicID)
	return err
}
