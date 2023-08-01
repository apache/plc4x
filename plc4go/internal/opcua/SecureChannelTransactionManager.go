/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package opcua

import (
	"sync"
	"sync/atomic"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

//go:generate go run ../../tools/plc4xgenerator/gen.go -type=SecureChannelTransactionManager
type SecureChannelTransactionManager struct {
	transactionIdentifierGenerator atomic.Int32
	requestIdentifierGenerator     atomic.Int32
	activeTransactionId            atomic.Int32
	queue                          map[int32]SecureChannelTransactionManagerTransaction

	lock sync.Mutex

	log zerolog.Logger `ignore:"true"`
}

func NewSecureChannelTransactionManager(log zerolog.Logger) *SecureChannelTransactionManager {
	return &SecureChannelTransactionManager{
		queue: map[int32]SecureChannelTransactionManagerTransaction{},
		log:   log,
	}
}

func (m *SecureChannelTransactionManager) submit(onSend func(transactionId int32), transactionId int32) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.log.Info().Int32("activeTransactionId", m.activeTransactionId.Load()).Msg("Active transaction Number")
	if m.activeTransactionId.Load() == transactionId {
		onSend(transactionId)
		newTransactionId := m.getActiveTransactionIdentifier()
		if len(m.queue) > 0 {
			t, ok := m.queue[newTransactionId]
			if !ok {
				m.log.Info().Int("queueLength", len(m.queue)).Msg("Length of Queue")
				m.log.Info().Int32("newTransactionId", newTransactionId).Msg("Transaction ID")
				m.log.Info().Interface("map", m.queue).Msg("Map is")
				return errors.Errorf("Transaction Id not found in queued messages %v", m.queue)
			}
			delete(m.queue, newTransactionId)
			if err := m.submit(t.consumer, t.transactionId); err != nil {
				return errors.Wrap(err, "Error submitting")
			}
		}
	} else {
		m.log.Info().Int32("transactionId", transactionId).Msg("Storing out of order transaction")
		m.queue[transactionId] = SecureChannelTransactionManagerTransaction{consumer: onSend, transactionId: transactionId}
	}
	return nil
}

func (m *SecureChannelTransactionManager) getTransactionIdentifier() int32 {
	return m.transactionIdentifierGenerator.Add(1) - 1
}

func (m *SecureChannelTransactionManager) getActiveTransactionIdentifier() int32 {
	return m.activeTransactionId.Add(1)
}

type SecureChannelTransactionManagerTransaction struct {
	transactionId int32
	consumer      func(transactionId int32)
}
