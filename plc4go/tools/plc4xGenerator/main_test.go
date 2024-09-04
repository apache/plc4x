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

package main

import (
	"github.com/stretchr/testify/require"
	"os"
	"path"
	"testing"
	"time"

	"github.com/apache/plc4x/plc4go/spi/utils"
)

func TestBanana(t *testing.T) {
	typeNames = utils.MakePtr("Example")
	suffix = utils.MakePtr("_test")
	tests = utils.MakePtr(true)
	pkgIndex = utils.MakePtr(1)
	output = utils.MakePtr(path.Join(t.TempDir(), "aTestOutput.go"))
	main()
	time.Sleep(1 * time.Second)
	file, err := os.ReadFile(outputFile)
	require.NoError(t, err)
	t.Logf("Output:\n%s", string(file))
}
