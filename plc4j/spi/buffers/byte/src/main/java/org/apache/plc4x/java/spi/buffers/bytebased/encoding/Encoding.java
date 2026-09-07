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

import java.math.BigDecimal;
import java.math.BigInteger;

/**
 * Both the normal and the raw encoders share the methods for encoding data.
 */
public interface Encoding {

    String getName();

    byte[] encodeByte(int numBits, byte value) throws BufferException;

    byte[] encodeShort(int numBits, short value) throws BufferException;

    byte[] encodeInt(int numBits, int value) throws BufferException;

    byte[] encodeLong(int numBits, long value) throws BufferException;

    byte[] encodeBigInteger(int numBits, BigInteger value) throws BufferException;

    byte[] encodeFloat(int numBits, float value) throws BufferException;

    byte[] encodeDouble(int numBits, double value) throws BufferException;

    byte[] encodeBigDecimal(int numBits, BigDecimal value) throws BufferException;

    byte[] encodeString(int numBits, String value) throws BufferException;

}
