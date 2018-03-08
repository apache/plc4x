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
package org.apache.plc4x.java.ads.protocol.util;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.ads.api.util.ByteReadable;

/**
 * TODO: temporary due to unclear licence
 * From https://stackoverflow.com/a/18333436/850036
 * // TODO: replace this with a better implementation. Maybe:
 * http://www.source-code.biz/snippets/java/crc16/
 * https://en.wikipedia.org/wiki/Cyclic_redundancy_check#Standards_and_common_use
 * or even better a netty supplied crc-16
 * https://github.com/openhab/jamod/blob/64cdbd16fbb7febd39470f873a00be986b052e39/src/main/java/net/wimpi/modbus/util/ModbusUtil.java
 */
public class DigestUtil {

    private DigestUtil() {
        // Utility class
    }

    public static int calculateCrc16(ByteReadable... byteReadables) {
        if (byteReadables.length == 1) {
            return calculateCrc16(byteReadables[0].getBytes());
        }
        return calculateCrc16(new ByteReadable() {
            @Override
            public ByteBuf getByteBuf() {
                return buildByteBuff(byteReadables);
            }
        }.getBytes());
    }

    public static int calculateCrc16(byte[] bytes) {
        // TODO: implement me
        return 0;
    }
}
