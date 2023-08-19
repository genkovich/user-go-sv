package cache

import (
	"github.com/mediocregopher/radix/v3"
)

type Redis struct {
	connection *radix.Pool
}

type Cache interface {
	Set(key string, value string, seconds string) error
	Get(key string) (string, error)
}

func NewCache(address string) (*Redis, error) {
	conn, err := radix.NewPool("tcp", address, 10)
	if err != nil {
		return nil, err
	}

	return &Redis{connection: conn}, nil
}

func (r *Redis) Set(key string, value string, seconds string) error {
	return r.connection.Do(radix.Cmd(nil, "SETEX", key, seconds, value))
}

func (r *Redis) Get(key string) (string, error) {
	var value string

	err := r.connection.Do(radix.Cmd(&value, "GET", key))
	if err != nil {
		return "", err
	}

	return value, nil
}
