// Copyright (c) 2020-2021 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/riversgo007/EvaBot/mtproto/blob/master/LICENSE for details

package mtproto

import (
	"reflect"

	"github.com/pkg/errors"
	"github.com/riversgo007/EvaBot/mtproto/internal/encoding/tl"
	"github.com/riversgo007/EvaBot/mtproto/internal/session"
	"github.com/riversgo007/EvaBot/mtproto/internal/utils"
)

// helper methods

// GetSessionID returns the current session id 🧐
func (m *MTProto) GetSessionID() int64 {
	return m.sessionId
}

// GetSeqNo returns seqno 🧐
func (m *MTProto) GetSeqNo() int32 {
	return m.seqNo
}

// GetServerSalt returns current server salt 🧐
func (m *MTProto) GetServerSalt() int64 {
	return m.serverSalt
}

// GetAuthKey returns decryption key of current session salt 🧐
func (m *MTProto) GetAuthKey() []byte {
	return m.authKey
}

func (m *MTProto) SetAuthKey(key []byte) {
	m.authKey = key
	m.authKeyHash = utils.AuthKeyHash(m.authKey)
}

func (m *MTProto) MakeRequest(msg tl.Object) (any, error) {
	return m.makeRequest(msg)
}

func (m *MTProto) MakeRequestWithHintToDecoder(msg tl.Object, expectedTypes ...reflect.Type) (any, error) {
	if len(expectedTypes) == 0 {
		return nil, errors.New("expected a few hints. If you don't need it, use m.MakeRequest")
	}
	return m.makeRequest(msg, expectedTypes...)
}

func (m *MTProto) AddCustomServerRequestHandler(handler customHandlerFunc) {
	m.serverRequestHandlers = append(m.serverRequestHandlers, handler)
}

func (m *MTProto) warnError(err error) {
	if m.Warnings != nil && err != nil {
		m.Warnings <- err
	}
}

func (m *MTProto) SaveSession() (err error) {
	s := new(session.Session)
	s.Key = m.authKey
	s.Hash = m.authKeyHash
	s.Salt = m.serverSalt
	s.Hostname = m.addr
	return session.SaveSession(s, m.tokensStorage)
}

func (m *MTProto) LoadSession(s *session.Session) {
	m.authKey = s.Key
	m.authKeyHash = s.Hash
	m.serverSalt = s.Salt
	m.addr = s.Hostname
}