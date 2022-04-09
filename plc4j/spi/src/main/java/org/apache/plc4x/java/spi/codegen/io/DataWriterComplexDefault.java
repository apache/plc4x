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
package org.apache.plc4x.java.spi.codegen.io;

import org.apache.commons.lang3.StringUtils;
import org.apache.plc4x.java.spi.generation.*;

public class DataWriterComplexDefault<T extends Message> implements DataWriterComplex<T> {

    protected final WriteBuffer writeBuffer;

    public DataWriterComplexDefault(WriteBuffer writeBuffer) {
        this.writeBuffer = writeBuffer;
    }

    @Override
    public ByteOrder getByteOrder() {
        return writeBuffer.getByteOrder();
    }

    @Override
    public void setByteOrder(ByteOrder byteOrder) {
        writeBuffer.setByteOrder(byteOrder);
    }

    @Override
    public void pushContext(String logicalName, WithWriterArgs... writerArgs) {
        writeBuffer.pushContext(logicalName, writerArgs);
    }

    @Override
    public void popContext(String logicalName, WithWriterArgs... writerArgs) {
        writeBuffer.popContext(logicalName, writerArgs);
    }

    @Override
    public void write(String logicalName, T value, WithWriterArgs... writerArgs) throws SerializationException {
        boolean hasLogicalName = StringUtils.isNotBlank(logicalName);
        if (hasLogicalName) {
            writeBuffer.pushContext(logicalName);
        }
        value.serialize(writeBuffer);
        if (hasLogicalName) {
            writeBuffer.popContext(logicalName);
        }
    }

}
