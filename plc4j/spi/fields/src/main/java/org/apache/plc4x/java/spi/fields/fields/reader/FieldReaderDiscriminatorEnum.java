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

public class FieldReaderDiscriminatorEnum<T> implements FieldCommons {

    private static final Logger LOGGER = LoggerFactory.getLogger(FieldReaderDiscriminatorEnum.class);

    public T readDiscriminatorEnumField(DataReader<T> dataReader, WithOption... options) throws BufferException {
        LOGGER.debug("reading field {}", getName(options));
        dataReader.pushContext(WithOption.AddOptions(options, WithOption.WithRenderAsList(true)));
        // The serializer names the inner element after the enum type (see FieldWriterDiscriminatorEnum),
        // which isn't statically known here - read with the wildcard name so name-validating
        // buffers (XML) accept it.
        T result = dataReader.read(WithOption.AddOptions(options, WithOption.WithName("*")));
        dataReader.popContext(WithOption.AddOptions(options, WithOption.WithRenderAsList(true)));
        return result;
    }

}
