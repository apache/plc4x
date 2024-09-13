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

package utils

import (
	"math/rand"
	"time"
)

var alphabet []rune

func init() {
	offset := 'a' - 'A'
	for i := range 26 {
		r := rune(i + 'A')
		alphabet = append(alphabet, r, r+offset)
	}
}

// RandomString returns a random string using the alphabet
func RandomString(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomString := make([]rune, length)
	for i := range randomString {
		randomString[i] = alphabet[r.Intn(len(alphabet))]
	}
	return string(randomString)
}

// AlternateStringer can be implemented by stuff using the generator to give an alternate fmt.Stringer representation
type AlternateStringer interface {
	AlternateString() (v string, use bool)
}
