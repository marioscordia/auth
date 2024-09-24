package redisgo

import (
	"context"
	"encoding/json"
	"log"

	"github.com/gomodule/redigo/redis"
	"github.com/marioscordia/auth/internal/client/cache"
	"github.com/marioscordia/auth/internal/model"
)

// NewUserCache returns a new Cache interface
func NewUserCache(conn redis.Conn) cache.Cache {
	return &userCache{
		conn: conn,
	}

}

type userCache struct {
	conn redis.Conn
}

// Save is a method to save user information in cache
func (u *userCache) Save(ctx context.Context, user *model.User) {
	v, err := json.Marshal(user)
	if err != nil {
		log.Printf("error saving user in cache: %v", err)
		return
	}

	if _, err = u.conn.Do("SET", user.ID, v); err != nil {
		log.Printf("error saving user in cache: %v", err)
	}
}

// Get is a method to get user informartion from cache
func (u *userCache) Get(ctx context.Context, userId int64) *model.User {
	v, err := redis.Bytes(u.conn.Do("GET", userId))
	if err != nil {
		log.Printf("error getting user from cache: %v", err)
		return nil
	}

	user := new(model.User)

	if err = json.Unmarshal(v, user); err != nil {
		log.Printf("error getting user from cache: %v", err)
		return nil
	}

	return user
}

// Delete is a method to delete user informartion from cache
func (u *userCache) Delete(ctx context.Context, userId int64) {
	if _, err := u.conn.Do("DEL", userId); err != nil {
		log.Printf("error deleting user from cache: %v", err)
	}
}

// Update is a method to update user informartion in cache
func (u *userCache) Update(ctx context.Context, update *model.UserUpdate) {
	user := u.Get(ctx, update.ID)

	if user != nil {
		if update.Name != "" {
			user.Name = update.Name
		}

		if update.Email != "" {
			user.Email = update.Email
		}

		u.Save(ctx, user)
	}
}
