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
package org.apache.plc4x.java.firmata.readwrite.context;

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.firmata.readwrite.*;
import org.apache.plc4x.java.firmata.readwrite.field.FirmataFieldAnalog;
import org.apache.plc4x.java.firmata.readwrite.field.FirmataFieldDigital;
import org.apache.plc4x.java.firmata.readwrite.types.PinMode;
import org.apache.plc4x.java.spi.context.DriverContext;

import java.util.*;

public class FirmataDriverContext implements DriverContext {

    private final List<PlcSubscriptionRequest> subscriptions;
    private final BitSet digitalInputPins;
    private final BitSet digitalOutputPins;
    private final BitSet analogInputPins;

    public FirmataDriverContext() {
        subscriptions = new LinkedList<>();
        digitalInputPins = new BitSet();
        digitalOutputPins = new BitSet();
        analogInputPins = new BitSet();
    }

    public List<FirmataMessage> processSubscriptionRequest(PlcSubscriptionRequest subscriptionRequest) {
        // Convert the request into maps of bit sets.
        BitSet requestDigitalFields = new BitSet();
        BitSet requestAnalogFields = new BitSet();
        for (String fieldName : subscriptionRequest.getFieldNames()) {
            final PlcField field = subscriptionRequest.getField(fieldName);
            if(field instanceof FirmataFieldDigital) {
                FirmataFieldDigital fieldDigital = (FirmataFieldDigital) field;
                for(int pin = fieldDigital.getAddress(); pin < fieldDigital.getAddress() + fieldDigital.getQuantity(); pin++) {
                    requestDigitalFields.set(pin, true);
                }
            } else if(field instanceof FirmataFieldAnalog) {
                FirmataFieldAnalog fieldAnalog = (FirmataFieldAnalog) field;
                for(int pin = fieldAnalog.getAddress(); pin < fieldAnalog.getAddress() + fieldAnalog.getQuantity(); pin++) {
                    requestAnalogFields.set(pin, true);
                }
            } else {
                throw new PlcRuntimeException("Unsupported field type " + field.getClass().getSimpleName());
            }
        }

        // If a requested digital pin is already subscribed, blank this out
        for(int pin = 0; pin < requestDigitalFields.length(); pin++) {
            if(requestDigitalFields.get(pin) && digitalInputPins.get(pin)) {
                requestDigitalFields.set(pin, false);
            }
        }
        // If a requested analog pin is already subscribed, blank this out
        for(int pin = 0; pin < requestAnalogFields.length(); pin++) {
            if(requestAnalogFields.get(pin) && !analogInputPins.get(pin)) {
                requestAnalogFields.set(pin, false);
            }
        }

        // Check if any of the remaining requested digital pins is already set to `output`
        for(int pin = 0; pin < requestDigitalFields.length(); pin++) {
            if(requestDigitalFields.get(pin) && digitalOutputPins.get(pin)) {
                throw new PlcRuntimeException("Pin " + pin + " already configured as output pin.");
            }
        }

        // Remember the subscription itself.
        subscriptions.add(subscriptionRequest);

        // Create a list of messages that need to be sent to achieve the desired subscriptions.
        List<FirmataMessage> messages = new LinkedList<>();
        for(int pin = 0; pin < requestDigitalFields.length(); pin++) {
            if(requestDigitalFields.get(pin)) {
                // Digital pins can be input and output, so first we have to set it to "input"
                messages.add(new FirmataMessageCommand(new FirmataCommandSetPinMode((byte) pin, PinMode.PinModeInput)));
                // And then tell the remote to send change of state information.
                messages.add(new FirmataMessageSubscribeDigitalPinValue((byte) pin, true));
            }
        }
        for(int pin = 0; pin < requestAnalogFields.length(); pin++) {
            if(requestAnalogFields.get(pin)) {
                // Tell the remote to send change of state information for this analog pin.
                messages.add(new FirmataMessageSubscribeAnalogPinValue((byte) pin, true));
            }
        }

        return messages;
    }

}
