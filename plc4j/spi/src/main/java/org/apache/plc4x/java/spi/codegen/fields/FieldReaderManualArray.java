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

import org.apache.commons.lang3.ArrayUtils;
import org.apache.commons.lang3.NotImplementedException;
import org.apache.plc4x.java.spi.codegen.FieldCommons;
import org.apache.plc4x.java.spi.codegen.io.DataReader;
import org.apache.plc4x.java.spi.codegen.io.ParseSupplier;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.WithReaderArgs;
import org.apache.plc4x.java.spi.generation.WithReaderWriterArgs;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.function.Function;
import java.util.function.Supplier;

public class FieldReaderManualArray<T> implements FieldCommons {

    private static final Logger LOGGER = LoggerFactory.getLogger(FieldReaderManualArray.class);

    public byte[] readManualByteArrayField(String logicalName, ReadBuffer readBuffer, Function<byte[],Boolean> termination, ParseSupplier<T> parse, WithReaderArgs... readerArgs) throws ParseException {
        LOGGER.debug("reading field {}", logicalName);
        // Ensure we have the render as list argument present
        readerArgs = ArrayUtils.add(readerArgs, WithReaderWriterArgs.WithRenderAsList(true));
        readBuffer.pullContext(logicalName, readerArgs);
        List<Byte> result = new ArrayList<>();
        while (!termination.apply(ArrayUtils.toPrimitive(result.toArray(new Byte[0])))) {
            //TODO: maybe switch to iterator here
            Byte element = (Byte) parse.get();
            LOGGER.debug("Adding element {}", element);
            result.add(element);
        }
        readBuffer.closeContext(logicalName, readerArgs);
        LOGGER.debug("done reading field {}", logicalName);
        return ArrayUtils.toPrimitive(result.toArray(new Byte[0]));
    }

    public List<T> readManualArrayField(String logicalName, ReadBuffer readBuffer, Function<List<T>,Boolean> termination, ParseSupplier<T> parse, WithReaderArgs... readerArgs) throws ParseException {
        LOGGER.debug("reading field {}", logicalName);
        // Ensure we have the render as list argument present
        readerArgs = ArrayUtils.add(readerArgs, WithReaderWriterArgs.WithRenderAsList(true));
        readBuffer.pullContext(logicalName, readerArgs);
        List<T> result = new ArrayList<>();
        while (!termination.apply(Collections.unmodifiableList(result))) {
            //TODO: maybe switch to iterator here
            T element = parse.get();
            LOGGER.debug("Adding element {}", element);
            result.add(element);
        }
        readBuffer.closeContext(logicalName, readerArgs);
        LOGGER.debug("done reading field {}", logicalName);
        return result;
    }

}
