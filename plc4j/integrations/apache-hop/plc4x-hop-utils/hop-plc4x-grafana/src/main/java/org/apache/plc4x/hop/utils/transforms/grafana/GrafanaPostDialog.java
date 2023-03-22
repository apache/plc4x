/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package org.apache.plc4x.hop.utils.transforms.grafana;

import org.apache.hop.core.Const;
import org.apache.hop.core.Props;
import org.apache.hop.core.exception.HopException;
import org.apache.hop.core.row.IRowMeta;
import org.apache.hop.core.util.Utils;
import org.apache.hop.core.variables.IVariables;
import org.apache.hop.i18n.BaseMessages;
import org.apache.hop.pipeline.PipelineMeta;
import org.apache.hop.pipeline.transform.BaseTransformMeta;
import org.apache.hop.pipeline.transform.ITransformDialog;
import org.apache.hop.pipeline.transform.TransformMeta;
import org.apache.hop.ui.core.dialog.BaseDialog;
import org.apache.hop.ui.core.dialog.ErrorDialog;
import org.apache.hop.ui.core.widget.*;
import org.apache.hop.ui.pipeline.transform.BaseTransformDialog;
import org.apache.hop.ui.pipeline.transform.ComponentSelectionListener;
import org.apache.hop.ui.pipeline.transform.ITableItemInsertListener;
import org.eclipse.swt.SWT;
import org.eclipse.swt.custom.CTabFolder;
import org.eclipse.swt.custom.CTabItem;
import org.eclipse.swt.events.*;
import org.eclipse.swt.graphics.Cursor;
import org.eclipse.swt.graphics.Point;
import org.eclipse.swt.layout.FormAttachment;
import org.eclipse.swt.layout.FormData;
import org.eclipse.swt.layout.FormLayout;
import org.eclipse.swt.widgets.*;

import java.nio.charset.Charset;
import java.util.List;
import java.util.*;

public class GrafanaPostDialog extends BaseTransformDialog implements ITransformDialog {
  private static final Class<?> PKG = GrafanaPostMeta.class; // For Translator

  private static final String[] YES_NO_COMBO =
      new String[] {
        BaseMessages.getString(PKG, "System.Combo.No"),
        BaseMessages.getString(PKG, "System.Combo.Yes")
      };
  private static final String YES = BaseMessages.getString(PKG, "System.Combo.Yes");
  private static final String NO = BaseMessages.getString(PKG, "System.Combo.No");

  private Label wlUrl;
  private TextVar wUrl;

  private TextVar wResult;

  private TextVar wResultCode;

  private TextVar wResponseTime;
  private TextVar wResponseHeader;

  private TableView wFields;

  private TableView wQuery;

  private Button wUrlInField;

  private Label wlUrlField;
  private ComboVar wUrlField;

  private ComboVar wRequestEntity;

  private TextVar wHttpLogin;

  private TextVar wHttpPassword;

  private TextVar wProxyHost;

  private TextVar wProxyPort;

  private final GrafanaPostMeta input;

  private final Map<String, Integer> inputFields;

  private ColumnInfo[] colinf;
  private ColumnInfo[] colinfquery;

  private String[] fieldNames;

  private boolean gotPreviousFields = false;

  private ComboVar wEncoding;

  private Button wPostAFile;

  private boolean gotEncodings = false;

  private TextVar wConnectionTimeOut;

  private TextVar wSocketTimeOut;

  private TextVar wCloseIdleConnectionsTime;

  public GrafanaPostDialog(
      Shell parent, IVariables variables, Object in, PipelineMeta pipelineMeta, String sname) {
    super(parent, variables, (BaseTransformMeta) in, pipelineMeta, sname);
    input = (GrafanaPostMeta) in;
    inputFields = new HashMap<>();
  }

  @Override
  public String open() {
    Shell parent = getParent();

    shell = new Shell(parent, SWT.DIALOG_TRIM | SWT.RESIZE | SWT.MAX | SWT.MIN);
    props.setLook(shell);
    setShellImage(shell, input);

    ModifyListener lsMod = e -> input.setChanged();

    changed = input.hasChanged();

    FormLayout formLayout = new FormLayout();
    formLayout.marginWidth = Const.FORM_MARGIN;
    formLayout.marginHeight = Const.FORM_MARGIN;

    shell.setLayout(formLayout);
    shell.setText(BaseMessages.getString(PKG, "HTTPPOSTDialog.Shell.Title"));

    int middle = props.getMiddlePct();
    int margin = props.getMargin();

    // THE BUTTONS
    wOk = new Button(shell, SWT.PUSH);
    wOk.setText(BaseMessages.getString(PKG, "System.Button.OK"));
    wOk.addListener(SWT.Selection, e -> ok());
    wCancel = new Button(shell, SWT.PUSH);
    wCancel.setText(BaseMessages.getString(PKG, "System.Button.Cancel"));
    wCancel.addListener(SWT.Selection, e -> cancel());
    setButtonPositions(new Button[] {wOk, wCancel}, margin, null);

    // TransformName line
    wlTransformName = new Label(shell, SWT.RIGHT);
    wlTransformName.setText(BaseMessages.getString(PKG, "HTTPPOSTDialog.TransformName.Label"));
    props.setLook(wlTransformName);
    fdlTransformName = new FormData();
    fdlTransformName.left = new FormAttachment(0, 0);
    fdlTransformName.right = new FormAttachment(middle, -margin);
    fdlTransformName.top = new FormAttachment(0, margin);
    wlTransformName.setLayoutData(fdlTransformName);
    wTransformName = new Text(shell, SWT.SINGLE | SWT.LEFT | SWT.BORDER);
    wTransformName.setText(transformName);
    props.setLook(wTransformName);
    wTransformName.addModifyListener(lsMod);
    fdTransformName = new FormData();
    fdTransformName.left = new FormAttachment(middle, 0);
    fdTransformName.top = new FormAttachment(0, margin);
    fdTransformName.right = new FormAttachment(100, 0);
    wTransformName.setLayoutData(fdTransformName);

    CTabFolder wTabFolder = new CTabFolder(shell, SWT.BORDER);
    props.setLook(wTabFolder, Props.WIDGET_STYLE_TAB);

    // ////////////////////////
    // START OF GENERAL TAB ///
    // ////////////////////////
    CTabItem wGeneralTab = new CTabItem(wTabFolder, SWT.NONE);
    wGeneralTab.setText(BaseMessages.getString(PKG, "HTTPPOSTDialog.GeneralTab.Title"));

    Composite wGeneralComp = new Composite(wTabFolder, SWT.NONE);
    props.setLook(wGeneralComp);

    FormLayout fileLayout = new FormLayout();
    fileLayout.marginWidth = 3;
    fileLayout.marginHeight = 3;
    wGeneralComp.setLayout(fileLayout);

    // ////////////////////////
    // START Settings GROUP

    Group gSettings = new Group(wGeneralComp, SWT.SHADOW_ETCHED_IN);
    gSettings.setText(BaseMessages.getString(PKG, "HTTPPOSTDialog.SettingsGroup.Label"));
    FormLayout settingsLayout = new FormLayout();
    settingsLayout.marginWidth = 3;
    settingsLayout.marginHeight = 3;
    gSettings.setLayout(settingsLayout);
    props.setLook(gSettings);

    wlUrl = new Label(gSettings, SWT.RIGHT);
    wlUrl.setText(BaseMessages.getString(PKG, "HTTPPOSTDialog.URL.Label"));
    props.setLook(wlUrl);
    FormData fdlUrl = new FormData();
    fdlUrl.left = new FormAttachment(0, 0);
    fdlUrl.right = new FormAttachment(middle, -margin);
    fdlUrl.top = new FormAttachment(wTransformName, margin);
    wlUrl.setLayoutData(fdlUrl);

    wUrl = new TextVar(variables, gSettings, SWT.SINGLE | SWT.LEFT | SWT.BORDER);
    props.setLook(wUrl);
    wUrl.addModifyListener(lsMod);
    FormData fdUrl = new FormData();
    fdUrl.left = new FormAttachment(middle, 0);
    fdUrl.top = new FormAttachment(wTransformName, margin);
    fdUrl.right = new FormAttachment(100, 0);
    wUrl.setLayoutData(fdUrl);

    // UrlInField line
    Label wlUrlInField = new Label(gSettings, SWT.RIGHT);
    wlUrlInField.setText(BaseMessages.getString(PKG, "HTTPPOSTDialog.UrlInField.Label"));
    props.setLook(wlUrlInField);
    FormData fdlUrlInField = new FormData();
    fdlUrlInField.left = new FormAttachment(0, 0);
    fdlUrlInField.top = new FormAttachment(wUrl, margin);
    fdlUrlInField.right = new FormAttachment(middle, -margin);
    wlUrlInField.setLayoutData(fdlUrlInField);
    wUrlInField = new Button(gSettings, SWT.CHECK);
    props.setLook(wUrlInField);
    FormData fdUrlInField = new FormData();
    fdUrlInField.left = new FormAttachment(middle, 0);
    fdUrlInField.top = new FormAttachment(wlUrlInField, 0, SWT.CENTER);
    fdUrlInField.right = new FormAttachment(100, 0);
    wUrlInField.setLayoutData(fdUrlInField);
    wUrlInField.addSelectionListener(
        new SelectionAdapter() {
          @Override
          public void widgetSelected(SelectionEvent e) {
            input.setChanged();
            activeUrlInfield();
          }
        });

    // UrlField Line
    wlUrlField = new Label(gSettings, SWT.RIGHT);
    wlUrlField.setText(BaseMessages.getString(PKG, "HTTPPOSTDialog.UrlField.Label"));
    props.setLook(wlUrlField);
    FormData fdlUrlField = new FormData();
    fdlUrlField.left = new FormAttachment(0, 0);
    fdlUrlField.right = new FormAttachment(middle, -margin);
    fdlUrlField.top = new FormAttachment(wUrlInField, margin);
    wlUrlField.setLayoutData(fdlUrlField);

    wUrlField = new ComboVar(variables, gSettings, SWT.BORDER | SWT.READ_ONLY);
    wUrlField.setEditable(true);
    props.setLook(wUrlField);
    wUrlField.addModifyListener(lsMod);
    FormData fdUrlField = new FormData();
    fdUrlField.left = new FormAttachment(middle, 0);
    fdUrlField.top = new FormAttachment(wUrlInField, margin);
    fdUrlField.right = new FormAttachment(100, -margin);
    wUrlField.setLayoutData(fdUrlField);
    wUrlField.addFocusListener(
        new FocusListener() {
          @Override
          public void focusLost(FocusEvent e) {
            // Disable focuslost event
          }

          @Override
          public void focusGained(FocusEvent e) {
            Cursor busy = new Cursor(shell.getDisplay(), SWT.CURSOR_WAIT);
            shell.setCursor(busy);
            setStreamFields();
            shell.setCursor(null);
            busy.dispose();
          }
        });

    Label wlEncoding = new Label(gSettings, SWT.RIGHT);
    wlEncoding.setText(BaseMessages.getString(PKG, "HTTPPOSTDialog.Encoding.Label"));
    props.setLook(wlEncoding);
    FormData fdlEncoding = new FormData();
    fdlEncoding.left = new FormAttachment(0, 0);
    fdlEncoding.top = new FormAttachment(wUrlField, margin);
    fdlEncoding.right = new FormAttachment(middle, -margin);
    wlEncoding.setLayoutData(fdlEncoding);
    wEncoding = new ComboVar(variables, gSettings, SWT.BORDER | SWT.READ_ONLY);
    wEncoding.setEditable(true);
    props.setLook(wEncoding);
    wEncoding.addModifyListener(lsMod);
    FormData fdEncoding = new FormData();
    fdEncoding.left = new FormAttachment(middle, 0);
    fdEncoding.top = new FormAttachment(wUrlField, margin);
    fdEncoding.right = new FormAttachment(100, -margin);
    wEncoding.setLayoutData(fdEncoding);
    wEncoding.addFocusListener(
        new FocusListener() {
          @Override
          public void focusLost(FocusEvent e) {
            // Disable focuslost event
          }

          @Override
          public void focusGained(FocusEvent e) {
            Cursor busy = new Cursor(shell.getDisplay(), SWT.CURSOR_WAIT);
            shell.setCursor(busy);
            setEncodings();
            shell.setCursor(null);
            busy.dispose();
          }
        });

    // requestEntity Line
    Label wlRequestEntity = new Label(gSettings, SWT.RIGHT);
    wlRequestEntity.setText(BaseMessages.getString(PKG, "HTTPPOSTDialog.requestEntity.Label"));
    props.setLook(wlRequestEntity);
    FormData fdlRequestEntity = new FormData();
    fdlRequestEntity.left = new FormAttachment(0, 0);
    fdlRequestEntity.right = new FormAttachment(middle, -margin);
    fdlRequestEntity.top = new FormAttachment(wEncoding, margin);
    wlRequestEntity.setLayoutData(fdlRequestEntity);

    wRequestEntity = new ComboVar(variables, gSettings, SWT.BORDER | SWT.READ_ONLY);
    wRequestEntity.setEditable(true);
    props.setLook(wRequestEntity);
    wRequestEntity.addModifyListener(lsMod);
    FormData fdRequestEntity = new FormData();
    fdRequestEntity.left = new FormAttachment(middle, 0);
    fdRequestEntity.top = new FormAttachment(wEncoding, margin);
    fdRequestEntity.right = new FormAttachment(100, -margin);
    wRequestEntity.setLayoutData(fdRequestEntity);
    wRequestEntity.addFocusListener(
        new FocusListener() {
          @Override
          public void focusLost(FocusEvent e) {
            // Disable focuslost event
          }

          @Override
          public void focusGained(FocusEvent e) {
            Cursor busy = new Cursor(shell.getDisplay(), SWT.CURSOR_WAIT);
            shell.setCursor(busy);
            setStreamFields();
            shell.setCursor(null);
            busy.dispose();
          }
        });

    // Post file?
    Label wlPostAFile = new Label(gSettings, SWT.RIGHT);
    wlPostAFile.setText(BaseMessages.getString(PKG, "HTTPPOSTDialog.postAFile.Label"));
    props.setLook(wlPostAFile);
    FormData fdlPostAFile = new FormData();
    fdlPostAFile.left = new FormAttachment(0, 0);
    fdlPostAFile.right = new FormAttachment(middle, -margin);
    fdlPostAFile.top = new FormAttachment(wRequestEntity, margin);
    wlPostAFile.setLayoutData(fdlPostAFile);
    wPostAFile = new Button(gSettings, SWT.CHECK);
    wPostAFile.setToolTipText(BaseMessages.getString(PKG, "HTTPPOSTDialog.postAFile.Tooltip"));
    props.setLook(wPostAFile);
    FormData fdPostAFile = new FormData();
    fdPostAFile.left = new FormAttachment(middle, 0);
    fdPostAFile.top = new FormAttachment(wlPostAFile, 0, SWT.CENTER);
    fdPostAFile.right = new FormAttachment(100, 0);
    wPostAFile.setLayoutData(fdPostAFile);
    wPostAFile.addSelectionListener(new ComponentSelectionListener(input));

    Label wlConnectionTimeOut = new Label(gSettings, SWT.RIGHT);
    wlConnectionTimeOut.setText(
        BaseMessages.getString(PKG, "HTTPPOSTDialog.ConnectionTimeOut.Label"));
    props.setLook(wlConnectionTimeOut);
    FormData fdlConnectionTimeOut = new FormData();
    fdlConnectionTimeOut.top = new FormAttachment(wPostAFile, margin);
    fdlConnectionTimeOut.left = new FormAttachment(0, 0);
    fdlConnectionTimeOut.right = new FormAttachment(middle, -margin);
    wlConnectionTimeOut.setLayoutData(fdlConnectionTimeOut);
    wConnectionTimeOut = new TextVar(variables, gSettings, SWT.SINGLE | SWT.LEFT | SWT.BORDER);
    wConnectionTimeOut.addModifyListener(lsMod);
    wConnectionTimeOut.setToolTipText(
        BaseMessages.getString(PKG, "HTTPPOSTDialog.ConnectionTimeOut.Tooltip"));
    props.setLook(wConnectionTimeOut);
    FormData fdConnectionTimeOut = new FormData();
    fdConnectionTimeOut.top = new FormAttachment(wPostAFile, margin);
    fdConnectionTimeOut.left = new FormAttachment(middle, 0);
    fdConnectionTimeOut.right = new FormAttachment(100, 0);
    wConnectionTimeOut.setLayoutData(fdConnectionTimeOut);

    Label wlSocketTimeOut = new Label(gSettings, SWT.RIGHT);
    wlSocketTimeOut.setText(BaseMessages.getString(PKG, "HTTPPOSTDialog.SocketTimeOut.Label"));
    props.setLook(wlSocketTimeOut);
    FormData fdlSocketTimeOut = new FormData();
    fdlSocketTimeOut.top = new FormAttachment(wConnectionTimeOut, margin);
    fdlSocketTimeOut.left = new FormAttachment(0, 0);
    fdlSocketTimeOut.right = new FormAttachment(middle, -margin);
    wlSocketTimeOut.setLayoutData(fdlSocketTimeOut);
    wSocketTimeOut = new TextVar(variables, gSettings, SWT.SINGLE | SWT.LEFT | SWT.BORDER);
    wSocketTimeOut.addModifyListener(lsMod);
    wSocketTimeOut.setToolTipText(
        BaseMessages.getString(PKG, "HTTPPOSTDialog.SocketTimeOut.Tooltip"));
    props.setLook(wSocketTimeOut);
    FormData fdSocketTimeOut = new FormData();
    fdSocketTimeOut.top = new FormAttachment(wConnectionTimeOut, margin);
    fdSocketTimeOut.left = new FormAttachment(middle, 0);
    fdSocketTimeOut.right = new FormAttachment(100, 0);
    wSocketTimeOut.setLayoutData(fdSocketTimeOut);

    Label wlCloseIdleConnectionsTime = new Label(gSettings, SWT.RIGHT);
    wlCloseIdleConnectionsTime.setText(
        BaseMessages.getString(PKG, "HTTPPOSTDialog.CloseIdleConnectionsTime.Label"));
    props.setLook(wlCloseIdleConnectionsTime);
    FormData fdlCloseIdleConnectionsTime = new FormData();
    fdlCloseIdleConnectionsTime.top = new FormAttachment(wSocketTimeOut, margin);
    fdlCloseIdleConnectionsTime.left = new FormAttachment(0, 0);
    fdlCloseIdleConnectionsTime.right = new FormAttachment(middle, -margin);
    wlCloseIdleConnectionsTime.setLayoutData(fdlCloseIdleConnectionsTime);
    wCloseIdleConnectionsTime =
        new TextVar(variables, gSettings, SWT.SINGLE | SWT.LEFT | SWT.BORDER);
    wCloseIdleConnectionsTime.addModifyListener(lsMod);
    wCloseIdleConnectionsTime.setToolTipText(
        BaseMessages.getString(PKG, "HTTPPOSTDialog.CloseIdleConnectionsTime.Tooltip"));
    props.setLook(wCloseIdleConnectionsTime);
    FormData fdCloseIdleConnectionsTime = new FormData();
    fdCloseIdleConnectionsTime.top = new FormAttachment(wSocketTimeOut, margin);
    fdCloseIdleConnectionsTime.left = new FormAttachment(middle, 0);
    fdCloseIdleConnectionsTime.right = new FormAttachment(100, 0);
    wCloseIdleConnectionsTime.setLayoutData(fdCloseIdleConnectionsTime);

    FormData fdSettings = new FormData();
    fdSettings.left = new FormAttachment(0, 0);
    fdSettings.right = new FormAttachment(100, 0);
    fdSettings.top = new FormAttachment(wTransformName, margin);
    gSettings.setLayoutData(fdSettings);

    // END Output Settings GROUP
    // ////////////////////////

    // ////////////////////////
    // START Output Fields GROUP

    Group gOutputFields = new Group(wGeneralComp, SWT.SHADOW_ETCHED_IN);
    gOutputFields.setText(BaseMessages.getString(PKG, "HTTPDialog.OutputFieldsGroup.Label"));
    FormLayout outputFieldsLayout = new FormLayout();
    outputFieldsLayout.marginWidth = 3;
    outputFieldsLayout.marginHeight = 3;
    gOutputFields.setLayout(outputFieldsLayout);
    props.setLook(gOutputFields);

    // Result line...
    Label wlResult = new Label(gOutputFields, SWT.RIGHT);
    wlResult.setText(BaseMessages.getString(PKG, "HTTPPOSTDialog.Result.Label"));
    props.setLook(wlResult);
    FormData fdlResult = new FormData();
    fdlResult.left = new FormAttachment(0, 0);
    fdlResult.right = new FormAttachment(middle, -margin);
    fdlResult.top = new FormAttachment(wPostAFile, margin);
    wlResult.setLayoutData(fdlResult);
    wResult = new TextVar(variables, gOutputFields, SWT.SINGLE | SWT.LEFT | SWT.BORDER);
    props.setLook(wResult);
    wResult.addModifyListener(lsMod);
    FormData fdResult = new FormData();
    fdResult.left = new FormAttachment(middle, 0);
    fdResult.top = new FormAttachment(wPostAFile, margin);
    fdResult.right = new FormAttachment(100, -margin);
    wResult.setLayoutData(fdResult);

    // Resultcode line...
    Label wlResultCode = new Label(gOutputFields, SWT.RIGHT);
    wlResultCode.setText(BaseMessages.getString(PKG, "HTTPPOSTDialog.ResultCode.Label"));
    props.setLook(wlResultCode);
    FormData fdlResultCode = new FormData();
    fdlResultCode.left = new FormAttachment(0, 0);
    fdlResultCode.right = new FormAttachment(middle, -margin);
    fdlResultCode.top = new FormAttachment(wResult, margin);
    wlResultCode.setLayoutData(fdlResultCode);
    wResultCode = new TextVar(variables, gOutputFields, SWT.SINGLE | SWT.LEFT | SWT.BORDER);
    props.setLook(wResultCode);
    wResultCode.addModifyListener(lsMod);
    FormData fdResultCode = new FormData();
    fdResultCode.left = new FormAttachment(middle, 0);
    fdResultCode.top = new FormAttachment(wResult, margin);
    fdResultCode.right = new FormAttachment(100, -margin);
    wResultCode.setLayoutData(fdResultCode);

    // Response time line...
    Label wlResponseTime = new Label(gOutputFields, SWT.RIGHT);
    wlResponseTime.setText(BaseMessages.getString(PKG, "HTTPPOSTDialog.ResponseTime.Label"));
    props.setLook(wlResponseTime);
    FormData fdlResponseTime = new FormData();
    fdlResponseTime.left = new FormAttachment(0, 0);
    fdlResponseTime.right = new FormAttachment(middle, -margin);
    fdlResponseTime.top = new FormAttachment(wResultCode, margin);
    wlResponseTime.setLayoutData(fdlResponseTime);
    wResponseTime = new TextVar(variables, gOutputFields, SWT.SINGLE | SWT.LEFT | SWT.BORDER);
    props.setLook(wResponseTime);
    wResponseTime.addModifyListener(lsMod);
    FormData fdResponseTime = new FormData();
    fdResponseTime.left = new FormAttachment(middle, 0);
    fdResponseTime.top = new FormAttachment(wResultCode, margin);
    fdResponseTime.right = new FormAttachment(100, 0);
    wResponseTime.setLayoutData(fdResponseTime);
    // Response header line...
    Label wlResponseHeader = new Label(gOutputFields, SWT.RIGHT);
    wlResponseHeader.setText(BaseMessages.getString(PKG, "HTTPPOSTDialog.ResponseHeader.Label"));
    props.setLook(wlResponseHeader);
    FormData fdlResponseHeader = new FormData();
    fdlResponseHeader.left = new FormAttachment(0, 0);
    fdlResponseHeader.right = new FormAttachment(middle, -margin);
    fdlResponseHeader.top = new FormAttachment(wResponseTime, margin);
    wlResponseHeader.setLayoutData(fdlResponseHeader);
    wResponseHeader = new TextVar(variables, gOutputFields, SWT.SINGLE | SWT.LEFT | SWT.BORDER);
    props.setLook(wResponseHeader);
    wResponseHeader.addModifyListener(lsMod);
    FormData fdResponseHeader = new FormData();
    fdResponseHeader.left = new FormAttachment(middle, 0);
    fdResponseHeader.top = new FormAttachment(wResponseTime, margin);
    fdResponseHeader.right = new FormAttachment(100, 0);
    wResponseHeader.setLayoutData(fdResponseHeader);

    FormData fdOutputFields = new FormData();
    fdOutputFields.left = new FormAttachment(0, 0);
    fdOutputFields.right = new FormAttachment(100, 0);
    fdOutputFields.top = new FormAttachment(gSettings, margin);
    gOutputFields.setLayoutData(fdOutputFields);

    // END Output Fields GROUP
    // ////////////////////////

    // ////////////////////////
    // START HTTP AUTH GROUP

    Group gHttpAuth = new Group(wGeneralComp, SWT.SHADOW_ETCHED_IN);
    gHttpAuth.setText(BaseMessages.getString(PKG, "HTTPPOSTDialog.HttpAuthGroup.Label"));
    FormLayout httpAuthLayout = new FormLayout();
    httpAuthLayout.marginWidth = 3;
    httpAuthLayout.marginHeight = 3;
    gHttpAuth.setLayout(httpAuthLayout);
    props.setLook(gHttpAuth);

    // HTTP Login
    Label wlHttpLogin = new Label(gHttpAuth, SWT.RIGHT);
    wlHttpLogin.setText(BaseMessages.getString(PKG, "HTTPPOSTDialog.HttpLogin.Label"));
    props.setLook(wlHttpLogin);
    FormData fdlHttpLogin = new FormData();
    fdlHttpLogin.top = new FormAttachment(0, margin);
    fdlHttpLogin.left = new FormAttachment(0, 0);
    fdlHttpLogin.right = new FormAttachment(middle, -margin);
    wlHttpLogin.setLayoutData(fdlHttpLogin);
    wHttpLogin = new TextVar(variables, gHttpAuth, SWT.SINGLE | SWT.LEFT | SWT.BORDER);
    wHttpLogin.addModifyListener(lsMod);
    wHttpLogin.setToolTipText(BaseMessages.getString(PKG, "HTTPPOSTDialog.HttpLogin.Tooltip"));
    props.setLook(wHttpLogin);
    FormData fdHttpLogin = new FormData();
    fdHttpLogin.top = new FormAttachment(0, margin);
    fdHttpLogin.left = new FormAttachment(middle, 0);
    fdHttpLogin.right = new FormAttachment(100, 0);
    wHttpLogin.setLayoutData(fdHttpLogin);

    // HTTP Password
    Label wlHttpPassword = new Label(gHttpAuth, SWT.RIGHT);
    wlHttpPassword.setText(BaseMessages.getString(PKG, "HTTPPOSTDialog.HttpPassword.Label"));
    props.setLook(wlHttpPassword);
    FormData fdlHttpPassword = new FormData();
    fdlHttpPassword.top = new FormAttachment(wHttpLogin, margin);
    fdlHttpPassword.left = new FormAttachment(0, 0);
    fdlHttpPassword.right = new FormAttachment(middle, -margin);
    wlHttpPassword.setLayoutData(fdlHttpPassword);
    wHttpPassword = new PasswordTextVar(variables, gHttpAuth, SWT.SINGLE | SWT.LEFT | SWT.BORDER);
    wHttpPassword.addModifyListener(lsMod);
    wHttpPassword.setToolTipText(
        BaseMessages.getString(PKG, "HTTPPOSTDialog.HttpPassword.Tooltip"));
    props.setLook(wHttpPassword);
    FormData fdHttpPassword = new FormData();
    fdHttpPassword.top = new FormAttachment(wHttpLogin, margin);
    fdHttpPassword.left = new FormAttachment(middle, 0);
    fdHttpPassword.right = new FormAttachment(100, 0);
    wHttpPassword.setLayoutData(fdHttpPassword);

    FormData fdHttpAuth = new FormData();
    fdHttpAuth.left = new FormAttachment(0, 0);
    fdHttpAuth.right = new FormAttachment(100, 0);
    fdHttpAuth.top = new FormAttachment(gOutputFields, margin);
    gHttpAuth.setLayoutData(fdHttpAuth);

    // END HTTP AUTH GROUP
    // ////////////////////////

    // ////////////////////////
    // START PROXY GROUP

    Group gProxy = new Group(wGeneralComp, SWT.SHADOW_ETCHED_IN);
    gProxy.setText(BaseMessages.getString(PKG, "HTTPPOSTDialog.ProxyGroup.Label"));
    FormLayout proxyLayout = new FormLayout();
    proxyLayout.marginWidth = 3;
    proxyLayout.marginHeight = 3;
    gProxy.setLayout(proxyLayout);
    props.setLook(gProxy);

    // HTTP Login
    Label wlProxyHost = new Label(gProxy, SWT.RIGHT);
    wlProxyHost.setText(BaseMessages.getString(PKG, "HTTPPOSTDialog.ProxyHost.Label"));
    props.setLook(wlProxyHost);
    FormData fdlProxyHost = new FormData();
    fdlProxyHost.top = new FormAttachment(0, margin);
    fdlProxyHost.left = new FormAttachment(0, 0);
    fdlProxyHost.right = new FormAttachment(middle, -margin);
    wlProxyHost.setLayoutData(fdlProxyHost);
    wProxyHost = new TextVar(variables, gProxy, SWT.SINGLE | SWT.LEFT | SWT.BORDER);
    wProxyHost.addModifyListener(lsMod);
    wProxyHost.setToolTipText(BaseMessages.getString(PKG, "HTTPPOSTDialog.ProxyHost.Tooltip"));
    props.setLook(wProxyHost);
    FormData fdProxyHost = new FormData();
    fdProxyHost.top = new FormAttachment(0, margin);
    fdProxyHost.left = new FormAttachment(middle, 0);
    fdProxyHost.right = new FormAttachment(100, 0);
    wProxyHost.setLayoutData(fdProxyHost);

    // HTTP Password
    Label wlProxyPort = new Label(gProxy, SWT.RIGHT);
    wlProxyPort.setText(BaseMessages.getString(PKG, "HTTPPOSTDialog.ProxyPort.Label"));
    props.setLook(wlProxyPort);
    FormData fdlProxyPort = new FormData();
    fdlProxyPort.top = new FormAttachment(wProxyHost, margin);
    fdlProxyPort.left = new FormAttachment(0, 0);
    fdlProxyPort.right = new FormAttachment(middle, -margin);
    wlProxyPort.setLayoutData(fdlProxyPort);
    wProxyPort = new TextVar(variables, gProxy, SWT.SINGLE | SWT.LEFT | SWT.BORDER);
    wProxyPort.addModifyListener(lsMod);
    wProxyPort.setToolTipText(BaseMessages.getString(PKG, "HTTPPOSTDialog.ProxyPort.Tooltip"));
    props.setLook(wProxyPort);
    FormData fdProxyPort = new FormData();
    fdProxyPort.top = new FormAttachment(wProxyHost, margin);
    fdProxyPort.left = new FormAttachment(middle, 0);
    fdProxyPort.right = new FormAttachment(100, 0);
    wProxyPort.setLayoutData(fdProxyPort);

    FormData fdProxy = new FormData();
    fdProxy.left = new FormAttachment(0, 0);
    fdProxy.right = new FormAttachment(100, 0);
    fdProxy.top = new FormAttachment(gHttpAuth, margin);
    gProxy.setLayoutData(fdProxy);

    // END HTTP AUTH GROUP
    // ////////////////////////

    FormData fdGeneralComp = new FormData();
    fdGeneralComp.left = new FormAttachment(0, 0);
    fdGeneralComp.top = new FormAttachment(wTransformName, margin);
    fdGeneralComp.right = new FormAttachment(100, 0);
    fdGeneralComp.bottom = new FormAttachment(100, 0);
    wGeneralComp.setLayoutData(fdGeneralComp);

    wGeneralComp.layout();
    wGeneralTab.setControl(wGeneralComp);

    // ///////////////////////////////////////////////////////////
    // / END OF GENERAL TAB
    // ///////////////////////////////////////////////////////////

    // Additional tab...
    //
    CTabItem wAdditionalTab = new CTabItem(wTabFolder, SWT.NONE);
    wAdditionalTab.setText(BaseMessages.getString(PKG, "HTTPPOSTDialog.FieldsTab.Title"));

    FormLayout addLayout = new FormLayout();
    addLayout.marginWidth = Const.FORM_MARGIN;
    addLayout.marginHeight = Const.FORM_MARGIN;

    Composite wAdditionalComp = new Composite(wTabFolder, SWT.NONE);
    wAdditionalComp.setLayout(addLayout);
    props.setLook(wAdditionalComp);

    Label wlFields = new Label(wAdditionalComp, SWT.NONE);
    wlFields.setText(BaseMessages.getString(PKG, "HTTPPOSTDialog.Parameters.Label"));
    props.setLook(wlFields);
    FormData fdlFields = new FormData();
    fdlFields.left = new FormAttachment(0, 0);
    fdlFields.top = new FormAttachment(gProxy, margin);
    wlFields.setLayoutData(fdlFields);

    int fieldsRows = 0;
    if (input.getLookupfield().get(0).getArgumentField() != null) {
      fieldsRows = input.getLookupfield().get(0).getArgumentField().size();
    }

    colinf =
        new ColumnInfo[] {
          new ColumnInfo(
              BaseMessages.getString(PKG, "HTTPPOSTDialog.ColumnInfo.Name"),
              ColumnInfo.COLUMN_TYPE_CCOMBO,
              new String[] {""},
              false),
          new ColumnInfo(
              BaseMessages.getString(PKG, "HTTPPOSTDialog.ColumnInfo.Parameter"),
              ColumnInfo.COLUMN_TYPE_TEXT,
              false),
          new ColumnInfo(
              BaseMessages.getString(PKG, "HTTPPOSTDialog.ColumnInfo.Header"),
              ColumnInfo.COLUMN_TYPE_CCOMBO,
              YES_NO_COMBO),
        };
    colinf[1].setUsingVariables(true);
    wFields =
        new TableView(
            variables,
            wAdditionalComp,
            SWT.BORDER | SWT.FULL_SELECTION | SWT.MULTI,
            colinf,
            fieldsRows,
            lsMod,
            props);

    Button wGetBodyParam = new Button(wAdditionalComp, SWT.PUSH);
    wGetBodyParam.setText(BaseMessages.getString(PKG, "HTTPPOSTDialog.GetFields.Button"));
    FormData fdGetBodyParam = new FormData();
    fdGetBodyParam.top = new FormAttachment(wlFields, margin);
    fdGetBodyParam.right = new FormAttachment(100, 0);
    wGetBodyParam.setLayoutData(fdGetBodyParam);

    FormData fdFields = new FormData();
    fdFields.left = new FormAttachment(0, 0);
    fdFields.top = new FormAttachment(wlFields, margin);
    fdFields.right = new FormAttachment(wGetBodyParam, -margin);
    fdFields.bottom = new FormAttachment(wlFields, 200);
    wFields.setLayoutData(fdFields);

    Label wlQuery = new Label(wAdditionalComp, SWT.NONE);
    wlQuery.setText(BaseMessages.getString(PKG, "HTTPPOSTDialog.QueryParameters.Label"));
    props.setLook(wlQuery);
    FormData fdlQuery = new FormData();
    fdlQuery.left = new FormAttachment(0, 0);
    fdlQuery.top = new FormAttachment(wFields, margin);
    wlQuery.setLayoutData(fdlQuery);

    int queryRows = 0;
    if (input.getLookupfield().get(0).getQueryField() != null) {
      queryRows = input.getLookupfield().get(0).getQueryField().size();
    }

    colinfquery =
        new ColumnInfo[] {
          new ColumnInfo(
              BaseMessages.getString(PKG, "HTTPPOSTDialog.ColumnInfo.QueryName"),
              ColumnInfo.COLUMN_TYPE_CCOMBO,
              new String[] {""},
              false),
          new ColumnInfo(
              BaseMessages.getString(PKG, "HTTPPOSTDialog.ColumnInfo.QueryParameter"),
              ColumnInfo.COLUMN_TYPE_TEXT,
              false),
        };
    colinfquery[1].setUsingVariables(true);
    wQuery =
        new TableView(
            variables,
            wAdditionalComp,
            SWT.BORDER | SWT.FULL_SELECTION | SWT.MULTI,
            colinfquery,
            queryRows,
            lsMod,
            props);

    wGet = new Button(wAdditionalComp, SWT.PUSH);
    wGet.setText(BaseMessages.getString(PKG, "HTTPPOSTDialog.GetFields.Button"));
    fdGet = new FormData();
    fdGet.top = new FormAttachment(wlQuery, margin);
    fdGet.right = new FormAttachment(100, 0);
    wGet.setLayoutData(fdGet);

    FormData fdQuery = new FormData();
    fdQuery.left = new FormAttachment(0, 0);
    fdQuery.top = new FormAttachment(wlQuery, margin);
    fdQuery.right = new FormAttachment(wGet, -margin);
    fdQuery.bottom = new FormAttachment(100, -margin);
    wQuery.setLayoutData(fdQuery);

    //
    // Search the fields in the background
    //

    final Runnable runnable =
        () -> {
          TransformMeta transformMeta = pipelineMeta.findTransform(transformName);
          if (transformMeta != null) {
            try {
              IRowMeta row = pipelineMeta.getPrevTransformFields(variables, transformMeta);

              // Remember these fields...
              for (int i = 0; i < row.size(); i++) {
                inputFields.put(row.getValueMeta(i).getName(), i);
              }

              setComboBoxes();
            } catch (HopException e) {
              logError(BaseMessages.getString(PKG, "System.Dialog.GetFieldsFailed.Message"));
            }
          }
        };
    new Thread(runnable).start();
    FormData fdAdditionalComp = new FormData();
    fdAdditionalComp.left = new FormAttachment(0, 0);
    fdAdditionalComp.top = new FormAttachment(wTransformName, margin);
    fdAdditionalComp.right = new FormAttachment(100, 0);
    fdAdditionalComp.bottom = new FormAttachment(100, 0);
    wAdditionalComp.setLayoutData(fdAdditionalComp);

    wAdditionalComp.layout();
    wAdditionalTab.setControl(wAdditionalComp);
    // ////// END of Additional Tab

    FormData fdTabFolder = new FormData();
    fdTabFolder.left = new FormAttachment(0, 0);
    fdTabFolder.top = new FormAttachment(wTransformName, margin);
    fdTabFolder.right = new FormAttachment(100, 0);
    fdTabFolder.bottom = new FormAttachment(wOk, -2 * margin);
    wTabFolder.setLayoutData(fdTabFolder);

    // Add listeners
    wGet.addListener(SWT.Selection, e -> getQueryFields());
    wGetBodyParam.addListener(SWT.Selection, e -> get());

    lsResize =
        event -> {
          Point size = shell.getSize();
          wFields.setSize(size.x - 10, size.y - 50);
          wFields.table.setSize(size.x - 10, size.y - 50);
          wFields.redraw();
        };
    shell.addListener(SWT.Resize, lsResize);

    wTabFolder.setSelection(0);
    getData();
    activeUrlInfield();
    input.setChanged(changed);

    BaseDialog.defaultShellHandling(shell, c -> ok(), c -> cancel());

    return transformName;
  }

  protected void setComboBoxes() {
    // Something was changed in the row.
    //
    final Map<String, Integer> fields = new HashMap<>();

    // Add the currentMeta fields...
    fields.putAll(inputFields);

    Set<String> keySet = fields.keySet();
    List<String> entries = new ArrayList<>(keySet);

    fieldNames = entries.toArray(new String[entries.size()]);

    Const.sortStrings(fieldNames);
    colinf[0].setComboValues(fieldNames);
    colinfquery[0].setComboValues(fieldNames);
  }

  private void setStreamFields() {
    if (!gotPreviousFields) {
      String urlfield = wUrlField.getText();
      wUrlField.removeAll();
      wUrlField.setItems(fieldNames);
      if (urlfield != null) {
        wUrlField.setText(urlfield);
      }

      String request = wRequestEntity.getText();
      wRequestEntity.removeAll();
      wRequestEntity.setItems(fieldNames);
      if (request != null) {
        wRequestEntity.setText(request);
      }

      gotPreviousFields = true;
    }
  }

  private void setEncodings() {
    // Encoding of the text file:
    if (!gotEncodings) {
      gotEncodings = true;

      wEncoding.removeAll();
      List<Charset> values = new ArrayList<>(Charset.availableCharsets().values());
      for (Charset charSet : values) {
        wEncoding.add(charSet.displayName());
      }

      // Now select the default!
      String defEncoding = Const.getEnvironmentVariable("file.encoding", "UTF-8");
      int idx = Const.indexOfString(defEncoding, wEncoding.getItems());
      if (idx >= 0) {
        wEncoding.select(idx);
      }
    }
  }

  private void activeUrlInfield() {
    wlUrlField.setEnabled(wUrlInField.getSelection());
    wUrlField.setEnabled(wUrlInField.getSelection());
    wlUrl.setEnabled(!wUrlInField.getSelection());
    wUrl.setEnabled(!wUrlInField.getSelection());
  }

  /** Copy information from the meta-data input to the dialog fields. */
  public void getData() {
    if (log.isDebug()) {
      logDebug(BaseMessages.getString(PKG, "HTTPPOSTDialog.Log.GettingKeyInfo"));
    }

    if (input.getLookupfield().get(0).getArgumentField() != null) {
      for (int i = 0; i < input.getLookupfield().get(0).getArgumentField().size(); i++) {
        TableItem item = wFields.table.getItem(i);
        if (input.getLookupfield().get(0).getArgumentField().get(i).getName() != null) {
          item.setText(1, input.getLookupfield().get(0).getArgumentField().get(i).getName());
        }
        if (input.getLookupfield().get(0).getArgumentField().get(i).getParameter() != null) {
          item.setText(2, input.getLookupfield().get(0).getArgumentField().get(i).getParameter());
        }
        item.setText(
            3, (input.getLookupfield().get(0).getArgumentField().get(i).isHeader()) ? YES : NO);
      }
    }
    if (input.getLookupfield().get(0).getQueryField() != null) {
      for (int i = 0; i < input.getLookupfield().get(0).getQueryField().size(); i++) {
        TableItem item = wQuery.table.getItem(i);
        if (input.getLookupfield().get(0).getQueryField().get(i).getName() != null) {
          item.setText(1, input.getLookupfield().get(0).getQueryField().get(i).getName());
        }
        if (input.getLookupfield().get(0).getQueryField().get(i).getParameter() != null) {
          item.setText(2, input.getLookupfield().get(0).getQueryField().get(i).getParameter());
        }
      }
    }
    if (input.getUrl() != null) {
      wUrl.setText(input.getUrl());
    }
    wUrlInField.setSelection(input.isUrlInField());
    if (input.getUrlField() != null) {
      wUrlField.setText(input.getUrlField());
    }
    if (input.getRequestEntity() != null) {
      wRequestEntity.setText(input.getRequestEntity());
    }
    if (input.getHttpPostResultField().get(0).getName() != null) {
      wResult.setText(input.getHttpPostResultField().get(0).getName());
    }
    if (input.getHttpPostResultField().get(0).getCode() != null) {
      wResultCode.setText(input.getHttpPostResultField().get(0).getCode());
    }
    if (input.getHttpPostResultField().get(0).getResponseTimeFieldName() != null) {
      wResponseTime.setText(input.getHttpPostResultField().get(0).getResponseTimeFieldName());
    }
    if (input.getEncoding() != null) {
      wEncoding.setText(input.getEncoding());
    }
    wPostAFile.setSelection(input.isPostAFile());

    if (input.getHttpLogin() != null) {
      wHttpLogin.setText(input.getHttpLogin());
    }
    if (input.getHttpPassword() != null) {
      wHttpPassword.setText(input.getHttpPassword());
    }
    if (input.getProxyHost() != null) {
      wProxyHost.setText(input.getProxyHost());
    }
    if (input.getProxyPort() != null) {
      wProxyPort.setText(input.getProxyPort());
    }
    if (input.getHttpPostResultField().get(0).getResponseHeaderFieldName() != null) {
      wResponseHeader.setText(input.getHttpPostResultField().get(0).getResponseHeaderFieldName());
    }

    wSocketTimeOut.setText(Const.NVL(input.getSocketTimeout(), ""));
    wConnectionTimeOut.setText(Const.NVL(input.getConnectionTimeout(), ""));
    wCloseIdleConnectionsTime.setText(Const.NVL(input.getCloseIdleConnectionsTime(), ""));

    wFields.setRowNums();
    wFields.optWidth(true);

    wTransformName.selectAll();
    wTransformName.setFocus();
  }

  private void cancel() {
    transformName = null;
    input.setChanged(changed);
    dispose();
  }

  private void ok() {
    GrafanaPostLoookupField loookupField = new GrafanaPostLoookupField();
    if (Utils.isEmpty(wTransformName.getText())) {
      return;
    }

    int nrargs = wFields.nrNonEmpty();

    if (log.isDebug()) {
      logDebug(
          BaseMessages.getString(PKG, "HTTPPOSTDialog.Log.FoundArguments", String.valueOf(nrargs)));
    }
    // CHECKSTYLE:Indentation:OFF
    for (int i = 0; i < nrargs; i++) {
      TableItem item = wFields.getNonEmpty(i);
      GrafanaPostArgumentField argumentField =
          new GrafanaPostArgumentField(item.getText(1), item.getText(2), YES.equals(item.getText(3)));
      loookupField.getArgumentField().add(argumentField);
    }

    int nrqueryparams = wQuery.nrNonEmpty();

    if (log.isDebug()) {
      logDebug(
          BaseMessages.getString(
              PKG, "HTTPPOSTDialog.Log.FoundQueryParameters", String.valueOf(nrqueryparams)));
    }
    // CHECKSTYLE:Indentation:OFF
    for (int i = 0; i < nrqueryparams; i++) {
      TableItem item = wQuery.getNonEmpty(i);
      input.getLookupfield().get(0).getQueryField().clear();
      GrafanaPostQuery httpPostQuery = new GrafanaPostQuery(item.getText(1), item.getText(2));
      loookupField.getQueryField().add(httpPostQuery);
    }

    List<GrafanaPostLoookupField> listLookupField = new ArrayList<>();
    listLookupField.add(loookupField);
    input.setLookupfield(listLookupField);

    input.setUrl(wUrl.getText());
    input.setUrlField(wUrlField.getText());
    input.setRequestEntity(wRequestEntity.getText());
    input.setUrlInField(wUrlInField.getSelection());

    GrafanaPostResultField httpPostResultField =
        new GrafanaPostResultField(
            wResultCode.getText(),
            wResult.getText(),
            wResponseTime.getText(),
            wResponseHeader.getText());

    List<GrafanaPostResultField> listHttpPostResultField = new ArrayList<>();
    listHttpPostResultField.add(httpPostResultField);
    input.setHttpPostResultField(listHttpPostResultField);

    input.setEncoding(wEncoding.getText());
    input.setPostAFile(wPostAFile.getSelection());
    input.setHttpLogin(wHttpLogin.getText());
    input.setHttpPassword(wHttpPassword.getText());
    input.setProxyHost(wProxyHost.getText());
    input.setProxyPort(wProxyPort.getText());
    input.setSocketTimeout(wSocketTimeOut.getText());
    input.setConnectionTimeout(wConnectionTimeOut.getText());
    input.setCloseIdleConnectionsTime(wCloseIdleConnectionsTime.getText());

    transformName = wTransformName.getText(); // return value

    dispose();
  }

  private void get() {
    try {
      IRowMeta r = pipelineMeta.getPrevTransformFields(variables, transformName);
      if (r != null && !r.isEmpty()) {
        ITableItemInsertListener listener =
            (tableItem, v) -> {
              tableItem.setText(3, NO); // default is "N"
              return true;
            };
        BaseTransformDialog.getFieldsFromPrevious(
            r, wFields, 1, new int[] {1, 2}, null, -1, -1, listener);
      }
    } catch (HopException ke) {
      new ErrorDialog(
          shell,
          BaseMessages.getString(PKG, "HTTPPOSTDialog.FailedToGetFields.DialogTitle"),
          BaseMessages.getString(PKG, "HTTPPOSTDialog.FailedToGetFields.DialogMessage"),
          ke);
    }
  }

  private void getQueryFields() {
    try {
      IRowMeta r = pipelineMeta.getPrevTransformFields(variables, transformName);
      if (r != null && !r.isEmpty()) {
        BaseTransformDialog.getFieldsFromPrevious(
            r, wQuery, 1, new int[] {1, 2}, new int[] {3}, -1, -1, null);
      }
    } catch (HopException ke) {
      new ErrorDialog(
          shell,
          BaseMessages.getString(PKG, "HTTPPOSTDialog.FailedToGetFields.DialogTitle"),
          BaseMessages.getString(PKG, "HTTPPOSTDialog.FailedToGetFields.DialogMessage"),
          ke);
    }
  }
}
