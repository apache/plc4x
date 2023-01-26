/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.examples.helloinflux;

import com.influxdb.client.InfluxDBClient;
import com.influxdb.client.InfluxDBClientFactory;
import com.influxdb.client.WriteApi;
import com.influxdb.client.domain.WritePrecision;
import com.influxdb.client.write.Point;
import org.apache.commons.configuration2.Configuration;
import org.apache.commons.configuration2.builder.fluent.Configurations;
import org.apache.commons.configuration2.ex.ConfigurationException;
import org.apache.plc4x.java.DefaultPlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.*;
import org.apache.plc4x.java.spi.messages.DefaultPlcSubscriptionEvent;
import org.apache.plc4x.java.spi.messages.utils.ResponseItem;
import org.apache.plc4x.java.spi.values.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.File;
import java.util.Map;
import java.util.concurrent.TimeUnit;

public class HelloInflux {

    private static final Logger logger = LoggerFactory.getLogger(HelloInflux.class);

    private Configuration configuration;

    public HelloInflux(File configFile) {
        Configurations configs = new Configurations();
        try {
            configuration = configs.properties(configFile);
        } catch (ConfigurationException cex) {
            throw new RuntimeException("Error reading configuration");
        }
    }

    public void run() {
        InfluxDBClient dbConnection = connectToDb();
        WriteApi writeApi = dbConnection.getWriteApi();
        try {
            PlcConnection plcConnection = connectToPlc();

            final PlcSubscriptionRequest subscriptionRequest =
                plcConnection.subscriptionRequestBuilder().addChangeOfStateTagAddress("query",
                    configuration.getString("plc.query")).build();
            final PlcSubscriptionResponse subscriptionResponse =
                subscriptionRequest.execute().get(10, TimeUnit.SECONDS);
            subscriptionResponse.getSubscriptionHandle("query").register(plcSubscriptionEvent -> {
                DefaultPlcSubscriptionEvent internalEvent = (DefaultPlcSubscriptionEvent) plcSubscriptionEvent;
                final Point point = Point.measurement(configuration.getString("influx.measurement"))
                    .time(plcSubscriptionEvent.getTimestamp().toEpochMilli(), WritePrecision.MS);
                final Map<String, ResponseItem<PlcValue>> values = internalEvent.getValues();
                values.forEach((tagName, tagResponsePair) -> {
                    final PlcResponseCode responseCode = tagResponsePair.getCode();
                    final PlcValue plcValue = tagResponsePair.getValue();
                    if(responseCode == PlcResponseCode.OK) {
                        PlcStruct structValue = (PlcStruct) plcValue;
                        for (String key : structValue.getKeys()) {
                            PlcValue subValue = structValue.getValue(key);
                            registerTags(point, key, subValue);
                        }
                    }
                });
                writeApi.writePoint(
                    configuration.getString("influx.bucket"), configuration.getString("influx.org"), point);
            });
        } catch (PlcException e) {
            logger.error("PLC Error", e);
        } catch (Exception e) {
            logger.error("General Error", e);
        }
    }

    private void registerTags(Point point, String contextName, PlcValue plcValue) {
        if (contextName.equals("address")) {
            point.addTag(contextName, plcValue.getString());
        } else {
            if (plcValue instanceof PlcBOOL) {
                point.addField(contextName, plcValue.getBoolean());
            } else if (plcValue instanceof PlcSINT) {
                point.addField(contextName, plcValue.getByte());
            } else if (plcValue instanceof PlcUSINT) {
                point.addField(contextName, plcValue.getShort());
            } else if (plcValue instanceof PlcINT) {
                point.addField(contextName, plcValue.getShort());
            } else if (plcValue instanceof PlcUINT) {
                point.addField(contextName, plcValue.getInteger());
            } else if (plcValue instanceof PlcDINT) {
                point.addField(contextName, plcValue.getInteger());
            } else if (plcValue instanceof PlcUDINT) {
                point.addField(contextName, plcValue.getLong());
            } else if (plcValue instanceof PlcLINT) {
                point.addField(contextName, plcValue.getLong());
            } else if (plcValue instanceof PlcULINT) {
                point.addField(contextName, plcValue.getBigInteger());
            } else if (plcValue instanceof PlcREAL) {
                point.addField(contextName, plcValue.getFloat());
            } else if (plcValue instanceof PlcLREAL) {
                point.addField(contextName, plcValue.getDouble());
            } else if (plcValue instanceof PlcSTRING) {
                point.addField(contextName, plcValue.getString());
            } else if (plcValue instanceof PlcStruct) {
                PlcStruct structValue = (PlcStruct) plcValue;
                for (String key : structValue.getKeys()) {
                    PlcValue subValue = structValue.getValue(key);
                    registerTags(point, contextName + "-" + key, subValue);
                }
            }
        }
    }

    private InfluxDBClient connectToDb() {
        char[] token = configuration.getString("influx.accessToken").toCharArray();
        return InfluxDBClientFactory.create(configuration.getString("influx.connectionString"), token);
    }

    private PlcConnection connectToPlc() throws PlcException {
        final PlcConnection connection =
            new DefaultPlcDriverManager().getConnection(configuration.getString("plc.connectionString"));
        connection.connect();
        return connection;
    }

    public static void main(String[] args) {
        if(args.length != 1) {
            System.out.println("Usage: HelloInflux {path-to-config-file}");
        }
        final File configFile = new File(args[0]);
        if(!configFile.exists() || !configFile.isFile()) {
            throw new PlcRuntimeException("Could not read config file");
        }
        new HelloInflux(configFile).run();
    }

}
