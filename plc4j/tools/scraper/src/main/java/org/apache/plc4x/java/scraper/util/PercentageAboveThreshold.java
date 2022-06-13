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
package org.apache.plc4x.java.scraper.util;

import org.apache.commons.math3.exception.MathIllegalArgumentException;
import org.apache.commons.math3.stat.descriptive.UnivariateStatistic;
import org.apache.plc4x.java.scraper.Scraper;

import java.util.Arrays;
import java.util.stream.IntStream;

public class PercentageAboveThreshold implements UnivariateStatistic {

    private final double threshold;

    public PercentageAboveThreshold(double threshold) {
        this.threshold = threshold;
    }

    @Override
    public double evaluate(double[] values) throws MathIllegalArgumentException {
        long below = Arrays.stream(values)
            .filter(val -> val <= threshold)
            .count();
        return (double) below / values.length;
    }

    @Override
    public double evaluate(double[] values, int begin, int length) throws MathIllegalArgumentException {
        long below = IntStream.range(begin, length)
            .mapToDouble(i -> values[i])
            .filter(val -> val > threshold)
            .count();
        return 100.0 * below / length;
    }

    @Override
    public UnivariateStatistic copy() {
        return new PercentageAboveThreshold(threshold);
    }
}
