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

import org.junit.jupiter.api.Test;

import java.util.Arrays;

class OptimizerTest {

    @Test
    void playWithOptimizing() {
        ByteArrayFragment fragment1 = new ByteArrayFragment(new byte[]{1, 2, 3, 4, 5});

        Optimizer<ByteArrayFragment> optimizer = new Optimizer<>(new MaximalSizeCondition(2));

        // Optimize the ensemble with Conditions
        optimizer.optimize(Arrays.asList(fragment1));
    }

    private static class MaximalSizeCondition implements SingleFragmentCondition {

        private final int maximalSize;

        private MaximalSizeCondition(int maximalSize) {
            this.maximalSize = maximalSize;
        }

        @Override public boolean applies(Fragment fragment) {
            if (!(fragment instanceof ByteArrayFragment)) {
                return false;
            }
            return fragment.getBytes().length > maximalSize;
        }

        @Override public void apply(Fragment fragment, Ensemble ensemble) {
            // Split it up to two fragments with the appropriate rule
            byte[] bytes = fragment.getBytes();
            byte[] first = Arrays.copyOf(bytes, maximalSize);
            byte[] second = Arrays.copyOfRange(bytes, maximalSize, bytes.length);
            // Remove our Fragment from the Ensemble and add the two new ones
            ensemble.remove(fragment);
            ByteArrayFragment fragment1 = new ByteArrayFragment(first);
            ByteArrayFragment fragment2 = new ByteArrayFragment(second);
            ensemble.addFragment(fragment1);
            ensemble.addFragment(fragment2);
            // Create Inverse Rule and Register it
            ensemble.addOperation(new MergeByteArrayFragmentsOp(fragment1.getId(), fragment2.getId(), fragment.getId()));
        }

        @Override public String toString() {
            return "MaximumPaketSize(" + maximalSize + ")";
        }

        @Override public String getName() {
            return "MaximumPaketSize(" + maximalSize + ")";
        }

        private static class MergeByteArrayFragmentsOp extends MergeOperation {

            private final int f1;
            private final int f2;
            private final int targetId;

            public MergeByteArrayFragmentsOp(int f1, int f2, int targetId) {
                this.f1 = f1;
                this.f2 = f2;
                this.targetId = targetId;
            }

            @Override public String getName() {
                return String.format("Merge #%d + #%d -> #%d", f1, f2, targetId);
            }

            @Override public boolean canMerge(Fragment f1, Fragment f2) {
                return false;
            }

            @Override public Fragment merge(Fragment f1, Fragment f2) {
                return null;
            }

            @Override public String toString() {
                return getName();
            }
        }
    }
}