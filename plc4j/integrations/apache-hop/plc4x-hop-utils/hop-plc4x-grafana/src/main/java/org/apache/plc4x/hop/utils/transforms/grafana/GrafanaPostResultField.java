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

public class GrafanaPostResultField {

  @HopMetadataProperty(injectionKeyDescription = "HTTPPOST.Injection.ResultFieldCode")
  private String code;

  @HopMetadataProperty(injectionKeyDescription = "HTTPPOST.Injection.ResultFieldName")
  private String name;

  @HopMetadataProperty(
      key = "response_time",
      injectionKeyDescription = "HTTPPOST.Injection.ResultFieldResponseTime")
  private String responseTimeFieldName;

  @HopMetadataProperty(
      key = "response_header",
      injectionKeyDescription = "HTTPPOST.Injection.ResultFieldResponseHeader")
  private String responseHeaderFieldName;

  public String getCode() {
    return code;
  }

  public void setCode(String code) {
    this.code = code;
  }

  public String getName() {
    return name;
  }

  public void setName(String name) {
    this.name = name;
  }

  public String getResponseTimeFieldName() {
    return responseTimeFieldName;
  }

  public void setResponseTimeFieldName(String responseTimeFieldName) {
    this.responseTimeFieldName = responseTimeFieldName;
  }

  public String getResponseHeaderFieldName() {
    return responseHeaderFieldName;
  }

  public void setResponseHeaderFieldName(String responseHeaderFieldName) {
    this.responseHeaderFieldName = responseHeaderFieldName;
  }

  public GrafanaPostResultField(
      String code, String name, String responseTimeFieldName, String responseHeaderFieldName) {
    this.code = code;
    this.name = name;
    this.responseTimeFieldName = responseTimeFieldName;
    this.responseHeaderFieldName = responseHeaderFieldName;
  }

  public GrafanaPostResultField(GrafanaPostResultField httpPostResultField) {
    this.code = httpPostResultField.code;
    this.name = httpPostResultField.name;
    this.responseTimeFieldName = httpPostResultField.responseTimeFieldName;
    this.responseHeaderFieldName = httpPostResultField.responseHeaderFieldName;
  }

  public GrafanaPostResultField() {
    this.code = "result";
  }
}
