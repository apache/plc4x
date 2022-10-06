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
package org.apache.plc4x.java.bacnetip.ede.layouts;

public interface EdeLayout {

    default int getKeyNamePos() { return -1; }

    default int getDeviceInstancePos() { return -1; }

    default int getObjectNamePos() { return -1; }

    default int getObjectTypePos() { return -1; }

    default int getObjectInstancePos() { return -1; }

    default int getDescriptionPos() { return -1; }

    default int getDefaultValuePos() { return -1; }

    default int getMinValuePos() { return -1; }

    default int getMaxValuePos() { return -1; }

    default int getCommandablePos() { return -1; }

    default int getSupportsCovPos() { return -1; }

    default int getHiLimitPos() { return -1; }

    default int getLowLimitPos() { return -1; }

    default int getStateTextReferencePos() { return -1; }

    default int getUnitCodePos() { return -1; }

    default int getVendorSpecificAddressPos() { return -1; }

    default int getNotificationClassPos() { return -1; }

}
