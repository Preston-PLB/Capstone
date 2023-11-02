package db

import (
	"context"

	"git.preston-baxter.com/Preston_PLB/capstone/frontend-service/config"
	"git.preston-baxter.com/Preston_PLB/capstone/frontend-service/db/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (db *DB) FindUserByEmail(email string) (*models.User, error) {
	conf := config.Config()

	opts := options.FindOne()
	res := db.client.Database(conf.Mongo.EntDb).Collection(conf.Mongo.EntCol).FindOne(context.Background(), bson.M{"email": email, "obj_info.ent": models.USER_TYPE}, opts)

	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, res.Err()
	}

	user := &models.User{}
	err := res.Decode(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (db *DB) FindUserById(id string) (*models.User, error) {
	conf := config.Config()

	opts := options.FindOne()
	res := db.client.Database(conf.Mongo.EntDb).Collection(conf.Mongo.EntCol).FindOne(context.Background(), bson.M{"_id": id, "obj_info.ent": models.USER_TYPE}, opts)

	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, res.Err()
	}

	user := &models.User{}
	err := res.Decode(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (db *DB) FindAllUsers() ([]models.User, error) {
	conf := config.Config()

	opts := options.Find()
	res, err := db.client.Database(conf.Mongo.EntDb).Collection(conf.Mongo.EntCol).Find(context.Background(), bson.M{"obj_info.ent": models.USER_TYPE}, opts)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, res.Err()
	}

	users := []models.User{}
	err = res.All(context.Background(), users)
	if err != nil {
		return nil, err
	}

	return users, nil
}
