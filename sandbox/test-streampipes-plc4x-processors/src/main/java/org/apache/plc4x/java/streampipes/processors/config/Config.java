package org.apache.plc4x.java.streampipes.processors.config;

import org.streampipes.config.SpConfig;
import org.streampipes.container.model.PeConfig;

public enum Config implements PeConfig {

  INSTANCE;

  private SpConfig config;

  public final static String serverUrl;

  private final static String SERVICE_ID= "pe/org.apache.plc4x.streampipes.processors";

  Config() {
    config = SpConfig.getSpConfig("pe/org.apache.plc4x.streampipes.processors");

    config.register(ConfigKeys.HOST, "processors-plc4x", "Hostname for the pe sinks");
    config.register(ConfigKeys.PORT, 8090, "Port for the pe sinks");

    config.register(ConfigKeys.SERVICE_NAME, "PLC4X Processors", "");
  }

  static {
    serverUrl = Config.INSTANCE.getHost() + ":" + Config.INSTANCE.getPort();
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
