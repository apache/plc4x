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
package org.apache.plc4x.java.can.generic.protocol;

import java.util.Collections;
import java.util.LinkedHashMap;
import java.util.Map;
import org.apache.plc4x.java.can.generic.field.GenericCANField;
import org.apache.plc4x.java.spi.messages.PlcSubscriber;
import org.apache.plc4x.java.spi.model.DefaultPlcSubscriptionHandle;

public class GenericCANSubscriptionHandle extends DefaultPlcSubscriptionHandle {

    private final Map<String, GenericCANField> fields = new LinkedHashMap<>();
    private final Integer nodeId;

    public GenericCANSubscriptionHandle(PlcSubscriber subscriber, Integer nodeId) {
        super(subscriber);
        this.nodeId = nodeId;
    }

    public boolean matches(int identifier) {
        return nodeId == identifier;
    }

    public String toString() {
        return "GenericCANSubscriptionHandle [node=" + nodeId + " " + intAndHex(nodeId) + "]";
    }

    public void add(String name, GenericCANField field) {
        fields.put(name, field);
    }

    private static String intAndHex(int val) {
        return val + "(0x" + Integer.toHexString(val) + ")";
    }

    public Map<String, GenericCANField> getFields() {
        return Collections.unmodifiableMap(fields);
    }
}

