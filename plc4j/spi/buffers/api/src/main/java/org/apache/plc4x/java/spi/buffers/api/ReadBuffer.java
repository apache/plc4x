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
package org.apache.plc4x.java.spi.buffers.api;

import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;

import java.math.BigDecimal;
import java.math.BigInteger;

/**
 * Interface for reading data from a buffer with bit-level precision.
 * Provides methods for reading various data types with specified bit lengths.
 */
public interface ReadBuffer extends Buffer {

    /**
     * Reads a single bit from the buffer.
     *
     * @return true if the bit is set, false otherwise
     */
    boolean readBit(WithOption... options) throws BufferException;

    /**
     * Reads a specified number of bits from the buffer.
     *
     * @param numBits the number of bits to read
     * @return a byte array containing the read bits
     */
    byte[] readBits(int numBits, WithOption... options) throws BufferException;

    /**
     * Reads an unsigned byte value with the specified bit length.
     *
     * @param numBits the number of bits to read
     * @param options additional options for reading
     * @return the unsigned byte value as a short value
     */
    byte readUnsignedByte(int numBits, WithOption... options) throws BufferException;

    /**
     * Reads an unsigned short value with the specified bit length.
     *
     * @param numBits the number of bits to read
     * @param options additional options for reading
     * @return the unsigned short value as an int value
     */
    short readUnsignedShort(int numBits, WithOption... options) throws BufferException;

    /**
     * Reads an unsigned int value with the specified bit length.
     *
     * @param numBits the number of bits to read
     * @param options additional options for reading
     * @return the unsigned int value as a long value
     */
    int readUnsignedInt(int numBits, WithOption... options) throws BufferException;

    /**
     * Reads an unsigned long value with the specified bit length.
     *
     * @param numBits the number of bits to read
     * @param options additional options for reading
     * @return the unsigned long value as a BigInteger value
     */
    long readUnsignedLong(int numBits, WithOption... options) throws BufferException;

    /**
     * Reads an unsigned big integer with the specified bit length.
     *
     * @param bitLength the number of bits to read
     * @param options   additional options for reading
     * @return the unsigned big integer value
     */
    BigInteger readUnsignedBigInteger(int bitLength, WithOption... options) throws BufferException;

    /**
     * Reads a signed byte value with the specified bit length.
     *
     * @param numBits the number of bits to read
     * @param options additional options for reading
     * @return the signed byte value
     */
    byte readSignedByte(int numBits, WithOption... options) throws BufferException;

    /**
     * Reads a signed short value with the specified bit length.
     *
     * @param numBits the number of bits to read
     * @param options additional options for reading
     * @return the signed short value
     */
    short readSignedShort(int numBits, WithOption... options) throws BufferException;

    /**
     * Reads a signed int value with the specified bit length.
     *
     * @param numBits the number of bits to read
     * @param options additional options for reading
     * @return the signed int value
     */
    int readSignedInt(int numBits, WithOption... options) throws BufferException;

    /**
     * Reads a signed long value with the specified bit length.
     *
     * @param numBits the number of bits to read
     * @param options additional options for reading
     * @return the signed long value
     */
    long readSignedLong(int numBits, WithOption... options) throws BufferException;

    /**
     * Reads a signed big integer with the specified bit length.
     *
     * @param bitLength the number of bits to read
     * @param options   additional options for reading
     * @return the signed big integer value
     */
    BigInteger readSignedBigInteger(int bitLength, WithOption... options) throws BufferException;

    /**
     * Reads a float value with the specified bit length.
     *
     * @param numBits the number of bits to read
     * @param options additional options for reading
     * @return the float value
     */
    float readFloat(int numBits, WithOption... options) throws BufferException;

    /**
     * Reads a double value with the specified bit length.
     *
     * @param numBits the number of bits to read
     * @param options additional options for reading
     * @return the double value
     */
    double readDouble(int numBits, WithOption... options) throws BufferException;

    /**
     * Reads a big decimal value with the specified bit length.
     *
     * @param numBits the number of bits to read
     * @param options additional options for reading
     * @return the big decimal value
     */
    BigDecimal readBigDecimal(int numBits, WithOption... options) throws BufferException;

    /**
     * Reads a string value with the specified bit length.
     *
     * @param numBits the number of bits to read
     * @param options additional options for reading
     * @return the string value
     */
    String readString(int numBits, WithOption... options) throws BufferException;

    /**
     * Creates a sub-buffer with the specified bit length.
     *
     * @param numBits the number of bits for the sub-buffer
     * @param options additional options for creating the sub-buffer
     * @return a new ReadBuffer instance representing the sub-buffer
     */
    ReadBuffer createSubBuffer(int numBits, WithOption... options) throws BufferException;

    /**
     * Gets the current position in the buffer in bits.
     *
     * @return the current position in bits
     */
    int getPositionInBits();

    /**
     * Gets the number of remaining bits in the buffer.
     *
     * @return the number of remaining bits
     */
    int getRemainingBits();

    /**
     * Sets the current position in the buffer in bits.
     *
     * @param positionInBits the new position in bits
     */
    void setPositionInBits(int positionInBits);

}
