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
package org.apache.plc4x.test.driver;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.dataformat.xml.XmlMapper;
import org.apache.plc4x.test.driver.model.api.TestField;
import org.apache.plc4x.test.driver.model.api.TestReadRequest;
import org.apache.plc4x.test.driver.model.api.TestRequest;

public class Lalala {
    public static void main(String[] args) throws Exception {
        ObjectMapper mapper = new XmlMapper().enableDefaultTyping();

        TestField[] testFields = new TestField[1];
        testFields[0] = new TestField("hurz", "%Q0.0:BOOL");
        TestRequest request = new TestReadRequest(testFields);

        String xmlString = mapper.writerWithDefaultPrettyPrinter().writeValueAsString(request);

        System.out.println(xmlString);

        TestRequest newRequest = mapper.readValue(xmlString, TestRequest.class);

        System.out.println(newRequest);
    }

}
