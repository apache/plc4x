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
package org.apache.plc4x.java.firmata.readwrite.context;

import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.model.PlcSubscriptionField;
import org.apache.plc4x.java.spi.model.DefaultPlcSubscriptionField;
import org.apache.plc4x.java.spi.values.PlcList;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.firmata.readwrite.*;
import org.apache.plc4x.java.firmata.readwrite.field.FirmataFieldAnalog;
import org.apache.plc4x.java.firmata.readwrite.field.FirmataFieldDigital;
import org.apache.plc4x.java.spi.context.DriverContext;

import java.util.*;

public class FirmataDriverContext implements DriverContext {

    private final List<PlcSubscriptionRequest> subscriptions;
    private final Map<Integer, PinMode> digitalPins;
    private final Map<Integer, PinMode> analogPins;

    public FirmataDriverContext() {
        subscriptions = new LinkedList<>();
        digitalPins = new HashMap<>();
        analogPins = new HashMap<>();
    }

    public List<FirmataMessage> processWriteRequest(PlcWriteRequest writeRequest) {
        List<FirmataMessage> messages = new LinkedList<>();

        for (String fieldName : writeRequest.getFieldNames()) {
            if (!(writeRequest.getField(fieldName) instanceof FirmataFieldDigital)) {
                throw new PlcRuntimeException("Writing only supported for digital pins");
            }

            FirmataFieldDigital digitalField = (FirmataFieldDigital) writeRequest.getField(fieldName);
            final PlcValue plcValue = writeRequest.getPlcValue(fieldName);
            if ((digitalField.getNumberOfElements() > 1) && plcValue.isList()) {
                final PlcList plcList = (PlcList) plcValue;
                if (plcList.getList().size() != digitalField.getNumberOfElements()) {
                    throw new PlcRuntimeException(
                        "Required " + digitalField.getNumberOfElements() + " but got " + plcList.getList().size());
                }
            }

            for (int i = 0; i < digitalField.getNumberOfElements(); i++) {
                int pin = digitalField.getAddress() + i;
                if (!digitalPins.containsKey(pin)) {
                    digitalPins.put(pin, PinMode.PinModeOutput);
                    messages.add(
                        new FirmataMessageCommand(
                            new FirmataCommandSetPinMode((byte) pin, PinMode.PinModeOutput, false), false
                        )
                    );
                }
                // Check that a requested output pin is currently not configured as 'input'.
                else if (!digitalPins.get(pin).equals(PinMode.PinModeOutput)) {
                    throw new PlcRuntimeException(
                        "Pin " + pin + " already configured as " + digitalPins.get(pin).name());
                }

                messages.add(
                    new FirmataMessageCommand(
                        new FirmataCommandSetDigitalPinValue((short) pin, plcValue.getIndex(i).getBoolean(), false),
                        false
                    )
                );
            }
        }
        return messages;
    }

    public List<FirmataMessage> processSubscriptionRequest(PlcSubscriptionRequest subscriptionRequest) {
        // Convert the request into maps of bit sets.
        Map<Integer, PinMode> requestDigitalFieldPinModes = new HashMap<>();
        Map<Integer, PinMode> requestAnalogFieldPinModes = new HashMap<>();
        for (String fieldName : subscriptionRequest.getFieldNames()) {
            final PlcField field = subscriptionRequest.getField(fieldName);
            DefaultPlcSubscriptionField subscriptionField = (DefaultPlcSubscriptionField) field;
            if (subscriptionField.getPlcField() instanceof FirmataFieldDigital) {
                FirmataFieldDigital fieldDigital = (FirmataFieldDigital) subscriptionField.getPlcField();
                PinMode fieldPinMode = (fieldDigital.getPinMode() != null) ?
                    fieldDigital.getPinMode() : PinMode.PinModeInput;
                if (!(fieldPinMode.equals(PinMode.PinModeInput) || fieldPinMode.equals(PinMode.PinModePullup))) {
                    throw new PlcInvalidFieldException("Subscription field must be of type 'INPUT' (default) or 'PULLUP'");
                }
                for (int pin = fieldDigital.getAddress(); pin < fieldDigital.getAddress() + fieldDigital.getNumberOfElements(); pin++) {
                    requestDigitalFieldPinModes.put(pin, fieldPinMode);
                }
            } else if (subscriptionField.getPlcField() instanceof FirmataFieldAnalog) {
                FirmataFieldAnalog fieldAnalog = (FirmataFieldAnalog) subscriptionField.getPlcField();
                for (int pin = fieldAnalog.getAddress(); pin < fieldAnalog.getAddress() + fieldAnalog.getNumberOfElements(); pin++) {
                    requestAnalogFieldPinModes.put(pin, PinMode.PinModeInput);
                }
            } else {
                throw new PlcRuntimeException("Unsupported field type " + field.getClass().getSimpleName());
            }
        }

        // If a requested digital pin is already subscribed, blank this out
        for (Map.Entry<Integer, PinMode> entry : requestDigitalFieldPinModes.entrySet()) {
            int pin = entry.getKey();
            PinMode pinMode = entry.getValue();
            if (digitalPins.containsKey(pin)) {
                if (!digitalPins.get(pin).equals(pinMode)) {
                    throw new PlcInvalidFieldException(String.format(
                        "Error setting digital pin to mode %s, pin is already set to mode %s",
                        pinMode.toString(), digitalPins.get(pin).toString()));
                } else {
                    requestDigitalFieldPinModes.remove(pin);
                }
            }
        }
        // If a requested analog pin is already subscribed, blank this out
        for (Map.Entry<Integer, PinMode> entry : requestAnalogFieldPinModes.entrySet()) {
            int pin = entry.getKey();
            if (analogPins.containsKey(pin)) {
                requestAnalogFieldPinModes.remove(pin);
            }
        }

        // Remember the subscription itself.
        subscriptions.add(subscriptionRequest);

        // Create a list of messages that need to be sent to achieve the desired subscriptions.
        List<FirmataMessage> messages = new LinkedList<>();
        for (Map.Entry<Integer, PinMode> entry : requestDigitalFieldPinModes.entrySet()) {
            int pin = entry.getKey();
            PinMode pinMode = entry.getValue();
            // Digital pins can be input and output, so first we have to set it to "input"
            messages.add(new FirmataMessageCommand(new FirmataCommandSetPinMode((byte) pin, pinMode, false), false));
            // And then tell the remote to send change of state information.
            messages.add(new FirmataMessageSubscribeDigitalPinValue((byte) pin, true, false));
        }
        for (Map.Entry<Integer, PinMode> entry : requestAnalogFieldPinModes.entrySet()) {
            int pin = entry.getKey();
            // Tell the remote to send change of state information for this analog pin.
            messages.add(new FirmataMessageSubscribeAnalogPinValue((byte) pin, true, false));
        }

        return messages;
    }

}
