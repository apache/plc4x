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

package org.apache.plc4x.java.streampipes.processors.config;

import org.streampipes.config.SpConfig;
import org.streampipes.container.model.PeConfig;

public enum Config implements PeConfig {

    INSTANCE;

    public final static String serverUrl;
    private final static String SERVICE_ID = "pe/org.apache.plc4x.streampipes.processors";

    static {
        serverUrl = Config.INSTANCE.getHost() + ":" + Config.INSTANCE.getPort();
    }

    private SpConfig config;

    Config() {
        // The name this config will be saved under
        // Can be found/edited/deleted here: http://localhost:8500/ui/dc1/kv/sp/v1/
        String name = "pe/org.apache.plc4x.streampipes.processors";
        config = SpConfig.getSpConfig(name);

        // Only used, if there is no configuration available in consul.
        config.register(ConfigKeys.HOST, "processors-plc4x", "Hostname for the pe sinks");
        config.register(ConfigKeys.PORT, 8090, "Port for the pe sinks");
        config.register(ConfigKeys.SERVICE_NAME, "PLC4X Processors", "");
    }

    public String getHost() {
        return config.getString(ConfigKeys.HOST);
    }

    public int getPort() {
        return config.getInteger(ConfigKeys.PORT);
    }

    public String getKafkaHost() {
        return config.getString(ConfigKeys.KAFKA_HOST);
    }

    public int getKafkaPort() {
        return config.getInteger(ConfigKeys.KAFKA_PORT);
    }

    public String getKafkaUrl() {
        return getKafkaHost() + ":" + getKafkaPort();
    }

    public String getZookeeperHost() {
        return config.getString(ConfigKeys.ZOOKEEPER_HOST);
    }

    public int getZookeeperPort() {
        return config.getInteger(ConfigKeys.ZOOKEEPER_PORT);
    }

    @Override
    public String getId() {
        return SERVICE_ID;
    }

    @Override
    public String getName() {
        return config.getString(ConfigKeys.SERVICE_NAME);
    }

}
