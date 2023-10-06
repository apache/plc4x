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

import nl.jqno.equalsverifier.EqualsVerifier;
import org.hamcrest.Matchers;
import org.junit.jupiter.api.Test;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.core.IsEqual.equalTo;

public class PlcUsernamePasswordAuthenticationTest {

    @Test
    public void authentication() {
        PlcUsernamePasswordAuthentication authentication = new PlcUsernamePasswordAuthentication("user", "password");

        assertThat("Unexpected user name", authentication.getUsername(), equalTo("user"));
        assertThat("Unexpected password", authentication.getPassword(), equalTo("password"));
    }

    /**
     * Usually in a toString method most properties are output.
     * However, the password field should never be output this way or the password could be leaked to a log-file
     * unintentionally.
     */
    @Test
    public void toStringDoesntLeakPassword() {
        PlcUsernamePasswordAuthentication authentication = new PlcUsernamePasswordAuthentication("user", "top-secret");

        assertThat(authentication.toString(), Matchers.not(Matchers.containsString("top-secret")));
    }

    @Test
    public void equalsContract() {
        EqualsVerifier.forClass(PlcUsernamePasswordAuthentication.class).verify();
    }

}