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

package test_utilities

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/tests"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/tests/state_machine"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/tests/time_machine"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/tests/trapped_classes"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type TPDU struct {
	x    []byte
	a, b int
}

var _ PDU = TPDU{}

func (t TPDU) GetLeafName() string {
	return "TPDU"
}

func (t TPDU) X() []byte {
	return t.x
}

func (t TPDU) A() int {
	return t.a
}

func (t TPDU) B() int {
	return t.b
}

func (t TPDU) String() string {
	content := ""
	if t.x != nil {
		content += fmt.Sprintf(" x=%v", t.x)
	}
	if t.a != 0 {
		content += fmt.Sprintf(" a=%v", t.a)
	}
	if t.b != 0 {
		content += fmt.Sprintf(" b=%v", t.b)
	}
	return fmt.Sprintf("<TPDU%v>", content)
}

func (t TPDU) SetRootMessage(message spi.Message) {
	//TODO implement me
	panic("implement me")
}

func (t TPDU) GetRootMessage() spi.Message {
	//TODO implement me
	panic("implement me")
}

func (t TPDU) SetPDUUserData(message spi.Message) {
	//TODO implement me
	panic("implement me")
}

func (t TPDU) GetPDUSource() *Address {
	panic("implement me")
}

func (t TPDU) SetPDUSource(source *Address) {
	panic("implement me")
}

func (t TPDU) GetPDUDestination() *Address {
	panic("implement me")
}

func (t TPDU) SetPDUDestination(address *Address) {
	panic("implement me")
}

func (t TPDU) SetExpectingReply(b bool) {
	panic("implement me")
}

func (t TPDU) GetExpectingReply() bool {
	panic("implement me")
}

func (t TPDU) SetNetworkPriority(priority readWriteModel.NPDUNetworkPriority) {
	panic("implement me")
}

func (t TPDU) GetNetworkPriority() readWriteModel.NPDUNetworkPriority {
	panic("implement me")
}

func (t TPDU) Serialize() ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (t TPDU) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	panic("implement me")
}

func (t TPDU) GetLengthInBytes(ctx context.Context) uint16 {
	//TODO implement me
	panic("implement me")
}

func (t TPDU) GetLengthInBits(ctx context.Context) uint16 {
	//TODO implement me
	panic("implement me")
}

func (t TPDU) GetPDUUserData() spi.Message {
	//TODO implement me
	panic("implement me")
}

func (t TPDU) Update(pci Arg) error {
	//TODO implement me
	panic("implement me")
}

func (t TPDU) SetPduData(bytes []byte) {
	//TODO implement me
	panic("implement me")
}

func (t TPDU) GetPduData() []byte {
	//TODO implement me
	panic("implement me")
}

func (t TPDU) Get() (byte, error) {
	//TODO implement me
	panic("implement me")
}

func (t TPDU) GetShort() (int16, error) {
	//TODO implement me
	panic("implement me")
}

func (t TPDU) GetLong() (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (t TPDU) GetData(dlen int) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (t TPDU) Put(b byte) {
	//TODO implement me
	panic("implement me")
}

func (t TPDU) PutData(b ...byte) {
	//TODO implement me
	panic("implement me")
}

func (t TPDU) PutShort(i uint16) {
	//TODO implement me
	panic("implement me")
}

func (t TPDU) PutLong(i uint32) {
	//TODO implement me
	panic("implement me")
}

func (t TPDU) GetPCI() PCI {
	//TODO implement me
	panic("implement me")
}

func (t TPDU) DeepCopy() any {
	panic("implement me")
}

type Anon struct {
	TPDU
}

func TestMatchPdu(t *testing.T) {
	testingLogger := testutils.ProduceTestingLogger(t)

	tpdu := TPDU{x: []byte{1}}
	anon := Anon{TPDU{x: []byte("Anon")}}

	// no criteria passes
	assert.True(t, MatchPdu(testingLogger, tpdu, nil, nil))
	assert.True(t, MatchPdu(testingLogger, anon, nil, nil))

	// matching/not matching types
	assert.True(t, MatchPdu(testingLogger, tpdu, TPDU{}, nil))
	assert.False(t, MatchPdu(testingLogger, tpdu, Anon{}, nil))
	// Note the other testcase is irrelevant as we don't have dynamic types

	// matching/not matching attributes
	assert.True(t, MatchPdu(testingLogger, tpdu, nil, map[KnownKey]any{"x": []byte{1}}))
	assert.False(t, MatchPdu(testingLogger, tpdu, nil, map[KnownKey]any{"x": []byte{2}}))
	assert.False(t, MatchPdu(testingLogger, tpdu, nil, map[KnownKey]any{"y": []byte{1}}))
	assert.False(t, MatchPdu(testingLogger, anon, nil, map[KnownKey]any{"x": []byte{1}}))

	// matching/not matching types and attributes
	assert.True(t, MatchPdu(testingLogger, tpdu, TPDU{}, map[KnownKey]any{"x": []byte{1}}))
	assert.False(t, MatchPdu(testingLogger, tpdu, TPDU{}, map[KnownKey]any{"x": []byte{2}}))
	assert.False(t, MatchPdu(testingLogger, tpdu, TPDU{}, map[KnownKey]any{"y": []byte{1}}))
	assert.False(t, MatchPdu(testingLogger, anon, Anon{}, map[KnownKey]any{"x": []byte{1}}))
}

func TestState(t *testing.T) {
	t.Run("test_state_doc", func(t *testing.T) {
		testingLogger := testutils.ProduceTestingLogger(t)

		// change the doc string
		ts := NewState(testingLogger, nil, "")
		ns := ts.Doc("test state")

		assert.Equal(t, "test state", ts.DocString())
		assert.Same(t, ts, ns)
	})
	t.Run("test_success", func(t *testing.T) {
		testingLogger := testutils.ProduceTestingLogger(t)

		ts := NewState(testingLogger, nil, "")
		ns := ts.Success("")
		assert.True(t, ts.IsSuccessState())
		assert.Same(t, ts, ns)

		assert.Panics(t, func() {
			ts.Success("")
		})
		assert.Panics(t, func() {
			ts.Fail("")
		})
	})
	t.Run("test_state_fail", func(t *testing.T) {
		testingLogger := testutils.ProduceTestingLogger(t)

		ts := NewState(testingLogger, nil, "")
		ns := ts.Fail("")
		assert.True(t, ts.IsFailState())
		assert.Same(t, ts, ns)

		assert.Panics(t, func() {
			ts.Success("")
		})
		assert.Panics(t, func() {
			ts.Fail("")
		})
	})
}

func TestStateMachine(t *testing.T) {
	t.Run("test_state_machine_run", func(t *testing.T) {
		testingLogger := testutils.ProduceTestingLogger(t)

		// create a state machine //TODO: fix nil requirement
		var init func()
		tsm, init := NewStateMachine(testingLogger, nil)
		init()

		// run the machine
		err := tsm.Run()
		assert.NoError(t, err)

		assert.True(t, tsm.IsRunning())
		assert.Same(t, tsm.GetStartState(), tsm.GetCurrentState())
	})
	t.Run("test_state_machine_success", func(t *testing.T) {
		testingLogger := testutils.ProduceTestingLogger(t)

		// create a state machine
		tsm := NewTrappedStateMachine(testingLogger)

		// make the start state a sucess
		tsm.GetStartState().Success("")

		// run the machine
		err := tsm.Run()
		assert.NoError(t, err)

		assert.False(t, tsm.IsRunning())
		assert.True(t, tsm.GetCurrentState().IsSuccessState())
	})
	t.Run("test_state_machine_fail", func(t *testing.T) {
		testingLogger := testutils.ProduceTestingLogger(t)

		// create a state machine
		tsm := NewTrappedStateMachine(testingLogger)

		// make the start state a sucess
		tsm.GetStartState().Fail("")

		// run the machine
		err := tsm.Run()
		assert.NoError(t, err)

		assert.False(t, tsm.IsRunning())
		assert.True(t, tsm.GetCurrentState().IsFailState())
	})
	t.Run("test_state_machine_send", func(t *testing.T) {
		testingLogger := testutils.ProduceTestingLogger(t)

		// create a state machine
		tsm := NewTrappedStateMachine(testingLogger)

		// Make a pdu object
		pdu := NewPDU(Nothing())

		// make a send transition from start to success, run the machine
		tsm.GetStartState().Send(pdu, nil).Success("")
		err := tsm.Run()
		assert.NoError(t, err)

		// check for success
		assert.False(t, tsm.IsRunning())
		assert.True(t, tsm.GetCurrentState().IsSuccessState())

		// check the callbacks
		assert.IsType(t, NewPDU(Nothing()), tsm.GetBeforeSendPdu())
		assert.IsType(t, NewPDU(Nothing()), tsm.GetAfterSendPdu())

		// make sure pdu was sent
		assert.Same(t, pdu, tsm.GetSent())

		// check the transaction log
		require.Equal(t, 1, len(tsm.GetTransactionLog()))
		assert.Equal(t, tsm.GetTransactionLog()[0].Pdu, pdu)
	})
	t.Run("test_state_machine_receive", func(t *testing.T) {
		testingLogger := testutils.ProduceTestingLogger(t)

		// create a state machine
		tsm := NewTrappedStateMachine(testingLogger)

		// Make a pdu object
		pdu := TPDU{}

		// make a send transition from start to success, run the machine
		tsm.GetStartState().Receive(NA(pdu), NoKWArgs()).Success("")
		err := tsm.Run()
		assert.NoError(t, err)

		// check for still running
		assert.True(t, tsm.IsRunning())

		// tell the machine it is receiving the pdu
		err = tsm.Receive(NA(pdu), NoKWArgs())
		require.NoError(t, err)

		// check for success
		assert.False(t, tsm.IsRunning())
		assert.True(t, tsm.GetCurrentState().IsSuccessState())

		// check the callbacks
		assert.IsType(t, TPDU{}, tsm.GetBeforeReceivePdu())
		assert.IsType(t, TPDU{}, tsm.GetAfterReceivePdu())

		// check the transaction log
		require.Equal(t, 1, len(tsm.GetTransactionLog()))
		assert.Equal(t, tsm.GetTransactionLog()[0].Pdu, pdu)
	})
	t.Run("test_state_machine_unexpected", func(t *testing.T) {
		testingLogger := testutils.ProduceTestingLogger(t)

		// create a state machine
		tsm := NewTrappedStateMachine(testingLogger)

		// Make a pdu object
		goodPdu := TPDU{a: 1}
		_ = goodPdu
		badPdu := TPDU{b: 2}

		// make a send transition from start to success, run the machine
		tsm.GetStartState().Receive(NA(TPDU{}), NKW(KnownKey("a"), 1)).Success("")
		err := tsm.Run()
		assert.NoError(t, err)

		// check for still running
		assert.True(t, tsm.IsRunning())

		// give the machine a bad pdu
		err = tsm.Receive(NA(badPdu), NoKWArgs())
		require.NoError(t, err)

		// check for fail
		assert.False(t, tsm.IsRunning())
		assert.True(t, tsm.GetCurrentState().IsFailState())
		assert.Same(t, tsm.GetCurrentState(), tsm.GetUnexpectedReceiveState())

		// check the callbacks
		assert.Equal(t, badPdu, tsm.GetUnexpectedReceivePDU())

		// check the transaction log
		require.Equal(t, 1, len(tsm.GetTransactionLog()))
		assert.Equal(t, tsm.GetTransactionLog()[0].Pdu, badPdu)
	})
	t.Run("test_state_machine_call", func(t *testing.T) {
		testingLogger := testutils.ProduceTestingLogger(t)

		// simpleHook
		called := false
		_called := func(args Args, kwArgs KWArgs) error {
			called = args[0].(bool)
			return nil
		}

		// create a state machine
		tsm := NewTrappedStateMachine(testingLogger)

		// make a send transition from start to success, run the machine
		tsm.GetStartState().Call(_called, NA(true), NoKWArgs()).Success("")
		err := tsm.Run()
		assert.NoError(t, err)

		// check for success
		assert.False(t, tsm.IsRunning())
		assert.True(t, tsm.IsSuccessState())

		// check for the call
		assert.True(t, called)
	})
	t.Run("test_state_machine_call_exception", func(t *testing.T) {
		testingLogger := testutils.ProduceTestingLogger(t)

		// simpleHook
		called := false
		_called := func(args Args, kwArgs KWArgs) error {
			called = args[0].(bool)
			return AssertionError{Message: "error"}
		}

		// create a state machine
		tsm := NewTrappedStateMachine(testingLogger)

		// make a send transition from start to success, run the machine
		tsm.GetStartState().Call(_called, NA(true), NoKWArgs())
		err := tsm.Run()
		assert.NoError(t, err)

		// check for success
		assert.False(t, tsm.IsRunning())
		assert.True(t, tsm.IsFailState())

		// check for the call
		assert.True(t, called)
	})
	t.Run("test_state_machine_loop_01", func(t *testing.T) {
		testingLogger := testutils.ProduceTestingLogger(t)

		// create a state machine
		tsm := NewTrappedStateMachine(testingLogger)

		// Make a pdu object
		firstPdu := TPDU{a: 1}
		t.Log(firstPdu)
		secondPdu := TPDU{a: 2}
		t.Log(secondPdu)

		// after sending the first pdu, wait for the second
		s0 := tsm.GetStartState()
		s1 := s0.Send(firstPdu, nil)
		s2 := s1.Receive(NA(TPDU{}), NKW(KnownKey("a"), 2))
		s2.Success("")

		// run the machine
		err := tsm.Run()
		assert.NoError(t, err)

		// check for still running and waiting
		assert.True(t, tsm.IsRunning())
		assert.Same(t, s1, tsm.GetCurrentState())

		// give the machine the second pdu
		err = tsm.Receive(NA(secondPdu), NoKWArgs())
		require.NoError(t, err)

		// check for success
		assert.False(t, tsm.IsRunning())
		assert.True(t, tsm.GetCurrentState().IsSuccessState())
		t.Log("success")

		// check the callbacks
		assert.Equal(t, firstPdu, tsm.GetBeforeSendPdu())
		assert.Equal(t, firstPdu, tsm.GetAfterSendPdu())
		assert.Equal(t, secondPdu, tsm.GetBeforeReceivePdu())
		assert.Equal(t, secondPdu, tsm.GetAfterReceivePdu())
		t.Log("callbacks passed")

		// check the transaction log
		assert.Equal(t, tsm.GetTransactionLog()[0].Pdu, firstPdu)
		assert.Equal(t, tsm.GetTransactionLog()[1].Pdu, secondPdu)
	})
	t.Run("test_state_machine_loop_02", func(t *testing.T) {
		testingLogger := testutils.ProduceTestingLogger(t)

		// create a state machine
		tsm := NewTrappedStateMachine(testingLogger)

		// Make a pdu object
		firstPdu := TPDU{a: 1}
		t.Log(firstPdu)
		secondPdu := TPDU{a: 2}
		t.Log(secondPdu)

		// when the first pdu is received, send the second
		s0 := tsm.GetStartState()
		s1 := s0.Receive(NA(TPDU{}), NKW(KnownKey("a"), 1))
		s2 := s1.Send(secondPdu, nil)
		s2.Success("")

		// run the machine
		err := tsm.Run()
		assert.NoError(t, err)

		// check for still running and waiting
		assert.True(t, tsm.IsRunning())

		// give the machine the first pdu
		err = tsm.Receive(NA(firstPdu), NoKWArgs())
		require.NoError(t, err)

		// check for success
		assert.False(t, tsm.IsRunning())
		assert.True(t, tsm.GetCurrentState().IsSuccessState())
		t.Log("success")

		// check the callbacks
		assert.Equal(t, firstPdu, tsm.GetBeforeReceivePdu())
		assert.Equal(t, firstPdu, tsm.GetAfterReceivePdu())
		assert.Equal(t, secondPdu, tsm.GetBeforeSendPdu())
		assert.Equal(t, secondPdu, tsm.GetAfterSendPdu())
		t.Log("callbacks passed")

		// check the transaction log
		assert.Equal(t, tsm.GetTransactionLog()[0].Pdu, firstPdu)
		assert.Equal(t, tsm.GetTransactionLog()[1].Pdu, secondPdu)
	})
}

func TestStateMachineTimeout1(t *testing.T) {
	ExclusiveGlobalTimeMachine(t)
	testingLogger := testutils.ProduceTestingLogger(t)

	// create a state machine
	tsm := NewTrappedStateMachine(testingLogger)

	// make a timeout transition from start to success
	tsm.GetStartState().Timeout(1*time.Second, nil).Success("")

	// reset the time machine
	ResetTimeMachine(StartTime)
	t.Log("time machine reset")

	err := tsm.Run()
	require.NoError(t, err)

	RunTimeMachine(testingLogger, 60*time.Second, time.Time{})
	t.Log("time machine finished")

	// check for success
	assert.False(t, tsm.IsRunning())
	assert.True(t, tsm.GetCurrentState().IsSuccessState())
}

func TestStateMachineTimeout2(t *testing.T) {
	testingLogger := testutils.ProduceTestingLogger(t)
	ExclusiveGlobalTimeMachine(t)

	// make some pdus
	firstPdu := TPDU{a: 1}
	t.Log(firstPdu)
	secondPdu := TPDU{a: 2}
	t.Log(secondPdu)

	// create a state machine
	tsm := NewTrappedStateMachine(testingLogger)
	s0 := tsm.GetStartState()

	// send something, wait, send something, wait, success
	s1 := s0.Send(firstPdu, nil)
	s2 := s1.Timeout(1*time.Millisecond, nil)
	s3 := s2.Send(secondPdu, nil)
	s4 := s3.Timeout(1*time.Millisecond, nil).Success("")
	_ = s4

	// reset the time machine
	ResetTimeMachine(StartTime)
	t.Log("time machine reset")

	err := tsm.Run()
	require.NoError(t, err)

	RunTimeMachine(testingLogger, 60*time.Millisecond, time.Time{})
	t.Log("time machine finished")

	// check for success
	assert.False(t, tsm.IsRunning())
	assert.True(t, tsm.GetCurrentState().IsSuccessState())

	// check the transaction log
	assert.Len(t, tsm.GetTransactionLog(), 2)
	assert.Equal(t, tsm.GetTransactionLog()[0].Pdu, firstPdu)
	assert.Equal(t, tsm.GetTransactionLog()[1].Pdu, secondPdu)
}

func TestStateMachineGroup(t *testing.T) {
	t.Run("test_state_machine_group_success", func(t *testing.T) {
		ExclusiveGlobalTimeMachine(t)

		testingLogger := testutils.ProduceTestingLogger(t)

		// create a state machine group
		smg := NewStateMachineGroup(testingLogger)

		// create a trapped state machine, start state is success
		tsm := NewTrappedStateMachine(testingLogger)
		tsm.GetStartState().Success("")

		// add it to the group
		smg.Append(tsm)

		// reset the time machine
		ResetTimeMachine(StartTime)
		t.Log("time machine reset")

		// tell the group to run
		err := smg.Run()
		require.NoError(t, err)

		RunTimeMachine(testingLogger, 60*time.Second, time.Time{})
		t.Log("time machine finished")

		// check for success
		assert.False(t, tsm.IsRunning())
		assert.True(t, tsm.GetCurrentState().IsSuccessState())
		assert.False(t, smg.IsRunning())
		assert.True(t, smg.IsSuccessState())
	})
	t.Run("test_state_machine_group_success", func(t *testing.T) {
		ExclusiveGlobalTimeMachine(t)

		testingLogger := testutils.ProduceTestingLogger(t)

		// create a state machine group
		smg := NewStateMachineGroup(testingLogger)

		// create a trapped state machine, start state is success
		tsm := NewTrappedStateMachine(testingLogger)
		tsm.GetStartState().Fail("")

		// add it to the group
		smg.Append(tsm)

		// reset the time machine
		ResetTimeMachine(StartTime)
		t.Log("time machine reset")

		// tell the group to run
		err := smg.Run()
		require.NoError(t, err)

		RunTimeMachine(testingLogger, 60*time.Second, time.Time{})
		t.Log("time machine finished")

		// check for success
		assert.False(t, tsm.IsRunning())
		assert.True(t, tsm.GetCurrentState().IsFailState())
		assert.False(t, smg.IsRunning())
		assert.True(t, smg.IsFailState())
	})
}

func TestStateMachineEvents(t *testing.T) {
	t.Run("test_state_machine_event_01", func(t *testing.T) {
		ExclusiveGlobalTimeMachine(t)

		testingLogger := testutils.ProduceTestingLogger(t)

		// create a state machine group
		smg := NewStateMachineGroup(testingLogger)

		// create a trapped state machine, start state is success
		tsm1 := NewTrappedStateMachine(testingLogger)
		tsm1.GetStartState().SetEvent("e").Success("")
		smg.Append(tsm1)

		// create another trapped state machine, waiting for the event
		tsm2 := NewTrappedStateMachine(testingLogger)
		tsm2.GetStartState().WaitEvent("e", nil).Success("")
		smg.Append(tsm2)

		// reset the time machine
		ResetTimeMachine(StartTime)
		t.Log("time machine reset")

		// tell the group to run
		err := smg.Run()
		require.NoError(t, err)

		RunTimeMachine(testingLogger, 60*time.Second, time.Time{})
		t.Log("time machine finished")

		// check for success
		assert.True(t, tsm1.GetCurrentState().IsSuccessState())
		assert.True(t, tsm2.GetCurrentState().IsSuccessState())
		assert.False(t, smg.IsRunning())
		assert.True(t, smg.IsSuccessState())
	})
	t.Run("test_state_machine_event_02", func(t *testing.T) {
		ExclusiveGlobalTimeMachine(t)

		testingLogger := testutils.ProduceTestingLogger(t)

		// create a state machine group
		smg := NewStateMachineGroup(testingLogger)

		// create a trapped state machine, waiting for an event
		tsm1 := NewTrappedStateMachine(testingLogger)
		tsm1.GetStartState().WaitEvent("e", nil).Success("")
		smg.Append(tsm1)

		// create another trapped state machine, start state is success
		tsm2 := NewTrappedStateMachine(testingLogger)
		tsm2.GetStartState().SetEvent("e").Success("")
		smg.Append(tsm2)

		// reset the time machine
		ResetTimeMachine(StartTime)
		t.Log("time machine reset")

		// tell the group to run
		err := smg.Run()
		require.NoError(t, err)

		RunTimeMachine(testingLogger, 60*time.Second, time.Time{})
		t.Log("time machine finished")

		// check for success
		assert.True(t, tsm1.GetCurrentState().IsSuccessState())
		assert.True(t, tsm2.GetCurrentState().IsSuccessState())
		assert.False(t, smg.IsRunning())
		assert.True(t, smg.IsSuccessState())
	})
}
