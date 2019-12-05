/*
 * Copyright 2019 FZI Forschungszentrum Informatik
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package org.apache.plc4x.java.streampipes.bacnetip.config;

import org.streampipes.config.SpConfig;
import org.streampipes.connect.init.Config;

public enum ConnectWorkerConfig {

    INSTANCE;

    private SpConfig config;

    ConnectWorkerConfig() {
        String name = "bacnetip-connect-worker-main";
        config = SpConfig.getSpConfig(name);

        config.register(ConfigKeys.KAFKA_HOST, "kafka", "Hostname for backend service for kafka");
        config.register(ConfigKeys.KAFKA_PORT, 9092, "Port for backend service for kafka");

        config.register(ConfigKeys.CONNECT_CONTAINER_WORKER_PORT, Config.WORKER_PORT, "The port of the connect container");
        config.register(ConfigKeys.CONNECT_CONTAINER_WORKER_HOST, name, "The hostname of the connect container");

        config.register(ConfigKeys.CONNECT_CONTAINER_MASTER_PORT, Config.MASTER_PORT, "The port of the connect container");
        config.register(ConfigKeys.CONNECT_CONTAINER_MASTER_HOST, Config.MASTER_HOST, "The hostname of the connect container");
    }

    public String getConnectContainerWorkerUrl() {
        return "http://" + config.getString(ConfigKeys.CONNECT_CONTAINER_WORKER_HOST) + ":" + config.getInteger(ConfigKeys.CONNECT_CONTAINER_WORKER_PORT) + "/";
    }

    public String getConnectContainerMasterUrl() {
        return "http://" + config.getString(ConfigKeys.CONNECT_CONTAINER_MASTER_HOST) + ":" + config.getInteger(ConfigKeys.CONNECT_CONTAINER_MASTER_PORT) + "/";
    }

    public String getKafkaHost() {
        return config.getString(ConfigKeys.KAFKA_HOST);
    }

    public void setKafkaHost(String s) {
        config.setString(ConfigKeys.KAFKA_HOST, s);
    }

    public int getKafkaPort() {
        return config.getInteger(ConfigKeys.KAFKA_PORT);
    }

    public String getKafkaUrl() {
        return getKafkaHost() + ":" + getKafkaPort();
    }

    public String getConnectContainerWorkerHost() {
        return config.getString(ConfigKeys.CONNECT_CONTAINER_WORKER_HOST);
    }

    public Integer getConnectContainerWorkerPort() {
        return config.getInteger(ConfigKeys.CONNECT_CONTAINER_WORKER_PORT);
    }

    public String getConnectContainerMasterHost() {
        return config.getString(ConfigKeys.CONNECT_CONTAINER_MASTER_HOST);
    }

    public Integer getConnectContainerMasterPort() {
        return config.getInteger(ConfigKeys.CONNECT_CONTAINER_MASTER_PORT);
    }

}
