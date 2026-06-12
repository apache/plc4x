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
package org.apache.plc4x.java.transport.rawsocket;

import org.pcap4j.core.PcapNetworkInterface;
import org.pcap4j.core.Pcaps;

import java.io.File;
import java.util.List;

import static org.junit.jupiter.api.Assumptions.assumeTrue;

/**
 * Helper for the pcap-based tests.
 * <p>
 * The pcap4j native library (libpcap on *nix/macOS, wpcap.dll on Windows) is not
 * available on all CI nodes - notably Windows runners without WinPcap/Npcap installed.
 * When the native library can't be loaded, pcap4j fails by throwing {@link Error}s
 * (e.g. {@link NoClassDefFoundError} / {@link UnsatisfiedLinkError} /
 * {@link ExceptionInInitializerError}) rather than checked exceptions, so a plain
 * {@code catch (Exception ...)} doesn't skip the affected tests and the build fails.
 * These helpers translate any {@link Throwable} into a JUnit assumption so the tests
 * are skipped instead.
 */
final class PcapTestSupport {

    private PcapTestSupport() {
    }

    /**
     * Points JNA at a libpcap installation in the well-known Homebrew locations on macOS.
     * On other platforms this is a no-op and the default library resolution is used.
     */
    static void configureNativeLibraryPath() {
        // For some reason it doesn't work if we pass this in from the outside.
        // On an Intel Mac the libs are in: "/usr/local/Cellar/libpcap/<version>/lib"
        // On an M1 Mac the libs are in: "/opt/homebrew/opt/libpcap/lib"
        if (new File("/usr/local/Cellar/libpcap/1.10.1/lib").exists()) {
            System.getProperties().setProperty("jna.library.path", "/usr/local/Cellar/libpcap/1.10.1/lib");
        } else if (new File("/usr/local/Cellar/libpcap/1.10.5/lib").exists()) {
            System.getProperties().setProperty("jna.library.path", "/usr/local/Cellar/libpcap/1.10.5/lib");
        } else if (new File("/opt/homebrew/opt/libpcap/lib").exists()) {
            System.getProperties().setProperty("jna.library.path", "/opt/homebrew/opt/libpcap/lib");
        }
    }

    /**
     * Returns the available pcap network interfaces, or aborts the current test (via a
     * JUnit assumption) if pcap is unavailable - either because the native library could
     * not be loaded or because no interfaces were found.
     *
     * @return the (non-empty) list of network interfaces
     */
    static List<PcapNetworkInterface> findAllDevsOrSkip() {
        configureNativeLibraryPath();
        List<PcapNetworkInterface> devs;
        try {
            devs = Pcaps.findAllDevs();
        } catch (Throwable t) {
            // Native library missing/unloadable (NoClassDefFoundError, UnsatisfiedLinkError, ...)
            // or pcap permission issues - skip rather than fail.
            assumeTrue(false, "pcap native library not available, skipping test: " + t);
            throw new AssertionError("unreachable");
        }
        assumeTrue(devs != null && !devs.isEmpty(), "No network interfaces found");
        return devs;
    }
}
