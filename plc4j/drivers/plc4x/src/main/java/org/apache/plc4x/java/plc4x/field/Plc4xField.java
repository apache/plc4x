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
package org.apache.plc4x.java.plc4x.field;

import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.apache.plc4x.java.spi.utils.Serializable;

public class Plc4xField implements PlcField, Serializable {

    private final String address;

    public Plc4xField(String address) {
        this.address = address;
    }

    public String getAddress() {
        return address;
    }

    @Override
    public String getPlcDataType() {
        return PlcField.super.getPlcDataType();
    }

    @Override
    public int getNumberOfElements() {
        return PlcField.super.getNumberOfElements();
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws SerializationException {
        writeBuffer.pushContext(getClass().getSimpleName());

        writeBuffer.writeUnsignedInt("addressLength", 16, address.length());
        writeBuffer.writeString("address",address.length() * 8, "UTF-8", address);

        writeBuffer.popContext(getClass().getSimpleName());
    }

}
