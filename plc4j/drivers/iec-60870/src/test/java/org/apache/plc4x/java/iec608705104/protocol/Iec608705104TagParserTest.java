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
package org.apache.plc4x.java.iec608705104.protocol;

import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.iec608705104.readwrite.BinaryCounterReading;
import org.apache.plc4x.java.iec608705104.readwrite.CauseOfInitialization;
import org.apache.plc4x.java.iec608705104.readwrite.DoubleCommand;
import org.apache.plc4x.java.iec608705104.readwrite.DoublePointInformation;
import org.apache.plc4x.java.iec608705104.readwrite.FixedTestBitPatternTwoOctet;
import org.apache.plc4x.java.iec608705104.readwrite.InformationObject;
import org.apache.plc4x.java.iec608705104.readwrite.InformationObjectWithoutTime_DOUBLE_COMMAND;
import org.apache.plc4x.java.iec608705104.readwrite.InformationObjectWithoutTime_DOUBLE_POINT_INFORMATION;
import org.apache.plc4x.java.iec608705104.readwrite.InformationObjectWithoutTime_END_OF_INITIALISATION;
import org.apache.plc4x.java.iec608705104.readwrite.InformationObjectWithoutTime_INTEGRATED_TOTALS;
import org.apache.plc4x.java.iec608705104.readwrite.InformationObjectWithoutTime_INTERROGATION_COMMAND;
import org.apache.plc4x.java.iec608705104.readwrite.InformationObjectWithoutTime_MEASURED_VALUE_NORMALISED_VALUE;
import org.apache.plc4x.java.iec608705104.readwrite.InformationObjectWithoutTime_MEASURED_VALUE_SCALED_VALUE;
import org.apache.plc4x.java.iec608705104.readwrite.InformationObjectWithoutTime_MEASURED_VALUE_SHORT_FLOATING_POINT_NUMBER;
import org.apache.plc4x.java.iec608705104.readwrite.InformationObjectWithoutTime_SINGLE_POINT_INFORMATION;
import org.apache.plc4x.java.iec608705104.readwrite.InformationObjectWithoutTime_STEP_POSITION_INFORMATION;
import org.apache.plc4x.java.iec608705104.readwrite.InformationObjectWithoutTime_TEST_COMMAND;
import org.apache.plc4x.java.iec608705104.readwrite.NormalizedValue;
import org.apache.plc4x.java.iec608705104.readwrite.QualifierOfCommand;
import org.apache.plc4x.java.iec608705104.readwrite.QualifierOfCounterInterrogationCommand;
import org.apache.plc4x.java.iec608705104.readwrite.QualifierOfInterrogation;
import org.apache.plc4x.java.iec608705104.readwrite.QualifierOfParameterActivation;
import org.apache.plc4x.java.iec608705104.readwrite.QualifierOfParameterOfMeasuredValues;
import org.apache.plc4x.java.iec608705104.readwrite.QualifierOfResetProcessCommand;
import org.apache.plc4x.java.iec608705104.readwrite.QualifierOfSetPointCommand;
import org.apache.plc4x.java.iec608705104.readwrite.QualityDescriptor;
import org.apache.plc4x.java.iec608705104.readwrite.QualityDescriptorForPointsOfProtectionEquipment;
import org.apache.plc4x.java.iec608705104.readwrite.ScaledValue;
import org.apache.plc4x.java.iec608705104.readwrite.SingleCommand;
import org.apache.plc4x.java.iec608705104.readwrite.SingleEventOfProtectionEquipment;
import org.apache.plc4x.java.iec608705104.readwrite.SinglePointInformation;
import org.apache.plc4x.java.iec608705104.readwrite.TypeIdentification;
import org.apache.plc4x.java.iec608705104.readwrite.ValueWithTransientStateIndication;
import org.apache.plc4x.java.iec608705104.readwrite.OutputCircuitInformation;
import org.apache.plc4x.java.iec608705104.protocol.Iec608705104TagParser;
import org.apache.plc4x.java.spi.values.PlcBOOL;
import org.apache.plc4x.java.spi.values.PlcList;
import org.apache.plc4x.java.spi.values.PlcREAL;
import org.apache.plc4x.java.spi.values.PlcStruct;
import org.apache.plc4x.java.spi.values.PlcUDINT;
import org.apache.plc4x.java.spi.values.PlcUINT;
import org.apache.plc4x.java.spi.values.PlcUSINT;
import org.apache.plc4x.java.spi.values.PlcValueAdapter;
import org.apache.plc4x.java.spi.values.PlcWORD;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertInstanceOf;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

class Iec608705104TagParserTest {

    private static QualityDescriptor goodQuality() {
        return new QualityDescriptor(false, false, false, false, false);
    }

    @Test
    void singlePointInformationCarriesQualityFlags() {
        SinglePointInformation siq = new SinglePointInformation(true, false, true, false, true);
        PlcValueAdapter v = Iec608705104TagParser.processSinglePointInformation(siq);
        assertTrue(v.getBoolean());
        assertTrue(v.getMetaData("invalid").getBoolean());
        assertTrue(v.getMetaData("substituted").getBoolean());
    }

    @Test
    void doublePointInformationSplitsBothCodeBits() {
        // dpiCode = 0b11 -> both points set.
        DoublePointInformation diq = new DoublePointInformation(false, false, false, false, (byte) 0x03);
        PlcValueAdapter v = Iec608705104TagParser.processDoublePointInformation(diq);
        PlcList list = assertInstanceOf(PlcList.class, v);
        assertTrue(list.getList().get(0).getBoolean());
        assertTrue(list.getList().get(1).getBoolean());
    }

    @Test
    void valueWithTransientStateIndicationCarriesValueAndFlag() {
        ValueWithTransientStateIndication vti = new ValueWithTransientStateIndication(true, (byte) 42);
        PlcValueAdapter v = Iec608705104TagParser.processValueWithTransientStateIndication(vti);
        assertEquals((short) 42, v.getShort());
        assertTrue(v.getMetaData("transientState").getBoolean());
    }

    @Test
    void normalizedAndScaledValueAreUnsignedIntegers() {
        assertInstanceOf(PlcUINT.class,
            Iec608705104TagParser.processNormalizedValue(new NormalizedValue(1234)));
        assertInstanceOf(PlcUINT.class,
            Iec608705104TagParser.processScaledValue(new ScaledValue((short) 5)));
    }

    @Test
    void binaryCounterReadingExposesCounterMetaData() {
        BinaryCounterReading bcr = new BinaryCounterReading(99L, true, false, true, (byte) 4);
        PlcValueAdapter v = Iec608705104TagParser.processBinaryCounterReading(bcr);
        assertInstanceOf(PlcUDINT.class, v);
        assertEquals(99L, v.getLong());
        assertTrue(v.getMetaData("counterValid").getBoolean());
        assertTrue(v.getMetaData("carry").getBoolean());
    }

    @Test
    void singleEventOfProtectionEquipmentExposesAllFlags() {
        SingleEventOfProtectionEquipment sep =
            new SingleEventOfProtectionEquipment(true, true, true, true, true, (byte) 2);
        PlcValueAdapter v = Iec608705104TagParser.processSingleEventOfProtectionEquipment(sep);
        assertEquals((short) 2, v.getShort());
        assertTrue(v.getMetaData("elapsedTimeInvalid").getBoolean());
    }

    @Test
    void outputCircuitInformationBecomesStruct() {
        OutputCircuitInformation oci = new OutputCircuitInformation(true, false, true, false);
        PlcValueAdapter v = Iec608705104TagParser.processOutputCircuitInformation(oci);
        PlcStruct struct = assertInstanceOf(PlcStruct.class, v);
        assertTrue(struct.getValue("stateOfOperationPhaseL1").getBoolean());
        assertTrue(struct.getValue("stateOfOperationPhaseL3").getBoolean());
    }

    @Test
    void qualityDescriptorAnnotatesPlcValue() {
        PlcValueAdapter v = PlcBOOL.of(true);
        Iec608705104TagParser.processQualityDescriptor(
            new QualityDescriptor(true, false, true, false, true), v);
        assertTrue(v.getMetaData("invalid").getBoolean());
        assertTrue(v.getMetaData("overflow").getBoolean());
    }

    @Test
    void qualityDescriptorForProtectionEquipmentAnnotatesPlcValue() {
        PlcValueAdapter v = PlcBOOL.of(false);
        Iec608705104TagParser.processQualityDescriptorForPointsOfProtectionEquipment(
            new QualityDescriptorForPointsOfProtectionEquipment(true, true, true, true, true), v);
        assertTrue(v.getMetaData("elapsedTimeInvalid").getBoolean());
    }

    @Test
    void interrogationQualifiersWrapTheirRawValues() {
        assertEquals(20, Iec608705104TagParser
            .processQualifierOfInterrogation(new QualifierOfInterrogation((short) 20)).getInt());
        PlcStruct counter = (PlcStruct) Iec608705104TagParser.processQualifierOfCounterInterrogationCommand(
            new QualifierOfCounterInterrogationCommand((byte) 1, (byte) 2));
        assertEquals((short) 1, counter.getValue("freeze").getShort());
        assertEquals((short) 2, counter.getValue("request").getShort());
    }

    @Test
    void parameterQualifiersAnnotateOrWrap() {
        PlcValueAdapter v = PlcBOOL.of(true);
        Iec608705104TagParser.processQualifierOfParameterOfMeasuredValues(
            new QualifierOfParameterOfMeasuredValues(true, false, (byte) 3), v);
        assertTrue(v.getMetaData("parameterInOperation").getBoolean());

        assertEquals(7, Iec608705104TagParser
            .processQualifierOfParameterActivation(new QualifierOfParameterActivation((short) 7))
            .getInt());
    }

    @Test
    void commandQualifiersAnnotateOrWrap() {
        PlcValueAdapter v = Iec608705104TagParser
            .processQualifierOfCommand(new QualifierOfCommand(true, (byte) 5));
        assertEquals((short) 5, v.getShort());
        assertTrue(v.getMetaData("select").getBoolean());

        assertEquals(9, Iec608705104TagParser
            .processQualifierOfResetProcessCommand(new QualifierOfResetProcessCommand((short) 9))
            .getInt());

        PlcValueAdapter target = PlcBOOL.of(false);
        Iec608705104TagParser.processQualifierOfSetPointCommand(
            new QualifierOfSetPointCommand(true, (byte) 1), target);
        assertTrue(target.getMetaData("select").getBoolean());
    }

    @Test
    void causeOfInitializationAnnotatesSelect() {
        PlcValueAdapter v = Iec608705104TagParser
            .processCauseOfInitialization(new CauseOfInitialization(true, (byte) 4));
        assertEquals((short) 4, v.getShort());
        assertTrue(v.getMetaData("select").getBoolean());
    }

    @Test
    void fixedTestBitPatternTwoOctetWrapsAsPlcWord() {
        assertInstanceOf(PlcWORD.class,
            Iec608705104TagParser.processFixedTestBitPatternTwoOctet(new FixedTestBitPatternTwoOctet(0x55AA)));
    }

    @Test
    void commandHelpersReturnNullUntilImplemented() {
        // The 60870-5-104 command path is not yet implemented in the
        // driver — these helpers are stubbed and intentionally return null.
        // This test pins the contract so the gap is visible in coverage.
        assertNull(Iec608705104TagParser.processSingleCommand(
            new SingleCommand(new QualifierOfCommand(false, (byte) 0), false)));
        assertNull(Iec608705104TagParser.processDoubleCommand(
            new DoubleCommand(new QualifierOfCommand(false, (byte) 0), (byte) 0)));
        assertNull(Iec608705104TagParser
            .processRegulatingStepCommand(new org.apache.plc4x.java.iec608705104.readwrite.RegulatingStepCommand(
                new QualifierOfCommand(false, (byte) 0), (byte) 0)));
        assertNull(Iec608705104TagParser.processSevenOctetBinaryTime(null));
        assertNull(Iec608705104TagParser.processThreeOctetBinaryTime(null));
        assertNull(Iec608705104TagParser.processTwoOctetBinaryTime(null));
        assertNull(Iec608705104TagParser.processBinaryStateInformation(null));
        assertNull(Iec608705104TagParser.processStatusChangeDetection(null));
    }

    @Test
    void parseTagDispatchesByTypeIdentification() {
        // Cover a representative slice of the typeSwitch — single-point,
        // step position (uses VTI + quality), normalised + scaled measured
        // values, short floating point, integrated totals, end-of-init,
        // interrogation, test command. Going through every one of the ~70
        // typeIdentifications would just bulk-replicate this shape.

        InformationObject obj;

        obj = new InformationObjectWithoutTime_SINGLE_POINT_INFORMATION(
            12, new SinglePointInformation(false, false, false, false, true));
        assertTrue(Iec608705104TagParser.parseTag(obj, TypeIdentification.SINGLE_POINT_INFORMATION)
            .getBoolean());

        obj = new InformationObjectWithoutTime_DOUBLE_POINT_INFORMATION(
            12, new DoublePointInformation(false, false, false, false, (byte) 1));
        assertInstanceOf(PlcList.class,
            Iec608705104TagParser.parseTag(obj, TypeIdentification.DOUBLE_POINT_INFORMATION));

        obj = new InformationObjectWithoutTime_STEP_POSITION_INFORMATION(12,
            new ValueWithTransientStateIndication(false, (byte) 5), goodQuality());
        assertEquals((short) 5,
            Iec608705104TagParser.parseTag(obj, TypeIdentification.STEP_POSITION_INFORMATION).getShort());

        obj = new InformationObjectWithoutTime_MEASURED_VALUE_NORMALISED_VALUE(
            12, new NormalizedValue(8), goodQuality());
        assertInstanceOf(PlcUINT.class, Iec608705104TagParser
            .parseTag(obj, TypeIdentification.MEASURED_VALUE_NORMALISED_VALUE));

        obj = new InformationObjectWithoutTime_MEASURED_VALUE_SCALED_VALUE(
            12, new ScaledValue((short) 3), goodQuality());
        assertInstanceOf(PlcUINT.class, Iec608705104TagParser
            .parseTag(obj, TypeIdentification.MEASURED_VALUE_SCALED_VALUE));

        obj = new InformationObjectWithoutTime_MEASURED_VALUE_SHORT_FLOATING_POINT_NUMBER(
            12, 1.5f, goodQuality());
        PlcValue floatVal = Iec608705104TagParser
            .parseTag(obj, TypeIdentification.MEASURED_VALUE_SHORT_FLOATING_POINT_NUMBER);
        assertInstanceOf(PlcREAL.class, floatVal);
        assertEquals(1.5f, floatVal.getFloat());

        obj = new InformationObjectWithoutTime_INTEGRATED_TOTALS(
            12, new BinaryCounterReading(100L, true, false, false, (byte) 0));
        assertInstanceOf(PlcUDINT.class,
            Iec608705104TagParser.parseTag(obj, TypeIdentification.INTEGRATED_TOTALS));

        obj = new InformationObjectWithoutTime_END_OF_INITIALISATION(
            12, new CauseOfInitialization(false, (byte) 1));
        assertEquals((short) 1,
            Iec608705104TagParser.parseTag(obj, TypeIdentification.END_OF_INITIALISATION).getShort());

        obj = new InformationObjectWithoutTime_INTERROGATION_COMMAND(
            12, new QualifierOfInterrogation((short) 20));
        assertEquals(20,
            Iec608705104TagParser.parseTag(obj, TypeIdentification.INTERROGATION_COMMAND).getInt());

        obj = new InformationObjectWithoutTime_TEST_COMMAND(
            12, new FixedTestBitPatternTwoOctet(0xAA55));
        assertInstanceOf(PlcWORD.class,
            Iec608705104TagParser.parseTag(obj, TypeIdentification.TEST_COMMAND));

        // Command path is still stubbed out; the case arm executes but
        // returns null from the helper.
        obj = new InformationObjectWithoutTime_DOUBLE_COMMAND(
            12, new DoubleCommand(new QualifierOfCommand(false, (byte) 0), (byte) 1));
        assertNull(Iec608705104TagParser.parseTag(obj, TypeIdentification.DOUBLE_COMMAND));
    }

    @Test
    void parseTagCoversRemainingWithoutTimeCases() {
        // Walks the rest of the typeSwitch over no-time information objects.
        // Many of these are stubbed (return null after the cast) — running
        // through them pins the case dispatch and protects against silent
        // reordering when new types are added.
        QualityDescriptor q = goodQuality();
        QualifierOfSetPointCommand qos = new QualifierOfSetPointCommand(false, (byte) 0);
        QualifierOfParameterOfMeasuredValues qpm =
            new QualifierOfParameterOfMeasuredValues(false, false, (byte) 0);
        org.apache.plc4x.java.iec608705104.readwrite.BinaryStateInformation bsi =
            new org.apache.plc4x.java.iec608705104.readwrite.BinaryStateInformation(0L);
        org.apache.plc4x.java.iec608705104.readwrite.StatusChangeDetection scd =
            new org.apache.plc4x.java.iec608705104.readwrite.StatusChangeDetection(0L);

        // BITSTRING_OF_32_BIT and PACKED_SINGLE_POINT_INFORMATION_WITH_*:
        // processBinaryStateInformation / processStatusChangeDetection are
        // still stubbed to null but the case arm then calls
        // processQualityDescriptor on the null — so it blows up in NPE.
        // Exercising the arm anyway covers the cast + dispatch + helper call.
        assertThrows(NullPointerException.class, () -> Iec608705104TagParser.parseTag(
            new org.apache.plc4x.java.iec608705104.readwrite.InformationObjectWithoutTime_BITSTRING_OF_32_BIT(
                0, bsi, q),
            TypeIdentification.BITSTRING_OF_32_BIT));
        assertThrows(NullPointerException.class, () -> Iec608705104TagParser.parseTag(
            new org.apache.plc4x.java.iec608705104.readwrite.InformationObjectWithoutTime_PACKED_SINGLE_POINT_INFORMATION_WITH_STATUS_CHANGE_DETECTION(
                0, scd, q),
            TypeIdentification.PACKED_SINGLE_POINT_INFORMATION_WITH_STATUS_CHANGE_DETECTION));
        Iec608705104TagParser.parseTag(
            new org.apache.plc4x.java.iec608705104.readwrite.InformationObjectWithoutTime_MEASURED_VALUE_NORMALIZED_VALUE_WITHOUT_QUALITY_DESCRIPTOR(
                0, new NormalizedValue(1)),
            TypeIdentification.MEASURED_VALUE_NORMALIZED_VALUE_WITHOUT_QUALITY_DESCRIPTOR);
        Iec608705104TagParser.parseTag(
            new org.apache.plc4x.java.iec608705104.readwrite.InformationObjectWithoutTime_SINGLE_COMMAND(
                0, new SingleCommand(new QualifierOfCommand(false, (byte) 0), false)),
            TypeIdentification.SINGLE_COMMAND);
        Iec608705104TagParser.parseTag(
            new org.apache.plc4x.java.iec608705104.readwrite.InformationObjectWithoutTime_REGULATING_STEP_COMMAND(
                0, new org.apache.plc4x.java.iec608705104.readwrite.RegulatingStepCommand(
                    new QualifierOfCommand(false, (byte) 0), (byte) 0)),
            TypeIdentification.REGULATING_STEP_COMMAND);
        Iec608705104TagParser.parseTag(
            new org.apache.plc4x.java.iec608705104.readwrite.InformationObjectWithoutTime_SET_POINT_COMMAND_NORMALISED_VALUE(
                0, new NormalizedValue(1), qos),
            TypeIdentification.SET_POINT_COMMAND_NORMALISED_VALUE);
        Iec608705104TagParser.parseTag(
            new org.apache.plc4x.java.iec608705104.readwrite.InformationObjectWithoutTime_SET_POINT_COMMAND_SCALED_VALUE(
                0, new ScaledValue((short) 1), qos),
            TypeIdentification.SET_POINT_COMMAND_SCALED_VALUE);
        Iec608705104TagParser.parseTag(
            new org.apache.plc4x.java.iec608705104.readwrite.InformationObjectWithoutTime_SET_POINT_COMMAND_SHORT_FLOATING_POINT_NUMBER(
                0, 1.0f, qos),
            TypeIdentification.SET_POINT_COMMAND_SHORT_FLOATING_POINT_NUMBER);
        Iec608705104TagParser.parseTag(
            new org.apache.plc4x.java.iec608705104.readwrite.InformationObjectWithoutTime_BITSTRING_32_BIT_COMMAND(
                0, bsi),
            TypeIdentification.BITSTRING_32_BIT_COMMAND);
        Iec608705104TagParser.parseTag(
            new org.apache.plc4x.java.iec608705104.readwrite.InformationObjectWithoutTime_COUNTER_INTERROGATION_COMMAND(
                0, new QualifierOfCounterInterrogationCommand((byte) 0, (byte) 0)),
            TypeIdentification.COUNTER_INTERROGATION_COMMAND);
        // READ_COMMAND, FILE_READY, SECTION_READY, and similar stubs only
        // execute the case arm and `break` to the throw — wrap them in
        // assertThrows since they fall through to NotImplementedException.
        assertThrows(RuntimeException.class, () -> Iec608705104TagParser.parseTag(
            new org.apache.plc4x.java.iec608705104.readwrite.InformationObjectWithoutTime_READ_COMMAND(0),
            TypeIdentification.READ_COMMAND));
        Iec608705104TagParser.parseTag(
            new org.apache.plc4x.java.iec608705104.readwrite.InformationObjectWithoutTime_CLOCK_SYNCHRONISATION_COMMAND(
                0, null),
            TypeIdentification.CLOCK_SYNCHRONISATION_COMMAND);
        Iec608705104TagParser.parseTag(
            new org.apache.plc4x.java.iec608705104.readwrite.InformationObjectWithoutTime_RESET_PROCESS_COMMAND(
                0, new QualifierOfResetProcessCommand((short) 0)),
            TypeIdentification.RESET_PROCESS_COMMAND);
        Iec608705104TagParser.parseTag(
            new org.apache.plc4x.java.iec608705104.readwrite.InformationObjectWithoutTime_DELAY_ACQUISITION_COMMAND(
                0, new org.apache.plc4x.java.iec608705104.readwrite.TwoOctetBinaryTime(0)),
            TypeIdentification.DELAY_ACQUISITION_COMMAND);
        Iec608705104TagParser.parseTag(
            new org.apache.plc4x.java.iec608705104.readwrite.InformationObjectWithoutTime_PARAMETER_OF_MEASURED_VALUES_NORMALIZED_VALUE(
                0, new NormalizedValue(1), qpm),
            TypeIdentification.PARAMETER_OF_MEASURED_VALUES_NORMALIZED_VALUE);
        Iec608705104TagParser.parseTag(
            new org.apache.plc4x.java.iec608705104.readwrite.InformationObjectWithoutTime_PARAMETER_OF_MEASURED_VALUES_SCALED_VALUE(
                0, new ScaledValue((short) 1), qpm),
            TypeIdentification.PARAMETER_OF_MEASURED_VALUES_SCALED_VALUE);
        Iec608705104TagParser.parseTag(
            new org.apache.plc4x.java.iec608705104.readwrite.InformationObjectWithoutTime_PARAMETER_OF_MEASURED_VALUES_SHORT_FLOATING_POINT_NUMBER(
                0, 1.0f, qpm),
            TypeIdentification.PARAMETER_OF_MEASURED_VALUES_SHORT_FLOATING_POINT_NUMBER);
        Iec608705104TagParser.parseTag(
            new org.apache.plc4x.java.iec608705104.readwrite.InformationObjectWithoutTime_PARAMETER_ACTIVATION(
                0, new QualifierOfParameterActivation((short) 0)),
            TypeIdentification.PARAMETER_ACTIVATION);
    }

    @Test
    void parseTagCoversTimeTaggedCases() {
        QualityDescriptor q = goodQuality();
        org.apache.plc4x.java.iec608705104.readwrite.ThreeOctetBinaryTime cp24 =
            new org.apache.plc4x.java.iec608705104.readwrite.ThreeOctetBinaryTime(0, false, (byte) 0);
        org.apache.plc4x.java.iec608705104.readwrite.SevenOctetBinaryTime cp56 =
            new org.apache.plc4x.java.iec608705104.readwrite.SevenOctetBinaryTime(
                0, false, false, (byte) 0, false, (byte) 0, (byte) 0, (byte) 0, (byte) 0, (byte) 0);
        org.apache.plc4x.java.iec608705104.readwrite.BinaryStateInformation bsi =
            new org.apache.plc4x.java.iec608705104.readwrite.BinaryStateInformation(0L);

        // Three-octet time variants (CP24Time2a).
        Iec608705104TagParser.parseTag(
            new org.apache.plc4x.java.iec608705104.readwrite.InformationObjectWithTreeByteTime_MEASURED_VALUE_SCALED_VALUE(
                0, new ScaledValue((short) 1), q, cp24),
            TypeIdentification.MEASURED_VALUE_SCALED_VALUE_WITH_TIME_TAG);
        Iec608705104TagParser.parseTag(
            new org.apache.plc4x.java.iec608705104.readwrite.InformationObjectWithTreeByteTime_SINGLE_POINT_INFORMATION(
                0, new SinglePointInformation(false, false, false, false, true), cp24),
            TypeIdentification.SINGLE_POINT_INFORMATION_WITH_TIME_TAG);
        Iec608705104TagParser.parseTag(
            new org.apache.plc4x.java.iec608705104.readwrite.InformationObjectWithTreeByteTime_DOUBLE_POINT_INFORMATION(
                0, new DoublePointInformation(false, false, false, false, (byte) 1), cp24),
            TypeIdentification.DOUBLE_POINT_INFORMATION_WITH_TIME_TAG);
        Iec608705104TagParser.parseTag(
            new org.apache.plc4x.java.iec608705104.readwrite.InformationObjectWithTreeByteTime_STEP_POSITION_INFORMATION(
                0, new ValueWithTransientStateIndication(false, (byte) 1), q, cp24),
            TypeIdentification.STEP_POSITION_INFORMATION_WITH_TIME_TAG);
        assertThrows(NullPointerException.class, () -> Iec608705104TagParser.parseTag(
            new org.apache.plc4x.java.iec608705104.readwrite.InformationObjectWithTreeByteTime_BITSTRING_OF_32_BIT(
                0, bsi, q, cp24),
            TypeIdentification.BITSTRING_OF_32_BIT_WITH_TIME_TAG));
        Iec608705104TagParser.parseTag(
            new org.apache.plc4x.java.iec608705104.readwrite.InformationObjectWithTreeByteTime_MEASURED_VALUE_NORMALIZED_VALUE(
                0, new NormalizedValue(1), q, cp24),
            TypeIdentification.MEASURED_VALUE_NORMALIZED_VALUE_WITH_TIME_TAG);
        Iec608705104TagParser.parseTag(
            new org.apache.plc4x.java.iec608705104.readwrite.InformationObjectWithTreeByteTime_MEASURED_VALUE_SHORT_FLOATING_POINT_NUMBER(
                0, 1.0f, q, cp24),
            TypeIdentification.MEASURED_VALUE_SHORT_FLOATING_POINT_NUMBER_WITH_TIME_TAG);
        Iec608705104TagParser.parseTag(
            new org.apache.plc4x.java.iec608705104.readwrite.InformationObjectWithTreeByteTime_INTEGRATED_TOTALS(
                0, new BinaryCounterReading(0L, false, false, false, (byte) 0), cp24),
            TypeIdentification.INTEGRATED_TOTALS_WITH_TIME_TAG);

        // Seven-octet time variants (CP56Time2a).
        Iec608705104TagParser.parseTag(
            new org.apache.plc4x.java.iec608705104.readwrite.InformationObjectWithSevenByteTime_SINGLE_POINT_INFORMATION(
                0, new SinglePointInformation(false, false, false, false, true), cp56),
            TypeIdentification.SINGLE_POINT_INFORMATION_WITH_TIME_TAG_CP56TIME2A);
        Iec608705104TagParser.parseTag(
            new org.apache.plc4x.java.iec608705104.readwrite.InformationObjectWithSevenByteTime_DOUBLE_POINT_INFORMATION(
                0, new DoublePointInformation(false, false, false, false, (byte) 1), cp56),
            TypeIdentification.DOUBLE_POINT_INFORMATION_WITH_TIME_TAG_CP56TIME2A);
        Iec608705104TagParser.parseTag(
            new org.apache.plc4x.java.iec608705104.readwrite.InformationObjectWithSevenByteTime_STEP_POSITION_INFORMATION(
                0, new ValueWithTransientStateIndication(false, (byte) 1), q, cp56),
            TypeIdentification.STEP_POSITION_INFORMATION_WITH_TIME_TAG_CP56TIME2A);
        assertThrows(NullPointerException.class, () -> Iec608705104TagParser.parseTag(
            new org.apache.plc4x.java.iec608705104.readwrite.InformationObjectWithSevenByteTime_BITSTRING_OF_32_BIT(
                0, bsi, q, cp56),
            TypeIdentification.BITSTRING_OF_32_BIT_WITH_TIME_TAG_CP56TIME2A));
        Iec608705104TagParser.parseTag(
            new org.apache.plc4x.java.iec608705104.readwrite.InformationObjectWithSevenByteTime_MEASURED_VALUE_NORMALISED_VALUE(
                0, new NormalizedValue(1), q, cp56),
            TypeIdentification.MEASURED_VALUE_NORMALISED_VALUE_WITH_TIME_TAG_CP56TIME2A);
        Iec608705104TagParser.parseTag(
            new org.apache.plc4x.java.iec608705104.readwrite.InformationObjectWithSevenByteTime_MEASURED_VALUE_SCALED_VALUE(
                0, new ScaledValue((short) 1), q, cp56),
            TypeIdentification.MEASURED_VALUE_SCALED_VALUE_WITH_TIME_TAG_CP56TIME2A);
        Iec608705104TagParser.parseTag(
            new org.apache.plc4x.java.iec608705104.readwrite.InformationObjectWithSevenByteTime_MEASURED_VALUE_SHORT_FLOATING_POINT_NUMBER(
                0, 1.0f, q, cp56),
            TypeIdentification.MEASURED_VALUE_SHORT_FLOATING_POINT_NUMBER_WITH_TIME_TAG_CP56TIME2A);
        Iec608705104TagParser.parseTag(
            new org.apache.plc4x.java.iec608705104.readwrite.InformationObjectWithSevenByteTime_INTEGRATED_TOTALS(
                0, new BinaryCounterReading(0L, false, false, false, (byte) 0), cp56),
            TypeIdentification.INTEGRATED_TOTALS_WITH_TIME_TAG_CP56TIME2A);
    }

    @Test
    void parseTagThrowsForUnimplementedTypeIdentification() {
        // NOT_USED falls through to the default 'throw' — the switch has
        // no return in that arm.
        assertThrows(RuntimeException.class,
            () -> Iec608705104TagParser.parseTag(
                new InformationObjectWithoutTime_SINGLE_POINT_INFORMATION(
                    0, new SinglePointInformation(false, false, false, false, false)),
                TypeIdentification.NOT_USED));
    }

}
