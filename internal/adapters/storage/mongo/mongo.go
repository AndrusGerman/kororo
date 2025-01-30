package mongodb

import (
	"context"
	"time"

	"kororo/internal/adapters/config"
	"kororo/internal/core/domain/types"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Mongo struct {
	client *mongo.Client
	config *config.Config
}

func (ctx *Mongo) Close() error {
	return ctx.client.Disconnect(context.TODO())
}

func (ctx *Mongo) initConnection(uri string) error {
	if ctx.client != nil {
		return nil
	}

	var appName = "KororoServices"
	var timeout = 10 * time.Second

	client, err := mongo.Connect(options.Client().ApplyURI(uri),
		&options.ClientOptions{
			AppName: &appName,
			//TLSConfig: &tls.Config{
			//	InsecureSkipVerify: true,
			//},
			Timeout: &timeout,
		},
	)

	if err != nil {
		return err
	}
	ctx.client = client
	return nil
}

func (ctx *Mongo) GetDB(database types.Database) *mongo.Database {
	return ctx.client.Database(database.String())
}

func (ctx *Mongo) Collection(name string) *mongo.Collection {
	return ctx.GetDB(ctx.config.Database()).Collection(name)
}

func (m *Mongo) ListCollectionNames() ([]string, error) {
	return m.GetDB(m.config.Database()).ListCollectionNames(context.TODO(), bson.D{})
}

func NewMongo(config *config.Config) (*Mongo, error) {
	var mongo = new(Mongo)
	mongo.config = config
	var err = mongo.initConnection(config.UriMongo())
	return mongo, err
}
