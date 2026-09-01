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
package org.apache.plc4x.java.opcua.context;

import org.apache.plc4x.java.opcua.readwrite.ApplicationDescription;
import org.apache.plc4x.java.opcua.readwrite.EndpointDescription;
import org.apache.plc4x.java.opcua.readwrite.MessageSecurityMode;
import org.apache.plc4x.java.opcua.readwrite.PascalByteString;
import org.apache.plc4x.java.opcua.readwrite.PascalString;
import org.apache.plc4x.java.opcua.readwrite.UserTokenPolicy;
import org.apache.plc4x.java.opcua.readwrite.UserTokenType;
import org.apache.plc4x.java.opcua.security.SecurityPolicy;
import org.junit.jupiter.api.Test;
import org.mockito.Mockito;

import java.util.List;
import java.util.Map.Entry;

import static java.util.Map.entry;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * A server commonly offers several endpoints at once, including a wide-open one. Which of them is
 * chosen decides what protection the connection actually gets, whatever the configuration asked
 * for.
 */
class EndpointSelectionTest {

    private static final PascalString NULL_STRING = new PascalString("");

    private static UserTokenPolicy tokenPolicy(String id, UserTokenType type, SecurityPolicy policy) {
        return new UserTokenPolicy(new PascalString(id), type, NULL_STRING, NULL_STRING,
            policy == null ? NULL_STRING : new PascalString(policy.getSecurityPolicyUri()));
    }

    private static Entry<EndpointDescription, UserTokenPolicy> endpoint(
        SecurityPolicy channelPolicy, MessageSecurityMode mode, int securityLevel,
        UserTokenPolicy userTokenPolicy) {
        EndpointDescription description = new EndpointDescription(
            new PascalString("opc.tcp://localhost:4840"),
            Mockito.mock(ApplicationDescription.class),
            new PascalByteString(0, new byte[0]),
            mode,
            new PascalString(channelPolicy.getSecurityPolicyUri()),
            List.of(userTokenPolicy),
            NULL_STRING,
            (short) securityLevel);
        return entry(description, userTokenPolicy);
    }

    @Test
    void theStrongestEndpointWinsNotTheWeakest() {
        Entry<EndpointDescription, UserTokenPolicy> wideOpen = endpoint(
            SecurityPolicy.NONE, MessageSecurityMode.messageSecurityModeNone, 0,
            tokenPolicy("anon", UserTokenType.userTokenTypeAnonymous, SecurityPolicy.NONE));
        Entry<EndpointDescription, UserTokenPolicy> protected_ = endpoint(
            SecurityPolicy.Basic256Sha256, MessageSecurityMode.messageSecurityModeSignAndEncrypt, 30,
            tokenPolicy("user", UserTokenType.userTokenTypeUserName, SecurityPolicy.Basic256Sha256));

        // Whichever order the server listed them in.
        assertEquals(protected_, SecureChannel.strongestOf(List.of(wideOpen, protected_)));
        assertEquals(protected_, SecureChannel.strongestOf(List.of(protected_, wideOpen)));
    }

    @Test
    void betweenEquallyStrongEndpointsTheOneProtectingTheTokenWins() {
        Entry<EndpointDescription, UserTokenPolicy> tokenInTheClear = endpoint(
            SecurityPolicy.Basic256Sha256, MessageSecurityMode.messageSecurityModeSignAndEncrypt, 20,
            tokenPolicy("user-plain", UserTokenType.userTokenTypeUserName, SecurityPolicy.NONE));
        Entry<EndpointDescription, UserTokenPolicy> tokenProtected = endpoint(
            SecurityPolicy.Basic256Sha256, MessageSecurityMode.messageSecurityModeSignAndEncrypt, 20,
            tokenPolicy("user-enc", UserTokenType.userTokenTypeUserName, SecurityPolicy.Basic256Sha256));

        assertEquals(tokenProtected, SecureChannel.strongestOf(List.of(tokenInTheClear, tokenProtected)));
        assertEquals(tokenProtected, SecureChannel.strongestOf(List.of(tokenProtected, tokenInTheClear)));
    }

    @Test
    void nothingMatchingIsNotSomethingMatching() {
        assertNull(SecureChannel.strongestOf(List.of()));
        assertNull(SecureChannel.strongestOf(null));
    }

    @Test
    void aTokenPolicyOfNoneDoesNotProtectTheToken() {
        assertFalse(SecureChannel.protectsUserToken(endpoint(
            SecurityPolicy.Basic256Sha256, MessageSecurityMode.messageSecurityModeSignAndEncrypt, 10,
            tokenPolicy("plain", UserTokenType.userTokenTypeUserName, SecurityPolicy.NONE))));
    }

    @Test
    void aTokenPolicyThatEncryptsProtectsTheToken() {
        assertTrue(SecureChannel.protectsUserToken(endpoint(
            SecurityPolicy.NONE, MessageSecurityMode.messageSecurityModeNone, 10,
            tokenPolicy("enc", UserTokenType.userTokenTypeUserName, SecurityPolicy.Basic256Sha256))));
    }

    @Test
    void sayingNothingLeavesItToTheChannel() {
        // A token policy naming no policy of its own is governed by the channel's.
        assertTrue(SecureChannel.protectsUserToken(endpoint(
            SecurityPolicy.Basic256Sha256, MessageSecurityMode.messageSecurityModeSignAndEncrypt, 10,
            tokenPolicy("inherit", UserTokenType.userTokenTypeUserName, null))));
        assertFalse(SecureChannel.protectsUserToken(endpoint(
            SecurityPolicy.NONE, MessageSecurityMode.messageSecurityModeNone, 10,
            tokenPolicy("inherit", UserTokenType.userTokenTypeUserName, null))));
    }
}
