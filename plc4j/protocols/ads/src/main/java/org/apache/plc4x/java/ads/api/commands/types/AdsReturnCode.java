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
package org.apache.plc4x.java.ads.api.commands.types;

import static java.lang.Long.toHexString;
import static java.util.Objects.requireNonNull;

/**
 * Based on spec from: https://infosys.beckhoff.com/content/1033/tcadscommon/html/ads_returncodes.htm
 */
// Error codes: 0x000..., 0x500..., 0x700..., 0x1000..., 0x274C...
public enum AdsReturnCode {
    // Global Error Codes
    ADS_CODE_0(0x0, 0, "no error"),
    ADS_CODE_1(0x1, 1, "Internal error"),
    ADS_CODE_2(0x2, 2, "No Rtime"),
    ADS_CODE_3(0x3, 3, "Allocation locked memory error"),
    ADS_CODE_4(0x4, 4, "Insert mailbox error", "No ADS mailbox was available to process this message.", "Reduce the number of ADS calls (e.g ADS-Sum commands or Max Delay Parameter)"),
    ADS_CODE_5(0x5, 5, "Wrong receive HMSG"),
    ADS_CODE_6(0x6, 6, "target port not found", "ADS Server not started"),
    ADS_CODE_7(0x7, 7, "target machine not found", "Missing ADS routes"),
    ADS_CODE_8(0x8, 8, "Unknown command ID"),
    ADS_CODE_9(0x9, 9, "Bad task ID"),
    ADS_CODE_10(0xA, 10, "No IO"),
    ADS_CODE_11(0xB, 11, "Unknown ADS command"),
    ADS_CODE_12(0xC, 12, "Win 32 error"),
    ADS_CODE_13(0xD, 13, "Port not connected"),
    ADS_CODE_14(0xE, 14, "Invalid ADS length"),
    ADS_CODE_15(0xF, 15, "Invalid ADS Net ID"),
    ADS_CODE_16(0x10, 16, "Low Installation level"),
    ADS_CODE_17(0x11, 17, "No debug available"),
    ADS_CODE_18(0x12, 18, "Port disabled"),
    ADS_CODE_19(0x13, 19, "Port already connected"),
    ADS_CODE_20(0x14, 20, "ADS Sync Win32 error"),
    ADS_CODE_21(0x15, 21, "ADS Sync Timeout"),
    ADS_CODE_22(0x16, 22, "ADS Sync AMS error"),
    ADS_CODE_23(0x17, 23, "ADS Sync no index map"),
    ADS_CODE_24(0x18, 24, "Invalid ADS port"),
    ADS_CODE_25(0x19, 25, "No memory"),
    ADS_CODE_26(0x1A, 26, "TCP send error"),
    ADS_CODE_27(0x1B, 27, "Host unreachable"),
    ADS_CODE_28(0x1C, 28, "Invalid AMS fragment"),


    // Router Error Codes
    ADS_CODE_1280(0x500, 1280, "ROUTERERR_NOLOCKEDMEMORY", "No locked memory can be allocated"),
    ADS_CODE_1281(0x501, 1281, "ROUTERERR_RESIZEMEMORY", "The size of the  router memory could not be changed"),
    ADS_CODE_1282(0x502, 1282, "ROUTERERR_MAILBOXFULL", "The mailbox has reached the maximum number of possible messages. The current sent message was rejected", "Check the connection between the communication partners"),
    ADS_CODE_1283(0x503, 1283, "ROUTERERR_DEBUGBOXFULL ", "The mailbox has reached the maximum number of possible messages.The sent message will not be displayed in the debug monitor", "Check the connection to the debug monitor"),
    ADS_CODE_1284(0x504, 1284, "ROUTERERR_UNKNOWNPORTTYPE ", "The port type is unknown"),
    ADS_CODE_1285(0x505, 1285, "ROUTERERR_NOTINITIALIZED", "Router is not initialised"),
    ADS_CODE_1286(0x506, 1286, "ROUTERERR_PORTALREADYINUSE ", "The desired port number is already assigned"),
    ADS_CODE_1287(0x507, 1287, "ROUTERERR_NOTREGISTERED ", "Port not registered"),
    ADS_CODE_1288(0x508, 1288, "ROUTERERR_NOMOREQUEUES", "The maximum number of Ports reached"),
    ADS_CODE_1289(0x509, 1289, "ROUTERERR_INVALIDPORT ", "The port is invalid."),
    ADS_CODE_1290(0x50A, 1290, "ROUTERERR_NOTACTIVATED ", "TwinCAT Router not active"),
    ADS_CODE_1291(0x50B, 1291, "ROUTERERR_FRAGMENTBOXFULL"),
    ADS_CODE_1292(0x50C, 1292, "ROUTERERR_FRAGMENTTIMEOUT"),
    ADS_CODE_1293(0x50D, 1293, "ROUTERERR_TOBEREMOVED"),


    // General ADS Error Codes
    ADS_CODE_1792(0x700, 1792, "error class <device error>"),
    ADS_CODE_1793(0x701, 1793, "Service is not supported by server"),
    ADS_CODE_1794(0x702, 1794, "invalid index group"),
    ADS_CODE_1795(0x703, 1795, "invalid index offset"),
    ADS_CODE_1796(0x704, 1796, "reading/writing not permitted"),
    ADS_CODE_1797(0x705, 1797, "parameter size not correct"),
    ADS_CODE_1798(0x706, 1798, "invalid parameter value(s)"),
    ADS_CODE_1799(0x707, 1799, "device is not in a ready state"),
    ADS_CODE_1800(0x708, 1800, "device is busy"),
    ADS_CODE_1801(0x709, 1801, "invalid context (must be in Windows)"),
    ADS_CODE_1802(0x70A, 1802, "out of memory"),
    ADS_CODE_1803(0x70B, 1803, "invalid parameter value(s)"),
    ADS_CODE_1804(0x70C, 1804, "not found (files, ...)"),
    ADS_CODE_1805(0x70D, 1805, "syntax error in command or file"),
    ADS_CODE_1806(0x70E, 1806, "objects do not match"),
    ADS_CODE_1807(0x70F, 1807, "object already exists"),
    ADS_CODE_1808(0x710, 1808, "symbol not found"),
    ADS_CODE_1809(0x711, 1809, "symbol version invalid", "Onlinechange", "Release handle and get a new one"),
    ADS_CODE_1810(0x712, 1810, "server is in invalid state"),
    ADS_CODE_1811(0x713, 1811, "AdsTransMode not supported"),
    ADS_CODE_1812(0x714, 1812, "Notification handle is invalid", "Onlinechange", "Release handle and get a new one"),
    ADS_CODE_1813(0x715, 1813, "Notification client not registered"),
    ADS_CODE_1814(0x716, 1814, "no more notification handles"),
    ADS_CODE_1815(0x717, 1815, "size for watch too big"),
    ADS_CODE_1816(0x718, 1816, "device not initialized"),
    ADS_CODE_1817(0x719, 1817, "device has a timeout"),
    ADS_CODE_1818(0x71A, 1818, "query interface failed"),
    ADS_CODE_1819(0x71B, 1819, "wrong interface required"),
    ADS_CODE_1820(0x71C, 1820, "class ID is invalid"),
    ADS_CODE_1821(0x71D, 1821, "object ID is invalid"),
    ADS_CODE_1822(0x71E, 1822, "request is pending"),
    ADS_CODE_1823(0x71F, 1823, "request is aborted"),
    ADS_CODE_1824(0x720, 1824, "signal warning"),
    ADS_CODE_1825(0x721, 1825, "invalid array index"),
    ADS_CODE_1826(0x722, 1826, "symbol not active", "Onlinechange", "Release handle and get a new one"),
    ADS_CODE_1827(0x723, 1827, "access denied"),
    ADS_CODE_1828(0x724, 1828, "missing license ", "", "Activate license for TwinCAT 3 function"),
    ADS_CODE_1836(0x72c, 1836, "exception occured during system start", "", "Check each device transistions"),
    ADS_CODE_1856(0x740, 1856, "Error class <client error>"),
    ADS_CODE_1857(0x741, 1857, "invalid parameter at service"),
    ADS_CODE_1858(0x742, 1858, "polling list is empty"),
    ADS_CODE_1859(0x743, 1859, "var connection already in use"),
    ADS_CODE_1860(0x744, 1860, "invoke ID in use"),
    ADS_CODE_1861(0x745, 1861, "timeout elapsed", "", "Check ADS routes of sender and receiver and your firewall setting"),
    ADS_CODE_1862(0x746, 1862, "error in win32 subsystem"),
    ADS_CODE_1863(0x747, 1863, "Invalid client timeout value"),
    ADS_CODE_1864(0x748, 1864, "ads-port not opened"),
    ADS_CODE_1872(0x750, 1872, "internal error in ads sync"),
    ADS_CODE_1873(0x751, 1873, "hash table overflow"),
    ADS_CODE_1874(0x752, 1874, "key not found in hash"),
    ADS_CODE_1875(0x753, 1875, "no more symbols in cache"),
    ADS_CODE_1876(0x754, 1876, "invalid response received"),
    ADS_CODE_1877(0x755, 1877, "sync port is locked"),


    // RTime Error Codes
    ADS_CODE_4096(0x1000, 4096, "RTERR_INTERNAL ", "Internal fatal error in the TwinCAT real-time system"),
    ADS_CODE_4097(0x1001, 4097, "RTERR_BADTIMERPERIODS", "Timer value not vaild"),
    ADS_CODE_4098(0x1002, 4098, "RTERR_INVALIDTASKPTR", "Task pointer has the invalid value ZERO"),
    ADS_CODE_4099(0x1003, 4099, "RTERR_INVALIDSTACKPTR", "Task stack pointer has the invalid value ZERO"),
    ADS_CODE_4100(0x1004, 4100, "RTERR_PRIOEXISTS", "The demand task priority is already assigned"),
    ADS_CODE_4101(0x1005, 4101, "RTERR_NOMORETCB", "No more free TCB (Task Control Block) available. Maximum number of TCBs is 64"),
    ADS_CODE_4102(0x1006, 4102, "RTERR_NOMORESEMAS", "No more free semaphores available. Maximum number of semaphores is 64"),
    ADS_CODE_4103(0x1007, 4103, "RTERR_NOMOREQUEUES", "No more free queue available. Maximum number of queue is 64"),
    ADS_CODE_4104(0x1008, 4104, "TwinCAT ", "reserved."),
    ADS_CODE_4105(0x1009, 4105, "TwinCAT ", "reserved."),
    ADS_CODE_4106(0x100A, 4106, "TwinCAT ", "reserved."),
    ADS_CODE_4107(0x100B, 4107, "TwinCAT ", "reserved."),
    ADS_CODE_4108(0x100C, 4108, "TwinCAT ", "reserved."),
    ADS_CODE_4109(0x100D, 4109, "RTERR_EXTIRQALREADYDEF", "An external synchronisation interrupt is already applied"),
    ADS_CODE_4110(0x100E, 4110, "RTERR_EXTIRQNOTDEF", "No external synchronsiation interrupt applied"),
    ADS_CODE_4111(0x100F, 4111, "RTERR_EXTIRQINSTALLFAILED", "The apply of the external synchronisation interrupt failed"),
    ADS_CODE_4112(0x1010, 4112, "RTERR_IRQLNOTLESSOREQUAL", "Call of a service function in the wrong context"),
    ADS_CODE_4119(0x1017, 4119, "RTERR_VMXNOTSUPPORTED", "Intel VT-x extension is not supported."),
    ADS_CODE_4120(0x1018, 4120, "RTERR_VMXDISABLED", "Intel VT-x extension is not enabled in BIOS."),
    ADS_CODE_4121(0x1019, 4121, "RTERR_VMXCONTROLSMISSING", "Missing feature in Intel VT-x extension."),
    ADS_CODE_4122(0x101A, 4122, "RTERR_VMXENABLEFAILS", "Enabling Intel VT-x fails."),


    // TCP Winsock Error Codes
    ADS_CODE_10060(0x274c, 10060, "A socket operation was attempted to an unreachable host", "Host unreachable", "Check network connection via ping"),
    ADS_CODE_10061(0x274d, 10061, "A connection attempt failed because the connected party did not properly respond after a period of time,or established connection failed because connected host has failed to respond.", "Host unreachable", "Check network connection via ping"),
    ADS_CODE_10065(0x2751, 10065, "No connection could be made because the target machine actively refused it"),
    UNKNOWN(0, 0, "", "", "");

    private final long hex;
    private final long dec;
    private final String description;
    private final String possibleCauses;
    private final String solution;

    AdsReturnCode(long hex, long dec, String description) {
        this(hex, dec, description, "");
    }

    AdsReturnCode(long hex, long dec, String description, String possibleCauses) {
        this(hex, dec, description, possibleCauses, "");
    }

    AdsReturnCode(long hex, long dec, String description, String possibleCauses, String solution) {
        this.hex = hex;
        this.dec = dec;
        if (hex != dec) {
            throw new IllegalArgumentException("hex " + hex + " is different from dec " + dec);
        }
        this.description = requireNonNull(description);
        this.possibleCauses = possibleCauses;
        this.solution = solution;
    }

    // TODO: optimize with map
    public static AdsReturnCode of(long hex) {
        for (AdsReturnCode adsReturnCode : values()) {
            if (adsReturnCode.hex == hex) {
                return adsReturnCode;
            }
        }
        return UNKNOWN;
    }

    public static AdsReturnCode of(String name) {
        return valueOf(name);
    }

    public static AdsReturnCode ofInt(String value) {
        return of(Long.parseLong(value));
    }

    public long getHex() {
        return hex;
    }

    public long getDec() {
        return dec;
    }

    public String getDescription() {
        return description;
    }

    public String getPossibleCauses() {
        return possibleCauses;
    }

    public String getSolution() {
        return solution;
    }

    @Override
    public String toString() {
        return "AdsReturnCode{" +
            "hex=" + String.format("0x%4s", toHexString(hex)).replace(' ', '0') +
            ", dec=" + dec +
            ", description='" + description + '\'' +
            ", possibleCauses='" + possibleCauses + '\'' +
            ", solution='" + solution + '\'' +
            '}';
    }
}
