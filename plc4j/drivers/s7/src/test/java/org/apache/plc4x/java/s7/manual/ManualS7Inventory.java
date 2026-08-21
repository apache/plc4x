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
import org.apache.plc4x.java.api.messages.PlcDiscoveryItem;
import org.apache.plc4x.java.api.messages.PlcDiscoveryResponse;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.s7.S7CotpConnection;
import org.apache.plc4x.java.s7.userdata.S7SzlService.S7DeviceIdentification;

import java.io.IOException;
import java.io.PrintWriter;
import java.net.InetSocketAddress;
import java.net.Socket;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.TimeUnit;

/**
 * Manual test that builds an inventory of the Siemens S7 devices on the local network: what
 * they are, which firmware they run and whether they are protected.
 *
 * <p>Two steps. First a PROFINET DCP {@code IdentifyAll} broadcast finds the S7Comm-capable
 * devices — the same discovery {@link ManualS7Discovery} prints. DCP answers carry the station
 * name, MAC and IP but no firmware, so the second step connects to each device over ISO-on-TCP
 * and reads the identification system status lists (SZL {@code 0x0011}, {@code 0x001C} and
 * {@code 0x0232}).
 *
 * <p><b>This actively connects to every device it finds.</b> Unlike plain discovery, which only
 * listens for broadcast answers, this opens a TCP session to port 102 on each device. Run it
 * only on networks you are authorised to scan. It reads nothing but identification data and
 * changes no controller state, but it is unmistakably an active probe.
 *
 * <p>Requires libpcap on the host (and root/admin privileges on most platforms — pcap4j needs
 * raw socket access). On macOS, run with {@code sudo}; on Windows, run as Administrator with
 * Npcap installed.
 *
 * <p>The UNAUTHENTICATED column records whether the device handed its identification to a
 * caller that presented no credentials at all. That is the finding behind the CISA/NSA
 * guidance on securing S7 devices: anything reachable on port 102 that answers is talking to
 * whoever asks. Note the limit of the claim - it covers identification data only. This tool
 * never reads process data, so it cannot tell you whether memory or blocks are readable.
 *
 * <p>Usage: {@code ManualS7Inventory [output.csv]} — with an argument, the inventory is also
 * written to that file as CSV.
 */
public class ManualS7Inventory {

    /** Rack/slot combinations to try, in order: S7-1200/1500, then S7-300, then S7-400. */
    private static final int[] CANDIDATE_SLOTS = {0, 1, 2, 3};

    /** How long to wait for port 102 before deciding a device isn't speaking ISO-on-TCP. */
    private static final int PORT_PROBE_TIMEOUT_MS = 1_500;

    private static final int READ_TIMEOUT_MS = 3_000;

    public static void main(String[] args) throws Exception {
        System.out.println("Discovering S7 devices via PROFINET DCP ...");
        PlcDiscoveryResponse discovery = PlcDriverManager.getDefault()
            .getDriver("s7")
            .discoveryRequestBuilder()
            .build()
            .execute()
            .get();

        List<PlcDiscoveryItem> devices = discovery.getValues() == null
            ? List.of() : discovery.getValues();
        if (devices.isEmpty()) {
            System.out.println("No S7-capable devices discovered.");
            return;
        }
        System.out.printf("Discovered %d device(s); reading identification from each ...%n%n",
            devices.size());

        List<InventoryEntry> inventory = new ArrayList<>();
        for (PlcDiscoveryItem device : devices) {
            inventory.add(inventory(device));
        }

        printTable(inventory);
        if (args.length > 0) {
            writeCsv(inventory, Path.of(args[0]));
            System.out.printf("%nWrote %d row(s) to %s%n", inventory.size(), args[0]);
        }
    }

    /**
     * Connect to one discovered device and read its identification. Never throws: a device that
     * can't be reached still belongs in the inventory, with the reason in place of the data.
     */
    private static InventoryEntry inventory(PlcDiscoveryItem device) {
        String ipAddress = attribute(device, "ipAddress");
        String stationName = attribute(device, "deviceName");
        String deviceType = attribute(device, "deviceTypeName");

        if (ipAddress == null || ipAddress.isEmpty() || ipAddress.startsWith("0.0.0.0")) {
            return InventoryEntry.unreachable(ipAddress, stationName, deviceType,
                "no IP address assigned");
        }
        if (!isPortOpen(ipAddress)) {
            // Devices that answer DCP but keep port 102 closed: PROFINET IO devices that
            // aren't CPUs, or CPUs with S7 communication disabled.
            return InventoryEntry.unreachable(ipAddress, stationName, deviceType,
                "port 102 closed or filtered");
        }
        for (int slot : CANDIDATE_SLOTS) {
            // The rack/slot parameters belong to the COTP transport, hence the prefix.
            String url = String.format("s7://%s?cotp.remote-slot=%d&read-timeout=%d",
                ipAddress, slot, READ_TIMEOUT_MS);
            // getConnection() hands back an already-connected connection — connecting again
            // would re-run the S7 handshake on a live session and fail.
            try (PlcConnection connection = PlcDriverManager.getDefault()
                    .getConnectionFactory().getConnection(url)) {
                if (!(connection instanceof S7CotpConnection s7Connection)) {
                    return InventoryEntry.unreachable(ipAddress, stationName, deviceType,
                        "not an S7 COTP connection");
                }
                S7DeviceIdentification id = s7Connection.readDeviceIdentification()
                    .get(4 * READ_TIMEOUT_MS, TimeUnit.MILLISECONDS);
                // Connected, but every SZL was rejected — typical of LOGO and other low-end
                // devices that don't implement the system status lists at all.
                InventoryEntry entry = new InventoryEntry(
                    ipAddress, stationName, deviceType, slot, id, null);
                return entry.wasIdentified() ? entry : new InventoryEntry(
                    ipAddress, stationName, deviceType, slot, id,
                    "connected on slot " + slot + " but answered no identification SZL");
            } catch (Exception e) {
                // Wrong rack/slot is the common case — the CPU refuses the COTP connection
                // request. Try the next candidate before giving up on the device.
                if (slot == CANDIDATE_SLOTS[CANDIDATE_SLOTS.length - 1]) {
                    return InventoryEntry.unreachable(ipAddress, stationName, deviceType,
                        "no answer on rack 0 slots 0-3 (" + rootCause(e) + ")");
                }
            }
        }
        return InventoryEntry.unreachable(ipAddress, stationName, deviceType, "not reachable");
    }

    /** Cheap pre-check so unreachable hosts don't cost four connection attempts each. */
    private static boolean isPortOpen(String ipAddress) {
        try (Socket socket = new Socket()) {
            socket.connect(new InetSocketAddress(ipAddress, 102), PORT_PROBE_TIMEOUT_MS);
            return true;
        } catch (IOException e) {
            return false;
        }
    }

    private static String attribute(PlcDiscoveryItem device, String name) {
        PlcValue value = device.getAttributes().get(name);
        return value == null ? null : value.getString();
    }

    private static String rootCause(Throwable t) {
        Throwable cause = t;
        while (cause.getCause() != null && cause.getCause() != cause) {
            cause = cause.getCause();
        }
        String message = cause.getMessage();
        return message == null ? cause.getClass().getSimpleName() : message;
    }

    // ------------------------------------------------------------------------
    // Output
    // ------------------------------------------------------------------------

    private static void printTable(List<InventoryEntry> inventory) {
        String format = "%-15s  %-18s  %-22s  %-10s  %-16s  %-24s  %-22s%n";
        System.out.printf(format,
            "IP", "STATION", "ORDER CODE", "FIRMWARE", "SERIAL", "UNAUTHENTICATED", "PROTECTION");
        System.out.println("-".repeat(140));
        for (InventoryEntry entry : inventory) {
            System.out.printf(format,
                nvl(entry.ipAddress()),
                truncate(nvl(entry.stationName()), 18),
                truncate(entry.orderCode(), 22),
                entry.firmwareVersion(),
                truncate(entry.serialNumber(), 16),
                entry.unauthenticatedAccess(),
                entry.protection());
        }

        long identified = inventory.stream().filter(InventoryEntry::wasIdentified).count();
        System.out.printf("%n%d of %d device(s) disclosed identification to an unauthenticated "
            + "session.%n", identified, inventory.size());
        long refused = inventory.stream()
            .filter(e -> e.identification() != null)
            .filter(e -> e.identification().identificationFailureReason() != null)
            .filter(e -> e.identification().identificationFailureReason().startsWith("access denied"))
            .count();
        if (refused > 0) {
            System.out.printf("%d refused an unauthenticated caller.%n", refused);
        }
        long unprotected = inventory.stream()
            .filter(InventoryEntry::wasIdentified)
            .filter(e -> Boolean.FALSE.equals(e.identification().passwordProtected()))
            .count();
        if (unprotected > 0) {
            System.out.printf("%d of them report no configured protection level.%n", unprotected);
        }
        for (InventoryEntry entry : inventory) {
            if (entry.failure() != null) {
                System.out.printf("  %s: %s%n", nvl(entry.ipAddress()), entry.failure());
            }
        }
    }

    private static void writeCsv(List<InventoryEntry> inventory, Path path) throws IOException {
        try (PrintWriter writer = new PrintWriter(Files.newBufferedWriter(path,
                StandardCharsets.UTF_8))) {
            writer.println("ip,station,device_type,slot,order_code,module_type,firmware,"
                + "hardware,serial,plant,protection_level,password_protected,key_switch,"
                + "protection_unavailable_reason,unauthenticated_access,"
                + "identification_failure_reason,status");
            for (InventoryEntry entry : inventory) {
                S7DeviceIdentification id = entry.identification();
                writer.println(String.join(",",
                    csv(entry.ipAddress()),
                    csv(entry.stationName()),
                    csv(entry.deviceType()),
                    entry.slot() < 0 ? "" : Integer.toString(entry.slot()),
                    csv(id == null ? null : id.orderCode()),
                    csv(id == null ? null : id.moduleTypeName()),
                    csv(id == null ? null : id.firmwareVersion()),
                    csv(id == null ? null : id.hardwareVersion()),
                    csv(id == null ? null : id.serialNumber()),
                    csv(id == null ? null : id.plantDesignation()),
                    id == null || id.protectionLevel() == null
                        ? "" : Integer.toString(id.protectionLevel()),
                    id == null || id.passwordProtected() == null
                        ? "" : Boolean.toString(id.passwordProtected()),
                    csv(id == null ? null : id.keySwitchPosition()),
                    csv(id == null ? null : id.protectionUnavailableReason()),
                    csv(entry.unauthenticatedAccess()),
                    csv(id == null ? null : id.identificationFailureReason()),
                    csv(entry.failure() == null ? "ok" : entry.failure())));
            }
        }
    }

    private static String csv(String value) {
        if (value == null) {
            return "";
        }
        return "\"" + value.replace("\"", "\"\"") + "\"";
    }

    private static String nvl(String value) {
        return value == null || value.isEmpty() ? "-" : value;
    }

    private static String truncate(String value, int max) {
        String text = nvl(value);
        return text.length() <= max ? text : text.substring(0, max - 1) + "…";
    }

    /** One row of the inventory: what discovery said, plus what the device told us. */
    private record InventoryEntry(String ipAddress, String stationName, String deviceType,
                                  int slot, S7DeviceIdentification identification,
                                  String failure) {

        static InventoryEntry unreachable(String ipAddress, String stationName,
                                          String deviceType, String reason) {
            return new InventoryEntry(ipAddress, stationName, deviceType, -1, null, reason);
        }

        /**
         * A device that connected but answered no SZL at all hasn't been identified. Judged on
         * the data itself: an identification can carry nothing but the reason its protection
         * list was unavailable, and that alone isn't an answer.
         */
        boolean wasIdentified() {
            return identification != null
                && (identification.orderCode() != null
                    || identification.firmwareVersion() != null
                    || identification.moduleName() != null
                    || identification.serialNumber() != null
                    || identification.protectionLevel() != null);
        }

        String orderCode() {
            return identification == null ? "-" : nvl(identification.orderCode());
        }

        String firmwareVersion() {
            return identification == null ? "-" : nvl(identification.firmwareVersion());
        }

        String hardwareVersion() {
            return identification == null ? "-" : nvl(identification.hardwareVersion());
        }

        String serialNumber() {
            return identification == null ? "-" : nvl(identification.serialNumber());
        }

        /**
         * The security-relevant summary: the protection level in force and whether a password
         * backs it. Devices that don't implement the protection SZL report {@code "n/a"} —
         * absence of an answer is not evidence of an unprotected CPU.
         */
        /**
         * What an unauthenticated caller got out of this device. "disclosed" is the finding
         * the CISA/NSA advisories care about: order number, firmware and serial handed over
         * with no credentials, which means S7 communication is open to any host that can
         * reach port 102.
         *
         * <p>It is a statement about identification only. It does not say whether process
         * data could be read or written - this tool never attempts that.
         */
        String unauthenticatedAccess() {
            if (identification == null) {
                return "no session";
            }
            if (wasIdentified()) {
                return "disclosed";
            }
            String reason = identification.identificationFailureReason();
            if (reason != null && reason.startsWith("access denied")) {
                return "refused";
            }
            return reason == null ? "nothing disclosed" : reason;
        }

        String protection() {
            if (identification == null) {
                return "-";
            }
            if (identification.protectionLevel() == null) {
                // Say which it is. "not implemented" is the S7-1200/1500 answer and says
                // nothing about access rights; a refusal would show its own error code here.
                String reason = identification.protectionUnavailableReason();
                return reason == null ? "-" : reason;
            }
            return "level " + identification.protectionLevel()
                + (Boolean.TRUE.equals(identification.passwordProtected()) ? " (pwd)" : " (no pwd)");
        }
    }
}
