package org.apache.plc4x.edgent.wip;

import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.connection.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcException;

/**
 * PlcConnectionAdapter encapsulates a PlcConnection.
 * <p>
 * The idea here is to use PlcConnectionAdapter to enable our Edgent Supplier/Consumer
 * instances to be isolated from some of the details of / variability of 
 * PlcConnection mgmt and such.
 * <p>
 * The current usage presumption (in Supplier and Consumer factories) is that
 * a single PlcConnector instance can be used in multiple Supplier/Consumer instances.
 * IS THAT VALID?
 */
public class PlcConnectionAdapter implements AutoCloseable{
  
  private String plcConnectionUrl;
  private PlcConnection plcConnection;
  
  // plcConnection must be connected
  public PlcConnectionAdapter(PlcConnection plcConnection) {
    this.plcConnection = plcConnection;
  }
  
  // getConnection() gets a PlcConnection and connects.
  public PlcConnectionAdapter(String plcConnectionUrl) {
    this.plcConnectionUrl = plcConnectionUrl;
  }
  
  // returns connected PlcConnection
  public PlcConnection getConnection() throws PlcException {
    synchronized(this) {
      if (plcConnection == null) {
        plcConnection = new PlcDriverManager().getConnection(plcConnectionUrl);
        plcConnection.connect();
      }
      return plcConnection;
    }
  }

  @Override
  public void close() throws Exception {
    // only close a connection this instance created/connected
    if (plcConnectionUrl != null) {
      if (plcConnection != null)
        plcConnection.close();
    }
  }

}
