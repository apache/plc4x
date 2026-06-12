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

import org.apache.plc4x.java.opcua.readwrite.OpcuaProtocolLimits;
import org.apache.plc4x.java.opcua.security.MessageSecurity;
import org.apache.plc4x.java.opcua.security.SecurityPolicy;

import java.security.cert.X509Certificate;

/**
 * Read-only window on the secure channel state that the encryption handlers
 * need. Conversation implements this; encryption handlers see only the slice
 * they actually use, which keeps the dependency between Conversation and the
 * encryption layer unidirectional.
 */
public interface SecureChannelState {

    SecurityPolicy getSecurityPolicy();

    MessageSecurity getMessageSecurity();

    byte[] getLocalNonce();

    byte[] getRemoteNonce();

    X509Certificate getLocalCertificate();

    X509Certificate getRemoteCertificate();

    boolean isSymmetricEncryptionEnabled();

    boolean isSymmetricSigningEnabled();

    OpcuaProtocolLimits getLimits();

}
