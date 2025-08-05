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
package org.apache.plc4x.test.driver.internal.validator;

import org.apache.plc4x.test.driver.exceptions.DriverTestsuiteException;
import org.apache.plc4x.test.driver.xmlunit.SkipAttributeFilter;
import org.apache.plc4x.test.driver.xmlunit.SkipDifferenceEvaluator;
import org.dom4j.Element;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.xmlunit.builder.DiffBuilder;
import org.xmlunit.diff.Diff;

public class ApiValidator {

    private static final Logger LOGGER = LoggerFactory.getLogger(ApiValidator.class);

    public static void validateApiMessage(Element referenceXml, String apiMessage) throws DriverTestsuiteException {
        final String referenceXmlString = referenceXml.asXML();
        final Diff diff = DiffBuilder.compare(referenceXmlString)
            .withTest(apiMessage).checkForSimilar().ignoreComments().ignoreWhitespace()
            .withDifferenceEvaluator(new SkipDifferenceEvaluator())
            .withAttributeFilter(new SkipAttributeFilter())
            .build();
        if (diff.hasDifferences()) {
            LOGGER.warn("got\n{}", apiMessage);
            LOGGER.warn("Expect\n{}", referenceXmlString);
            LOGGER.warn("diff\n{}", diff);
            throw new DriverTestsuiteException("Differences were found after parsing.\n" + diff);
        }
    }
}
