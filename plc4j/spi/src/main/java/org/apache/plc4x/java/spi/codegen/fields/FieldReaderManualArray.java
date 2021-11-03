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
import org.apache.plc4x.java.spi.codegen.io.DataReader;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.WithReaderArgs;
import org.apache.plc4x.java.spi.generation.WithReaderWriterArgs;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.ArrayList;
import java.util.List;
import java.util.function.Supplier;

public class FieldReaderManualArray<T> implements FieldReader<T> {

    private static final Logger LOGGER = LoggerFactory.getLogger(FieldReaderManualArray.class);

    @Override
    public T readField(String logicalName, DataReader<T> dataReader, WithReaderArgs... readerArgs) throws ParseException {
        throw new NotImplementedException();
    }

    public byte[] readManualByteArrayField(String logicalName, ReadBuffer readBuffer, Supplier<Boolean> termination, Supplier<T> parse, WithReaderArgs... readerArgs) throws ParseException {
        // Ensure we have the render as list argument present
        readerArgs = ArrayUtils.add(readerArgs, WithReaderWriterArgs.WithRenderAsList(true));
        readBuffer.pullContext(logicalName, readerArgs);
        List<Byte> result = new ArrayList<>();
        while (!termination.get()) {
            result.add((Byte) parse.get());
        }
        readBuffer.closeContext(logicalName, readerArgs);
        return ArrayUtils.toPrimitive(result.toArray(new Byte[0]));
    }

    public List<T> readManualArrayField(String logicalName, ReadBuffer readBuffer, Supplier<Boolean> termination, Supplier<T> parse, WithReaderArgs... readerArgs) throws ParseException {
        // Ensure we have the render as list argument present
        readerArgs = ArrayUtils.add(readerArgs, WithReaderWriterArgs.WithRenderAsList(true));
        readBuffer.pullContext(logicalName, readerArgs);
        List<T> result = new ArrayList<>();
        while (!termination.get()) {
            result.add(parse.get());
        }
        readBuffer.closeContext(logicalName, readerArgs);
        return result;
    }

}
