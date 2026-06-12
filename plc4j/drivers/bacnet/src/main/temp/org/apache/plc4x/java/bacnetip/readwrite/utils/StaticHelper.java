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
package org.apache.plc4x.java.bacnetip.readwrite.utils;

import org.apache.commons.lang3.reflect.FieldUtils;
import org.apache.commons.lang3.tuple.Pair;
import org.apache.plc4x.java.bacnetip.readwrite.*;
import org.apache.plc4x.java.spi.buffers.api.ReadBuffer;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.WriteBuffer;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Method;
import java.math.BigInteger;
import java.util.*;

public class StaticHelper {

    public static final Logger LOGGER = LoggerFactory.getLogger(StaticHelper.class);

    public static Object readEnumGenericFailing(ReadBuffer readBuffer, Long actualLength, Enum<?> template) throws BufferException {
        int bitsToRead = (int) (actualLength * 8);
        long rawValue = readBuffer.readUnsignedLong(bitsToRead, WithOption.WithName("value"));
        Class<?> declaringClass = template.getDeclaringClass();
        if (declaringClass == BACnetConfirmedServiceRequestDeviceCommunicationControlEnableDisable.class) {
            if (!BACnetConfirmedServiceRequestDeviceCommunicationControlEnableDisable.isDefined((short) rawValue))
                throw new BufferException("Invalid value " + rawValue + " for " + BACnetConfirmedServiceRequestDeviceCommunicationControlEnableDisable.class.getSimpleName());
            return BACnetConfirmedServiceRequestDeviceCommunicationControlEnableDisable.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDevice.class) {
            if (!BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDevice.isDefined((short) rawValue))
                throw new BufferException("Invalid value " + rawValue + " for " + BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDevice.class.getSimpleName());
            return BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDevice.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetSegmentation.class) {
            if (!BACnetSegmentation.isDefined((short) rawValue))
                throw new BufferException("Invalid value " + rawValue + " for " + BACnetSegmentation.class.getSimpleName());
            return BACnetSegmentation.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetAction.class) {
            if (!BACnetAction.isDefined((short) rawValue))
                throw new BufferException("Invalid value " + rawValue + " for " + BACnetAction.class.getSimpleName());
            return BACnetAction.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetNotifyType.class) {
            if (!BACnetNotifyType.isDefined((short) rawValue))
                throw new BufferException("Invalid value " + rawValue + " for " + BACnetBinaryPV.class.getSimpleName());
            return BACnetNotifyType.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetBinaryPV.class) {
            if (!BACnetBinaryPV.isDefined((short) rawValue))
                throw new BufferException("Invalid value " + rawValue + " for " + BACnetBinaryPV.class.getSimpleName());
            return BACnetBinaryPV.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetLockStatus.class) {
            if (!BACnetLockStatus.isDefined((short) rawValue))
                throw new BufferException("Invalid value " + rawValue + " for " + BACnetLockStatus.class.getSimpleName());
            return BACnetLockStatus.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetDoorSecuredStatus.class) {
            if (!BACnetDoorSecuredStatus.isDefined((short) rawValue))
                throw new BufferException("Invalid value " + rawValue + " for " + BACnetDoorSecuredStatus.class.getSimpleName());
            return BACnetDoorSecuredStatus.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetNodeType.class) {
            if (!BACnetNodeType.isDefined((short) rawValue))
                throw new BufferException("Invalid value " + rawValue + " for " + BACnetNodeType.class.getSimpleName());
            return BACnetNodeType.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetProgramState.class) {
            if (!BACnetProgramState.isDefined((short) rawValue))
                throw new BufferException("Invalid value " + rawValue + " for " + BACnetProgramState.class.getSimpleName());
            return BACnetProgramState.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetProgramRequest.class) {
            if (!BACnetProgramRequest.isDefined((short) rawValue))
                throw new BufferException("Invalid value " + rawValue + " for " + BACnetProgramRequest.class.getSimpleName());
            return BACnetProgramRequest.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetFileAccessMethod.class) {
            if (!BACnetFileAccessMethod.isDefined((short) rawValue))
                throw new BufferException("Invalid value " + rawValue + " for " + BACnetFileAccessMethod.class.getSimpleName());
            return BACnetFileAccessMethod.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetAccumulatorRecordAccumulatorStatus.class) {
            if (!BACnetAccumulatorRecordAccumulatorStatus.isDefined((short) rawValue))
                throw new BufferException("Invalid value " + rawValue + " for " + BACnetAccumulatorRecordAccumulatorStatus.class.getSimpleName());
            return BACnetAccumulatorRecordAccumulatorStatus.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetPolarity.class) {
            if (!BACnetPolarity.isDefined((short) rawValue))
                throw new BufferException("Invalid value " + rawValue + " for " + BACnetPolarity.class.getSimpleName());
            return BACnetPolarity.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetShedState.class) {
            if (!BACnetShedState.isDefined((short) rawValue))
                throw new BufferException("Invalid value " + rawValue + " for " + BACnetShedState.class.getSimpleName());
            return BACnetShedState.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetDoorValue.class) {
            if (!BACnetDoorValue.isDefined((short) rawValue))
                throw new BufferException("Invalid value " + rawValue + " for " + BACnetDoorValue.class.getSimpleName());
            return BACnetDoorValue.enumForValue((short) rawValue);
        }
        throw new BufferException("Unmapped type " + declaringClass);
    }

    public static Object readEnumGeneric(ReadBuffer readBuffer, Long actualLength, Enum<?> template) throws BufferException {
        int bitsToRead = (int) (actualLength * 8);
        long rawValue = readBuffer.readUnsignedLong(bitsToRead, WithOption.WithName("value"));
        // TODO: map types here for better performance which doesn't use reflection
        Class<?> declaringClass = template.getDeclaringClass();
        if (declaringClass == BACnetNodeType.class) {
            return BACnetNodeType.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetDataType.class) {
            return BACnetDataType.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetMaintenance.class) {
            if (!BACnetMaintenance.isDefined((short) rawValue)) return BACnetMaintenance.VENDOR_PROPRIETARY_VALUE;
            return BACnetMaintenance.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetDoorAlarmState.class) {
            if (!BACnetDoorAlarmState.isDefined((short) rawValue)) return BACnetDoorAlarmState.VENDOR_PROPRIETARY_VALUE;
            return BACnetDoorAlarmState.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetLiftGroupMode.class) {
            return BACnetLiftGroupMode.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetNotifyType.class) {
            return BACnetNotifyType.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetEscalatorOperationDirection.class) {
            if (!BACnetEscalatorOperationDirection.isDefined((short) rawValue))
                return BACnetEscalatorOperationDirection.VENDOR_PROPRIETARY_VALUE;
            return BACnetEscalatorOperationDirection.enumForValue((int) rawValue);
        } else if (declaringClass == BACnetDeviceStatus.class) {
            if (!BACnetDeviceStatus.isDefined((short) rawValue)) return BACnetDeviceStatus.VENDOR_PROPRIETARY_VALUE;
            return BACnetDeviceStatus.enumForValue((int) rawValue);
        } else if (declaringClass == BACnetRouterEntryStatus.class) {
            return BACnetRouterEntryStatus.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetAccumulatorRecordAccumulatorStatus.class) {
            return BACnetAccumulatorRecordAccumulatorStatus.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetLiftCarDirection.class) {
            if (!BACnetLiftCarDirection.isDefined((short) rawValue))
                return BACnetLiftCarDirection.VENDOR_PROPRIETARY_VALUE;
            return BACnetLiftCarDirection.enumForValue((int) rawValue);
        } else if (declaringClass == BACnetEventType.class) {
            if (!BACnetEventType.isDefined((short) rawValue)) return BACnetEventType.VENDOR_PROPRIETARY_VALUE;
            return BACnetEventType.enumForValue((int) rawValue);
        } else if (declaringClass == BACnetBinaryPV.class) {
            return BACnetBinaryPV.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetAccessRuleLocationSpecifier.class) {
            return BACnetAccessRuleLocationSpecifier.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetConfirmedServiceChoice.class) {
            return BACnetConfirmedServiceChoice.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetAccessCredentialDisable.class) {
            if (!BACnetAccessCredentialDisable.isDefined((short) rawValue))
                return BACnetAccessCredentialDisable.VENDOR_PROPRIETARY_VALUE;
            return BACnetAccessCredentialDisable.enumForValue((int) rawValue);
        } else if (declaringClass == BACnetWriteStatus.class) {
            return BACnetWriteStatus.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilter.class) {
            return BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilter.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetEventTransitionBits.class) {
            return BACnetEventTransitionBits.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetAccessZoneOccupancyState.class) {
            if (!BACnetAccessZoneOccupancyState.isDefined((short) rawValue))
                return BACnetAccessZoneOccupancyState.VENDOR_PROPRIETARY_VALUE;
            return BACnetAccessZoneOccupancyState.enumForValue((int) rawValue);
        } else if (declaringClass == BACnetNetworkNumberQuality.class) {
            return BACnetNetworkNumberQuality.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetTimerTransition.class) {
            return BACnetTimerTransition.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetAccessAuthenticationFactorDisable.class) {
            if (!BACnetAccessAuthenticationFactorDisable.isDefined((short) rawValue))
                return BACnetAccessAuthenticationFactorDisable.VENDOR_PROPRIETARY_VALUE;
            return BACnetAccessAuthenticationFactorDisable.enumForValue((int) rawValue);
        } else if (declaringClass == BACnetAccessUserType.class) {
            if (!BACnetAccessUserType.isDefined((short) rawValue)) return BACnetAccessUserType.VENDOR_PROPRIETARY_VALUE;
            return BACnetAccessUserType.enumForValue((int) rawValue);
        } else if (declaringClass == BACnetIPMode.class) {
            return BACnetIPMode.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetFaultType.class) {
            return BACnetFaultType.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetSecurityPolicy.class) {
            return BACnetSecurityPolicy.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetAuthorizationMode.class) {
            if (!BACnetAuthorizationMode.isDefined((short) rawValue))
                return BACnetAuthorizationMode.VENDOR_PROPRIETARY_VALUE;
            return BACnetAuthorizationMode.enumForValue((int) rawValue);
        } else if (declaringClass == BACnetNetworkType.class) {
            if (!BACnetNetworkType.isDefined((short) rawValue)) return BACnetNetworkType.VENDOR_PROPRIETARY_VALUE;
            return BACnetNetworkType.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetLifeSafetyMode.class) {
            if (!BACnetLifeSafetyMode.isDefined((short) rawValue)) return BACnetLifeSafetyMode.VENDOR_PROPRIETARY_VALUE;
            return BACnetLifeSafetyMode.enumForValue((int) rawValue);
        } else if (declaringClass == BACnetEventState.class) {
            if (!BACnetEventState.isDefined((short) rawValue)) return BACnetEventState.VENDOR_PROPRIETARY_VALUE;
            return BACnetEventState.enumForValue((int) rawValue);
        } else if (declaringClass == BACnetAuthorizationExemption.class) {
            if (!BACnetAuthorizationExemption.isDefined((short) rawValue))
                return BACnetAuthorizationExemption.VENDOR_PROPRIETARY_VALUE;
            return BACnetAuthorizationExemption.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetObjectTypesSupported.class) {
            return BACnetObjectTypesSupported.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetBackupState.class) {
            return BACnetBackupState.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetDoorStatus.class) {
            if (!BACnetDoorStatus.isDefined((short) rawValue)) return BACnetDoorStatus.VENDOR_PROPRIETARY_VALUE;
            return BACnetDoorStatus.enumForValue((int) rawValue);
        } else if (declaringClass == BACnetReliability.class) {
            if (!BACnetReliability.isDefined((short) rawValue)) return BACnetReliability.VENDOR_PROPRIETARY_VALUE;
            return BACnetReliability.enumForValue((int) rawValue);
        } else if (declaringClass == BACnetResultFlags.class) {
            return BACnetResultFlags.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetAbortReason.class) {
            if (!BACnetAbortReason.isDefined((short) rawValue)) return BACnetAbortReason.VENDOR_PROPRIETARY_VALUE;
            return BACnetAbortReason.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetVTClass.class) {
            if (!BACnetVTClass.isDefined((short) rawValue)) return BACnetVTClass.VENDOR_PROPRIETARY_VALUE;
            return BACnetVTClass.enumForValue((int) rawValue);
        } else if (declaringClass == BACnetTimerState.class) {
            return BACnetTimerState.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetConfirmedServiceRequestGetEnrollmentSummaryEventStateFilter.class) {
            return BACnetConfirmedServiceRequestGetEnrollmentSummaryEventStateFilter.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetUnconfirmedServiceChoice.class) {
            return BACnetUnconfirmedServiceChoice.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetEngineeringUnits.class) {
            if (!BACnetEngineeringUnits.isDefined((short) rawValue))
                return BACnetEngineeringUnits.VENDOR_PROPRIETARY_VALUE;
            return BACnetEngineeringUnits.enumForValue((int) rawValue);
        } else if (declaringClass == BACnetNetworkPortCommand.class) {
            if (!BACnetNetworkPortCommand.isDefined((short) rawValue))
                return BACnetNetworkPortCommand.VENDOR_PROPRIETARY_VALUE;
            return BACnetNetworkPortCommand.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetAuthenticationStatus.class) {
            return BACnetAuthenticationStatus.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetLockStatus.class) {
            return BACnetLockStatus.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetShedState.class) {
            return BACnetShedState.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetSilencedState.class) {
            if (!BACnetSilencedState.isDefined((short) rawValue)) return BACnetSilencedState.VENDOR_PROPRIETARY_VALUE;
            return BACnetSilencedState.enumForValue((int) rawValue);
        } else if (declaringClass == BACnetDoorValue.class) {
            return BACnetDoorValue.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetDoorSecuredStatus.class) {
            return BACnetDoorSecuredStatus.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetFileAccessMethod.class) {
            return BACnetFileAccessMethod.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetLiftCarMode.class) {
            if (!BACnetLiftCarMode.isDefined((short) rawValue)) return BACnetLiftCarMode.VENDOR_PROPRIETARY_VALUE;
            return BACnetLiftCarMode.enumForValue((int) rawValue);
        } else if (declaringClass == BACnetStatusFlags.class) {
            return BACnetStatusFlags.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetLightingInProgress.class) {
            return BACnetLightingInProgress.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetRestartReason.class) {
            return BACnetRestartReason.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetAccessPassbackMode.class) {
            return BACnetAccessPassbackMode.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetLiftCarDoorCommand.class) {
            return BACnetLiftCarDoorCommand.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetSegmentation.class) {
            if (!BACnetSegmentation.isDefined((short) rawValue))
                LOGGER.error("{} not defined for segmentation falling back to no segmentation", rawValue);
            return BACnetSegmentation.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetAccessCredentialDisableReason.class) {
            if (!BACnetAccessCredentialDisableReason.isDefined((short) rawValue))
                return BACnetAccessCredentialDisableReason.VENDOR_PROPRIETARY_VALUE;
            return BACnetAccessCredentialDisableReason.enumForValue((int) rawValue);
        } else if (declaringClass == BACnetPolarity.class) {
            return BACnetPolarity.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetLifeSafetyState.class) {
            if (!BACnetLifeSafetyState.isDefined((short) rawValue))
                return BACnetLifeSafetyState.VENDOR_PROPRIETARY_VALUE;
            return BACnetLifeSafetyState.enumForValue((int) rawValue);
        } else if (declaringClass == BACnetSecurityLevel.class) {
            return BACnetSecurityLevel.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetBinaryLightingPV.class) {
            return BACnetBinaryLightingPV.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetLifeSafetyOperation.class) {
            if (!BACnetLifeSafetyOperation.isDefined((short) rawValue))
                return BACnetLifeSafetyOperation.VENDOR_PROPRIETARY_VALUE;
            return BACnetLifeSafetyOperation.enumForValue((int) rawValue);
        } else if (declaringClass == BACnetDaysOfWeek.class) {
            return BACnetDaysOfWeek.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetRejectReason.class) {
            if (!BACnetRejectReason.isDefined((short) rawValue)) return BACnetRejectReason.VENDOR_PROPRIETARY_VALUE;
            return BACnetRejectReason.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetLogStatus.class) {
            return BACnetLogStatus.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetLiftFault.class) {
            if (!BACnetLiftFault.isDefined((short) rawValue)) return BACnetLiftFault.VENDOR_PROPRIETARY_VALUE;
            return BACnetLiftFault.enumForValue((int) rawValue);
        } else if (declaringClass == BACnetAction.class) {
            return BACnetAction.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetPropertyIdentifier.class) {
            if (!BACnetPropertyIdentifier.isDefined((short) rawValue))
                return BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE;
            return BACnetPropertyIdentifier.enumForValue((int) rawValue);
        } else if (declaringClass == BACnetEscalatorMode.class) {
            if (!BACnetEscalatorMode.isDefined((short) rawValue)) return BACnetEscalatorMode.VENDOR_PROPRIETARY_VALUE;
            return BACnetEscalatorMode.enumForValue((int) rawValue);
        } else if (declaringClass == BACnetProgramRequest.class) {
            return BACnetProgramRequest.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetConfirmedServiceRequestConfirmedTextMessageMessagePriority.class) {
            return BACnetConfirmedServiceRequestConfirmedTextMessageMessagePriority.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetProtocolLevel.class) {
            return BACnetProtocolLevel.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetRelationship.class) {
            if (!BACnetRelationship.isDefined((short) rawValue)) return BACnetRelationship.VENDOR_PROPRIETARY_VALUE;
            return BACnetRelationship.enumForValue((int) rawValue);
        } else if (declaringClass == BACnetLoggingType.class) {
            return BACnetLoggingType.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetAccessRuleTimeRangeSpecifier.class) {
            return BACnetAccessRuleTimeRangeSpecifier.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetConfirmedServiceRequestDeviceCommunicationControlEnableDisable.class) {
            return BACnetConfirmedServiceRequestDeviceCommunicationControlEnableDisable.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetCharacterEncoding.class) {
            return BACnetCharacterEncoding.enumForValue((byte) rawValue);
        } else if (declaringClass == BACnetServicesSupported.class) {
            return BACnetServicesSupported.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetLimitEnable.class) {
            return BACnetLimitEnable.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetEscalatorFault.class) {
            if (!BACnetEscalatorFault.isDefined((short) rawValue)) return BACnetEscalatorFault.VENDOR_PROPRIETARY_VALUE;
            return BACnetEscalatorFault.enumForValue((int) rawValue);
        } else if (declaringClass == BACnetVendorId.class) {
            if (!BACnetVendorId.isDefined((short) rawValue)) return BACnetVendorId.UNKNOWN_VENDOR;
            return BACnetVendorId.enumForValue((int) rawValue);
        } else if (declaringClass == BACnetAuthenticationFactorType.class) {
            return BACnetAuthenticationFactorType.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetLightingTransition.class) {
            if (!BACnetLightingTransition.isDefined((short) rawValue))
                return BACnetLightingTransition.VENDOR_PROPRIETARY_VALUE;
            return BACnetLightingTransition.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetLightingOperation.class) {
            if (!BACnetLightingOperation.isDefined((short) rawValue))
                return BACnetLightingOperation.VENDOR_PROPRIETARY_VALUE;
            return BACnetLightingOperation.enumForValue((int) rawValue);
        } else if (declaringClass == BACnetProgramError.class) {
            if (!BACnetProgramError.isDefined((short) rawValue)) return BACnetProgramError.VENDOR_PROPRIETARY_VALUE;
            return BACnetProgramError.enumForValue((int) rawValue);
        } else if (declaringClass == BACnetLiftCarDriveStatus.class) {
            if (!BACnetLiftCarDriveStatus.isDefined((short) rawValue))
                return BACnetLiftCarDriveStatus.VENDOR_PROPRIETARY_VALUE;
            return BACnetLiftCarDriveStatus.enumForValue((int) rawValue);
        } else if (declaringClass == BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDevice.class) {
            if (!BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDevice.isDefined((short) rawValue))
                return BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDevice.VENDOR_PROPRIETARY_VALUE;
            return BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDevice.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetObjectType.class) {
            if (!BACnetObjectType.isDefined((short) rawValue)) return BACnetObjectType.VENDOR_PROPRIETARY_VALUE;
            return BACnetObjectType.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetAccessEvent.class) {
            if (!BACnetAccessEvent.isDefined((short) rawValue)) return BACnetAccessEvent.VENDOR_PROPRIETARY_VALUE;
            return BACnetAccessEvent.enumForValue((int) rawValue);
        } else if (declaringClass == ErrorClass.class) {
            if (!ErrorClass.isDefined((short) rawValue)) return ErrorClass.VENDOR_PROPRIETARY_VALUE;
            return ErrorClass.enumForValue((short) rawValue);
        } else if (declaringClass == ErrorCode.class) {
            if (!ErrorCode.isDefined((short) rawValue)) return ErrorCode.VENDOR_PROPRIETARY_VALUE;
            return ErrorCode.enumForValue((short) rawValue);
        } else {
            LOGGER.warn("read: using reflection for {}", declaringClass);
            Optional<Method> enumForValue = Arrays.stream(declaringClass.getDeclaredMethods()).filter(method -> method.getName().equals("enumForValue")).findAny();
            if (!enumForValue.isPresent()) {
                throw new BufferException("No enumForValue available");
            }
            Method method = enumForValue.get();
            try {
                Class<?> parameterType = method.getParameterTypes()[0];
                Object paramValue = null;
                if (parameterType == byte.class || parameterType == Byte.class) {
                    paramValue = (byte) rawValue;
                } else if (parameterType == short.class || parameterType == Short.class) {
                    paramValue = (short) rawValue;
                } else if (parameterType == int.class || parameterType == Integer.class) {
                    paramValue = (int) rawValue;
                } else if (parameterType == long.class || parameterType == Long.class) {
                    paramValue = (int) rawValue;
                }
                Object result = method.invoke(null, paramValue);
                return Objects.requireNonNullElseGet(result, () -> Enum.valueOf(template.getDeclaringClass(), "VENDOR_PROPRIETARY_VALUE"));
            } catch (IllegalAccessException | InvocationTargetException e) {
                throw new BufferException("error invoking method", e);
            }
        }
    }

    public static long readProprietaryEnumGeneric(ReadBuffer readBuffer, Long actualLength, boolean shouldRead) throws BufferException {
        if (!shouldRead) {
            return 0L;
        }
        // We need to reset our reader to the position we read before
        readBuffer.setPositionInBits((int) (readBuffer.getPositionInBits() - actualLength));
        int bitsToRead = (int) (actualLength * 8);
        return readBuffer.readUnsignedLong(bitsToRead, WithOption.WithName("proprietaryValue"));
    }

    public static void writeEnumGeneric(WriteBuffer writeBuffer, Enum<?> value) throws BufferException {
        if (value == null) {
            return;
        }
        // TODO: generic name mapping for now
        if (value.name().equals("VENDOR_PROPRIETARY_VALUE")) return;
        int bitsToWrite;
        long valueValue;
        // TODO: map types here for better performance which doesn't use reflection
        if (value.getDeclaringClass() == BACnetNodeType.class) {
            valueValue = ((BACnetNodeType) value).getValue();
        } else if (value.getDeclaringClass() == BACnetDataType.class) {
            valueValue = ((BACnetDataType) value).getValue();
        } else if (value.getDeclaringClass() == BACnetMaintenance.class) {
            valueValue = ((BACnetMaintenance) value).getValue();
        } else if (value.getDeclaringClass() == BACnetDoorAlarmState.class) {
            valueValue = ((BACnetDoorAlarmState) value).getValue();
        } else if (value.getDeclaringClass() == BACnetLiftGroupMode.class) {
            valueValue = ((BACnetLiftGroupMode) value).getValue();
        } else if (value.getDeclaringClass() == BACnetNotifyType.class) {
            valueValue = ((BACnetNotifyType) value).getValue();
        } else if (value.getDeclaringClass() == BACnetEscalatorOperationDirection.class) {
            valueValue = ((BACnetEscalatorOperationDirection) value).getValue();
        } else if (value.getDeclaringClass() == BACnetDeviceStatus.class) {
            valueValue = ((BACnetDeviceStatus) value).getValue();
        } else if (value.getDeclaringClass() == BACnetRouterEntryStatus.class) {
            valueValue = ((BACnetRouterEntryStatus) value).getValue();
        } else if (value.getDeclaringClass() == BACnetAccumulatorRecordAccumulatorStatus.class) {
            valueValue = ((BACnetAccumulatorRecordAccumulatorStatus) value).getValue();
        } else if (value.getDeclaringClass() == BACnetLiftCarDirection.class) {
            valueValue = ((BACnetLiftCarDirection) value).getValue();
        } else if (value.getDeclaringClass() == BACnetEventType.class) {
            valueValue = ((BACnetEventType) value).getValue();
        } else if (value.getDeclaringClass() == BACnetBinaryPV.class) {
            valueValue = ((BACnetBinaryPV) value).getValue();
        } else if (value.getDeclaringClass() == BACnetAccessRuleLocationSpecifier.class) {
            valueValue = ((BACnetAccessRuleLocationSpecifier) value).getValue();
        } else if (value.getDeclaringClass() == BACnetConfirmedServiceChoice.class) {
            valueValue = ((BACnetConfirmedServiceChoice) value).getValue();
        } else if (value.getDeclaringClass() == BACnetAccessCredentialDisable.class) {
            valueValue = ((BACnetAccessCredentialDisable) value).getValue();
        } else if (value.getDeclaringClass() == BACnetWriteStatus.class) {
            valueValue = ((BACnetWriteStatus) value).getValue();
        } else if (value.getDeclaringClass() == BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilter.class) {
            valueValue = ((BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilter) value).getValue();
        } else if (value.getDeclaringClass() == BACnetProgramState.class) {
            valueValue = ((BACnetProgramState) value).getValue();
        } else if (value.getDeclaringClass() == BACnetEventTransitionBits.class) {
            valueValue = ((BACnetEventTransitionBits) value).getValue();
        } else if (value.getDeclaringClass() == BACnetAccessZoneOccupancyState.class) {
            valueValue = ((BACnetAccessZoneOccupancyState) value).getValue();
        } else if (value.getDeclaringClass() == BACnetNetworkNumberQuality.class) {
            valueValue = ((BACnetNetworkNumberQuality) value).getValue();
        } else if (value.getDeclaringClass() == BACnetTimerTransition.class) {
            valueValue = ((BACnetTimerTransition) value).getValue();
        } else if (value.getDeclaringClass() == BACnetAccessAuthenticationFactorDisable.class) {
            valueValue = ((BACnetAccessAuthenticationFactorDisable) value).getValue();
        } else if (value.getDeclaringClass() == BACnetAccessUserType.class) {
            valueValue = ((BACnetAccessUserType) value).getValue();
        } else if (value.getDeclaringClass() == BACnetIPMode.class) {
            valueValue = ((BACnetIPMode) value).getValue();
        } else if (value.getDeclaringClass() == BACnetFaultType.class) {
            valueValue = ((BACnetFaultType) value).getValue();
        } else if (value.getDeclaringClass() == BACnetSecurityPolicy.class) {
            valueValue = ((BACnetSecurityPolicy) value).getValue();
        } else if (value.getDeclaringClass() == BACnetAuthorizationMode.class) {
            valueValue = ((BACnetAuthorizationMode) value).getValue();
        } else if (value.getDeclaringClass() == BACnetNetworkType.class) {
            valueValue = ((BACnetNetworkType) value).getValue();
        } else if (value.getDeclaringClass() == BACnetLifeSafetyMode.class) {
            valueValue = ((BACnetLifeSafetyMode) value).getValue();
        } else if (value.getDeclaringClass() == BACnetEventState.class) {
            valueValue = ((BACnetEventState) value).getValue();
        } else if (value.getDeclaringClass() == BACnetAuthorizationExemption.class) {
            valueValue = ((BACnetAuthorizationExemption) value).getValue();
        } else if (value.getDeclaringClass() == BACnetObjectTypesSupported.class) {
            valueValue = ((BACnetObjectTypesSupported) value).getValue();
        } else if (value.getDeclaringClass() == BACnetBackupState.class) {
            valueValue = ((BACnetBackupState) value).getValue();
        } else if (value.getDeclaringClass() == BACnetDoorStatus.class) {
            valueValue = ((BACnetDoorStatus) value).getValue();
        } else if (value.getDeclaringClass() == BACnetReliability.class) {
            valueValue = ((BACnetReliability) value).getValue();
        } else if (value.getDeclaringClass() == BACnetResultFlags.class) {
            valueValue = ((BACnetResultFlags) value).getValue();
        } else if (value.getDeclaringClass() == BACnetAbortReason.class) {
            valueValue = ((BACnetAbortReason) value).getValue();
        } else if (value.getDeclaringClass() == BACnetVTClass.class) {
            valueValue = ((BACnetVTClass) value).getValue();
        } else if (value.getDeclaringClass() == BACnetTimerState.class) {
            valueValue = ((BACnetTimerState) value).getValue();
        } else if (value.getDeclaringClass() == BACnetConfirmedServiceRequestGetEnrollmentSummaryEventStateFilter.class) {
            valueValue = ((BACnetConfirmedServiceRequestGetEnrollmentSummaryEventStateFilter) value).getValue();
        } else if (value.getDeclaringClass() == BACnetUnconfirmedServiceChoice.class) {
            valueValue = ((BACnetUnconfirmedServiceChoice) value).getValue();
        } else if (value.getDeclaringClass() == BACnetEngineeringUnits.class) {
            valueValue = ((BACnetEngineeringUnits) value).getValue();
        } else if (value.getDeclaringClass() == BACnetNetworkPortCommand.class) {
            valueValue = ((BACnetNetworkPortCommand) value).getValue();
        } else if (value.getDeclaringClass() == BACnetAuthenticationStatus.class) {
            valueValue = ((BACnetAuthenticationStatus) value).getValue();
        } else if (value.getDeclaringClass() == BACnetLockStatus.class) {
            valueValue = ((BACnetLockStatus) value).getValue();
        } else if (value.getDeclaringClass() == BACnetShedState.class) {
            valueValue = ((BACnetShedState) value).getValue();
        } else if (value.getDeclaringClass() == BACnetSilencedState.class) {
            valueValue = ((BACnetSilencedState) value).getValue();
        } else if (value.getDeclaringClass() == BACnetDoorValue.class) {
            valueValue = ((BACnetDoorValue) value).getValue();
        } else if (value.getDeclaringClass() == BACnetDoorSecuredStatus.class) {
            valueValue = ((BACnetDoorSecuredStatus) value).getValue();
        } else if (value.getDeclaringClass() == BACnetFileAccessMethod.class) {
            valueValue = ((BACnetFileAccessMethod) value).getValue();
        } else if (value.getDeclaringClass() == BACnetLiftCarMode.class) {
            valueValue = ((BACnetLiftCarMode) value).getValue();
        } else if (value.getDeclaringClass() == BACnetStatusFlags.class) {
            valueValue = ((BACnetStatusFlags) value).getValue();
        } else if (value.getDeclaringClass() == BACnetLightingInProgress.class) {
            valueValue = ((BACnetLightingInProgress) value).getValue();
        } else if (value.getDeclaringClass() == BACnetRestartReason.class) {
            valueValue = ((BACnetRestartReason) value).getValue();
        } else if (value.getDeclaringClass() == BACnetAccessPassbackMode.class) {
            valueValue = ((BACnetAccessPassbackMode) value).getValue();
        } else if (value.getDeclaringClass() == BACnetLiftCarDoorCommand.class) {
            valueValue = ((BACnetLiftCarDoorCommand) value).getValue();
        } else if (value.getDeclaringClass() == BACnetSegmentation.class) {
            valueValue = ((BACnetSegmentation) value).getValue();
        } else if (value.getDeclaringClass() == BACnetAccessCredentialDisableReason.class) {
            valueValue = ((BACnetAccessCredentialDisableReason) value).getValue();
        } else if (value.getDeclaringClass() == BACnetPolarity.class) {
            valueValue = ((BACnetPolarity) value).getValue();
        } else if (value.getDeclaringClass() == BACnetLifeSafetyState.class) {
            valueValue = ((BACnetLifeSafetyState) value).getValue();
        } else if (value.getDeclaringClass() == BACnetSecurityLevel.class) {
            valueValue = ((BACnetSecurityLevel) value).getValue();
        } else if (value.getDeclaringClass() == BACnetBinaryLightingPV.class) {
            valueValue = ((BACnetBinaryLightingPV) value).getValue();
        } else if (value.getDeclaringClass() == BACnetLifeSafetyOperation.class) {
            valueValue = ((BACnetLifeSafetyOperation) value).getValue();
        } else if (value.getDeclaringClass() == BACnetDaysOfWeek.class) {
            valueValue = ((BACnetDaysOfWeek) value).getValue();
        } else if (value.getDeclaringClass() == BACnetRejectReason.class) {
            valueValue = ((BACnetRejectReason) value).getValue();
        } else if (value.getDeclaringClass() == BACnetLogStatus.class) {
            valueValue = ((BACnetLogStatus) value).getValue();
        } else if (value.getDeclaringClass() == BACnetLiftFault.class) {
            valueValue = ((BACnetLiftFault) value).getValue();
        } else if (value.getDeclaringClass() == BACnetAction.class) {
            valueValue = ((BACnetAction) value).getValue();
        } else if (value.getDeclaringClass() == BACnetPropertyIdentifier.class) {
            valueValue = ((BACnetPropertyIdentifier) value).getValue();
        } else if (value.getDeclaringClass() == BACnetEscalatorMode.class) {
            valueValue = ((BACnetEscalatorMode) value).getValue();
        } else if (value.getDeclaringClass() == BACnetProgramRequest.class) {
            valueValue = ((BACnetProgramRequest) value).getValue();
        } else if (value.getDeclaringClass() == BACnetConfirmedServiceRequestConfirmedTextMessageMessagePriority.class) {
            valueValue = ((BACnetConfirmedServiceRequestConfirmedTextMessageMessagePriority) value).getValue();
        } else if (value.getDeclaringClass() == BACnetProtocolLevel.class) {
            valueValue = ((BACnetProtocolLevel) value).getValue();
        } else if (value.getDeclaringClass() == BACnetRelationship.class) {
            valueValue = ((BACnetRelationship) value).getValue();
        } else if (value.getDeclaringClass() == BACnetLoggingType.class) {
            valueValue = ((BACnetLoggingType) value).getValue();
        } else if (value.getDeclaringClass() == BACnetAccessRuleTimeRangeSpecifier.class) {
            valueValue = ((BACnetAccessRuleTimeRangeSpecifier) value).getValue();
        } else if (value.getDeclaringClass() == BACnetConfirmedServiceRequestDeviceCommunicationControlEnableDisable.class) {
            valueValue = ((BACnetConfirmedServiceRequestDeviceCommunicationControlEnableDisable) value).getValue();
        } else if (value.getDeclaringClass() == BACnetCharacterEncoding.class) {
            valueValue = ((BACnetCharacterEncoding) value).getValue();
        } else if (value.getDeclaringClass() == BACnetServicesSupported.class) {
            valueValue = ((BACnetServicesSupported) value).getValue();
        } else if (value.getDeclaringClass() == BACnetLimitEnable.class) {
            valueValue = ((BACnetLimitEnable) value).getValue();
        } else if (value.getDeclaringClass() == BACnetEscalatorFault.class) {
            valueValue = ((BACnetEscalatorFault) value).getValue();
        } else if (value.getDeclaringClass() == BACnetVendorId.class) {
            valueValue = ((BACnetVendorId) value).getValue();
        } else if (value.getDeclaringClass() == BACnetAuthenticationFactorType.class) {
            valueValue = ((BACnetAuthenticationFactorType) value).getValue();
        } else if (value.getDeclaringClass() == BACnetLightingTransition.class) {
            valueValue = ((BACnetLightingTransition) value).getValue();
        } else if (value.getDeclaringClass() == BACnetLightingOperation.class) {
            valueValue = ((BACnetLightingOperation) value).getValue();
        } else if (value.getDeclaringClass() == BACnetProgramError.class) {
            valueValue = ((BACnetProgramError) value).getValue();
        } else if (value.getDeclaringClass() == BACnetLiftCarDriveStatus.class) {
            valueValue = ((BACnetLiftCarDriveStatus) value).getValue();
        } else if (value.getDeclaringClass() == BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDevice.class) {
            valueValue = ((BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDevice) value).getValue();
        } else if (value.getDeclaringClass() == BACnetObjectType.class) {
            valueValue = ((BACnetObjectType) value).getValue();
        } else if (value.getDeclaringClass() == BACnetAccessEvent.class) {
            valueValue = ((BACnetAccessEvent) value).getValue();
        } else if (value.getDeclaringClass() == ErrorClass.class) {
            valueValue = ((ErrorClass) value).getValue();
        } else if (value.getDeclaringClass() == ErrorCode.class) {
            valueValue = ((ErrorCode) value).getValue();
        } else {
            LOGGER.warn("write: using reflection for {}", value.getDeclaringClass());
            try {
                valueValue = ((Number) FieldUtils.getDeclaredField(value.getDeclaringClass(), "value", true).get(value)).longValue();
            } catch (IllegalAccessException e) {
                throw new BufferException("error accessing value", e);
            }
        }

        if (valueValue <= 0xffL) {
            bitsToWrite = 8;
        } else if (valueValue <= 0xffffL) {
            bitsToWrite = 16;
        } else if (valueValue <= 0xffffffffL) {
            bitsToWrite = 32;
        } else {
            bitsToWrite = 32;
        }
        writeBuffer.writeUnsignedLong(bitsToWrite, valueValue, WithOption.WithName("value"), WithOption.WithAdditionalStringRepresentation(value.name()), WithOption.WithEncoding(value.name()));
    }

    public static void writeProprietaryEnumGeneric(WriteBuffer writeBuffer, long value, boolean shouldWrite) throws BufferException {
        if (!shouldWrite) {
            return;
        }
        int bitsToWrite;
        if (value <= 0xffL) {
            bitsToWrite = 8;
        } else if (value <= 0xffffL) {
            bitsToWrite = 16;
        } else if (value <= 0xffffffffL) {
            bitsToWrite = 32;
        } else {
            bitsToWrite = 32;
        }
        writeBuffer.writeUnsignedLong(bitsToWrite, value, WithOption.WithName("proprietaryValue"), WithOption.WithEncoding("VENDOR_PROPRIETARY_VALUE"));
    }

    @Deprecated
    public static BACnetObjectType readObjectType(ReadBuffer readBuffer) throws BufferException {
        short readUnsignedShort = readBuffer.readUnsignedShort(10, WithOption.WithName("objectType"));
        if (!BACnetObjectType.isDefined(readUnsignedShort)) {
            return BACnetObjectType.VENDOR_PROPRIETARY_VALUE;
        }
        return BACnetObjectType.enumForValue(readUnsignedShort);
    }

    @Deprecated
    public static Short readProprietaryObjectType(ReadBuffer readBuffer, BACnetObjectType value) throws BufferException {
        if (value != null && value != BACnetObjectType.VENDOR_PROPRIETARY_VALUE) {
            return 0;
        }
        // We need to reset our reader to the position we read before
        // TODO: maybe we reset to much here because pos is byte based
        // we consume the leftover bits before we reset to avoid trouble
        // TODO: we really need bit precision on resetting
        readBuffer.readUnsignedInt(6);
        readBuffer.setPositionInBits(readBuffer.getPositionInBits() - 2);
        return readBuffer.readUnsignedShort(10, WithOption.WithName("proprietaryObjectType"));
    }

    @Deprecated
    public static void writeObjectType(WriteBuffer writeBuffer, BACnetObjectType value) throws BufferException {
        if (value == null || value == BACnetObjectType.VENDOR_PROPRIETARY_VALUE) {
            return;
        }
        writeBuffer.writeUnsignedLong(10, value.getValue(), WithOption.WithName("objectType"), WithOption.WithAdditionalStringRepresentation(value.name()), WithOption.WithEncoding(value.name()));
    }

    @Deprecated
    public static void writeProprietaryObjectType(WriteBuffer writeBuffer, BACnetObjectType objectType, int value) throws BufferException {
        if (objectType != null && objectType != BACnetObjectType.VENDOR_PROPRIETARY_VALUE) {
            return;
        }
        writeBuffer.writeUnsignedInt(10, value, WithOption.WithName("proprietaryObjectType"), WithOption.WithEncoding(BACnetObjectType.VENDOR_PROPRIETARY_VALUE.name()));
    }

    @Deprecated
    public static BACnetObjectType mapBACnetObjectType(BACnetContextTagEnumerated rawObjectType) {
        if (rawObjectType == null) return null;
        BACnetObjectType baCnetObjectType = BACnetObjectType.enumForValue((short) rawObjectType.getActualValue());
        if (baCnetObjectType == null) return BACnetObjectType.VENDOR_PROPRIETARY_VALUE;
        return baCnetObjectType;
    }

    public static boolean isBACnetConstructedDataClosingTag(ReadBuffer readBuffer, boolean instantTerminate, int expectedTagNumber) {
        if (instantTerminate) {
            return true;
        }
        int oldPos = readBuffer.getPositionInBits();
        try {
            // TODO: add graceful exit if we know already that we are at the end (we might need to add available bytes to reader)
            int tagNumber = readBuffer.readUnsignedInt(4);
            boolean isContextTag = readBuffer.readBit();
            int tagValue = readBuffer.readUnsignedInt(3);

            boolean foundOurClosingTag = isContextTag && tagNumber == expectedTagNumber && tagValue == 0x7;
            LOGGER.debug("Checking at pos pos:{}: tagNumber:{}, isContextTag:{}, tagValue:{}, expectedTagNumber:{}. foundOurClosingTag:{}", oldPos, tagNumber, isContextTag, tagValue, expectedTagNumber, foundOurClosingTag);
            return foundOurClosingTag;
        } catch (BufferException e) {
            LOGGER.warn("Error reading termination bit", e);
            return true;
        } catch (ArrayIndexOutOfBoundsException e) {
            LOGGER.debug("Reached EOF at {}", oldPos, e);
            return true;
        } finally {
            readBuffer.setPositionInBits(oldPos);
        }
    }

    public static long parseVarUint(byte[] data) {
        if (data.length == 0) {
            return 0;
        }
        return new BigInteger(data).longValue();
    }

    public static byte[] writeVarUint(long value) {
        return BigInteger.valueOf(value).toByteArray();
    }

    public static BACnetTagHeader createBACnetTagHeaderBalanced(boolean isContext, short id, long value) {
        TagClass tagClass = TagClass.APPLICATION_TAGS;
        if (isContext) {
            tagClass = TagClass.CONTEXT_SPECIFIC_TAGS;
        }

        byte tagNumber;
        Short extTagNumber = null;
        if (id <= 14) {
            tagNumber = (byte) id;
        } else {
            tagNumber = 0xF;
            extTagNumber = id;
        }

        byte lengthValueType;
        Short extLength = null;
        Integer extExtLength = null;
        Long extExtExtLength = null;
        if (value <= 4) {
            lengthValueType = (byte) value;
        } else {
            lengthValueType = 5;
            // Depending on the length, we will either write it as an 8 bit, 32 bit, or 64 bit integer
            if (value <= 253) {
                extLength = (short) value;
            } else if (value <= 65535) {
                extLength = 254;
                extExtLength = (int) value;
            } else {
                extLength = 255;
                extExtExtLength = value;
            }
        }

        return new BACnetTagHeader(tagNumber, tagClass, lengthValueType, extTagNumber, extLength, extExtLength, extExtExtLength);
    }

    public static BACnetApplicationTagNull createBACnetApplicationTagNull() {
        BACnetTagHeader header = createBACnetTagHeaderBalanced(false, BACnetDataType.NULL.getValue(), 0);
        return new BACnetApplicationTagNull(header);
    }

    public static BACnetContextTagNull createBACnetContextTagNull(byte tagNumber) {
        BACnetTagHeader header = createBACnetTagHeaderBalanced(true, tagNumber, 0);
        return new BACnetContextTagNull(header);
    }

    public static BACnetOpeningTag createBACnetOpeningTag(short tagNum) {
        byte tagNumber;
        Short extTagNumber = null;
        if (tagNum <= 14) {
            tagNumber = (byte) tagNum;
        } else {
            tagNumber = 0xF;
            extTagNumber = tagNum;
        }
        BACnetTagHeader header = new BACnetTagHeader(tagNumber, TagClass.CONTEXT_SPECIFIC_TAGS, (byte) 0x6, extTagNumber, null, null, null);
        return new BACnetOpeningTag(header);
    }

    public static BACnetClosingTag createBACnetClosingTag(short tagNum) {
        byte tagNumber;
        Short extTagNumber = null;
        if (tagNum <= 14) {
            tagNumber = (byte) tagNum;
        } else {
            tagNumber = 0xF;
            extTagNumber = tagNum;
        }
        BACnetTagHeader header = new BACnetTagHeader(tagNumber, TagClass.CONTEXT_SPECIFIC_TAGS, (byte) 0x7, extTagNumber, null, null, null);
        return new BACnetClosingTag(header);
    }

    public static BACnetApplicationTagObjectIdentifier createBACnetApplicationTagObjectIdentifier(short objectType, int instance) {
        BACnetTagHeader header = new BACnetTagHeader((byte) BACnetDataType.BACNET_OBJECT_IDENTIFIER.getValue(), TagClass.APPLICATION_TAGS, (byte) 4, null, null, null, null);
        BACnetObjectType objectTypeEnum = BACnetObjectType.enumForValue(objectType);
        short proprietaryValue = 0;
        if (objectType >= 128 || !BACnetObjectType.isDefined(objectType)) {
            objectTypeEnum = BACnetObjectType.VENDOR_PROPRIETARY_VALUE;
            proprietaryValue = objectType;
        }
        BACnetTagPayloadObjectIdentifier payload = new BACnetTagPayloadObjectIdentifier(objectTypeEnum, proprietaryValue, instance);
        return new BACnetApplicationTagObjectIdentifier(header, payload);
    }

    public static BACnetContextTagObjectIdentifier createBACnetContextTagObjectIdentifier(byte tagNum, short objectType, int instance) {
        BACnetTagHeader header = new BACnetTagHeader(tagNum, TagClass.CONTEXT_SPECIFIC_TAGS, (byte) 4, null, null, null, null);
        BACnetObjectType objectTypeEnum = BACnetObjectType.enumForValue(objectType);
        short proprietaryValue = 0;
        if (objectType >= 128 || !BACnetObjectType.isDefined(objectType)) {
            objectTypeEnum = BACnetObjectType.VENDOR_PROPRIETARY_VALUE;
            proprietaryValue = objectType;
        }
        BACnetTagPayloadObjectIdentifier payload = new BACnetTagPayloadObjectIdentifier(objectTypeEnum, proprietaryValue, instance);
        return new BACnetContextTagObjectIdentifier(header, payload);
    }

    public static BACnetPropertyIdentifierTagged createBACnetPropertyIdentifierTagged(byte tagNum, int propertyType) {
        BACnetPropertyIdentifier propertyIdentifier = BACnetPropertyIdentifier.enumForValue(propertyType);
        long proprietaryValue = 0;
        if (!BACnetPropertyIdentifier.isDefined(propertyType)) {
            propertyIdentifier = BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE;
            proprietaryValue = propertyType;
        }
        BACnetTagHeader header = new BACnetTagHeader(tagNum, TagClass.CONTEXT_SPECIFIC_TAGS, (byte) requiredLength(propertyType), null, null, null, null);
        return new BACnetPropertyIdentifierTagged(header, propertyIdentifier, proprietaryValue);
    }

    public static BACnetVendorIdTagged createBACnetVendorIdApplicationTagged(int vendorId) {
        BACnetVendorId baCnetVendorId = BACnetVendorId.enumForValue(vendorId);
        long unknownVendorId = 0;
        if (!BACnetVendorId.isDefined(vendorId)) {
            baCnetVendorId = BACnetVendorId.UNKNOWN_VENDOR;
            unknownVendorId = vendorId;
        }
        BACnetTagHeader header = new BACnetTagHeader((byte) 0x2, TagClass.APPLICATION_TAGS, (byte) requiredLength(vendorId), null, null, null, null);
        return new BACnetVendorIdTagged(header, baCnetVendorId, unknownVendorId);
    }

    public static BACnetVendorIdTagged createBACnetVendorIdContextTagged(byte tagNum, int vendorId) {
        BACnetVendorId baCnetVendorId = BACnetVendorId.enumForValue(vendorId);
        long unknownVendorId = 0;
        if (!BACnetVendorId.isDefined(vendorId)) {
            baCnetVendorId = BACnetVendorId.UNKNOWN_VENDOR;
            unknownVendorId = vendorId;
        }
        BACnetTagHeader header = new BACnetTagHeader(tagNum, TagClass.CONTEXT_SPECIFIC_TAGS, (byte) requiredLength(vendorId), null, null, null, null);
        return new BACnetVendorIdTagged(header, baCnetVendorId, unknownVendorId);
    }

    public static BACnetSegmentationTagged creatBACnetSegmentationTagged(BACnetSegmentation value) {
        BACnetTagHeader header = createBACnetTagHeaderBalanced(false, (byte) 9, 1);
        return new BACnetSegmentationTagged(header, value);
    }

    public static BACnetApplicationTagBoolean createBACnetApplicationTagBoolean(boolean value) {
        BACnetTagHeader header = createBACnetTagHeaderBalanced(false, BACnetDataType.BOOLEAN.getValue(), value ? 1L : 0L);
        return new BACnetApplicationTagBoolean(header, new BACnetTagPayloadBoolean(value ? 1L : 0L));
    }

    public static BACnetContextTagBoolean createBACnetContextTagBoolean(byte tagNumber, boolean value) {
        BACnetTagHeader header = createBACnetTagHeaderBalanced(true, tagNumber, 1);
        return new BACnetContextTagBoolean(header, (short) (value ? 1 : 0), new BACnetTagPayloadBoolean(value ? 1L : 0L));
    }

    public static BACnetApplicationTagUnsignedInteger createBACnetApplicationTagUnsignedInteger(long value) {
        Pair<Long, BACnetTagPayloadUnsignedInteger> lengthPayload = createUnsignedPayload(value);
        BACnetTagHeader header = createBACnetTagHeaderBalanced(false, BACnetDataType.UNSIGNED_INTEGER.getValue(), lengthPayload.getLeft());
        return new BACnetApplicationTagUnsignedInteger(header, lengthPayload.getRight());
    }

    public static BACnetContextTagUnsignedInteger createBACnetContextTagUnsignedInteger(byte tagNumber, long value) {
        Pair<Long, BACnetTagPayloadUnsignedInteger> lengthPayload = createUnsignedPayload(value);
        BACnetTagHeader header = createBACnetTagHeaderBalanced(true, tagNumber, lengthPayload.getLeft());
        return new BACnetContextTagUnsignedInteger(header, lengthPayload.getRight());
    }

    public static Pair<Long, BACnetTagPayloadUnsignedInteger> createUnsignedPayload(long value) {
        long length;
        Short valueUint8 = null;
        Integer valueUint16 = null;
        Integer valueUint24 = null;
        Long valueUint32 = null;
        Long valueUint40 = null;
        Long valueUint48 = null;
        Long valueUint56 = null;
        BigInteger valueUint64 = null;
        if (value < 0x100) {
            length = 1;
            valueUint8 = (short) value;
        } else if (value < 0x10000) {
            length = 2;
            valueUint16 = (int) value;
        } else if (value < 0x1000000) {
            length = 3;
            valueUint24 = (int) value;
        } else {
            length = 4;
            valueUint32 = value;
        }
        BACnetTagPayloadUnsignedInteger payload = new BACnetTagPayloadUnsignedInteger(length, valueUint8, valueUint16, valueUint24, valueUint32, valueUint40, valueUint48, valueUint56, valueUint64);
        return Pair.of(length, payload);
    }

    public static BACnetApplicationTagSignedInteger createBACnetApplicationTagSignedInteger(long value) {
        Pair<Long, BACnetTagPayloadSignedInteger> lengthPayload = createSignedPayload(value);
        BACnetTagHeader header = createBACnetTagHeaderBalanced(false, BACnetDataType.SIGNED_INTEGER.getValue(), lengthPayload.getLeft());
        return new BACnetApplicationTagSignedInteger(header, lengthPayload.getRight());
    }

    public static BACnetContextTagSignedInteger createBACnetContextTagSignedInteger(short tagNumber, long value) {
        Pair<Long, BACnetTagPayloadSignedInteger> lengthPayload = createSignedPayload(value);
        BACnetTagHeader header = createBACnetTagHeaderBalanced(false, (byte) tagNumber, lengthPayload.getLeft());
        return new BACnetContextTagSignedInteger(header, lengthPayload.getRight());
    }

    public static Pair<Long, BACnetTagPayloadSignedInteger> createSignedPayload(long value) {
        long length;
        Byte valueInt8 = null;
        Short valueInt16 = null;
        Integer valueInt24 = null;
        Integer valueInt32 = null;
        if (value < 0x100) {
            length = 1;
            valueInt8 = (byte) value;
        } else if (value < 0x10000) {
            length = 2;
            valueInt16 = (short) value;
        } else if (value < 0x1000000) {
            length = 3;
            valueInt24 = (int) value;
        } else {
            length = 4;
            valueInt32 = (int) value;
        }
        BACnetTagPayloadSignedInteger payload = new BACnetTagPayloadSignedInteger(length, valueInt8, valueInt16, valueInt24, valueInt32, null, null, null, null);
        return Pair.of(length, payload);
    }

    public static BACnetApplicationTagReal createBACnetApplicationTagReal(float value) {
        BACnetTagHeader header = createBACnetTagHeaderBalanced(false, BACnetDataType.REAL.getValue(), 4);
        return new BACnetApplicationTagReal(header, new BACnetTagPayloadReal(value));
    }

    public static BACnetContextTagReal createBACnetContextTagReal(byte tagNumber, float value) {
        BACnetTagHeader header = createBACnetTagHeaderBalanced(true, tagNumber, 4);
        return new BACnetContextTagReal(header, new BACnetTagPayloadReal(value));
    }

    public static BACnetApplicationTagDouble createBACnetApplicationTagDouble(double value) {
        BACnetTagHeader header = createBACnetTagHeaderBalanced(false, BACnetDataType.DOUBLE.getValue(), 8);
        return new BACnetApplicationTagDouble(header, new BACnetTagPayloadDouble(value));
    }

    public static BACnetContextTagDouble createBACnetContextTagDouble(byte tagNumber, double value) {
        BACnetTagHeader header = createBACnetTagHeaderBalanced(true, tagNumber, 8);
        return new BACnetContextTagDouble(header, new BACnetTagPayloadDouble(value));
    }

    public static BACnetApplicationTagOctetString createBACnetApplicationTagOctetString(byte[] octets) {
        BACnetTagHeader header = createBACnetTagHeaderBalanced(false, BACnetDataType.OCTET_STRING.getValue(), octets.length + 1);
        return new BACnetApplicationTagOctetString(header, new BACnetTagPayloadOctetString(octets));
    }

    public static BACnetContextTagOctetString createBACnetContextTagOctetString(byte tagNumber, byte[] octets) {
        BACnetTagHeader header = createBACnetTagHeaderBalanced(true, tagNumber, octets.length + 1);
        return new BACnetContextTagOctetString(header, new BACnetTagPayloadOctetString(octets));
    }

    public static BACnetApplicationTagCharacterString createBACnetApplicationTagCharacterString(BACnetCharacterEncoding baCnetCharacterEncoding, String value) {
        BACnetTagHeader header = createBACnetTagHeaderBalanced(false, BACnetDataType.CHARACTER_STRING.getValue(), value.length() + 1);
        return new BACnetApplicationTagCharacterString(header, new BACnetTagPayloadCharacterString((long) value.length() + 1, baCnetCharacterEncoding, value));
    }

    public static BACnetContextTagCharacterString createBACnetContextTagCharacterString(byte tagNumber, BACnetCharacterEncoding baCnetCharacterEncoding, String value) {
        BACnetTagHeader header = createBACnetTagHeaderBalanced(true, tagNumber, value.length() + 1);
        return new BACnetContextTagCharacterString(header, new BACnetTagPayloadCharacterString((long) (value.length() + 1), baCnetCharacterEncoding, value));
    }

    public static BACnetApplicationTagBitString createBACnetApplicationTagBitString(List<Boolean> value) {
        long numberOfBytesNeeded = (value.size() + 7) / 8;
        short unusedBits = (short) (8 - (value.size() % 8));
        if (unusedBits == 8) {
            unusedBits = 0;
        }
        BACnetTagHeader header = createBACnetTagHeaderBalanced(false, BACnetDataType.BIT_STRING.getValue(), numberOfBytesNeeded + 1);
        return new BACnetApplicationTagBitString(header, new BACnetTagPayloadBitString(unusedBits, value, new ArrayList<>(unusedBits)));
    }

    public static BACnetContextTagBitString createBACnetContextTagBitString(byte tagNumber, List<Boolean> value) {
        long numberOfBytesNeeded = (value.size() + 7) / 8;
        short unusedBits = (short) (8 - (value.size() % 8));
        if (unusedBits == 8) {
            unusedBits = 0;
        }
        BACnetTagHeader header = createBACnetTagHeaderBalanced(true, tagNumber, numberOfBytesNeeded + 1);
        return new BACnetContextTagBitString(header, new BACnetTagPayloadBitString(unusedBits, value, new ArrayList<>(unusedBits)));
    }

    public static BACnetApplicationTagEnumerated createBACnetApplicationTagEnumerated(long value) {
        Pair<Long, BACnetTagPayloadEnumerated> lengthPayload = CreateEnumeratedPayload(value);
        BACnetTagHeader header = createBACnetTagHeaderBalanced(false, BACnetDataType.ENUMERATED.getValue(), lengthPayload.getLeft());
        return new BACnetApplicationTagEnumerated(header, lengthPayload.getRight());
    }

    public static BACnetContextTagEnumerated createBACnetContextTagEnumerated(byte tagNumber, long value) {
        Pair<Long, BACnetTagPayloadEnumerated> lengthPayload = CreateEnumeratedPayload(value);
        BACnetTagHeader header = createBACnetTagHeaderBalanced(true, tagNumber, lengthPayload.getLeft());
        return new BACnetContextTagEnumerated(header, lengthPayload.getRight());
    }

    public static Pair<Long, BACnetTagPayloadEnumerated> CreateEnumeratedPayload(long value) {
        long length = requiredLength(value);
        byte[] data = writeVarUint(value);
        BACnetTagPayloadEnumerated payload = new BACnetTagPayloadEnumerated(data);
        return Pair.of(length, payload);
    }

    public static BACnetApplicationTagDate createBACnetApplicationTagDate(int year, short month, short dayOfMonth, short dayOfWeek) {
        BACnetTagHeader header = createBACnetTagHeaderBalanced(false, BACnetDataType.DATE.getValue(), 4);
        short yearMinus1900 = (short) (year - 1900);
        if (year == 0xFF) {
            yearMinus1900 = 0xFF;
        }
        return new BACnetApplicationTagDate(header, new BACnetTagPayloadDate(yearMinus1900, month, dayOfMonth, dayOfWeek));
    }

    public static BACnetContextTagDate createBACnetContextTagDate(byte tagNumber, int year, short month, short dayOfMonth, short dayOfWeek) {
        BACnetTagHeader header = createBACnetTagHeaderBalanced(true, tagNumber, 4);
        short yearMinus1900 = (short) (year - 1900);
        if (year == 0xFF) {
            yearMinus1900 = 0xFF;
        }
        return new BACnetContextTagDate(header, new BACnetTagPayloadDate(yearMinus1900, month, dayOfMonth, dayOfWeek));
    }

    public static BACnetApplicationTagTime createBACnetApplicationTagTime(short hour, short minute, short second, short fractional) {
        BACnetTagHeader header = createBACnetTagHeaderBalanced(false, BACnetDataType.TIME.getValue(), 4);
        return new BACnetApplicationTagTime(header, new BACnetTagPayloadTime(hour, minute, second, fractional));
    }

    public static BACnetContextTagTime createBACnetContextTagTime(byte tagNumber, short hour, short minute, short second, short fractional) {
        BACnetTagHeader header = createBACnetTagHeaderBalanced(true, tagNumber, 4);
        return new BACnetContextTagTime(header, new BACnetTagPayloadTime(hour, minute, second, fractional));
    }

    private static long requiredLength(long value) {
        long length;
        if (value < 0x100) length = 1;
        else if (value < 0x10000) length = 2;
        else if (value < 0x1000000) length = 3;
        else length = 4;
        return length;
    }

}
