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
package org.apache.plc4x.java.eip.base;

import org.apache.plc4x.java.eip.base.EipTcpConnection.Correlator;
import org.apache.plc4x.java.eip.readwrite.EipPacket;
import org.junit.jupiter.api.Test;
import org.mockito.Mockito;

import java.util.concurrent.CompletableFuture;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertSame;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * Responses are matched to requests by send order, because EIP carries no transaction id. That is
 * only sound while every request gets answered - so the interesting cases are the ones where a
 * request does not.
 */
class EipCorrelationTest {

    private static EipPacket packet() {
        return Mockito.mock(EipPacket.class);
    }

    @Test
    void aResponseAnswersTheRequestInFrontOfIt() {
        Correlator correlator = new Correlator();
        CompletableFuture<EipPacket> first = new CompletableFuture<>();
        CompletableFuture<EipPacket> second = new CompletableFuture<>();
        correlator.register(first);
        correlator.register(second);

        EipPacket firstResponse = packet();
        EipPacket secondResponse = packet();
        correlator.deliver(firstResponse);
        correlator.deliver(secondResponse);

        assertSame(firstResponse, first.getNow(null));
        assertSame(secondResponse, second.getNow(null));
    }

    /**
     * The case that used to hand one caller another caller's data: a request gives up, its answer
     * arrives late, and by then somebody else is at the head of the queue.
     */
    @Test
    void aLateResponseIsNotGivenToTheNextRequest() {
        Correlator correlator = new Correlator();
        CompletableFuture<EipPacket> abandoned = new CompletableFuture<>();
        correlator.register(abandoned);
        correlator.timedOut(abandoned);

        CompletableFuture<EipPacket> next = new CompletableFuture<>();
        correlator.register(next);

        // The abandoned request's answer turns up now.
        correlator.deliver(packet());

        assertFalse(next.isDone(),
            "a response that may answer an abandoned request must not complete a later one");
    }

    @Test
    void afterTheLateResponseIsDiscardedCorrelationLinesUpAgain() {
        Correlator correlator = new Correlator();
        CompletableFuture<EipPacket> abandoned = new CompletableFuture<>();
        correlator.register(abandoned);
        correlator.timedOut(abandoned);

        CompletableFuture<EipPacket> next = new CompletableFuture<>();
        correlator.register(next);
        correlator.deliver(packet());          // discarded as possibly stale

        EipPacket own = packet();
        correlator.deliver(own);               // this one is genuinely next's
        assertSame(own, next.getNow(null),
            "once the stale response is out of the way, ordering holds again");
    }

    @Test
    void onlyOneResponseIsDiscardedPerRequestThatGaveUp() {
        Correlator correlator = new Correlator();
        CompletableFuture<EipPacket> abandoned = new CompletableFuture<>();
        correlator.register(abandoned);
        correlator.timedOut(abandoned);

        CompletableFuture<EipPacket> a = new CompletableFuture<>();
        CompletableFuture<EipPacket> b = new CompletableFuture<>();
        correlator.register(a);
        correlator.register(b);

        correlator.deliver(packet());          // discarded
        EipPacket forA = packet();
        EipPacket forB = packet();
        correlator.deliver(forA);
        correlator.deliver(forB);

        assertSame(forA, a.getNow(null));
        assertSame(forB, b.getNow(null));
    }

    @Test
    void aRequestThatWasNeverSentIsJustForgotten() {
        Correlator correlator = new Correlator();
        CompletableFuture<EipPacket> unsent = new CompletableFuture<>();
        correlator.register(unsent);
        correlator.forget(unsent);
        assertEquals(0, correlator.pendingCount());

        // Nothing is owed, so nothing is discarded: the next response belongs to the next request.
        CompletableFuture<EipPacket> next = new CompletableFuture<>();
        correlator.register(next);
        EipPacket response = packet();
        correlator.deliver(response);
        assertSame(response, next.getNow(null));
    }

    @Test
    void closingFailsEveryoneStillWaiting() {
        Correlator correlator = new Correlator();
        CompletableFuture<EipPacket> first = new CompletableFuture<>();
        CompletableFuture<EipPacket> second = new CompletableFuture<>();
        correlator.register(first);
        correlator.register(second);

        correlator.failAll(new IllegalStateException("closed"));

        assertTrue(first.isCompletedExceptionally());
        assertTrue(second.isCompletedExceptionally());
        assertEquals(0, correlator.pendingCount());
    }
}
