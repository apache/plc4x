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

import org.apache.plc4x.java.streampipes.shared.source.knxnetip.Constants;
import org.streampipes.model.DataProcessorType;
import org.streampipes.model.graph.DataProcessorDescription;
import org.streampipes.model.graph.DataProcessorInvocation;
import org.streampipes.model.schema.PropertyScope;
import org.streampipes.sdk.builder.PrimitivePropertyBuilder;
import org.streampipes.sdk.builder.ProcessingElementBuilder;
import org.streampipes.sdk.builder.StreamRequirementsBuilder;
import org.streampipes.sdk.extractor.ProcessingElementParameterExtractor;
import org.streampipes.sdk.helpers.EpRequirements;
import org.streampipes.sdk.helpers.Labels;
import org.streampipes.sdk.helpers.Locales;
import org.streampipes.sdk.helpers.OutputStrategies;
import org.streampipes.sdk.utils.Assets;
import org.streampipes.sdk.utils.Datatypes;
import org.streampipes.wrapper.standalone.ConfiguredEventProcessor;
import org.streampipes.wrapper.standalone.declarer.StandaloneEventProcessingDeclarer;

import java.io.IOException;

public class Ets5DataEnrichmentController extends StandaloneEventProcessingDeclarer<Ets5DataEnrichmentParameters> {

    public static final String ID = "org.apache.plc4x.streampipes.processors.enrich.knxnetip.ets5";

    private static final String PROJECT_FILE_KEY = "project-fIle";
    private static final String DESTINATION_ID_MAPPING = "destination-id-mapping";
    private static final String PAYLOAD_ID_MAPPING = "payload-id-mapping";

    public static final String MAPPING_FIELD_DECODED_GROUP_ADDRESS = "decodedGroupAddress";
    public static final String MAPPING_FIELD_TYPE = "type";
    public static final String MAPPING_FIELD_LOCATION = "location";
    public static final String MAPPING_FIELD_FUNCTION = "function";
    public static final String MAPPING_FIELD_MEANING = "meaning";
    public static final String MAPPING_FIELD_DECODED_PROPERTY_VALUE = "decodedPropertyValue";

    @Override
    public DataProcessorDescription declareModel() {
        return ProcessingElementBuilder
            .create(ID, "ETS5", "Processor that interprets a data stream from a KXNnet/IP Datasource according to the settings in the ETS5 'knxproj' file")
            .category(DataProcessorType.ENRICH)
            //.withAssets(Assets.DOCUMENTATION, Assets.ICON)
            .withLocales(Locales.EN)
            .requiredStream(StreamRequirementsBuilder
                .create()
                .requiredPropertyWithUnaryMapping(EpRequirements.domainPropertyReq(Constants.KNXNET_ID_DESTINATION_ADDRESS), Labels.withId(DESTINATION_ID_MAPPING), PropertyScope.NONE)
                .requiredPropertyWithUnaryMapping(EpRequirements.domainPropertyReq(Constants.KNXNET_ID_PAYLOAD), Labels.withId(PAYLOAD_ID_MAPPING), PropertyScope.NONE)
                .build())
            .outputStrategy(OutputStrategies.append(
                PrimitivePropertyBuilder.create(Datatypes.String, MAPPING_FIELD_DECODED_GROUP_ADDRESS).build(),
                PrimitivePropertyBuilder.create(Datatypes.String, MAPPING_FIELD_TYPE).build(),
                PrimitivePropertyBuilder.create(Datatypes.String, MAPPING_FIELD_LOCATION).build(),
                PrimitivePropertyBuilder.create(Datatypes.String, MAPPING_FIELD_FUNCTION).build(),
                PrimitivePropertyBuilder.create(Datatypes.String, MAPPING_FIELD_MEANING).build(),
                PrimitivePropertyBuilder.create(Datatypes.String, MAPPING_FIELD_DECODED_PROPERTY_VALUE).build()
            ))
            .requiredFile(Labels.withId(PROJECT_FILE_KEY))
            .build();
    }

    @Override
    public ConfiguredEventProcessor<Ets5DataEnrichmentParameters> onInvocation(DataProcessorInvocation graph, ProcessingElementParameterExtractor extractor) {
        String destinationIdFieldName = extractor.mappingPropertyValue(DESTINATION_ID_MAPPING);
        String payloadIdFieldName = extractor.mappingPropertyValue(PAYLOAD_ID_MAPPING);
        String fileContents = null;
        try {
            fileContents = extractor.fileContentsAsString(PROJECT_FILE_KEY);
        } catch (IOException e) {
            logger.error("Error accessing the " + PROJECT_FILE_KEY + " contents", e);
        }
        Ets5DataEnrichmentParameters params = new Ets5DataEnrichmentParameters(graph, destinationIdFieldName, payloadIdFieldName);
        return new ConfiguredEventProcessor<>(params, Ets5DataEnrichment::new);
    }

}
