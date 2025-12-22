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

package org.apache.plc4x.java.spi.buffers.asciiboxbased.utils.lang3;

import java.util.Collections;
import java.util.Set;
import java.util.StringJoiner;
import java.util.function.BiConsumer;
import java.util.function.BinaryOperator;
import java.util.function.Function;
import java.util.function.Supplier;
import java.util.stream.Collector;

public class LangCollectors {

    /**
     * Simple implementation class for {@code Collector}.
     *
     * @param <T> the type of elements to be collected
     * @param <R> the type of the result
     */
    private static final class SimpleCollector<T, A, R> implements Collector<T, A, R> {

        private final BiConsumer<A, T> accumulator;
        private final Set<Characteristics> characteristics;
        private final BinaryOperator<A> combiner;
        private final Function<A, R> finisher;
        private final Supplier<A> supplier;

        private SimpleCollector(final Supplier<A> supplier, final BiConsumer<A, T> accumulator, final BinaryOperator<A> combiner, final Function<A, R> finisher,
                                final Set<Characteristics> characteristics) {
            this.supplier = supplier;
            this.accumulator = accumulator;
            this.combiner = combiner;
            this.finisher = finisher;
            this.characteristics = characteristics;
        }

        @Override
        public BiConsumer<A, T> accumulator() {
            return accumulator;
        }

        @Override
        public Set<Characteristics> characteristics() {
            return characteristics;
        }

        @Override
        public BinaryOperator<A> combiner() {
            return combiner;
        }

        @Override
        public Function<A, R> finisher() {
            return finisher;
        }

        @Override
        public Supplier<A> supplier() {
            return supplier;
        }
    }

    private static final Set<Collector.Characteristics> CH_NOID = Collections.emptySet();

    public static Collector<Object, ?, String> joining(final CharSequence delimiter, final CharSequence prefix, final CharSequence suffix,
                                                       final Function<Object, String> toString) {
        return new SimpleCollector<>(() -> new StringJoiner(delimiter, prefix, suffix), (a, t) -> a.add(toString.apply(t)), StringJoiner::merge,
            StringJoiner::toString, CH_NOID);
    }

}
