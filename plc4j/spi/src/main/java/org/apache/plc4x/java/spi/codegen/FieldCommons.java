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
package org.apache.plc4x.java.spi.codegen;

import org.apache.plc4x.java.spi.codegen.io.ByteOrderAware;
import org.apache.plc4x.java.spi.generation.*;

import java.util.Objects;
import java.util.Optional;
import java.util.stream.Stream;

public interface FieldCommons {

    default Optional<ByteOrder> extractByteOder(WithReaderArgs... readerArgs) {
        return extractByteOder(Stream.of(readerArgs).map(WithReaderWriterArgs.class::cast).toArray(WithReaderWriterArgs[]::new));
    }

    default Optional<ByteOrder> extractByteOder(WithWriterArgs... writerArgs) {
        return extractByteOder(Stream.of(writerArgs).map(WithReaderWriterArgs.class::cast).toArray(WithReaderWriterArgs[]::new));
    }

    default Optional<ByteOrder> extractByteOder(WithReaderWriterArgs... readerWriterArgs) {
        for (WithReaderWriterArgs arg : readerWriterArgs) {
            if (arg instanceof withOptionByteOrder) {
                return Optional.of(((withOptionByteOrder) arg).byteOrder());
            }
        }
        return Optional.empty();
    }

    default <T> T switchParseByteOrderIfNecessary(RunParseWrapped<T> runnable, ByteOrderAware byteOrderAware, ByteOrder wantedByteOrder) throws ParseException {
        Objects.requireNonNull(runnable);
        Objects.requireNonNull(byteOrderAware);
        ByteOrder currentByteOrder = byteOrderAware.getByteOrder();
        if (wantedByteOrder == null || currentByteOrder == wantedByteOrder) {
            return runnable.run();
        }
        try {
            byteOrderAware.setByteOrder(wantedByteOrder);
            return runnable.run();
        } finally {
            byteOrderAware.setByteOrder(currentByteOrder);
        }
    }

    default void switchSerializeByteOrderIfNecessary(RunSerializeWrapped runnable, ByteOrderAware byteOrderAware, ByteOrder wantedByteOrder) throws SerializationException {
        Objects.requireNonNull(runnable);
        Objects.requireNonNull(byteOrderAware);
        ByteOrder currentByteOrder = byteOrderAware.getByteOrder();
        if (wantedByteOrder == null || currentByteOrder == wantedByteOrder) {
            runnable.run();
        } else {
            try {
                byteOrderAware.setByteOrder(wantedByteOrder);
                runnable.run();
            } finally {
                byteOrderAware.setByteOrder(currentByteOrder);
            }
        }
    }

    @FunctionalInterface
    interface RunParseWrapped<T> {
        T run() throws ParseException;
    }

    @FunctionalInterface
    interface RunSerializeWrapped {
        void run() throws SerializationException;
    }

    @FunctionalInterface
    interface ConsumeSerializeWrapped<T> {
        void consume(T value) throws SerializationException;
    }

}