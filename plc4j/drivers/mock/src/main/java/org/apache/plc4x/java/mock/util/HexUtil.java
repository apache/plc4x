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
import org.apache.commons.lang3.ArrayUtils;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;

import java.io.ByteArrayOutputStream;
import java.io.IOException;

public class HexUtil {

    /**
     * Can be used to display a byte[] as hex string.
     *
     * @param bytes to be displayed.
     * @return hexed string.
     */
    public static String toHex(Byte[] bytes) {
        return toHex(ArrayUtils.toPrimitive(bytes));
    }

    public static String toHex(byte[] bytes) {
        try (ByteArrayOutputStream byteArrayOutputStream = new ByteArrayOutputStream()) {
            HexDump.dump(bytes, 0, byteArrayOutputStream, 0);
            return byteArrayOutputStream.toString();
        } catch (IOException e) {
            throw new PlcRuntimeException(e);
        }
    }
}
