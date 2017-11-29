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
