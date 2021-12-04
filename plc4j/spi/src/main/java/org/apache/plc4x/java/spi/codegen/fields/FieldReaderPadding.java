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
import org.apache.plc4x.java.spi.codegen.io.DataReader;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.WithReaderArgs;
import org.apache.plc4x.java.spi.generation.WithReaderWriterArgs;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class FieldReaderPadding<T> implements FieldCommons {

    public void readPaddingField(DataReader<T> dataReader, int timesPadding, WithReaderArgs... readerArgs) {
        dataReader.pullContext("padding", WithReaderWriterArgs.WithRenderAsList(true));
        while (timesPadding-- > 0) {
            // Just read the padding data and ignore it
            try {
                dataReader.read("value", readerArgs);
            } catch (Exception e) {
                // Ignore ...
                // This could simply be that we're out of data to read for padding.
                // In protocols like the S7 protocol, this can happen if this is the
                // last field item, then the packet might end here.
            }
        }
        dataReader.closeContext("padding", WithReaderWriterArgs.WithRenderAsList(true));
    }

}
