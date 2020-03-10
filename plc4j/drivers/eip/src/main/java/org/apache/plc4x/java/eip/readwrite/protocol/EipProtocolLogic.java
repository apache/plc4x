/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.eip.readwrite.protocol;

import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.generation.Message;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.time.Duration;

public class EipProtocolLogic extends Plc4xProtocolBase<EipPkt> {

    private static final Logger logger = LoggerFactory.getLogger(EipProtocolLogic.class);
    public static final Duration REQUEST_TIMEOUT = Duration.ofMillis(10000);

    @Override
    public void onConnect(ConversationContext<Message> context) {
        //TODO Send a ENIP Message with Register Session Code '0x0065'
        // , empty Session Handle and Sender Context
        // Then we need to accept the response with the same Code
        // and save the assigned Session Handle
        // PS: Check Status for Success : 0x00000000
    }





    @Override
    public void close(ConversationContext<Message> context) {
        //TODO Send a ENIP Message with Unregister Session Code '0x0066'
        // PS: Check Status for Success : 0x00000000
    }

}
