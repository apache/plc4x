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
<<<<<<<< HEAD:plc4j/drivers/eip/src/main/java/org/apache/plc4x/java/eip/base/field/EipFieldHandler.java
package org.apache.plc4x.java.eip.base.field;
========
package org.apache.plc4x.java.mock.tag;
>>>>>>>> develop:plc4j/drivers/mock/src/main/java/org/apache/plc4x/java/mock/tag/MockTagHandler.java

import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.model.PlcQuery;
import org.apache.plc4x.java.spi.connection.PlcTagHandler;

public class MockTagHandler implements PlcTagHandler {

    @Override
    public PlcTag parseTag(String tagAddress) {
        return new MockTag(tagAddress);
    }

    @Override
    public PlcQuery parseQuery(String query) {
        throw new UnsupportedOperationException("This driver doesn't support browsing");
    }

}
