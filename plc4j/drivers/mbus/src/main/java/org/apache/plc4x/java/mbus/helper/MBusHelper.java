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

import org.apache.plc4x.java.mbus.readwrite.DataInformationFieldExtension;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBuffer;

public class MBusHelper {

    public static byte crc(byte[] frame) {
        return crc(frame, 3);
    }

    public static byte crc(byte[] frame, int headerLength) {
        return crc(frame, headerLength, 0);
    }

    public static byte crc(byte[] frame, int headerLength, int skipBytes) {
        int checksum = 0;
        for (int index = headerLength; index < (frame.length - skipBytes); index++) {
            checksum += frame[index];
        }
        return (byte) (checksum & 0xff);
    }

    public static boolean isLastExtension(ReadBuffer io) throws ParseException {
        return (io.peekByte(-1) & 0x80) != 0x80;
    }

    public static DataInformationFieldExtension mergeDife(int storage, DataInformationFieldExtension[] difes) {
        DataInformationFieldExtension ext = null;
        for (DataInformationFieldExtension dife : difes) {
            if (ext == null) {
                ext =  dife;
            } else {
                /*
                ext = new DataInformationFieldExtension(
                    (byte) (ext.getSubunit() + dife.getSubunit()),
                    (byte) (ext.getTariff() + dife.getTariff()),
                    ext.getStorage() + dife.getStorage()
                );
                */
            }

        }
        return ext;
    }

    public static int difeSize(DataInformationFieldExtension dife) {
        return 0;
    }

    public static DataInformationFieldExtension readDife(ReadBuffer io, boolean extension, long storage) throws ParseException {
        if (!extension) {
            return null;
        }

        int subunit = 0;
        int tariff = 0;

        int numDife = 0;
        byte val;
        do {
            val = io.readByte(8);
            subunit += (((val & 0x40) >> 6) << numDife);
            tariff += ((val & 0x30) >> 4) << (numDife * 2);
            storage += ((val & 0x0f) << ((numDife * 4) + 1));
            numDife++;
        } while ((val & 0x80) == 0x80);

//        return new DataInformationFieldExtension(
//            subunit, tariff, BigInteger.valueOf(storage)
//        );
        return null;
    }

    public static void writeDife(Object value) {

    }
}
