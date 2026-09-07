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
package org.apache.plc4x.java.ads.manual;

import org.apache.plc4x.java.api.authentication.PlcNullAuthentication;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.values.*;
import org.apache.plc4x.java.utils.testutils.manual.BasicPlcTest;

import java.math.BigInteger;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class ManualFactoryAdsDriverTestTC3 extends BasicPlcTest {

    public ManualFactoryAdsDriverTestTC3(String connectionString) {
        super(connectionString, new PlcNullAuthentication(), true, true, true, true, 100);
    }

    public static void main(String[] args) throws Exception {
        String connectionString = "ads:tcp://192.168.23.20:48898?target-ams-port=851&source-ams-port=65534&source-ams-net-id=192.168.23.220.1.1&target-ams-net-id=192.168.23.20.1.1";
        ManualFactoryAdsDriverTestTC3 test = new ManualFactoryAdsDriverTestTC3(connectionString);

        // =========================================================
        // 1) Scalars — all primitive PlcValueType you declared
        // =========================================================
        test.addTestCase("MAIN.g_b1",  new PlcBOOL(true));
        test.addTestCase("MAIN.g_b8",  new PlcBYTE(0xAB));
        test.addTestCase("MAIN.g_s8",  new PlcSINT(-12));
        test.addTestCase("MAIN.g_u8",  new PlcUSINT(250));
        test.addTestCase("MAIN.g_b16", new PlcWORD(0xBEEF));
        test.addTestCase("MAIN.g_s16", new PlcINT(-1234));
        test.addTestCase("MAIN.g_u16", new PlcUINT(54321));
        test.addTestCase("MAIN.g_b32", new PlcDWORD(0xDEADBEEFL));
        test.addTestCase("MAIN.g_s32", new PlcDINT(-12345678));
        test.addTestCase("MAIN.g_u32", new PlcUDINT(305_419_896L)); // 0x12345678
        test.addTestCase("MAIN.g_b64", new PlcLWORD(0x0123_4567_89AB_CDEFL));
        test.addTestCase("MAIN.g_s64", new PlcLINT(-9_223_372_036_854_770_000L));
        test.addTestCase("MAIN.g_u64", new PlcULINT(new BigInteger("18446744073709551000", 10)));

        test.addTestCase("MAIN.g_r32", new PlcREAL(3.14159f));
        test.addTestCase("MAIN.g_r64", new PlcLREAL(2.718281828459045d));

        test.addTestCase("MAIN.g_tim",     new PlcTIME(java.time.Duration.parse("PT2.5S")));
        test.addTestCase("MAIN.g_ltim",    new PlcLTIME(java.time.Duration.parse("PT3.20005S"))); // 3s 200ms 50µs
        test.addTestCase("MAIN.g_dat",     new PlcDATE(java.time.LocalDate.parse("2025-11-12")));
        test.addTestCase("MAIN.g_timoday", new PlcTIME_OF_DAY(java.time.LocalTime.parse("14:33:21.250")));
        test.addTestCase("MAIN.g_dattim",  new PlcDATE_AND_TIME(java.time.LocalDateTime.parse("2025-11-12T14:33:21")));

        test.addTestCase("MAIN.g_str",  new PlcSTRING("Hello PLC4X"));
        test.addTestCase("MAIN.g_wstr", new PlcWSTRING("Grüße aus TwinCAT 👋"));

        // =========================================================
        // 2) 1-D arrays — element reads and full-array reads (PlcList)
        // =========================================================
        test.addTestCase("MAIN.g_arrBool[1]", new PlcBOOL(true));
        test.addTestCase("MAIN.g_arrBool[8]", new PlcBOOL(false));
        test.addTestCase("MAIN.g_arrBool",
            new PlcList(java.util.Arrays.asList(
                new PlcBOOL(true), new PlcBOOL(false), new PlcBOOL(true), new PlcBOOL(true),
                new PlcBOOL(false), new PlcBOOL(false), new PlcBOOL(true), new PlcBOOL(false)
            ))
        );

        test.addTestCase("MAIN.g_arrByte[0]", new PlcBYTE(0xDE));
        test.addTestCase("MAIN.g_arrByte[7]", new PlcBYTE(0x78));
        test.addTestCase("MAIN.g_arrByte",
            new PlcList(java.util.Arrays.asList(
                new PlcBYTE(0xDE), new PlcBYTE(0xAD), new PlcBYTE(0xBE), new PlcBYTE(0xEF),
                new PlcBYTE(0x12), new PlcBYTE(0x34), new PlcBYTE(0x56), new PlcBYTE(0x78)
            ))
        );

        test.addTestCase("MAIN.g_arrInt[1]",  new PlcINT(-3));
        test.addTestCase("MAIN.g_arrInt[5]",  new PlcINT(3));
        test.addTestCase("MAIN.g_arrInt",
            new PlcList(java.util.Arrays.asList(
                new PlcINT(-3), new PlcINT(-1), new PlcINT(0), new PlcINT(1), new PlcINT(3)
            ))
        );

        test.addTestCase("MAIN.g_arrUInt[1]", new PlcUINT(1));
        test.addTestCase("MAIN.g_arrUInt[5]", new PlcUINT(10_000));

        test.addTestCase("MAIN.g_arrDInt[1]", new PlcDINT(-1000));
        test.addTestCase("MAIN.g_arrDInt[4]", new PlcDINT(2_000_000));

        test.addTestCase("MAIN.g_arrUDInt[1]", new PlcUDINT(0L));
        test.addTestCase("MAIN.g_arrUDInt[4]", new PlcUDINT(0x1234_5678L));

        test.addTestCase("MAIN.g_arrLReal[1]", new PlcLREAL(1.5d));
        test.addTestCase("MAIN.g_arrLReal[2]", new PlcLREAL(-2.0d));
        test.addTestCase("MAIN.g_arrLReal[3]", new PlcLREAL(0.125d));

        test.addTestCase("MAIN.g_arrTime[1]", new PlcTIME(java.time.Duration.parse("PT0.01S")));
        test.addTestCase("MAIN.g_arrTime[3]", new PlcTIME(java.time.Duration.parse("PT10S")));

        test.addTestCase("MAIN.g_arrString[1]", new PlcSTRING("alpha"));
        test.addTestCase("MAIN.g_arrString[3]", new PlcSTRING("gamma"));
        test.addTestCase("MAIN.g_arrString",
            new PlcList(java.util.Arrays.asList(
                new PlcSTRING("alpha"), new PlcSTRING("beta"), new PlcSTRING("gamma")
            ))
        );

        test.addTestCase("MAIN.g_arrWString[1]", new PlcWSTRING("Äpfel"));
        test.addTestCase("MAIN.g_arrWString[2]", new PlcWSTRING("Öl"));

        // =========================================================
        // 3) Multidimensional arrays — element + partial dimension -> PlcList
        // =========================================================
        // ---- INT[2,3] ----
        test.addTestCase("MAIN.g_matI16_2x3[1][1]", new PlcINT(10));
        test.addTestCase("MAIN.g_matI16_2x3[1][3]", new PlcINT(12));
        test.addTestCase("MAIN.g_matI16_2x3[2][1]", new PlcINT(-10));
        test.addTestCase("MAIN.g_matI16_2x3[2][3]", new PlcINT(-12));

        // row 1 as PlcList (3 elements)
        test.addTestCase("MAIN.g_matI16_2x3[1]",
            new PlcList(java.util.Arrays.asList(
                new PlcINT(10), new PlcINT(11), new PlcINT(12)
            ))
        );
        // row 2 as PlcList (3 elements)
        test.addTestCase("MAIN.g_matI16_2x3[2]",
            new PlcList(java.util.Arrays.asList(
                new PlcINT(-10), new PlcINT(-11), new PlcINT(-12)
            ))
        );
        // full 2D matrix as List<List<PlcINT>>
        test.addTestCase("MAIN.g_matI16_2x3",
            new PlcList(java.util.Arrays.asList(
                new PlcList(java.util.Arrays.asList(new PlcINT(10),  new PlcINT(11),  new PlcINT(12))),
                new PlcList(java.util.Arrays.asList(new PlcINT(-10), new PlcINT(-11), new PlcINT(-12)))
            ))
        );

        // ---- REAL[3,2] ----
        test.addTestCase("MAIN.g_matR32_3x2[1][1]", new PlcREAL(1.0f));
        test.addTestCase("MAIN.g_matR32_3x2[3][2]", new PlcREAL(3.5f));
        // row 2 as PlcList (2 elements)
        test.addTestCase("MAIN.g_matR32_3x2[2]",
            new PlcList(java.util.Arrays.asList(new PlcREAL(2.0f), new PlcREAL(2.5f)))
        );
        // full matrix (3x2)
        test.addTestCase("MAIN.g_matR32_3x2",
            new PlcList(java.util.Arrays.asList(
                new PlcList(java.util.Arrays.asList(new PlcREAL(1.0f), new PlcREAL(1.5f))),
                new PlcList(java.util.Arrays.asList(new PlcREAL(2.0f), new PlcREAL(2.5f))),
                new PlcList(java.util.Arrays.asList(new PlcREAL(3.0f), new PlcREAL(3.5f)))
            ))
        );

        // ---- UINT[2,2,2] (3D) ----
        test.addTestCase("MAIN.g_cubeU16_2x2x2[1][1][1]", new PlcUINT(1));
        test.addTestCase("MAIN.g_cubeU16_2x2x2[1][2][2]", new PlcUINT(4));
        test.addTestCase("MAIN.g_cubeU16_2x2x2[2][1][1]", new PlcUINT(5));
        test.addTestCase("MAIN.g_cubeU16_2x2x2[2][2][2]", new PlcUINT(8));

        // Fix first dimension => 2D plane (as List<List<PlcUINT>>)
        test.addTestCase("MAIN.g_cubeU16_2x2x2[1]",
            new PlcList(java.util.Arrays.asList(
                new PlcList(java.util.Arrays.asList(new PlcUINT(1), new PlcUINT(2))), // [1,1,*]
                new PlcList(java.util.Arrays.asList(new PlcUINT(3), new PlcUINT(4)))  // [1,2,*]
            ))
        );
        // Plane for index 2
        test.addTestCase("MAIN.g_cubeU16_2x2x2[2]",
            new PlcList(java.util.Arrays.asList(
                new PlcList(java.util.Arrays.asList(new PlcUINT(5), new PlcUINT(6))),
                new PlcList(java.util.Arrays.asList(new PlcUINT(7), new PlcUINT(8)))
            ))
        );
        // Full 3D cube => List<List<List<PlcUINT>>>
        test.addTestCase("MAIN.g_cubeU16_2x2x2",
            new PlcList(java.util.Arrays.asList(
                new PlcList(java.util.Arrays.asList(
                    new PlcList(java.util.Arrays.asList(new PlcUINT(1), new PlcUINT(2))),
                    new PlcList(java.util.Arrays.asList(new PlcUINT(3), new PlcUINT(4)))
                )),
                new PlcList(java.util.Arrays.asList(
                    new PlcList(java.util.Arrays.asList(new PlcUINT(5), new PlcUINT(6))),
                    new PlcList(java.util.Arrays.asList(new PlcUINT(7), new PlcUINT(8)))
                ))
            ))
        );

        // =========================================================
        // 4) Structs (TSimpleStruct, TMatrixI16, TSignalPack)
        // =========================================================
        test.addTestCase("MAIN.g_simple.s8",  new PlcSINT(-8));
        test.addTestCase("MAIN.g_simple.u64", new PlcULINT(64000L));
        test.addTestCase("MAIN.g_simple.b16", new PlcWORD(0xCAFE));
        test.addTestCase("MAIN.g_simple.r64", new PlcLREAL(-0.125d));
        test.addTestCase("MAIN.g_simple.tim", new PlcTIME(java.time.Duration.parse("PT0.123S")));
        test.addTestCase("MAIN.g_simple.str", new PlcSTRING("struct-string"));
        test.addTestCase("MAIN.g_simple.wstr",new PlcWSTRING("Struktur-WSTRING"));

        // Whole struct as PlcStruct (keys match your field names)
        {
            Map<String, PlcValue> s = new HashMap<>();
            s.put("s8", new PlcSINT(-8)); s.put("u8", new PlcUSINT(200));
            s.put("s16", new PlcINT(-1600)); s.put("u16", new PlcUINT(1600));
            s.put("s32", new PlcDINT(-32000)); s.put("u32", new PlcUDINT(32000L));
            s.put("s64", new PlcLINT(-64000L)); s.put("u64", new PlcULINT(64000L));
            s.put("b1", new PlcBOOL(true)); s.put("b8", new PlcBYTE(0x5A));
            s.put("b16", new PlcWORD(0xCAFE)); s.put("b32", new PlcDWORD(0xC0FFEE00L));
            s.put("b64", new PlcLWORD(new BigInteger("DEADBEEFF00DCAFE", 16)));
            s.put("r32", new PlcREAL(0.5f)); s.put("r64", new PlcLREAL(-0.125d));
            s.put("tim", new PlcTIME(java.time.Duration.parse("PT0.123S")));
            s.put("ltim", new PlcLTIME(java.time.Duration.parse("PT1.002003S")));
            s.put("dat", new PlcDATE(java.time.LocalDate.parse("2025-11-12")));
            s.put("timoday", new PlcTIME_OF_DAY(java.time.LocalTime.parse("06:07:08.009")));
            s.put("dattim", new PlcDATE_AND_TIME(java.time.LocalDateTime.parse("2025-11-12T06:07:08")));
            s.put("str", new PlcSTRING("struct-string"));
            s.put("wstr", new PlcWSTRING("Struktur-WSTRING"));
            test.addTestCase("MAIN.g_simple", new PlcStruct(s));
        }

        // Struct with 2D array inside
        test.addTestCase("MAIN.g_matrixI16.m[1][2]", new PlcINT(101));

        test.addTestCase("MAIN.g_matrixI16.m",
            new PlcList(java.util.Arrays.asList(
                new PlcList(java.util.Arrays.asList(new PlcINT(100), new PlcINT(101), new PlcINT(102))),
                new PlcList(java.util.Arrays.asList(new PlcINT(200), new PlcINT(201), new PlcINT(202)))
            ))
        );

        // TSignalPack (0-based arrays!)
        test.addTestCase("MAIN.g_signals.flags[0]", new PlcBOOL(true));
        test.addTestCase("MAIN.g_signals.flags[7]", new PlcBOOL(false));
        test.addTestCase("MAIN.g_signals.words[0]", new PlcWORD(0x1111));
        test.addTestCase("MAIN.g_signals.words[3]", new PlcWORD(0x4444));
        // full payload as PlcList of PlcBYTE
        test.addTestCase("MAIN.g_signals.payload",
            new PlcList(java.util.Arrays.asList(
                new PlcBYTE(0x00), new PlcBYTE(0x11), new PlcBYTE(0x22), new PlcBYTE(0x33),
                new PlcBYTE(0x44), new PlcBYTE(0x55), new PlcBYTE(0x66), new PlcBYTE(0x77),
                new PlcBYTE(0x88), new PlcBYTE(0x99), new PlcBYTE(0xAA), new PlcBYTE(0xBB),
                new PlcBYTE(0xCC), new PlcBYTE(0xDD), new PlcBYTE(0xEE), new PlcBYTE(0xFF)
            ))
        );

        // =========================================================
        // 5) Nested struct: TPlantSnapshot (mix of everything)
        // =========================================================
        test.addTestCase("MAIN.g_plant.meta.b16",  new PlcWORD(0xABCD));
        test.addTestCase("MAIN.g_plant.meta.r64",  new PlcLREAL(-98.765d));
        test.addTestCase("MAIN.g_plant.meta.timoday", new PlcTIME_OF_DAY(java.time.LocalTime.parse("12:34:56.789")));
        test.addTestCase("MAIN.g_plant.meta.str",  new PlcSTRING("meta-ok"));

        test.addTestCase("MAIN.g_plant.gridI16.m[1][1]", new PlcINT(1));
        test.addTestCase("MAIN.g_plant.gridI16.m[2][3]", new PlcINT(6));

        // signals inside g_plant
        test.addTestCase("MAIN.g_plant.signals.words[0]",   new PlcWORD(0xDEAD));
        test.addTestCase("MAIN.g_plant.signals.payload[14]", new PlcBYTE(0xF0));

        // channels array (2 elements)
        test.addTestCase("MAIN.g_plant.channels[1].id",       new PlcUDINT(1L));
        test.addTestCase("MAIN.g_plant.channels[1].name",     new PlcSTRING("CH-A"));
        test.addTestCase("MAIN.g_plant.channels[1].enabled",  new PlcBOOL(true));
        test.addTestCase("MAIN.g_plant.channels[1].setpoints[4]", new PlcREAL(30.0f));
        // lut is assigned in your init code
        test.addTestCase("MAIN.g_plant.channels[1].lut[2][2][2]", new PlcINT(-4));

        test.addTestCase("MAIN.g_plant.channels[2].id",       new PlcUDINT(2L));
        test.addTestCase("MAIN.g_plant.channels[2].name",     new PlcSTRING("CH-B"));
        test.addTestCase("MAIN.g_plant.channels[2].enabled",  new PlcBOOL(false));
        test.addTestCase("MAIN.g_plant.channels[2].setpoints[1]", new PlcREAL(5.5f));
        test.addTestCase("MAIN.g_plant.channels[2].lut[2][2][2]",   new PlcINT(203));

        test.addTestCase("MAIN.g_plant.stamps[1]", new PlcDATE_AND_TIME(java.time.LocalDateTime.parse("2025-06-01T00:00:00")));
        test.addTestCase("MAIN.g_plant.stamps[2]", new PlcDATE_AND_TIME(java.time.LocalDateTime.parse("2025-06-01T12:00:00")));
        test.addTestCase("MAIN.g_plant.tags[1]",   new PlcSTRING("MAIN.speed"));
        test.addTestCase("MAIN.g_plant.tags[3]",   new PlcSTRING("MAIN.pressure"));
        test.addTestCase("MAIN.g_plant.wtags[1]",  new PlcWSTRING("Δv"));
        test.addTestCase("MAIN.g_plant.wtags[2]",  new PlcWSTRING("Ölstand"));

        // Whole nested struct "meta" again as PlcStruct (sanity)
        {
            Map<String, PlcValue> meta = new HashMap<>();
            meta.put("s8", new PlcSINT(-1));   meta.put("u8", new PlcUSINT(255));
            meta.put("s16", new PlcINT(-2));   meta.put("u16", new PlcUINT(2));
            meta.put("s32", new PlcDINT(-3));  meta.put("u32", new PlcUDINT(3L));
            meta.put("s64", new PlcLINT(-4L)); meta.put("u64", new PlcULINT(4L));
            meta.put("b1", new PlcBOOL(true)); meta.put("b8", new PlcBYTE(0xAA));
            meta.put("b16", new PlcWORD(0xABCD)); meta.put("b32", new PlcDWORD(0x01020304L));
            meta.put("b64", new PlcLWORD(0x0A0B0C0D0E0F1011L));
            meta.put("r32", new PlcREAL(12.5f)); meta.put("r64", new PlcLREAL(-98.765d));
            meta.put("tim", new PlcTIME(java.time.Duration.parse("PT2S")));
            meta.put("ltim", new PlcLTIME(java.time.Duration.parse("PT2.1S")));
            meta.put("dat", new PlcDATE(java.time.LocalDate.parse("2025-01-01")));
            meta.put("timoday", new PlcTIME_OF_DAY(java.time.LocalTime.parse("12:34:56.789")));
            // TODO: In the PLC there is a ".007" at the end ...
            meta.put("dattim", new PlcDATE_AND_TIME(java.time.LocalDateTime.parse("2025-02-03T04:05:06")));
            meta.put("str", new PlcSTRING("meta-ok"));
            meta.put("wstr", new PlcWSTRING("Meta-OK"));
            test.addTestCase("MAIN.g_plant.meta", new PlcStruct(meta));
        }

        // =========================================================
        // 6) Multidimensional array of structs: g_chanGrid[2,2] of TChannel
        //    - element field reads
        //    - partial dimensions => PlcList of structs / lists
        // =========================================================
        test.addTestCase("MAIN.g_chanGrid[1][1].name",         new PlcSTRING("A1"));
        test.addTestCase("MAIN.g_chanGrid[1][1].enabled",      new PlcBOOL(true));
        test.addTestCase("MAIN.g_chanGrid[1][1].lut[2][2][2]",   new PlcINT(8));

        test.addTestCase("MAIN.g_chanGrid[1][2].id",           new PlcUDINT(11L));
        test.addTestCase("MAIN.g_chanGrid[1][2].enabled",      new PlcBOOL(false));
        test.addTestCase("MAIN.g_chanGrid[1][2].lut[2][2][2]",   new PlcINT(17));

        test.addTestCase("MAIN.g_chanGrid[2][1].name",         new PlcSTRING("B1"));
        test.addTestCase("MAIN.g_chanGrid[2][1].lut[2][2][1]",   new PlcINT(27));

        test.addTestCase("MAIN.g_chanGrid[2][2].id",           new PlcUDINT(21L));
        test.addTestCase("MAIN.g_chanGrid[2][2].enabled",      new PlcBOOL(true));

        // Fix first dimension => row of 2 TChannel as List<PlcStruct-like> (driver may surface as PlcList of nested values)
        {
            // We only assert a couple of representative subfields in each element of the PlcList row.
            test.addTestCase("MAIN.g_chanGrid[1]",
                new PlcList(java.util.Arrays.asList(
                    // element [1,1]
                    new PlcStruct(new HashMap<>() {{
                        put("id", new PlcUDINT(10L));
                        put("name", new PlcSTRING("A1"));
                        put("enabled", new PlcBOOL(true));
                        put("setpoints", new PlcList(List.of(
                            new PlcREAL(11.1),
                            new PlcREAL(11.2),
                            new PlcREAL(11.3),
                            new PlcREAL(11.4)
                        )));
                        put("lut",
                            new PlcList(
                                List.of(
                                    new PlcList(
                                        List.of(
                                            new PlcList(
                                                List.of(
                                                    new PlcINT(1),
                                                    new PlcINT(2)
                                                )
                                            ),
                                            new PlcList(
                                                List.of(
                                                    new PlcINT(3),
                                                    new PlcINT(4)
                                                )
                                            )
                                        )
                                    ),
                                    new PlcList(
                                        List.of(
                                            new PlcList(
                                                List.of(
                                                    new PlcINT(5),
                                                    new PlcINT(6)
                                                )
                                            ),
                                            new PlcList(
                                                List.of(
                                                    new PlcINT(7),
                                                    new PlcINT(8)
                                                )
                                            )
                                        )
                                    )
                                )
                            )
                        );
                    }}),
                    // element [1,2]
                    new PlcStruct(new HashMap<>() {{
                        put("id", new PlcUDINT(11L));
                        put("name", new PlcSTRING("A2"));
                        put("enabled", new PlcBOOL(false));
                        put("setpoints", new PlcList(List.of(
                            new PlcREAL(12.1),
                            new PlcREAL(12.2),
                            new PlcREAL(12.3),
                            new PlcREAL(12.4)
                        )));
                        put("lut",
                            new PlcList(
                                List.of(
                                    new PlcList(
                                        List.of(
                                            new PlcList(
                                                List.of(
                                                    new PlcINT(10),
                                                    new PlcINT(11)
                                                )
                                            ),
                                            new PlcList(
                                                List.of(
                                                    new PlcINT(12),
                                                    new PlcINT(13)
                                                )
                                            )
                                        )
                                    ),
                                    new PlcList(
                                        List.of(
                                            new PlcList(
                                                List.of(
                                                    new PlcINT(14),
                                                    new PlcINT(15)
                                                )
                                            ),
                                            new PlcList(
                                                List.of(
                                                    new PlcINT(16),
                                                    new PlcINT(17)
                                                )
                                            )
                                        )
                                    )
                                )
                            )
                        );
                    }})
                ))
            );
        }

        // Full 2x2 matrix of channels as List<List<PlcStruct>> — keep it light, assert subset
        {
            PlcStruct ch11 = new PlcStruct(new HashMap<>() {{
                put("id", new PlcUDINT(10L));
                put("name", new PlcSTRING("A1"));
                put("enabled", new PlcBOOL(true));
                put("setpoints", new PlcList(List.of(
                    new PlcREAL(11.1),
                    new PlcREAL(11.2),
                    new PlcREAL(11.3),
                    new PlcREAL(11.4)
                )));
                put("lut",
                    new PlcList(
                        List.of(
                            new PlcList(
                                List.of(
                                    new PlcList(
                                        List.of(
                                            new PlcINT(1),
                                            new PlcINT(2)
                                        )
                                    ),
                                    new PlcList(
                                        List.of(
                                            new PlcINT(3),
                                            new PlcINT(4)
                                        )
                                    )
                                )
                            ),
                            new PlcList(
                                List.of(
                                    new PlcList(
                                        List.of(
                                            new PlcINT(5),
                                            new PlcINT(6)
                                        )
                                    ),
                                    new PlcList(
                                        List.of(
                                            new PlcINT(7),
                                            new PlcINT(8)
                                        )
                                    )
                                )
                            )
                        )
                    )
                );
            }});
            PlcStruct ch12 = new PlcStruct(new HashMap<>() {{
                put("id", new PlcUDINT(11L));
                put("name", new PlcSTRING("A2"));
                put("enabled", new PlcBOOL(false));
                put("setpoints", new PlcList(List.of(
                    new PlcREAL(12.1),
                    new PlcREAL(12.2),
                    new PlcREAL(12.3),
                    new PlcREAL(12.4)
                )));
                put("lut",
                    new PlcList(
                        List.of(
                            new PlcList(
                                List.of(
                                    new PlcList(
                                        List.of(
                                            new PlcINT(10),
                                            new PlcINT(11)
                                        )
                                    ),
                                    new PlcList(
                                        List.of(
                                            new PlcINT(12),
                                            new PlcINT(13)
                                        )
                                    )
                                )
                            ),
                            new PlcList(
                                List.of(
                                    new PlcList(
                                        List.of(
                                            new PlcINT(14),
                                            new PlcINT(15)
                                        )
                                    ),
                                    new PlcList(
                                        List.of(
                                            new PlcINT(16),
                                            new PlcINT(17)
                                        )
                                    )
                                )
                            )
                        )
                    )
                );
            }});
            PlcStruct ch21 = new PlcStruct(new HashMap<>() {{
                put("id", new PlcUDINT(20L));
                put("name", new PlcSTRING("B1"));
                put("enabled", new PlcBOOL(true));
                put("setpoints", new PlcList(List.of(
                    new PlcREAL(21.1),
                    new PlcREAL(21.2),
                    new PlcREAL(21.3),
                    new PlcREAL(21.4)
                )));
                put("lut",
                    new PlcList(
                        List.of(
                            new PlcList(
                                List.of(
                                    new PlcList(
                                        List.of(
                                            new PlcINT(21),
                                            new PlcINT(22)
                                        )
                                    ),
                                    new PlcList(
                                        List.of(
                                            new PlcINT(23),
                                            new PlcINT(24)
                                        )
                                    )
                                )
                            ),
                            new PlcList(
                                List.of(
                                    new PlcList(
                                        List.of(
                                            new PlcINT(25),
                                            new PlcINT(26)
                                        )
                                    ),
                                    new PlcList(
                                        List.of(
                                            new PlcINT(27),
                                            new PlcINT(28)
                                        )
                                    )
                                )
                            )
                        )
                    )
                );
            }});
            PlcStruct ch22 = new PlcStruct(new HashMap<>() {{
                put("id", new PlcUDINT(21L));
                put("name", new PlcSTRING("B2"));
                put("enabled", new PlcBOOL(true));
                put("setpoints", new PlcList(List.of(
                    new PlcREAL(22.1),
                    new PlcREAL(22.2),
                    new PlcREAL(22.3),
                    new PlcREAL(22.4)
                )));
                put("lut",
                    new PlcList(
                        List.of(
                            new PlcList(
                                List.of(
                                    new PlcList(
                                        List.of(
                                            new PlcINT(31),
                                            new PlcINT(32)
                                        )
                                    ),
                                    new PlcList(
                                        List.of(
                                            new PlcINT(33),
                                            new PlcINT(34)
                                        )
                                    )
                                )
                            ),
                            new PlcList(
                                List.of(
                                    new PlcList(
                                        List.of(
                                            new PlcINT(35),
                                            new PlcINT(36)
                                        )
                                    ),
                                    new PlcList(
                                        List.of(
                                            new PlcINT(37),
                                            new PlcINT(38)
                                        )
                                    )
                                )
                            )
                        )
                    )
                );
            }});
            test.addTestCase("MAIN.g_chanGrid",
                new PlcList(java.util.Arrays.asList(
                    new PlcList(java.util.Arrays.asList(ch11, ch12)),
                    new PlcList(java.util.Arrays.asList(ch21, ch22))
                ))
            );
        }

        // =========================================================
        // 7) Mixed addressing (struct-in-array, array-in-struct, etc.)
        // =========================================================
        test.addTestCase("MAIN.g_plant.channels[1].setpoints[1]", new PlcREAL(0.0f));
        test.addTestCase("MAIN.g_plant.channels[2].lut[1][2][2]",   new PlcINT(103));
        test.addTestCase("MAIN.g_matrixI16.m[2][1]",               new PlcINT(200));
        // =========================================================
        // Run
        // =========================================================
        long startMillis = System.currentTimeMillis();
        test.run();
        long endMillis = System.currentTimeMillis();
        System.out.println("Test executed in " + (endMillis - startMillis) + "ms");
    }

}
