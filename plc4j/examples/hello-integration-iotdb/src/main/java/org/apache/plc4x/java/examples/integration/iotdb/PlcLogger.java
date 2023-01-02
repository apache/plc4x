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
package org.apache.plc4x.java.examples.integration.iotdb;

import org.apache.edgent.function.Supplier;
import org.apache.edgent.providers.direct.DirectProvider;
import org.apache.edgent.topology.TStream;
import org.apache.edgent.topology.Topology;
import org.apache.plc4x.edgent.PlcConnectionAdapter;
import org.apache.plc4x.edgent.PlcFunctions;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.concurrent.TimeUnit;

/**
 * using this example, you can store data one by one into IoTDB.
 *
 * modified according to hello-integration-edgent
 *
 * arguments example:
 * --connection-string simulated://127.0.0.1 --tsg-address RANDOM/foo:Integer  --polling-interval 1000
 * --iotdb-address 127.0.0.1:6667 --iotdb-user-name root --iotdb-user-password root --iotdb-sg mi
 * --iotdb-device d1 --iotdb-datatype INT32
 */
public class PlcLogger {

    private static final Logger LOGGER = LoggerFactory.getLogger(PlcLogger.class);

    //Time series ID
    static String timeSeries;

    //device ID
    static String deviceId;

    //sensor ID
    static String sensor;

    static String dataType;

    static IIoTDBWriter ioTDBWriter = null;

    static boolean useJDBC;

    public static void main(String[] args) throws Exception {
        CliOptions options = CliOptions.fromArgs(args);
        if (options == null) {
            CliOptions.printHelp();
            // Could not parse.
            System.exit(1);
        }
        useJDBC = options.isUseJDBC();
        deviceId = String.format("root.%s.%s", options.getStorageGroup(), options.getDevice());
        sensor = options.getTagAddress().replace("/", "_").replace(":", "_");
        timeSeries = String.format("%s.%s", deviceId, sensor);
        dataType = options.getDevice();

        // Get IoTDB connection
        if (useJDBC) {
            ioTDBWriter = new IoTDBWriterWithJDBC(options.getIotdbIpPort(), options.getUser(), options.getPassword());
        } else {
            ioTDBWriter = new IoTDBWriterWithSession(options.getIotdbIpPort(), options.getUser(), options.getPassword());
        }

        //as we do not know whether the storage group is created or not, we can init the storage
        // group in force.
        //from v0.9.0 on, IoTDB can automatically create storage group if you enable the parameter
        //`enable_auto_create_schema` and set the suitable `default_storage_group_level`.
        // If so, this method can be ignored.
        ioTDBWriter.initStorageGroup("root." + options.getStorageGroup());

        //register the time series. from v0.9.0 on, IoTDB can automatically create timeseries
        //if you enable the parameter `enable_auto_create_schema`.
        //But if you want to explicitly define the series data type, you have to call the method.

        //ioTDBWriter.createTimeseries(timeSeries, dataType);

        // Get a plc connection.
        try (PlcConnectionAdapter plcAdapter = new PlcConnectionAdapter(options.getConnectionString())) {
            // Initialize the Edgent core.
            DirectProvider dp = new DirectProvider();
            Topology top = dp.newTopology();

            // Define the event stream.
            // 1) PLC4X source generating a stream of bytes.
            Supplier<Integer> plcSupplier = PlcFunctions.integerSupplier(plcAdapter,
                options.getTagAddress());
            // 2) Use polling to get an item from the byte-stream in regular intervals.
            TStream<Integer> source = top.poll(plcSupplier, options.getPollingInterval(),
                TimeUnit.MILLISECONDS);
            // 3) Output the events in the stream to IoTDB.
            source.peek(value -> ioTDBWriter.writeData(deviceId, sensor, System.currentTimeMillis(), value));
            // Submit the topology and hereby start the event streams.
            dp.submit(top);
        }
        //close IoTDB client.
        Runtime.getRuntime().addShutdownHook(new Thread(ioTDBWriter::close));
    }


}
