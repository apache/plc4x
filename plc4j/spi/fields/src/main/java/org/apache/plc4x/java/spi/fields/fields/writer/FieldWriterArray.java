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

package org.apache.plc4x.java.spi.fields.fields.writer;

import org.apache.plc4x.java.spi.buffers.api.Message;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.WriteBuffer;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.fields.data.writer.DataWriter;
import org.apache.plc4x.java.spi.fields.fields.FieldCommons;
import org.apache.plc4x.java.spi.fields.utils.ThreadLocalHelper;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.List;

public class FieldWriterArray<T> implements FieldCommons {

    private static final Logger LOGGER = LoggerFactory.getLogger(FieldWriterArray.class);

    public void writeByteArrayField(byte[] values, DataWriter<byte[]> dataWriter, WithOption... options) throws BufferException {
        LOGGER.debug("write field {}", getName(options));
        if (values != null) {
            dataWriter.write(values, options);
        }
    }

    public void writeSimpleTypeArrayField(List<T> values, DataWriter<T> dataWriter, WithOption... options) throws BufferException {
        LOGGER.debug("write field {}", getName(options));
        if (values != null) {
            dataWriter.pushContext(WithOption.AddOptions(options, WithOption.WithRenderAsList(true)));
            for (int curItem = 0; curItem < values.size(); curItem++) {
                T value = values.get(curItem);
                ThreadLocalHelper.lastItemThreadLocal.set(curItem == values.size() - 1);
                dataWriter.write(value, WithOption.UpdateOptions(options, WithOption.WithName("value")));
            }
            dataWriter.popContext(WithOption.AddOptions(options, WithOption.WithRenderAsList(true)));
        }
    }

    public void writeComplexTypeArrayField(List<? extends Message> values, WriteBuffer writeBuffer, WithOption... options) throws BufferException {
        LOGGER.debug("write field {}", getName(options));
        if (values != null) {
            try {
                writeBuffer.pushContext(WithOption.AddOptions(options, WithOption.WithRenderAsList(true)));
                for (int curItem = 0; curItem < values.size(); curItem++) {
                    Message value = values.get(curItem);
                    ThreadLocalHelper.lastItemThreadLocal.set(curItem == values.size() - 1);
                    value.serialize(writeBuffer);
                }
                writeBuffer.popContext(WithOption.AddOptions(options, WithOption.WithRenderAsList(true)));
            } catch (BufferException e) {
                throw new BufferException("Error writing complex type array field", e);
            }
        }
    }

}
