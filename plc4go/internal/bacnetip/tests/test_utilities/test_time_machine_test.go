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
	"fmt"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/apache/plc4x/plc4go/spi/testutils"

	"github.com/apache/plc4x/plc4go/internal/bacnetip"
	"github.com/apache/plc4x/plc4go/internal/bacnetip/tests"
)

type TimeMachineSuite struct {
	suite.Suite

	// flag to make sure the function was called
	sampleTaskFunctionCalled []time.Time

	log zerolog.Logger
}

func (suite *TimeMachineSuite) SetupTest() {
	t := suite.T()
	suite.log = testutils.ProduceTestingLogger(t)
	tests.LockGlobalTimeMachine(t)
	tests.NewGlobalTimeMachine(t)
}

func (suite *TimeMachineSuite) TearDownTest() {
	tests.ClearGlobalTimeMachine(suite.T())
}

type SampleOneShotTask struct {
	*bacnetip.OneShotTask

	processTaskCalled []time.Time

	log zerolog.Logger
}

func NewSampleOneShotTask(localLog zerolog.Logger) *SampleOneShotTask {
	s := &SampleOneShotTask{
		log: localLog,
	}
	s.OneShotTask = bacnetip.NewOneShotTask(s, nil)
	return s
}

func (s *SampleOneShotTask) ProcessTask() error {
	s.log.Debug().Time("current_time", tests.GlobalTimeMachineCurrentTime()).Msg("processing task")

	// add the current time
	s.processTaskCalled = append(s.processTaskCalled, tests.GlobalTimeMachineCurrentTime())
	return nil
}

func (suite *TimeMachineSuite) SampleTaskFunction() func(args bacnetip.Args, kwArgs bacnetip.KWArgs) error {
	return func(args bacnetip.Args, kwArgs bacnetip.KWArgs) error {
		currentTime := tests.GlobalTimeMachineCurrentTime()
		suite.log.Debug().Stringer("args", args).Stringer("kwArgs", kwArgs).Time("current_time", currentTime).Msg("sample_task_function")

		suite.sampleTaskFunctionCalled = append(suite.sampleTaskFunctionCalled, currentTime)
		return nil
	}
}

type SampleRecurringTask struct {
	*bacnetip.RecurringTask

	processTaskCalled []time.Time

	log zerolog.Logger
}

func NewSampleRecurringTask(localLog zerolog.Logger) *SampleRecurringTask {
	s := &SampleRecurringTask{
		log: localLog,
	}
	s.RecurringTask = bacnetip.NewRecurringTask(localLog, s)
	return s
}

func (s *SampleRecurringTask) ProcessTask() error {
	s.log.Debug().Time("current_time", tests.GlobalTimeMachineCurrentTime()).Msg("processing task")

	// add the current time
	s.processTaskCalled = append(s.processTaskCalled, tests.GlobalTimeMachineCurrentTime())
	return nil
}

func (s *SampleRecurringTask) String() string {
	return fmt.Sprintf("SampleRecurringTask{%s, processTaskCalled: %v}", s.RecurringTask, s.processTaskCalled)
}

func (suite *TimeMachineSuite) TestTimeMachineExists() {
	assert.True(suite.T(), tests.IsGlobalTimeMachineSet())
}

func (suite *TimeMachineSuite) TestEmptyRun() {
	// reset the time machine
	tests.ResetTimeMachine(tests.StartTime)

	// let it run
	tests.RunTimeMachine(suite.log, 60*time.Second, time.Time{})

	// 60 seconds have passed
	suite.Equal(60*time.Second, tests.GlobalTimeMachineCurrentTime().Sub(tests.StartTime))
}

func (suite *TimeMachineSuite) TestOneShotImmediate1() {
	// create a function task
	ft := NewSampleOneShotTask(suite.log)

	// Reset time machine
	tests.ResetTimeMachine(tests.StartTime)
	var startTime time.Time
	ft.InstallTask(bacnetip.WithInstallTaskOptionsWhen(startTime))
	tests.RunTimeMachine(suite.log, 60*time.Second, time.Time{})

	// function called, 60 seconds passed
	suite.Contains(ft.processTaskCalled, startTime)
	suite.Equal(60*time.Second, tests.GlobalTimeMachineCurrentTime().Sub(tests.StartTime))
}

func (suite *TimeMachineSuite) TestOneShotImmediate2() {
	// create a function task
	ft := NewSampleOneShotTask(suite.log)

	// run the functions sometime later
	t1, err := time.Parse("2006-01-02", "2000-06-06")
	suite.Require().NoError(err)
	suite.T().Log(t1)

	// reset the time machine to midnight, install the task, let it run
	startTime, err := time.Parse("2006-01-02", "2000-01-01")
	suite.Require().NoError(err)
	tests.ResetTimeMachine(startTime)
	ft.InstallTask(bacnetip.WithInstallTaskOptionsWhen(t1))
	stopTime, err := time.Parse("2006-01-02", "2001-01-01")
	suite.Require().NoError(err)
	tests.RunTimeMachine(suite.log, 0, stopTime)

	// function called, 60 seconds passed
	suite.Contains(ft.processTaskCalled, t1)
}

func (suite *TimeMachineSuite) TestFunctionTaskImmediate() {
	// create a function task
	ft := bacnetip.FunctionTask(suite.SampleTaskFunction(), bacnetip.NoArgs, bacnetip.NoKWArgs)
	suite.sampleTaskFunctionCalled = nil

	// reset the time machine to midnight, install the task, let it run
	tests.ResetTimeMachine(tests.StartTime)
	var now time.Time
	ft.InstallTask(bacnetip.WithInstallTaskOptionsWhen(now))
	tests.RunTimeMachine(suite.log, 60*time.Second, time.Time{})

	// function called, 60 seconds passed
	suite.Contains(suite.sampleTaskFunctionCalled, now)
	suite.Equal(60*time.Second, tests.GlobalTimeMachineCurrentTime().Sub(tests.StartTime))
}

func (suite *TimeMachineSuite) TestFunctionTaskDelay() {
	sampleDelay := 10 * time.Second

	// create a function task
	ft := bacnetip.FunctionTask(suite.SampleTaskFunction(), bacnetip.NoArgs, bacnetip.NoKWArgs)
	suite.sampleTaskFunctionCalled = nil

	// reset the time machine to midnight, install the task, let it run
	tests.ResetTimeMachine(tests.StartTime)
	var now time.Time
	when := now.Add(sampleDelay)
	ft.InstallTask(bacnetip.WithInstallTaskOptionsWhen(when))
	tests.RunTimeMachine(suite.log, 60*time.Second, time.Time{})

	// function called, 60 seconds passed
	suite.Contains(suite.sampleTaskFunctionCalled, when)
	suite.Equal(60*time.Second, tests.GlobalTimeMachineCurrentTime().Sub(tests.StartTime))
}

func (suite *TimeMachineSuite) TestRecurringTask1() {
	// create a function task
	ft := NewSampleRecurringTask(suite.log)

	// reset the time machine to midnight, install the task, let it run
	now := tests.StartTime
	tests.ResetTimeMachine(now)
	ft.InstallTask(bacnetip.WithInstallTaskOptionsInterval(1 * time.Second))
	tests.RunTimeMachine(suite.log, 5*time.Second, time.Time{})

	// function called, 60 seconds passed
	suite.Equal(now.Add(1*time.Second), ft.processTaskCalled[0])
	suite.Equal(now.Add(2*time.Second), ft.processTaskCalled[1])
	suite.Equal(now.Add(3*time.Second), ft.processTaskCalled[2])
	suite.Equal(now.Add(4*time.Second), ft.processTaskCalled[3])
	suite.Equal(now.Add(5*time.Second), ft.processTaskCalled[4])
	suite.Equal(now.Add(6*time.Second), ft.processTaskCalled[5])
	suite.InDelta(5*time.Second, tests.GlobalTimeMachineCurrentTime().Sub(now), float64(100*time.Millisecond))
}

func (suite *TimeMachineSuite) TestRecurringTask2() {
	// create a function task
	ft1 := NewSampleRecurringTask(suite.log)
	ft2 := NewSampleRecurringTask(suite.log)

	// reset the time machine to midnight, install the task, let it run
	tests.ResetTimeMachine(tests.StartTime)
	ft1.InstallTask(bacnetip.WithInstallTaskOptionsInterval(1000 * time.Millisecond))
	ft2.InstallTask(bacnetip.WithInstallTaskOptionsInterval(1500 * time.Millisecond))
	tests.RunTimeMachine(suite.log, 5*time.Second, time.Time{})

	// function called, 60 seconds passed
	suite.Require().Greater(len(ft1.processTaskCalled), 4)
	suite.Equal(tests.StartTime.Add(1*time.Second), ft1.processTaskCalled[0])
	suite.Equal(tests.StartTime.Add(2*time.Second), ft1.processTaskCalled[1])
	suite.Equal(tests.StartTime.Add(3*time.Second), ft1.processTaskCalled[2])
	suite.Equal(tests.StartTime.Add(4*time.Second), ft1.processTaskCalled[3])
	suite.Require().Greater(len(ft2.processTaskCalled), 3)
	suite.Equal(tests.StartTime.Add(1500*time.Millisecond), ft2.processTaskCalled[0])
	suite.Equal(tests.StartTime.Add(3000*time.Millisecond), ft2.processTaskCalled[1])
	suite.Equal(tests.StartTime.Add(4500*time.Millisecond), ft2.processTaskCalled[2])
	suite.InDelta(5*time.Second, tests.GlobalTimeMachineCurrentTime().Sub(tests.StartTime), float64(100*time.Millisecond))
}

func (suite *TimeMachineSuite) TestRecurringTask3() {
	// create a function task
	ft := NewSampleRecurringTask(suite.log)

	// reset the time machine to midnight, install the task, let it run
	startTime := time.Time{}.Add(1 * time.Hour) // We add an hour to avoid underflow
	tests.ResetTimeMachine(startTime)
	ft.InstallTask(bacnetip.WithInstallTaskOptionsInterval(1000 * time.Millisecond).WithOffset(100 * time.Millisecond))
	tests.RunTimeMachine(suite.log, 5*time.Second, time.Time{})

	// function called, 60 seconds passed
	suite.Equal(startTime.Add(100*time.Millisecond), ft.processTaskCalled[0])
	suite.Equal(startTime.Add(1100*time.Millisecond), ft.processTaskCalled[1])
	suite.Equal(startTime.Add(2100*time.Millisecond), ft.processTaskCalled[2])
	suite.Equal(startTime.Add(3100*time.Millisecond), ft.processTaskCalled[3])
	suite.Equal(startTime.Add(4100*time.Millisecond), ft.processTaskCalled[4])
	suite.InDelta(5*time.Second, tests.GlobalTimeMachineCurrentTime().Sub(startTime), float64(100*time.Millisecond))
}

func (suite *TimeMachineSuite) TestRecurringTask4() {
	// create a function task
	ft := NewSampleRecurringTask(suite.log)

	// reset the time machine to midnight, install the task, let it run
	tests.ResetTimeMachine(tests.StartTime)
	ft.InstallTask(bacnetip.WithInstallTaskOptionsInterval(1000 * time.Millisecond).WithOffset(-100 * time.Millisecond))
	tests.RunTimeMachine(suite.log, 5*time.Second, time.Time{})

	// function called, 60 seconds passed
	suite.Equal(tests.StartTime.Add(900*time.Millisecond), ft.processTaskCalled[0])
	suite.Equal(tests.StartTime.Add(1900*time.Millisecond), ft.processTaskCalled[1])
	suite.Equal(tests.StartTime.Add(2900*time.Millisecond), ft.processTaskCalled[2])
	suite.Equal(tests.StartTime.Add(3900*time.Millisecond), ft.processTaskCalled[3])
	suite.Equal(tests.StartTime.Add(4900*time.Millisecond), ft.processTaskCalled[4])
	suite.InDelta(5*time.Second, tests.GlobalTimeMachineCurrentTime().Sub(tests.StartTime), float64(100*time.Millisecond))
}

func (suite *TimeMachineSuite) TestRecurringTask5() {
	// create a function task
	ft := NewSampleRecurringTask(suite.log)

	// reset the time machine, install the task, let it run
	now, err := time.Parse("2006-01-02", "2000-01-01")
	suite.Require().NoError(err)
	tests.ResetTimeMachine(now)
	ft.InstallTask(bacnetip.WithInstallTaskOptionsInterval(86400 * time.Second))
	stopTime, err := time.Parse("2006-01-02", "2000-02-01")
	suite.Require().NoError(err)
	tests.RunTimeMachine(suite.log, 0, stopTime)

	// function called every day
	suite.Equal(32, len(ft.processTaskCalled))
}

func TestTimeMachine(t *testing.T) {
	suite.Run(t, new(TimeMachineSuite))
}
