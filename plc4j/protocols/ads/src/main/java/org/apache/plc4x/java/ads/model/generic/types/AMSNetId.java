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
package org.apache.plc4x.java.ads.model.generic.types;

import org.apache.plc4x.java.ads.model.util.ByteValue;

/**
 * The AMSNetId consists of 6 bytes and addresses the transmitter or receiver. One possible AMSNetId would be e.g.. 172.16.17.10.1.1. The storage arrangement in this example is as follows:
 * <p>
 * _____0     1     2     3     4     5
 * __+-----------------------------------+
 * 0 | 127 |  16 |  17 |  10 |   1 |   1 |
 * __+-----------------------------------+
 * <p>
 * <p>
 * The AMSNetId is purely logical and has usually no relation to the IP address. The AMSNetId is configurated at the target system. At the PC for this the TwinCAT System Control is used. If you use other hardware, see the considering documentation for notes about settings of the AMS NetId.
 */
public class AMSNetId extends ByteValue {

    public AMSNetId(int octed1, int octed2, int octed3, int octed4, int octed5, int octed6) {
        super((byte) octed1, (byte) octed2, (byte) octed3, (byte) octed4, (byte) octed5, (byte) octed6);
    }

    public AMSNetId(byte... value) {
        super(value);
        assertLength(6);
    }
}
