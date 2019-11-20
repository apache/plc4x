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

import org.streampipes.commons.exceptions.SpRuntimeException;
import org.streampipes.model.runtime.Event;
import org.streampipes.wrapper.context.EventProcessorRuntimeContext;
import org.streampipes.wrapper.routing.SpOutputCollector;
import org.streampipes.wrapper.runtime.EventProcessor;

public class Ets5DataEnrichment implements EventProcessor<Ets5DataEnrichmentParameters> {
    @Override
    public void onInvocation(Ets5DataEnrichmentParameters ets5DataEnrichmentParameters, SpOutputCollector spOutputCollector, EventProcessorRuntimeContext eventProcessorRuntimeContext) throws SpRuntimeException {

    }

    @Override
    public void onEvent(Event event, SpOutputCollector spOutputCollector) throws SpRuntimeException {
        spOutputCollector.collect(event);
    }

    @Override
    public void onDetach() throws SpRuntimeException {

    }
}
