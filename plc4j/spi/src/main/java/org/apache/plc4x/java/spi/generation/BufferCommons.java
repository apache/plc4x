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
package org.apache.plc4x.java.spi.generation;

import java.util.stream.Stream;

public interface BufferCommons {
    String rwDataTypeKey = "dataType";
    String rwBitLengthKey = "bitLength";
    String rwStringRepresentationKey = "stringRepresentation";
    String rwBitKey = "bit";
    String rwByteKey = "byte";
    String rwUintKey = "uint";
    String rwIntKey = "int";
    String rwFloatKey = "float";
    String rwStringKey = "string";
    String rwEncodingKey = "encoding";
    String rwIsListKey = "isList";

    default String sanitizeLogicalName(String logicalName) {
        if (logicalName.equals("")) {
            return "value";
        }
        return logicalName;
    }

    default boolean isToBeRenderedAsList(WithReaderArgs... readerArgs) {
        return isToBeRenderedAsList(Stream.of(readerArgs).map(WithReaderWriterArgs.class::cast).toArray(WithReaderWriterArgs[]::new));
    }

    default boolean isToBeRenderedAsList(WithWriterArgs... writerArgs) {
        return isToBeRenderedAsList(Stream.of(writerArgs).map(WithReaderWriterArgs.class::cast).toArray(WithReaderWriterArgs[]::new));
    }

    default boolean isToBeRenderedAsList(WithReaderWriterArgs... readerWriterArgs) {
        for (WithReaderWriterArgs arg : readerWriterArgs) {
            if (arg instanceof withRenderAsList) {
                return ((withRenderAsList) arg).renderAsList();
            }
        }
        return false;
    }

    default String extractAdditionalStringRepresentation(WithReaderArgs... readerArgs) {
        return extractAdditionalStringRepresentation(Stream.of(readerArgs).map(WithReaderWriterArgs.class::cast).toArray(WithReaderWriterArgs[]::new));
    }

    default String extractAdditionalStringRepresentation(WithWriterArgs... writerArgs) {
        return extractAdditionalStringRepresentation(Stream.of(writerArgs).map(WithReaderWriterArgs.class::cast).toArray(WithReaderWriterArgs[]::new));
    }

    default String extractAdditionalStringRepresentation(WithReaderWriterArgs... readerWriterArgs) {
        for (WithReaderWriterArgs arg : readerWriterArgs) {
            if (arg instanceof withAdditionalStringRepresentation) {
                return ((withAdditionalStringRepresentation) arg).stringRepresentation();
            }
        }
        return null;
    }
}
