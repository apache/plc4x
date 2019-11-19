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
package org.apache.plc4x.java.streampipes.knxnetip.processors.ets5;

import org.apache.plc4x.java.streampipes.knxnetip.source.KnxNetIpAdapter;
import org.streampipes.model.DataProcessorType;
import org.streampipes.model.graph.DataProcessorDescription;
import org.streampipes.model.graph.DataProcessorInvocation;
import org.streampipes.sdk.builder.ProcessingElementBuilder;
import org.streampipes.sdk.builder.StreamRequirementsBuilder;
import org.streampipes.sdk.extractor.ProcessingElementParameterExtractor;
import org.streampipes.sdk.helpers.EpRequirements;
import org.streampipes.sdk.helpers.Labels;
import org.streampipes.sdk.helpers.Locales;
import org.streampipes.sdk.helpers.OutputStrategies;
import org.streampipes.sdk.utils.Assets;
import org.streampipes.wrapper.standalone.ConfiguredEventProcessor;
import org.streampipes.wrapper.standalone.declarer.StandaloneEventProcessingDeclarer;

public class Ets5DataEnrichmentController extends StandaloneEventProcessingDeclarer<Ets5DataEnrichmentParameters> {

    public static final String ID = "http://plc4x.apache.org/streampipes/processor/ets5";

    @Override
    public DataProcessorDescription declareModel() {
        final DataProcessorDescription description = ProcessingElementBuilder
            .create(ID, "ETS5", "Processor that interprets a data stream from a KXNnet/IP Datasource according to the settings in the ETS5 'knxproj' file")
            .category(DataProcessorType.ENRICH)
            .withAssets(Assets.DOCUMENTATION, Assets.ICON)
            .withLocales(Locales.EN)
            .requiredStream(StreamRequirementsBuilder
                .create()
                .requiredProperty(EpRequirements.domainPropertyReq(KnxNetIpAdapter.ID_DESTINATION_ADDRESS))
                .requiredProperty(EpRequirements.domainPropertyReq(KnxNetIpAdapter.ID_PAYLOAD))
                .build())
            .outputStrategy(OutputStrategies.keep())
            .requiredFile(Labels.from("knxprojFile", "ETS5 Project File", "ETS5 Project File (.knxproj)"))
            .build();
        return description;
    }

    @Override
    public ConfiguredEventProcessor<Ets5DataEnrichmentParameters> onInvocation(DataProcessorInvocation dataProcessorInvocation, ProcessingElementParameterExtractor processingElementParameterExtractor) {
        return null;
    }

}
