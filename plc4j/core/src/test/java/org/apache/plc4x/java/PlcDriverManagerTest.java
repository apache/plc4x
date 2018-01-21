/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/
package org.apache.plc4x.java;

import org.apache.plc4x.java.api.authentication.PlcUsernamePasswordAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.mock.MockConnection;
import org.assertj.core.api.SoftAssertions;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import java.io.File;
import java.net.MalformedURLException;
import java.net.URL;
import java.net.URLClassLoader;

import static org.assertj.core.api.Assertions.*;

class PlcDriverManagerTest {

    /**
     * Tries to get the mock plc driver which is part of this testsuite.
     */
    @Test
    @Tag("fast")
    void getExistingDriverTest() throws PlcException {
        MockConnection mockConnection = (MockConnection) new PlcDriverManager().getConnection("mock://some-cool-url");
        SoftAssertions softAssertions = new SoftAssertions();
        softAssertions.assertThat(mockConnection.getAuthentication())
            .as("check authentication object")
            .isNull();
        softAssertions.assertThat(mockConnection.isConnected())
            .as("check connection state")
            .isTrue();
        softAssertions.assertThat(mockConnection.isClosed())
            .as("check closed state")
            .isFalse();
        softAssertions.assertAll();
    }

    /**
     * Tries to get the mock plc driver with authentication which is part of this testsuite.
     */
    @Test
    @Tag("fast")
    void getExistingDriverWithAuthenticationTest() throws PlcException {
        PlcUsernamePasswordAuthentication authentication =
            new PlcUsernamePasswordAuthentication("user", "pass");
        MockConnection mockConnection = (MockConnection) new PlcDriverManager().getConnection("mock://some-cool-url", authentication);

        SoftAssertions softAssertions = new SoftAssertions();
        softAssertions.assertThat(mockConnection.getAuthentication())
            .as("check authentication object")
            .isNotNull();
        softAssertions.assertThat(mockConnection.getAuthentication())
            .as("check authentication object type")
            .isInstanceOf(PlcUsernamePasswordAuthentication.class);
        softAssertions.assertThat(mockConnection.isConnected())
            .as("check connection state")
            .isTrue();
        softAssertions.assertThat(mockConnection.isClosed())
            .as("check closed state")
            .isFalse();
        softAssertions.assertAll();
    }

    /**
     * In this test case a driver is requested which is not registered with the {@link PlcDriverManager}.
     */
    @Test
    @Tag("fast")
    void getNotExistingDriverTest() throws PlcException {
        assertThatThrownBy(() -> new PlcDriverManager().getConnection("non-existing-protocol://some-cool-url"))
            .as("check rejection of invalid protocol")
            .isInstanceOf(PlcConnectionException.class);
    }

    /**
     * In this test case a driver is requested which is not registered with the {@link PlcDriverManager}.
     */
    @Test
    @Tag("fast")
    void getInvalidUriTest() throws PlcException {
        assertThatThrownBy(() -> new PlcDriverManager().getConnection("The quick brown fox jumps over the lazy dog"))
            .as("check rejection of invalid uri")
            .isInstanceOf(PlcConnectionException.class);
    }

    /**
     * In this test the {@link PlcDriverManager} will be configured with a service list that
     * contains multiple implementation instances of the same protocol. This should result in
     * an error.
     */
    @Test
    @Tag("fast")
    void getDuplicateDriver() throws PlcException, MalformedURLException {
        // Save and replace the context classloader as we need to force the ServiceLoader to
        // use a different service file.
        ClassLoader originalClassloader = Thread.currentThread().getContextClassLoader();
        URL[] urls = new URL[1];
        urls[0] = new File("src/test/resources/test").toURI().toURL();
        ClassLoader fakeClassLoader = new URLClassLoader(urls, originalClassloader);

        assertThatThrownBy(() -> new PlcDriverManager(fakeClassLoader).getConnection("mock://some-cool-url"))
            .as("check detection of duplicated driver detection")
            .isInstanceOf(IllegalStateException.class);
    }

}
