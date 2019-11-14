package main

import (
	"context"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go"
	"google.golang.org/api/option"
)

func getFireStore() (*firestore.Client, context.Context, error) {

	ctx := context.Background()
	accountkey := option.WithCredentialsFile("./mailinggo-firebase-adminsdk-55h8l-27d95cb552.json")
	app, err := firebase.NewApp(ctx, nil, accountkey)
	if err != nil {
		return nil, nil, err
	}
	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, nil, err
	}
	return client, ctx, nil
}
