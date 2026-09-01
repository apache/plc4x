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

import org.apache.plc4x.java.spi.config.Configuration;
import org.apache.plc4x.java.spi.config.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.config.annotations.Secret;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * What is redacted is decided by the {@link Secret} markings on the configuration classes, with a
 * name-based backstop for parameters no configuration declares.
 */
class ConnectionStringRedactorTest {

    private static final String REDACTED = "******";
    private static final String NESTED = "s7%3A%2F%2Foperator%3Ahunter2%40plc%3A102%3Fpassword%3Dhunter2";

    static class DriverConfig implements Configuration {
        @Secret
        @ConfigurationParameter("password")
        public String password;

        @ConfigurationParameter("username")
        public String username;

        @ConfigurationParameter("remote-connection-string")
        public String remoteConnectionString;
    }

    static class TlsConfig implements Configuration {
        @Secret
        @ConfigurationParameter("keystore-password")
        public String keystorePassword;

        @ConfigurationParameter("keystore")
        public String keystore;

        @ConfigurationParameter("verify")
        public boolean verifySsl;
    }

    static class PskConfig implements Configuration {
        @Secret
        @ConfigurationParameter("psk-key")
        public String pskKey;

        // Not marked: the identity says which key was refused.
        @ConfigurationParameter("psk-identity")
        public String pskIdentity;
    }

    private static String redact(String url) {
        return ConnectionStringRedactor.redact(url, DriverConfig.class, "tls", TlsConfig.class);
    }

    @Test
    @DisplayName("masks a marked parameter's value")
    void masksAMarkedParameter() {
        assertEquals(
            "plc4x:tls://host:59837?remote-connection-string=s7&username=op&password=******&tls.verify=false",
            redact("plc4x:tls://host:59837?remote-connection-string=s7&username=op&password=hunter2&tls.verify=false"));
    }

    @Test
    @DisplayName("masks a marked parameter that leads the query string")
    void masksALeadingParameter() {
        assertEquals("plc4x://host?password=******&username=op",
            redact("plc4x://host?password=hunter2&username=op"));
    }

    @Test
    @DisplayName("masks a transport parameter under the transport's prefix")
    void masksATransportParameterUnderItsPrefix() {
        assertEquals("plc4x:tls://host?password=******&tls.keystore-password=******",
            redact("plc4x:tls://host?password=abc&tls.keystore-password=def"));
    }

    @Test
    @DisplayName("masks a pre-shared key, whose name no word list would catch")
    void masksAPreSharedKey() {
        // The reason the marking exists: this name contains none of password, secret or token,
        // and it is the credential for the whole connection.
        assertEquals("plc4x:tls-psk://host?tls-psk.psk-identity=plc4x&tls-psk.psk-key=******",
            ConnectionStringRedactor.redact(
                "plc4x:tls-psk://host?tls-psk.psk-identity=plc4x&tls-psk.psk-key=0011deadbeef",
                DriverConfig.class, "tls-psk", PskConfig.class));
    }

    @Test
    @DisplayName("leaves unmarked parameters untouched")
    void leavesUnmarkedParametersUntouched() {
        String url = "plc4x:tls://host:59837?remote-connection-string=s7&username=op&tls.verify=false";
        assertEquals(url, redact(url));
    }

    @Test
    @DisplayName("leaves the things named like keys that are not keys")
    void leavesKeyShapedNonSecretsUntouched() {
        // Masking these would cost an operator the diagnosis and protect nothing: a path, a store
        // type, a key size, a boolean, and the identity that says *which* key failed.
        String url = "plc4x:tls://host?keystore=/etc/plc4x/client.p12"
            + "&keystore-type=pkcs12&generated-key-size=2048&log-session-keys=false"
            + "&tls-psk.psk-identity=plc4x&allow-insecure-credentials=false";
        assertEquals(url, redact(url));
    }

    @Test
    @DisplayName("does not match a parameter that merely looks similar")
    void doesNotMaskUsername() {
        assertEquals("plc4x://host?username=secretive-bob", redact("plc4x://host?username=secretive-bob"));
    }

    @Test
    @DisplayName("masks a marked parameter whose case does not match the declaration")
    void masksAWronglyCasedParameter() {
        // It would not bind - and is reported as unknown - but the user typed a real password, and
        // it must not reach the log on its way to being ignored.
        assertEquals("plc4x://host?PassWord=******", redact("plc4x://host?PassWord=hunter2"));
    }

    @Test
    @DisplayName("masks a secret-looking parameter no configuration declares")
    void masksAnUndeclaredSecretLookingParameter() {
        // The backstop. The markings can say nothing about a parameter nothing declares, and a
        // credential passed under a name this build does not know is still a credential.
        assertEquals("plc4x://host?api-token=******&passphrase=******&my-password=******",
            redact("plc4x://host?api-token=abc&passphrase=open-sesame&my-password=hunter2"));
    }

    @Test
    @DisplayName("the backstop reads names, so a misspelling it cannot recognise gets through")
    void doesNotCatchAMisspellingThatLooksLikeNothing() {
        // Worth stating rather than discovering: the backstop matches names *containing* a
        // credential word. "passwrod" contains none, so nothing can tell it from any other
        // unknown parameter, and its value is logged. The unknown-parameter warning names it,
        // which is the only signal available - a marking cannot exist for a name nobody declared.
        assertEquals("plc4x://host?passwrod=hunter2", redact("plc4x://host?passwrod=hunter2"));
    }

    @Test
    @DisplayName("masks credentials in the URI's userinfo, which have no parameter name")
    void masksUserinfoCredentials() {
        assertEquals("s7://operator:******@plc:102?username=op",
            redact("s7://operator:hunter2@plc:102?username=op"));
    }

    @Test
    @DisplayName("masks a userinfo password that contains a colon")
    void masksUserinfoCredentialsContainingAColon() {
        // A colon is legal inside a password. The user segment stops at the first colon, so
        // everything after it is the password; a greedy one would publish "operator:hun".
        assertEquals("s7://operator:******@plc:102",
            redact("s7://operator:hun:ter2@plc:102"));
    }

    @Test
    @DisplayName("masks a percent-encoded parameter name")
    void masksAPercentEncodedName() {
        // The driver reads this as "password" once the query is decoded, so the value is a
        // credential however the name was written.
        assertEquals("plc4x://host?%70assword=" + REDACTED,
            redact("plc4x://host?%70assword=hunter2"));
    }

    @Test
    @DisplayName("masks credentials nested in a connection string passed as a value")
    void masksCredentialsInsideANestedConnectionString() {
        // The PLC4X proxy driver takes a whole PLC URL as a parameter. Nothing about the outer
        // name says "secret", and the encoding hides the inner one from every pattern.
        String redacted = redact("plc4x://proxy?remote-connection-string=" + NESTED);

        assertFalse(redacted.contains("hunter2"), redacted);
        assertTrue(redacted.contains("remote-connection-string="), "the parameter is still named");
    }

    @Test
    @DisplayName("keeps a nested connection string readable apart from its credentials")
    void keepsTheNestedEndpointVisible() {
        String redacted = redact("plc4x://proxy?remote-connection-string=s7://operator:hunter2@plc:102");

        assertFalse(redacted.contains("hunter2"), redacted);
        assertTrue(redacted.contains("plc:102"), "which PLC the proxy talks to is the diagnosis");
    }

    @Test
    @DisplayName("leaves a connection string without parameters untouched")
    void leavesNoParamStringUntouched() {
        String url = "plc4x:tls://host:59837";
        assertEquals(url, redact(url));
    }

    @Test
    @DisplayName("handles null")
    void handlesNull() {
        assertNull(redact(null));
    }

    @Test
    @DisplayName("handles a driver with no configuration class")
    void handlesNoConfigurationClass() {
        assertEquals("plc4x://host?password=******",
            ConnectionStringRedactor.redact("plc4x://host?password=hunter2", null, null, null));
    }
}
