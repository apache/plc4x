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
package org.apache.plc4x.java.mock.util;

import org.apache.commons.io.HexDump;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.hamcrest.Description;
import org.hamcrest.Matcher;
import org.hamcrest.core.IsEqual;

import java.io.ByteArrayOutputStream;
import java.io.IOException;

import static java.lang.Integer.toHexString;
import static org.apache.commons.lang3.StringUtils.leftPad;
import static org.junit.jupiter.api.Assertions.assertEquals;

public class Assert {

    public static void assertByteEquals(byte expected, byte actual) {
        String expectedHex = "0x" + leftPad(toHexString(expected), 2, '0');
        String actualHex = "0x" + leftPad(toHexString(actual), 2, '0');
        assertEquals(expectedHex, actualHex);
    }

    public static void assertByteEquals(byte[] expected, byte[] actual) throws IOException {
        String expectedHex = cleanHexDump(dump(expected));
        String actualHex = cleanHexDump(dump(actual));
        assertEquals(expectedHex, actualHex);
    }

    public static Matcher<byte[]> byteArrayEqualsTo(byte[] expected) {
        return new IsEqual<byte[]>(expected) {
            @Override
            public void describeTo(Description description) {
                try {
                    String dump = dump(expected);
                    description.appendText("\n").appendText(cleanHexDump(dump));
                } catch (IOException e) {
                    throw new PlcRuntimeException(e);
                }
            }

            @Override
            public void describeMismatch(Object item, Description description) {
                if (item == null || !(item instanceof byte[])) {
                    super.describeMismatch(item, description);
                    return;
                }
                try {
                    String dump = dump((byte[]) item);
                    description.appendText("was ").appendText("\n").appendText(cleanHexDump(dump));
                } catch (IOException e) {
                    throw new PlcRuntimeException(e);
                }
            }
        };
    }

    private static String cleanHexDump(String in) {
        return in.replaceAll("@.*\\{", "@XXXXXXXX{");
    }

    private static String dump(byte[] bytes) throws IOException {
        try (ByteArrayOutputStream byteArrayOutputStream = new ByteArrayOutputStream()) {
            HexDump.dump(bytes, 0, byteArrayOutputStream, 0);
            return HexDump.EOL + byteArrayOutputStream.toString();
        }
    }
}
