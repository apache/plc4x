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
package org.apache.plc4x.java.streampipes.processors;

import org.apache.plc4x.java.streampipes.processors.config.Config;
import org.apache.plc4x.java.streampipes.processors.processors.ets5.Ets5DataEnrichmentController;
import org.streampipes.container.init.DeclarersSingleton;
import org.streampipes.container.standalone.init.StandaloneModelSubmitter;
import org.streampipes.dataformat.cbor.CborDataFormatFactory;
import org.streampipes.dataformat.fst.FstDataFormatFactory;
import org.streampipes.dataformat.json.JsonDataFormatFactory;
import org.streampipes.dataformat.smile.SmileDataFormatFactory;
import org.streampipes.messaging.kafka.SpKafkaProtocolFactory;

public class Plc4xProcessorsInit extends StandaloneModelSubmitter {

    public static void main(String[] args) {
        // Declare the processors.
        DeclarersSingleton
            .getInstance()
            .add(new Ets5DataEnrichmentController());

        // Declare the data formats the
        DeclarersSingleton.getInstance().registerDataFormats(
            new JsonDataFormatFactory(),
            new CborDataFormatFactory(),
            new SmileDataFormatFactory(),
            new FstDataFormatFactory());

        // Declare which methods can be used to actually read the data
        // as well as how to transport it to the next pipeline-element.
        DeclarersSingleton.getInstance().registerProtocol(new SpKafkaProtocolFactory());

        // Initialize these modules components.
        new Plc4xProcessorsInit().init(Config.INSTANCE);
    }

}
