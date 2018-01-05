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
package org.apache.plc4x.java.applications.plclogger;

import java.util.Calendar;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.concurrent.atomic.AtomicLong;

import org.apache.edgent.function.Supplier;
import org.apache.edgent.providers.direct.DirectProvider;
import org.apache.edgent.topology.TStream;
import org.apache.edgent.topology.Topology;
import org.apache.plc4x.edgent.PlcConnectionAdapter;
import org.apache.plc4x.edgent.PlcFunctions;

public class PlcLogger {

    private final PlcConnectionAdapter plcAdapter;
    private final String addressStr;
    private final int interval;

    private PlcLogger(PlcConnectionAdapter plcAdapter, String addressString, int interval) {
        this.plcAdapter = plcAdapter;
        this.addressStr = addressString;
        this.interval = interval;
    }

    private void run() throws Exception {
        AtomicInteger counter = new AtomicInteger(0);
        AtomicLong totalTime = new AtomicLong(0);
        DirectProvider dp = new DirectProvider();
        Topology top = dp.newTopology();
        // normally would just do the following
        //   TStream<Byte> source = top.poll(
        //       PlcFunctions.byteSupplier(adapter, addressStr),
        //       interval, TimeUnit.MILLISECONDS);
        // but in this case we want to make some timing measurements
        Supplier<Byte> plcSupplier = PlcFunctions.byteSupplier(plcAdapter, addressStr);
        TStream<Byte> source = top.poll(() -> {
            long start = Calendar.getInstance().getTimeInMillis();
            Byte value = plcSupplier.get();
            long end = Calendar.getInstance().getTimeInMillis();
            long time = end - start;
            System.out.println("Time: " + time);
            int curCounter = counter.incrementAndGet();
            long curTotalTime = totalTime.addAndGet(time);
            System.out.println("Avg:  " + (curTotalTime / curCounter));
            return value;
          },
          interval, TimeUnit.MILLISECONDS);
        source.print();
        dp.submit(top);
    }

    public static void main(String[] args) {
        if(args.length != 3) {
            System.out.println("Usage: PlcLogger {connection-string} {resource-address-string} {interval-ms}");
            System.out.println("Example: PlcLogger s7://192.168.0.1/0/0 INPUTS/0 10");
        }

        String connectionString = args[0];
        String addressString = args[1];
        Integer interval = Integer.valueOf(args[2]);
        try (PlcConnectionAdapter plcAdapter = new PlcConnectionAdapter(connectionString)) {

            // Initialize the logger itself
            PlcLogger logger = new PlcLogger(plcAdapter, addressString, interval);

            // Start the logging ...
            logger.run();

            // Yeah ... well prevent the application from exiting ;-)
            while (true) {
                Thread.sleep(1000);
            }
        } catch (Exception e) {
            e.printStackTrace();
        }
    }


}
