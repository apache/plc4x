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

import org.apache.commons.lang3.StringUtils;
import org.junit.Test;

import java.util.stream.IntStream;

public class DigestUtilTest {

    @Test
    public void displayValue() throws Exception {
        Object[] crcs = IntStream.range(0, 256)
            .map(value -> DigestUtil.calculateCrc16(new byte[]{(byte) value}))
            .mapToObj(Integer::toHexString)
            .map(s -> StringUtils.leftPad(s, 4, '0'))
            .toArray(String[]::new);
        for (int i = 0; i < 32; i++) {
            String col1 = "" + (i) + '\t' + crcs[i];
            String col2 = "" + (i + 32) + '\t' + crcs[i + 32];
            String col3 = "" + (i + 64) + '\t' + crcs[i + 64];
            String col4 = "" + (i + 96) + '\t' + crcs[i + 96];
            String col5 = "" + (i + 128) + '\t' + crcs[i + 128];
            String col6 = "" + (i + 160) + '\t' + crcs[i + 160];
            String col7 = "" + (i + 192) + '\t' + crcs[i + 192];
            String col8 = "" + (i + 224) + '\t' + crcs[i + 224];
            System.out.println(col1 + "\t\t" + col2 + "\t\t" + col3 + "\t\t" + col4 + "\t\t" + col5 + "\t\t" + col6 + "\t\t" + col7 + "\t\t" + col8);
        }
    }
}