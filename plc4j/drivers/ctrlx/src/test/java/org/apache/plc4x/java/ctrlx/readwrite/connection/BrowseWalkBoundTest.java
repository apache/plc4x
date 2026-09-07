/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.ctrlx.readwrite.connection;

import org.junit.jupiter.api.Test;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

/**
 * A browse asks the device for a node's children and then asks the same of each child, so how much
 * work it is belongs to the device. These cover the two ways that ends badly and the one answer
 * that used to end it with an exception.
 */
public class BrowseWalkBoundTest {

    @Test
    void aNodeNamingNoChildrenIsALeafRatherThanAFailure() {
        // The device can say "no children" by returning nothing at all, and reading that as a list
        // is what used to throw - inside a Runnable, where the exception went nowhere and the
        // browse future was simply never completed.
        List<String> leaves = new ArrayList<>();
        int read = CtrlXConnection.walkNodes(List.of("a"), node -> null, 0, 0, leaves::add);
        assertEquals(1, read);
        assertEquals(List.of("a"), leaves);
    }

    @Test
    void anEmptyChildListIsAlsoALeaf() {
        List<String> leaves = new ArrayList<>();
        CtrlXConnection.walkNodes(List.of("a"), node -> Collections.emptyList(), 0, 0, leaves::add);
        assertEquals(List.of("a"), leaves);
    }

    @Test
    void aTreeWithNoBottomIsStoppedByTheDepth() {
        // Every node names one child, forever. Paths never repeat, so remembering where the walk
        // has been would not notice; only the depth ends it.
        List<String> leaves = new ArrayList<>();
        int read = CtrlXConnection.walkNodes(List.of("a"), node -> List.of("down"), 0, 5, leaves::add);
        assertEquals(6, read, "the root plus five levels below it");
        assertTrue(leaves.isEmpty(), "nothing along the way was a leaf");
    }

    @Test
    void aTreeThatLeadsBackOnItselfIsAlsoStoppedByTheDepth() {
        // "a" names "b", "b" names "a" - as paths, a/b, a/b/a, a/b/a/b, each one new.
        List<String> leaves = new ArrayList<>();
        int read = CtrlXConnection.walkNodes(List.of("a"),
            node -> List.of(node.endsWith("b") ? "a" : "b"), 0, 4, leaves::add);
        assertEquals(5, read);
        assertTrue(leaves.isEmpty());
    }

    @Test
    void aTreeTooWideIsStoppedByTheTotal() {
        // Ten children each, so the queue outgrows any budget quickly.
        List<String> leaves = new ArrayList<>();
        int read = CtrlXConnection.walkNodes(List.of("a"),
            node -> List.of("0", "1", "2", "3", "4", "5", "6", "7", "8", "9"), 100, 0, leaves::add);
        assertEquals(100, read, "stops at the configured total rather than running on");
    }

    @Test
    void anOrdinaryTreeIsWalkedWhole() {
        List<String> leaves = new ArrayList<>();
        int read = CtrlXConnection.walkNodes(List.of("root"), node -> switch (node) {
            case "root" -> List.of("x", "y");
            case "root/x" -> List.of("leaf1", "leaf2");
            default -> Collections.emptyList();
        }, 0, 0, leaves::add);
        assertEquals(5, read);
        assertEquals(List.of("root/y", "root/x/leaf1", "root/x/leaf2"), leaves);
    }

    /**
     * The browse runs on an executor, and anything thrown out of a Runnable handed to an executor
     * goes to that thread's uncaught handler - not to the caller. The caller is waiting on the
     * future, so a failure that escapes leaves it waiting for good rather than telling it what
     * went wrong. Whatever the failure is, the future has to end up completed.
     */
    @Test
    void aFailureInsideTheBrowseCompletesTheFutureInsteadOfHanging() throws Exception {
        CtrlXConnection connection = new CtrlXConnection(
            "https://192.0.2.1", "user", "pass",
            new org.apache.plc4x.java.ctrlx.readwrite.configuration.CtrlXConfiguration());
        // Never connected, so the browse fails on the first thing it touches.
        java.util.concurrent.CompletableFuture<?> future =
            connection.browse(connection.browseRequestBuilder().build());
        java.util.concurrent.ExecutionException thrown =
            assertThrows(java.util.concurrent.ExecutionException.class,
                () -> future.get(10, java.util.concurrent.TimeUnit.SECONDS));
        assertNotNull(thrown.getCause());
    }

    @Test
    void aNodeThatCannotBeReadIsSkippedNotFatal() {
        List<String> leaves = new ArrayList<>();
        int read = CtrlXConnection.walkNodes(List.of("good", "bad"), node -> {
            if ("bad".equals(node)) {
                throw new org.apache.plc4x.java.ctrlx.readwrite.rest.datalayer.ApiException("nope");
            }
            return Collections.emptyList();
        }, 0, 0, leaves::add);
        assertEquals(2, read);
        assertEquals(List.of("good"), leaves, "the readable node is still reported");
    }
}
