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
package org.apache.plc4x.java.s7.manual;

import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.PlcDriverManager;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;

import java.util.concurrent.TimeUnit;

/**
 * Manual smoke-test for the S7H dual-path failover using two independent S7-1500 CPUs as
 * stand-ins for an actual S7-400H redundant pair. The CPUs need not be the same model nor
 * have synchronised data — for failover state-machine validation we only require both to
 * respond to {@code %MW10:WORD} reads, which is universal on every S7-1500.
 *
 * <p>The reads will return <em>different values</em> from each CPU when failover happens
 * (each PLC has its own MW10), which is a useful visual signal that the swap actually
 * took place.</p>
 *
 * <p><b>How to use:</b></p>
 * <ol>
 *   <li>Run the program. It connects via the dual-path URL and starts reading
 *       {@code %MW10:WORD} every 2 seconds.</li>
 *   <li>While it's running, unplug the Ethernet cable from the <b>primary</b> S7-1500
 *       ({@code 192.168.24.66}). Within ~2 s of the next read attempt you should see
 *       {@code S7H operation failed on primary, retrying on secondary} in the log,
 *       followed by reads succeeding via the secondary at {@code .64} (with that PLC's
 *       MW10 value, possibly different from the primary's).</li>
 *   <li>Re-plug the primary cable. Within 4 s (next heartbeat tick), the wrapper pings
 *       the primary, sees it's healthy, and switches active back. You should see
 *       {@code primary recovered (ping succeeded)} → {@code switching active back to
 *       primary}, followed by reads at the primary's MW10 value again.</li>
 *   <li>Repeat with the secondary cable to verify the symmetric path.</li>
 * </ol>
 */
public class ManualWallS71500HFailover {

    /** Primary = S7-1511C-1 PN, Secondary = S7-1516-3 PN/DP. Different models, same family. */
    private static final String CONNECTION_URL =
        "s7://192.168.24.66/192.168.24.64?local-device-group=PG_OR_PC&remote-rack=0&remote-slot=1&ha-failover-timeout-ms=500&ha-heartbeat-interval-ms=1000";
    private static final long RUN_DURATION_MS = 120_000L;
    private static final long READ_INTERVAL_MS = 2_000L;

    public static void main(String[] args) throws Exception {
        try (PlcConnection connection = PlcDriverManager.getDefault()
                .getConnectionFactory()
                .getConnection(CONNECTION_URL)) {
            System.out.println("Connected to S7H dual-path (two independent S7-1500 CPUs). "
                + "metadata.readSupported=" + connection.getMetadata().isReadSupported());

            long deadline = System.currentTimeMillis() + RUN_DURATION_MS;
            int seq = 0;
            while (System.currentTimeMillis() < deadline) {
                seq++;
                long t0 = System.currentTimeMillis();
                try {
                    PlcReadResponse resp = connection.readRequestBuilder()
                        .addTagAddress("mw10", "%MW10:WORD")
                        .build()
                        .execute()
                        .get(5, TimeUnit.SECONDS);
                    long elapsed = System.currentTimeMillis() - t0;
                    PlcResponseCode code = resp.getResponseCode("mw10");
                    if (code == PlcResponseCode.OK) {
                        System.out.printf("[#%03d %4dms] mw10 = %s%n", seq, elapsed, resp.getPlcValue("mw10"));
                    } else {
                        System.out.printf("[#%03d %4dms] FAIL response code=%s%n", seq, elapsed, code);
                    }
                } catch (Exception e) {
                    long elapsed = System.currentTimeMillis() - t0;
                    System.out.printf("[#%03d %4dms] FAIL %s%n", seq, elapsed, e.getMessage());
                }
                Thread.sleep(READ_INTERVAL_MS);
            }
            System.out.println("Run complete.");
        }
    }
}
