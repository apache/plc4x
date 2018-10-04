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
package org.apache.plc4x.java.examples.plclogger;

import org.apache.edgent.function.Supplier;
import org.apache.edgent.providers.direct.DirectProvider;
import org.apache.edgent.topology.TStream;
import org.apache.edgent.topology.Topology;
import org.apache.plc4x.edgent.PlcConnectionAdapter;
import org.apache.plc4x.edgent.PlcFunctions;

import java.util.concurrent.TimeUnit;

public class PlcLogger {

    public static void main(String[] args) throws Exception {
        if(args.length != 3) {
            System.out.println("Usage: PlcLogger {connection-string} {resource-address-string} {interval-ms}");
            System.out.println("Example: PlcLogger s7://10.10.64.20/0/0 INPUTS/0 10");
        }

        String connectionString = args[0];
        String addressString = args[1];
        Integer interval = Integer.valueOf(args[2]);

        // Get a plc connection.
        try (PlcConnectionAdapter plcAdapter = new PlcConnectionAdapter(connectionString)) {
            // Initialize the Edgent core.
            DirectProvider dp = new DirectProvider();
            Topology top = dp.newTopology();

            // Define the event stream.
            // 1) PLC4X source generating a stream of bytes.
            Supplier<Byte> plcSupplier = PlcFunctions.byteSupplier(plcAdapter, addressString);
            // 2) Use polling to get an item from the byte-stream in regular intervals.
            TStream<Byte> source = top.poll(plcSupplier, interval, TimeUnit.MILLISECONDS);
            // 3) Output the events in the stream on the console.
            source.print();

            // Submit the topology and hereby start the event streams.
            dp.submit(top);
        }
    }

}
