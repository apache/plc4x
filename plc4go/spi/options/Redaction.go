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

package options

import "regexp"

// RedactConnectionString removes credentials from a connection string before it is logged.
//
// plc4go logged connection strings verbatim, so a password in one reached the log in clear at
// debug level - the driver manager and the connection cache both do it, on every connect.
//
// Which parameters are secret is decided by name here, unlike plc4j, where a @Secret marking on
// the configuration field decides it. plc4go has no equivalent to read at this point: the manager
// logs before any driver has looked at the string, and a driver's options are parsed by hand
// rather than described anywhere this code can see. Rendering a configuration *is* marking-driven
// (the generator's `secret:"true"` tag); this covers the raw string on its way past.
//
// psk-identity is deliberately not matched: it says which key was refused, which is what an
// operator needs when a handshake fails, and hiding it protects nothing.
var (
	secretParameter = regexp.MustCompile(`(?i)([?&][^=&]*(password|passwd|secret|token|psk-key|passphrase)[^=&]*=)[^&]*`)
	// Credentials in a URI authority have no parameter name to match, so they are removed
	// structurally: everything between the ':' of the userinfo and the '@' that ends it.
	userinfo = regexp.MustCompile(`(//[^/@\s]*:)([^/@\s]*)(@)`)
)

// Redacted is what a removed value is replaced with.
const Redacted = "******"

// RedactConnectionString returns the connection string with every credential replaced.
func RedactConnectionString(connectionString string) string {
	redacted := userinfo.ReplaceAllString(connectionString, "${1}"+Redacted+"${3}")
	return secretParameter.ReplaceAllString(redacted, "${1}"+Redacted)
}
