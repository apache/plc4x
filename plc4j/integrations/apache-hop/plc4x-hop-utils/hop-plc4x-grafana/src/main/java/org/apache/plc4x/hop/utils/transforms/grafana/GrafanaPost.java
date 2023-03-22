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

import com.google.common.annotations.VisibleForTesting;
import org.apache.commons.lang.StringUtils;
import org.apache.hop.core.Const;
import org.apache.hop.core.exception.HopException;
import org.apache.hop.core.exception.HopTransformException;
import org.apache.hop.core.row.RowDataUtil;
import org.apache.hop.core.util.HttpClientManager;
import org.apache.hop.core.util.StringUtil;
import org.apache.hop.core.util.Utils;
import org.apache.hop.i18n.BaseMessages;
import org.apache.hop.pipeline.Pipeline;
import org.apache.hop.pipeline.PipelineMeta;
import org.apache.hop.pipeline.transform.BaseTransform;
import org.apache.hop.pipeline.transform.ITransform;
import org.apache.hop.pipeline.transform.TransformMeta;
import org.apache.http.*;
import org.apache.http.client.AuthCache;
import org.apache.http.client.entity.UrlEncodedFormEntity;
import org.apache.http.client.methods.CloseableHttpResponse;
import org.apache.http.client.protocol.HttpClientContext;
import org.apache.http.client.utils.URIBuilder;
import org.apache.http.entity.ByteArrayEntity;
import org.apache.http.entity.ContentType;
import org.apache.http.entity.InputStreamEntity;
import org.apache.http.entity.StringEntity;
import org.apache.http.impl.auth.BasicScheme;
import org.apache.http.impl.client.BasicAuthCache;
import org.apache.http.impl.client.CloseableHttpClient;
import org.apache.http.message.BasicNameValuePair;
import org.apache.http.util.EntityUtils;
import org.json.simple.JSONObject;

import java.io.File;
import java.io.FileInputStream;
import java.io.InputStreamReader;
import java.io.UnsupportedEncodingException;
import java.net.HttpURLConnection;
import java.net.URLEncoder;
import java.net.UnknownHostException;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;

import static org.apache.plc4x.hop.utils.transforms.grafana.GrafanaPostMeta.DEFAULT_ENCODING;

/** Make a HTTP Post call */
public class GrafanaPost extends BaseTransform<GrafanaPostMeta, GrafanaPostData>  implements ITransform {

  private static final Class<?> PKG = GrafanaPostMeta.class; // For Translator

  private static final String CONTENT_TYPE = "Content-type";
  private static final String CONTENT_TYPE_TEXT_XML = "text/xml";
  private static final String PKG_HEADER_VALUE = "HTTPPOST.Log.HeaderValue";
  private static final String PKG_ERROR_FINDING_FIELD = "HTTPPOST.Log.ErrorFindingField";

  public GrafanaPost(
      TransformMeta transformMeta,
      GrafanaPostMeta meta,
      GrafanaPostData data,
      int copyNr,
      PipelineMeta pipelineMeta,
      Pipeline pipeline) {
    super(transformMeta, meta, data, copyNr, pipelineMeta, pipeline);
  }

  private Object[] callHttpPOST(Object[] rowData) throws HopException {
    HttpClientManager.HttpClientBuilderFacade clientBuilder =
        HttpClientManager.getInstance().createBuilder();

    if (data.realConnectionTimeout > -1) {
      clientBuilder.setConnectionTimeout(data.realConnectionTimeout);
    }
    if (data.realSocketTimeout > -1) {
      clientBuilder.setSocketTimeout(data.realSocketTimeout);
    }
    if (StringUtils.isNotBlank(data.realHttpLogin)) {
      clientBuilder.setCredentials(data.realHttpLogin, data.realHttpPassword);
    }
    if (StringUtils.isNotBlank(data.realProxyHost)) {
      clientBuilder.setProxy(data.realProxyHost, data.realProxyPort);
    }

    CloseableHttpClient httpClient = clientBuilder.build();

    // get dynamic url ?
    if (meta.isUrlInField()) {
      data.realUrl = data.inputRowMeta.getString(rowData, data.indexOfUrlField);
    }
    // Prepare HTTP POST
    FileInputStream fis = null;
    try {
      if (isDetailed()) {
        logDetailed(BaseMessages.getString(PKG, "HTTPPOST.Log.ConnectingToURL", data.realUrl));
      }
      URIBuilder uriBuilder = new URIBuilder(data.realUrl);
      org.apache.http.client.methods.HttpPost post =
          new org.apache.http.client.methods.HttpPost(uriBuilder.build());

      // Specify content type and encoding
      // If content encoding is not explicitly specified
      // ISO-8859-1 is assumed by the POSTMethod
      if (!data.contentTypeHeaderOverwrite) { // can be overwritten now
        if (Utils.isEmpty(data.realEncoding)) {
          post.setHeader(CONTENT_TYPE, CONTENT_TYPE_TEXT_XML);
          if (isDebug()) {
            logDebug(
                BaseMessages.getString(PKG, PKG_HEADER_VALUE, CONTENT_TYPE, CONTENT_TYPE_TEXT_XML));
          }
        } else {
          post.setHeader(CONTENT_TYPE, CONTENT_TYPE_TEXT_XML + "; " + data.realEncoding);
          if (isDebug()) {
            logDebug(
                BaseMessages.getString(
                    PKG,
                    PKG_HEADER_VALUE,
                    CONTENT_TYPE,
                    CONTENT_TYPE_TEXT_XML + "; " + data.realEncoding));
          }
        }
      }

      // HEADER PARAMETERS
      if (data.useHeaderParameters) {
        // set header parameters that we want to send
        for (int i = 0; i < data.header_parameters_nrs.length; i++) {
          post.addHeader(
              data.headerParameters[i].getName(),
              data.inputRowMeta.getString(rowData, data.header_parameters_nrs[i]));
          if (isDebug()) {
            logDebug(
                BaseMessages.getString(
                    PKG,
                    PKG_HEADER_VALUE,
                    data.headerParameters[i].getName(),
                    data.inputRowMeta.getString(rowData, data.header_parameters_nrs[i])));
          }
        }
      }

      // BODY PARAMETERS
      if (data.useBodyParameters) {
        // set body parameters that we want to send
        for (int i = 0; i < data.body_parameters_nrs.length; i++) {
          String bodyParameterName = data.bodyParameters[i].getName();
          String bodyParameterValue =
              data.inputRowMeta.getString(rowData, data.body_parameters_nrs[i]);
          data.bodyParameters[i] = new BasicNameValuePair(bodyParameterName, bodyParameterValue);
          if (isDebug()) {
            logDebug(
                BaseMessages.getString(
                    PKG, "HTTPPOST.Log.BodyValue", bodyParameterName, bodyParameterValue));
          }
        }
        String bodyParams = getRequestBodyParamsAsStr(data.bodyParameters, data.realEncoding);
        post.setEntity(
            (new StringEntity(bodyParams, ContentType.TEXT_XML.withCharset("US-ASCII"))));
      }

      // QUERY PARAMETERS
      if (data.useQueryParameters) {
        for (int i = 0; i < data.query_parameters_nrs.length; i++) {
          String queryParameterName = data.queryParameters[i].getName();
          String queryParameterValue =
              data.inputRowMeta.getString(rowData, data.query_parameters_nrs[i]);
          data.queryParameters[i] = new BasicNameValuePair(queryParameterName, queryParameterValue);
          if (isDebug()) {
            logDebug(
                BaseMessages.getString(
                    PKG, "HTTPPOST.Log.QueryValue", queryParameterName, queryParameterValue));
          }
        }
        post.setEntity(new UrlEncodedFormEntity(Arrays.asList(data.queryParameters)));
      }

      // Set request entity?
      if (data.indexOfRequestEntity >= 0) {
        String tmp = data.inputRowMeta.getString(rowData, data.indexOfRequestEntity);
        // Request content will be retrieved directly
        // from the input stream
        // Per default, the request content needs to be buffered
        // in order to determine its length.
        // Request body buffering can be avoided when
        // content length is explicitly specified

        if (meta.isPostAFile()) {
          File input = new File(tmp);
          fis = new FileInputStream(input);
          post.setEntity(new InputStreamEntity(fis, input.length()));
        } else {
          byte[] bytes;
          if ((data.realEncoding != null) && (data.realEncoding.length() > 0)) {
            bytes = tmp.getBytes(data.realEncoding);
          } else {
            bytes = tmp.getBytes();
          }
          post.setEntity(new ByteArrayEntity(bytes));
        }
      }

      // Execute request
      Object[] newRow = null;
      if (rowData != null) {
        newRow = rowData.clone();
      }
      CloseableHttpResponse httpResponse = null;
      try {
        // used for calculating the responseTime
        long startTime = System.currentTimeMillis();

        // Execute the POST method
        if (StringUtils.isNotBlank(data.realProxyHost)) {
          HttpHost target = new HttpHost(data.realProxyHost, data.realProxyPort, "http");
          // Create AuthCache instance
          AuthCache authCache = new BasicAuthCache();
          // Generate BASIC scheme object and add it to the local
          // auth cache
          BasicScheme basicAuth = new BasicScheme();
          authCache.put(target, basicAuth);
          // Add AuthCache to the execution context
          HttpClientContext localContext = HttpClientContext.create();
          localContext.setAuthCache(authCache);
          httpResponse = httpClient.execute(target, post, localContext);
        } else {
          httpResponse = httpClient.execute(post);
        }
        int statusCode = requestStatusCode(httpResponse);

        // calculate the responseTime
        long responseTime = System.currentTimeMillis() - startTime;

        if (isDetailed()) {
          logDetailed(
              BaseMessages.getString(PKG, "HTTPPOST.Log.ResponseTime", responseTime, data.realUrl));
        }

        // Display status code
        if (isDebug()) {
          logDebug(
              BaseMessages.getString(PKG, "HTTPPOST.Log.ResponseCode", String.valueOf(statusCode)));
        }

        String body;
        String headerString = "";
        switch (statusCode) {
          case HttpURLConnection.HTTP_UNAUTHORIZED:
            throw new HopTransformException(
                BaseMessages.getString(PKG, "HTTPPOST.Exception.Authentication", data.realUrl));
          case -1:
            throw new HopTransformException(
                BaseMessages.getString(PKG, "HTTPPOST.Exception.IllegalStatusCode", data.realUrl));
          case HttpURLConnection.HTTP_NO_CONTENT:
            body = "";
            break;
          default:
            HttpEntity entity = httpResponse.getEntity();
            if (entity != null) {
              body = EntityUtils.toString(entity);
            } else {
              body = "";
            }
            Header[] headers = searchForHeaders(httpResponse);
            // Use request encoding if specified in component to avoid strange response encodings

            JSONObject json = new JSONObject();
            for (Header header : headers) {
              Object previousValue = json.get(header.getName());
              if (previousValue == null) {
                json.put(header.getName(), header.getValue());
              } else if (previousValue instanceof List) {
                List<String> list = (List<String>) previousValue;
                list.add(header.getValue());
              } else {
                ArrayList<String> list = new ArrayList<>();
                list.add((String) previousValue);
                list.add(header.getValue());
                json.put(header.getName(), list);
              }
            }
            headerString = json.toJSONString();
        }

        if (isDebug()) {
          logDebug(BaseMessages.getString(PKG, "HTTPPOST.Log.ResponseBody", body));
        }

        int returnFieldsOffset = data.inputRowMeta.size();
        if (!Utils.isEmpty(meta.getHttpPostResultField().get(0).getName())) {
          newRow = RowDataUtil.addValueData(newRow, returnFieldsOffset, body);
          returnFieldsOffset++;
        }

        if (!Utils.isEmpty(meta.getHttpPostResultField().get(0).getCode())) {
          newRow = RowDataUtil.addValueData(newRow, returnFieldsOffset, Long.valueOf(statusCode));
          returnFieldsOffset++;
        }
        if (!Utils.isEmpty(meta.getHttpPostResultField().get(0).getResponseTimeFieldName())) {
          newRow = RowDataUtil.addValueData(newRow, returnFieldsOffset, Long.valueOf(responseTime));
          returnFieldsOffset++;
        }
        if (!Utils.isEmpty(meta.getHttpPostResultField().get(0).getResponseHeaderFieldName())) {
          newRow = RowDataUtil.addValueData(newRow, returnFieldsOffset, headerString);
        }
      } finally {
        // Release current connection to the connection pool once you are done
        post.releaseConnection();
        if (httpResponse != null) {
          httpResponse.close();
        }
      }
      return newRow;
    } catch (UnknownHostException uhe) {
      throw new HopException(
          BaseMessages.getString(PKG, "HTTPPOST.Error.UnknownHostException", uhe.getMessage()));
    } catch (Exception e) {
      throw new HopException(
          BaseMessages.getString(PKG, "HTTPPOST.Error.CanNotReadURL", data.realUrl), e);

    } finally {
      if (fis != null) {
        BaseTransform.closeQuietly(fis);
      }
    }
  }

  protected int requestStatusCode(HttpResponse httpResponse) {
    return httpResponse.getStatusLine().getStatusCode();
  }

  protected InputStreamReader openStream(String encoding, HttpResponse httpResponse)
      throws Exception {
    if (!Utils.isEmpty(encoding)) {
      return new InputStreamReader(httpResponse.getEntity().getContent(), encoding);
    } else {
      return new InputStreamReader(httpResponse.getEntity().getContent());
    }
  }

  protected Header[] searchForHeaders(HttpResponse response) {
    return response.getAllHeaders();
  }

  @Override
  public boolean processRow() throws HopException {

    Object[] r = getRow(); // Get row from input rowset & set row busy!
    if (r == null) { // no more input to be expected...
      setOutputDone();
      return false;
    }
    if (first) {
      first = false;
      data.inputRowMeta = getInputRowMeta();
      data.outputRowMeta = getInputRowMeta().clone();
      meta.getFields(data.outputRowMeta, getTransformName(), null, null, this, metadataProvider);

      if (meta.isUrlInField()) {
        if (Utils.isEmpty(meta.getUrlField())) {
          logError(BaseMessages.getString(PKG, "HTTPPOST.Log.NoField"));
          throw new HopException(BaseMessages.getString(PKG, "HTTPPOST.Log.NoField"));
        }

        // cache the position of the field
        if (data.indexOfUrlField < 0) {
          String realUrlfieldName = resolve(meta.getUrlField());
          data.indexOfUrlField = data.inputRowMeta.indexOfValue((realUrlfieldName));
          if (data.indexOfUrlField < 0) {
            // The field is unreachable !
            logError(BaseMessages.getString(PKG, PKG_ERROR_FINDING_FIELD, realUrlfieldName));
            throw new HopException(
                BaseMessages.getString(
                    PKG, "HTTPPOST.Exception.ErrorFindingField", realUrlfieldName));
          }
        }
      } else {
        data.realUrl = resolve(meta.getUrl());
      }
      // set body parameters
      int nrargs = meta.getLookupfield().get(0).getArgumentField().size();
      if (nrargs > 0) {
        data.useBodyParameters = false;
        data.useHeaderParameters = false;
        data.contentTypeHeaderOverwrite = false;
        int nrheader = 0;
        int nrbody = 0;
        for (int i = 0; i < nrargs; i++) { // split into body / header
          if (meta.getLookupfield().get(0).getArgumentField().get(i).isHeader()) {
            data.useHeaderParameters = true; // at least one header parameter
            nrheader++;
          } else {
            data.useBodyParameters = true; // at least one body parameter
            nrbody++;
          }
        }
        data.header_parameters_nrs = new int[nrheader];
        data.headerParameters = new NameValuePair[nrheader];
        data.body_parameters_nrs = new int[nrbody];
        data.bodyParameters = new NameValuePair[nrbody];
        int posHeader = 0;
        int posBody = 0;
        for (int i = 0; i < nrargs; i++) {
          int fieldIndex =
              data.inputRowMeta.indexOfValue(
                  meta.getLookupfield().get(0).getArgumentField().get(i).getName());
          if (fieldIndex < 0) {
            logError(
                BaseMessages.getString(PKG, PKG_ERROR_FINDING_FIELD)
                    + meta.getLookupfield().get(0).getArgumentField().get(i).getName()
                    + "]");
            throw new HopTransformException(
                BaseMessages.getString(
                    PKG,
                    "HTTPPOST.Exception.CouldnotFindField",
                    meta.getLookupfield().get(0).getArgumentField().get(i).getName()));
          }
          if (meta.getLookupfield().get(0).getArgumentField().get(i).isHeader()) {
            data.header_parameters_nrs[posHeader] = fieldIndex;
            data.headerParameters[posHeader] =
                new BasicNameValuePair(
                    resolve(meta.getLookupfield().get(0).getArgumentField().get(i).getParameter()),
                    data.outputRowMeta.getString(r, data.header_parameters_nrs[posHeader]));
            posHeader++;
            if (CONTENT_TYPE.equalsIgnoreCase(
                meta.getLookupfield().get(0).getArgumentField().get(i).getParameter())) {
              data.contentTypeHeaderOverwrite = true; // Content-type will be overwritten
            }
          } else {
            data.body_parameters_nrs[posBody] = fieldIndex;
            data.bodyParameters[posBody] =
                new BasicNameValuePair(
                    resolve(meta.getLookupfield().get(0).getArgumentField().get(i).getParameter()),
                    data.outputRowMeta.getString(r, data.body_parameters_nrs[posBody]));
            posBody++;
          }
        }
      }
      // set query parameters
      int nrQuery = meta.getLookupfield().get(0).getQueryField().size();
      if (nrQuery > 0) {
        data.useQueryParameters = true;
        data.query_parameters_nrs = new int[nrQuery];
        data.queryParameters = new NameValuePair[nrQuery];
        for (int i = 0; i < nrQuery; i++) {
          data.query_parameters_nrs[i] =
              data.inputRowMeta.indexOfValue(
                  meta.getLookupfield().get(0).getQueryField().get(i).getName());
          if (data.query_parameters_nrs[i] < 0) {
            logError(
                BaseMessages.getString(PKG, PKG_ERROR_FINDING_FIELD)
                    + meta.getLookupfield().get(0).getQueryField().get(i).getName()
                    + "]");
            throw new HopTransformException(
                BaseMessages.getString(
                    PKG,
                    "HTTPPOST.Exception.CouldnotFindField",
                    meta.getLookupfield().get(0).getQueryField().get(i).getName()));
          }
          data.queryParameters[i] =
              new BasicNameValuePair(
                  resolve(meta.getLookupfield().get(0).getQueryField().get(i).getParameter()),
                  data.outputRowMeta.getString(r, data.query_parameters_nrs[i]));
        }
      }
      // set request entity?
      if (!Utils.isEmpty(meta.getRequestEntity())) {
        data.indexOfRequestEntity =
            data.inputRowMeta.indexOfValue(resolve(meta.getRequestEntity()));
        if (data.indexOfRequestEntity < 0) {
          throw new HopTransformException(
              BaseMessages.getString(
                  PKG,
                  "HTTPPOST.Exception.CouldnotFindRequestEntityField",
                  meta.getRequestEntity()));
        }
      }
      data.realEncoding = resolve(meta.getEncoding());
    } // end if first

    try {
      Object[] outputRowData = callHttpPOST(r);
      putRow(data.outputRowMeta, outputRowData); // copy row to output rowset(s)

      if (checkFeedback(getLinesRead()) && isDetailed()) {
        logDetailed(BaseMessages.getString(PKG, "HTTPPOST.LineNumber") + getLinesRead());
      }
    } catch (HopException e) {
      boolean sendToErrorRow = false;
      String errorMessage = null;

      if (getTransformMeta().isDoingErrorHandling()) {
        sendToErrorRow = true;
        errorMessage = e.toString();
      } else {
        logError(BaseMessages.getString(PKG, "HTTPPOST.ErrorInTransformRunning") + e.getMessage());
        setErrors(1);
        logError(Const.getStackTracker(e));
        stopAll();
        setOutputDone(); // signal end to receiver(s)
        return false;
      }

      if (sendToErrorRow) {
        // Simply add this row to the error row
        putError(getInputRowMeta(), r, 1, errorMessage, null, "HTTPPOST001");
      }
    }

    return true;
  }

  @VisibleForTesting
  String getRequestBodyParamsAsStr(NameValuePair[] pairs, String charset) throws HopException {
    StringBuilder buf = new StringBuilder();
    try {
      for (int i = 0; i < pairs.length; ++i) {
        NameValuePair pair = pairs[i];
        if (pair.getName() != null) {
            if (i > 0) {
                buf.append("&");
            }
            System.out.println("Caracter set: " + charset);
            
            if (!charset.matches("US-ASCII")) {  
              buf.append(
                  URLEncoder.encode(
                      pair.getName(), !StringUtil.isEmpty(charset) ? charset : DEFAULT_ENCODING));

              if (!StringUtil.isEmpty(pair.getName()))
                buf.append("=");

              if (pair.getValue() != null) {
                buf.append(
                    URLEncoder.encode(
                        pair.getValue(), !StringUtil.isEmpty(charset) ? charset : DEFAULT_ENCODING));
              }
            } else {
                buf.append(pair.getName());
              
                if (!StringUtil.isEmpty(pair.getName()))
                    buf.append("=");

                buf.append(pair.getValue());                           
            }
        }
      }
      return buf.toString();
    } catch (UnsupportedEncodingException e) {
      throw new HopException(e.getMessage(), e.getCause());
    }
  }

  @Override
  public boolean init() {

    if (super.init()) {
      // get authentication settings once
      data.realProxyHost = resolve(meta.getProxyHost());
      data.realProxyPort = Const.toInt(resolve(meta.getProxyPort()), 8080);
      data.realHttpLogin = resolve(meta.getHttpLogin());
      data.realHttpPassword = Utils.resolvePassword(variables, meta.getHttpPassword());

      data.realSocketTimeout = Const.toInt(resolve(meta.getSocketTimeout()), -1);
      data.realConnectionTimeout = Const.toInt(resolve(meta.getSocketTimeout()), -1);
      data.realcloseIdleConnectionsTime =
          Const.toInt(resolve(meta.getCloseIdleConnectionsTime()), -1);

      return true;
    }
    return false;
  }
}
