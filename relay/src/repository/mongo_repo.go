package repository

import (
	"context"

	"github.com/brunobdc/nostr/relay/src/infra"
	"github.com/brunobdc/nostr/relay/src/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoEventsRepository struct {
	collection *mongo.Collection
}

func NewMongoEventsRepository() *MongoEventsRepository {
	return &MongoEventsRepository{collection: infra.DB.Collection("Events")}
}

type MongoEventCursor struct {
	cursor *mongo.Cursor
}

func (c *MongoEventCursor) Close(ctx context.Context) {
	c.cursor.Close(ctx)
}

func (c *MongoEventCursor) Next(ctx context.Context) bool {
	return c.cursor.Next(ctx)
}

func (c *MongoEventCursor) Decode(event *model.Event) error {
	return c.cursor.Decode(event)
}

func (repo *MongoEventsRepository) Save(event model.Event) {
	repo.collection.InsertOne(context.TODO(), event)
}

func (repo *MongoEventsRepository) SaveLatest(event model.Event) {
	var eventFound model.Event
	err := repo.collection.FindOne(
		context.TODO(),
		bson.M{
			"PublicKey": event.PublicKey,
			"Kind":      event.Kind,
		},
	).Decode(&eventFound)
	if err == mongo.ErrNoDocuments {
		repo.collection.InsertOne(context.TODO(), event)
		return
	}
	if err != nil {
		panic(err)
	}

	if eventFound.CreatedAt > event.CreatedAt {
		return
	}
	repo.collection.InsertOne(context.TODO(), event)
	repo.collection.DeleteOne(context.TODO(), bson.M{"_id": eventFound.ID})
}

func (repo *MongoEventsRepository) SaveParemeterizedLatest(event model.Event) {
	var tagValue string
	if tagValues, ok := event.Tags["d"]; ok {
		tagValue = tagValues[0]
	}
	var eventFound model.Event
	err := repo.collection.FindOne(
		context.TODO(),
		bson.M{
			"PublicKey": event.PublicKey,
			"Kind":      event.Kind,
			"Tags.d.0":  tagValue,
		},
	).Decode(&eventFound)
	if err == mongo.ErrNoDocuments {
		repo.collection.InsertOne(context.TODO(), event)
		return
	}
	if err != nil {
		panic(err)
	}

	if eventFound.CreatedAt > event.CreatedAt {
		return
	}
	repo.collection.InsertOne(context.TODO(), event)
	repo.collection.DeleteOne(context.TODO(), bson.M{"_id": eventFound.ID})
}

func (repo *MongoEventsRepository) FindWithFilters(filters []*model.Filters) EventCursor {
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

	opts := options.Find()
	if len(filters) > 0 && filters[0].Limit > 0 {
		opts = opts.SetLimit(int64(filters[0].Limit))
	}

	cursor, err := repo.collection.Find(context.TODO(), bson.M{"$or": mongoFilters}, opts)
	if err != nil {
		panic(err)
	}

	return &MongoEventCursor{cursor: cursor}
}
