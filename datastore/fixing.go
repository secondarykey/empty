package datastore

import (
	"context"

	"cloud.google.com/go/datastore"
)

type OAuthValue struct {
	ClientID     string
	ClientSecret string
}

const OAuthValueKindName = "OAuthValue"

func GetOAuthValue(ctx context.Context) (*OAuthValue, error) {
	client, err := datastore.NewClient(ctx, ProjectId)
	if err != nil {
		return nil, err
	}
	k := datastore.NameKey(OAuthValueKindName, fixingKey, nil)
	e := new(OAuthValue)
	if err := client.Get(ctx, k, e); err != nil {
		if err == datastore.ErrNoSuchEntity {
			e.ClientID = ""
			e.ClientSecret = ""
			if _, pe := client.Put(ctx, k, e); pe != nil {
				return nil, pe
			}
			return nil, fmt.Errorf("Put Empty Value(OAuthValue)")
		}
		return nil, err
	}
	return e, nil
}
