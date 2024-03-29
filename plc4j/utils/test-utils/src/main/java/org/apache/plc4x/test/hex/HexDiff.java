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

package org.apache.plc4x.test.hex;

import org.apache.plc4x.java.spi.utils.ascii.AsciiBox;
import org.apache.plc4x.java.spi.utils.ascii.AsciiBoxWriter;

import java.util.LinkedList;
import java.util.List;

public class HexDiff {
    public static AsciiBox diffHex(byte[] expectedBytes, byte[] actualBytes) {
        int numBytes = Math.min(expectedBytes.length, actualBytes.length);
        int brokenAt = -1;
        List<Integer> diffIndexes = new LinkedList<>();
        for (int i = 0; i < numBytes; i++) {
            if (expectedBytes[i] != actualBytes[i]) {
                if (brokenAt < 0) {
                    brokenAt = i;
                }
                diffIndexes.add(i);
            }
        }
        String expectedHex = org.apache.plc4x.java.spi.utils.hex.Hex.dump(expectedBytes, 46, diffIndexes.stream().mapToInt(integer -> integer).toArray());
        String actialHex = org.apache.plc4x.java.spi.utils.hex.Hex.dump(actualBytes, 46, diffIndexes.stream().mapToInt(integer -> integer).toArray());
        return AsciiBoxWriter.DEFAULT.boxSideBySide(AsciiBoxWriter.DEFAULT.boxString("expected", expectedHex, 0), AsciiBoxWriter.DEFAULT.boxString("actual", actialHex, 0));
    }
}
