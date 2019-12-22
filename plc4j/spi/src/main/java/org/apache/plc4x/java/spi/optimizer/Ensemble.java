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

package org.apache.plc4x.java.spi.optimizer;

import java.util.ArrayDeque;
import java.util.Deque;
import java.util.HashSet;
import java.util.LinkedList;
import java.util.List;
import java.util.Queue;
import java.util.Set;
import java.util.concurrent.atomic.AtomicInteger;

/** Ensemble of fragments that share a stack **/
public class Ensemble {

    private final AtomicInteger fragmentCounter = new AtomicInteger(0);

    private final Set<Fragment> fragments = new HashSet<>();
    private final Deque<OptimizerOperation> operations = new LinkedList<>();

    public int getFragmendId() {
        return fragmentCounter.getAndIncrement();
    }

    public void addFragment(Fragment fragment) {
        if (fragment.getParent() != null) {
            throw new IllegalArgumentException("This Fragment alredy belongs to an ensemble!");
        }
        fragment.setId(getFragmendId());
        fragment.setParent(this);
        this.fragments.add(fragment);
    }

    public Set<Fragment> getFragments() {
        return this.fragments;
    }

    public void remove(Fragment fragment) {
        assert fragment.getParent() == this;
        // Remove by Id
        this.fragments.remove(fragment);
    }

    public void addOperation(OptimizerOperation operation) {
        this.operations.addFirst(operation);
    }

    @Override public String toString() {
        return "Ensemble{" +
            "fragments=" + fragments +
            ", operations=" + operations +
            '}';
    }

}
