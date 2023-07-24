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

import org.apache.hop.core.CheckResult;
import org.apache.hop.core.ICheckResult;
import org.apache.hop.core.annotations.Transform;
import org.apache.hop.core.exception.HopTransformException;
import org.apache.hop.core.row.IRowMeta;
import org.apache.hop.core.row.IValueMeta;
import org.apache.hop.core.row.value.ValueMetaInteger;
import org.apache.hop.core.row.value.ValueMetaString;
import org.apache.hop.core.util.Utils;
import org.apache.hop.core.variables.IVariables;
import org.apache.hop.i18n.BaseMessages;
import org.apache.hop.metadata.api.HopMetadataProperty;
import org.apache.hop.metadata.api.IHopMetadataProvider;
import org.apache.hop.pipeline.Pipeline;
import org.apache.hop.pipeline.PipelineMeta;
import org.apache.hop.pipeline.transform.BaseTransformMeta;
import org.apache.hop.pipeline.transform.ITransformMeta;
import org.apache.hop.pipeline.transform.TransformMeta;

import java.util.ArrayList;
import java.util.List;

/*
* This code is taken from the "HttpPost" transform to be adapted for use 
* with the Grafana APIs.
* In the future, its integration with the original version should be proposed
* to the working group of Apache Hop
*/

@Transform(
    id = "grafana-post",
    image = "grafana.svg",
    name = "i18n::BaseTransform.TypeLongDesc.GrafanaPost",
    description = "i18n::BaseTransform.TypeTooltipDesc.HTTPPOST",
    categoryDescription = "i18n:org.apache.plc4x.hop.transforms.plc4xinput.messages:Plc4x.Category.plc4x",
    keywords = "i18n::HttpPostMeta.keyword",
    documentationUrl = "/pipeline/transforms/httppost.html")
public class GrafanaPostMeta extends BaseTransformMeta<GrafanaPost, GrafanaPostData> implements ITransformMeta {
  private static final Class<?> PKG = GrafanaPostMeta.class; // For Translator

  // the timeout for waiting for data (milliseconds)
  public static final int DEFAULT_SOCKET_TIMEOUT = 10000;

  // the timeout until a connection is established (milliseconds)
  public static final int DEFAULT_CONNECTION_TIMEOUT = 10000;

  // the time to wait till a connection is closed (milliseconds)? -1 is no not close.
  public static final int DEFAULT_CLOSE_CONNECTIONS_TIME = -1;

  public static final String DEFAULT_ENCODING = "UTF-8";

  @HopMetadataProperty(injectionKeyDescription = "HTTPPOST.Injection.socketTimeout")
  private String socketTimeout;

  @HopMetadataProperty(injectionKeyDescription = "HTTPPOST.Injection.connectionTimeout")
  private String connectionTimeout;

  @HopMetadataProperty(injectionKeyDescription = "HTTPPOST.Injection.closeIdleConnectionsTime")
  private String closeIdleConnectionsTime;

  /** URL / service to be called */
  @HopMetadataProperty(injectionKeyDescription = "HTTPPOST.Injection.url")
  private String url;

  @HopMetadataProperty(key = "lookup", injectionGroupDescription = "HTTPPOST.Injection.lookupfield")
  private List<GrafanaPostLoookupField> lookupfield = new ArrayList<>();

  @HopMetadataProperty(injectionKeyDescription = "HTTPPOST.Injection.urlInField")
  private boolean urlInField;

  @HopMetadataProperty(injectionKeyDescription = "HTTPPOST.Injection.urlField")
  private String urlField;

  @HopMetadataProperty(injectionKeyDescription = "HTTPPOST.Injection.requestEntity")
  private String requestEntity;

  @HopMetadataProperty(injectionKeyDescription = "HTTPPOST.Injection.encoding")
  private String encoding;

  @HopMetadataProperty(key = "postafile", injectionKeyDescription = "HTTPPOST.Injection.postAFile")
  private boolean postAFile;

  @HopMetadataProperty(injectionKeyDescription = "HTTPPOST.Injection.proxyHost")
  private String proxyHost;

  @HopMetadataProperty(injectionKeyDescription = "HTTPPOST.Injection.proxyPort")
  private String proxyPort;

  @HopMetadataProperty(injectionKeyDescription = "HTTPPOST.Injection.httpLogin")
  private String httpLogin;

  @HopMetadataProperty(password = true, injectionKeyDescription = "HTTPPOST.Injection.httpPassword")
  private String httpPassword;

  @HopMetadataProperty(
      key = "result",
      injectionGroupDescription = "HTTPPOST.Injection.httpPostResultField")
  private List<GrafanaPostResultField> httpPostResultField = new ArrayList<>();

  public GrafanaPostMeta() {
    super(); // allocate BaseTransformMeta
  }

  public String getEncoding() {
    return encoding;
  }

  public void setEncoding(String encoding) {
    this.encoding = encoding;
  }

  /** @return Returns the connectionTimeout. */
  public String getConnectionTimeout() {
    return connectionTimeout;
  }

  /** @param connectionTimeout The connectionTimeout to set. */
  public void setConnectionTimeout(String connectionTimeout) {
    this.connectionTimeout = connectionTimeout;
  }

  /** @return Returns the closeIdleConnectionsTime. */
  public String getCloseIdleConnectionsTime() {
    return closeIdleConnectionsTime;
  }

  /** @param closeIdleConnectionsTime The connectionTimeout to set. */
  public void setCloseIdleConnectionsTime(String closeIdleConnectionsTime) {
    this.closeIdleConnectionsTime = closeIdleConnectionsTime;
  }

  /** @return Returns the socketTimeout. */
  public String getSocketTimeout() {
    return socketTimeout;
  }

  /** @param socketTimeout The socketTimeout to set. */
  public void setSocketTimeout(String socketTimeout) {
    this.socketTimeout = socketTimeout;
  }

  /** @return Returns the procedure. */
  public String getUrl() {
    return url;
  }

  /** @param procedure The procedure to set. */
  public void setUrl(String procedure) {
    this.url = procedure;
  }

  /** @return Is the url coded in a field? */
  public boolean isUrlInField() {
    return urlInField;
  }

  public boolean isPostAFile() {
    return postAFile;
  }

  public void setPostAFile(boolean postafile) {
    this.postAFile = postafile;
  }

  /** @param urlInField Is the url coded in a field? */
  public void setUrlInField(boolean urlInField) {
    this.urlInField = urlInField;
  }

  /** @return The field name that contains the url. */
  public String getUrlField() {
    return urlField;
  }

  /** @param urlField name of the field that contains the url */
  public void setUrlField(String urlField) {
    this.urlField = urlField;
  }

  /** @param requestEntity the requestEntity to set */
  public void setRequestEntity(String requestEntity) {
    this.requestEntity = requestEntity;
  }

  /** @return requestEntity */
  public String getRequestEntity() {
    return requestEntity;
  }

  public List<GrafanaPostLoookupField> getLookupfield() {
    return lookupfield;
  }

  public void setLookupfield(List<GrafanaPostLoookupField> lookupfield) {
    this.lookupfield = lookupfield;
  }

  public List<GrafanaPostResultField> getHttpPostResultField() {
    return httpPostResultField;
  }

  public void setHttpPostResultField(List<GrafanaPostResultField> httpPostResultField) {
    this.httpPostResultField = httpPostResultField;
  }

  @Override
  public Object clone() {
    GrafanaPostMeta retval = (GrafanaPostMeta) super.clone();

    return retval;
  }

  @Override
  public void setDefault() {
    encoding = DEFAULT_ENCODING;
    postAFile = false;
    lookupfield.add(new GrafanaPostLoookupField());
    httpPostResultField.add(new GrafanaPostResultField());
    socketTimeout = String.valueOf(DEFAULT_SOCKET_TIMEOUT);
    connectionTimeout = String.valueOf(DEFAULT_CONNECTION_TIMEOUT);
    closeIdleConnectionsTime = String.valueOf(DEFAULT_CLOSE_CONNECTIONS_TIME);
  }

  @Override
  public void getFields(
      IRowMeta inputRowMeta,
      String name,
      IRowMeta[] info,
      TransformMeta nextTransform,
      IVariables variables,
      IHopMetadataProvider metadataProvider)
      throws HopTransformException {
    if (!Utils.isEmpty(httpPostResultField.get(0).getName())) {
      IValueMeta v = new ValueMetaString(httpPostResultField.get(0).getName());
      inputRowMeta.addValueMeta(v);
    }

    if (!Utils.isEmpty(httpPostResultField.get(0).getCode())) {
      IValueMeta v = new ValueMetaInteger(httpPostResultField.get(0).getCode());
      inputRowMeta.addValueMeta(v);
    }
    if (!Utils.isEmpty(httpPostResultField.get(0).getResponseTimeFieldName())) {
      IValueMeta v =
          new ValueMetaInteger(
              variables.resolve(httpPostResultField.get(0).getResponseTimeFieldName()));
      inputRowMeta.addValueMeta(v);
    }
    String headerFieldName =
        variables.resolve(httpPostResultField.get(0).getResponseHeaderFieldName());
    if (!Utils.isEmpty(headerFieldName)) {
      IValueMeta v = new ValueMetaString(headerFieldName);
      v.setOrigin(name);
      inputRowMeta.addValueMeta(v);
    }
  }

  @Override
  public void check(
      List<ICheckResult> remarks,
      PipelineMeta pipelineMeta,
      TransformMeta transformMeta,
      IRowMeta prev,
      String[] input,
      String[] output,
      IRowMeta info,
      IVariables variables,
      IHopMetadataProvider metadataProvider) {
    CheckResult cr;

    // See if we have input streams leading to this transform!
    if (input.length > 0) {
      cr =
          new CheckResult(
              ICheckResult.TYPE_RESULT_OK,
              BaseMessages.getString(
                  PKG, "HTTPPOSTMeta.CheckResult.ReceivingInfoFromOtherTransforms"),
              transformMeta);
      remarks.add(cr);
    } else {
      cr =
          new CheckResult(
              ICheckResult.TYPE_RESULT_ERROR,
              BaseMessages.getString(PKG, "HTTPPOSTMeta.CheckResult.NoInpuReceived"),
              transformMeta);
      remarks.add(cr);
    }

    // check Url
    if (urlInField) {
      if (Utils.isEmpty(urlField)) {
        cr =
            new CheckResult(
                ICheckResult.TYPE_RESULT_ERROR,
                BaseMessages.getString(PKG, "HTTPPOSTMeta.CheckResult.UrlfieldMissing"),
                transformMeta);
      } else {
        cr =
            new CheckResult(
                ICheckResult.TYPE_RESULT_ERROR,
                BaseMessages.getString(PKG, "HTTPPOSTMeta.CheckResult.UrlfieldOk"),
                transformMeta);
      }

    } else {
      if (Utils.isEmpty(url)) {
        cr =
            new CheckResult(
                ICheckResult.TYPE_RESULT_ERROR,
                BaseMessages.getString(PKG, "HTTPPOSTMeta.CheckResult.UrlMissing"),
                transformMeta);
      } else {
        cr =
            new CheckResult(
                ICheckResult.TYPE_RESULT_OK,
                BaseMessages.getString(PKG, "HTTPPOSTMeta.CheckResult.UrlOk"),
                transformMeta);
      }
    }
    remarks.add(cr);
  }

  public GrafanaPost createTransform(
      TransformMeta transformMeta,
      GrafanaPostData data,
      int cnr,
      PipelineMeta pipelineMeta,
      Pipeline pipeline) {
    return new GrafanaPost(transformMeta, this, data, cnr, pipelineMeta, pipeline);
  }

  public GrafanaPostData getTransformData() {
    return new GrafanaPostData();
  }

  @Override
  public boolean supportsErrorHandling() {
    return true;
  }

  /**
   * ISetter
   *
   * @param proxyHost
   */
  public void setProxyHost(String proxyHost) {
    this.proxyHost = proxyHost;
  }

  /**
   * IGetter
   *
   * @return
   */
  public String getProxyHost() {
    return proxyHost;
  }

  /**
   * ISetter
   *
   * @param proxyPort
   */
  public void setProxyPort(String proxyPort) {
    this.proxyPort = proxyPort;
  }

  /**
   * IGetter
   *
   * @return
   */
  public String getProxyPort() {
    return this.proxyPort;
  }

  /**
   * ISetter
   *
   * @param httpLogin
   */
  public void setHttpLogin(String httpLogin) {
    this.httpLogin = httpLogin;
  }

  /**
   * IGetter
   *
   * @return
   */
  public String getHttpLogin() {
    return httpLogin;
  }

  /**
   * ISetter
   *
   * @param httpPassword
   */
  public void setHttpPassword(String httpPassword) {
    this.httpPassword = httpPassword;
  }

  /** @return */
  public String getHttpPassword() {
    return httpPassword;
  }
}
