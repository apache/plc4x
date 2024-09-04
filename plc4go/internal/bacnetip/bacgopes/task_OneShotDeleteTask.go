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

package bacgopes

import (
	"fmt"
	"time"
)

type OneShotDeleteTask struct {
	*Task
}

func NewOneShotDeleteTask(taskRequirements TaskRequirements, when *time.Time) *OneShotDeleteTask {
	o := &OneShotDeleteTask{}
	o.Task = NewTask(taskRequirements, func(task *Task) {
		task.taskTime = when
	})
	return o
}

func (r *OneShotDeleteTask) IsOneShotDeleteTask() bool {
	return true
}

func (r *OneShotDeleteTask) String() string {
	return fmt.Sprintf("OneShotDeleteTask(%v)", r.Task)
}
