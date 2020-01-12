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
package org.apache.plc4x.java.bacnetip.protocol;

import org.apache.plc4x.java.bacnetip.configuration.PassiveBacNetIpConfiguration;
import org.apache.plc4x.java.bacnetip.readwrite.BVLC;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.configuration.HasConfiguration;

public class PassiveBacNetIpProtocolLogic extends Plc4xProtocolBase<BVLC> implements HasConfiguration<PassiveBacNetIpConfiguration> {

    private PassiveBacNetIpConfiguration configuration;

    @Override
    public void setConfiguration(PassiveBacNetIpConfiguration configuration) {
        this.configuration = configuration;
    }

    @Override
    public void close(ConversationContext<BVLC> context) {
        // Nothing to do here ...
    }

    @Override
    protected void decode(ConversationContext<BVLC> context, BVLC msg) throws Exception {
        System.out.println(msg);
        super.decode(context, msg);
    }

}
