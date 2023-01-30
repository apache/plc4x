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
package org.apache.plc4x.java.firmata.readwrite.context;

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.firmata.readwrite.tag.FirmataTagAnalog;
import org.apache.plc4x.java.firmata.readwrite.tag.FirmataTagDigital;
import org.apache.plc4x.java.spi.model.DefaultPlcSubscriptionTag;
import org.apache.plc4x.java.spi.values.PlcList;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.firmata.readwrite.*;
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

        for (String tagName : writeRequest.getTagNames()) {
            if (!(writeRequest.getTag(tagName) instanceof FirmataTagDigital)) {
                throw new PlcRuntimeException("Writing only supported for digital pins");
            }

            FirmataTagDigital digitalTag = (FirmataTagDigital) writeRequest.getTag(tagName);
            final PlcValue plcValue = writeRequest.getPlcValue(tagName);
            if ((digitalTag.getNumberOfElements() > 1) && plcValue.isList()) {
                final PlcList plcList = (PlcList) plcValue;
                if (plcList.getList().size() != digitalTag.getNumberOfElements()) {
                    throw new PlcRuntimeException(
                        "Required " + digitalTag.getNumberOfElements() + " but got " + plcList.getList().size());
                }
            }

            for (int i = 0; i < digitalTag.getNumberOfElements(); i++) {
                int pin = digitalTag.getAddress() + i;
                if (!digitalPins.containsKey(pin)) {
                    digitalPins.put(pin, PinMode.PinModeOutput);
                    messages.add(
                        new FirmataMessageCommand(
                            new FirmataCommandSetPinMode((byte) pin, PinMode.PinModeOutput)
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
                        new FirmataCommandSetDigitalPinValue((short) pin, plcValue.getIndex(i).getBoolean())
                    )
                );
            }
        }
        return messages;
    }

    public List<FirmataMessage> processSubscriptionRequest(PlcSubscriptionRequest subscriptionRequest) {
        // Convert the request into maps of bit sets.
        Map<Integer, PinMode> requestDigitalTagPinModes = new HashMap<>();
        Map<Integer, PinMode> requestAnalogTagPinModes = new HashMap<>();
        for (String tagName : subscriptionRequest.getTagNames()) {
            final PlcTag tag = subscriptionRequest.getTag(tagName);
            DefaultPlcSubscriptionTag subscriptionTag = (DefaultPlcSubscriptionTag) tag;
            if (subscriptionTag.getTag() instanceof FirmataTagDigital) {
                FirmataTagDigital tagDigital = (FirmataTagDigital) subscriptionTag.getTag();
                PinMode tagPinMode = (tagDigital.getPinMode() != null) ?
                    tagDigital.getPinMode() : PinMode.PinModeInput;
                if (!(tagPinMode.equals(PinMode.PinModeInput) || tagPinMode.equals(PinMode.PinModePullup))) {
                    throw new PlcInvalidTagException("Subscription tag must be of type 'INPUT' (default) or 'PULLUP'");
                }
                for (int pin = tagDigital.getAddress(); pin < tagDigital.getAddress() + tagDigital.getNumberOfElements(); pin++) {
                    requestDigitalTagPinModes.put(pin, tagPinMode);
                }
            } else if (subscriptionTag.getTag() instanceof FirmataTagAnalog) {
                FirmataTagAnalog tagAnalog = (FirmataTagAnalog) subscriptionTag.getTag();
                for (int pin = tagAnalog.getAddress(); pin < tagAnalog.getAddress() + tagAnalog.getNumberOfElements(); pin++) {
                    requestAnalogTagPinModes.put(pin, PinMode.PinModeInput);
                }
            } else {
                throw new PlcRuntimeException("Unsupported tag type " + tag.getClass().getSimpleName());
            }
        }

        // If a requested digital pin is already subscribed, blank this out
        for (Map.Entry<Integer, PinMode> entry : requestDigitalTagPinModes.entrySet()) {
            int pin = entry.getKey();
            PinMode pinMode = entry.getValue();
            if (digitalPins.containsKey(pin)) {
                if (!digitalPins.get(pin).equals(pinMode)) {
                    throw new PlcInvalidTagException(String.format(
                        "Error setting digital pin to mode %s, pin is already set to mode %s",
                        pinMode.toString(), digitalPins.get(pin).toString()));
                } else {
                    requestDigitalTagPinModes.remove(pin);
                }
            }
        }
        // If a requested analog pin is already subscribed, blank this out
        for (Map.Entry<Integer, PinMode> entry : requestAnalogTagPinModes.entrySet()) {
            int pin = entry.getKey();
            if (analogPins.containsKey(pin)) {
                requestAnalogTagPinModes.remove(pin);
            }
        }

        // Remember the subscription itself.
        subscriptions.add(subscriptionRequest);

        // Create a list of messages that need to be sent to achieve the desired subscriptions.
        List<FirmataMessage> messages = new LinkedList<>();
        for (Map.Entry<Integer, PinMode> entry : requestDigitalTagPinModes.entrySet()) {
            int pin = entry.getKey();
            PinMode pinMode = entry.getValue();
            // Digital pins can be input and output, so first we have to set it to "input"
            messages.add(new FirmataMessageCommand(new FirmataCommandSetPinMode((byte) pin, pinMode)));
            // And then tell the remote to send change of state information.
            messages.add(new FirmataMessageSubscribeDigitalPinValue((byte) pin, true));
        }
        for (Map.Entry<Integer, PinMode> entry : requestAnalogTagPinModes.entrySet()) {
            int pin = entry.getKey();
            // Tell the remote to send change of state information for this analog pin.
            messages.add(new FirmataMessageSubscribeAnalogPinValue((byte) pin, true));
        }

        return messages;
    }

}
