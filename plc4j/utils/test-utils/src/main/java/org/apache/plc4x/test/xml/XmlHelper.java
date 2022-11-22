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
package org.apache.plc4x.test.xml;

import java.util.Collections;
import java.util.HashMap;
import java.util.Map;
import org.dom4j.Element;
import org.dom4j.QName;

/**
 * Common XML processing helpers used in driver and serializer test suites.
 */
public class XmlHelper {

    public static String extractText(Element element, String name) {
        Element child = element.element(new QName(name));
        if (child == null) {
            throw new RuntimeException("Required element " + name + " not present");
        }
        return child.getTextTrim();
    }

    // generic parser block for <parameter> tag with name/value tags
    public static Map<String, String> parseParameters(Element driverParametersElement) {
        if (driverParametersElement == null) {
            return Collections.emptyMap();
        }
        Map<String, String> driverParameters = new HashMap<>();
        for (Element parameter : driverParametersElement.elements(new QName("parameter"))) {
            String parameterName = extractText(parameter, "name");
            String parameterValue = extractText(parameter, "value");
            driverParameters.put(parameterName, parameterValue);
        }
        return driverParameters;
    }
}
