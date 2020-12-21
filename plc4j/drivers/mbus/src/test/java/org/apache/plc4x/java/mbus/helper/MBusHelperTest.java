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
package org.apache.plc4x.java.mbus.helper;

import org.apache.commons.codec.DecoderException;
import org.apache.commons.codec.binary.Hex;
import org.junit.jupiter.api.Test;

import static org.assertj.core.api.Assertions.assertThat;

class MBusHelperTest {

    String FRAME = "681F1F680802727856341224400107550000000313153100DA023B13018B60043718021816";
    /*
    "68 1F 1F 68" + // header + len=31
    "08 02 72" + // C=08 (RSP), node=2, CI=72
    "78 56 34 12" + // identification=12345678
    "24 40 01 07" + // manufacturer id=4024, version =1, device type = water meter
    "55 00 00 00" +
    "03 13 15 31 00" + // data block 1, unit 0, storage 0, instant value, volume  = 123456 l
    "DA 02 3B 13 01" + // data block 2, unit 0, storage 5, max value, volume flow = 113 l/h
    "8B 60 04 37 18 02" + // data block 3, unit 1, storage 0, tariff 2, instant value, energy = 218,37 KwH
    "18 16" // checksum, stop sign
     */

    @Test
    void calculateChecksum() throws DecoderException {
        byte[] frame = Hex.decodeHex(FRAME);
        byte crc = MBusHelper.crc(frame, 4, 2);
        assertThat(Integer.toHexString(crc)).isEqualTo("18");
    }
}