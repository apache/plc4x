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
package org.apache.plc4x.java.examples.integration.iotdb;

import java.sql.Connection;
import java.sql.DriverManager;
import java.sql.ResultSet;
import java.sql.SQLException;
import java.sql.Statement;
import org.apache.edgent.function.Supplier;
import org.apache.edgent.providers.direct.DirectProvider;
import org.apache.edgent.topology.TStream;
import org.apache.edgent.topology.Topology;
import org.apache.iotdb.jdbc.IoTDBSQLException;
import org.apache.plc4x.edgent.PlcConnectionAdapter;
import org.apache.plc4x.edgent.PlcFunctions;

import java.util.concurrent.TimeUnit;

/**
 * using this example, you can store data one by one into IoTDB.
 *
 * modified according to hello-integration-edgent
 *
 * arguments example:
 * --connection-string test://127.0.0.1 --field-address RANDOM/foo:BYTE  --polling-interval 1000
 * --iotdb-address 127.0.0.1:6667 --iotdb-user-name root --iotdb-user-password root --iotdb-sg mi
 * --iotdb-device d1 --iotdb-datatype INT32
 */
public class PlcLogger {

    //IoTDB JDBC connection
    static Connection connection;

    //IoTDB JDBC Statement
    static Statement statement;

    //Time series ID
    static String timeSeries;

    //device ID
    static String deviceId;

    //sensor ID
    static String sensor;

    public static void main(String[] args) throws Exception {
        CliOptions options = CliOptions.fromArgs(args);
        if (options == null) {
            CliOptions.printHelp();
            // Could not parse.
            System.exit(1);
        }

        deviceId = String.format("root.%s.%s", options.getStorageGroup(), options.getDevice());
        sensor = options.getFieldAddress().replace("/", "_").replace(":", "_");
        timeSeries = String.format("%s.%s", deviceId, sensor);

        // Get IoTDB connection
        Class.forName("org.apache.iotdb.jdbc.IoTDBDriver");
        connection = DriverManager.getConnection("jdbc:iotdb://" + options.getIotdbIpPort()+"/",
            options.getUser(), options.getPassword());
        statement = connection.createStatement();

        //as we do not know whether the storage group is created or not, we can init the storage
        // group in force.
        try {
            statement.execute("SET STORAGE GROUP TO root." + options.getStorageGroup());
        } catch (IoTDBSQLException e) {
            //from v0.9.0, you can use the error code to check whether the sg exists.
            System.err.println(e.getMessage());
        }

        //before IoTDB v0.9, we have to create timeseries manually
        try {
            statement.execute(String.format("CREATE TIMESERIES %s WITH DATATYPE=%s, ENCODING=RLE",
                timeSeries, options.getDatatype()));
        } catch (IoTDBSQLException e) {
            //from v0.9.0, you can use the error code to check whether the time series exists.
            System.err.println(e.getMessage());
        }

        // Get a plc connection.
        try (PlcConnectionAdapter plcAdapter = new PlcConnectionAdapter(options.getConnectionString())) {
            // Initialize the Edgent core.
            DirectProvider dp = new DirectProvider();
            Topology top = dp.newTopology();

            // Define the event stream.
            // 1) PLC4X source generating a stream of bytes.
            Supplier<Byte> plcSupplier = PlcFunctions.byteSupplier(plcAdapter,
                options.getFieldAddress());
            // 2) Use polling to get an item from the byte-stream in regular intervals.
            TStream<Byte> source = top.poll(plcSupplier, options.getPollingInterval(),
                TimeUnit.MILLISECONDS);
            // 3) Output the events in the stream to IoTDB.
            source.peek(x->storeData(x));

            // Submit the topology and hereby start the event streams.
            dp.submit(top);
        }
        //close IoTDB client.
        Runtime.getRuntime().addShutdownHook(new Thread() {
            public void run() {
                try {
                    connection.close();
                } catch (SQLException e) {
                    e.printStackTrace();
                }
            }
        });
    }

    private static void storeData(Byte x) {
        try {
            statement.execute(String.format("insert into %s  (timestamp, %s) values (%d, %s)",
                deviceId, sensor, System.currentTimeMillis(), x.byteValue()+""));
        } catch (SQLException e) {
            e.printStackTrace();
        }
    }

}
