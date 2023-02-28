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
package org.apache.plc4x.hop.metadata;

import org.apache.hop.core.gui.plugin.GuiElementType;
import org.apache.hop.core.gui.plugin.GuiPlugin;
import org.apache.hop.core.gui.plugin.GuiWidgetElement;
import org.apache.hop.metadata.api.HopMetadata;
import org.apache.hop.metadata.api.HopMetadataBase;
import org.apache.hop.metadata.api.HopMetadataProperty;
import org.apache.hop.metadata.api.IHopMetadata;

@GuiPlugin
@HopMetadata(
    key = "plc4x-connection",
    name = "Plc4x Connection",
    description = "Describes a connection to a PLC",
    image = "plc4x_toddy.svg",
    documentationUrl = "/metadata-types/plc4x/plc4x-connection.html")
public class Plc4xConnection extends HopMetadataBase implements IHopMetadata {

  public static final String WIDGET_ID_URL = "10000-url";    
  
  @HopMetadataProperty
  @GuiWidgetElement(
      id = WIDGET_ID_URL,
      type = GuiElementType.TEXT,
      parentId = Plc4xConnectionEditor.PARENT_WIDGET_ID,
      label = "URL",
      toolTip = "Specify the hostname of your cassandra server")
  private String url;  
  
  
    public Plc4xConnection() {}
  
  public Plc4xConnection(Plc4xConnection p) {
    super(p.name);
    this.url = p.url;
  }  
  
  
  /**
   * Gets url
   *
   * @return value of hostname
   */
  public String getUrl() {
    return url;
  }

  /** @param hostname The hostname to set */
  public void setUrl(String url) {
    this.url = url;
  }  
    
    
}
