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
package org.apache.plc4x.java.ads.adslib;

import org.apache.commons.codec.binary.Hex;
import org.apache.commons.lang3.tuple.ImmutablePair;
import org.apache.plc4x.java.ads.api.commands.types.*;
import org.apache.plc4x.java.ads.api.generic.types.AmsNetId;
import org.apache.plc4x.java.ads.api.generic.types.AmsPort;
import org.apache.plc4x.java.ads.api.util.UnsignedIntLEByteValue;

import java.io.Console;
import java.io.PrintStream;
import java.util.concurrent.TimeUnit;

import static org.apache.plc4x.java.ads.adslib.AdsLib.*;

/**
 * Test ported from <a href="https://github.com/Beckhoff/ADS/blob/master/example/example.cpp">github</a>
 * <p>
 * On github there is a test project which can be used to test the protocol. The example there uses a c++ implementation.
 */
@SuppressWarnings("all")
public class ADSClientNotificationExample {

    private static void NotifyCallback(ImmutablePair<AmsNetId, AmsPort> pAddr, TimeStamp timeStamp, AdsNotificationSample notificationSample, int hUser) {
        Data notificationSampleData = notificationSample.getData();
        byte[] data = notificationSampleData.getBytes();
        System.out.print("NetId: " + pAddr.left +
            " hUser 0x" + Integer.toHexString(hUser) +
            " sample time: " + timeStamp.getAsDate() +
            " sample size: " + notificationSample.getSampleSize().getAsLong() +
            " value: ");
        assert data.length == notificationSample.getSampleSize().getAsLong();
        System.out.println(" 0x" + Hex.encodeHexString(data));
    }

    private static SymbolHandle getHandleByNameExample(PrintStream out, long port, ImmutablePair<AmsNetId, AmsPort> server, String handleName) {
        Output<byte[]> handle = new Output<>();
        Result handleStatus = AdsSyncReadWriteReqEx2(port,
            server,
            IndexGroup.ReservedGroups.ADSIGRP_SYM_HNDBYNAME,
            IndexOffset.NONE,
            ReadLength.of(4),
            handle,
            WriteLength.of(handleName.length()),
            Data.of(handleName),
            new Output<>());
        if (handleStatus.toAdsReturnCode() != AdsReturnCode.ADS_CODE_0) {
            out.println("Create handle for '" + handleName + "' failed with: " + handleStatus);
            return null;
        }
        return SymbolHandle.of(handle.value);
    }

    private static void releaseHandleExample(PrintStream out, long port, ImmutablePair<AmsNetId, AmsPort> server, SymbolHandle handle) {
        Result releaseHandle = AdsSyncWriteReqEx(port, server, IndexGroup.ReservedGroups.ADSIGRP_SYM_RELEASEHND, IndexOffset.NONE, WriteLength.of(4), Data.of(handle.getBytes()));
        if (releaseHandle.toAdsReturnCode() != AdsReturnCode.ADS_CODE_0) {
            out.println("Release handle " + handle + "' failed with: " + releaseHandle);
        }
    }


    private static void notificationExample(PrintStream out, long port, ImmutablePair<AmsNetId, AmsPort> server) {
        AdsNotificationAttribute adsNotificationAttribute = AdsNotificationAttribute.of(
            Length.of(1),
            AdsNotificationAttribute.ADSTRANS_SERVERCYCLE,
            MaxDelay.of(0),
            CycleTime.of(4000000)
        );

        Output<Long> notificationHandle = new Output<>();
        int hUser = 0;

        Result addStatus = AdsSyncAddDeviceNotificationReqEx(port,
            server,
            IndexGroup.of(0x4020),
            IndexOffset.of(4),
            adsNotificationAttribute,
            ADSClientNotificationExample::NotifyCallback,
            hUser,
            notificationHandle);

        if (addStatus.toAdsReturnCode() != AdsReturnCode.ADS_CODE_0) {
            out.println("Add device notification failed with: " + addStatus);
            return;
        }

        out.println("Hit ENTER to stop notifications");
        tryInteractiveWait();

        Result delStatus = AdsSyncDelDeviceNotificationReqEx(port, server, NotificationHandle.of(notificationHandle.value));
        if (delStatus.toAdsReturnCode() != AdsReturnCode.ADS_CODE_0) {
            out.println("Delete device notification failed with: " + delStatus);
            return;
        }
    }

    private static void notificationByNameExample(PrintStream out, long port, ImmutablePair<AmsNetId, AmsPort> server) {
        AdsNotificationAttribute adsNotificationAttribute = AdsNotificationAttribute.of(
            Length.of(1),
            AdsNotificationAttribute.ADSTRANS_SERVERCYCLE,
            MaxDelay.of(0),
            CycleTime.of(4000000)
        );
        Output<Long> notificationHandle = new Output<>();
        int hUser = 0;

        SymbolHandle handle;

        out.println("notificationByNameExample():");
        handle = getHandleByNameExample(out, port, server, "MAIN.byByte[4]");
        if (handle == null) {
            return;
        }

        Result addStatus = AdsSyncAddDeviceNotificationReqEx(port,
            server,
            IndexGroup.ReservedGroups.ADSIGRP_SYM_VALBYHND,
            IndexOffset.of(handle.getAsLong()),
            adsNotificationAttribute,
            ADSClientNotificationExample::NotifyCallback,
            hUser,
            notificationHandle);
        if (addStatus.toAdsReturnCode() != AdsReturnCode.ADS_CODE_0) {
            out.println("Add device notification failed with: " + addStatus);
            return;
        }

        out.println("Hit ENTER to stop by name notifications");
        tryInteractiveWait();

        Result delStatus = AdsSyncDelDeviceNotificationReqEx(port, server, NotificationHandle.of(notificationHandle.value));
        if (delStatus.toAdsReturnCode() != AdsReturnCode.ADS_CODE_0) {
            out.println("Delete device notification failed with: " + delStatus);
            return;
        }
        releaseHandleExample(out, port, server, handle);
    }

    private static void readExample(PrintStream out, long port, ImmutablePair<AmsNetId, AmsPort> server) {
        Output<Integer> bytesRead = new Output<>();
        Output<byte[]> buffer = new Output<>();

        out.println("readExample():");
        for (int i = 0; i < 8; ++i) {
            Result status = AdsSyncReadReqEx2(port, server, IndexGroup.of(0x4020), IndexOffset.NONE, Length.of(4), buffer, bytesRead);
            if (status.toAdsReturnCode() != AdsReturnCode.ADS_CODE_0) {
                out.println("ADS read failed with: " + status);
                return;
            }
            out.println("ADS read " + bytesRead.value + " bytes, value: 0x" + Hex.encodeHexString(buffer.value));
        }
    }

    private static void readByNameExample(PrintStream out, long port, ImmutablePair<AmsNetId, AmsPort> server) {
        Output<Integer> bytesRead = new Output<>();
        Output<byte[]> buffer = new Output<>();
        SymbolHandle handle;

        out.println("readByNameExample():");
        handle = getHandleByNameExample(out, port, server, "MAIN.byByte[4]");
        if (handle == null) {
            return;
        }

        for (int i = 0; i < 8; ++i) {
            Result status = AdsSyncReadReqEx2(port, server, IndexGroup.ReservedGroups.ADSIGRP_SYM_VALBYHND, IndexOffset.of(handle.getAsLong()), Length.of(4), buffer, bytesRead);
            if (status.toAdsReturnCode() != AdsReturnCode.ADS_CODE_0) {
                out.println("ADS read failed with: " + status);
                return;
            }
            out.println("ADS read " + bytesRead.value + " bytes, value: 0x" + Hex.encodeHexString(buffer.value));
        }
        releaseHandleExample(out, port, server, handle);
    }

    private static void readStateExample(PrintStream out, long port, ImmutablePair<AmsNetId, AmsPort> server) {
        Output<Integer> adsState = new Output<>();
        Output<Integer> devState = new Output<>();

        Result status = AdsSyncReadStateReqEx(port, server, adsState, devState);
        if (status.toAdsReturnCode() != AdsReturnCode.ADS_CODE_0) {
            out.println("ADS read failed with: " + status);
            return;
        }
        out.println("ADS state: " + adsState.value + " devState: " + devState.value);
    }

    private static void runExample(String remoteNetIdString, String remoteIpV4, PrintStream out) {
        AmsNetId remoteNetId = AmsNetId.of(remoteNetIdString);

        // uncomment and adjust if automatic AmsNetId deduction is not working as expected
        AdsSetLocalAddress(AmsNetId.of("10.10.56.23.1.1"));

        // add local route to your EtherCAT Master
        if (AdsAddRoute(remoteNetId, remoteIpV4).toAdsReturnCode() != AdsReturnCode.ADS_CODE_0) {
            out.println("Adding ADS route failed, did you specified valid addresses?");
            return;
        }

        // open a new ADS port
        long port = AdsPortOpenEx();
        if (port < 0) {
            out.println("Open ADS port failed");
            return;
        }


        ImmutablePair<AmsNetId, AmsPort> remote = ImmutablePair.of(remoteNetId, AmsPort.of(851));

        notificationExample(out, port, remote);
        notificationByNameExample(out, port, remote);
        readExample(out, port, remote);
        readByNameExample(out, port, remote);
        readStateExample(out, port, remote);

        Result closeStatus = AdsPortCloseEx(port);
        if (closeStatus.toAdsReturnCode() != AdsReturnCode.ADS_CODE_0) {
            out.println("Close ADS port failed with: " + closeStatus);
        }

        AdsDelRoute(remoteNetId);
    }

    public static void main(String... args) {
        String remoteNetIdString = "10.10.64.40.1.1";
        String remoteIpV4 = "10.10.64.40";
        if (args.length == 2) {
            remoteIpV4 = args[0];
            remoteNetIdString = args[1];
            System.out.println("Using supplied arguments " + remoteIpV4 + "/" + remoteNetIdString);
        }
        runExample(remoteNetIdString, remoteIpV4, System.out);
        System.exit(0);
    }

    /////
    // Utils
    private static void tryInteractiveWait() {
        Console console = System.console();
        if (console != null) {
            console.readLine();
        } else {
            try {
                int timeout = Integer.valueOf(System.getProperty("input.timeout", "3"));
                System.out.println("Using timeout of " + timeout + "Seconds as System.console() is not available. Override with -Dinput.timeout=3");
                TimeUnit.SECONDS.sleep(timeout);
                System.out.println("Timeout reached enter pressed");
            } catch (InterruptedException e) {
                throw new RuntimeException(e);
            }
        }
    }

    static class SymbolHandle extends UnsignedIntLEByteValue {

        public SymbolHandle(byte... bytes) {
            super(bytes);
        }

        public static SymbolHandle of(byte... bytes) {
            return new SymbolHandle(bytes);
        }
    }
}
