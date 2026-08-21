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

import org.apache.plc4x.java.spi.buffers.api.ReadBuffer;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferUnderflowException;
import org.apache.plc4x.java.spi.fields.data.reader.ParseSupplier;
import org.apache.plc4x.java.spi.fields.fields.FieldCommons;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.function.Function;

public class FieldReaderManualArray<T> implements FieldCommons {

    private static final Logger LOGGER = LoggerFactory.getLogger(FieldReaderManualArray.class);

    public byte[] readManualByteArrayField(ReadBuffer readBuffer, Function<List<Byte>, Boolean> termination, ParseSupplier<T> parse, WithOption... options) throws BufferException {
        LOGGER.debug("reading field {}", getName(options));
        // Ensure we have the render as list argument present
        WithOption[] newOptions = WithOption.AddOptions(options, WithOption.WithRenderAsList(true));
        readBuffer.pushContext(newOptions);
        List<Byte> result = new ArrayList<>();
        while (!termination.apply(result)) {
            //TODO: maybe switch to iterator here
            int positionBefore = readBuffer.getPositionInBits();
            Byte element = (Byte) parse.get();
            LOGGER.debug("Adding element {}", element);
            result.add(element);
            // The loop ends only when the termination sequence turns up, so a parse supplier that
            // reports the end of the data as a value would spin forever on data that never contains
            // it. Termination by count is legitimate and consumes nothing, so only give up once
            // there is also nothing left to read.
            if (readBuffer.getPositionInBits() == positionBefore && readBuffer.getRemainingBits() <= 0) {
                throw new BufferUnderflowException("ran out of data reading " + getName(options) + " after " + result.size() + " items");
            }
        }
        readBuffer.popContext(newOptions);
        LOGGER.debug("done reading field {}", getName(options));

        byte[] resultArray = new byte[result.size()];
        for (int i = 0; i < result.size(); i++) {
            resultArray[i] = result.get(i);
        }
        return resultArray;
    }

    public List<T> readManualArrayField(ReadBuffer readBuffer, Function<List<T>, Boolean> termination, ParseSupplier<T> parse, WithOption... options) throws BufferException {
        LOGGER.debug("reading field {}", getName(options));
        // Ensure we have the render as list argument present
        WithOption[] newOptions = WithOption.AddOptions(options, WithOption.WithRenderAsList(true));
        readBuffer.pushContext(newOptions);
        List<T> result = new ArrayList<>();
        while (!termination.apply(Collections.unmodifiableList(result))) {
            //TODO: maybe switch to iterator here
            int positionBefore = readBuffer.getPositionInBits();
            T element = parse.get();
            LOGGER.debug("Adding element {}", element);
            result.add(element);
            // The loop ends only when the termination sequence turns up, so a parse supplier that
            // reports the end of the data as a value would spin forever on data that never contains
            // it. Termination by count is legitimate and consumes nothing, so only give up once
            // there is also nothing left to read.
            if (readBuffer.getPositionInBits() == positionBefore && readBuffer.getRemainingBits() <= 0) {
                throw new BufferUnderflowException("ran out of data reading " + getName(options) + " after " + result.size() + " items");
            }
        }
        readBuffer.popContext(newOptions);
        LOGGER.debug("done reading field {}", getName(options));
        return result;
    }

}
