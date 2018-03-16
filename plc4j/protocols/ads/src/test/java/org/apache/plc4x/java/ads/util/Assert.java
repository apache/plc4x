/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing,
 software distributed under the License is distributed on an
 "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 KIND, either express or implied.  See the License for the
 specific language governing permissions and limitations
 under the License.
 */
package org.apache.plc4x.java.ads.util;

import org.apache.plc4x.java.ads.api.util.ByteValue;
import org.hamcrest.Description;
import org.hamcrest.Matcher;
import org.hamcrest.core.IsEqual;

import java.io.IOException;

import static java.lang.Integer.toHexString;
import static org.apache.commons.lang3.StringUtils.leftPad;
import static org.junit.Assert.assertEquals;

public class Assert {

    public static void assertByteEquals(byte expected, byte actual) throws Exception {
        String expectedHex = "0x" + leftPad(toHexString(expected), 2, '0');
        String actualHex = "0x" + leftPad(toHexString(actual), 2, '0');
        assertEquals(expectedHex, actualHex);
    }

    public static void assertByteEquals(byte[] expected, byte[] actual) throws Exception {
        String expectedHex = cleanHexDump(ByteValue.of(expected).dump());
        String actualHex = cleanHexDump(ByteValue.of(actual).dump());
        assertEquals(expectedHex, actualHex);
    }

    public static Matcher<byte[]> byteArrayEqualsTo(byte[] expected) {
        return new IsEqual<byte[]>(expected) {
            @Override
            public void describeTo(Description description) {
                try {
                    String dump = ByteValue.of(expected).dump();
                    description.appendText("\n").appendText(cleanHexDump(dump));
                } catch (IOException e) {
                    throw new RuntimeException(e);
                }
            }

            @Override
            public void describeMismatch(Object item, Description description) {
                if (item == null || !(item instanceof byte[])) {
                    super.describeMismatch(item, description);
                    return;
                }
                try {
                    String dump = ByteValue.of((byte[]) item).dump();
                    description.appendText("was ").appendText("\n").appendText(cleanHexDump(dump));
                } catch (IOException e) {
                    throw new RuntimeException(e);
                }
            }
        };
    }

    private static String cleanHexDump(String in) {
        return in.replaceAll("@.*\\{", "@XXXXXXXX{");
    }
}
