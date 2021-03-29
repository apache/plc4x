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
package org.apache.plc4x.java.profinet.dcp.protocol;

import org.apache.plc4x.java.profinet.dcp.field.ProfinetDcpField;
import org.apache.plc4x.java.profinet.dcp.readwrite.ProfinetFrame;
import org.apache.plc4x.java.spi.messages.PlcSubscriber;
import org.apache.plc4x.java.spi.model.DefaultPlcSubscriptionHandle;

public class ProfinetDCPSubscriptionHandle extends DefaultPlcSubscriptionHandle {

    private final String name;
    private final ProfinetDcpField field;

    public ProfinetDCPSubscriptionHandle(PlcSubscriber subscriber, String name, ProfinetDcpField field) {
        super(subscriber);
        this.name = name;
        this.field = field;
    }

    public String getName() {
        return name;
    }

    boolean matches(ProfinetFrame pdu) {
        // TODO implement matching logic
        return false;
    }

    public ProfinetDcpField getField() {
        return field;
    }

    public String toString() {
        return "ProfinetDCPSubscriptionHandle [service=" + field + "]";
    }

}
