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
package org.apache.plc4x.java.spi.codegen.fields;

import org.apache.plc4x.java.spi.codegen.FieldCommons;
import org.apache.plc4x.java.spi.codegen.io.DataReader;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.WithReaderArgs;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Objects;

public class FieldReaderReserved<T> implements FieldCommons {

    private static final Logger LOGGER = LoggerFactory.getLogger(FieldReaderReserved.class);

    public T readReservedField(String logicalName, DataReader<T> dataReader, T referenceValue, WithReaderArgs... readerArgs) throws ParseException {
        LOGGER.debug("reading field {}", logicalName);
        T reserved = dataReader.read(logicalName, readerArgs);
        if (!Objects.equals(reserved, referenceValue)) {
            LOGGER.info("Expected constant value {} but got {} for reserved field.", referenceValue, reserved);
            return reserved;
        }
        return null;
    }

}
