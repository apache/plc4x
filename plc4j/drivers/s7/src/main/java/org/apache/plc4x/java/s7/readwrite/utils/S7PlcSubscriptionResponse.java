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
package org.apache.plc4x.java.s7.readwrite.utils;

import java.util.Map;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.spi.messages.DefaultPlcSubscriptionResponse;
import org.apache.plc4x.java.spi.messages.utils.ResponseItem;

/**
 *
 * @author cgarcia
 */
public class S7PlcSubscriptionResponse extends DefaultPlcSubscriptionResponse {
    
        private final short jobid;
    
    public S7PlcSubscriptionResponse(PlcSubscriptionRequest request,
            short jobid,
            Map<String, ResponseItem<PlcSubscriptionHandle>> values) {
        super(request, values);
        this.jobid = jobid;
    }
    
    public short getJobId(){
        return jobid;
    }
    
    
}
