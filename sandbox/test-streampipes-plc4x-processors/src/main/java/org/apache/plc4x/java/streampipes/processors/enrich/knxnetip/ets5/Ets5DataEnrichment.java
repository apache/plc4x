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
package org.apache.plc4x.java.streampipes.processors.enrich.knxnetip.ets5;

import org.apache.commons.codec.DecoderException;
import org.apache.commons.codec.binary.Hex;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.knxnetip.ets5.Ets5Parser;
import org.apache.plc4x.java.knxnetip.ets5.model.Ets5Model;
import org.apache.plc4x.java.knxnetip.ets5.model.GroupAddress;
import org.apache.plc4x.java.knxnetip.readwrite.io.KnxDatapointIO;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.ReadBufferByteBased;
import org.apache.streampipes.model.runtime.Event;
import org.apache.streampipes.wrapper.context.EventProcessorRuntimeContext;
import org.apache.streampipes.wrapper.routing.SpOutputCollector;
import org.apache.streampipes.wrapper.runtime.EventProcessor;

import java.io.File;

public class Ets5DataEnrichment implements EventProcessor<Ets5DataEnrichmentParameters> {

    private String destinationIdFieldName;
    private String payloadIdFieldName;
    // TODO: Make this dynamic.
    private static final Ets5Model ets5Model = new Ets5Parser().parse(new File("/Users/christofer.dutz/Projects/Apache/PLC4X-Documents/KNX/Stettiner Str. 13/StettinerStr-Soll-Ist-Temperatur.knxproj"));

    @Override
    public void onInvocation(Ets5DataEnrichmentParameters params, SpOutputCollector spOutputCollector,
                             EventProcessorRuntimeContext eventProcessorRuntimeContext) {
        destinationIdFieldName = params.getDestinationIdFieldName();
        payloadIdFieldName = params.getPayloadIdFieldName();
    }

    @Override
    public void onEvent(Event event, SpOutputCollector spOutputCollector) {
        try {
            // Get the raw group address data.
            String destinationFieldValue = event.getFieldBySelector(this.destinationIdFieldName).getAsPrimitive().getAsString();
            byte[] destinationGroupAddress = Hex.decodeHex(destinationFieldValue);
            // Decode the group address depending on the project settings.
            String destinationAddress = Ets5Model.parseGroupAddress(ets5Model.getGroupAddressType(), destinationGroupAddress);
            final GroupAddress groupAddress = ets5Model.getGroupAddresses().get(destinationAddress);

            // Get the raw HEX-encoded data
            String hexEncodedRawData = event.getFieldBySelector(this.payloadIdFieldName).getAsPrimitive().getAsString();
            // Convert the HEX-encoded data back to byte[]
            byte[] rawData = Hex.decodeHex(hexEncodedRawData);
            ReadBuffer rawDataReader = new ReadBufferByteBased(rawData);

            if (groupAddress != null) {
                // Decode the raw data.
                /*final PlcValue plcValue = KnxDatapointIO.staticParse(rawDataReader, groupAddress.getType().getFormatName());

                // Serialize the decoded object to json
                final String jsonDatapoint = datapoint.toString(ToStringStyle.JSON_STYLE);

                // Add the additional properties.
                event.addField(Ets5DataEnrichmentController.MAPPING_FIELD_DECODED_GROUP_ADDRESS,
                    toGroupAddressString(destinationAddress));
                event.addField(Ets5DataEnrichmentController.MAPPING_FIELD_TYPE,
                    groupAddress.getType().getName());
                event.addField(Ets5DataEnrichmentController.MAPPING_FIELD_LOCATION,
                    groupAddress.getFunction().getSpaceName());
                event.addField(Ets5DataEnrichmentController.MAPPING_FIELD_FUNCTION,
                    groupAddress.getFunction().getName());
                event.addField(Ets5DataEnrichmentController.MAPPING_FIELD_MEANING,
                    groupAddress.getName());
                event.addField(Ets5DataEnrichmentController.MAPPING_FIELD_DECODED_PROPERTY_VALUE,
                    jsonDatapoint);

                System.out.println(hexEncodedRawData + " decoded to: " + jsonDatapoint + " " + groupAddress.getType().getMainType() + " " + groupAddress.getType().getSubType());

                //Event enrichedEvent = new Event()
                spOutputCollector.collect(event);*/
            } else {
                System.out.println("Couldn't decode group address " + destinationAddress);
            }
        /*} catch (ParseException e) {
            // Driver Decoding
            e.printStackTrace();*/
        } catch (DecoderException e) {
            // Hex Decoding
            e.printStackTrace();
        }
    }

    @Override
    public void onDetach() {

    }

}
