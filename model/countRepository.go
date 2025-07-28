package model

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Counter struct {
	Id         string
	SequenceNo uint32
}

// if you define any new table(collection) then please set an initial value in db mnually
const (
	COUNTER_ID_USER_ID = "userId"
)

func GetNextSequence(ctx context.Context, db *mongo.Database, counterName string) (int64, error) {
    filter := bson.M{"_id": counterName}
    update := bson.M{"$inc": bson.M{"seq": 1}}
    opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

    var result struct {
        Seq int64 `bson:"seq"`
    }

    err := db.Collection("counters").FindOneAndUpdate(ctx, filter, update, opts).Decode(&result)
    if err != nil {
        return 0, err
    }

    return result.Seq, nil
}