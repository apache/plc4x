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

package org.apache.plc4x.java.spi.fields.fields.reader;

import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.fields.data.reader.DataReader;
import org.apache.plc4x.java.spi.fields.fields.FieldCommons;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Objects;

public class FieldReaderConst<T> implements FieldCommons {

    private static final Logger LOGGER = LoggerFactory.getLogger(FieldReaderConst.class);

    public T readConstField(DataReader<T> dataReader, T expectedValue, WithOption... options) throws BufferException {
        LOGGER.debug("reading field {}", getName(options));
        T constValue = dataReader.read(options);
        if (!Objects.equals(constValue, expectedValue)) {
            throw new BufferException("Actual value " + constValue + " for field " + getName(options) + " doesn't match expected " + expectedValue + ". Bit position: " + dataReader.getPositionInBits() + " Byte position: " + (dataReader.getPositionInBits() / 8));
        }
        return constValue;
    }

}
