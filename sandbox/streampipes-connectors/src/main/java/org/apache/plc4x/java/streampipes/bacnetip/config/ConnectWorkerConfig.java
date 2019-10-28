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
package org.apache.plc4x.java.streampipes.bacnetip.config;

import org.streampipes.config.SpConfig;
import org.streampipes.connect.init.Config;

public class ConnectWorkerConfig {

    private static ConnectWorkerConfig instance = new ConnectWorkerConfig();
    public static ConnectWorkerConfig getInstance() {
        return instance;
    }

    private SpConfig config;

    private ConnectWorkerConfig() {
        String name = "bacnetip-connect-worker-main";

        config = SpConfig.getSpConfig(name);

        config.register(ConfigKeys.SP_KAFKA_HOST, "kafka", "");
        config.register(ConfigKeys.SP_KAFKA_PORT, 9092, "");

        config.register(ConfigKeys.SP_CONNECT_CONTAINER_MASTER_HOST, name, "");
        config.register(ConfigKeys.SP_CONNECT_CONTAINER_MASTER_PORT, Config.WORKER_PORT, "");

        config.register(ConfigKeys.SP_CONNECT_CONTAINER_WORKER_HOST, Config.MASTER_HOST, "");
        config.register(ConfigKeys.SP_CONNECT_CONTAINER_WORKER_PORT, Config.MASTER_PORT, "");

    }

    public String getConnectContainerMasterUrl() {
        return "http://" + config.getString(ConfigKeys.SP_CONNECT_CONTAINER_MASTER_HOST) + ":" +
            config.getInteger(ConfigKeys.SP_CONNECT_CONTAINER_MASTER_PORT) + "/";
    }

    public String getConnectContainerWorkerUrl() {
        return "http://" + config.getString(ConfigKeys.SP_CONNECT_CONTAINER_WORKER_HOST) + ":" +
            config.getInteger(ConfigKeys.SP_CONNECT_CONTAINER_WORKER_PORT) + "/";
    }

    public int getConnectContainerWorkerPort() {
        return config.getInteger(ConfigKeys.SP_CONNECT_CONTAINER_WORKER_PORT);
    }

}
