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

import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WithWriterArgs;
import org.apache.plc4x.java.spi.generation.WriteBuffer;

import java.util.Objects;

public class DataWriterSimpleByteArray extends DataWriterSimpleBase<byte[]> {

    public DataWriterSimpleByteArray(WriteBuffer writeBuffer, int bitLength) {
        super(writeBuffer, bitLength);
    }

    @Override
    public byte[] write(String logicalName, byte[] value, WithWriterArgs... writerArgs) throws SerializationException {
        if (bitLength != 8) {
            throw new SerializationException("ByteArray fields only support bitLength of 8");
        }
        // TODO: Get a WirReaderArgs parameter for the number of bytes ...
        writeBuffer.writeByteArray(logicalName, value, writerArgs);
        return value;
    }

}
