package datastore

import (
	"context"
	"fmt"

	"cloud.google.com/go/datastore"
)

var ProjectId string

type Entity struct {
	Value string
}

func Put(v string) error {
	ctx := context.Background()

	dsClient, err := datastore.NewClient(ctx, ProjectId)
	if err != nil {
		return err
	}

	k := datastore.NameKey("Entity", "stringID", nil)
	e := new(Entity)
	if err := dsClient.Get(ctx, k, e); err != nil {
		//return err
	}

	old := e.Value
	e.Value = v

	if _, err := dsClient.Put(ctx, k, e); err != nil {
		return err
	}

	fmt.Printf("Updated value from %q to %q\n", old, e.Value)
	return nil
}
