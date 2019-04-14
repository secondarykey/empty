package datastore

import (
	"context"

	"cloud.google.com/go/datastore"
	"golang.org/x/xerrors"
)

type OAuthValue struct {
	ClientID     string
	ClientSecret string
}

const OAuthValueKindName = "OAuthValue"

func GetOAuthValue(ctx context.Context) (*OAuthValue, error) {

	client, err := datastore.NewClient(ctx, "")
	if err != nil {
		return nil, xerrors.Errorf("error in datastore.NewClient(): %w", err)
	}

	k := datastore.NameKey(OAuthValueKindName, fixingKey, nil)
	e := new(OAuthValue)
	if err := client.Get(ctx, k, e); err != nil {
		if err == datastore.ErrNoSuchEntity {
			e.ClientID = ""
			e.ClientSecret = ""
			if _, pe := client.Put(ctx, k, e); pe != nil {
				return nil, xerrors.Errorf("error in client.Put(): %w", err)
			}
			return nil, xerrors.Errorf("error Put() successed,but you need set the OAuthValue.ClientID")
		}
		return nil, err
	}

	if e.ClientID == "" {
		return nil, xerrors.Errorf("error OAuthValue.ClientID is empty.")
	}

	return e, nil
}
