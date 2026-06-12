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

public class FieldReaderPadding<T> implements FieldCommons {

    private static final Logger LOGGER = LoggerFactory.getLogger(FieldReaderPadding.class);

    public void readPaddingField(DataReader<T> dataReader, int timesPadding, WithOption... options) throws BufferException {
        LOGGER.debug("reading field padding");
        dataReader.pushContext(WithOption.AddOptions(options, WithOption.WithRenderAsList(true)));
        while (timesPadding-- > 0) {
            // Just read the padding data and ignore it
            try {
                dataReader.read(options);
            } catch (Exception e) {
                // Ignore ...
                // This could simply be that we're out of data to read for padding.
                // In protocols like the S7 protocol, this can happen if this is the
                // last field item, then the packet might end here.
            }
        }
        dataReader.popContext(WithOption.AddOptions(options, WithOption.WithRenderAsList(true)));
    }

}
