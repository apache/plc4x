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
package org.apache.plc4x.edgent.mock;

import org.apache.plc4x.java.api.model.Address;

public class MockAddress implements Address {
  private final String address;
  
  public MockAddress(String address) {
    this.address = address;
  }
  
  public String getAddress() {
    return address;
  }
  
  @Override
  public String toString() {
    return "mock address: "+address;
  }
  
  @Override
  public boolean equals(Object o) {
    return o != null
        && o instanceof MockAddress
        && ((MockAddress)o).address.equals(this.address);
  }

  @Override
  public int hashCode() {
    return address.hashCode();
  }

}
