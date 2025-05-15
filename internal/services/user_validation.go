package services

import (
	"context"
	"fmt"
)

type Cache interface {
	IsValueTaken(key string, value string) (bool, bool)
	SetValue(key string, value string, exists bool)
	DeleteField(key string, value string)
}

func (o *UserService) IsUniqueField(ctx context.Context, key string, value string) (bool, error) {

	taken, found := o.cache.IsValueTaken(key, value)

	switch found {
	case true:
		switch taken {
		case true:
			return false, fmt.Errorf("%s: %s already taken", key, value)
		case false:
			return true, nil
		}
	case false:
		switch key {
		case uniqueFields["Phone"]:
			return o.db.IsPhoneTaken(ctx, value)
		case uniqueFields["Email"]:
			return o.db.IsEmailTaken(ctx, value)
		default:
			return false, fmt.Errorf("%s: invalid key", key)
		}
	}

	return false, fmt.Errorf("validation failure for %s: %s ", key, value)
}
