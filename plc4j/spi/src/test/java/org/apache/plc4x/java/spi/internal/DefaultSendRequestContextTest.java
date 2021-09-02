/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.spi.internal;

import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.Plc4xNettyWrapper;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.time.Duration;
import java.util.function.Consumer;

import static org.hamcrest.Matchers.*;
import static org.hamcrest.MatcherAssert.assertThat;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.verify;

@SuppressWarnings({"rawtypes", "unchecked"})
@ExtendWith(MockitoExtension.class)
class DefaultSendRequestContextTest {

    @Mock
    Consumer<HandlerRegistration> finisher;

    @Mock
    Plc4xNettyWrapper.DefaultConversationContext context;

    DefaultSendRequestContext<Object> SUT;

    @BeforeEach
    void setUp() {
        SUT = new DefaultSendRequestContext(finisher, null, context);
    }

    @Test
    void expectResponse() {
        SUT.expectResponse(Object.class, Duration.ZERO);
        assertThat(SUT.expectClazz, is(Object.class));
        assertThat(SUT.commands, hasSize(1));
        assertThat(SUT.commands.getLast().get(), notNullValue());
        assertThrows(ConversationContext.PlcWiringException.class, () -> SUT.expectResponse(Object.class, Duration.ZERO));
    }

    @Test
    void check() {
        SUT.check(o -> true);
        assertThat(SUT.commands, hasSize(1));
        assertThat(SUT.commands.getLast().get(), notNullValue());
    }

    @Test
    void handle() {
        assertThat(SUT.packetConsumer, nullValue());
        SUT.handle(Object::notify);
        assertThat(SUT.commands, hasSize(0));
        assertThat(SUT.packetConsumer, notNullValue());
        assertThrows(ConversationContext.PlcWiringException.class, () -> SUT.handle(Object::notify));
    }

    @Test
    void onTimeout() {
        SUT.onTimeout(Throwable::printStackTrace);
        assertThat(SUT.commands, hasSize(0));
        assertThat(SUT.onTimeoutConsumer, notNullValue());
        assertThrows(ConversationContext.PlcWiringException.class, () -> SUT.onTimeout(Throwable::printStackTrace));
    }

    @Test
    void onError() {
        SUT.onError((p, e) -> e.printStackTrace());
        assertThat(SUT.commands, hasSize(0));
        assertThat(SUT.errorConsumer, notNullValue());
        assertThrows(ConversationContext.PlcWiringException.class, () -> SUT.onError((p, e) -> e.printStackTrace()));
    }

    @Test
    void unwrap() {
        assertThat(SUT.expectClazz, nullValue());
        assertThat(SUT.onTimeoutConsumer, nullValue());
        assertThrows(ConversationContext.PlcWiringException.class, () -> SUT.unwrap(Object::toString));
        SUT.expectResponse(Object.class, null);
        SUT.onTimeout(Throwable::printStackTrace);
        ConversationContext.SendRequestContext<String> unwrap = SUT.unwrap(Object::toString);
        assertThat(unwrap, notNullValue());
        assertThat(SUT.commands, hasSize(2));
        assertThat(SUT.expectClazz, notNullValue());
        assertThat(SUT.onTimeoutConsumer, notNullValue());
    }

    @Test
    void finish() {
        SUT.handle(o -> {
        });
        verify(finisher).accept(any());
    }
}