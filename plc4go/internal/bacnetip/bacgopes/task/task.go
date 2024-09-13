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

package task

import (
	"time"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/debugging"
)

var _debug = CreateDebugPrinter()

type InstallTaskOptions struct {
	When     *time.Time
	Delta    *time.Duration
	Interval *time.Duration
	Offset   *time.Duration
}

func WithInstallTaskOptionsNone() InstallTaskOptions {
	return InstallTaskOptions{}
}

func WithInstallTaskOptionsWhen(when time.Time) InstallTaskOptions {
	return InstallTaskOptions{
		When: &when,
	}
}

func (i InstallTaskOptions) WithWhen(when time.Time) InstallTaskOptions {
	i.When = &when
	return i
}

func WithInstallTaskOptionsDelta(delta time.Duration) InstallTaskOptions {
	return InstallTaskOptions{
		Delta: &delta,
	}
}

func (i InstallTaskOptions) WithDelta(delta time.Duration) InstallTaskOptions {
	i.Delta = &delta
	return i
}

func WithInstallTaskOptionsInterval(interval time.Duration) InstallTaskOptions {
	return InstallTaskOptions{
		Interval: &interval,
	}
}

func (i InstallTaskOptions) WithInterval(interval time.Duration) InstallTaskOptions {
	i.Interval = &interval
	return i
}

func WithInstallTaskOptionsOffset(offset time.Duration) InstallTaskOptions {
	return InstallTaskOptions{
		Offset: &offset,
	}
}

func (i InstallTaskOptions) WithOffset(offset time.Duration) InstallTaskOptions {
	i.Offset = &offset
	return i
}
