// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package messages

import (
	"github.com/go-redis/redis"
	"github.com/mainflux/mainflux/errors"
)

const (
	msgPrefix = "prefix"
	streamLen = 1000
)

var _ Cache = (*cache)(nil)

var (
	errGrpCreateGroupMissing  = errors.New("Group not created, group not being set")
	errGrpCreateStreamMissing = errors.New("Group not created, stream not being set")
)

type cache struct {
	client *redis.Client
}

// NewMessageCache returns redis message cache implementation.
func NewRedisCache(client *redis.Client) Cache {
	return &cache{client: client}
}

func (cc *cache) Add(stream, topic string, payload []byte) (string, error) {
	m := msg{topic: topic, payload: string(payload)}.encode()
	return cc.add(stream, m)
}

func (cc *cache) Remove(msgID string) error {
	return cc.client.Del(msgID).Err()
}

func (cc *cache) GroupCreate(stream, group string) (string, error) {
	if stream == "" {
		return "", errGrpCreateStreamMissing
	}
	if group == "" {
		return "", errGrpCreateGroupMissing
	}
	return cc.client.XGroupCreateMkStream(stream, group, "$").Result()
}

func (cc *cache) add(stream string, m map[string]interface{}) (string, error) {

	record := &redis.XAddArgs{
		Stream:       stream,
		MaxLenApprox: streamLen,
		Values:       m,
	}

	return cc.client.XAdd(record).Result()
}

func (cc *cache) ReadGroup(streams []string, group, consumer string) ([]redis.XStream, error) {

	xReadGroupArgs := &redis.XReadGroupArgs{
		Group:    group,
		Consumer: consumer,
		Streams:  streams,
		Block:    0,
	}
	return cc.client.XReadGroup(xReadGroupArgs).Result() //Get Results from XRead command
}