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
	cryptoRand "crypto/rand"
	"fmt"
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

// RandomBytes returns length cryptographically secure random bytes. It is intended for
// security-sensitive values such as nonces and fails hard (panics) if the system CSPRNG is
// unavailable, because silently continuing with predictable values would undermine any key
// material derived from them.
func RandomBytes(length int) []byte {
	randomBytes := make([]byte, length)
	if _, err := cryptoRand.Read(randomBytes); err != nil {
		panic(fmt.Sprintf("unable to generate %d secure random bytes: %v", length, err))
	}
	return randomBytes
}

// RandomString returns a random string using the alphabet.
// Not cryptographically secure: use RandomBytes for nonces and other security-sensitive values.
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
