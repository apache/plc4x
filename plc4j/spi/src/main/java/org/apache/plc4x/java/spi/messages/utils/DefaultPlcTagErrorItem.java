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

package org.apache.plc4x.java.spi.messages.utils;

import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.apache.plc4x.java.spi.utils.Serializable;

public class DefaultPlcTagErrorItem<T extends PlcTag> implements PlcTagItem<T>, PlcTagValueItem<T>, Serializable {

    private final PlcResponseCode responseCode;

    public DefaultPlcTagErrorItem(PlcResponseCode responseCode) {
        this.responseCode = responseCode;
    }

    public PlcResponseCode getResponseCode() {
        return responseCode;
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws SerializationException {
        writeBuffer.pushContext("PlcTagErrorItem");
        writeBuffer.writeString("responseCode", responseCode.name().length() * 8, responseCode.name());
        writeBuffer.popContext("PlcTagErrorItem");
    }

}
