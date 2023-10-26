package zmisc

import (
	"context"

	"github.com/pkg/errors"
	"go.elastic.co/apm/module/apmmongo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoClient(mongoDbHost string) (*mongo.Client, error) {
	ctx := context.Background()

	opt := options.Client()
	opt.ApplyURI(mongoDbHost)
	opt.SetMonitor(apmmongo.CommandMonitor())

	mongo, err := mongo.Connect(ctx, opt)
	if err != nil {
		return nil, errors.Wrap(err, "error creating client for mongodb")
	}

	return mongo, nil
}
