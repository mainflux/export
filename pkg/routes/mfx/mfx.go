// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package mfx

import (
	"github.com/gogo/protobuf/proto"
	"github.com/mainflux/export/pkg/messages"
	"github.com/mainflux/export/pkg/routes"
	"github.com/mainflux/mainflux"
	"github.com/mainflux/mainflux/logger"
	"github.com/mainflux/senml"
	nats "github.com/nats-io/nats.go"
)

const (
	// ContentTypeJSON represents SenML in JSON format content type.
	ContentTypeJSON = "application/senml+json"

	// ContentTypeCBOR represents SenML in CBOR format content type.
	ContentTypeCBOR = "application/senml+cbor"
)

var formats = map[string]senml.Format{
	ContentTypeJSON: senml.JSON,
	ContentTypeCBOR: senml.CBOR,
}

type mfxRoute struct {
	route routes.Route
}

func NewRoute(n, m, s string, w int, log logger.Logger, pub messages.Publisher) routes.Route {
	return &mfxRoute{
		route: routes.NewRoute(n, m, s, w, log, pub),
	}
}

func (mr *mfxRoute) Workers() int {
	return mr.route.Workers()
}

func (mr *mfxRoute) NatsTopic() string {
	return mr.route.NatsTopic()
}

func (mr *mfxRoute) MqttTopic() string {
	return mr.route.MqttTopic()
}

func (mr *mfxRoute) Subtopic() string {
	return mr.route.Subtopic()
}

func (mr *mfxRoute) Consume() {
	mr.route.Consume()
}

func (mr *mfxRoute) Process(data []byte) ([]byte, error) {
	msg := mainflux.Message{}
	err := proto.Unmarshal(data, &msg)
	if err != nil {
		return []byte{}, err
	}
	format, ok := formats[msg.ContentType]
	if !ok {
		format = senml.JSON
	}

	raw, err := senml.Decode(msg.Payload, format)
	if err != nil {
		return []byte{}, err
	}

	payload, err := senml.Encode(raw, senml.JSON)
	if err != nil {
		return []byte{}, err
	}
	return payload, nil
}

func (mr *mfxRoute) MessagesBuffer() chan *nats.Msg {
	return mr.route.MessagesBuffer()
}