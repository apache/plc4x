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

package org.apache.plc4x.java.spi.connection;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.spi.GeneratedDriverByteToMessageCodec;
import org.apache.plc4x.java.spi.generation.Message;
import org.apache.plc4x.java.spi.generation.MessageIO;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.WriteBuffer;

import java.util.function.Consumer;
import java.util.function.Function;

class GeneratedProtocolMessageCodec<BASE_PAKET_CLASS extends Message> extends GeneratedDriverByteToMessageCodec<BASE_PAKET_CLASS> {

    private final Function<ByteBuf, Integer> packetSizeEstimator;
    private final Consumer<ByteBuf> corruptPackageRemover;

    public GeneratedProtocolMessageCodec(Class<BASE_PAKET_CLASS> basePaketClass,
                                         Parser<BASE_PAKET_CLASS> parser,
                                         Serializer<BASE_PAKET_CLASS> serializer) {
        this(basePaketClass, parser, serializer, null, null);
    }

    public GeneratedProtocolMessageCodec(
        Class<BASE_PAKET_CLASS> basePaketClass,
        Parser<BASE_PAKET_CLASS> parser,
        Serializer<BASE_PAKET_CLASS> serializer,
        Function<ByteBuf, Integer> packetSizeEstimator,
        Consumer<ByteBuf> corruptPackageRemover) {
        super(new MessageIO<BASE_PAKET_CLASS, BASE_PAKET_CLASS>() {
            @Override
            public BASE_PAKET_CLASS parse(ReadBuffer io) throws ParseException {
                // return TPKTPacketIO.parse(io);
                return parser.parse(io);
            }

            @Override
            public void serialize(WriteBuffer io, BASE_PAKET_CLASS value) throws ParseException {
                serializer.serialize(io, value);
            }

        }, basePaketClass);
        this.packetSizeEstimator = packetSizeEstimator;
        this.corruptPackageRemover = corruptPackageRemover;
    }

    @Override protected int getPacketSize(ByteBuf byteBuf) {
        if (this.packetSizeEstimator == null) {
            return -1;
        }
        return packetSizeEstimator.apply(byteBuf);
    }

    @Override protected void removeRestOfCorruptPackage(ByteBuf byteBuf) {
        if (this.corruptPackageRemover == null) {
            throw new IllegalStateException("This Implementation does not support Corrupt Package Removal!");
        }
        this.corruptPackageRemover.accept(byteBuf);
    }
}
