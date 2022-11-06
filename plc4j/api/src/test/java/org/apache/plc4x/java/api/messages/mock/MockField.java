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
package org.apache.plc4x.java.api.messages.mock;

import org.apache.plc4x.java.api.model.ArrayInfo;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.types.PlcValueType;

import java.util.Collections;
import java.util.List;

public class MockField implements PlcField {

    private final String field;

    public MockField(String field) {
        this.field = field;
    }

    @Override
    public String getAddressString() {
        return field;
    }

    @Override
    public PlcValueType getPlcValueType() {
        return PlcValueType.INT;
    }

    @Override
    public List<ArrayInfo> getArrayInfo() {
        return Collections.emptyList();
    }

    public String toString() {
        return "mock field: " + field;
    }

    @Override
    public boolean equals(Object o) {
        return o != null
            && o instanceof MockField
            && ((MockField) o).field.equals(this.field);
    }

    @Override
    public int hashCode() {
        return field.hashCode();
    }

}
