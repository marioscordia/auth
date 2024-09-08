package redisgo

import (
	"context"
	"encoding/json"
	"log"

	"github.com/gomodule/redigo/redis"
	"github.com/marioscordia/auth/internal/model"
)

func NewUserCache(conn redis.Conn) *userCache {
	return &userCache{
		conn: conn,
	}

}

type userCache struct {
	conn redis.Conn
}

func (u *userCache) Save(ctx context.Context, user *model.User) {
	v, err := json.Marshal(user)
	if err != nil {
		log.Printf("error saving user in cache: %v", err)
	} else {
		if _, err = u.conn.Do("SET", user.ID, v); err != nil {
			log.Printf("error saving user in cache: %v", err)
		}
	}
}

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

func (u *userCache) Delete(ctx context.Context, userId int64) {
	if _, err := u.conn.Do("DEL", userId); err != nil {
		log.Printf("error deleting user from cache: %v", err)
	}
}

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
