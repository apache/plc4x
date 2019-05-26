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
import org.apache.plc4x.java.ads.api.commands.AdsAddDeviceNotificationResponse;
import org.apache.plc4x.java.ads.api.commands.AdsDeleteDeviceNotificationResponse;
import org.apache.plc4x.java.ads.api.commands.types.NotificationHandle;
import org.apache.plc4x.java.ads.api.commands.types.Result;
import org.apache.plc4x.java.ads.api.generic.types.AmsNetId;
import org.apache.plc4x.java.ads.api.generic.types.AmsPort;
import org.apache.plc4x.java.ads.model.AdsDataType;
import org.apache.plc4x.java.ads.model.AdsSubscriptionHandle;
import org.apache.plc4x.java.ads.model.DirectAdsField;
import org.apache.plc4x.java.api.types.PlcSubscriptionType;
import org.apache.plc4x.java.base.messages.*;
import org.apache.plc4x.java.base.model.SubscriptionPlcField;
import org.assertj.core.api.WithAssertions;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Nested;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.Mockito;
import org.mockito.junit.jupiter.MockitoExtension;

import java.net.InetAddress;
import java.time.Duration;
import java.time.temporal.ChronoUnit;
import java.util.Collections;
import java.util.LinkedHashMap;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

import static org.junit.Assert.assertEquals;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class AdsTcpPlcConnectionTests implements WithAssertions {

    AdsTcpPlcConnection SUT;

    @Mock
    Channel channelMock;

    @Mock
    PlcSubscriber plcSubscriber;

    ExecutorService executorService;

    @BeforeEach
    void setUp() throws Exception {
        SUT = AdsTcpPlcConnection.of(InetAddress.getByName("localhost"), AmsNetId.of("0.0.0.0.0.0"), AmsPort.of(13));
        // TODO: Refactor this to use the TestChannelFactory instead.
        channelMock = Mockito.mock(Channel.class, RETURNS_DEEP_STUBS);
        FieldUtils.writeField(SUT, "channel", channelMock, true);
        executorService = Executors.newFixedThreadPool(10);
    }

    @AfterEach
    void tearDown() {
        executorService.shutdownNow();
        SUT = null;
    }

    @Nested
    class Lifecycle {
        @Test
        void initialState() {
            assertEquals(SUT.getTargetAmsNetId().toString(), "0.0.0.0.0.0");
            assertEquals(SUT.getTargetAmsPort().toString(), "13");
        }
    }

    @Nested
    class Subscription {
        @Test
        void subscribe() {
            when(channelMock.writeAndFlush(any(PlcRequestContainer.class))).then(invocation -> {
                PlcRequestContainer plcRequestContainer = invocation.getArgument(0);
                PlcProprietaryResponse plcProprietaryResponse = Mockito.mock(InternalPlcProprietaryResponse.class, RETURNS_DEEP_STUBS);
                AdsAddDeviceNotificationResponse adsAddDeviceNotificationResponse = Mockito.mock(AdsAddDeviceNotificationResponse.class, RETURNS_DEEP_STUBS);
                when(adsAddDeviceNotificationResponse.getResult()).thenReturn(Result.of(0));
                when(adsAddDeviceNotificationResponse.getNotificationHandle()).thenReturn(NotificationHandle.of(1));
                when(plcProprietaryResponse.getResponse()).thenReturn(adsAddDeviceNotificationResponse);
                plcRequestContainer.getResponseFuture().complete(plcProprietaryResponse);
                return mock(ChannelFuture.class);
            });

            SUT.subscribe(new DefaultPlcSubscriptionRequest(
                plcSubscriber,
                new LinkedHashMap<>(
                    Collections.singletonMap("field1",
                        new SubscriptionPlcField(PlcSubscriptionType.CYCLIC, DirectAdsField.of("0/0:BOOL"), Duration.of(1, ChronoUnit.SECONDS)))
                )
            ));
        }

        @Test
        void unsubscribe() {
            when(channelMock.writeAndFlush(any(PlcRequestContainer.class))).then(invocation -> {
                PlcRequestContainer plcRequestContainer = invocation.getArgument(0);
                PlcProprietaryResponse plcProprietaryResponse = Mockito.mock(InternalPlcProprietaryResponse.class, RETURNS_DEEP_STUBS);
                AdsDeleteDeviceNotificationResponse adsDeleteDeviceNotificationResponse = Mockito.mock(AdsDeleteDeviceNotificationResponse.class, RETURNS_DEEP_STUBS);
                when(adsDeleteDeviceNotificationResponse.getResult()).thenReturn(Result.of(0));
                when(plcProprietaryResponse.getResponse()).thenReturn(adsDeleteDeviceNotificationResponse);
                plcRequestContainer.getResponseFuture().complete(plcProprietaryResponse);
                return mock(ChannelFuture.class);
            });

            SUT.unsubscribe(new DefaultPlcUnsubscriptionRequest(plcSubscriber,
                Collections.singletonList(new AdsSubscriptionHandle(plcSubscriber, "hurz", AdsDataType.BYTE, NotificationHandle.of(1))))
            );
        }
    }

    // TODO: Commented out as it was causing problems with Java 11
    /*@Nested
    class Registration {
        @Captor
        ArgumentCaptor<Consumer<AdsDeviceNotificationRequest>> consumerArgumentCaptor;

        @Test
        @Disabled("This test seems to be causing problems in Java 11")
        void register() throws Exception {
            Plc4x2AdsProtocol plc4x2AdsProtocol = mock(Plc4x2AdsProtocol.class);
            when(channelMock.pipeline().get(Plc4x2AdsProtocol.class)).thenReturn(plc4x2AdsProtocol);

            AtomicReference<PlcSubscriptionEvent> plcSubscriptionEventAtomicReference = new AtomicReference<>();
            SUT.register(plcSubscriptionEventAtomicReference::set);
            verify(plc4x2AdsProtocol).addConsumer(consumerArgumentCaptor.capture());

            consumerArgumentCaptor.getValue().accept(AdsDeviceNotificationRequest.of(mock(AmsHeader.class), Length.of(1), Stamps.of(1), Collections.singletonList(AdsStampHeader.of(TimeStamp.of(1), Collections.singletonList(AdsNotificationSample.of(NotificationHandle.of(1), Data.of("Hello World!")))))));
            TimeUnit.MILLISECONDS.sleep(100);
            assertThat(plcSubscriptionEventAtomicReference).isNotNull();
        }

        @Test
        void unregister() {
            SUT.unregister(mock(InternalPlcConsumerRegistration.class));
        }
    }*/

    @Nested
    class Misc {
        @Test
        void remainingMethods() {
            assertThat(SUT.canSubscribe()).isTrue();
            assertThat(SUT.subscriptionRequestBuilder()).isNotNull();
            assertThat(SUT.unsubscriptionRequestBuilder()).isNotNull();
        }
    }

}