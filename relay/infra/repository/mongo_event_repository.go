package repository

import (
	"context"

	"github.com/brunobdc/nostr/relay/infra/db"
	"github.com/brunobdc/nostr/relay/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoEventsRepository struct {
	collection *mongo.Collection
}

func MakeMongoEventsRepository() *MongoEventsRepository {
	return &MongoEventsRepository{
		collection: db.MongoRelayDB().Collection("Events"),
	}
}

func (repo *MongoEventsRepository) Save(ctx context.Context, event model.Event) error {
	_, err := repo.collection.InsertOne(ctx, event)
	return err
}

func (repo *MongoEventsRepository) SaveLatest(ctx context.Context, event model.Event) error {
	eventFound := model.NewEvent()
	err := repo.collection.FindOne(
		context.TODO(),
		bson.M{
			"PublicKey": event.PublicKey,
			"Kind":      event.Kind,
		},
	).Decode(&eventFound)
	if err == mongo.ErrNoDocuments {
		repo.collection.InsertOne(ctx, event)
		return nil
	}
	if err != nil {
		return err
	}

	if eventFound.CreatedAt > event.CreatedAt {
		return nil
	}
	repo.collection.InsertOne(ctx, event)
	repo.collection.DeleteOne(ctx, bson.M{"_id": eventFound.ID})

	return nil
}

func (repo *MongoEventsRepository) SaveParemeterizedLatest(ctx context.Context, event model.Event) error {
	var tagValue string
	if tagValues, ok := event.Tags["d"]; ok {
		tagValue = tagValues[0]
	}
	eventFound := model.NewEvent()
	err := repo.collection.FindOne(
		ctx,
		bson.M{
			"PublicKey": event.PublicKey,
			"Kind":      event.Kind,
			"Tags.d.0":  tagValue,
		},
	).Decode(eventFound)
	if err == mongo.ErrNoDocuments {
		repo.collection.InsertOne(ctx, event)
		return nil
	}
	if err != nil {
		return err
	}

	if eventFound.CreatedAt > event.CreatedAt {
		return nil
	}
	repo.collection.InsertOne(ctx, event)
	repo.collection.DeleteOne(ctx, bson.M{"_id": eventFound.ID})

	return nil
}

func (repo *MongoEventsRepository) FindWithFilters(
	ctx context.Context,
	filters []*model.Filters,
	foreachCb func(event *model.Event) error,
) error {
	mongoFilters := []bson.M{}
	for _, filter := range filters {
		f := bson.M{}
		if len(filter.IDs) > 0 {
			f["_id"] = bson.M{"$in": filter.IDs}
		}
		if len(filter.Authors) > 0 {
			f["PublicKey"] = bson.M{"$in": filter.Authors}
		}
		if len(filter.Kinds) > 0 {
			f["Kind"] = bson.M{"$in": filter.Kinds}
		}
		f["CreatedAt"] = bson.M{"$gte": filter.Since}
		if filter.Until > 0 {
			f["CreatedAt"] = bson.M{"$lte": filter.Until}
		}
		if len(filter.Tags) > 0 {
			for tag, values := range filter.Tags {
				f["Tags."+tag] = bson.M{"$all": values}
			}
		}
		mongoFilters = append(mongoFilters, f)
	}

	opts := options.Find().SetSort(bson.M{"CreatedAt": -1})
	if len(filters) > 0 && filters[0].Limit > 0 {
		opts = opts.SetLimit(int64(filters[0].Limit))
	}

	queryFilters := bson.M{}
	if len(mongoFilters) > 0 {
		queryFilters["$or"] = mongoFilters
	}
	cursor, err := repo.collection.Find(ctx, queryFilters, opts)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		event := model.NewEvent()
		cursor.Decode(event)
		err := foreachCb(event)
		if err != nil {
			return err
		}
	}

	return nil
}
