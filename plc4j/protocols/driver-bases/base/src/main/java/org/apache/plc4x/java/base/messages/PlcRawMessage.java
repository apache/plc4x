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
package org.apache.plc4x.java.base.messages;

import io.netty.buffer.ByteBuf;

public class PlcRawMessage implements PlcProtocolMessage {

    private final ByteBuf userData;
    private final PlcProtocolMessage parent;

    public PlcRawMessage(ByteBuf userData) {
        this(userData, null);
    }

    public PlcRawMessage(ByteBuf userData, PlcProtocolMessage parent) {
        this.userData = userData;
        this.parent = parent;
    }

    public ByteBuf getUserData() {
        return userData;
    }

    @Override
    public PlcProtocolMessage getParent() {
        return parent;
    }

}
