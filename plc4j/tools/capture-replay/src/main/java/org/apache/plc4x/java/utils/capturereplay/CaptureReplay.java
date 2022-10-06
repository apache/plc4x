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
package org.apache.plc4x.java.utils.capturereplay;

import org.pcap4j.core.*;
import org.pcap4j.packet.Packet;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.File;
import java.sql.Timestamp;
import java.util.concurrent.TimeUnit;

public class CaptureReplay {

    private static final Logger LOGGER = LoggerFactory.getLogger(CaptureReplay.class);

    private final File inputFile;
    private final String outputDevice;
    private final float replaySpeed;

    public CaptureReplay(CliOptions options) {
        inputFile = options.getInputFile();
        if (!(inputFile.exists() && inputFile.isFile())) {
            throw new IllegalArgumentException("Could not open file " + inputFile.getPath());
        }
        outputDevice = options.getOutputDevice();
        replaySpeed = options.getReplaySpeed();
    }

    public void run() throws PcapNativeException {
        try (PcapHandle readHandle = Pcaps.openOffline(inputFile.getAbsolutePath(), PcapHandle.TimestampPrecision.NANO)) {

            PcapNetworkInterface sendDevice = Pcaps.getDevByName(outputDevice);
            if (sendDevice == null) {
                throw new IllegalArgumentException("Could not open output device " + outputDevice);
            }
            try (PcapHandle sendHandle = sendDevice.openLive(65536, PcapNetworkInterface.PromiscuousMode.PROMISCUOUS, 100)) {
                // Start a thread that processes the callbacks from the raw socket and simply
                // forwards the bytes read to the buffer.
                readHandle.loop(-1, new PacketListener() {
                    private Timestamp lastPacketTime = null;

                    @Override
                    public void gotPacket(Packet packet) {
                        Timestamp curPacketTime = readHandle.getTimestamp();

                        // Only enable the throttling if it is not disabled.
                        // If last-time is not null, wait for the given number of nano-seconds.
                        if ((replaySpeed > 0) && (lastPacketTime != null)) {
                            int numMicrosecondsSleep = (int)
                                ((curPacketTime.getNanos() - lastPacketTime.getNanos()) * replaySpeed);
                            nanoSecondSleep(numMicrosecondsSleep);
                        }

                        // Send the packet to the output device ...
                        try {
                            sendHandle.sendPacket(packet);
                        } catch (PcapNativeException | NotOpenException e) {
                            LOGGER.error("Error sending packet", e);
                        }

                        // Remember the timestamp of the current packet.
                        lastPacketTime = curPacketTime;
                    }
                });
            } catch (PcapNativeException | NotOpenException e) {
                LOGGER.error("PCAP sending loop thread died!", e);
            } catch (InterruptedException e) {
                LOGGER.warn("PCAP sending loop thread was interrupted (hopefully intentionally)", e);
                Thread.currentThread().interrupt();
            }
        }
    }

    private void nanoSecondSleep(long numNanos) {
        try {
            TimeUnit.NANOSECONDS.sleep(numNanos);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }
    }

    public static void main(String[] args) throws Exception {
        CliOptions options = CliOptions.fromArgs(args);
        if (options == null) {
            CliOptions.printHelp();
            // Could not parse.
            System.exit(1);
        }

        CaptureReplay replay = new CaptureReplay(options);
        do {
            replay.run();
        } while (options.isLoop());
    }

}
