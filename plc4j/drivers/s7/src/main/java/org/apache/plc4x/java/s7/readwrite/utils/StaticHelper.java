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
package org.apache.plc4x.java.s7.readwrite.utils;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.ByteBufUtil;
import io.netty.buffer.Unpooled;
import org.apache.commons.lang3.NotImplementedException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.s7.events.S7AlarmEvent;
import org.apache.plc4x.java.s7.events.S7ModeEvent;
import org.apache.plc4x.java.s7.events.S7SysEvent;
import org.apache.plc4x.java.s7.readwrite.DataTransportErrorCode;
import org.apache.plc4x.java.s7.readwrite.ModeTransitionType;
import org.apache.plc4x.java.s7.utils.S7DiagnosticEventId;
import org.apache.plc4x.java.spi.codegen.WithOption;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;

import java.nio.charset.Charset;
import java.nio.charset.StandardCharsets;
import java.time.Duration;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.time.temporal.ChronoUnit;
import java.util.HashMap;
import java.util.Map;
import java.util.logging.Level;
import java.util.logging.Logger;
import java.util.regex.Matcher;
import java.util.regex.Pattern;
import org.apache.plc4x.java.s7.readwrite.DataTransportSize;

/**
 * +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
 * |15|14|13|12|11|10| 9| 8| 7| 6| 5| 4| 3| 2| 1|
 * +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
 * \__________/\__________/\____________________/
 * Module      Number of   Number of the partial
 * class       the partial list
 * list
 * extract
 *
 * <b>Module Class:</b>
 * +--------------+-----------------+
 * | Module class | Coding (Binary) |
 * +--------------|-----------------+
 * |     CPU      |      0000       |
 * +--------------|-----------------+
 * |     IM       |      0100       |
 * +--------------|-----------------+
 * |     FM       |      1000       |
 * +--------------|-----------------+
 * |     CP       |      1100       |
 * +--------------|-----------------+
 *
 * <b>Possible SSL Partial Lists:</b>
 * +-----------------------------------------------------------|--------------+
 * | Module class                                              |    SSL-ID    |
 * +-----------------------------------------------------------|--------------+
 * | Module identification                                     |    16#xy11   |
 * +-----------------------------------------------------------|--------------+
 * | CPU characteristics                                       |    16#xy12   |
 * +-----------------------------------------------------------|--------------+
 * | User memory areas                                         |    16#xy13   |
 * +-----------------------------------------------------------|--------------+
 * | System areas                                              |    16#xy14   |
 * +-----------------------------------------------------------|--------------+
 * | Block types                                               |    16#xy15   |
 * +-----------------------------------------------------------|--------------+
 * | Interrupt status                                          |    16#xy22   |
 * +-----------------------------------------------------------|--------------+
 * | Assignment between process image partitions and OBs       |    16#xy25   |
 * +-----------------------------------------------------------|--------------+
 * | Communication status data                                 |    16#xy32   |
 * +-----------------------------------------------------------|--------------+
 * | H CPU group information                                   |    16#xy71   |
 * +-----------------------------------------------------------|--------------+
 * | Status of the module LEDs                                 |    16#xy74   |
 * +-----------------------------------------------------------|--------------+
 * | Switched DP slaves in the H-system                        |    16#xy75   |
 * +-----------------------------------------------------------|--------------+
 * | Module status information                                 |    16#xy91   |
 * +-----------------------------------------------------------|--------------+
 * | Rack / station status information                         |    16#xy92   |
 * +-----------------------------------------------------------|--------------+
 * | Rack / station status information                         |    16#xy94   |
 * +-----------------------------------------------------------|--------------+
 * | Extended DP master system / PROFINET IO system information|    16#xy95   |
 * +-----------------------------------------------------------|--------------+
 * | Module status information, PROFINET IO and PROFIBUS DP    |    16#xy96   |
 * +-----------------------------------------------------------|--------------+
 * | Tool changer information (PROFINET IO)                    |    16#xy9C   |
 * +-----------------------------------------------------------|--------------+
 * | Diagnostic buffer of the CPU                              |    16#xyA0   |
 * +-----------------------------------------------------------|--------------+
 * | Module diagnostic information (data record 0)             |    16#xyB1   |
 * +-----------------------------------------------------------|--------------+
 * | Module diagnostic information (data record 1),            |    16#xyB2   |
 * | geographical address                                      |              |
 * +-----------------------------------------------------------|--------------+
 * | Module diagnostic information (data record 1),            |    16#xyB3   |
 * | logical address                                           |              |
 * +-----------------------------------------------------------|--------------+
 * | Diagnostic data of a DP slave                             |    16#xyB4   |
 * +-----------------------------------------------------------|--------------+
 */

public class StaticHelper {

    public enum OB {
        FREE_CYC(0X0000, "OB1 Free cycle"),

        //Time of day
        TOD_INT0(0X000A, "OB10 Time of day interrupt"),
        TOD_INT1(0X000B, "OB11 Time of day interrupt"),
        TOD_INT2(0X000C, "OB12 Time of day interrupt"),
        TOD_INT3(0X000D, "OB13 Time of day interrupt"),
        TOD_INT4(0X000E, "OB14 Time of day interrupt"),
        TOD_INT5(0X000F, "OB15 Time of day interrupt"),
        TOD_INT6(0X0010, "OB16 Time of day interrupt"),
        TOD_INT7(0X0011, "OB17 Time of day interrupt"),

        //Time delay
        DEL_INT0(0X0014, "OB20 Time delay interrupt"),
        DEL_INT1(0X0015, "OB21 Time delay interrupt"),
        DEL_INT2(0X0016, "OB22 Time delay interrupt"),
        DEL_INT3(0X0017, "OB23 Time delay interrupt"),

        //Cyclic
        CYC_INT0(0X001E, "OB30 Cyclic interrupt"),
        CYC_INT1(0X001F, "OB31 Cyclic interrupt"),
        CYC_INT2(0X0020, "OB32 Cyclic interrupt"),
        CYC_INT3(0X0021, "OB33 Cyclic interrupt"),
        CYC_INT4(0X0022, "OB34 Cyclic interrupt"),
        CYC_INT5(0X0023, "OB35 Cyclic interrupt"),
        CYC_INT6(0X0024, "OB36 Cyclic interrupt"),
        CYC_INT7(0X0025, "OB37 Cyclic interrupt"),
        CYC_INT8(0X0026, "OB38 Cyclic interrupt"),

        //Hardware interrupts
        HW_INT0(0X0028, "OB40 Hardware interrupt"),
        HW_INT1(0X0029, "OB41 Hardware interrupt"),
        HW_INT2(0X002A, "OB42 Hardware interrupt"),
        HW_INT3(0X002B, "OB43 Hardware interrupt"),
        HW_INT4(0X002C, "OB44 Hardware interrupt"),
        HW_INT5(0X002D, "OB45 Hardware interrupt"),
        HW_INT6(0X002E, "OB46 Hardware interrupt"),
        HW_INT7(0X002F, "OB47 Hardware interrupt"),

        //Startup
        BACKGROUND(0X005A, "OB90 Background"),
        COMPLETE_RESTART(0X0064, "OB100 Startup"),
        RESTART(0X0064, "OB101 Background"),
        COLD_RESTART(0X005A, "OB101 Background"),

        //Fault interrupts
        CYC_FLT(0X0051, "OB80 Time execution error interrupt"),
        PS_FLT(0X0051, "OB81 Power supply interrupt"),
        IO_FLT1(0X0052, "OB82 Module diagnostic interrupt"),
        IO_FLT2(0X0053, "OB83 Module change interrupt"),
        CPU_FLT(0X0054, "OB84 CPU hardware error interrupt"),
        OBNL_FLT(0X0055, "OB85 Program execution error interrupt"),
        RACK_FLT(0X0056, "OB86 Rack fault interrupt"),
        COMM_FLT(0X0057, "OB87 Communication error interrupt"),
        BREAKUP_ERR(0X0058, "OB88 Process interrupt"),
        SYNC_ERR(0X0079, "OB120 Synchronous error interrupt"),
        PROG_ERR(0X0079, "OB121 Program error interrupt"),
        MOD_ERR(0X007A, "OB122 Module error interrupt");

        private final int code;
        private final String description;

        private static final Map<Integer, OB> map;

        static {
            map = new HashMap<>();
            for (OB obid : OB.values()) {
                map.put(obid.code, obid);
            }
        }

        OB(final int code, final String description) {
            this.code = code;
            this.description = description;
        }

        public int getCode() {
            return code;
        }

        public String getDescription() {
            return description;
        }

        public static OB valueOf(int code) {
            return map.get(code);
        }


    }


    public enum MODULE {
        CPU(0X00),
        IM(0X04),
        FM(0X80),
        CP(0XC0);

        private final int code;

        MODULE(final int code) {
            this.code = code;
        }

        public int getCode() {
            return code;
        }
    }

    public enum LED_ID {
        SF(0X0001, "Group error"),
        INTF(0X0002, "Internal error"),
        EXTF(0X0003, "External error"),
        RUN(0X0004, "RUN"),
        STOP(0X0005, "STOP"),
        FRCE(0X0006, "Force"),
        CRST(0X0007, "Cold restart"),
        BAF(0X0008, "Battery fault"),
        USR(0X0009, "User defined"),
        USR1(0X000A, "User defined"),
        BUS1F(0X000B, "Bus error interface 1"),
        BUS2F(0X000C, "Bus error interface 2"),
        REDF(0X000D, "Redundancy error"),
        MSTR(0X000E, "Master"),
        RACK0(0X000F, "Rack number 0"),
        RACK1(0X0010, "Rack number 1"),
        RACK2(0X0011, "Rack number 2"),
        IFM1F(0X0012, "Interface error interface module 1"),
        IFM2F(0X0013, "Interface error interface module 2"),
        BUS3F(0X0014, "Bus error interface 3"),
        MAINT(0X0015, "Maintenance demand"),
        DC24V(0X0016, "DC24V"),
        BUS5F(0X0017, "Bus error interface 5"),
        BUS8F(0X0018, "Bus error interface 8"),
        IF(0X0080, "Init failure"),
        UF(0X0081, "User failure"),
        MF(0X0082, "Monitoring failure"),
        CF(0X0083, "Communication failure"),
        TF(0X0084, "Task failure"),
        APPL_STATE_RED(0X00EC, "APPL_STATE_RED"),
        APPL_STATE_GREEN(0X00ED, "APPL_STATE_GREEN");

        private final int code;
        private final String description;

        private static final Map<Integer, LED_ID> map;

        static {
            map = new HashMap<>();
            for (LED_ID ledid : LED_ID.values()) {
                map.put(ledid.code, ledid);
            }
        }

        LED_ID(final int code, final String description) {
            this.code = code;
            this.description = description;
        }

        public int getCode() {
            return code;
        }

        public String getDescription() {
            return description;
        }

        public static LED_ID valueOf(int code) {
            return map.get(code);
        }


    }

    /**
     *
     */
    public enum CPU_CHARACTERISTICS {
        CH_0x0000(0X0000, "MC7 processing unit (group with index 0000)"),
        CH_0x0001(0X0001, "MC7 processing generating code"),
        CH_0x0002(0X0002, "MC7 interpreter"),

        CH_0x0100(0X0100, "Time system (group with index 0100)"),
        CH_0x0101(0X0101, "1 ms resolution"),
        CH_0x0102(0X0102, "10 ms resolution"),
        CH_0x0103(0X0103, "No real time clock"),
        CH_0x0104(0X0104, "BCD time-of-day format"),
        CH_0x0105(0X0105, "All time-of-day functions"),
        CH_0x0106(0X0106, "SFC 78 \"OB_RT\" is available"),

        CH_0x0200(0X0200, "System response (group with index 0200)"),
        CH_0x0201(0X0201, "Capable of multiprocessor mode"),
        CH_0x0202(0X0202, "Cold restart, warm restart and hot restart possible"),
        CH_0x0203(0X0203, "Cold restart and hot restart possible"),
        CH_0x0204(0X0204, "Warm restart and hot restart possible"),
        CH_0x0205(0X0205, "Only warm restart possible"),
        CH_0x0206(0X0206, "New distributed I/O configuration is possible during\n" +
            "RUN by using predefined resources"),
        CH_0x0207(0X0207, "H-CPU in stand-alone mode: New distributed I/O configuration\n" +
            "is possible during RUN by using predefined resources"),
        CH_0x0208(0X0208, "For taking motion control functionality into account"),

        CH_0x0300(0X0300, "MC7 Language description of the CPU (group with index 0300)"),
        CH_0x0301(0X0301, "Reserved"),
        CH_0x0302(0X0302, "All 32 bit fixed-point instructions"),
        CH_0x0303(0X0303, "All floating-point instructions"),
        CH_0x0304(0X0304, "sin, asin, cos, acos, tan, atan, sqr, sqrt, ln, exp"),
        CH_0x0305(0X0305, "Accumulator 3/accumulator 4 with corresponding instructions\n" +
            "(ENT,PUSH,POP,LEAVE)"),
        CH_0x0306(0X0306, "Master Control Relay instructions"),
        CH_0x0307(0X0307, "Address register 1 exists with corresponding instructions"),
        CH_0x0308(0X0308, "Address register 2 exists with corresponding instructions"),
        CH_0x0309(0X0309, "Operations for area-crossing addressing"),
        CH_0x030A(0X030A, "Operations for area-internal addressing"),
        CH_0x030B(0X030B, "All memory-indirect addressing instructions for bit memory (M)"),
        CH_0x030C(0X030C, "All memory-indirect addressing instructions for data blocks (DB)"),
        CH_0x030D(0X030D, "All memory-indirect addressing instructions for data blocks (DI)"),
        CH_0x030E(0X030E, "All memory-indirect addressing instructions for local data (L)"),
        CH_0x030F(0X030F, "All instructions for parameter transfer in FCs"),
        CH_0x0310(0X0310, "Memory bit edge instructions for process image input (I)"),
        CH_0x0311(0X0311, "Memory bit edge instructions for process image output (Q)"),
        CH_0x0312(0X0312, "Memory bit edge instructions for bit memory (M)"),
        CH_0x0313(0X0313, "Memory bit edge instructions for data blocks (DB)"),
        CH_0x0314(0X0314, "Memory bit edge instructions for data blocks (DI)"),
        CH_0x0315(0X0315, "Memory bit edge instructions for local data (L)"),
        CH_0x0316(0X0316, "Dynamic evaluation of the FC bit"),
        CH_0x0317(0X0317, "Dynamic local data area with the corresponding instructions"),
        CH_0x0318(0X0318, "Reserved"),
        CH_0x0319(0X0319, "Reserved"),

        CH_0x0401(0X0401, "SFC 87 \"C_DIAG\" is available"),
        CH_0x0402(0X0402, "SFC 88 \"C_CNTRL\" is available)");

        private final int code;
        private final String description;

        private static final Map<Integer, CPU_CHARACTERISTICS> map;

        static {
            map = new HashMap<>();
            for (CPU_CHARACTERISTICS cpuc : CPU_CHARACTERISTICS.values()) {
                map.put(cpuc.code, cpuc);
            }
        }

        CPU_CHARACTERISTICS(final int code, final String description) {
            this.code = code;
            this.description = description;
        }

        public int getCode() {
            return code;
        }

        public String getDescription() {
            return description;
        }

        public static CPU_CHARACTERISTICS valueOf(int code) {
            return map.get(code);
        }


    }

    /**
     *
     */
    public enum SZL {
        ID_0x0011(0x0011, "Module identification.") {
            @Override
            public StringBuilder execute(ByteBuf data) {
                return ID_0xXY11(data);
            }

        },
        ID_0x0012(0x0012, "CPU characteristics .") {
            @Override
            public StringBuilder execute(ByteBuf data) {
                return ID_0xXY12(data);
            }

        },
        ID_0x0013(0x0013, "User memory areas.") {
            @Override
            public StringBuilder execute(ByteBuf data) {
                return ID_0xXY13(data);
            }

        },
        ID_0x0014(0x0014, "System areas.") {
            @Override
            public StringBuilder execute(ByteBuf data) {
                return ID_0xXY14(data);
            }

        },
        ID_0x0015(0x0015, "Block types.") {
            @Override
            public StringBuilder execute(ByteBuf data) {
                return ID_0xXY15(data);
            }

        },
        ID_0x001C(0x001C, "Component identification.") {
            @Override
            public StringBuilder execute(ByteBuf data) {
                return ID_0xXY1C(data);
            }

        },
        ID_0x0022(0x0022, "Interrupt status.") {
            @Override
            public StringBuilder execute(ByteBuf data) {
                return ID_0xXY22(data);
            }

        },
        ID_0x0025(0x0025, "Assignment between process image partitions and OBs.") {
            @Override
            public StringBuilder execute(ByteBuf data) {
                return ID_0xXY25(data);
            }

        },
        ID_0x0032(0x0032, "Communication status data.") {
            @Override
            public StringBuilder execute(ByteBuf data) {
                return ID_0xXY32(data);
            }

        },
        ID_0x0071(0x0071, "H CPU group information.") {
            @Override
            public StringBuilder execute(ByteBuf data) {
                return ID_0xXY71(data);
            }

        },
        ID_0x0074(0x0074, "Status of the module LEDs.") {
            @Override
            public StringBuilder execute(ByteBuf data) {
                return ID_0xXY74(data);
            }

        },
        ID_0x0075(0x0075, "Switched DP slaves in the H-system.") {
            @Override
            public StringBuilder execute(ByteBuf data) {
                return ID_0xXY75(data);
            }

        },
        ID_0x0090(0x0090, "DP Master System Information.") {
            @Override
            public StringBuilder execute(ByteBuf data) {
                return ID_0xXY90(data);
            }

        },
        ID_0x0091(0x0091, "Module status information.") {
            @Override
            public StringBuilder execute(ByteBuf data) {
                return ID_0xXY91(data);
            }

        },
        ID_0x0092(0x0092, "Rack / station status information.") {
            @Override
            public StringBuilder execute(ByteBuf data) {
                return ID_0xXY92(data);
            }

        },
        ID_0x0094(0x0094, "Rack / station status information.") {
            @Override
            public StringBuilder execute(ByteBuf data) {
                return ID_0xXY94(data);
            }

        },
        ID_0x0095(0x0095, "Extended DP master system / PROFINET IO system information.") {
            @Override
            public StringBuilder execute(ByteBuf data) {
                return ID_0xXY95(data);
            }

        },
        ID_0x0096(0x0096, "Module status information, PROFINET IO and PROFIBUS DP.") {
            @Override
            public StringBuilder execute(ByteBuf data) {
                return ID_0xXY96(data);
            }

        },
        ID_0x009C(0x009C, "Tool changer information (PROFINET IO).") {
            @Override
            public StringBuilder execute(ByteBuf data) {
                return ID_0xXY9C(data);
            }

        },
        ID_0x00A0(0x00A0, "Diagnostic buffer of the CPU.") {
            @Override
            public StringBuilder execute(ByteBuf data) {
                return ID_0xXYA0(data);
            }

        },
        ID_0x00B1(0x00B1, "Module diagnostic information (data record 0) .") {
            @Override
            public StringBuilder execute(ByteBuf data) {
                return ID_0xXYB1(data);
            }

        },
        ID_0x00B2(0x00B2, "Module diagnostic information (data record 1),geographical address.") {
            @Override
            public StringBuilder execute(ByteBuf data) {
                return ID_0xXYB2(data);
            }

        },
        ID_0x00B3(0x00B3, "Module diagnostic information (data record 1),logical address.") {
            @Override
            public StringBuilder execute(ByteBuf data) {
                return ID_0xXYB3(data);
            }

        },
        ID_0x00B4(0x00B4, "Diagnostic data of a DP slave.") {
            @Override
            public StringBuilder execute(ByteBuf data) {
                return ID_0xXYB4(data);
            }

        };

        private final int code;
        private final String description;

        private static final Map<Integer, SZL> map;

        static {
            map = new HashMap<>();
            for (SZL subszl : SZL.values()) {
                map.put(subszl.code, subszl);
            }
        }

        SZL(final int code, final String description) {
            this.code = code;
            this.description = description;
        }

        public int getCode() {
            return code;
        }

        public String getDescription() {
            return description;
        }

        public static SZL valueOf(int code) {
            return map.get(code);
        }

        public abstract StringBuilder execute(ByteBuf data);

        private static StringBuilder ID_0xXY11(ByteBuf data) {
            StringBuilder sb = new StringBuilder();
            int ssl_length = data.readShort();
            int ssl_count = data.readShort();
            try {
                sb.append("SSL partial list length in bytes: ").append(ssl_length).append("\r\n");
                sb.append("SSL partial list count: ").append(ssl_count).append("\r\n");
                for (int i = 1; i <= ssl_count; i++) {
                    sb.append("Data Record: ").append(i).append("\r\n");

                    sb.append("Index: ").append(data.readShort()).append("\r\n");
                    byte[] bytestr = new byte[20];
                    data.readBytes(bytestr, 0, 20);
                    sb.append("Mlfb (Order number): ").append(new String(bytestr)).append("\r\n");
                    sb.append("BGTyp (Moduler type ID): ").append(data.readShort()).append("\r\n");
                    sb.append("Ausbg (Version of the module/hardware version): ").append(data.readShort()).append("\r\n");
                    sb.append("Ausbe (Firmware version): ").append(data.readShort()).append("\r\n");
                }
            } catch (Exception ex) {
                sb.append(ex);
            }
            return sb;
        }

        private static StringBuilder ID_0xXY12(ByteBuf data) {
            StringBuilder sb = new StringBuilder();
            int ssl_length = data.readShort();
            int ssl_count = data.readShort();
            try {
                sb.append("SSL partial list length in bytes: ").append(ssl_length).append("\r\n");
                sb.append("SSL partial list count: ").append(ssl_count).append("\r\n");
                for (int i = 1; i <= ssl_count; i++) {
                    sb.append("Data Record: ").append(i).append("\r\n");
                    int code = data.readShort();
                    sb.append(CPU_CHARACTERISTICS.valueOf(code)).append(": ").
                        append(CPU_CHARACTERISTICS.valueOf(code).getDescription()).append("\r\n");
                }
            } catch (Exception ex) {
                sb.append(ex);
            }
            return sb;
        }

        private static StringBuilder ID_0xXY13(ByteBuf data) {
            StringBuilder sb = new StringBuilder();
            int ssl_length = data.readShort();
            int ssl_count = data.readShort();
            try {
                sb.append("SSL partial list length in bytes: ").append(ssl_length).append("\r\n");
                sb.append("SSL partial list count: ").append(ssl_count).append("\r\n");
                for (int i = 1; i <= ssl_count; i++) {
                    sb.append("Data Record: ").append(i).append("\r\n");
                    sb.append("Index: ").append(data.readShort()).append("\r\n");
                    int code = data.readShort();
                    int size = data.readInt();
                    int mode = data.readShort();
                    int granu = data.readShort();
                    int ber1 = data.readInt();
                    int be1egt1 = data.readInt();
                    int block1 = data.readInt();
                    int ber2 = data.readInt();
                    int be1egt2 = data.readInt();
                    int block2 = data.readInt();

                    switch (code) {
                        case 0x0001:
                            sb.append("Code (Memory Type): ").append("Volatile memory (RAM)\r\n");
                            break;
                        case 0x0002:
                            sb.append("Code (Memory Type): ").append("Non-volatile memory (FEPROM)\r\n");
                            break;
                        case 0x0003:
                            sb.append("Code (Memory Type): ").append("mixed memory (RAM + FEPROM)\r\n");
                            break;
                    }
                    sb.append("Size (Total size of the selected memory (total of area 1 and area 2)): ").append(size).append("\r\n");

                    sb.append("Mode (Volatile memory area): ").append(((mode & 0x0001) != 0)).append("\r\n");
                    sb.append("Mode (Non-volatile memory area): ").append(((mode & 0x0002) != 0)).append("\r\n");
                    sb.append("Mode (Mixed memory area): ").append(((mode & 0x0004) != 0)).append("\r\n");
                    sb.append("Mode (Code and data separate): ").append(((mode & 0x0008) != 0)).append("\r\n");
                    sb.append("Mode (Code and data together): ").append(((mode & 0x0010) != 0)).append("\r\n");

                    sb.append("Granu (Always zero): ").append(granu).append("\r\n");

                    sb.append("Ber1 (Size of the volatile memory area in bytes): ").append(ber1).append("\r\n");
                    sb.append("Belegt1 (Size of the volatile memory area being used): ").append(be1egt1).append("\r\n");
                    sb.append("Block1 (Largest free block in the volatile memory area): ").append(block1).append("\r\n");

                    sb.append("Ber2 (Size of the non-volatile memory area in bytes): ").append(ber2).append("\r\n");
                    sb.append("Belegt2 (Size of the non-volatile memory area being used): ").append(be1egt2).append("\r\n");
                    sb.append("Block2 (Largest free block in the non-volatile memory area): ").append(block2).append("\r\n");

                }
            } catch (Exception ex) {
                sb.append(ex);
            }
            return sb;
        }

        private static StringBuilder ID_0xXY14(ByteBuf data) {
            StringBuilder sb = new StringBuilder();
            int ssl_length = data.readShort();
            int ssl_count = data.readShort();
            try {
                sb.append("SSL partial list length in bytes: ").append(ssl_length).append("\r\n");
                sb.append("SSL partial list count: ").append(ssl_count).append("\r\n");
                for (int i = 1; i <= ssl_count; i++) {
                    int index = data.readShort();
                    int code = data.readShort();
                    int quantity = data.readShort();
                    int reman = data.readShort();

                    sb.append("Data Record: ").append(i).append("\r\n");
                    sb.append("Index: ").append(index).append(" System area: ");
                    switch (index) {
                        case 0x0001:
                            sb.append("PII (number in bytes)").append("\r\n");
                            break;
                        case 0x0002:
                            sb.append("PIQ (number in bytes))").append("\r\n");
                            break;
                        case 0x0003:
                            sb.append("memory (number in bits)").append("\r\n");
                            break;
                        case 0x0004:
                            sb.append("timers (number)").append("\r\n");
                            break;
                        case 0x0005:
                            sb.append("counter (number)").append("\r\n");
                            break;
                        case 0x0006:
                            sb.append("number of bytes in the logical address area)").append("\r\n");
                            break;
                        case 0x0007:
                            sb.append("local data (entire local data area of the CPU in bytes)").append("\r\n");
                            break;
                        case 0x0008:
                            sb.append("memory (number in bytes)").append("\r\n");
                            break;
                        case 0x0009:
                            sb.append("local data (entire local data area of the CPU in Kbytes)").append("\r\n");
                            break;
                    }
                    sb.append("Code: ").append(code).append(" Memory type: ");
                    switch (code) {
                        case 0x0001:
                            sb.append("volatile memory (RAM)").append("\r\n");
                            break;
                        case 0x0002:
                            sb.append("non-volatile memory (FEPROM)").append("\r\n");
                            break;
                        case 0x0003:
                            sb.append("mixed memory (RAM and FEPROM)").append("\r\n");
                            break;
                    }
                    sb.append("Quantity: ").append(quantity).append("\r\n");
                    sb.append("Reman: ").append(reman).append("\r\n");
                }
            } catch (Exception ex) {
                sb.append(ex);
            }
            return sb;
        }

        private static StringBuilder ID_0xXY15(ByteBuf data) {
            StringBuilder sb = new StringBuilder();
            int ssl_length = data.readShort();
            int ssl_count = data.readShort();
            try {
                sb.append("SSL partial list length in bytes: ").append(ssl_length).append("\r\n");
                sb.append("SSL partial list count: ").append(ssl_count).append("\r\n");

                for (int i = 1; i <= ssl_count; i++) {
                    sb.append("Data Record: ").append(i).append("\r\n");
                    int type = data.readShort();
                    switch (type) {
                        case 0x0800:
                            sb.append("Block type: OB\r\n");
                            break;
                        case 0x0A00:
                            sb.append("Block type: DB\r\n");
                            break;
                        case 0x0B00:
                            sb.append("Block type: SDB\r\n");
                            break;
                        case 0x0C00:
                            sb.append("Block type: FC\r\n");
                            break;
                        case 0x0E00:
                            sb.append("Block type: FB\r\n");
                            break;
                        default:
                    }
                    sb.append("Maximum number of blocks of the type: ").append(data.readShort()).append("\r\n");
                    sb.append("Maximum total size of the object to be loaded in Kbytes: ").append(data.readShort()).append("\r\n");
                    sb.append("Maximum length of the work memory part of a block in bytes: ").append(data.readInt()).append("\r\n");
                }
            } catch (Exception ex) {
                sb.append(ex);
            }
            return sb;
        }

        private static StringBuilder ID_0xXY1C(ByteBuf data) {
            StringBuilder sb = new StringBuilder();
            int ssl_length = data.readShort();
            int ssl_count = data.readShort();
            try {
                sb.append("SSL partial list length in bytes: ").append(ssl_length).append("\r\n");
                sb.append("SSL partial list count: ").append(ssl_count).append("\r\n");

                for (int i = 1; i <= ssl_count; i++) {
                    sb.append("Data Record: ").append(i).append("\r\n");
                    int index = data.getShort(data.readerIndex());
                    int index_b0 = data.readByte();
                    int index_b1 = data.readByte();
                    sb.append("Index: ").append(index).append("\r\n");
                    sb.append("Index b1: ").append(index_b1).append("\r\n");

                    switch (index_b1) {
                        case 0x01:
                        case 0x02: {
                            byte[] strbyte = new byte[24];
                            data.readBytes(strbyte, 0, 24);
                            sb.append("Name: ").append(new String(strbyte)).append("\r\n");
                            ByteBuf data2 = data.readSlice(8);
                            sb.append("Reserved: ").append(ByteBufUtil.hexDump(data2)).append("\r\n");
                        }
                        ;
                        break;
                        case 0x03: {
                            byte[] strbyte = new byte[32];
                            data.readBytes(strbyte, 0, 32);
                            sb.append("Tag: ").append(new String(strbyte)).append("\r\n");
                        }
                        break;
                        case 0x04: {
                            byte[] strbyte = new byte[26];
                            data.readBytes(strbyte, 0, 26);
                            sb.append("Copyright: ").append(new String(strbyte)).append("\r\n");
                            ByteBuf data2 = data.readSlice(6);
                            sb.append("Reserved: ").append(ByteBufUtil.hexDump(data2)).append("\r\n");
                        }
                        break;
                        case 0x05: {
                            byte[] strbyte = new byte[24];
                            data.readBytes(strbyte, 0, 24);
                            sb.append("Serial: ").append(new String(strbyte)).append("\r\n");
                            ByteBuf data2 = data.readSlice(8);
                            sb.append("Reserved: ").append(ByteBufUtil.hexDump(data2)).append("\r\n");
                        }
                        break;
                        case 0x06: {

                        }
                        break;
                        case 0x07: {
                            byte[] strbyte = new byte[32];
                            data.readBytes(strbyte, 0, 32);
                            sb.append("Cpu type: ").append(new String(strbyte)).append("\r\n");
                        }
                        break;
                        case 0x08: {
                            byte[] strbyte = new byte[32];
                            data.readBytes(strbyte, 0, 32);
                            sb.append("Serial mc/mmc: ").append(new String(strbyte)).append("\r\n");
                        }
                        break;
                        case 0x09: {
                            sb.append("Manufacturer ID: ").append(data.readShort()).append("\r\n");
                            sb.append("Profile ID: ").append(data.readShort()).append("\r\n");
                            sb.append("Profile specific type: ").append(data.readShort()).append("\r\n");
                            ByteBuf data2 = data.readSlice(26);
                            sb.append("Reserved: ").append(ByteBufUtil.hexDump(data2)).append("\r\n");
                        }
                        break;
                        case 0x0A: {
                            byte[] strbyte = new byte[26];
                            data.readBytes(strbyte, 0, 26);
                            sb.append("OEM copyright: ").append(new String(strbyte)).append("\r\n");
                            sb.append("OEM ID: ").append(data.readShort()).append("\r\n");
                            sb.append("OEM add ID: ").append(data.readInt()).append("\r\n");
                        }
                        break;
                        case 0x0B: {
                            byte[] strbyte = new byte[32];
                            data.readBytes(strbyte, 0, 32);
                            sb.append("Location ID: ").append(new String(strbyte)).append("\r\n");
                        }
                        break;
                        case 0x0C: {
                            byte[] strbyte = new byte[10];
                            data.readBytes(strbyte, 0, 10);
                            sb.append("Order number sync module 1: ").append(new String(strbyte)).append("\r\n");
                            data.readShort();
                            byte[] strbyte2 = new byte[2];
                            data.readBytes(strbyte2, 0, 2);
                            sb.append("Product version: ").append(new String(strbyte2)).append("\r\n");
                            data.readByte();
                            ByteBuf data2 = data.readSlice(17);
                            sb.append("Vendor serial: ").append(ByteBufUtil.hexDump(data2)).append("\r\n");
                        }
                        break;
                        case 0x0D: {
                            byte[] strbyte = new byte[10];
                            data.readBytes(strbyte, 0, 10);
                            sb.append("Order number sync module 2: ").append(new String(strbyte)).append("\r\n");
                            data.readShort();
                            byte[] strbyte2 = new byte[2];
                            data.readBytes(strbyte2, 0, 2);
                            sb.append("Product version: ").append(new String(strbyte2)).append("\r\n");
                            data.readByte();
                            ByteBuf data2 = data.readSlice(17);
                            sb.append("Vendor serial: ").append(ByteBufUtil.hexDump(data2)).append("\r\n");
                        }
                        break;
                        case 0x0E: {
                            byte[] strbyte = new byte[18];
                            data.readBytes(strbyte, 0, 18);
                            sb.append("Serial number: ").append(new String(strbyte)).append("\r\n");
                            ByteBuf data2 = data.readSlice(14);
                            sb.append("Reserved: ").append(ByteBufUtil.hexDump(data2)).append("\r\n");
                        }
                        break;
                    }

                }
            } catch (Exception ex) {
                sb.append(ex);
            }
            return sb;
        }

        private static StringBuilder ID_0xXY22(ByteBuf data) {
            StringBuilder sb = new StringBuilder();
            int ssl_length = data.readShort();
            int ssl_count = data.readShort();

            try {
                sb.append("SSL partial list length in bytes: ").append(ssl_length).append("\r\n");
                ;
                sb.append("SSL partial list count: ").append(ssl_count).append("\r\n");
                ;

                for (int i = 1; i <= ssl_count; i++) {
                    sb.append("Data Record: ").append(i).append("\r\n");
                    ByteBuf infobytes = data.readSlice(20);
                    short al1 = data.readShort();
                    short al2 = data.readShort();
                    int al3 = data.readInt();
                    sb.append("Info: ").append(ByteBufUtil.hexDump(infobytes)).append("\r\n");
                    sb.append("al 1: ").append(Integer.toHexString(al1)).append("\r\n");
                    sb.append("al 2: ").append(Integer.toHexString(al2)).append("\r\n");
                    sb.append("al 3: ").append(Integer.toHexString(al3)).append("\r\n");
                }
            } catch (Exception ex) {
                sb.append(ex);
            }
            return sb;

        }

        private static StringBuilder ID_0xXY25(ByteBuf data) {
            StringBuilder sb = new StringBuilder();
            sb.append("It is not implemeted yet!");
            return sb;
        }

        private static StringBuilder ID_0xXY32(ByteBuf data) {
            StringBuilder sb = new StringBuilder();
            sb.append("It is not implemeted yet!");
            return sb;
        }

        private static StringBuilder ID_0xXY71(ByteBuf data) {
            StringBuilder sb = new StringBuilder();
            sb.append("It is not implemeted yet!");
            return sb;
        }

        private static StringBuilder ID_0xXY74(ByteBuf data) {
            StringBuilder sb = new StringBuilder();
            int ssl_length = data.readShort();
            int ssl_count = data.readShort();
            try {
                sb.append("SSL partial list length in bytes: ").append(ssl_length).append("\r\n");
                sb.append("SSL partial list count: ").append(ssl_count).append("\r\n");
                for (int i = 1; i <= ssl_count; i++) {
                    sb.append("Data Record: ").append(i).append("\r\n");
                    int cpu_type = data.readByte();
                    int led_id = data.readByte();
                    boolean led_onoff = (data.readByte() != 0x00);
                    int led_flashing = data.readByte();

                    sb.append("Rack number: ").append((cpu_type & 0x07)).append("\r\n");
                    sb.append("CPU Type: ").append((cpu_type & 0x08) == 0x00 ? "Standby" : "Master").append("\r\n");
                    sb.append("LED ID: ").append(LED_ID.valueOf(led_id)).append(" (").append(LED_ID.valueOf(led_id).description).append(")\r\n");
                    sb.append("LED status: ").append(led_onoff).append("\r\n");
                    switch (led_flashing) {
                        case 0x00:
                            sb.append("LED flashing: not.").append("\r\n");
                            break;
                        case 0x01:
                            sb.append("LED blinking: normally (2 Hz).").append("\r\n");
                            break;
                        case 0x02:
                            sb.append("LED blinking: slowly (0.5 Hz).").append("\r\n");
                            break;
                        default:
                            sb.append("LED blinking: no information.").append("\r\n");
                    }
                }
            } catch (Exception ex) {
                sb.append(ex);
            }
            return sb;
        }

        private static StringBuilder ID_0xXY75(ByteBuf data) {
            StringBuilder sb = new StringBuilder();
            sb.append("It is not implemeted yet!");
            return sb;
        }

        private static StringBuilder ID_0xXY90(ByteBuf data) {
            StringBuilder sb = new StringBuilder();
            sb.append("It is not implemeted yet!");
            return sb;
        }

        private static StringBuilder ID_0xXY91(ByteBuf data) {
            StringBuilder sb = new StringBuilder();
            sb.append("It is not implemeted yet!");
            return sb;
        }

        private static StringBuilder ID_0xXY92(ByteBuf data) {
            StringBuilder sb = new StringBuilder();
            sb.append("It is not implemeted yet!");
            return sb;
        }

        private static StringBuilder ID_0xXY94(ByteBuf data) {
            StringBuilder sb = new StringBuilder();
            sb.append("It is not implemeted yet!");
            return sb;
        }

        private static StringBuilder ID_0xXY95(ByteBuf data) {
            StringBuilder sb = new StringBuilder();
            sb.append("It is not implemeted yet!");
            return sb;
        }

        private static StringBuilder ID_0xXY96(ByteBuf data) {
            StringBuilder sb = new StringBuilder();
            sb.append("It is not implemeted yet!");
            return sb;
        }

        private static StringBuilder ID_0xXY9C(ByteBuf data) {
            StringBuilder sb = new StringBuilder();
            sb.append("It is not implemeted yet!");
            return sb;
        }

        private static StringBuilder ID_0xXYA0(ByteBuf data) {
            StringBuilder sb = new StringBuilder();
            int ssl_length = data.readShort();
            int ssl_count = data.readShort();

            try {
                sb.append("SSL partial list length in bytes: ").append(ssl_length).append("\r\n");
                sb.append("SSL partial list count: ").append(ssl_count).append("\r\n");
                int i = 1;
                while (data.isReadable()) {
                    sb.append("Data Record: ").append(i).append("\r\n");
                    Short id = data.readShort();

                    sb.append("ID: ").append(id).append("\r\n");
                    ByteBuf infobytes = data.readSlice(10);
                    sb.append("Info: ").append(ByteBufUtil.hexDump(infobytes)).append("\r\n");
                    infobytes = data.readSlice(8);
                    sb.append("Timestamp: ").append(readDateAndTime(infobytes).toString()).append("\r\n");
                    sb.append("Desc: ").append(S7DiagnosticEventId.valueOf(id).getDescription()).append("\r\n");
                    i++;
                }
            } catch (Exception ex) {
                sb.append(ex);
            }

            return sb;
        }

        private static StringBuilder ID_0xXYB1(ByteBuf data) {
            StringBuilder sb = new StringBuilder();
            sb.append("It is not implemeted yet!");
            return sb;
        }

        private static StringBuilder ID_0xXYB2(ByteBuf data) {
            StringBuilder sb = new StringBuilder();
            sb.append("It is not implemeted yet!");
            return sb;
        }

        private static StringBuilder ID_0xXYB3(ByteBuf data) {
            StringBuilder sb = new StringBuilder();
            sb.append("It is not implemeted yet!");
            return sb;
        }

        private static StringBuilder ID_0xXYB4(ByteBuf data) {
            StringBuilder sb = new StringBuilder();
            sb.append("It is not implemeted yet!");
            return sb;
        }

        /*
         * Date and time of day (BCD coded).
         *          +----------------+
         * Byte n   | Year   0 to 99 |
         *          +----------------+
         * Byte n+1 | Month  1 to 12 |
         *          +----------------+
         * Byte n+2 | Day    1 to 31 |
         *          +----------------+
         * Byte n+3 | Hour   0 to 23 |
         *          +----------------+
         * Byte n+4 | Minute 0 to 59 |
         *          +----------------+
         * Byte n+5 | Second 0 to 59 |
         *          +----------------+
         * Byte n+6 | ms    0 to 999 |
         * Byte n+7 | X X X X X D O W|
         *          +----------------+
         * DOW: Day of weed (last 3 bits)
         */
        private static LocalDateTime readDateAndTime(ByteBuf data) {
            //from Plc4XS7Protocol
            int year = convertByteToBcd(data.readByte());
            byte themonth = data.readByte();
            int month = convertByteToBcd(themonth == 0x00 ? 0x01 : themonth);
            byte theday = data.readByte();
            int day = convertByteToBcd(theday == 0x00 ? 0x01 : theday);
            int hour = convertByteToBcd(data.readByte());
            int minute = convertByteToBcd(data.readByte());
            int second = convertByteToBcd(data.readByte());
            int milliseconds = (data.readShort() & 0xfff0) >> 4;

            int cen = ((milliseconds & 0x0f00) >> 8) * 100;
            int dec = ((milliseconds & 0x00f0) >> 4) * 10;
            milliseconds = cen + dec + (milliseconds & 0x000f);
            int nanoseconds = milliseconds * 1000000;

            //data-type ranges from 1990 up to 2089
            if (year >= 90) {
                year += 1900;
            } else {
                year += 2000;
            }

            return LocalDateTime.of(year, month, day, hour, minute, second, nanoseconds);
        }

        /**
         * converts incoming byte to an integer regarding used BCD format
         *
         * @param incomingByte the incoming byte
         * @return converted BCD number
         */
        private static int convertByteToBcd(byte incomingByte) {
            int dec = (incomingByte >> 4) * 10;
            return dec + (incomingByte & 0x0f);
        }

        /**
         * converts incoming Short to an integer regarding used BCD format
         *
         * @param incomingShort the incoming byte
         * @return converted BCD number
         */
        private static short convertShortToBcd(short incomingShort) {
            return (short) ((incomingShort >> 8) * 100 +
                (incomingShort >> 4) * 10 +
                (incomingShort & 0x0f));
        }

    }

    public static Duration S5TimeToDuration(Short data) {
        Duration res;
        short t = data;
        long tv = (short) (((t & 0x000F)) + ((t & 0x00F0) >> 4) * 10 + ((t & 0x0F00) >> 8) * 100);
        long tb = (short) (10 * Math.pow(10, ((t & 0xF000) >> 12)));
        long totalms = tv * tb;
        if (totalms <= 9990000) {
            res = Duration.ofMillis(totalms);
        } else {
            res = Duration.ofMillis(9990000);
        }
        return res;
    }

    public static Short DurationToS5Time(Duration duration) {
        short tv = 0;
        short tb = 0x0000_0000;
        short s5time = 0x0000;
        long totalms = duration.toMillis();

        if ((totalms >= 0) && (totalms <= 9990000)) {
            if (totalms <= 9990) {
                tb = 0x0000_0000; //10 ms
                tv = (short) (totalms / 10);
            } else if (totalms <= 99900) {
                tb = 0x0000_0001;// 100 ms
                tv = (short) (totalms / 100);
            } else if (totalms <= 999000) {
                tb = 0x0000_0002;//1000 ms
                tv = (short) (totalms / 1000);
            } else if (totalms > 999000) {
                tb = 0x0000_0003;//10000 ms
                tv = (short) (totalms / 10000);
            }

            short uni = (short) (tv % 10);
            short dec = (short) ((tv / 10) % 10);
            short cen = (short) ((tv / 100) % 10);

            return (short) (((tb) << 12) | (cen << 8) | (dec << 4) | (uni));
        }
        return s5time;
    }

    public static Duration S7TimeToDuration(Integer data) {
        Duration res = Duration.ZERO;
        if (data >= 0) {
            res = res.plusMillis((long) data);
        } else {
            long ms = 0x8000_0000 - (data & 0x8000_0000);
            res = res.minusMillis((long) data);
        }

        return res;
    }

    public static Integer DurationToS7Time(Duration data) {
        Integer res = 0x0000_0000;
        if (data.isNegative()) {
            res = (int) data.toMillis() + 0x8000_0000;
        } else {
            res = (int) data.toMillis();
        }
        return res;
    }

    public static LocalTime S7TodToLocalTime(Integer data) {
        if (data > 0x0526_5bff) data = 0x0526_5bff;
        if (data < 0) data = 0x0000_0000;
        return LocalTime.MIDNIGHT.plusNanos((long) data * 1_000_000);
    }

    public static Integer LocalTimeToS7Tod(LocalTime data) {
        return (int) (data.toNanoOfDay() / 1_000_000);
    }

    public static LocalDate S7DateToLocalDate(Short data) {
        LocalDate res = LocalDate.of(1990, 1, 1);
        res = res.plusDays((long) data);
        return res;
    }

    public static Short LocalDateToS7Date(LocalDate data) {
        LocalDate ini = LocalDate.of(1990, 1, 1);
        long resl = ChronoUnit.DAYS.between(ini, data);
        return (short) resl;
    }

    /*
     * Date and time of day (BCD coded).
     *          +----------------+
     * Byte n   | Year   0 to 99 |
     *          +----------------+
     * Byte n+1 | Month  1 to 12 |
     *          +----------------+
     * Byte n+2 | Day    1 to 31 |
     *          +----------------+
     * Byte n+3 | Hour   0 to 23 |
     *          +----------------+
     * Byte n+4 | Minute 0 to 59 |
     *          +----------------+
     * Byte n+5 | Second 0 to 59 |
     *          +----------------+
     * Byte n+6 | ms    0 to 999 |
     * Byte n+7 | X X X X X D O W|
     *          +----------------+
     * DOW: Day of weed (last 3 bits)
     */
    public static LocalDateTime S7DateTimeToLocalDateTime(ByteBuf data) {
        //from Plc4XS7Protocol
        int year = BcdToInt(data.readByte());
        int month = BcdToInt(data.readByte());
        int day = BcdToInt(data.readByte());
        int hour = BcdToInt(data.readByte());
        int minute = BcdToInt(data.readByte());
        int second = BcdToInt(data.readByte());
        int millih = BcdToInt(data.readByte()) * 10;

        int milll = (data.readByte() >> 4);

        int milliseconds = millih + milll;
        int nanoseconds = milliseconds * 1000000;
        //At this point a dont need the day of week
        //data-type ranges from 1990 up to 2089
        if (year >= 90) {
            year += 1900;
        } else {
            year += 2000;
        }

        return LocalDateTime.of(year, month, day, hour, minute, second, nanoseconds);
    }

    public static LocalDateTime S7DateAndTimeToLocalDateTime(int year, int month, int day,
                                                             int hour, int min, int sec, int msec) {
        int nanoseconds = msec * 1000000;
        //At this point a dont need the day of week
        //data-type ranges from 1990 up to 2089
        if (year >= 90) {
            year += 1900;
        } else {
            year += 2000;
        }
        return LocalDateTime.of(year, month, day, hour, min, sec, nanoseconds);
    }

    public static byte[] LocalDateTimeToS7DateTime(LocalDateTime data) {
        byte[] res = new byte[8];

        res[0] = ByteToBcd((data.getYear() % 100));
        res[1] = ByteToBcd(data.getMonthValue());
        res[2] = ByteToBcd(data.getDayOfMonth());
        res[3] = ByteToBcd(data.getHour());
        res[4] = ByteToBcd(data.getMinute());
        res[5] = ByteToBcd(data.getSecond());

        long ms = (long) (data.getNano() / 1_000_000);
        res[6] = (byte) ((int) (((ms / 100) << 4) | ((ms / 10) % 10)));
        //Java:1 (Monday) to 7 (Sunday)->S7:1 (Sunday) to 7 (Saturday)
        byte dayofweek = (byte) ((data.getDayOfWeek().getValue() < 7) ?
            data.getDayOfWeek().getValue() + 1 :
            (byte) 0x01);
        res[7] = (byte) (((ms % 10) << 4) | ((byte) (dayofweek)));

        return res;
    }


    /**
     * converts incoming byte to an integer regarding used BCD format
     *
     * @param incomingByte the incoming byte
     * @return converted BCD number
     */
    private static byte ByteToBcd(int incomingByte) {
        byte dec = (byte) ((incomingByte / 10) % 10);
        return (byte) ((dec << 4) | (incomingByte % 10));
    }

    private static int BcdToInt(byte bcd) {
        return (bcd >> 4) * 10 + (bcd & 0x0f);
    }

    public static void ByteToBcd(final WriteBuffer buffer, short _value) throws SerializationException {
        short incomingByte = _value;
        byte outputByte = 0;
        byte dec = (byte) ((incomingByte / 10) % 10);
        outputByte = (byte) ((dec << 4) | (incomingByte % 10));
        buffer.writeByte(outputByte);
    }

    public static int BcdToInt(final ReadBuffer buffer) throws ParseException {
        byte bcd = buffer.readByte();
        return (bcd >> 4) * 10 + (bcd & 0x0f);
    }

    public static int S7msecToInt(final ReadBuffer buffer) throws ParseException {
        int centenas = BcdToInt(buffer.readUnsignedByte(4));
        int decenas = BcdToInt(buffer.readUnsignedByte(4));
        int unidad = BcdToInt(buffer.readUnsignedByte(4));
        return centenas * 100 + decenas * 10 + unidad;
    }

    public static void IntToS7msec(final WriteBuffer buffer, int _value) throws SerializationException {
        int local = 0;
        if (_value > 999) {
            local = 999;
        } else local = _value;

        int centenas = local / 100;
        int residual = (local - centenas * 100);
        int decenas = (residual) / 10;
        int unidad = residual - (decenas * 10);

        buffer.writeUnsignedByte(4, (byte) centenas);
        buffer.writeUnsignedByte(4, (byte) decenas);
        buffer.writeUnsignedByte(4, (byte) unidad);
    }

    public static void LeftShift3(final WriteBuffer buffer, int _value) throws SerializationException {
        int valor = _value << 3;
        buffer.writeUnsignedInt(16, valor);
    }

    public static int RightShift3(final ReadBuffer buffer) throws ParseException {
        return buffer.readUnsignedInt(16) >> 3;
    }
    
    public static int RightShift3(final ReadBuffer buffer, DataTransportSize tsize) throws ParseException {
        int value = 0;
        if (tsize == DataTransportSize.OCTET_STRING){
            value = buffer.readUnsignedInt(16);
        } else {
            value = buffer.readUnsignedInt(16) >> 3;
        }
        return value;    
    }

    //TODO: apply only if not the last item
    public static int EventItemLength(final ReadBuffer buffer, int valueLength) {
        return ((valueLength % 2 == 0) || (!buffer.hasMore((valueLength + 1) * 8))) ? valueLength : valueLength + 1;
    }


    public static PlcResponseCode decodeResponseCode(DataTransportErrorCode dataTransportErrorCode) {
        if (dataTransportErrorCode == null) {
            return PlcResponseCode.INTERNAL_ERROR;
        }
        switch (dataTransportErrorCode) {
            case OK:
                return PlcResponseCode.OK;
            case NOT_FOUND:
                return PlcResponseCode.NOT_FOUND;
            case INVALID_ADDRESS:
                return PlcResponseCode.INVALID_ADDRESS;
            case DATA_TYPE_NOT_SUPPORTED:
                return PlcResponseCode.INVALID_DATATYPE;
            case ACCESS_DENIED:
                return PlcResponseCode.ACCESS_DENIED;
            default:
                return PlcResponseCode.INTERNAL_ERROR;
        }
    }


    private static byte[] wordToBytes(long data) {
        return new byte[]{
            (byte) ((data >> 8) & 0xff),
            (byte) ((data >> 0) & 0xff),
        };
    }

    private static byte[] dwordToBytes(long data) {
        return new byte[]{
            (byte) ((data >> 24) & 0xff),
            (byte) ((data >> 16) & 0xff),
            (byte) ((data >> 8) & 0xff),
            (byte) ((data >> 0) & 0xff),
        };
    }

    public static String ModeEventProcessing(final S7ModeEvent mode) {
        StringBuilder sb = new StringBuilder("CPU is in : ");
        if (ModeTransitionType.isDefined((short) mode.getMap().get("CURRENT_MODE"))) {
            short currentmode = (short) mode.getMap().get("CURRENT_MODE");
            sb.append(ModeTransitionType.enumForValue(currentmode).name());
        } else {
            sb.append("UNDEFINED");
        }
        return sb.toString();
    }

    public static String SysEventProcessing(final S7SysEvent event, String eventtext, HashMap<String, HashMap<String, String>> textlists) {
        final Pattern EVENT_SIG =
            Pattern.compile("(@[\\d]{0,3}[bycwixdrBYCWIXDR](%([\\d]{0,2}[duxbs]){1}|(\\d\\.\\df){1}|(t#[a-zA-Z0-9]+){1})@)");

        final Pattern FIELDS =
            Pattern.compile("@(?<sig>[\\d]{0,3})(?<type>[bycwixdrBYCWIXDR])(?<format>%([\\d]{0,2}[duxbs]){1}|(\\d\\.\\df){1}|(t#[a-zA-Z0-9]+){1})@");

        final Pattern FIELD_FORMAT =
            Pattern.compile("%([\\d]{0,2})([duxbsDUXBS]{1})");

        Map<String, Object> map = event.getMap();
        Matcher matcher = EVENT_SIG.matcher(eventtext);
        Matcher fields = null;
        Matcher fieldformat = null;

        String strSig = null;
        ByteBuf bytebuf = null;
        int length = 0;
        int sig = 0;
        long value = 0;
        String strOut = eventtext;
        String strField = null;

        while (matcher.find()) {
            fields = FIELDS.matcher(matcher.group(0));
            if (!fields.find()) break;
            sig = fields.group(1) == "" ? 1 : Integer.parseInt(fields.group(1));
            if ((sig == 0) || (sig > 2)) break;
            String infofield = (sig == 1) ? "INFO1" : "INFO2";
            long infovalue = (long) event.getMap().get(infofield);
            String format = fields.group(3).toUpperCase();
            bytebuf = (sig == 1) ?
                Unpooled.wrappedBuffer(wordToBytes(infovalue)) :
                Unpooled.wrappedBuffer(dwordToBytes(infovalue));
            switch (fields.group(2).toUpperCase()) {
                case "B":
                    if (bytebuf.capacity() < Byte.BYTES) break;
                    strField = String.valueOf(bytebuf.getBoolean(0));
                    strOut = strOut.replaceAll(matcher.group(0), strField);
                    ;
                    break;
                case "Y":
                    if (bytebuf.capacity() < Byte.BYTES) break;
                    if (format.contains("U")) {
                        value = bytebuf.getUnsignedByte(0);
                        String strReplace = format.replace("U", "d");
                        strField = String.format(strReplace, value);
                    } else if (format.contains("D")) {
                        value = bytebuf.getByte(0);
                        strField = String.format(format, value);
                    } else if (format.contains("B")) {
                        value = bytebuf.getUnsignedByte(0);
                        strField = Integer.toBinaryString((byte) value);
                    } else {
                        value = bytebuf.getByte(0);
                        strField = String.format(format, value);
                    }
                    strOut = strOut.replaceAll(matcher.group(0), strField);
                    ;
                    break;
                case "C":
                    if (format.contains("%T#")) {

                    } else {
                        if (bytebuf.capacity() < Byte.BYTES) break;
                        fieldformat = FIELD_FORMAT.matcher(format);
                        if (fieldformat.find()) {
                            length = Integer.parseInt(fieldformat.group(1));
                            length = (length > bytebuf.capacity()) ? bytebuf.capacity() : length;
                            strField =
                                bytebuf.readCharSequence(length, Charset.forName("utf-8")).toString();
                        }
                    }
                    strOut = strOut.replaceAll(matcher.group(0), strField);
                    ;
                    break;
                case "W":
                    if (bytebuf.capacity() < Short.BYTES) break;
                    if (format.contains("U")) {
                        value = bytebuf.getUnsignedShort(0);
                        String strReplace = format.replace("U", "d");
                        strField = String.format(strReplace, value);
                    } else if (format.contains("D")) {
                        value = bytebuf.getShort(0);
                        strField = String.format(format, value);
                    } else if (format.contains("B")) {
                        value = bytebuf.getUnsignedShort(0);
                        strField = Integer.toBinaryString((short) value);
                    } else {
                        value = bytebuf.getShort(0);
                        strField = String.format(format, value);
                    }
                    strOut = strOut.replaceAll(matcher.group(0), strField);
                    ;
                    break;
                case "I":
                    if (bytebuf.capacity() < Integer.BYTES) break;
                    if (format.contains("U")) {
                        value = bytebuf.getUnsignedInt(0);
                        String strReplace = format.replace("U", "d");
                        strField = String.format(strReplace, value);
                    } else if (format.contains("D")) {
                        value = bytebuf.getInt(0);
                        strField = String.format(format, value);
                    } else if (format.contains("B")) {
                        value = bytebuf.getUnsignedInt(0);
                        strField = Long.toBinaryString(value);
                    } else {
                        value = bytebuf.getInt(0);
                        strField = String.format(format, value);
                    }
                    strOut = strOut.replaceAll(matcher.group(0), strField);
                    break;
                case "X":
                    if (bytebuf.capacity() < Long.BYTES) break;
                    if (format.contains("U")) {
                        value = bytebuf.getUnsignedInt(0);
                        String strReplace = format.replace("U", "d");
                        strField = String.format(strReplace, value);
                    } else if (format.contains("D")) {
                        value = bytebuf.getLong(0);
                        strField = String.format(format, value);
                    } else if (format.contains("B")) {
                        value = bytebuf.getUnsignedInt(0);
                        strField = Long.toBinaryString(value);
                    } else {
                        value = bytebuf.getLong(0);
                        strField = String.format(format, value);
                    }
                    strOut = strOut.replaceAll(matcher.group(0), strField);
                    break;
                case "D":
                    if (bytebuf.capacity() < Double.BYTES) break;
                    if (format.contains("U")) {
                        value = bytebuf.getUnsignedInt(0);
                        String strReplace = format.replace("U", "d");
                        strField = String.format(strReplace, value);
                    } else if (format.contains("D")) {
                        value = bytebuf.getLong(0);
                        strField = String.format(format, value);
                    } else if (format.contains("B")) {
                        value = bytebuf.getUnsignedInt(0);
                        strField = Long.toBinaryString(value);
                    } else {
                        value = bytebuf.getLong(0);
                        strField = String.format(format, value);
                    }
                    strOut = strOut.replaceAll(matcher.group(0), strField);
                    break;
                case "R":
                    if (bytebuf.capacity() < Float.BYTES) break;
                    if (format.contains("F")) {
                        strField = String.format(format, value);
                        strOut = strOut.replaceAll(matcher.group(0), strField);
                    }
                    ;
                    break;
            }
        }

        return strOut;
    }

    /**
     * Symbol       Meaning
     * ======       =======
     *
     * @param alarm     Alarm type from PLC.
     * @param alarmText The text string to be processed.
     * @param textlists List of texts for indexed replacement.
     * @return The text string with the replacement values.
     * @ Beginning of format string
     * Pv           Process value number (optional)
     * Typ          Type of process value (optional)
     * % width      Width of the signaling field on the OP
     * .precision   Number of decimal places (optional with format = "f")
     * format       Process value representation mode
     * @ End of format string
     * <p>
     * Type		Meaning
     * ====             =======
     * B		BOOL
     * Y		BYTE
     * C		CHARACTER
     * W		WORD
     * I		INTEGER
     * X		DWORD
     * D		DINT
     * R		REAL
     * <p>
     * Representation       Format
     * ==============       ======
     * %[i]d                Decimal signed
     * %[i]u                Decimal without sign
     * %[i]x                Hexadecimal
     * %[i].[y]f            Signed fixed point
     * %[i]b                Binary
     * %[i]s                Character string (STRING ANSI)
     * %t#<Library name>    Access to text library
     * <p>
     * General representation:
     * @@&gt;Associated Value>&gt;Type>&gt;Format>&gt;Library name>@
     */
    public static String AlarmProcessing(final S7AlarmEvent alarm, String alarmText, HashMap<String, HashMap<String, String>> textlists) {
        final Pattern ALARM_SIG =
            Pattern.compile("(@[\\d]{0,3}[bycwixdrBYCWIXDR](%([\\d]{0,2}[duxbs]){1}|(\\d\\.\\df){1}|(t#[a-zA-Z0-9]+){1})@)");

        final Pattern FIELDS =
            Pattern.compile("@(?<sig>[\\d]{0,3})(?<type>[bycwixdrBYCWIXDR])(?<format>%([\\d]{0,2}[duxbs]){1}|(\\d\\.\\df){1}|(t#[a-zA-Z0-9]+){1})@");

        final Pattern FIELD_FORMAT =
            Pattern.compile("%([\\d]{0,2})([duxbsDUXBS]{1})");

        Map<String, Object> map = alarm.getMap();
        Matcher matcher = ALARM_SIG.matcher(alarmText);
        Matcher fields = null;
        Matcher fieldformat = null;

        String strSig = null;
        ByteBuf bytebuf = null;
        int length = 0;
        int sig = 0;
        long value = 0;
        String strOut = new String(alarmText);
        String strField = null;

        while (matcher.find()) {
            fields = FIELDS.matcher(matcher.group(0));
            if (!fields.find()) break;
            sig = fields.group(1) == "" ? 1 : Integer.parseInt(fields.group(1));
            strSig = "SIG_" + sig + "_DATA";
            if ((((short) map.get("ASSOCIATED_VALUES")) == 0) ||
                (sig > ((short) map.get("ASSOCIATED_VALUES")))) break;
            bytebuf = Unpooled.wrappedBuffer((byte[]) map.get(strSig));
            String format = fields.group(3).toUpperCase();
            switch (fields.group(2).toUpperCase()) {
                case "B":
                    if (bytebuf.capacity() < Byte.BYTES) break;
                    strField = String.valueOf(bytebuf.getBoolean(0));
                    strOut = strOut.replaceAll(matcher.group(0), strField);
                    ;
                    break;
                case "Y":
                    if (bytebuf.capacity() < Byte.BYTES) break;
                    if (format.contains("U")) {
                        value = bytebuf.getUnsignedByte(0);
                        String strReplace = format.replace("U", "d");
                        strField = String.format(strReplace, value);
                    } else if (format.contains("D")) {
                        value = bytebuf.getByte(0);
                        strField = String.format(format, value);
                    } else if (format.contains("B")) {
                        value = bytebuf.getUnsignedByte(0);
                        strField = Integer.toBinaryString((byte) value);
                    } else {
                        value = bytebuf.getByte(0);
                        strField = String.format(format, value);
                    }
                    strOut = strOut.replaceAll(matcher.group(0), strField);
                    ;
                    break;
                case "C":
                    if (format.contains("%T#")) {

                    } else {
                        if (bytebuf.capacity() < Byte.BYTES) break;
                        fieldformat = FIELD_FORMAT.matcher(format);
                        if (fieldformat.find()) {
                            length = Integer.parseInt(fieldformat.group(1));
                            length = (length > bytebuf.capacity()) ? bytebuf.capacity() : length;
                            strField =
                                bytebuf.readCharSequence(length, Charset.forName("utf-8")).toString();
                        }
                    }
                    strOut = strOut.replaceAll(matcher.group(0), strField);
                    ;
                    break;
                case "W":
                    if (bytebuf.capacity() < Short.BYTES) break;
                    if (format.contains("U")) {
                        value = bytebuf.getUnsignedShort(0);
                        String strReplace = format.replace("U", "d");
                        strField = String.format(strReplace, value);
                    } else if (format.contains("D")) {
                        value = bytebuf.getShort(0);
                        strField = String.format(format, value);
                    } else if (format.contains("B")) {
                        value = bytebuf.getUnsignedShort(0);
                        strField = Integer.toBinaryString((short) value);
                    } else {
                        value = bytebuf.getShort(0);
                        strField = String.format(format, value);
                    }
                    strOut = strOut.replaceAll(matcher.group(0), strField);
                    ;
                    break;
                case "I":
                    if (bytebuf.capacity() < Integer.BYTES) break;
                    if (format.contains("U")) {
                        value = bytebuf.getUnsignedInt(0);
                        String strReplace = format.replace("U", "d");
                        strField = String.format(strReplace, value);
                    } else if (format.contains("D")) {
                        value = bytebuf.getInt(0);
                        strField = String.format(format, value);
                    } else if (format.contains("B")) {
                        value = bytebuf.getUnsignedInt(0);
                        strField = Long.toBinaryString(value);
                    } else {
                        value = bytebuf.getInt(0);
                        strField = String.format(format, value);
                    }
                    ;
                    strOut = strOut.replaceAll(matcher.group(0), strField);
                    ;
                    break;
                case "X":
                    if (bytebuf.capacity() < Long.BYTES) break;
                    if (format.contains("U")) {
                        value = bytebuf.getUnsignedInt(0);
                        String strReplace = format.replace("U", "d");
                        strField = String.format(strReplace, value);
                    } else if (format.contains("D")) {
                        value = bytebuf.getLong(0);
                        strField = String.format(format, value);
                    } else if (format.contains("B")) {
                        value = bytebuf.getUnsignedInt(0);
                        strField = Long.toBinaryString(value);
                    } else {
                        value = bytebuf.getLong(0);
                        strField = String.format(format, value);
                    }
                    ;
                    strOut = strOut.replaceAll(matcher.group(0), strField);
                    ;
                    break;
                case "D":
                    if (bytebuf.capacity() < Double.BYTES) break;
                    if (format.contains("U")) {
                        value = bytebuf.getUnsignedInt(0);
                        String strReplace = format.replace("U", "d");
                        strField = String.format(strReplace, value);
                    } else if (format.contains("D")) {
                        value = bytebuf.getLong(0);
                        strField = String.format(format, value);
                    } else if (format.contains("B")) {
                        value = bytebuf.getUnsignedInt(0);
                        strField = Long.toBinaryString(value);
                    } else {
                        value = bytebuf.getLong(0);
                        strField = String.format(format, value);
                    }
                    ;
                    strOut = strOut.replaceAll(matcher.group(0), strField);
                    ;
                    break;
                case "R":
                    if (bytebuf.capacity() < Float.BYTES) break;
                    if (format.contains("F")) {
                        strField = String.format(format, value);
                        strOut = strOut.replaceAll(matcher.group(0), strField);
                    }
                    ;
                    break;

            }
        }

        return strOut;
    }

    public static LocalTime parseTiaTime(ReadBuffer io) {
        try {
            int millisSinceMidnight = io.readInt(32);
            return LocalTime.now().withHour(0).withMinute(0).withSecond(0).withNano(0).plus(
                millisSinceMidnight, ChronoUnit.MILLIS);
        } catch (ParseException e) {
            return null;
        }
    }

    public static void serializeTiaTime(WriteBuffer io, PlcValue value) {
        throw new NotImplementedException("Serializing TIME not implemented");
    }

    public static LocalTime parseS5Time(ReadBuffer io) {
        try {
            int stuff = io.readInt(16);
            // TODO: Implement this correctly.
            throw new NotImplementedException("S5TIME not implemented");
        } catch (ParseException e) {
            return null;
        }
    }

    public static void serializeS5Time(WriteBuffer io, PlcValue value) {
        throw new NotImplementedException("Serializing S5TIME not implemented");

    }

    public static LocalTime parseTiaLTime(ReadBuffer io) {
        throw new NotImplementedException("LTIME not implemented");
    }

    public static void serializeTiaLTime(WriteBuffer io, PlcValue value) {
        throw new NotImplementedException("Serializing LTIME not implemented");
    }

    public static LocalTime parseTiaTimeOfDay(ReadBuffer io) {
        try {
            long millisSinceMidnight = io.readUnsignedLong(32);
            return LocalTime.now().withHour(0).withMinute(0).withSecond(0).withNano(0).plus(
                millisSinceMidnight, ChronoUnit.MILLIS);
        } catch (ParseException e) {
            return null;
        }
    }

    public static void serializeTiaTimeOfDay(WriteBuffer io, PlcValue value) {
        throw new NotImplementedException("Serializing TIME_OF_DAY not implemented");
    }

    public static LocalDate parseTiaDate(ReadBuffer io) {
        try {
            int daysSince1990 = io.readUnsignedInt(16);
            return LocalDate.now().withYear(1990).withDayOfMonth(1).withMonth(1).plus(daysSince1990, ChronoUnit.DAYS);
        } catch (ParseException e) {
            return null;
        }
    }

    public static void serializeTiaDate(WriteBuffer io, PlcValue value) {
        throw new NotImplementedException("Serializing DATE not implemented");
    }

    //TODO: Call BCD converter
    public static LocalDateTime parseTiaDateTime(ReadBuffer io) {
        try {
            int year = io.readUnsignedInt(16);
            int month = io.readUnsignedInt(8);
            int day = io.readUnsignedInt(8);
            // Skip day-of-week
            io.readByte();
            int hour = io.readByte();
            int minute = io.readByte();
            int second = io.readByte();
            int nanosecond = io.readUnsignedInt(24);
            // Skip day-of-week
            io.readByte();

            return LocalDateTime.of(year, month, day, hour, minute, second, nanosecond);
        } catch (Exception e) {
            return null;
        }
    }

    public static void serializeTiaDateTime(WriteBuffer io, PlcValue value) {
        throw new NotImplementedException("Serializing DATE_AND_TIME not implemented");
    }

    public static String parseS7Char(ReadBuffer io, String encoding) throws ParseException {
        if ("UTF-8".equalsIgnoreCase(encoding)) {
            return io.readString(8, WithOption.WithEncoding(encoding));
        } else if ("UTF-16".equalsIgnoreCase(encoding)) {
            return io.readString(16, WithOption.WithEncoding(encoding));
        } else {
            throw new PlcRuntimeException("Unsupported encoding");
        }
    }

    public static String parseS7String(ReadBuffer io, int stringLength, String encoding) {
        try {
            if ("UTF-8".equalsIgnoreCase(encoding)) {
                // This is the maximum number of bytes a string can be long.
                short maxLength = io.readUnsignedShort(8);
                // This is the total length of the string on the PLC (Not necessarily the number of characters read)
                short totalStringLength = io.readUnsignedShort(8);

                final byte[] byteArray = new byte[totalStringLength];
                for (int i = 0; (i < stringLength) && io.hasMore(8); i++) {
                    final byte curByte = io.readByte();
                    if (i < totalStringLength) {
                        byteArray[i] = curByte;
                    } else {
                        // Gobble up the remaining data, which is not added to the string.
                        i++;
                        for (; (i < stringLength) && io.hasMore(8); i++) {
                            io.readByte();
                        }
                        break;
                    }
                }
                return new String(byteArray, StandardCharsets.UTF_8);
            } else if ("UTF-16".equalsIgnoreCase(encoding)) {
                // This is the maximum number of bytes a string can be long.
                int maxLength = io.readUnsignedInt(16);
                // This is the total length of the string on the PLC (Not necessarily the number of characters read)
                int totalStringLength = io.readUnsignedInt(16);

                final byte[] byteArray = new byte[totalStringLength * 2];
                for (int i = 0; (i < stringLength) && io.hasMore(16); i++) {
                    final short curShort = io.readShort(16);
                    if (i < totalStringLength) {
                        byteArray[i * 2] = (byte) (curShort >>> 8);
                        byteArray[(i * 2) + 1] = (byte) (curShort & 0xFF);
                    } else {
                        // Gobble up the remaining data, which is not added to the string.
                        i++;
                        for (; (i < stringLength) && io.hasMore(16); i++) {
                            io.readShort(16);
                        }
                        break;
                    }
                }
                return new String(byteArray, StandardCharsets.UTF_16);
            } else {
                throw new PlcRuntimeException("Unsupported string encoding " + encoding);
            }
        } catch (ParseException e) {
            throw new PlcRuntimeException("Error parsing string", e);
        }
    }

    /*
     * A variable of data type CHAR (character) occupies one byte.
     */
    public static void serializeS7Char(WriteBuffer io, PlcValue value, String encoding) {
        // TODO: Need to implement the serialization or we can't write strings
        if ("UTF-8".equalsIgnoreCase(encoding)) {
            //return io.readString(8, encoding);
        } else if ("UTF-16".equalsIgnoreCase(encoding)) {
            //return io.readString(16, encoding);
        } else {
            throw new PlcRuntimeException("Unsupported encoding");
        }
    }

    /*           +-------------------+
     * Byte n     | Maximum length    | (k)
     *            +-------------------+
     * Byte n+1   | Current Length    | (m)
     *            +-------------------+
     * Byte n+2   | 1st character     | \         \
     *            +-------------------+  |         |
     * Byte n+3   | 2st character     |  | Current |
     *            +-------------------+   >        |
     * Byte ...   | ...               |  | length  |  Maximum
     *            +-------------------+  |          >
     * Byte n+m+1 | mth character     | /          |  length
     *            +-------------------+            |
     * Byte ...   | ...               |            |
     *            +-------------------+            |
     * Byte ...   | ...               |           /
     *            +-------------------+
     * For this version, the user must read the maximum acceptable length in
     * the string in a first instance.
     * Then the user application should avoid the envelope of the adjacent
     * fields passing the maximum length in "stringLength".
     * If your application does not handle S7string, you can handle
     * the String as char arrays from your application.
     */
    public static void serializeS7String(WriteBuffer io, PlcValue value, int stringLength, String encoding) {
        int k = 0xFF & ((stringLength > 250) ? 250 : stringLength);
        int m = 0xFF & value.getString().length();
        m = (m > k) ? k : m;
        byte[] chars = new byte[m];
        for (int i = 0; i < m; ++i) {
            char c = value.getString().charAt(i);
            chars[i] = (byte) c;
        }

        try {
            io.writeByte((byte)(k & 0xFF));
            io.writeByte((byte)(m & 0xFF));
            io.writeByteArray(chars);
        } catch (SerializationException ex) {
            Logger.getLogger(StaticHelper.class.getName()).log(Level.SEVERE, null, ex);
        }
    }

}
