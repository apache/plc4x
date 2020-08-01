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
package org.apache.plc4x.java.amsads.attic.connection;

import org.apache.plc4x.java.amsads.readwrite.AmsNetId;

import java.net.InetAddress;
import java.util.Objects;

@Deprecated
public class AdsConnectionFactory {

    public AdsTcpPlcConnection adsTcpPlcConnectionOf(InetAddress address, Integer port, AmsNetId targetAmsNetId, Integer targetAmsPort, AmsNetId sourceAmsNetId, Integer sourceAmsPort) {
        Objects.requireNonNull(address);
        Objects.requireNonNull(targetAmsNetId);
        Objects.requireNonNull(targetAmsPort);
        if (sourceAmsNetId == null || sourceAmsPort == null) {
            if (port == null) {
                return AdsTcpPlcConnection.of(address, targetAmsNetId, targetAmsPort);
            } else {
                return AdsTcpPlcConnection.of(address, port, targetAmsNetId, targetAmsPort);
            }
        } else {
            if (port == null) {
                return AdsTcpPlcConnection.of(address, targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort);
            } else {
                return AdsTcpPlcConnection.of(address, port, targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort);
            }
        }
    }

    public AdsSerialPlcConnection adsSerialPlcConnectionOf(String serialPort, AmsNetId targetAmsNetId, Integer targetAmsPort, AmsNetId sourceAmsNetId, Integer sourceAmsPort) {
        Objects.requireNonNull(serialPort);
        Objects.requireNonNull(targetAmsNetId);
        Objects.requireNonNull(targetAmsPort);
        if (sourceAmsNetId == null || sourceAmsPort == null) {
            return AdsSerialPlcConnection.of(serialPort, targetAmsNetId, targetAmsPort);
        } else {
            return AdsSerialPlcConnection.of(serialPort, targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort);
        }
    }
}
