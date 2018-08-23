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

package org.apache.plc4x.java.ads.connection;

import io.netty.channel.Channel;
import io.netty.channel.ChannelFuture;
import org.apache.commons.lang3.reflect.FieldUtils;
import org.apache.plc4x.java.ads.api.commands.*;
import org.apache.plc4x.java.ads.api.commands.types.*;
import org.apache.plc4x.java.ads.api.generic.types.AmsNetId;
import org.apache.plc4x.java.ads.api.generic.types.AmsPort;
import org.apache.plc4x.java.ads.model.AdsAddress;
import org.apache.plc4x.java.ads.model.SymbolicAdsAddress;
import org.apache.plc4x.java.ads.protocol.Plc4x2AdsProtocol;
import org.apache.plc4x.java.api.exceptions.PlcInvalidAddressException;
import org.apache.plc4x.java.api.messages.PlcProprietaryRequest;
import org.apache.plc4x.java.api.messages.PlcProprietaryResponse;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.apache.plc4x.java.api.messages.items.SubscriptionEventItem;
import org.apache.plc4x.java.api.messages.items.SubscriptionRequestChangeOfStateItem;
import org.apache.plc4x.java.base.messages.PlcRequestContainer;
import org.hamcrest.Matchers;
import org.junit.After;
import org.junit.Before;
import org.junit.Test;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.InetAddress;
import java.util.Collections;
import java.util.Date;
import java.util.List;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.TimeUnit;
import java.util.function.Consumer;

import static org.hamcrest.core.IsEqual.equalTo;
import static org.hamcrest.core.IsNull.notNullValue;
import static org.junit.Assert.*;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.*;

public class AdsTcpPlcConnectionTests {

    private static final Logger LOGGER = LoggerFactory.getLogger(AdsTcpPlcConnectionTests.class);

    private AdsTcpPlcConnection SUT;

    private Channel channelMock;

    private ExecutorService executorService;

    @Before
    public void setUp() throws Exception {
        SUT = AdsTcpPlcConnection.of(InetAddress.getByName("localhost"), AmsNetId.of("0.0.0.0.0.0"), AmsPort.of(13));
        // TODO: Refactor this to use the TestChannelFactory instead.
        channelMock = mock(Channel.class, RETURNS_DEEP_STUBS);
        FieldUtils.writeField(SUT, "channel", channelMock, true);
        executorService = Executors.newFixedThreadPool(10);
    }

    @After
    public void tearDown() {
        executorService.shutdownNow();
        SUT = null;
    }

    @Test
    public void initialState() {
        assertEquals(SUT.getTargetAmsNetId().toString(), "0.0.0.0.0.0");
        assertEquals(SUT.getTargetAmsPort().toString(), "13");
    }

    @Test
    public void emptyParseAddress() {
        try {
            SUT.parseAddress("");
        } catch (PlcInvalidAddressException exception) {
            assertThat(exception.getMessage(), Matchers.startsWith(" invalid"));
        }
    }

    @Test
    public void parseAddress() throws Exception {
        try {
            AdsAddress address = (AdsAddress) SUT.parseAddress("1/1");
            assertEquals(address.getIndexGroup(), 1);
            assertEquals(address.getIndexOffset(), 1);
        } catch (IllegalArgumentException exception) {
            fail("valid data block address");
        }
    }

    @Test
    public void parseSymbolicAddress() throws Exception {
        try {
            SymbolicAdsAddress address = (SymbolicAdsAddress) SUT.parseAddress("Main.variable");
            assertEquals(address.getSymbolicAddress(), "Main.variable");
        } catch (IllegalArgumentException exception) {
            fail("valid data block address");
        }
    }

    @Test
    public void subscribe() throws Exception {
        // TODO: Does this really test the driver implementation?
        when(channelMock.writeAndFlush(any(PlcRequestContainer.class)))
            .then(invocationOnMock -> {
                PlcRequestContainer plcRequestContainer = invocationOnMock.getArgument(0);
                PlcProprietaryResponse plcProprietaryResponse = mock(PlcProprietaryResponse.class, RETURNS_DEEP_STUBS);

                PlcProprietaryRequest plcProprietaryRequest = (PlcProprietaryRequest) plcRequestContainer.getRequest();
                if (plcProprietaryRequest.getRequest() instanceof AdsAddDeviceNotificationRequest) {
                    AdsAddDeviceNotificationResponse adsAddDeviceNotificationResponse = mock(AdsAddDeviceNotificationResponse.class, RETURNS_DEEP_STUBS);
                    when(adsAddDeviceNotificationResponse.getResult().toAdsReturnCode()).thenReturn(AdsReturnCode.ADS_CODE_0);
                    when(adsAddDeviceNotificationResponse.getNotificationHandle()).thenReturn(NotificationHandle.of(0));
                    when(plcProprietaryResponse.getResponse()).thenReturn(adsAddDeviceNotificationResponse);
                } else if (plcProprietaryRequest.getRequest() instanceof AdsReadWriteRequest) {
                    AdsReadWriteResponse adsReadWriteResponse = mock(AdsReadWriteResponse.class, RETURNS_DEEP_STUBS);
                    when(adsReadWriteResponse.getData().getBytes()).thenReturn(new byte[]{0, 0, 0, 0});
                    when(adsReadWriteResponse.getResult().toAdsReturnCode()).thenReturn(AdsReturnCode.ADS_CODE_0);
                    when(plcProprietaryResponse.getResponse()).thenReturn(adsReadWriteResponse);
                }

                plcRequestContainer.getResponseFuture().complete(plcProprietaryResponse);
                return mock(ChannelFuture.class);
            });
        Plc4x2AdsProtocol plc4x2AdsProtocol = mock(Plc4x2AdsProtocol.class);
        when(plc4x2AdsProtocol.addConsumer(any())).then(invocation -> {
            Consumer<AdsDeviceNotificationRequest> consumer = invocation.getArgument(0);
            executorService.submit(() -> {
                while (!Thread.currentThread().isInterrupted()) {
                    AdsDeviceNotificationRequest mock = mock(AdsDeviceNotificationRequest.class);
                    AdsStampHeader adsStampHeader = mock(AdsStampHeader.class, RETURNS_DEEP_STUBS);
                    when(adsStampHeader.getTimeStamp()).thenReturn(TimeStamp.of(new Date()));
                    AdsNotificationSample adsNotificationSample = mock(AdsNotificationSample.class, RETURNS_DEEP_STUBS);
                    when(adsNotificationSample.getNotificationHandle()).thenReturn(NotificationHandle.of(0));
                    when(adsNotificationSample.getData()).thenReturn(Data.of("Hello " + consumer));
                    when(adsStampHeader.getAdsNotificationSamples()).thenReturn(Collections.singletonList(adsNotificationSample));
                    List<AdsStampHeader> adsStampHeaders = Collections.singletonList(adsStampHeader);
                    when(mock.getAdsStampHeaders()).thenReturn(adsStampHeaders);
                    consumer.accept(mock);
                }
            });
            return true;
        });
        when(channelMock.pipeline().get(Plc4x2AdsProtocol.class)).thenReturn(plc4x2AdsProtocol);

        CompletableFuture<?> notificationReceived = new CompletableFuture<>();
        Consumer<SubscriptionEventItem<String>> plcNotificationConsumer = plcNotification -> {
            LOGGER.info("Received {}", plcNotification);
            notificationReceived.complete(null);
        };
        PlcSubscriptionRequest subscriptionRequest = new PlcSubscriptionRequest();
        subscriptionRequest.addItem(new SubscriptionRequestChangeOfStateItem(
            String.class, SUT.parseAddress("0/0"), plcNotificationConsumer));
        /*subscriptionRequest.addItem(new SubscriptionRequestItem<>(
            String.class, SUT.parseAddress("Main.by[0]"), plcNotificationConsumer));*/
        CompletableFuture<? extends PlcSubscriptionResponse> subscriptionFuture = SUT.subscribe(subscriptionRequest);
        PlcSubscriptionResponse subscriptionResponse = subscriptionFuture.get(5, TimeUnit.SECONDS);
        //notificationReceived.get(3, TimeUnit.SECONDS);
        assertThat(subscriptionResponse, notNullValue());
        assertThat(subscriptionResponse.getNumberOfItems(), equalTo(1));

        // Now unsubscribe again ...

        // TODO: Setup the mock to actually perform the unsubscription.
        /*PlcUnsubscriptionRequest unsubscriptionRequest = new PlcUnsubscriptionRequest();
        for (SubscriptionResponseItem<?> subscriptionResponseItem : subscriptionResponse.getResponseItems()) {
            unsubscriptionRequest.addItem(subscriptionResponseItem.getSubscriptionHandle());
        }
        CompletableFuture<? extends PlcUnsubscriptionResponse> unsubscriptionFuture = SUT.unsubscribe(unsubscriptionRequest);
        PlcUnsubscriptionResponse plcUnsubscriptionResponse = unsubscriptionFuture.get(5, TimeUnit.SECONDS);
        assertThat(plcUnsubscriptionResponse, notNullValue());*/
    }
}