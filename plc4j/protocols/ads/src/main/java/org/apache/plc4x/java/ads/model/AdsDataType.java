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
package org.apache.plc4x.java.ads.model;

import org.apache.commons.lang3.ArrayUtils;

import java.io.ByteArrayOutputStream;
import java.math.BigInteger;
import java.time.Duration;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.time.temporal.ChronoUnit;
import java.util.stream.IntStream;

/**
 * Documentation can be found here:
 *
 * @see <a href="https://infosys.beckhoff.com/english.php?content=../content/1033/tcsystemmanager/basics/TcSysMgr_DatatypeComparison.htm&id=">TwinCAT System Manager - I/O Variables</a>
 * @see <a href="https://infosys.beckhoff.com/english.php?content=../content/1033/tcplccontrol/html/tcplcctrl_plc_data_types_overview.htm&id">TwinCAT PLC Control: Data Types</a>
 */
public enum AdsDataType {
    // TODO: maybe this are just types for the plc ide and can be removed
    // https://infosys.beckhoff.com/english.php?content=../content/1033/tcsystemmanager/basics/TcSysMgr_DatatypeComparison.htm&id=
    BIT(8),
    BIT8(8),
    BITARR8(8),
    BITARR16(16),
    BITARR32(32),
    INT8(8),
    INT16(Short.MIN_VALUE, Short.MAX_VALUE, 16),
    INT32(Integer.MIN_VALUE, Integer.MAX_VALUE, 32),
    INT64(Long.MIN_VALUE, Long.MAX_VALUE, 64),
    UINT8(0, Short.MAX_VALUE, 8),
    UINT16(0, Integer.MAX_VALUE, 16),
    UINT32(0, produceUnsignedMaxValue(32), 32),
    UINT64(0, produceUnsignedMaxValue(64), 64),
    FLOAT(Float.MIN_VALUE, Float.MAX_VALUE, 32),
    DOUBLE(Double.MIN_VALUE, Double.MAX_VALUE, 64),
    // https://infosys.beckhoff.com/english.php?content=../content/1033/tcplccontrol/html/tcplcctrl_plc_data_types_overview.htm&id
    // Standard Data Types
    /**
     * BOOL type variables may be given the values TRUE and FALSE.
     * <p>
     * Type	Memory use
     * BOOL	8 Bit
     * Note:
     * <p>
     * A BOOL type variable is true, if the least significant bit in the memory is set (e.g. 2#00000001 ). If no bit is set in the memory, the variable is FALSE (2#00000000). All other values can´t be interpeted accurately and be displayed (***INVALID: 16#xy *** in the Online View). Such problems may appear, if for example overlapped memory ranges are used in the PLC program.
     * <p>
     * Example:
     * <p>
     * The boolean variable is in the same memory range as the byte variable.
     */
    BOOL(8),
    /**
     * BYTE
     * <p>
     * Integer data type.
     * <p>
     * Type	Lower bound	Upper bound	Memory use
     * BYTE	0	255	8 Bit
     */
    BYTE(0, 255, 8),
    /**
     * WORD
     * Integer data type.
     * <p>
     * Type	Lower bound	Upper bound	Memory use
     * WORD	0	65535	16 Bit
     */
    WORD(0, 65535, 16),
    /**
     * DWORD
     * Integer data type.
     * <p>
     * Type	Lower bound	Upper bound	Memory use
     * DWORD	0	4294967295	32 Bit
     */
    DWORD(0, 4294967295L, 32),
    /**
     * SINT
     * (Short) signed integer data type.
     * <p>
     * Type	Lower bound	Upper bound	Memory use
     * SINT	-128	127	8 Bit
     */
    SINT(-128, 127, 8),
    /**
     * USINT
     * Unsigned (short) integer data type.
     * <p>
     * Type	Lower bound	Upper bound	Memory use
     * USINT	0	255	8 Bit
     */
    USINT(0, 255, 8),
    /**
     * INT
     * Signed integer data type.
     * <p>
     * Type	Lower bound	Upper bound	Memory use
     * INT	-32768	32767	16 Bit
     */
    INT(-32768, 32767, 16),
    /**
     * UINT
     * Unsigned integer data type.
     * <p>
     * Type	Lower bound	Upper bound	Memory use
     * UINT	0	65535	16 Bit
     */
    UINT(0, 65535, 16),
    /**
     * DINT
     * Signed integer data type.
     * <p>
     * Type	Lower bound	Upper bound	Memory use
     * DINT	-2147483648	2147483647	32 Bit
     */
    DINT(-2147483648, 2147483647, 32),
    /**
     * UDINT
     * Unsigned integer data type.
     * <p>
     * Type	Lower bound	Upper bound	Memory use
     * UDINT	0	4294967295	32 Bit
     */
    UDINT(0, 4294967295L, 32),
    /**
     * LINT  (64 bit integer, currently not supported by TwinCAT)
     */
    LINT(64),
    /**
     * ULINT (Unsigned 64 bit integer, currently not supported by TwinCAT)
     */
    ULINT(64),
    /**
     * REAL
     * 32 Bit floating point data type. It is required to represent rational numbers.
     * <p>
     * Type	Lower bound	Upper bound	Memory use
     * REAL	~ -3.402823 x 1038	~ 3.402823 x 1038	32 Bit
     */
    REAL(Float.MIN_VALUE, Float.MAX_VALUE, 32),
    /**
     * LREAL
     * 64 Bit floating point data type. It is required to represent rational numbers.
     * <p>
     * Type	Lower bound	Upper bound	Memory use
     * LREAL	~ -1.79769313486231E308	~ 1.79769313486232E308	64 Bit
     */
    LREAL(Double.MIN_VALUE, Double.MAX_VALUE, 64),
    /**
     * STRING
     * A STRING type variable can contain any string of characters. The size entry in the declaration determines how much memory space should be reserved for the variable. It refers to the number of characters in the string and can be placed in parentheses or square brackets.
     * <p>
     * Example of a string declaration with 35 characters:
     * <p>
     * str:STRING(35):='This is a String';
     * Type	Memory use
     * STRING
     * If no size specification is given, the default size of 80 characters will be used: Memory use [Bytes] =  80 + 1 Byte for string terminated Null character;
     * If string size specification is given: Memory use [Bytes] = String Size + 1 Byte for string terminated Null character);
     */
    STRING(81 * 8),
    /**
     * TIME
     * Duration time. The most siginificant digit is one millisecond. The data type is handled internally like DWORD.
     * <p>
     * Type	Lower bound	Upper bound	Memory use
     * TIME	T#0ms	T#71582m47s295ms	32 Bit
     */
    TIME(0, Duration.ofMinutes(71582).plusSeconds(47).plusMillis(295).toMillis(), 32),
    /**
     * TIME_OF_DAY
     * TOD
     * Time of day. The most siginificant digit is one millisecond. The data type is handled internally like DWORD.
     * <p>
     * Type	Lower bound	Upper bound	Memory use
     * TIME_OF_DAY
     * TOD
     * <p>
     * TOD#00:00	TOD#1193:02:47.295	32 Bit
     *///TODO: strange maximum
    TIME_OF_DAY(0, ChronoUnit.MILLIS.between(LocalTime.of(0, 0), LocalTime.of(23, 59, 59)), 32),
    /**
     * DATE
     * Date. The most significant digit is one second. The data type is handled internally like DWORD.
     * <p>
     * Type	Lower bound	Upper bound	Memory use
     * DATE	D#1970-01-01	D#2106-02-06	32 Bit
     */
    DATE(0, ChronoUnit.SECONDS.between(LocalDateTime.of(1970, 1, 1, 0, 0), LocalDateTime.of(2106, 2, 6, 0, 0)), 32),
    /**
     * DATE_AND_TIME
     * DT
     * Date and time. The most siginificant digit is one second. The data type is handled internally like DWORD.
     * <p>
     * Type	Lower  bound	Upper  bound	Memory use
     * DATE_AND_TIME
     * DT
     * <p>
     * DT#1970-01-01-00:00	DT#2106-02-06-06:28:15	32 Bit
     *////TODO: calculate max
    DATE_AND_TIME(0, -1, 32),
    //User-defined Data Types
    /**
     * Arrays
     * One-, two-, and three-dimensional fields (arrays) are supported as elementary data types. Arrays can be defined both in the declaration part of a POU and in the global variable lists.
     * <p>
     * Syntax:
     * <p>
     * &lt;Field_Name&gt;:ARRAY [&lt;LowLim1&gt;..&lt;UpLim1&gt;, &lt;LowLim2&gt;..&lt;UpLim2&gt;] OF &lt;elem. Type&gt;
     * <p>
     * LowLim1, LowLim2 identify the lower limit of the field range; UpLim1 and UpLim2 identify the upper limit. The range values must be integers.
     * <p>
     * Example:
     * <p>
     * Card_game: ARRAY [1..13, 1..4] OF INT;
     * <p>
     * <p>
     * <p>
     * Initializing of Arrays
     * You can initialize either all of the elements in an array or none of them.
     * <p>
     * Example for initializing arrays:
     * <p>
     * arr1 : ARRAY [1..5] OF INT := 1,2,3,4,5;
     * arr2 : ARRAY [1..2,3..4] OF INT := 1,3(7); (* short for 1,7,7,7 *)
     * arr3 : ARRAY [1..2,2..3,3..4] OF INT := 2(0),4(4),2,3; (* short for 0,0,4,4,4,4,2,3 *)
     * <p>
     * <p>
     * <p>
     * Example for the initialization of an array of a structure:
     * <p>
     * TYPE STRUCT1
     * STRUCT
     * p1:int;
     * p2:int;
     * p3:dword;
     * END_STRUCT
     * arr1 : ARRAY[1..3] OF STRUCT1:= (p1:=1,p2:=10,p3:=4723), (p1:=2,p2:=0,p3:=299), (p1:=14,p2:=5,p3:=112);
     * <p>
     * <p>
     * Example of the partial initialization of an Array:
     * <p>
     * arr1 : ARRAY [1..10] OF INT := 1,2;
     * Elements to which no value is pre-assigned are initialized with the default initial value of the basic type. In the example above, the elements arr1[3]  to arr1[10] are therefore initialized with 0.
     * <p>
     * <p>
     * <p>
     * Array components are accessed in a two-dimensional array using the following syntax:
     * <p>
     * &lt;Field_Name&gt;[Index1,Index2]
     * <p>
     * Example:
     * <p>
     * Card_game[9,2]
     * <p>
     * <p>
     * <p>
     * Note:
     * <p>
     * If you define a function in your project with the name CheckBounds, you can automatically check for out-of-range errors in arrays ! The name of the function is fixed and can only have this designation.
     */
    ARRAY(-1),//TODO: implement me
    /**
     * Pointer
     * Variable or function block addresses are saved in pointers while a program is running. Pointer declarations have the following syntax:
     * <p>
     * &lt;Identifier&gt;: POINTER TO &lt;Datatype/Functionblock&gt;;
     * A pointer can point to any data type or function block even to user-defined types. The function of the Address Operator ADR is to assign the address of a variable or function block to the pointer.
     * A pointer can be dereferenced by adding the content operator "^" after the pointer identifier. With the help of the SIZEOF Operator, e.g. a pointer increment can be done.
     * <p>
     * <p>
     * Please note: A pointer is counted up byte-wise ! You can get it counted up like it is usual in the C-Compiler by using the instruction p=p+SIZEOF(p^);.
     * <p>
     * <p>
     * <p>
     * Attention:
     * After an Online Change there might be changes concerning the data on certain addresses. Please regard this in case of using pointers on addresses.
     * <p>
     * <p>
     * <p>
     * Example:
     * <p>
     * pt:POINTER TO INT;
     * var_int1:INT := 5;
     * var_int2:INT;
     * <p>
     * <p>
     * pt := ADR(var_int1);
     * var_int2:= pt^; (* var_int2 is now 5 *)
     * <p>
     * <p>
     * Example 2 (Pointer increment):
     * <p>
     * ptByCurrDataOffs : POINTER TO BYTE;
     * udiAddress       : UDINT;
     * <p>
     * <p>
     * <p>
     * (*--- pointer increment ---*)
     * udiAddress := ptByCurrDataOffs;
     * udiAddress := udiAddress + SIZEOF(ptByCurrDataOffs^);
     * ptByCurrDataOffs := udiAddress;
     * (* -- end of pointer increment ---*)
     */
    POINTER(-1),//TODO: implement me,
    /**
     * Enumeration (ENUM)
     * Enumeration is a user-defined data type that is made up of a number of string constants. These constants are referred to as enumeration values. Enumeration values are recognized in all areas of the project even if they were locally declared within aPOU. It is best to create your enumerations as objects in the Object Organizer under the register card Data types. They begin with the keyword TYPE and end with END_TYPE.
     * <p>
     * Syntax:
     * <p>
     * TYPE &lt;Identifier&gt;:(&lt;Enum_0&gt; ,&lt;Enum_1&gt;, ...,&lt;Enum_n&gt;);END_TYPE
     * <p>
     * The &lt;Identifier&gt; can take on one of the enumeration values and will be initialized with the first one. These values are compatible with whole numbers which means that you can perform operations with them just as you would with INT. You can assign a number x to the &lt;Identifier&gt;. If the enumeration values are not initialized, counting will begin with 0. When initializing, make certain the initial values are increasing. The validity of the number will be reviewed at the time it is run.
     * <p>
     * Example:
     * <p>
     * TRAFFIC_SIGNAL: (Red, Yellow, Green:=10); (*The initial value for each of the colors is red 0, yellow 1, green 10 *)
     * TRAFFIC_SIGNAL:=0; (* The value of the traffic signal is red*)
     * FOR i:= Red TO Green DO
     * i := i + 1;
     * END_FOR;
     * <p>
     * You may not use the same enumeration value more than once.
     * <p>
     * Example:
     * <p>
     * TRAFFIC_SIGNAL: (red, yellow, green);
     * COLOR: (blue, white, red);
     * <p>
     * Error: red may not be used for both TRAFFIC_SIGNAL and COLOR.
     */
    ENUM(-1),//TODO: implement me,
    /**
     * Structures (STRUCT)
     * Structures are created as objects in the Object Organizer under the register card Data types. They begin with the keyword TYPE and end with END_TYPE.The syntax for structure declarations is as follows:
     * <p>
     * TYPE &lt;Structurename&gt;:
     * STRUCT
     * &lt;Declaration of Variables 1&gt;
     * .
     * .
     * &lt;Declaration of Variables n&gt;
     * END_STRUCT
     * END_TYPE
     * <p>
     * &lt;Structurename&gt; is a type that is recognized throughout the project and can be used like a standard data type. Interlocking structures are allowed. The only restriction is that variables may not be placed at addresses (the AT declaration is not allowed!).
     * <p>
     * Example for a structure definition named Polygonline:
     * <p>
     * TYPE Polygonline:
     * STRUCT
     * Start:ARRAY [1..2] OF INT;
     * Point1:ARRAY [1..2] OF INT;
     * Point2:ARRAY [1..2] OF INT;
     * Point3:ARRAY [1..2] OF INT;
     * Point4:ARRAY [1..2] OF INT;
     * End:ARRAY [1..2] OF INT;
     * END_STRUCT
     * END_TYPE
     * <p>
     * You can gain access to structure components using the following syntax:
     * <p>
     * &lt;Structure_Name&gt;.&lt;Componentname&gt;
     * <p>
     * For example, if you have a structure named "Week" that contains a component named "Monday", you can get to it by doing the following: Week.Monday
     * <p>
     * <p>
     * <p>
     * Note:
     * Due to different alignments, structures and arrays may have different configurations and sizes on different hardware platforms (e.g. CX1000 and CX90xx).
     * <p>
     * During data exchange the size and structure alignment must be identical!
     * <p>
     * <p>
     * <p>
     * Example for a structure definition with name ST_ALIGN_SAMPLE:
     * <p>
     * TYPE ST_ALIGN_SAMPLE:
     * STRUCT
     * _diField1   : DINT;
     * _byField1   : BYTE;
     * _iField     : INT;
     * _byField2   : BYTE;
     * _diField2   : DINT;
     * _pField     : POINTER TO BYTE;
     * END_STRUCT
     * END_TYPE
     * <p>
     * On CX90xx (RISC) platforms the member components of structure ST_ALIGN_SAMPLE have the following sizes and offsets:
     * <p>
     * _diField1 (DINT), Offset = 0 (16#0),   Size = 4
     * _byField1 (BYTE), Offset = 4 (16#4),   Size = 1
     * _iField (INT), Offset = 6 (16#6),   Size = 2
     * _byField2 (BYTE), Offset = 8 (16#8),   Size = 1
     * _diField2 (DINT), Offset = 12 (16#C),  Size = 4
     * _pField (POINTER TO BYTE), Offset = 16 (16#10), Size = 4
     * <p>
     * Overall size through natural alignment with Pack(4) and so-called padding bytes: 20
     * <p>
     * <p>
     * <p>
     * On CX10xx platforms the member components of structure ST_ALIGN_SAMPLE have the following sizes and offsets:
     * <p>
     * _diField1 (DINT), Offset = 0 (16#0),   Size = 4
     * _byField1 (BYTE), Offset = 4 (16#4),   Size = 1
     * _iField (INT), Offset = 5 (16#5),   Size = 2
     * _byField2 (BYTE), Offset = 7 (16#7),   Size = 1
     * _diField2 (DINT), Offset = 8 (16#8),  Size = 4
     * _pField (POINTER TO BYTE), Offset = 12 (16#C), Size = 4
     * <p>
     * Overall size: 16
     * <p>
     * <p>
     * <p>
     * Display of structure ST_ALIGN_SAMPLE for CX90xx platforms (RISC) with representation of the padding bytes:
     * <p>
     * TYPE ST_ALIGN_SAMPLE:
     * STRUCT
     * _diField1    : DINT;
     * _byField1    : BYTE;
     * _byPadding   : BYTE;
     * _iField      : INT;
     * _byField2    : BYTE;
     * _a_byPadding : ARRAY[0..2] OF BYTE;
     * _diField2    : DINT;
     * _pField      : POINTER TO BYTE;
     * END_STRUCT
     * END_TYPE
     */
    STRUCT(-1),//TODO: implement me,
    /**
     * References (Alias types)
     * You can use the user-defined derived data type to create an alternative name for a variable, constant or function block. Create your references as objects in the Object Organizer under the register card Data types. They begin with the keyword TYPE and end with END_TYPE.
     * <p>
     * Syntax:
     * <p>
     * TYPE &lt;Identifier&gt;: &lt;Assignment term&gt;;
     * END_TYPE
     * <p>
     * Example:
     * <p>
     * TYPE message:STRING[50];
     * END_TYPE;
     */
    ALIAS(-1),//TODO: implement me,
    /**
     * Subrange types
     * A sub-range data type is a type whose range of values is only a subset of that of the basic type. The declaration can be carried out in the data types register, but a variable can also be directly declared with a subrange type:
     * Syntax for the declaration in the 'Data types' register:
     * <p>
     * TYPE &lt;Name&gt; : &lt;Inttype&gt; (&lt;ug&gt;..&lt;og&gt;) END_TYPE;
     * Type	Description
     * &lt;Name&gt;	must be a valid IEC identifier
     * &lt;Inttype&gt;	is one of the data types SINT, USINT, INT, UINT, DINT, UDINT, BYTE, WORD, DWORD (LINT, ULINT, LWORD).
     * &lt;ug&gt;	Is a constant which must be compatible with the basic type and which sets the lower boundary of the range types. The lower boundary itself is included in this range.
     * &lt;og&gt;	Is a constant that must be compatible with the basic type, and sets the upper boundary of the range types. The upper boundary itself is included in this basic type.
     * Example:
     * <p>
     * TYPE
     * SubInt : INT (-4095..4095);
     * END_TYPE
     * Direct declaration of a variable with a subrange type:
     * <p>
     * VAR
     * i1 : INT (-4095..4095);
     * i2: INT (5...10):=5;
     * ui : UINT (0..10000);
     * END_VAR
     * If a constant is assigned to a subrange type (in the declaration or in the implementation) that does not apply to this range (e.g. 1:=5000), an error message is issued.
     * In order to check for observance of range boundaries at runtime, the functions CheckRangeSigned or CheckRangeUnsigned must be introduced.
     */
    SUB_RANGE_DATA_TYPE(-1),//TODO: implement me,

    UNKNOWN(-1);

    private final String typeName;

    private final double lowerBound;

    private final double upperBound;

    private final int memoryUse;

    private final int targetByteSize;

    // TODO: BYTE.MAX default might not be the best....
    AdsDataType(int memoryUse) {
        this(0, Byte.MAX_VALUE, memoryUse);
    }

    AdsDataType(double lowerBound, double upperBound, int memoryUse) {
        this.lowerBound = lowerBound;
        this.upperBound = upperBound;
        this.typeName = name();
        this.memoryUse = memoryUse;
        this.targetByteSize = this.memoryUse / 8;
    }

    public String getTypeName() {
        return typeName;
    }

    public double getLowerBound() {
        return lowerBound;
    }

    public double getUpperBound() {
        return upperBound;
    }

    public int getMemoryUse() {
        return memoryUse;
    }

    public int getTargetByteSize() {
        return targetByteSize;
    }

    public boolean withinBounds(double other) {
        return other >= lowerBound && other <= upperBound;
    }

    @Override
    public String toString() {
        return "AdsDataType{" +
            "typeName='" + typeName + '\'' +
            ", lowerBound=" + lowerBound +
            ", upperBound=" + upperBound +
            ", memoryUse=" + memoryUse +
            ", targetByteSize=" + targetByteSize +
            "} " + super.toString();
    }

    private static double produceUnsignedMaxValue(int numberOfBytes) {
        return new BigInteger(
            ArrayUtils.insert(
                0,
                IntStream.range(0, numberOfBytes)
                    .map(ignore -> 0xff)
                    .collect(
                        ByteArrayOutputStream::new,
                        (baos, i) -> baos.write((byte) i),
                        (baos1, baos2) -> baos1.write(baos2.toByteArray(), 0, baos2.size())
                    )
                    .toByteArray(),
                (byte) 0x0)
        ).doubleValue();
    }
}
