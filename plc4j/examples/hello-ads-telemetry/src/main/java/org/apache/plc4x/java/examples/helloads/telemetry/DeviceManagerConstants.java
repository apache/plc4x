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
package org.apache.plc4x.java.examples.helloads.telemetry;

// Extracted from https://infosys.beckhoff.com/index.php?content=../content/1031/devicemanager/262982923.html

public enum DeviceManagerConstants {
    NIC(0x00000002),
    Time(0x00000003),
    UserManagement(0x00000004),
    RAS(0x00000005),
    FTP(0x00000006),
    SMB(0x00000007),
    TwinCat(0x00000008),
    Software(0x0000000A),
    CPU(0x0000000B),
    Memory(0x0000000C),
    FirewallWinCE(0x0000000E),
    FileSystemObject(0x00000010),
    DisplayDevice(0x00000013),
    EWF(0x00000014),
    FBWF(0x00000015),
    OS(0x00000018),
    RAID(0x00000019),
    Fan(0x0000001B),
    Mainboard(0x0000001C),
    DiskManagement(0x0000001D),
    UPS(0x0000001E),
    PhysicalDrive(0x0000001F),
    MassStorageDrive(0x00000020),
    UnifiedWriteFilter(0x00000021),
    IO(0x00000022),
    Misc(0x00000100);

    final int typeNumber;
    DeviceManagerConstants(int typeNumber) {
        this.typeNumber = typeNumber;
    }
}
