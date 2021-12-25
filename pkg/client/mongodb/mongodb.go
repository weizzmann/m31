package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//func NewClient(ctx context.Context, host, port, username, password, database, authDB string) (db *mongo.Database, err error) {
//	var mongoDBURL string
//	var isAuth bool
//	if username == "" && password == "" {
//		mongoDBURL = fmt.Sprintf("mongodb://%s:%s", host, port)
//	} else {
//		isAuth = true
//		mongoDBURL = fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, host, port)
//	}
//
//	clientOptions := options.Client().ApplyURI(mongoDBURL)
//	if isAuth {
//		if authDB == "" {
//			authDB = database
//		}
//		clientOptions.SetAuth(options.Credential{
//			AuthSource: authDB,
//			Username:   username,
//			Password:   password,
//		})
//	}
//
//	client, err := mongo.Connect(ctx, clientOptions)
//	if err != nil {
//		return nil, fmt.Errorf("failed to connect to mongoDB due to error: %v", err)
//	}
//
//	if err = client.Ping(ctx, nil); err != nil {
//		return nil, fmt.Errorf("failed to ping mongoDB due to error: %v", err)
//	}
//
//	return client.Database(database), nil
//}

func NewClientCloud(ctx context.Context) (db *mongo.Database, err error) {
	clientOptions := options.Client().
		ApplyURI("mongodb+srv://alex:a4LZAoI4MGtTUh3e@cluster0.vo553.mongodb.net/module31?retryWrites=true&w=majority")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return
	}
	return client.Database("module31"), nil
}
