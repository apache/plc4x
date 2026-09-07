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
package org.apache.plc4x.java.spi.buffers.asciiboxbased.utils.hex;

import org.apache.plc4x.java.spi.buffers.asciiboxbased.utils.utils.StringUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.io.ObjectOutputStream;
import java.util.AbstractMap;
import java.util.Arrays;
import java.util.Map;
import java.util.Set;
import java.util.stream.Collectors;

/**
 * Utility class for creating hexadecimal dump representations of byte arrays.
 *
 * <p>This class produces human-readable hex dumps similar to the output of tools like
 * {@code hexdump} or {@code xxd}. The output format includes:
 * <ul>
 *   <li>A row index (byte offset from the start)</li>
 *   <li>Hexadecimal representation of each byte</li>
 *   <li>ASCII representation of printable characters (non-printable shown as '.')</li>
 * </ul>
 *
 * <p>Example output:
 * <pre>
 * 000|48 65 6c 6c 6f 20 57 6f 72 6c 64 21 'Hello World!'
 * </pre>
 *
 * <p>The class also supports highlighting specific byte positions using ANSI color codes,
 * which is useful for debugging and visualizing specific parts of binary data.
 */
public class Hex {

    private static final Logger LOGGER = LoggerFactory.getLogger(Hex.class);

    /**
     * Default width for hex dump output in characters.
     * This value (46) allows approximately 10 bytes per line for arrays with less than 999 bytes.
     */
    public static final int DefaultWidth = 46;

    /**
     * Overhead per line when the hex dump is rendered inside ASCII boxes.
     * Accounts for left and right border characters.
     */
    public static final int boxLineOverheat = 1 + 1;

    /**
     * Width of a single blank/space character.
     */
    public static final int blankWidth = 1;

    /**
     * Width required to display one byte: 2 hex digits + 1 space separator.
     */
    public static final int byteWidth = 2 + 1;

    /**
     * Width of the pipe character '|' used as separator between index and hex data.
     */
    public static final int pipeWidth = 1;

    /**
     * Debug flag. When set to {@code true}, enables detailed debug logging
     * of the hex dump calculation process.
     */
    public static boolean DebugHex;

    // Private constructor to prevent instantiation of utility class
    private Hex() {
    }

    /**
     * Creates a hex dump of the given byte array using the default width.
     *
     * @param data the byte array to dump; may be {@code null} or empty
     * @return a formatted hex dump string, or empty string if data is null/empty
     * @see #dump(byte[], int, int...)
     */
    public static String dump(byte[] data) {
        return dump(data, DefaultWidth);
    }

    /**
     * Creates a hex dump of the given byte array with configurable width and optional highlighting.
     *
     * <p>The output format for each row is:
     * <pre>
     * INDEX|HH HH HH ... 'ASCII'
     * </pre>
     * Where:
     * <ul>
     *   <li>INDEX is the zero-padded byte offset</li>
     *   <li>HH are two-digit lowercase hex values separated by spaces</li>
     *   <li>ASCII is the printable character representation (non-printable chars shown as '.')</li>
     * </ul>
     *
     * @param data             the byte array to dump; may be {@code null} or empty
     * @param desiredCharWidth the desired output width in characters (minimum 18)
     * @param highlights       optional byte indices to highlight with ANSI red color
     * @return a formatted hex dump string, or empty string if data is null/empty
     */
    public static String dump(byte[] data, int desiredCharWidth, int... highlights) {
        // Handle null or empty input
        if (data == null || data.length < 1) {
            return "";
        }

        // Convert highlight indices to a Set for O(1) lookup
        Set<Integer> highlightsSet = Arrays.stream(highlights).boxed().collect(Collectors.toSet());

        // Copy the array to avoid mutating the original during maskString()
        data = Arrays.copyOf(data, data.length);

        StringBuilder hexString = new StringBuilder();

        // Calculate how many bytes fit per row and the width needed for the index column
        Map.Entry<Integer, Integer> rowIndexCalculation = calculateBytesPerRowAndIndexWidth(data.length, desiredCharWidth);
        int maxBytesPerRow = rowIndexCalculation.getKey();
        int indexWidth = rowIndexCalculation.getValue();

        // Iterate through the data, processing one row at a time
        for (int byteIndex = 0, rowIndex = 0; byteIndex < data.length; byteIndex = byteIndex + maxBytesPerRow, rowIndex = rowIndex + 1) {

            // Build the index prefix (e.g., "000|" or "0012|")
            String indexString = String.format("%1$" + indexWidth + "s|", byteIndex).replace(' ', '0');
            hexString.append(indexString);

            // Output each byte in this row as hex
            for (int columnIndex = 0; columnIndex < maxBytesPerRow; columnIndex++) {
                int absoluteIndex = byteIndex + columnIndex;

                if (absoluteIndex < data.length) {
                    // Apply ANSI red color for highlighted bytes
                    if (highlightsSet.contains(absoluteIndex)) {
                        hexString.append("\033[0;31m");  // ANSI red
                    }

                    // Append the hex value (2 digits + space)
                    hexString.append(String.format("%02x ", data[absoluteIndex]));

                    // Reset color after highlighted byte
                    if (highlightsSet.contains(absoluteIndex)) {
                        hexString.append("\033[0m");  // ANSI reset
                    }
                } else {
                    // Pad with spaces in case of an incomplete last row
                    hexString.append(" ".repeat(byteWidth));
                }
            }

            // Calculate the end index for the ASCII representation
            int endIndex = Math.min(byteIndex + maxBytesPerRow, data.length);

            // Create ASCII representation (non-printable chars replaced with '.')
            String stringRepresentation = maskString(Arrays.copyOfRange(data, byteIndex, endIndex));

            // Pad the ASCII representation if this is a partial row
            if (stringRepresentation.length() < maxBytesPerRow) {
                stringRepresentation += StringUtils.repeat(" ", (maxBytesPerRow - stringRepresentation.length()) % maxBytesPerRow);
            }

            // Append the quoted ASCII representation
            hexString.append(String.format("'%s'\n", stringRepresentation));
        }

        // Remove the trailing newline
        return hexString.substring(0, hexString.length() - 1);
    }

    /**
     * Calculates the optimal number of bytes per row and the index column width.
     *
     * <p>This method determines how to best fit the hex dump within the desired width by:
     * <ol>
     *   <li>Calculating the number of digits needed for the row index</li>
     *   <li>Determining the minimum required width for a valid output</li>
     *   <li>Solving for the maximum bytes that can fit per row</li>
     * </ol>
     *
     * <p>The layout formula accounts for:
     * <ul>
     *   <li>Index column: variable width based on total byte count</li>
     *   <li>Pipe separator: 1 character</li>
     *   <li>Hex bytes: 3 characters each (2 hex digits + space)</li>
     *   <li>ASCII column: 1 character per byte + 2 quote characters</li>
     * </ul>
     *
     * @param numberOfBytes      total number of bytes in the data
     * @param desiredStringWidth desired output width in characters
     * @return a Map.Entry where key = bytes per row, value = index digit width
     */
    static Map.Entry<Integer, Integer> calculateBytesPerRowAndIndexWidth(int numberOfBytes, int desiredStringWidth) {
        if (DebugHex) {
            LOGGER.debug("Calculating max row and index for {} number of bytes and a desired string width of {}", numberOfBytes, desiredStringWidth);
        }

        // Calculate how many digits we need for the index (e.g., 3 digits for 100-999 bytes)
        int indexDigits = (int) (Math.log10(numberOfBytes) + 1);
        int requiredIndexWidth = indexDigits + pipeWidth;

        if (DebugHex) {
            LOGGER.debug("index width {} for indexDigits {} for bytes {}", requiredIndexWidth, indexDigits, numberOfBytes);
        }

        // Account for the quote characters around the ASCII representation
        int quoteRune = 1;
        int potentialStringRenderRune = 1;

        // Calculate the minimum number of spaces needed for at least one byte: "0|00 '.'"
        int availableSpace = requiredIndexWidth + byteWidth + quoteRune + potentialStringRenderRune + quoteRune;

        if (DebugHex) {
            LOGGER.debug("calculated {} minimal width for number of bytes {}", availableSpace, numberOfBytes);
        }

        // Use the larger of desired width or minimum required width
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

        // Solve for maximum bytes per row using the layout equation:
        // totalWidth = indexWidth + (bytesPerRow * byteWidth) + quote + (bytesPerRow * 1) + quote
        // Rearranging: bytesPerRow = (totalWidth - indexWidth - 2*quote) / (byteWidth + 1)
        double z = availableSpace;       // total available width
        double y = requiredIndexWidth;   // index column width
        double a = byteWidth;            // width per byte in hex section (3)
        double b = quoteRune;            // quote character width (1)

        // Formula: x = (-2*b - y + z) / (a + 1) where x = bytesPerRow
        double x = ((-2 * b) - y + z) / (a + 1);

        if (DebugHex) {
            LOGGER.debug("Calculated number of bytes per row {} in int {}", x, (int) x);
        }

        return new AbstractMap.SimpleEntry<>((int) x, indexDigits);
    }

    /**
     * Converts a byte array to a printable ASCII string.
     *
     * <p>Non-printable characters (bytes outside the range 32-126) are replaced
     * with a period ('.') character. This is the standard convention for hex dump
     * ASCII representations.
     *
     * <p><strong>Note:</strong> This method modifies the input array in place.
     *
     * @param data the byte array to convert (will be modified)
     * @return a String representation with non-printable characters masked
     */
    static String maskString(byte[] data) {
        for (int i = 0; i < data.length; i++) {
            // ASCII printable range is 32 (space) to 126 (~)
            if (data[i] < 32 || data[i] > 126) {
                data[i] = '.';
            }
        }
        return new String(data);
    }

    /**
     * Serializes a Java object to a byte array using standard Java serialization.
     *
     * <p>This utility method can be used to convert objects to bytes for hex dumping.
     *
     * @param obj the object to serialize (must implement {@link java.io.Serializable})
     * @return the serialized byte array
     * @throws RuntimeException if serialization fails due to an IOException
     */
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
