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

import java.util.ArrayList;
import java.util.List;

public class GrafanaPostLoookupField {

  @HopMetadataProperty(
      key = "query",
      injectionGroupDescription = "HTTPPOST.Injection.LookupQueryField")
  private List<GrafanaPostQuery> queryField = new ArrayList<>();

  @HopMetadataProperty(
      key = "arg",
      injectionGroupDescription = "HTTPPOST.Injection.LookupArgumentField")
  private List<GrafanaPostArgumentField> argumentField = new ArrayList<>();

  public List<GrafanaPostQuery> getQueryField() {
    return queryField;
  }

  public void setQueryField(List<GrafanaPostQuery> queryField) {
    this.queryField = queryField;
  }

  public List<GrafanaPostArgumentField> getArgumentField() {
    return argumentField;
  }

  public void setArgumentField(List<GrafanaPostArgumentField> argumentField) {
    this.argumentField = argumentField;
  }

  public GrafanaPostLoookupField(
      List<GrafanaPostQuery> postQuery, List<GrafanaPostArgumentField> argumentField) {
    this.queryField = postQuery;
    this.argumentField = argumentField;
  }

  public GrafanaPostLoookupField(GrafanaPostLoookupField httpPostLoookupField) {
    this.queryField = httpPostLoookupField.queryField;
    this.argumentField = httpPostLoookupField.argumentField;
  }

  public GrafanaPostLoookupField() {}
}
