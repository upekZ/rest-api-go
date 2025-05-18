package services

import (
	"context"
	"fmt"
)

type Cache interface {
	IsValueTaken(key string, value string) bool
	SetValue(key string, value string, exists bool)
	DeleteField(key string, value string)
}

func (o *UserService) IsUniqueField(ctx context.Context, key string, value string) (bool, error) {

	found := o.cache.IsValueTaken(key, value)

	if found {
		return false, fmt.Errorf("duplicate %s: %s", key, value)
	} else {
		switch key {
		//To Do: Implement a mapping with unique field --> function ptr
		case uniqueFields["Phone"]:
			return o.db.IsPhoneUnique(ctx, value)
		case uniqueFields["Email"]:
			return o.db.IsEmailUnique(ctx, value)
		default:
			return true, nil
		}
	}
}
