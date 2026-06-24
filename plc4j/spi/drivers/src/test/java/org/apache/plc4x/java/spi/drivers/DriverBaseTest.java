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
package org.apache.plc4x.java.spi.drivers;

import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Nested;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNull;

public class DriverBaseTest {

    @Nested
    @DisplayName("redactSecrets")
    class RedactSecrets {

        @Test
        @DisplayName("masks the password parameter value")
        void masksPassword() {
            assertEquals(
                "plc4x:tls://host:59837?remote-connection-string=s7&username=op&password=***&tls.verify-ssl=false",
                DriverBase.redactSecrets(
                    "plc4x:tls://host:59837?remote-connection-string=s7&username=op&password=hunter2&tls.verify-ssl=false"));
        }

        @Test
        @DisplayName("masks a password that is the first query parameter")
        void masksLeadingPassword() {
            assertEquals(
                "plc4x://host?password=***&username=op",
                DriverBase.redactSecrets("plc4x://host?password=hunter2&username=op"));
        }

        @Test
        @DisplayName("masks transport-prefixed and multiple secret parameters")
        void masksMultipleSecrets() {
            assertEquals(
                "plc4x:tls://host?password=***&tls.keystore-password=***&api-token=***",
                DriverBase.redactSecrets(
                    "plc4x:tls://host?password=abc&tls.keystore-password=def&api-token=ghi"));
        }

        @Test
        @DisplayName("is case-insensitive on the parameter name")
        void caseInsensitive() {
            assertEquals(
                "plc4x://host?PassWord=***&Secret=***",
                DriverBase.redactSecrets("plc4x://host?PassWord=abc&Secret=xyz"));
        }

        @Test
        @DisplayName("leaves non-secret parameters untouched")
        void leavesNonSecretsUntouched() {
            String url = "plc4x:tls://host:59837?remote-connection-string=s7&username=op&tls.verify-ssl=false";
            assertEquals(url, DriverBase.redactSecrets(url));
        }

        @Test
        @DisplayName("leaves a connection string without parameters untouched")
        void leavesNoParamStringUntouched() {
            String url = "plc4x:tls://host:59837";
            assertEquals(url, DriverBase.redactSecrets(url));
        }

        @Test
        @DisplayName("does not match parameter names that merely look similar (e.g. username)")
        void doesNotMaskUsername() {
            assertEquals(
                "plc4x://host?username=secretive-bob",
                DriverBase.redactSecrets("plc4x://host?username=secretive-bob"));
        }

        @Test
        @DisplayName("handles null")
        void handlesNull() {
            assertNull(DriverBase.redactSecrets(null));
        }
    }
}
