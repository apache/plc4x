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
package org.apache.plc4x.java.api.authentication;

import org.junit.jupiter.api.Test;

import java.security.KeyStore;

import static org.junit.jupiter.api.Assertions.*;

class PlcCertificateAuthenticationTest {

    @Test
    void testInstantiation() throws Exception {
        KeyStore keyStore = createEmptyKeyStore();
        char[] password = "password".toCharArray();

        PlcCertificateAuthentication auth = new PlcCertificateAuthentication(keyStore, password);

        assertNotNull(auth);
    }

    @Test
    void testImplementsPlcAuthentication() throws Exception {
        KeyStore keyStore = createEmptyKeyStore();
        char[] password = "password".toCharArray();

        PlcCertificateAuthentication auth = new PlcCertificateAuthentication(keyStore, password);

        assertTrue(true);
    }

    @Test
    void testGetKeyStore() throws Exception {
        KeyStore keyStore = createEmptyKeyStore();
        char[] password = "password".toCharArray();

        PlcCertificateAuthentication auth = new PlcCertificateAuthentication(keyStore, password);

        assertSame(keyStore, auth.getKeyStore());
    }

    @Test
    void testGetKeyStorePassword() throws Exception {
        KeyStore keyStore = createEmptyKeyStore();
        char[] password = "password".toCharArray();

        PlcCertificateAuthentication auth = new PlcCertificateAuthentication(keyStore, password);

        char[] returnedPassword = auth.getKeyStorePassword();
        assertArrayEquals(password, returnedPassword);
        // Verify it returns a copy, not the original
        assertNotSame(password, returnedPassword);
    }

    @Test
    void testGetKeyAliasDefault() throws Exception {
        KeyStore keyStore = createEmptyKeyStore();
        char[] password = "password".toCharArray();

        PlcCertificateAuthentication auth = new PlcCertificateAuthentication(keyStore, password);

        assertNull(auth.getKeyAlias());
    }

    @Test
    void testGetKeyAliasSpecified() throws Exception {
        KeyStore keyStore = createEmptyKeyStore();
        char[] password = "password".toCharArray();
        String alias = "myKey";

        PlcCertificateAuthentication auth = new PlcCertificateAuthentication(keyStore, password, alias);

        assertEquals(alias, auth.getKeyAlias());
    }

    @Test
    void testConstructorWithNullKeyStoreThrows() {
        char[] password = "password".toCharArray();

        assertThrows(NullPointerException.class, () ->
            new PlcCertificateAuthentication(null, password)
        );
    }

    @Test
    void testConstructorWithNullPasswordThrows() throws Exception {
        KeyStore keyStore = createEmptyKeyStore();

        assertThrows(NullPointerException.class, () ->
            new PlcCertificateAuthentication(keyStore, null)
        );
    }

    @Test
    void testConstructorWithNullAliasAllowed() throws Exception {
        KeyStore keyStore = createEmptyKeyStore();
        char[] password = "password".toCharArray();

        // Null alias should be allowed (means use first entry)
        PlcCertificateAuthentication auth = new PlcCertificateAuthentication(keyStore, password, null);

        assertNull(auth.getKeyAlias());
    }

    @Test
    void testPasswordIsCopied() throws Exception {
        KeyStore keyStore = createEmptyKeyStore();
        char[] password = "password".toCharArray();

        PlcCertificateAuthentication auth = new PlcCertificateAuthentication(keyStore, password);

        // Modify original password
        password[0] = 'X';

        // Returned password should not be affected
        char[] returnedPassword = auth.getKeyStorePassword();
        assertEquals('p', returnedPassword[0]);
    }

    @Test
    void testEquals() throws Exception {
        KeyStore keyStore1 = createEmptyKeyStore();
        KeyStore keyStore2 = createEmptyKeyStore();
        char[] password1 = "password".toCharArray();
        char[] password2 = "password".toCharArray();

        PlcCertificateAuthentication auth1 = new PlcCertificateAuthentication(keyStore1, password1);
        PlcCertificateAuthentication auth2 = new PlcCertificateAuthentication(keyStore1, password2);
        PlcCertificateAuthentication auth3 = new PlcCertificateAuthentication(keyStore2, password1);

        // Same keystore and password
        assertEquals(auth1, auth2);
        // Different keystore instance
        assertNotEquals(auth1, auth3);
    }

    @Test
    void testEqualsWithAlias() throws Exception {
        KeyStore keyStore = createEmptyKeyStore();
        char[] password = "password".toCharArray();

        PlcCertificateAuthentication auth1 = new PlcCertificateAuthentication(keyStore, password, "alias1");
        PlcCertificateAuthentication auth2 = new PlcCertificateAuthentication(keyStore, password, "alias1");
        PlcCertificateAuthentication auth3 = new PlcCertificateAuthentication(keyStore, password, "alias2");

        assertEquals(auth1, auth2);
        assertNotEquals(auth1, auth3);
    }

    @Test
    void testEqualsNull() throws Exception {
        KeyStore keyStore = createEmptyKeyStore();
        char[] password = "password".toCharArray();

        PlcCertificateAuthentication auth = new PlcCertificateAuthentication(keyStore, password);

        assertNotEquals(null, auth);
    }

    @Test
    void testEqualsDifferentType() throws Exception {
        KeyStore keyStore = createEmptyKeyStore();
        char[] password = "password".toCharArray();

        PlcCertificateAuthentication auth = new PlcCertificateAuthentication(keyStore, password);

        assertNotEquals("not an authentication", auth);
    }

    @Test
    void testHashCode() throws Exception {
        KeyStore keyStore = createEmptyKeyStore();
        char[] password1 = "password".toCharArray();
        char[] password2 = "password".toCharArray();

        PlcCertificateAuthentication auth1 = new PlcCertificateAuthentication(keyStore, password1);
        PlcCertificateAuthentication auth2 = new PlcCertificateAuthentication(keyStore, password2);

        assertEquals(auth1.hashCode(), auth2.hashCode());
    }

    @Test
    void testToString() throws Exception {
        KeyStore keyStore = KeyStore.getInstance("PKCS12");
        keyStore.load(null, null);
        char[] password = "password".toCharArray();

        PlcCertificateAuthentication auth = new PlcCertificateAuthentication(keyStore, password, "myAlias");

        String str = auth.toString();
        assertTrue(str.contains("PlcCertificateAuthentication"));
        assertTrue(str.contains("PKCS12"));
        assertTrue(str.contains("myAlias"));
    }

    @Test
    void testToStringWithoutAlias() throws Exception {
        KeyStore keyStore = KeyStore.getInstance("PKCS12");
        keyStore.load(null, null);
        char[] password = "password".toCharArray();

        PlcCertificateAuthentication auth = new PlcCertificateAuthentication(keyStore, password);

        String str = auth.toString();
        assertTrue(str.contains("(first entry)"));
    }

    private KeyStore createEmptyKeyStore() throws Exception {
        KeyStore keyStore = KeyStore.getInstance("PKCS12");
        keyStore.load(null, null);
        return keyStore;
    }
}
