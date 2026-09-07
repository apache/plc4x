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
package org.apache.plc4x.java.spi.config;

import org.apache.plc4x.java.spi.config.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.config.annotations.Secret;
import org.junit.jupiter.api.Test;

import java.util.Set;

import static org.junit.jupiter.api.Assertions.*;

class SecretParametersTest {

    static class Credentials implements Configuration {
        @Secret
        @ConfigurationParameter("password")
        public String password;

        @ConfigurationParameter("username")
        public String username;

        /** Marked, but not a parameter: covered by the toString rule, nothing to match in a URL. */
        @Secret
        public String derivedKey;

        /** A parameter with no explicit name is addressed by its field name. */
        @Secret
        @ConfigurationParameter
        public String sessionToken;
    }

    static class InheritedCredentials extends Credentials {
        @Secret
        @ConfigurationParameter("psk-key")
        public String pskKey;
    }

    @Test
    void reportsTheNamesOfMarkedParameters() {
        assertEquals(Set.of("password", "sessionToken"), SecretParametersTest.namesOf(Credentials.class));
    }

    @Test
    void leavesUnmarkedParametersAlone() {
        assertFalse(SecretParameters.namesFor(Credentials.class).contains("username"));
    }

    @Test
    void aMarkedFieldThatIsNotAParameterContributesNoName() {
        // It has no name to match in a connection string; the toString rule still covers it.
        assertFalse(SecretParameters.namesFor(Credentials.class).contains("derivedKey"));
        assertTrue(SecretParameters.fieldsOf(Credentials.class).stream()
            .anyMatch(field -> field.getName().equals("derivedKey")));
    }

    @Test
    void includesInheritedMarkings() {
        Set<String> names = SecretParameters.namesFor(InheritedCredentials.class);
        assertTrue(names.contains("psk-key"), "its own");
        assertTrue(names.contains("password"), "and its parent's");
    }

    @Test
    void handlesAClassWithNoMarkings() {
        assertTrue(SecretParameters.namesFor(String.class).isEmpty());
    }

    @Test
    void handlesNull() {
        assertTrue(SecretParameters.namesFor(null).isEmpty());
    }

    private static Set<String> namesOf(Class<?> type) {
        return SecretParameters.namesFor(type);
    }
}
