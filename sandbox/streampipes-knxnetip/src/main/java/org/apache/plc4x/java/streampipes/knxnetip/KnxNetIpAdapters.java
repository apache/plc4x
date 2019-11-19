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
package org.apache.plc4x.java.streampipes.knxnetip;

import org.apache.plc4x.java.streampipes.knxnetip.config.BacNetIpAdapterInit;
import org.apache.plc4x.java.streampipes.knxnetip.processors.ets5.Ets5DataEnrichmentController;
import org.apache.plc4x.java.streampipes.knxnetip.source.KnxNetIpAdapter;
import org.apache.plc4x.java.streampipes.knxnetip.config.ConnectWorkerConfig;
import org.streampipes.connect.init.AdapterDeclarerSingleton;
import org.streampipes.container.init.DeclarersSingleton;

public class KnxNetIpAdapters {

    public static void main(String[] args) {
        // Declare the data-sources.
        AdapterDeclarerSingleton
            .getInstance()
            .add(new KnxNetIpAdapter());

        // Declare the processors.
        /*DeclarersSingleton
            .getInstance()
            .add(new Ets5DataEnrichmentController());*/

        // Initialize the adapter itself.
        String workerUrl = ConnectWorkerConfig.INSTANCE.getConnectContainerWorkerUrl();
        String masterUrl = ConnectWorkerConfig.INSTANCE.getConnectContainerMasterUrl();
        Integer workerPort = ConnectWorkerConfig.INSTANCE.getConnectContainerWorkerPort();
        new BacNetIpAdapterInit().init(workerUrl, masterUrl, workerPort);
    }

}
