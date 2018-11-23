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

import org.apache.commons.lang3.tuple.ImmutablePair;
import org.apache.commons.lang3.tuple.Triple;
import org.apache.plc4x.java.ads.api.commands.*;
import org.apache.plc4x.java.ads.api.commands.types.*;
import org.apache.plc4x.java.ads.api.generic.types.AmsError;
import org.apache.plc4x.java.ads.api.generic.types.AmsNetId;
import org.apache.plc4x.java.ads.api.generic.types.AmsPort;
import org.apache.plc4x.java.ads.api.generic.types.Invoke;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.Inet4Address;
import java.net.UnknownHostException;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;

/**
 * Ported from <a href="https://github.com/Beckhoff/ADS">github AdsLib</a>
 */
@SuppressWarnings("all")
public class AdsLib {

    private static final Logger LOGGER = LoggerFactory.getLogger(AdsLib.class);

    private static final AmsRouter amsRouter = new AmsRouter();

    private static AmsRouter GetRouter() {
        return amsRouter;
    }

    /**
     * Add new ams route to target system
     *
     * @param ams address of the target system
     * @param ip  address of the target system
     * @return [ADS Return Code](http://infosys.beckhoff.de/content/1033/tc3_adsdll2/html/ads_returncodes.htm?id=17663)
     */
    public static Result AdsAddRoute(AmsNetId ams, String ip) {
        try {
            return GetRouter().AddRoute(ams, (Inet4Address) Inet4Address.getByName(ip));
        } catch (UnknownHostException e) {
            LOGGER.error("Error adding route", e);
            return Result.of(AdsReturnCode.ADS_CODE_1);
        }
    }

    /**
     * Delete ams route that had previously been added with AdsAddRoute().
     *
     * @param ams address of the target system
     */
    public static void AdsDelRoute(AmsNetId ams) {
        GetRouter().DelRoute(ams);
    }

    /**
     * The connection (communication port) to the message router is
     * closed. The port to be closed must previously have been opened via
     * an AdsPortOpenEx() call.
     *
     * @param port port number of an Ads port that had previously been opened with AdsPortOpenEx().
     * @return [ADS Return Code](http://infosys.beckhoff.de/content/1033/tc3_adsdll2/html/ads_returncodes.htm?id=17663)
     */
    public static Result AdsPortCloseEx(long port) {
        return GetRouter().ClosePort((int) port);
    }

    /**
     * Establishes a connection (communication port) to the message
     * router. The port number returned by AdsPortOpenEx() is required as
     * parameter for further AdsLib function calls.
     *
     * @return port number of a new Ads port or 0 if no more ports available
     */
    public static long AdsPortOpenEx() {
        return GetRouter().OpenPort();
    }

    /**
     * Returns the local NetId and port number.
     *
     * @param port  port number of an Ads port that had previously been opened with AdsPortOpenEx().
     * @param pAddr Pointer to the structure of type {@code ImmutablePair<AmsNetId, AmsPort>}.
     * @return [ADS Return Code](http://infosys.beckhoff.de/content/1033/tc3_adsdll2/html/ads_returncodes.htm?id=17663)
     */
    public static Result AdsGetLocalAddressEx(long port, Output<ImmutablePair<AmsNetId, AmsPort>> pAddr) {
        return GetRouter().GetLocalAddress((int) port, pAddr);
    }

    /**
     * Change local NetId
     *
     * @param ams local AmsNetId
     */
    public static void AdsSetLocalAddress(AmsNetId ams) {
        GetRouter().SetLocalAddress(ams);
    }

    /**
     * Reads data synchronously from an ADS server.
     *
     * @param port         port number of an Ads port that had previously been opened with AdsPortOpenEx().
     * @param pAddr        Structure with NetId and port number of the ADS server.
     * @param indexGroup   Index Group.
     * @param indexOffset  Index Offset.
     * @param bufferLength Length of the data in bytes.
     * @param buffer       Pointer to a data buffer that will receive the data.
     * @param bytesRead    pointer to a variable. If successful, this variable will return the number of actually read data bytes.
     * @return [ADS Return Code](http://infosys.beckhoff.de/content/1033/tc3_adsdll2/html/ads_returncodes.htm?id=17663)
     */
    public static Result AdsSyncReadReqEx2(long port,
                                           ImmutablePair<AmsNetId, AmsPort> pAddr,
                                           IndexGroup indexGroup,
                                           IndexOffset indexOffset,
                                           Length bufferLength,
                                           Output<byte[]> buffer,
                                           Output<Integer> bytesRead) {
        Output<ImmutablePair<AmsNetId, AmsPort>> localAmsNet = new Output<>();
        Result getLocalAddressResult = GetRouter().GetLocalAddress((int) port, localAmsNet);
        if (getLocalAddressResult.toAdsReturnCode() != AdsReturnCode.ADS_CODE_0) {
            return getLocalAddressResult;
        }
        AdsReadRequest adsReadRequest = AdsReadRequest.of(
            pAddr.left,
            pAddr.right,
            localAmsNet.value.left,
            localAmsNet.value.right,
            Invoke.NONE,
            indexGroup,
            indexOffset,
            bufferLength
        );

        AmsRequest<AdsReadRequest, AdsReadResponse> request = AmsRequest.of(adsReadRequest);
        AmsError adsRequestResult = GetRouter().AdsRequest(request);
        if (adsRequestResult.toAdsReturnCode() != AdsReturnCode.ADS_CODE_0) {
            return Result.of(adsRequestResult.getAsLong());
        }
        try {
            AdsReadResponse response = request.getResponseFuture().get(3, TimeUnit.SECONDS);
            byte[] bytes = response.getData().getBytes();
            buffer.value = bytes;
            bytesRead.value = bytes.length;
            return response.getResult();
        } catch (ExecutionException | TimeoutException e) {
            e.printStackTrace();
            return Result.of(AdsReturnCode.ADS_CODE_1);
        } catch (InterruptedException e) {
            e.printStackTrace();
            Thread.currentThread().interrupt();
            return Result.of(AdsReturnCode.ADS_CODE_1);
        }
    }

    /**
     * Reads the identification and version number of an ADS server.
     *
     * @param port    port number of an Ads port that had previously been opened with AdsPortOpenEx().
     * @param pAddr   Structure with NetId and port number of the ADS server.
     * @param devName Pointer to a character string of at least 16 bytes, that will receive the name of the ADS device.
     * @param version Address of a variable of type AdsVersion, which will receive the version number, revision number and the build number.
     * @return [ADS Return Code](http://infosys.beckhoff.de/content/1033/tc3_adsdll2/html/ads_returncodes.htm?id=17663)
     */
    public static Result AdsSyncReadDeviceInfoReqEx(long port, ImmutablePair<AmsNetId, AmsPort> pAddr, Output<String> devName, Output<Triple<Byte, Byte, Integer>> version) {
        Output<ImmutablePair<AmsNetId, AmsPort>> localAmsNet = new Output<>();
        Result getLocalAddressResult = GetRouter().GetLocalAddress((int) port, localAmsNet);
        if (getLocalAddressResult.toAdsReturnCode() != AdsReturnCode.ADS_CODE_0) {
            return getLocalAddressResult;
        }
        AdsReadDeviceInfoRequest adsReadDeviceInfoRequest = AdsReadDeviceInfoRequest.of(
            pAddr.left,
            pAddr.right,
            localAmsNet.value.left,
            localAmsNet.value.right,
            Invoke.NONE
        );
        AmsRequest<AdsReadDeviceInfoRequest, AdsReadDeviceInfoResponse> request = AmsRequest.of(adsReadDeviceInfoRequest);
        AmsError adsRequestResult = GetRouter().AdsRequest(request);
        if (adsRequestResult.toAdsReturnCode() != AdsReturnCode.ADS_CODE_0) {
            return Result.of(adsRequestResult.getAsLong());
        }
        try {
            AdsReadDeviceInfoResponse response = request.getResponseFuture().get(3, TimeUnit.SECONDS);
            devName.value = new String(response.getDevice().getBytes());
            version.value = Triple.of(response.getMajorVersion().getAsByte(), response.getMinorVersion().getAsByte(), response.getVersion().getAsInt());
            return response.getResult();
        } catch (InterruptedException | ExecutionException | TimeoutException e) {
            e.printStackTrace();
            return Result.of(AdsReturnCode.ADS_CODE_1);
        }
    }

    /**
     * Reads the ADS status and the device status from an ADS server.
     *
     * @param port     port number of an Ads port that had previously been opened with AdsPortOpenEx().
     * @param pAddr    Structure with NetId and port number of the ADS server.
     * @param adsState Address of a variable that will receive the ADS status (see data type [ADSSTATE](http://infosys.beckhoff.de/content/1033/tc3_adsdll2/html/tcadsdll_enumadsstate.htm?id=17630)).
     * @param devState Address of a variable that will receive the device status.
     * @return [ADS Return Code](http://infosys.beckhoff.de/content/1033/tc3_adsdll2/html/ads_returncodes.htm?id=17663)
     */
    public static Result AdsSyncReadStateReqEx(long port, ImmutablePair<AmsNetId, AmsPort> pAddr, Output<Integer> adsState, Output<Integer> devState) {
        Output<ImmutablePair<AmsNetId, AmsPort>> localAmsNet = new Output<>();
        Result getLocalAddressResult = GetRouter().GetLocalAddress((int) port, localAmsNet);
        if (getLocalAddressResult.toAdsReturnCode() != AdsReturnCode.ADS_CODE_0) {
            return getLocalAddressResult;
        }
        AdsReadStateRequest adsReadStateRequest = AdsReadStateRequest.of(
            pAddr.left,
            pAddr.right,
            localAmsNet.value.left,
            localAmsNet.value.right,
            Invoke.NONE
        );
        AmsRequest<AdsReadStateRequest, AdsReadStateResponse> request = AmsRequest.of(adsReadStateRequest);
        AmsError adsRequestResult = GetRouter().AdsRequest(request);
        if (adsRequestResult.toAdsReturnCode() != AdsReturnCode.ADS_CODE_0) {
            return Result.of(adsRequestResult.getAsLong());
        }
        try {
            AdsReadStateResponse response = request.getResponseFuture().get(3, TimeUnit.SECONDS);
            adsState.value = response.getAdsState().getAsInt();
            devState.value = response.getDeviceState().getAsInt();
            return response.getResult();
        } catch (ExecutionException | TimeoutException e) {
            e.printStackTrace();
            return Result.of(AdsReturnCode.ADS_CODE_1);
        } catch (InterruptedException e) {
            e.printStackTrace();
            Thread.currentThread().interrupt();
            return Result.of(AdsReturnCode.ADS_CODE_1);
        }
    }

    /**
     * Writes data synchronously into an ADS server and receives data back from the ADS server.
     *
     * @param port        port number of an Ads port that had previously been opened with AdsPortOpenEx().
     * @param pAddr       Structure with NetId and port number of the ADS server.
     * @param indexGroup  Index Group.
     * @param indexOffset Index Offset.
     * @param readLength  Length, in bytes, of the read buffer readData.
     * @param readData    Buffer for data read from the ADS server.
     * @param writeLength Length of the data, in bytes, send to the ADS server.
     * @param writeData   Buffer with data send to the ADS server.
     * @param bytesRead   pointer to a variable. If successful, this variable will return the number of actually read data bytes.
     * @return [ADS Return Code](http://infosys.beckhoff.de/content/1033/tc3_adsdll2/html/ads_returncodes.htm?id=17663)
     */
    public static Result AdsSyncReadWriteReqEx2(long port,
                                                ImmutablePair<AmsNetId, AmsPort> pAddr,
                                                IndexGroup indexGroup,
                                                IndexOffset indexOffset,
                                                ReadLength readLength,
                                                Output<byte[]> readData,
                                                WriteLength writeLength,
                                                Data writeData,
                                                Output<Integer> bytesRead) {
        Output<ImmutablePair<AmsNetId, AmsPort>> localAmsNet = new Output<>();
        Result getLocalAddressResult = GetRouter().GetLocalAddress((int) port, localAmsNet);
        if (getLocalAddressResult.toAdsReturnCode() != AdsReturnCode.ADS_CODE_0) {
            return getLocalAddressResult;
        }
        AdsReadWriteRequest adsReadWriteRequest = AdsReadWriteRequest.of(
            pAddr.left,
            pAddr.right,
            localAmsNet.value.left,
            localAmsNet.value.right,
            Invoke.NONE,
            indexGroup,
            indexOffset,
            readLength,
            writeData
        );
        AmsRequest<AdsReadWriteRequest, AdsReadWriteResponse> request = AmsRequest.of(adsReadWriteRequest);
        AmsError adsRequestResult = GetRouter().AdsRequest(request);
        if (adsRequestResult.toAdsReturnCode() != AdsReturnCode.ADS_CODE_0) {
            return Result.of(adsRequestResult.getAsLong());
        }
        try {
            AdsReadWriteResponse response = request.getResponseFuture().get(3, TimeUnit.SECONDS);
            byte[] bytes = response.getData().getBytes();
            readData.value = bytes;
            bytesRead.value = bytes.length;
            return response.getResult();
        } catch (InterruptedException | ExecutionException | TimeoutException e) {
            e.printStackTrace();
            return Result.of(AdsReturnCode.ADS_CODE_1);
        }
    }

    /**
     * Writes data synchronously to an ADS server.
     *
     * @param port         port number of an Ads port that had previously been opened with AdsPortOpenEx().
     * @param pAddr        Structure with NetId and port number of the ADS server.
     * @param indexGroup   Index Group.
     * @param indexOffset  Index Offset.
     * @param bufferLength Length of the data, in bytes, send to the ADS server.
     * @param buffer       Buffer with data send to the ADS server.
     * @return [ADS Return Code](http://infosys.beckhoff.de/content/1033/tc3_adsdll2/html/ads_returncodes.htm?id=17663)
     */
    public static Result AdsSyncWriteReqEx(long port,
                                           ImmutablePair<AmsNetId, AmsPort> pAddr,
                                           IndexGroup indexGroup,
                                           IndexOffset indexOffset,
                                           WriteLength bufferLength,
                                           Data buffer) {
        Output<ImmutablePair<AmsNetId, AmsPort>> localAmsNet = new Output<>();
        Result getLocalAddressResult = GetRouter().GetLocalAddress((int) port, localAmsNet);
        if (getLocalAddressResult.toAdsReturnCode() != AdsReturnCode.ADS_CODE_0) {
            return getLocalAddressResult;
        }
        AdsWriteRequest adsWriteRequest = AdsWriteRequest.of(
            pAddr.left,
            pAddr.right,
            localAmsNet.value.left,
            localAmsNet.value.right,
            Invoke.NONE,
            indexGroup,
            indexOffset,
            buffer
        );
        AmsRequest<AdsWriteRequest, AdsWriteResponse> request = AmsRequest.of(adsWriteRequest);
        AmsError adsRequestResult = GetRouter().AdsRequest(request);
        if (adsRequestResult.toAdsReturnCode() != AdsReturnCode.ADS_CODE_0) {
            return Result.of(adsRequestResult.getAsLong());
        }
        try {
            AdsWriteResponse response = request.getResponseFuture().get(3, TimeUnit.SECONDS);
            return response.getResult();
        } catch (ExecutionException | TimeoutException e) {
            e.printStackTrace();
            return Result.of(AdsReturnCode.ADS_CODE_1);
        } catch (InterruptedException e) {
            e.printStackTrace();
            Thread.currentThread().interrupt();
            return Result.of(AdsReturnCode.ADS_CODE_1);
        }
    }

    /**
     * Changes the ADS status and the device status of an ADS server.
     *
     * @param port         port number of an Ads port that had previously been opened with AdsPortOpenEx().
     * @param pAddr        Structure with NetId and port number of the ADS server.
     * @param adsState     New ADS status.
     * @param devState     New device status.
     * @param bufferLength Length of the additional data, in bytes, send to the ADS server.
     * @param buffer       Buffer with additional data send to the ADS server.
     * @return [ADS Return Code](http://infosys.beckhoff.de/content/1033/tc3_adsdll2/html/ads_returncodes.htm?id=17663)
     */
    public static Result AdsSyncWriteControlReqEx(long port,
                                                  ImmutablePair<AmsNetId, AmsPort> pAddr,
                                                  short adsState,
                                                  short devState,
                                                  int bufferLength,
                                                  byte[] buffer) {
        Output<ImmutablePair<AmsNetId, AmsPort>> localAmsNet = new Output<>();
        Result getLocalAddressResult = GetRouter().GetLocalAddress((int) port, localAmsNet);
        if (getLocalAddressResult.toAdsReturnCode() != AdsReturnCode.ADS_CODE_0) {
            return getLocalAddressResult;
        }
        AdsWriteControlRequest adsWriteControlRequest = AdsWriteControlRequest.of(
            pAddr.left,
            pAddr.right,
            localAmsNet.value.left,
            localAmsNet.value.right,
            Invoke.NONE,
            AdsState.of(adsState),
            DeviceState.of(devState),
            Data.of(buffer)
        );
        AmsRequest<AdsWriteControlRequest, AdsWriteControlResponse> request = AmsRequest.of(adsWriteControlRequest);
        AmsError adsRequestResult = GetRouter().AdsRequest(request);
        if (adsRequestResult.toAdsReturnCode() != AdsReturnCode.ADS_CODE_0) {
            return Result.of(adsRequestResult.getAsLong());
        }
        try {
            AdsWriteControlResponse response = request.getResponseFuture().get(3, TimeUnit.SECONDS);
            return response.getResult();
        } catch (InterruptedException | ExecutionException | TimeoutException e) {
            e.printStackTrace();
            return Result.of(AdsReturnCode.ADS_CODE_1);
        }
    }

    /**
     * A notification is defined within an ADS server (e.g. PLC). When a
     * certain event occurs a function (the callback function) is invoked in
     * the ADS client (C program).
     *
     * @param port          port number of an Ads port that had previously been opened with AdsPortOpenEx().
     * @param pAddr         Structure with NetId and port number of the ADS server.
     * @param indexGroup    Index Group.
     * @param indexOffset   Index Offset.
     * @param attribute     Pointer to the structure that contains further information.
     * @param pFunc         Pointer to the structure describing the callback function.
     * @param hUser         32-bit value that is passed to the callback function.
     * @param pNotification Address of the variable that will receive the handle of the notification.
     * @return [ADS Return Code](http://infosys.beckhoff.de/content/1033/tc3_adsdll2/html/ads_returncodes.htm?id=17663)
     */
    public static Result AdsSyncAddDeviceNotificationReqEx(long port,
                                                           ImmutablePair<AmsNetId, AmsPort> pAddr,
                                                           IndexGroup indexGroup,
                                                           IndexOffset indexOffset,
                                                           AdsNotificationAttribute attribute,
                                                           PAdsNotificationFuncEx pFunc,
                                                           int hUser,
                                                           Output<Long> pNotification) {
        Output<ImmutablePair<AmsNetId, AmsPort>> localAmsNet = new Output<>();
        Result getLocalAddressResult = GetRouter().GetLocalAddress((int) port, localAmsNet);
        if (getLocalAddressResult.toAdsReturnCode() != AdsReturnCode.ADS_CODE_0) {
            return getLocalAddressResult;
        }
        AdsAddDeviceNotificationRequest adsAddDeviceNotificationRequest = AdsAddDeviceNotificationRequest.of(
            pAddr.left,
            pAddr.right,
            localAmsNet.value.left,
            localAmsNet.value.right,
            Invoke.NONE,
            indexGroup,
            indexOffset,
            attribute.length,
            attribute.transmissionMode,
            attribute.maxDelay,
            attribute.cycleTime
        );
        AmsRouter.Notification notify = AmsRouter.Notification.of(pFunc, hUser, attribute.length, pAddr, port);
        AmsRequest<AdsAddDeviceNotificationRequest, AdsAddDeviceNotificationResponse> request = AmsRequest.of(adsAddDeviceNotificationRequest);
        AmsError adsRequestResult = GetRouter().AddNotification(request, pNotification, notify);
        if (adsRequestResult.toAdsReturnCode() != AdsReturnCode.ADS_CODE_0) {
            return Result.of(adsRequestResult.getAsLong());
        }
        try {
            AdsAddDeviceNotificationResponse response = request.getResponseFuture().get(3, TimeUnit.SECONDS);
            return response.getResult();
        } catch (ExecutionException | TimeoutException e) {
            e.printStackTrace();
            return Result.of(AdsReturnCode.ADS_CODE_1);
        } catch (InterruptedException e) {
            e.printStackTrace();
            Thread.currentThread().interrupt();
            return Result.of(AdsReturnCode.ADS_CODE_1);
        }
    }

    /**
     * A notification defined previously is deleted from an ADS server.
     *
     * @param port               port number of an Ads port that had previously been opened with AdsPortOpenEx().
     * @param pAddr              Structure with NetId and port number of the ADS server.
     * @param notificationHandle Address of the variable that contains the handle of the notification.
     * @return [ADS Return Code](http://infosys.beckhoff.de/content/1033/tc3_adsdll2/html/ads_returncodes.htm?id=17663)
     */
    public static Result AdsSyncDelDeviceNotificationReqEx(long port, ImmutablePair<AmsNetId, AmsPort> pAddr, NotificationHandle notificationHandle) {
        Output<ImmutablePair<AmsNetId, AmsPort>> localAmsNet = new Output<>();
        Result getLocalAddressResult = GetRouter().GetLocalAddress((int) port, localAmsNet);
        if (getLocalAddressResult.toAdsReturnCode() != AdsReturnCode.ADS_CODE_0) {
            return getLocalAddressResult;
        }
        AdsDeleteDeviceNotificationRequest adsDeleteDeviceNotificationRequest = AdsDeleteDeviceNotificationRequest.of(
            pAddr.left,
            pAddr.right,
            localAmsNet.value.left,
            localAmsNet.value.right,
            Invoke.NONE,
            notificationHandle
        );
        AmsRequest<AdsDeleteDeviceNotificationRequest, AdsDeleteDeviceNotificationResponse> request = AmsRequest.of(adsDeleteDeviceNotificationRequest);
        AmsError adsRequestResult = GetRouter().DelNotification((int) port, pAddr, request);
        if (adsRequestResult.toAdsReturnCode() != AdsReturnCode.ADS_CODE_0) {
            return Result.of(adsRequestResult.getAsLong());
        }
        try {
            AdsDeleteDeviceNotificationResponse response = request.getResponseFuture().get(3, TimeUnit.SECONDS);
            return response.getResult();
        } catch (ExecutionException | TimeoutException e) {
            e.printStackTrace();
            return Result.of(AdsReturnCode.ADS_CODE_1);
        } catch (InterruptedException e) {
            e.printStackTrace();
            Thread.currentThread().interrupt();
            return Result.of(AdsReturnCode.ADS_CODE_1);
        }
    }

    /**
     * Read the configured timeout for the ADS functions. The standard value is 5000 ms.
     *
     * @param port    port number of an Ads port that had previously been opened with AdsPortOpenEx().
     * @param timeout Buffer to store timeout value in ms.
     * @return [ADS Return Code](http://infosys.beckhoff.de/content/1033/tc3_adsdll2/html/ads_returncodes.htm?id=17663)
     */
    public static Result AdsSyncGetTimeoutEx(long port, Output<Integer> timeout) {
        return GetRouter().GetTimeout((int) port, timeout);
    }

    /**
     * Alters the timeout for the ADS functions. The standard value is 5000 ms.
     *
     * @param port    port number of an Ads port that had previously been opened with AdsPortOpenEx().
     * @param timeout Timeout in ms.
     * @return [ADS Return Code](http://infosys.beckhoff.de/content/1033/tc3_adsdll2/html/ads_returncodes.htm?id=17663)
     */
    public static Result AdsSyncSetTimeoutEx(long port, int timeout) {
        return GetRouter().SetTimeout((int) port, timeout);
    }

    ////
    // Utils

    /**
     * Ported from <a href="https://github.com/Beckhoff/ADS">github AdsLib</a>
     */
    @FunctionalInterface
    public interface PAdsNotificationFuncEx {
        void notifyCallback(ImmutablePair<AmsNetId, AmsPort> pAddr, TimeStamp timeStamp, AdsNotificationSample notificationSample, int hUser);
    }

    /**
     * Ported from <a href="https://github.com/Beckhoff/ADS">github AdsLib</a>
     */
    public static class AdsNotificationAttribute {

        public static TransmissionMode ADSTRANS_SERVERCYCLE = TransmissionMode.of(3);

        public final Length length;
        public final TransmissionMode transmissionMode;
        public final MaxDelay maxDelay;
        public final CycleTime cycleTime;

        public AdsNotificationAttribute(Length length, TransmissionMode transmissionMode, MaxDelay maxDelay, CycleTime cycleTime) {
            this.length = length;
            this.transmissionMode = transmissionMode;
            this.maxDelay = maxDelay;
            this.cycleTime = cycleTime;
        }

        public static AdsNotificationAttribute of(Length length, TransmissionMode transmissionMode, MaxDelay maxDelay, CycleTime cycleTime) {
            return new AdsNotificationAttribute(length, transmissionMode, maxDelay, cycleTime);
        }
    }

}
