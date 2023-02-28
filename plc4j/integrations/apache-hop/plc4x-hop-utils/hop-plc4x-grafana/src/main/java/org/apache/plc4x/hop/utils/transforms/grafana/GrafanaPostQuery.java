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

public class GrafanaPostQuery {

  @HopMetadataProperty(injectionKeyDescription = "HTTPPOST.Injection.QueryFieldName")
  private String name;

  @HopMetadataProperty(injectionKeyDescription = "HTTPPOST.Injection.QueryFieldParameter")
  private String parameter;

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

  public GrafanaPostQuery(String name, String parameter) {
    this.name = name;
    this.parameter = parameter;
  }

  public GrafanaPostQuery(GrafanaPostQuery httpPostQuery) {
    this.name = httpPostQuery.name;
    this.parameter = httpPostQuery.parameter;
  }

  public GrafanaPostQuery() {}
}
