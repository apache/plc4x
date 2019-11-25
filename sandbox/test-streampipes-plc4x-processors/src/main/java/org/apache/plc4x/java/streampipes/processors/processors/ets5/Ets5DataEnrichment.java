/*
Copyright 2019 FZI Forschungszentrum Informatik

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package org.apache.plc4x.java.streampipes.processors.processors.ets5;

import org.apache.commons.codec.DecoderException;
import org.apache.commons.codec.binary.Hex;
import org.apache.commons.lang3.builder.ToStringStyle;
import org.apache.plc4x.java.ets5.passive.*;
import org.apache.plc4x.java.ets5.passive.io.KNXGroupAddressIO;
import org.apache.plc4x.java.ets5.passive.io.KnxDatapointIO;
import org.apache.plc4x.java.knxnetip.ets5.Ets5Parser;
import org.apache.plc4x.java.knxnetip.ets5.model.Ets5Model;
import org.apache.plc4x.java.knxnetip.ets5.model.GroupAddress;
import org.apache.plc4x.java.utils.ParseException;
import org.apache.plc4x.java.utils.ReadBuffer;
import org.streampipes.model.runtime.Event;
import org.streampipes.wrapper.context.EventProcessorRuntimeContext;
import org.streampipes.wrapper.routing.SpOutputCollector;
import org.streampipes.wrapper.runtime.EventProcessor;

import java.io.File;

public class Ets5DataEnrichment implements EventProcessor<Ets5DataEnrichmentParameters> {

    private String destinationIdFieldName;
    private String payloadIdFieldName;
    // TODO: Make this dynamic.
    private static final Ets5Model ets5Model = new Ets5Parser().parse(new File("/Users/christofer.dutz/Projects/Apache/PLC4X-Documents/KNX/Stettiner Str. 13/StettinerStr-Soll-Ist-Temperatur.knxproj"));;

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
            ReadBuffer addressReadBuffer = new ReadBuffer(destinationGroupAddress);
            // Decode the group address depending on the project settings.
            KNXGroupAddress destinationAddress =
                KNXGroupAddressIO.parse(addressReadBuffer, ets5Model.getGroupAddressType());
            final GroupAddress groupAddress = ets5Model.getGroupAddresses().get(destinationAddress);

            // Get the raw HEX-encoded data
            String hexEncodedRawData = event.getFieldBySelector(this.payloadIdFieldName).getAsPrimitive().getAsString();
            // Convert the HEX-encoded data back to byte[]
            byte[] rawData = Hex.decodeHex(hexEncodedRawData);
            ReadBuffer rawDataReader = new ReadBuffer(rawData);

            if (groupAddress != null) {
                // Decode the raw data.
                final KnxDatapoint datapoint = KnxDatapointIO.parse(rawDataReader, groupAddress.getType().getMainType(), groupAddress.getType().getSubType());

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

                //Event enrichedEvent = new Event()
                spOutputCollector.collect(event);
            } else {
                System.out.println("Couldn't decode group address " + toGroupAddressString(destinationAddress));
            }
        } catch (ParseException e) {
            // Driver Decoding
            e.printStackTrace();
        } catch (DecoderException e) {
            // Hex Decoding
            e.printStackTrace();
        }
    }

    @Override
    public void onDetach() {

    }

    private String toGroupAddressString(KNXGroupAddress groupAddress) {
        if(groupAddress instanceof KNXGroupAddress3Level) {
            KNXGroupAddress3Level castedAddress = (KNXGroupAddress3Level) groupAddress;
            return castedAddress.getMainGroup() + "/" + castedAddress.getMiddleGroup() + "/" + castedAddress.getSubGroup();
        } else if(groupAddress instanceof KNXGroupAddress2Level) {
            KNXGroupAddress2Level castedAddress = (KNXGroupAddress2Level) groupAddress;
            return castedAddress.getMainGroup() + "/" + castedAddress.getSubGroup();
        } else if(groupAddress instanceof KNXGroupAddressFreeLevel) {
            KNXGroupAddressFreeLevel castedAddress = (KNXGroupAddressFreeLevel) groupAddress;
            return castedAddress.getSubGroup() + "";
        }
        throw new RuntimeException("Unknown group address type");
    }

}
