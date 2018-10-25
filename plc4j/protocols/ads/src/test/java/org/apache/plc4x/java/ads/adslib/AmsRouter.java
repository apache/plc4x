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

import org.apache.commons.lang3.mutable.MutableInt;
import org.apache.commons.lang3.tuple.ImmutablePair;
import org.apache.plc4x.java.ads.api.commands.*;
import org.apache.plc4x.java.ads.api.commands.types.*;
import org.apache.plc4x.java.ads.api.generic.AmsPacket;
import org.apache.plc4x.java.ads.api.generic.types.AmsError;
import org.apache.plc4x.java.ads.api.generic.types.AmsNetId;
import org.apache.plc4x.java.ads.api.generic.types.AmsPort;
import org.apache.plc4x.java.ads.connection.AdsTcpPlcConnection;
import org.apache.plc4x.java.ads.protocol.Plc4x2AdsProtocol;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.base.messages.PlcProprietaryRequest;
import org.apache.plc4x.java.base.messages.PlcProprietaryResponse;

import java.io.IOException;
import java.net.*;
import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.*;
import java.util.function.Consumer;

/**
 * Ported from <a href="https://github.com/Beckhoff/ADS">github AdsLib</a>
 */
@SuppressWarnings("all")
public class AmsRouter {

    private static final int NUM_PORTS_MAX = 128;
    private static final int PORT_BASE = 30000;

    private AmsNetId localAddr;

    private Map<InetAddress, AdsTcpPlcConnection> connections = new HashMap<>();
    private Map<AmsNetId, AdsTcpPlcConnection> mapping = new HashMap<>();
    private Map<AdsTcpPlcConnection, MutableInt> refCounts = new HashMap<>();

    private Map<Integer, AdsLibPort> ports = new HashMap<>();

    public AmsRouter() {
        try {
            // TODO: ensure we have ipv4
            this.localAddr = AmsNetId.of(Inet4Address.getLocalHost().getHostAddress() + ".1.1");
        } catch (UnknownHostException e) {
            throw new RuntimeException(e);
        }
    }

    public AmsRouter(AmsNetId localAddr) {
        this.localAddr = localAddr;
    }

    synchronized Result AddRoute(AmsNetId ams, Inet4Address ip) {
        AdsTcpPlcConnection oldConnection = GetConnection(ams);
        if (oldConnection != null && !(ip.equals(((InetSocketAddress) oldConnection.getChannel().remoteAddress()).getAddress()))) {
            /*
             There is already a route for this AmsNetId, but with
             a different IP. The old route has to be deleted, first!
             */
            return Result.of(AdsReturnCode.ADS_CODE_1286);
        }

        AdsTcpPlcConnection conn = connections.get(ip);
        if (conn == null) {
            AdsTcpPlcConnection newConnection = AdsTcpPlcConnection.of(ip, null, null);
            // TODO: add listener here:
            try {
                newConnection.connect();
            } catch (PlcConnectionException e) {
                throw new RuntimeException(e);
            }
            connections.put(ip, newConnection);
            conn = newConnection;

            /* in case no local AmsNetId was set previously, we derive one */
            if (localAddr == null) {
                localAddr = conn.getSourceAmsNetId();
            }
        }

        MutableInt refCounter = refCounts.getOrDefault(conn, new MutableInt());
        refCounter.increment();
        mapping.put(ams, conn);
        return localAddr == null ? Result.of(AdsReturnCode.ADS_CODE_1) : Result.of(AdsReturnCode.ADS_CODE_0);
    }

    synchronized void DelRoute(AmsNetId ams) {
        AdsTcpPlcConnection route = mapping.get(ams);
        if (route != null) {
            AdsTcpPlcConnection conn = route;
            try {
                conn.close();
            } catch (PlcConnectionException e) {
                throw new RuntimeException(e);
            }
            MutableInt refCounter = refCounts.getOrDefault(conn, new MutableInt());
            if (0 == refCounter.decrementAndGet()) {
                mapping.remove(ams);
                refCounts.remove(conn);
                DeleteIfLastConnection(conn);
            }
        }
    }

    void DeleteIfLastConnection(AdsTcpPlcConnection conn) {
        if (conn != null) {
            if (mapping.containsValue(conn)) {
                return;
            }
            connections.remove(conn.getRemoteAddress());
        }
    }

    synchronized int OpenPort() {
        try {
            AdsLibPort adsLibPort = new AdsLibPort(mapping);
            int localPort = adsLibPort.getLocalPort();
            ports.put(localPort, adsLibPort);
            return localPort;
        } catch (IOException e) {
            return 0;
        }
    }

    synchronized Result ClosePort(int port) {
        AdsLibPort serverSocket = ports.get(port);
        if (serverSocket == null || (serverSocket.getLocalPort() < PORT_BASE) || (serverSocket.getLocalPort() >= PORT_BASE + NUM_PORTS_MAX)) {
            return Result.of(AdsReturnCode.ADS_CODE_1864);
        }
        serverSocket.close();
        return Result.of(AdsReturnCode.ADS_CODE_0);
    }

    synchronized Result GetLocalAddress(int port, Output<ImmutablePair<AmsNetId, AmsPort>> pAddr) {
        if ((port < PORT_BASE) || (port >= PORT_BASE + NUM_PORTS_MAX)) {
            return Result.of(AdsReturnCode.ADS_CODE_1864);
        }

        AdsLibPort serverSocket = ports.get(port);
        if (serverSocket == null) {
            return Result.of(AdsReturnCode.ADS_CODE_1864);
        }

        AmsPort amsPort = AmsPort.of(serverSocket.getLocalPort());
        pAddr.value = ImmutablePair.of(localAddr, amsPort);
        return Result.of(AdsReturnCode.ADS_CODE_0);
    }

    synchronized void SetLocalAddress(AmsNetId netId) {
        localAddr = netId;
    }

    synchronized Result GetTimeout(int port, Output<Integer> timeout) {
        if ((port < PORT_BASE) || (port >= PORT_BASE + NUM_PORTS_MAX)) {
            return Result.of(AdsReturnCode.ADS_CODE_1864);
        }
        AdsLibPort serverSocket = ports.get(port);
        timeout.value = serverSocket.getSoTimeout();
        return Result.of(AdsReturnCode.ADS_CODE_0);
    }

    synchronized Result SetTimeout(int port, int timeout) {
        if ((port < PORT_BASE) || (port >= PORT_BASE + NUM_PORTS_MAX)) {
            return Result.of(AdsReturnCode.ADS_CODE_1864);
        }
        AdsLibPort serverSocket = ports.get(port);
        serverSocket.setSoTimeout(timeout);
        return Result.of(AdsReturnCode.ADS_CODE_0);
    }

    synchronized AdsTcpPlcConnection GetConnection(AmsNetId amsDest) {
        AdsTcpPlcConnection conn = mapping.get(amsDest);
        if (conn == null) {
            return null;
        }
        return connections.get(conn.getRemoteAddress());
    }

    <T extends AmsPacket, R extends AmsPacket> AmsError AdsRequest(AmsRequest<T, R> request) {
        PlcProprietaryRequest<T> plcProprietaryRequest = request.getRequest();

        AdsTcpPlcConnection ads = GetConnection(plcProprietaryRequest.getProprietaryRequest().getAmsHeader().getTargetAmsNetId());
        if (ads == null) {
            return AmsError.of(AdsReturnCode.ADS_CODE_7);
        }
        CompletableFuture<PlcProprietaryResponse<R>> completableFuture = ads.send(plcProprietaryRequest);
        try {
            PlcProprietaryResponse<R> response = completableFuture.get(3, TimeUnit.SECONDS);
            request.getResponseFuture().complete(response.getResponse());
            return response.getResponse().getAmsHeader().getCode();
        } catch (ExecutionException | TimeoutException e) {
            e.printStackTrace();
            return AmsError.of(AdsReturnCode.ADS_CODE_1864);
        } catch (InterruptedException e) {
            e.printStackTrace();
            Thread.currentThread().interrupt();
            return AmsError.of(AdsReturnCode.ADS_CODE_1864);
        }
    }

    AmsError AddNotification(AmsRequest<AdsAddDeviceNotificationRequest, AdsAddDeviceNotificationResponse> request, Output<Long> pNotification, Notification notify) {
        PlcProprietaryRequest<AdsAddDeviceNotificationRequest> plcProprietaryRequest = request.getRequest();
        //if (request.bytesRead) {
        //    request.bytesRead = 0;
        //}

        AdsTcpPlcConnection ads = GetConnection(plcProprietaryRequest.getProprietaryRequest().getAmsHeader().getTargetAmsNetId());
        if (ads == null) {
            return AmsError.of(AdsReturnCode.ADS_CODE_7);
        }

        AdsLibPort port = ports.get(plcProprietaryRequest.getProprietaryRequest().getAmsHeader().getSourceAmsPort().getAsInt());
        CompletableFuture<PlcProprietaryResponse<AdsAddDeviceNotificationResponse>> send = ads.send(plcProprietaryRequest);
        try {
            PlcProprietaryResponse<AdsAddDeviceNotificationResponse> response = send.get(3, TimeUnit.SECONDS);
            if (response.getResponse().getResult().toAdsReturnCode() != AdsReturnCode.ADS_CODE_0) {
                return AmsError.of(response.getResponse().getResult().getAsLong());
            }
            pNotification.value = response.getResponse().getNotificationHandle().getAsLong();
            port.AddNotification(pNotification.value, notify);
            request.getResponseFuture().complete(response.getResponse());
            return response.getResponse().getAmsHeader().getCode();
        } catch (ExecutionException | TimeoutException e) {
            e.printStackTrace();
            return AmsError.of(AdsReturnCode.ADS_CODE_1);
        } catch (InterruptedException e) {
            e.printStackTrace();
            Thread.currentThread().interrupt();
            return AmsError.of(AdsReturnCode.ADS_CODE_1);
        }
    }

    AmsError DelNotification(int port, ImmutablePair<AmsNetId, AmsPort> pAddr, AmsRequest<AdsDeleteDeviceNotificationRequest, AdsDeleteDeviceNotificationResponse> request) {
        PlcProprietaryRequest<AdsDeleteDeviceNotificationRequest> plcProprietaryRequest = request.getRequest();

        AdsTcpPlcConnection ads = GetConnection(plcProprietaryRequest.getProprietaryRequest().getAmsHeader().getTargetAmsNetId());
        if (ads == null) {
            return AmsError.of(AdsReturnCode.ADS_CODE_7);
        }

        AdsLibPort adsLibPort = ports.get(port);
        CompletableFuture<PlcProprietaryResponse<AdsDeleteDeviceNotificationResponse>> send = ads.send(plcProprietaryRequest);
        try {
            PlcProprietaryResponse<AdsDeleteDeviceNotificationResponse> response = send.get(3, TimeUnit.SECONDS);

            adsLibPort.DelNotification(pAddr, plcProprietaryRequest.getProprietaryRequest().getNotificationHandle());
            request.getResponseFuture().complete(response.getResponse());
            return response.getResponse().getAmsHeader().getCode();
        } catch (ExecutionException | TimeoutException e) {
            e.printStackTrace();
            return AmsError.of(AdsReturnCode.ADS_CODE_1);
        } catch (InterruptedException e) {
            e.printStackTrace();
            Thread.currentThread().interrupt();
            return AmsError.of(AdsReturnCode.ADS_CODE_1);
        }
    }

    static class Notification {

        AdsLib.PAdsNotificationFuncEx pFunc;
        int hUser;
        Length length;
        ImmutablePair<AmsNetId, AmsPort> pAddr;
        long port;

        Notification(AdsLib.PAdsNotificationFuncEx pFunc, int hUser, Length length, ImmutablePair<AmsNetId, AmsPort> pAddr, long port) {
            this.pFunc = pFunc;
            this.hUser = hUser;
            this.length = length;
            this.pAddr = pAddr;
            this.port = port;
        }

        static Notification of(AdsLib.PAdsNotificationFuncEx pFunc, int hUser, Length length, ImmutablePair<AmsNetId, AmsPort> pAddr, long port) {
            return new Notification(pFunc, hUser, length, pAddr, port);
        }

    }

    static class AdsLibPort implements Consumer<AdsDeviceNotificationRequest> {

        private ServerSocket serverSocket;
        private Map<Long, Notification> notificationMap;
        private Map<AmsNetId, AdsTcpPlcConnection> mapping;

        public AdsLibPort(Map<AmsNetId, AdsTcpPlcConnection> mapping) throws IOException {
            notificationMap = new ConcurrentHashMap<>();
            this.mapping = mapping;
            ServerSocket socket = null;
            for (int i = PORT_BASE; i < PORT_BASE + NUM_PORTS_MAX; i++) {
                try {
                    socket = new ServerSocket(i);
                } catch (IOException ignore) {
                }
            }
            if (socket == null) {
                throw new IOException("Unable to open server in port range(" + PORT_BASE + '-' + PORT_BASE + NUM_PORTS_MAX + ')');
            }
            socket.setReuseAddress(true);
            serverSocket = socket;
        }

        public int getLocalPort() {
            return serverSocket.getLocalPort();
        }

        public void close() {
            try {
                serverSocket.close();
            } catch (IOException e) {
                e.printStackTrace();
            }
            for (AdsTcpPlcConnection adsTcpPlcConnection : mapping.values()) {
                adsTcpPlcConnection.getChannel().pipeline().get(Plc4x2AdsProtocol.class).removeConsumer(this);
            }
        }

        public Integer getSoTimeout() {
            try {
                return serverSocket.getSoTimeout();
            } catch (IOException e) {
                e.printStackTrace();
                return -1;
            }
        }

        public void setSoTimeout(int soTimeout) {
            try {
                serverSocket.setSoTimeout(soTimeout);
            } catch (SocketException e) {
                e.printStackTrace();
            }
        }

        public void AddNotification(Long notificationHandle, Notification notify) {
            notificationMap.put(notificationHandle, notify);
            AdsTcpPlcConnection connection = mapping.get(notify.pAddr.left);
            // TODO: filter for addr
            connection.getChannel().pipeline().get(Plc4x2AdsProtocol.class).addConsumer(this);
        }

        public void DelNotification(ImmutablePair<AmsNetId, AmsPort> pAddr, NotificationHandle notificationHandle) {
            // Note. pAddr is not used for anything.
            notificationMap.remove(notificationHandle.getAsLong());
            AdsTcpPlcConnection connection = mapping.get(pAddr.left);
            // TODO: filter for addr
            connection.getChannel().pipeline().get(Plc4x2AdsProtocol.class).addConsumer(this);
        }

        @Override
        public void accept(AdsDeviceNotificationRequest adsDeviceNotificationRequest) {
            for (AdsStampHeader adsStampHeader : adsDeviceNotificationRequest.getAdsStampHeaders()) {
                TimeStamp timeStamp = adsStampHeader.getTimeStamp();
                for (AdsNotificationSample adsNotificationSample : adsStampHeader.getAdsNotificationSamples()) {
                    NotificationHandle notificationHandle = adsNotificationSample.getNotificationHandle();
                    Notification notification = notificationMap.get(notificationHandle.getAsLong());
                    notification.pFunc.notifyCallback(notification.pAddr, timeStamp, adsNotificationSample, notification.hUser);
                }
            }
        }
    }


}
