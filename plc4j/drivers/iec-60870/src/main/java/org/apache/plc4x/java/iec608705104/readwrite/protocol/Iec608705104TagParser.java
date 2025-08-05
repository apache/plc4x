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

package org.apache.plc4x.java.iec608705104.readwrite.protocol;

import org.apache.commons.lang3.NotImplementedException;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.iec608705104.readwrite.*;
import org.apache.plc4x.java.spi.values.*;

import java.util.Arrays;
import java.util.HashMap;
import java.util.Map;

public class Iec608705104TagParser {

    public static PlcValue parseTag(InformationObject informationObject, TypeIdentification typeIdentification) {
        switch (typeIdentification) {
            case NOT_USED: {
                break;
            }
            case SINGLE_POINT_INFORMATION: {
                InformationObjectWithoutTime_SINGLE_POINT_INFORMATION castedObject = (InformationObjectWithoutTime_SINGLE_POINT_INFORMATION) informationObject;
                return processSinglePointInformation(castedObject.getSiq());
            }
            case DOUBLE_POINT_INFORMATION: {
                InformationObjectWithoutTime_DOUBLE_POINT_INFORMATION castedObject = (InformationObjectWithoutTime_DOUBLE_POINT_INFORMATION) informationObject;
                return processDoublePointInformation(castedObject.getDiq());
            }
            case STEP_POSITION_INFORMATION: {
                InformationObjectWithoutTime_STEP_POSITION_INFORMATION castedObject = (InformationObjectWithoutTime_STEP_POSITION_INFORMATION) informationObject;
                PlcValueAdapter plcValue = processValueWithTransientStateIndication(castedObject.getVti());
                processQualityDescriptor(castedObject.getQds(), plcValue);
                return plcValue;
            }
            case BITSTRING_OF_32_BIT: {
                InformationObjectWithoutTime_BITSTRING_OF_32_BIT castedObject = (InformationObjectWithoutTime_BITSTRING_OF_32_BIT) informationObject;
                PlcValueAdapter plcValue = processBinaryStateInformation(castedObject.getBsi());
                processQualityDescriptor(castedObject.getQds(), plcValue);
                return plcValue;
            }
            case MEASURED_VALUE_NORMALISED_VALUE: {
                InformationObjectWithoutTime_MEASURED_VALUE_NORMALISED_VALUE castedObject = (InformationObjectWithoutTime_MEASURED_VALUE_NORMALISED_VALUE) informationObject;
                PlcValueAdapter plcValue = processNormalizedValue(castedObject.getNva());
                processQualityDescriptor(castedObject.getQds(), plcValue);
                return plcValue;
            }
            case MEASURED_VALUE_SCALED_VALUE: {
                InformationObjectWithoutTime_MEASURED_VALUE_SCALED_VALUE castedObject = (InformationObjectWithoutTime_MEASURED_VALUE_SCALED_VALUE) informationObject;
                PlcValueAdapter plcValue = processScaledValue(castedObject.getSva());
                processQualityDescriptor(castedObject.getQds(), plcValue);
                return plcValue;
            }
            case MEASURED_VALUE_SHORT_FLOATING_POINT_NUMBER: {
                InformationObjectWithoutTime_MEASURED_VALUE_SHORT_FLOATING_POINT_NUMBER castedObject = (InformationObjectWithoutTime_MEASURED_VALUE_SHORT_FLOATING_POINT_NUMBER) informationObject;
                PlcValueAdapter plcValue = PlcREAL.of(castedObject.getValue());
                processQualityDescriptor(castedObject.getQds(), plcValue);
                return plcValue;
            }
            case INTEGRATED_TOTALS: {
                InformationObjectWithoutTime_INTEGRATED_TOTALS castedObject = (InformationObjectWithoutTime_INTEGRATED_TOTALS) informationObject;
                return processBinaryCounterReading(castedObject.getBcr());
            }
            case PACKED_SINGLE_POINT_INFORMATION_WITH_STATUS_CHANGE_DETECTION: {
                InformationObjectWithoutTime_PACKED_SINGLE_POINT_INFORMATION_WITH_STATUS_CHANGE_DETECTION castedObject = (InformationObjectWithoutTime_PACKED_SINGLE_POINT_INFORMATION_WITH_STATUS_CHANGE_DETECTION) informationObject;
                PlcValueAdapter plcValue = processStatusChangeDetection(castedObject.getScd());
                processQualityDescriptor(castedObject.getQds(), plcValue);
                return plcValue;
            }
            case MEASURED_VALUE_NORMALIZED_VALUE_WITHOUT_QUALITY_DESCRIPTOR: {
                InformationObjectWithoutTime_MEASURED_VALUE_NORMALIZED_VALUE_WITHOUT_QUALITY_DESCRIPTOR castedObject = (InformationObjectWithoutTime_MEASURED_VALUE_NORMALIZED_VALUE_WITHOUT_QUALITY_DESCRIPTOR) informationObject;
                return processNormalizedValue(castedObject.getNva());
            }
            case SINGLE_COMMAND: {
                InformationObjectWithoutTime_SINGLE_COMMAND castedObject = (InformationObjectWithoutTime_SINGLE_COMMAND) informationObject;
                return processSingleCommand(castedObject.getSco());
            }
            case DOUBLE_COMMAND: {
                InformationObjectWithoutTime_DOUBLE_COMMAND castedObject = (InformationObjectWithoutTime_DOUBLE_COMMAND) informationObject;
                return processDoubleCommand(castedObject.getDco());
            }
            case REGULATING_STEP_COMMAND: {
                InformationObjectWithoutTime_REGULATING_STEP_COMMAND castedObject = (InformationObjectWithoutTime_REGULATING_STEP_COMMAND) informationObject;
                return processRegulatingStepCommand(castedObject.getRco());
            }
            case SET_POINT_COMMAND_NORMALISED_VALUE: {
                InformationObjectWithoutTime_SET_POINT_COMMAND_NORMALISED_VALUE castedObject = (InformationObjectWithoutTime_SET_POINT_COMMAND_NORMALISED_VALUE) informationObject;
                PlcValueAdapter plcValue = processNormalizedValue(castedObject.getNva());
                processQualifierOfSetPointCommand(castedObject.getQos(), plcValue);
                return plcValue;
            }
            case SET_POINT_COMMAND_SCALED_VALUE: {
                InformationObjectWithoutTime_SET_POINT_COMMAND_SCALED_VALUE castedObject = (InformationObjectWithoutTime_SET_POINT_COMMAND_SCALED_VALUE) informationObject;
                PlcValueAdapter plcValue = processScaledValue(castedObject.getSva());
                processQualifierOfSetPointCommand(castedObject.getQos(), plcValue);
                return plcValue;
            }
            case SET_POINT_COMMAND_SHORT_FLOATING_POINT_NUMBER: {
                InformationObjectWithoutTime_SET_POINT_COMMAND_SHORT_FLOATING_POINT_NUMBER castedObject = (InformationObjectWithoutTime_SET_POINT_COMMAND_SHORT_FLOATING_POINT_NUMBER) informationObject;
                PlcValueAdapter plcValue = PlcREAL.of(castedObject.getValue());
                processQualifierOfSetPointCommand(castedObject.getQos(), plcValue);
                return plcValue;
            }
            case BITSTRING_32_BIT_COMMAND: {
                InformationObjectWithoutTime_BITSTRING_32_BIT_COMMAND castedObject = (InformationObjectWithoutTime_BITSTRING_32_BIT_COMMAND) informationObject;
                return processBinaryStateInformation(castedObject.getBsi());
            }
            case END_OF_INITIALISATION: {
                InformationObjectWithoutTime_END_OF_INITIALISATION castedObject = (InformationObjectWithoutTime_END_OF_INITIALISATION) informationObject;
                return processCauseOfInitialization(castedObject.getCoi());
            }
            case INTERROGATION_COMMAND: {
                InformationObjectWithoutTime_INTERROGATION_COMMAND castedObject = (InformationObjectWithoutTime_INTERROGATION_COMMAND) informationObject;
                return processQualifierOfInterrogation(castedObject.getQoi());
            }
            case COUNTER_INTERROGATION_COMMAND: {
                InformationObjectWithoutTime_COUNTER_INTERROGATION_COMMAND castedObject = (InformationObjectWithoutTime_COUNTER_INTERROGATION_COMMAND) informationObject;
                return processQualifierOfCounterInterrogationCommand(castedObject.getQcc());
            }
            case READ_COMMAND: {
                InformationObjectWithoutTime_READ_COMMAND castedObject = (InformationObjectWithoutTime_READ_COMMAND) informationObject;
                // TODO: Implement
                break;
            }
            case CLOCK_SYNCHRONISATION_COMMAND: {
                InformationObjectWithoutTime_CLOCK_SYNCHRONISATION_COMMAND castedObject = (InformationObjectWithoutTime_CLOCK_SYNCHRONISATION_COMMAND) informationObject;
                return processSevenOctetBinaryTime(castedObject.getCp56Time2a());
            }
            case TEST_COMMAND: {
                InformationObjectWithoutTime_TEST_COMMAND castedObject = (InformationObjectWithoutTime_TEST_COMMAND) informationObject;
                return processFixedTestBitPatternTwoOctet(castedObject.getFbp());
            }
            case RESET_PROCESS_COMMAND: {
                InformationObjectWithoutTime_RESET_PROCESS_COMMAND castedObject = (InformationObjectWithoutTime_RESET_PROCESS_COMMAND) informationObject;
                return processQualifierOfResetProcessCommand(castedObject.getQrp());
            }
            case DELAY_ACQUISITION_COMMAND: {
                InformationObjectWithoutTime_DELAY_ACQUISITION_COMMAND castedObject = (InformationObjectWithoutTime_DELAY_ACQUISITION_COMMAND) informationObject;
                return processTwoOctetBinaryTime(castedObject.getCp16Time2a());
            }
            case PARAMETER_OF_MEASURED_VALUES_NORMALIZED_VALUE: {
                InformationObjectWithoutTime_PARAMETER_OF_MEASURED_VALUES_NORMALIZED_VALUE castedObject = (InformationObjectWithoutTime_PARAMETER_OF_MEASURED_VALUES_NORMALIZED_VALUE) informationObject;
                PlcValueAdapter plcValue = processNormalizedValue(castedObject.getNva());
                processQualifierOfParameterOfMeasuredValues(castedObject.getQpm(), plcValue);
                return plcValue;
            }
            case PARAMETER_OF_MEASURED_VALUES_SCALED_VALUE: {
                InformationObjectWithoutTime_PARAMETER_OF_MEASURED_VALUES_SCALED_VALUE castedObject = (InformationObjectWithoutTime_PARAMETER_OF_MEASURED_VALUES_SCALED_VALUE) informationObject;
                PlcValueAdapter plcValue = processScaledValue(castedObject.getSva());
                processQualifierOfParameterOfMeasuredValues(castedObject.getQpm(), plcValue);
                return plcValue;
            }
            case PARAMETER_OF_MEASURED_VALUES_SHORT_FLOATING_POINT_NUMBER: {
                InformationObjectWithoutTime_PARAMETER_OF_MEASURED_VALUES_SHORT_FLOATING_POINT_NUMBER castedObject = (InformationObjectWithoutTime_PARAMETER_OF_MEASURED_VALUES_SHORT_FLOATING_POINT_NUMBER) informationObject;
                PlcValueAdapter plcValue = PlcREAL.of(castedObject.getValue());
                processQualifierOfParameterOfMeasuredValues(castedObject.getQpm(), plcValue);
                return plcValue;
            }
            case PARAMETER_ACTIVATION: {
                InformationObjectWithoutTime_PARAMETER_ACTIVATION castedObject = (InformationObjectWithoutTime_PARAMETER_ACTIVATION) informationObject;
                return processQualifierOfParameterActivation(castedObject.getQpa());
            }
            case FILE_READY: {
                InformationObjectWithoutTime_FILE_READY castedObject = (InformationObjectWithoutTime_FILE_READY) informationObject;
                // TODO: Implement
                break;
            }
            case SECTION_READY: {
                InformationObjectWithoutTime_SECTION_READY castedObject = (InformationObjectWithoutTime_SECTION_READY) informationObject;
                // TODO: Implement
                break;
            }
            case CALL_DIRECTORY_SELECT_FILE_CALL_FILE_CALL_SECTION: {
                InformationObjectWithoutTime_CALL_DIRECTORY_SELECT_FILE_CALL_FILE_CALL_SECTION castedObject = (InformationObjectWithoutTime_CALL_DIRECTORY_SELECT_FILE_CALL_FILE_CALL_SECTION) informationObject;
                // TODO: Implement
                break;
            }
            case LAST_SECTION_LAST_SEGMENT: {
                InformationObjectWithoutTime_LAST_SECTION_LAST_SEGMENT castedObject = (InformationObjectWithoutTime_LAST_SECTION_LAST_SEGMENT) informationObject;
                // TODO: Implement
                break;
            }
            case ACK_FILE_ACK_SECTION: {
                InformationObjectWithoutTime_ACK_FILE_ACK_SECTION castedObject = (InformationObjectWithoutTime_ACK_FILE_ACK_SECTION) informationObject;
                // TODO: Implement
                break;
            }
            case SEGMENT: {
                InformationObjectWithoutTime_SEGMENT castedObject = (InformationObjectWithoutTime_SEGMENT) informationObject;
                // TODO: Implement
                break;
            }
            case DIRECTORY: {
                InformationObjectWithoutTime_DIRECTORY castedObject = (InformationObjectWithoutTime_DIRECTORY) informationObject;
                // TODO: Implement
                break;
            }




            case MEASURED_VALUE_SCALED_VALUE_WITH_TIME_TAG: {
                InformationObjectWithTreeByteTime_MEASURED_VALUE_SCALED_VALUE castedObject = (InformationObjectWithTreeByteTime_MEASURED_VALUE_SCALED_VALUE) informationObject;
                PlcValueAdapter plcValue = processScaledValue(castedObject.getSva());
                processQualityDescriptor(castedObject.getQds(), plcValue);
                return plcValue;
            }
            case SINGLE_POINT_INFORMATION_WITH_TIME_TAG: {
                InformationObjectWithTreeByteTime_SINGLE_POINT_INFORMATION castedObject = (InformationObjectWithTreeByteTime_SINGLE_POINT_INFORMATION) informationObject;
                return processSinglePointInformation(castedObject.getSiq());
            }
            case DOUBLE_POINT_INFORMATION_WITH_TIME_TAG: {
                InformationObjectWithTreeByteTime_DOUBLE_POINT_INFORMATION castedObject = (InformationObjectWithTreeByteTime_DOUBLE_POINT_INFORMATION) informationObject;
                return processDoublePointInformation(castedObject.getDiq());
            }
            case STEP_POSITION_INFORMATION_WITH_TIME_TAG: {
                InformationObjectWithTreeByteTime_STEP_POSITION_INFORMATION castedObject = (InformationObjectWithTreeByteTime_STEP_POSITION_INFORMATION) informationObject;
                PlcValueAdapter plcValue = processValueWithTransientStateIndication(castedObject.getVti());
                processQualityDescriptor(castedObject.getQds(), plcValue);
                return plcValue;
            }
            case BITSTRING_OF_32_BIT_WITH_TIME_TAG: {
                InformationObjectWithTreeByteTime_BITSTRING_OF_32_BIT castedObject = (InformationObjectWithTreeByteTime_BITSTRING_OF_32_BIT) informationObject;
                PlcValueAdapter plcValue = processBinaryStateInformation(castedObject.getBsi());
                processQualityDescriptor(castedObject.getQds(), plcValue);
                return plcValue;
            }
            case MEASURED_VALUE_NORMALIZED_VALUE_WITH_TIME_TAG: {
                InformationObjectWithTreeByteTime_MEASURED_VALUE_NORMALIZED_VALUE castedObject = (InformationObjectWithTreeByteTime_MEASURED_VALUE_NORMALIZED_VALUE) informationObject;
                PlcValueAdapter plcValue = processNormalizedValue(castedObject.getNva());
                processQualityDescriptor(castedObject.getQds(), plcValue);
                return plcValue;
            }
            case MEASURED_VALUE_SHORT_FLOATING_POINT_NUMBER_WITH_TIME_TAG: {
                InformationObjectWithTreeByteTime_MEASURED_VALUE_SHORT_FLOATING_POINT_NUMBER castedObject = (InformationObjectWithTreeByteTime_MEASURED_VALUE_SHORT_FLOATING_POINT_NUMBER) informationObject;
                PlcValueAdapter plcValue = PlcREAL.of(castedObject.getValue());
                processQualityDescriptor(castedObject.getQds(), plcValue);
                return plcValue;
            }
            case INTEGRATED_TOTALS_WITH_TIME_TAG: {
                InformationObjectWithTreeByteTime_INTEGRATED_TOTALS castedObject = (InformationObjectWithTreeByteTime_INTEGRATED_TOTALS) informationObject;
                return processBinaryCounterReading(castedObject.getBcr());
            }
            case EVENT_OF_PROTECTION_EQUIPMENT_WITH_TIME_TAG: {
                InformationObjectWithTreeByteTime_EVENT_OF_PROTECTION_EQUIPMENT castedObject = (InformationObjectWithTreeByteTime_EVENT_OF_PROTECTION_EQUIPMENT) informationObject;
                // TODO: Implement
                break;
            }
            case PACKED_START_EVENTS_OF_PROTECTION_EQUIPMENT_WITH_TIME_TAG: {
                InformationObjectWithTreeByteTime_PACKED_START_EVENTS_OF_PROTECTION_EQUIPMENT castedObject = (InformationObjectWithTreeByteTime_PACKED_START_EVENTS_OF_PROTECTION_EQUIPMENT) informationObject;
                // TODO: Implement
                break;
            }
            case PACKED_OUTPUT_CIRCUIT_INFORMATION_OF_PROTECTION_EQUIPMENT_WITH_TIME_TAG: {
                InformationObjectWithTreeByteTime_PACKED_OUTPUT_CIRCUIT_INFORMATION_OF_PROTECTION_EQUIPMENT castedObject = (InformationObjectWithTreeByteTime_PACKED_OUTPUT_CIRCUIT_INFORMATION_OF_PROTECTION_EQUIPMENT) informationObject;
                // TODO: Implement
                break;
            }


            case SINGLE_POINT_INFORMATION_WITH_TIME_TAG_CP56TIME2A: {
                InformationObjectWithSevenByteTime_SINGLE_POINT_INFORMATION castedObject = (InformationObjectWithSevenByteTime_SINGLE_POINT_INFORMATION ) informationObject;
                return processSinglePointInformation(castedObject.getSiq());
            }
            case DOUBLE_POINT_INFORMATION_WITH_TIME_TAG_CP56TIME2A: {
                InformationObjectWithSevenByteTime_DOUBLE_POINT_INFORMATION castedObject = (InformationObjectWithSevenByteTime_DOUBLE_POINT_INFORMATION ) informationObject;
                return processDoublePointInformation(castedObject.getDiq());
            }
            case STEP_POSITION_INFORMATION_WITH_TIME_TAG_CP56TIME2A: {
                InformationObjectWithSevenByteTime_STEP_POSITION_INFORMATION castedObject = (InformationObjectWithSevenByteTime_STEP_POSITION_INFORMATION ) informationObject;
                PlcValueAdapter plcValue = processValueWithTransientStateIndication(castedObject.getVti());
                processQualityDescriptor(castedObject.getQds(), plcValue);
                return plcValue;
            }
            case BITSTRING_OF_32_BIT_WITH_TIME_TAG_CP56TIME2A: {
                InformationObjectWithSevenByteTime_BITSTRING_OF_32_BIT castedObject = (InformationObjectWithSevenByteTime_BITSTRING_OF_32_BIT ) informationObject;
                PlcValueAdapter plcValue = processBinaryStateInformation(castedObject.getBsi());
                processQualityDescriptor(castedObject.getQds(), plcValue);
                return plcValue;
            }
            case MEASURED_VALUE_NORMALISED_VALUE_WITH_TIME_TAG_CP56TIME2A: {
                InformationObjectWithSevenByteTime_MEASURED_VALUE_NORMALISED_VALUE castedObject = (InformationObjectWithSevenByteTime_MEASURED_VALUE_NORMALISED_VALUE ) informationObject;
                PlcValueAdapter plcValue = processNormalizedValue(castedObject.getNva());
                processQualityDescriptor(castedObject.getQds(), plcValue);
                return plcValue;
            }
            case MEASURED_VALUE_SCALED_VALUE_WITH_TIME_TAG_CP56TIME2A: {
                InformationObjectWithSevenByteTime_MEASURED_VALUE_SCALED_VALUE castedObject = (InformationObjectWithSevenByteTime_MEASURED_VALUE_SCALED_VALUE ) informationObject;
                PlcValueAdapter plcValue = processScaledValue(castedObject.getSva());
                processQualityDescriptor(castedObject.getQds(), plcValue);
                return plcValue;
            }
            case MEASURED_VALUE_SHORT_FLOATING_POINT_NUMBER_WITH_TIME_TAG_CP56TIME2A: {
                InformationObjectWithSevenByteTime_MEASURED_VALUE_SHORT_FLOATING_POINT_NUMBER castedObject = (InformationObjectWithSevenByteTime_MEASURED_VALUE_SHORT_FLOATING_POINT_NUMBER ) informationObject;
                PlcValueAdapter plcValue = PlcREAL.of(castedObject.getValue());
                processQualityDescriptor(castedObject.getQds(), plcValue);
                return plcValue;
            }
            case INTEGRATED_TOTALS_WITH_TIME_TAG_CP56TIME2A: {
                InformationObjectWithSevenByteTime_INTEGRATED_TOTALS castedObject = (InformationObjectWithSevenByteTime_INTEGRATED_TOTALS ) informationObject;
                return processBinaryCounterReading(castedObject.getBcr());
            }
            case EVENT_OF_PROTECTION_EQUIPMENT_WITH_TIME_TAG_CP56TIME2A: {
                InformationObjectWithSevenByteTime_EVENT_OF_PROTECTION_EQUIPMENT castedObject = (InformationObjectWithSevenByteTime_EVENT_OF_PROTECTION_EQUIPMENT ) informationObject;
                // TODO: Implement
                break;
            }
            case PACKED_START_EVENTS_OF_PROTECTION_EQUIPMENT_WITH_TIME_TAG_CP56TIME2A: {
                InformationObjectWithSevenByteTime_PACKED_START_EVENTS_OF_PROTECTION_EQUIPMENT castedObject = (InformationObjectWithSevenByteTime_PACKED_START_EVENTS_OF_PROTECTION_EQUIPMENT ) informationObject;
                // TODO: Implement
                break;
            }
            case PACKED_OUTPUT_CIRCUIT_INFORMATION_OF_PROTECTION_EQUIPMENT_WITH_TIME_TAG_CP56TIME2A: {
                InformationObjectWithSevenByteTime_PACKED_OUTPUT_CIRCUIT_INFORMATION_OF_PROTECTION_EQUIPMENT castedObject = (InformationObjectWithSevenByteTime_PACKED_OUTPUT_CIRCUIT_INFORMATION_OF_PROTECTION_EQUIPMENT ) informationObject;
                // TODO: Implement
                break;
            }
            /*case SINGLE_COMMAND_WITH_TIME_TAG_CP56TIME2A: {
                InformationObjectWithSevenByteTime_SINGLE_COMMAND castedObject = (InformationObjectWithSevenByteTime_SINGLE_COMMAND ) informationObject;
            }
            case DOUBLE_COMMAND_WITH_TIME_TAG_CP56TIME2A: {
                InformationObjectWithSevenByteTime_DOUBLE_COMMAND castedObject = (InformationObjectWithSevenByteTime_DOUBLE_COMMAND ) informationObject;
            }
            case REGULATING_STEP_COMMAND_WITH_TIME_TAG_CP56TIME2A: {
                InformationObjectWithSevenByteTime_REGULATING_STEP_COMMAND castedObject = (InformationObjectWithSevenByteTime_REGULATING_STEP_COMMAND ) informationObject;
            }
            case MEASURED_VALUE_NORMALISED_VALUE_COMMAND_WITH_TIME_TAG_CP56TIME2A: {
                InformationObjectWithSevenByteTime_MEASURED_VALUE_NORMALISED_VALUE_COMMAND castedObject = (InformationObjectWithSevenByteTime_MEASURED_VALUE_NORMALISED_VALUE_COMMAND ) informationObject;
            }
            case MEASURED_VALUE_SCALED_VALUE_COMMAND_WITH_TIME_TAG_CP56TIME2A: {
                InformationObjectWithSevenByteTime_MEASURED_VALUE_SCALED_VALUE_COMMAND castedObject = (InformationObjectWithSevenByteTime_MEASURED_VALUE_SCALED_VALUE_COMMAND ) informationObject;
            }
            case MEASURED_VALUE_SHORT_FLOATING_POINT_NUMBER_COMMAND_WITH_TIME_TAG_CP56TIME2A: {
                InformationObjectWithSevenByteTime_MEASURED_VALUE_SHORT_FLOATING_POINT_NUMBER_COMMAND castedObject = (InformationObjectWithSevenByteTime_MEASURED_VALUE_SHORT_FLOATING_POINT_NUMBER_COMMAND ) informationObject;
            }
            case BITSTRING_OF_32_BIT_COMMAND_WITH_TIME_TAG_CP56TIME2A: {
                InformationObjectWithSevenByteTime_BITSTRING_OF_32_BIT_COMMAND castedObject = (InformationObjectWithSevenByteTime_BITSTRING_OF_32_BIT_COMMAND ) informationObject;
            }
            case TEST_COMMAND_WITH_TIME_TAG_CP56TIME2A: {
                InformationObjectWithSevenByteTime_TEST_COMMAND castedObject = (InformationObjectWithSevenByteTime_TEST_COMMAND ) informationObject;
            }*/

        }
        throw new NotImplementedException("Datatype " + typeIdentification.name() +  " not implemented.");
    }

    public static PlcValueAdapter processSinglePointInformation(SinglePointInformation siq) {
        PlcBOOL plcValue = PlcBOOL.of(siq.getStausOn());
        plcValue.addMetaData("invalid", PlcBOOL.of(siq.getInvalid()));
        plcValue.addMetaData("notTopical", PlcBOOL.of(siq.getNotTopical()));
        plcValue.addMetaData("substituted", PlcBOOL.of(siq.getSubstituted()));
        plcValue.addMetaData("blocked", PlcBOOL.of(siq.getBlocked()));
        return plcValue;
    }

    public static PlcValueAdapter processDoublePointInformation(DoublePointInformation diq) {
        boolean firstPoint = (diq.getDpiCode() & 0x01) != 0;
        boolean secondPoint = (diq.getDpiCode() & 0x02) != 0;
        PlcList plcValue = new PlcList(Arrays.asList(PlcBOOL.of(firstPoint), PlcBOOL.of(secondPoint)));
        plcValue.addMetaData("invalid", PlcBOOL.of(diq.getInvalid()));
        plcValue.addMetaData("notTopical", PlcBOOL.of(diq.getNotTopical()));
        plcValue.addMetaData("substituted", PlcBOOL.of(diq.getSubstituted()));
        plcValue.addMetaData("blocked", PlcBOOL.of(diq.getBlocked()));
        return plcValue;
    }

    public static PlcValueAdapter processBinaryStateInformation(BinaryStateInformation bsi) {
        // TODO: Implement
        return null;
    }

    public static PlcValueAdapter processStatusChangeDetection(StatusChangeDetection scd) {
        // TODO: Implement
        return null;
    }

    public static void processQualityDescriptor(QualityDescriptor qds, PlcValueAdapter plcValue) {
        plcValue.addMetaData("invalid", PlcBOOL.of(qds.getInvalid()));
        plcValue.addMetaData("notTopical", PlcBOOL.of(qds.getNotTopical()));
        plcValue.addMetaData("substituted", PlcBOOL.of(qds.getSubstituted()));
        plcValue.addMetaData("blocked", PlcBOOL.of(qds.getBlocked()));
        plcValue.addMetaData("overflow", PlcBOOL.of(qds.getOverflow()));
    }

    public static PlcValueAdapter processValueWithTransientStateIndication(ValueWithTransientStateIndication vti) {
        PlcValueAdapter plcValue = PlcUSINT.of(vti.getValue());
        plcValue.addMetaData("transientState", PlcBOOL.of(vti.getTransientState()));
        return plcValue;
    }

    public static PlcValueAdapter processNormalizedValue(NormalizedValue nva) {
        return PlcUINT.of(nva.getValue());
    }

    public static PlcValueAdapter processScaledValue(ScaledValue sva) {
        return PlcUINT.of(sva.getValue());
    }

    public static PlcValueAdapter processBinaryCounterReading(BinaryCounterReading bcr) {
        PlcValueAdapter plcValue = PlcUDINT.of(bcr.getCounterValue());
        plcValue.addMetaData("counterValid", PlcBOOL.of(bcr.getCounterValid()));
        plcValue.addMetaData("counterAdjusted", PlcBOOL.of(bcr.getCounterAdjusted()));
        plcValue.addMetaData("carry", PlcBOOL.of(bcr.getCarry()));
        plcValue.addMetaData("sequenceNumber", PlcUSINT.of(bcr.getSequenceNumber()));
        return plcValue;
    }

    public static PlcValueAdapter processSingleEventOfProtectionEquipment(SingleEventOfProtectionEquipment sep) {
        PlcValueAdapter plcValue = PlcUSINT.of(sep.getEventState());
        plcValue.addMetaData("invalid", PlcBOOL.of(sep.getInvalid()));
        plcValue.addMetaData("notTopical", PlcBOOL.of(sep.getNotTopical()));
        plcValue.addMetaData("substituted", PlcBOOL.of(sep.getSubstituted()));
        plcValue.addMetaData("blocked", PlcBOOL.of(sep.getBlocked()));
        plcValue.addMetaData("elapsedTimeInvalid", PlcBOOL.of(sep.getElapsedTimeInvalid()));
        return plcValue;
    }

    public static PlcValueAdapter processOutputCircuitInformation(OutputCircuitInformation oci) {
        Map<String, PlcValue> plcValues = new HashMap<>();
        plcValues.put("stateOfOperationPhaseL3", PlcBOOL.of(oci.getStateOfOperationPhaseL3()));
        plcValues.put("stateOfOperationPhaseL2", PlcBOOL.of(oci.getStateOfOperationPhaseL2()));
        plcValues.put("stateOfOperationPhaseL1", PlcBOOL.of(oci.getStateOfOperationPhaseL1()));
        plcValues.put("generalStartOfOperation", PlcBOOL.of(oci.getGeneralStartOfOperation()));
        return new PlcStruct(plcValues);
    }

    public static void processQualityDescriptorForPointsOfProtectionEquipment(QualityDescriptorForPointsOfProtectionEquipment qdp, PlcValueAdapter plcValue) {
        plcValue.addMetaData("invalid", PlcBOOL.of(qdp.getInvalid()));
        plcValue.addMetaData("notTopical", PlcBOOL.of(qdp.getNotTopical()));
        plcValue.addMetaData("substituted", PlcBOOL.of(qdp.getSubstituted()));
        plcValue.addMetaData("blocked", PlcBOOL.of(qdp.getBlocked()));
        plcValue.addMetaData("elapsedTimeInvalid", PlcBOOL.of(qdp.getElapsedTimeInvalid()));
    }

    public static PlcValueAdapter processSingleCommand(SingleCommand sco) {
        return null;
    }

    public static PlcValueAdapter processDoubleCommand(DoubleCommand dco) {
        return null;
    }

    public static PlcValueAdapter processRegulatingStepCommand(RegulatingStepCommand rco) {
        return null;
    }

    public static PlcValueAdapter processSevenOctetBinaryTime(SevenOctetBinaryTime cp56Time2a) {
        // TODO: Implement ...
        return null;
    }

    public static PlcValueAdapter processThreeOctetBinaryTime(ThreeOctetBinaryTime cp24Time2a) {
        // TODO: Implement ...
        return null;
    }

    public static PlcValueAdapter processTwoOctetBinaryTime(TwoOctetBinaryTime cp16Time2a) {
        // TODO: Implement ...
        return null;
    }

    public static PlcValueAdapter processQualifierOfInterrogation(QualifierOfInterrogation qoi) {
        return PlcUINT.of(qoi.getQualifierOfCommand());
    }

    public static PlcValueAdapter processQualifierOfCounterInterrogationCommand(QualifierOfCounterInterrogationCommand qcc) {
        return new PlcStruct(Map.of(
            "freeze", PlcUSINT.of(qcc.getFreeze()),
            "request", PlcUSINT.of(qcc.getRequest())));
    }

    public static PlcValueAdapter processQualifierOfParameterOfMeasuredValues(QualifierOfParameterOfMeasuredValues qpm, PlcValueAdapter plcValue) {
        plcValue.addMetaData("parameterInOperation", PlcBOOL.of(qpm.getParameterInOperation()));
        plcValue.addMetaData("localParameterChange", PlcBOOL.of(qpm.getLocalParameterChange()));
        plcValue.addMetaData("kindOfParameter", PlcUSINT.of(qpm.getKindOfParameter()));
        return plcValue;
    }

    public static PlcValueAdapter processQualifierOfParameterActivation(QualifierOfParameterActivation qpa) {
        return PlcUINT.of(qpa.getQualifier());
    }

    public static PlcValueAdapter processQualifierOfCommand(QualifierOfCommand qoc) {
        PlcUSINT plcValue = PlcUSINT.of(qoc.getQualifier());
        plcValue.addMetaData("select", PlcBOOL.of(qoc.getSelect()));
        return plcValue;
    }

    public static PlcValueAdapter processQualifierOfResetProcessCommand(QualifierOfResetProcessCommand qrp) {
        return PlcUINT.of(qrp.getQualifier());
    }

    public static void processQualifierOfSetPointCommand(QualifierOfSetPointCommand qos, PlcValueAdapter plcValue) {
        plcValue.addMetaData("select", PlcBOOL.of(qos.getSelect()));
        plcValue.addMetaData("qualifier", PlcUSINT.of(qos.getQualifier()));
    }

    // Skipping the file handles ...

    public static PlcValueAdapter processCauseOfInitialization(CauseOfInitialization coi) {
        PlcUSINT plcValue = PlcUSINT.of(coi.getQualifier());
        plcValue.addMetaData("select", PlcBOOL.of(coi.getSelect()));
        return plcValue;
    }

    public static PlcValueAdapter processFixedTestBitPatternTwoOctet(FixedTestBitPatternTwoOctet fbp) {
        return PlcWORD.of(fbp.getPattern());
    }

}
