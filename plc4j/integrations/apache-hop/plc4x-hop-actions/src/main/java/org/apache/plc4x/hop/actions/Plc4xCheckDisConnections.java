/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package org.apache.plc4x.hop.actions;

import org.apache.hop.core.ICheckResult;
import org.apache.hop.core.Result;
import org.apache.hop.core.annotations.Action;
import org.apache.hop.core.exception.HopXmlException;
import org.apache.hop.core.variables.IVariables;
import org.apache.hop.core.xml.XmlHandler;
import org.apache.hop.i18n.BaseMessages;
import org.apache.hop.metadata.api.IHopMetadataProvider;
import org.apache.hop.workflow.WorkflowMeta;
import org.apache.hop.workflow.action.ActionBase;
import org.apache.hop.workflow.action.IAction;


import java.util.List;
import org.apache.hop.core.Const;
import org.apache.hop.core.exception.HopException;
import org.apache.plc4x.hop.metadata.Plc4xConnection;
import org.w3c.dom.Node;


@Action(
    id = "CHECK_PLC4X_DISCONNECTIONS",
    name = "i18n::Plc4xActionDisConnections.Name",
    description = "i18n::Plc4xActionDisConnections.Description",
    image = "plc4x_toddy.svg",
    categoryDescription = "i18n:org.apache.hop.workflow:ActionCategory.Category.Conditions",
    keywords = "i18n::Plc4xActionDisConnections.keyword",
    documentationUrl = "/workflow/actions/checkdbconnection.html")
public class Plc4xCheckDisConnections extends ActionBase implements Cloneable, IAction {
  private static final Class<?> PKG = Plc4xCheckDisConnections.class; // Needed by Translator

  
  private Plc4xConnection[] connections;  
  
  protected static final String[] unitTimeDesc =
      new String[] {
        BaseMessages.getString(PKG, "Plc4xCheckConnections.UnitTimeMilliSecond.Label"),
        BaseMessages.getString(PKG, "Plc4xCheckConnections.UnitTimeSecond.Label"),
        BaseMessages.getString(PKG, "Plc4xCheckConnections.UnitTimeMinute.Label"),
        BaseMessages.getString(PKG, "Plc4xCheckConnections.UnitTimeHour.Label"),
      };
  protected static final String[] unitTimeCode =
      new String[] {"millisecond", "second", "minute", "hour"};

  public static final int UNIT_TIME_MILLI_SECOND = 0;
  public static final int UNIT_TIME_SECOND = 1;
  public static final int UNIT_TIME_MINUTE = 2;
  public static final int UNIT_TIME_HOUR = 3;

  private String[] waitfors;
  private int[] waittimes;

  private long timeStart;
  private long now;  
     
  public Plc4xCheckDisConnections( String name) {
    super(name, "");
    //connections = null;
    waitfors = null;
    waittimes = null;
  }

  public Plc4xCheckDisConnections() {
    this( "");
  }
  
  public Object clone() {
    Plc4xCheckDisConnections c = (Plc4xCheckDisConnections) super.clone();
    return c;
  }


  public Plc4xConnection[] getConnections() {
    return connections;
  }

  public void setConnections(Plc4xConnection[] connections) {
    this.connections = connections;
  }

  public String[] getWaitfors() {
    return waitfors;
  }

  public void setWaitfors(String[] waitfors) {
    this.waitfors = waitfors;
  }

  public int[] getWaittimes() {
    return waittimes;
  }

  public void setWaittimes(int[] waittimes) {
    this.waittimes = waittimes;
  }

  public long getTimeStart() {
    return timeStart;
  }

  public long getNow() {
    return now;
  }

  private static String getWaitTimeCode(int i) {
    if (i < 0 || i >= unitTimeCode.length) {
      return unitTimeCode[0];
    }
    return unitTimeCode[i];
  }

  public static String getWaitTimeDesc(int i) {
    if (i < 0 || i >= unitTimeDesc.length) {
      return unitTimeDesc[0];
    }
    return unitTimeDesc[i];
  }

  public static int getWaitTimeByDesc(String tt) {
    if (tt == null) {
      return 0;
    }

    for (int i = 0; i < unitTimeDesc.length; i++) {
      if (unitTimeDesc[i].equalsIgnoreCase(tt)) {
        return i;
      }
    }

    // If this fails, try to match using the code.
    return getWaitTimeByCode(tt);
  }

  private static int getWaitTimeByCode(String tt) {
    if (tt == null) {
      return 0;
    }

    for (int i = 0; i < unitTimeCode.length; i++) {
      if (unitTimeCode[i].equalsIgnoreCase(tt)) {
        return i;
      }
    }
    return 0;
  }
  
  
  
  /**
   *
   * Save values to XML
   *
   * @return
   */
  @Override
  public String getXml() {
    StringBuilder xml = new StringBuilder(120);
    xml.append(super.getXml());
    xml.append("      <connections>").append(Const.CR);
    if (connections != null) {
      for (int i = 0; i < connections.length; i++) {
        xml.append("        <connection>").append(Const.CR);
        xml.append("          ")
            .append(
                XmlHandler.addTagValue(
                    "name", connections[i] == null ? null : connections[i].getName()));
        xml.append("          ").append(XmlHandler.addTagValue("waitfor", waitfors[i]));
        xml.append("          ")
            .append(XmlHandler.addTagValue("waittime", getWaitTimeCode(waittimes[i])));
        xml.append("        </connection>").append(Const.CR);
      }
    }
    xml.append("      </connections>").append(Const.CR);

    return xml.toString();
  }

  /**
   *
   * Read the XML and get the values needed for the acton
   *
   * @param entrynode
   * @param metadataProvider
   * @throws HopXmlException
   */
  @Override
  public void loadXml( Node entrynode, IHopMetadataProvider metadataProvider, IVariables variables ) throws HopXmlException {
    try {
      super.loadXml(entrynode);
      Node fields = XmlHandler.getSubNode(entrynode, "connections");

      // How many hosts?
      int nrFields = XmlHandler.countNodes(fields, "connection");
      connections = new Plc4xConnection[nrFields];
      waitfors = new String[nrFields];
      waittimes = new int[nrFields];
      // Read them all...
      for (int i = 0; i < nrFields; i++) {
        Node fnode = XmlHandler.getSubNodeByNr(fields, "connection", i);
        String dbname = XmlHandler.getTagValue(fnode, "name");
        
        //connections[i] = Plc4xConnection.loadDatabase(metadataProvider, dbname);
        if (dbname != null) {
            connections[i] = metadataProvider.getSerializer(Plc4xConnection.class).load(dbname);
            waitfors[i] = XmlHandler.getTagValue(fnode, "waitfor");
            waittimes[i] = getWaitTimeByCode(Const.NVL(XmlHandler.getTagValue(fnode, "waittime"), ""));
        };
      }
    } catch (HopXmlException xe) {
      throw new HopXmlException(
          BaseMessages.getString(
              PKG,
              "Plc4xCheckConnections.ERROR_0001_Cannot_Load_Job_Entry_From_Xml_Node",
              xe.getMessage()));
    } catch (HopException ex) {
       throw new HopXmlException(
          BaseMessages.getString(
              PKG,
              "Plc4xCheckConnections.ERROR_0001_Cannot_Load_Job_Entry_From_Xml_Node",
              ex.getMessage()));
      }
  }

  /**
   * Execute this action and return the result. In this case it means, just set the result boolean in the Result
   * class.
   *
   * @param result The result of the previous execution
   * @return The Result of the execution.
   */
  @Override
  public Result execute( Result result, int nr ) {
      result.setResult(true);
      System.out.println("NR: " + nr);

    return result;
  }

  /**
   *
   * Add checks to report warnings
   *
   * @param remarks
   * @param workflowMeta
   * @param variables
   * @param metadataProvider
   */
  @Override
  public void check( List<ICheckResult> remarks, WorkflowMeta workflowMeta, IVariables variables,
                     IHopMetadataProvider metadataProvider ) {
  }
}
