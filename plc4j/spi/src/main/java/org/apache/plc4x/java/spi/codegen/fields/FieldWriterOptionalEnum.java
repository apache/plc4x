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
package org.apache.plc4x.java.spi.codegen.fields;

import org.apache.plc4x.java.spi.codegen.FieldCommons;
import org.apache.plc4x.java.spi.codegen.io.DataWriter;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WithWriterArgs;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class FieldWriterOptionalEnum<T> implements FieldCommons {

    private static final Logger LOGGER = LoggerFactory.getLogger(FieldWriterOptionalEnum.class);

    public void writeOptionalEnumField(String logicalName, String innerName, T value, DataWriter<T> dataWriter, boolean condition, WithWriterArgs... writerArgs) throws SerializationException {
        LOGGER.debug("write field {}", logicalName);
        dataWriter.pushContext(logicalName);
        if (condition && value != null) {
            switchSerializeByteOrderIfNecessary(() -> dataWriter.write(innerName, value, writerArgs), dataWriter, extractByteOrder(writerArgs).orElse(null));
        } else {
            LOGGER.debug("field {} not written because value is null({}) or condition({}) didn't evaluate to true", logicalName, value != null, condition);
        }
        dataWriter.popContext(logicalName);
    }

}
