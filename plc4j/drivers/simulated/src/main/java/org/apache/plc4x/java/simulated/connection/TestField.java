/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing,
 software distributed under the License is distributed on an
 "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 KIND, either express or implied.  See the License for the
 specific language governing permissions and limitations
 under the License.
 */
package org.apache.plc4x.java.simulated.connection;

import org.apache.commons.text.WordUtils;
import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.model.PlcField;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Objects;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * Test address for accessing values in virtual devices.
 */
public class TestField implements PlcField {

    private static final Logger logger = LoggerFactory.getLogger(TestField.class);

    /**
     * Examples:
     * - {@code RANDOM/foo:INTEGER}
     * - {@code STDOUT/foo:STRING}
     */
    private static final Pattern ADDRESS_PATTERN = Pattern.compile("^(?<type>\\w+)/(?<name>\\w+):(?<dataType>[a-zA-Z]++)(\\[(?<numElements>\\d+)])?$");

    private final TestType type;
    private final String name;
    private final Class<?> dataType;
    private final int numElements;

    private TestField(TestType type, String name, Class<?> dataType, int numElements) {
        this.type = type;
        this.name = name;
        this.dataType = dataType;
        this.numElements = numElements;
    }

    public static TestField of(String fieldString) throws PlcInvalidFieldException {
        Matcher matcher = ADDRESS_PATTERN.matcher(fieldString);
        if (matcher.matches()) {
            TestType type = TestType.valueOf(matcher.group("type"));
            String name = matcher.group("name");
            String dataTypeName = WordUtils.capitalizeFully(matcher.group("dataType"));
            int numElements = 1;
            if (matcher.group("numElements") != null) {
                numElements = Integer.parseInt(matcher.group("numElements"));
            }
            try {
                Class<?> dataType = Class.forName("java.lang." + dataTypeName);
                return new TestField(type, name, dataType, numElements);
            } catch (ClassNotFoundException e) {
                logger.error("Unsupported type: " + dataTypeName, e);
                throw new PlcInvalidFieldException("Unsupported type: " + dataTypeName);
            }
        }
        throw new PlcInvalidFieldException("Unable to parse address: " + fieldString);
    }

    static boolean matches(String fieldString) {
        return ADDRESS_PATTERN.matcher(fieldString).matches();
    }

    public TestType getType() {
        return type;
    }

    public String getName() {
        return name;
    }

    public Class<?> getDataType() {
        return dataType;
    }

    public int getNumElements() {
        return numElements;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (o == null || getClass() != o.getClass()) {
            return false;
        }
        TestField testField = (TestField) o;
        return numElements == testField.numElements &&
            type == testField.type &&
            Objects.equals(name, testField.name) &&
            Objects.equals(dataType, testField.dataType);
    }

    @Override
    public int hashCode() {
        return Objects.hash(type, name, dataType, numElements);
    }

    @Override
    public String toString() {
        return "TestField{" +
            "type=" + type +
            ", name='" + name + '\'' +
            ", dataType=" + dataType +
            ", numElements=" + numElements +
            '}';
    }

}
