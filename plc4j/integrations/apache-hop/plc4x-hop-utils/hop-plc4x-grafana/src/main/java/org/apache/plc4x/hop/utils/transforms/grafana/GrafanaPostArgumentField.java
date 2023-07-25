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

package org.apache.plc4x.hop.utils.transforms.grafana;

import org.apache.hop.metadata.api.HopMetadataProperty;

public class GrafanaPostArgumentField {

  @HopMetadataProperty(injectionKeyDescription = "HTTPPOST.Injection.ArgumentFieldName")
  private String name;

  @HopMetadataProperty(injectionKeyDescription = "HTTPPOST.Injection.ArgumentFieldParameter")
  private String parameter;

  @HopMetadataProperty(injectionKeyDescription = "HTTPPOST.Injection.ArgumentFieldHeader")
  private boolean header;

  public String getName() {
    return name;
  }

  public void setName(String name) {
    this.name = name;
  }

  public String getParameter() {
    return parameter;
  }

  public void setParameter(String parameter) {
    this.parameter = parameter;
  }

  public boolean isHeader() {
    return header;
  }

  public void setHeader(boolean header) {
    this.header = header;
  }

  public GrafanaPostArgumentField(String name, String parameter, boolean header) {
    this.name = name;
    this.parameter = parameter;
    this.header = header;
  }

  public GrafanaPostArgumentField(GrafanaPostArgumentField httpPostArgumentField) {
    this.name = httpPostArgumentField.name;
    this.parameter = httpPostArgumentField.parameter;
    this.header = httpPostArgumentField.header;
  }

  public GrafanaPostArgumentField() {}
}
