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

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Arrays;
import java.util.List;

/**
 * General Principle
 *
 * API
 * S_Req = Set{Req_1{I, ..., I}, ..., Req_n{I, ..., I}}
 * ------------------
 * Optimizer
 * - S_opt = Set{Req*_1{I, ..., I}, ..., Req*_m{I, ..., I}}
 * - [O_1, ..., O_k] with O_i : Set of Response -> Set of Response
 * - O_1 : Set{Res*_1, ..., Res*_m} with Res*_1 response to Req*_1, ...
 * - O_k : Set{Res_1, ..., Res_n} with Res_1 response to Req_1, ...
 * - Im(O_i) = Sup(O_i+1)
 * - S* = Set{Res*_1, ..., Res*_m} -- O_1 --> S*_1 ---> .... -- O_k --> S = Set{Res_1, ..., Res_n}
 * ------------------
 * Protocol Layer (S*)
 * ------------------
 * Deoptimizer Layer (O, S_Res*
 *
 *
 *
 *
 *
 *
 *
 *
 * <code>
 *     Max Packet Size 2
 *     Read: (0 = [1,2,3,4,5])
 *
 *     Operation 1 : Split
 *      split(1, 2, -> 0) + (1 = [1, 2]), (2 = [3, 4, 5])
 *        Operation 2: Split
 *          merge(3, 4, -> 2) + (1 =)[1, 2]), (3 = [3, 4]), (4 = [5])
 * </code>
 */
public class Optimizer<T extends Fragment> {

    private static final Logger logger = LoggerFactory.getLogger(Optimizer.class);

    private final List<OptimizerCondition> conditions;

    public Optimizer(OptimizerCondition... conditions) {
        this.conditions = Arrays.asList(conditions);
    }

    public Ensemble optimize(List<T> fragments) {
        Ensemble ensemble = new Ensemble();
        fragments.forEach(ensemble::addFragment);

        logger.info("Starting optimization with Ensembke {}", ensemble);
        // Check if all Conditions match (return directly?!)
        boolean violation;
        do {
            // TODO this should also skip
            violation = false;
            for (OptimizerCondition condition : conditions) {
                if (condition instanceof SingleFragmentCondition) {
                    for (Fragment fragment : ensemble.getFragments()) {
                        if (((SingleFragmentCondition) condition).applies(fragment)) {
                            logger.warn("The Fragment {} violated constraint {}", fragment, condition);
                            ((SingleFragmentCondition) condition).apply(fragment, ensemble);
                            logger.info("Ensemble is now {}", ensemble);
                            violation = true;
                        }
                    }
                }
            }
        } while (violation);
        return null;
    }

    /**
     *  Important, the Fragments have to have identical Fragment ID!
     */
    public List<Fragment> deoptimize(Ensemble ensemble, List<ResponseFragment> fragments) {
        // Check if they match with their id's
        // Now we create a map id -> fragments to merge them in the appropriate order
        // Operation (Set<F> -> Set<F>)
    }

}
