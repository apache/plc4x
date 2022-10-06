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
package org.apache.plc4x.java.canopen.protocol;

import org.apache.plc4x.java.canopen.field.CANOpenSubscriptionField;
import org.apache.plc4x.java.canopen.readwrite.CANOpenService;
import org.apache.plc4x.java.spi.messages.PlcSubscriber;
import org.apache.plc4x.java.spi.model.DefaultPlcSubscriptionHandle;

public class CANOpenSubscriptionHandle extends DefaultPlcSubscriptionHandle {
    private final String name;
    private final CANOpenSubscriptionField field;

    public CANOpenSubscriptionHandle(PlcSubscriber subscriber, String name, CANOpenSubscriptionField field) {
        super(subscriber);
        this.name = name;
        this.field = field;
    }

    public boolean matches(CANOpenService service, int identifier) {
        if (field.getService() != service) {
            return false;
        }
        return field.isWildcard() || field.getNodeId() == identifier;
    }

    public String getName() {
        return name;
    }

    public CANOpenSubscriptionField getField() {
        return field;
    }

    public String toString() {
        return "CANopenSubscriptionHandle [service=" + field.getService() + ", node=" + intAndHex(field.getNodeId()) + ", cob=" + intAndHex(field.getService().getMin() + field.getNodeId()) + "]";
    }

    private static String intAndHex(int val) {
        return val + "(0x" + Integer.toHexString(val) + ")";
    }

}
