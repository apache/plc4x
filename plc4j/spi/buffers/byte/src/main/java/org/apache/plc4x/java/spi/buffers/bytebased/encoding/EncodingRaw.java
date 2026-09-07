/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.spi.buffers.bytebased.encoding;

import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferRaw;

import java.math.BigDecimal;
import java.math.BigInteger;

/**
 * Interface for the decoding part of raw encodings.
 */
public interface EncodingRaw extends Encoding {

    byte decodeByte(int numBits, ReadBufferRaw readBuffer) throws BufferException;

    short decodeShort(int numBits, ReadBufferRaw readBuffer) throws BufferException;

    int decodeInt(int numBits, ReadBufferRaw readBuffer) throws BufferException;

    long decodeLong(int numBits, ReadBufferRaw readBuffer) throws BufferException;

    BigInteger decodeBigInteger(int numBits, ReadBufferRaw readBuffer) throws BufferException;

    float decodeFloat(int numBits, ReadBufferRaw readBuffer) throws BufferException;

    double decodeDouble(int numBits, ReadBufferRaw readBuffer) throws BufferException;

    BigDecimal decodeBigDecimal(int numBits, ReadBufferRaw readBuffer) throws BufferException;

    String decodeString(int numBits, ReadBufferRaw readBuffer) throws BufferException;

}
