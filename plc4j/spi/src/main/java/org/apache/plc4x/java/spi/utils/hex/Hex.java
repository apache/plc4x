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

package org.apache.plc4x.java.spi.utils.hex;

import org.apache.commons.lang3.ArrayUtils;
import org.apache.commons.lang3.StringUtils;
import org.apache.commons.lang3.tuple.Pair;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.io.ObjectOutputStream;
import java.util.Arrays;
import java.util.Set;
import java.util.stream.Collectors;

public class Hex {

    private static final Logger LOGGER = LoggerFactory.getLogger(Hex.class);

    /**
     * DefaultWidth defaults to a default screen dumps size
     */
    public static final int DefaultWidth = 46; // 10 bytes per line on a []byte < 999

    /**
     * boxLineOverheat Overheat per line when drawing boxes
     */
    public static final int boxLineOverheat = 1 + 1;

    /**
     * blankWidth blank size of blank
     */
    public static final int blankWidth = 1;

    /**
     * byteWidth required size of runes required to print one bytes 2 hex digits + 1 blanks
     */
    public static final int byteWidth = 2 + 1;

    /**
     * pipeWidth size of the pipe char
     */
    public static final int pipeWidth = 1;

    /**
     * DebugHex set to true to get debug messages
     */
    public static boolean DebugHex;

    /**
     * dumps a 56 char wide hex string
     *
     * @param data bytes to dump
     * @return the hex string
     */
    public static String dump(byte[] data) {
        return dump(data, DefaultWidth);
    }

    // DumpFixedWidth dumps hex as hex string. Min width of string returned is 18 up to supplied charWidth
    public static String dump(byte[] data, int desiredCharWidth, int... highlights) {
        if (data == null || data.length < 1) {
            return "";
        }
        Set<Integer> highlightsSet = Arrays.stream(highlights).boxed().collect(Collectors.toSet());
        // We copy the array to avoid mutations
        data = Arrays.copyOf(data, data.length);
        StringBuilder hexString = new StringBuilder();
        Pair<Integer, Integer> rowIndexCalculation = calculateBytesPerRowAndIndexWidth(data.length, desiredCharWidth);
        int maxBytesPerRow = rowIndexCalculation.getLeft();
        int indexWidth = rowIndexCalculation.getRight();
        for (int byteIndex = 0, rowIndex = 0; byteIndex < data.length; byteIndex = byteIndex + maxBytesPerRow, rowIndex = rowIndex + 1) {
            String indexString = String.format("%1$" + indexWidth + "s|", byteIndex).replace(' ', '0');
            hexString.append(indexString);
            for (int columnIndex = 0; columnIndex < maxBytesPerRow; columnIndex++) {
                int absoluteIndex = byteIndex + columnIndex;
                if (absoluteIndex < data.length) {
                    if(highlightsSet.contains(absoluteIndex)) {
                        hexString.append("\033[0;31m");
                    }
                    hexString.append(String.format("%02x ", data[absoluteIndex]));
                    if(highlightsSet.contains(absoluteIndex)) {
                        hexString.append("\033[0m");
                    }
                } else {
                    // align with empty byte representation
                    hexString.append(StringUtils.repeat(" ", byteWidth));
                }
            }
            int endIndex = byteIndex + maxBytesPerRow;
            if (endIndex >= data.length) {
                endIndex = data.length;
            }
            String stringRepresentation = maskString(ArrayUtils.subarray(data, byteIndex, endIndex));
            if (stringRepresentation.length() < maxBytesPerRow) {
                stringRepresentation += StringUtils.repeat(" ", (maxBytesPerRow - stringRepresentation.length()) % maxBytesPerRow);
            }
            hexString.append(String.format("'%s'\n", stringRepresentation));
        }
        // remove last newline
        return hexString.substring(0, hexString.length() - 1);
    }


    static Pair<Integer, Integer> calculateBytesPerRowAndIndexWidth(int numberOfBytes, int desiredStringWidth) {
        if (DebugHex) {
            LOGGER.debug("Calculating max row and index for {} number of bytes and a desired string width of {}", numberOfBytes, desiredStringWidth);
        }
        int indexDigits = (int) (Math.log10(numberOfBytes) + 1);
        int requiredIndexWidth = indexDigits + pipeWidth;
        if (DebugHex) {
            LOGGER.debug("index width {} for indexDigits {} for bytes {}", requiredIndexWidth, indexDigits, numberOfBytes);
        }
        // strings get quoted by 2 chars
        int quoteRune = 1;
        int potentialStringRenderRune = 1;
        // 0 00  '.'
        int availableSpace = requiredIndexWidth + byteWidth + quoteRune + potentialStringRenderRune + quoteRune;
        if (DebugHex) {
            LOGGER.debug("calculated {} minimal width for number of bytes {}", availableSpace, numberOfBytes);
        }
        if (desiredStringWidth >= availableSpace) {
            availableSpace = desiredStringWidth;
        } else {
            if (DebugHex) {
                LOGGER.debug("Overflow by {} runes", desiredStringWidth - availableSpace);
            }
        }
        if (DebugHex) {
            LOGGER.debug("Actual space {}", availableSpace);
        }

        double z = availableSpace;
        double y = requiredIndexWidth;
        double a = byteWidth;
        double b = quoteRune;
        // c = needed space for bytes x * byteWidth
        // x = maxBytesPerRow
        // x = (z - (y + b + x * 1 + b)) / a == x = (-2 * b - y + z)/(a + 1) and a + 1!=0 and a!=0
        double x = ((-2 * b) - y + z) / (a + 1);
        if (DebugHex) {
            LOGGER.debug("Calculated number of bytes per row {} in int {}", x, (int) x);
        }
        return Pair.of((int) x, indexDigits);
    }

    static String maskString(byte[] data) {
        for (int i = 0; i < data.length; i++) {
            if (data[i] < 32 || data[i] > 126) {
                data[i] = '.';
            }
        }
        return new String(data);
    }

    static byte[] toBytes(Object obj) {
        ByteArrayOutputStream boas = new ByteArrayOutputStream();
        try (ObjectOutputStream ois = new ObjectOutputStream(boas)) {
            ois.writeObject(obj);
            return boas.toByteArray();
        } catch (IOException ioe) {
            throw new RuntimeException(ioe);
        }
    }
}
