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

package org.apache.plc4x.hop.transforms.plc4xsubs;

import org.apache.plc4x.hop.transforms.plc4xinput.*;
import org.apache.hop.core.Const;
import org.apache.hop.core.exception.HopException;
import org.apache.hop.core.row.value.ValueMetaFactory;
import org.apache.hop.core.util.Utils;
import org.apache.hop.core.variables.IVariables;
import org.apache.hop.i18n.BaseMessages;
import org.apache.hop.pipeline.Pipeline;
import org.apache.hop.pipeline.PipelineMeta;
import org.apache.hop.pipeline.PipelinePreviewFactory;
import org.apache.hop.pipeline.transform.BaseTransformMeta;
import org.apache.hop.pipeline.transform.ITransformDialog;
import org.apache.hop.ui.core.ConstUi;
import org.apache.hop.ui.core.dialog.BaseDialog;
import org.apache.hop.ui.core.dialog.EnterNumberDialog;
import org.apache.hop.ui.core.dialog.EnterTextDialog;
import org.apache.hop.ui.core.dialog.ErrorDialog;
import org.apache.hop.ui.core.dialog.PreviewRowsDialog;
import org.apache.hop.ui.core.widget.ColumnInfo;
import org.apache.hop.ui.core.widget.MetaSelectionLine;
import org.apache.hop.ui.core.widget.TableView;
import org.apache.hop.ui.core.widget.TextVar;
import org.apache.hop.ui.pipeline.dialog.PipelinePreviewProgressDialog;
import org.apache.hop.ui.pipeline.transform.BaseTransformDialog;
import org.apache.hop.ui.util.SwtSvgImageUtil;
import org.apache.plc4x.hop.metadata.Plc4xConnection;
import org.apache.plc4x.hop.transforms.util.Plc4xGeneratorField;
import org.eclipse.swt.SWT;
import org.eclipse.swt.events.*;
import org.eclipse.swt.graphics.Image;
import org.eclipse.swt.graphics.Point;
import org.eclipse.swt.layout.FormAttachment;
import org.eclipse.swt.layout.FormData;
import org.eclipse.swt.layout.FormLayout;
import org.eclipse.swt.widgets.Button;
import org.eclipse.swt.widgets.Control;
import org.eclipse.swt.widgets.Label;
import org.eclipse.swt.widgets.Shell;
import org.eclipse.swt.widgets.TableItem;
import org.eclipse.swt.widgets.Text;

public class Plc4xSubsDialog extends BaseTransformDialog implements ITransformDialog {
  private static final Class<?> PKG = Plc4xSubsDialog.class; // Needed by Translator

  private MetaSelectionLine<Plc4xConnection> wConnection;  
  
  private Label wlLimit;
  private TextVar wLimit;

  private Button wNeverEnding;

  private Label wlMaxwait;
  private TextVar wMaxwait;  
  
  private Label wlInterval;
  private TextVar wInterval;

  private Label wlRowTimeField;
  private TextVar wRowTimeField;

  private Label wlLastTimeField;
  private TextVar wLastTimeField;

  private TableView wFields;
  
  private final Plc4xSubsMeta input;

  public Plc4xSubsDialog(Shell parent, IVariables variables , Object in, PipelineMeta pipelineMeta, String sname ) {
    super( parent, variables, (BaseTransformMeta) in, pipelineMeta, sname );
    input = (Plc4xSubsMeta) in;
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
    shell.setText(BaseMessages.getString(PKG, "Plc4x.Read.Meta.Dialog.DialogTitle"));

    int middle = props.getMiddlePct();
    int margin = props.getMargin();

    // Filename line
    wlTransformName = new Label(shell, SWT.RIGHT);
    wlTransformName.setText(BaseMessages.getString(PKG, "System.Label.TransformName"));
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

    // Connection line
    wConnection =
        new MetaSelectionLine<>(
            variables,
            metadataProvider,
            Plc4xConnection.class,
            shell,
            SWT.NONE,
            BaseMessages.getString(PKG, "Plc4x.Read.Meta.Dialog.Connection.Label"),
            BaseMessages.getString(PKG, "Plc4x.Read.Meta.Dialog.Connection.Tooltip"));
    FormData fdConnection = new FormData();
    fdConnection.left = new FormAttachment(0, 0);
    fdConnection.right = new FormAttachment(100, 0);
    fdConnection.top = new FormAttachment(wTransformName, margin);
    wConnection.setLayoutData(fdConnection);    
    props.setLook(wConnection);
    try {
      wConnection.fillItems();
    } catch (Exception e) {
      new ErrorDialog(shell, "Error", "Error listing Cassandra connection metadata objects", e);
    }    
    
    Control lastControl = wConnection;    
    
    wlLimit = new Label(shell, SWT.RIGHT);
    wlLimit.setText(BaseMessages.getString(PKG, "Plc4x.Read.Meta.Dialog.Limit.Label"));
    props.setLook(wlLimit);
    FormData fdlLimit = new FormData();
    fdlLimit.left = new FormAttachment(0, 0);
    fdlLimit.right = new FormAttachment(middle, -margin);
    fdlLimit.top = new FormAttachment(lastControl, margin);
    wlLimit.setLayoutData(fdlLimit);
    wLimit = new TextVar(variables, shell, SWT.SINGLE | SWT.LEFT | SWT.BORDER);
    props.setLook(wLimit);
    wLimit.addModifyListener(lsMod);
    FormData fdLimit = new FormData();
    fdLimit.left = new FormAttachment(middle, 0);
    fdLimit.top = new FormAttachment(lastControl, margin);
    fdLimit.right = new FormAttachment(100, 0);
    wLimit.setLayoutData(fdLimit);
    lastControl = wLimit;

    Label wlNeverEnding = new Label(shell, SWT.RIGHT);
    wlNeverEnding.setText(BaseMessages.getString(PKG, "Plc4x.Read.Meta.Dialog.NeverEnding.Label"));
    props.setLook(wlNeverEnding);
    FormData fdlNeverEnding = new FormData();
    fdlNeverEnding.left = new FormAttachment(0, 0);
    fdlNeverEnding.right = new FormAttachment(middle, -margin);
    fdlNeverEnding.top = new FormAttachment(lastControl, margin);
    wlNeverEnding.setLayoutData(fdlNeverEnding);
    wNeverEnding = new Button(shell, SWT.CHECK);
    props.setLook(wNeverEnding);
    wNeverEnding.addSelectionListener(
        new SelectionAdapter() {
          @Override
          public void widgetSelected(SelectionEvent e) {
            setActive();
            input.setChanged();
          }
        });
    FormData fdNeverEnding = new FormData();
    fdNeverEnding.left = new FormAttachment(middle, 0);
    fdNeverEnding.top = new FormAttachment(wlNeverEnding, 0, SWT.CENTER);
    fdNeverEnding.right = new FormAttachment(100, 0);
    wNeverEnding.setLayoutData(fdNeverEnding);
    lastControl = wlNeverEnding;
    
    wlMaxwait = new Label(shell, SWT.RIGHT);
    wlMaxwait.setText(BaseMessages.getString(PKG, "Plc4x.Read.Meta.Dialog.Maxwait.Label"));
    props.setLook(wlMaxwait);
    FormData fdlMaxwait = new FormData();
    fdlMaxwait.left = new FormAttachment(0, 0);
    fdlMaxwait.right = new FormAttachment(middle, -margin);
    fdlMaxwait.top = new FormAttachment(lastControl, margin);
    wlMaxwait.setLayoutData(fdlMaxwait);
    wMaxwait = new TextVar(variables, shell, SWT.SINGLE | SWT.LEFT | SWT.BORDER);
    props.setLook(wMaxwait);
    wMaxwait.addModifyListener(lsMod);
    FormData fdMaxwait = new FormData();
    fdMaxwait.left = new FormAttachment(middle, 0);
    fdMaxwait.top = new FormAttachment(lastControl, margin);
    fdMaxwait.right = new FormAttachment(100, 0);
    wMaxwait.setLayoutData(fdMaxwait);
    lastControl = wMaxwait;    

    wlInterval = new Label(shell, SWT.RIGHT);
    wlInterval.setText(BaseMessages.getString(PKG, "Plc4x.Read.Meta.Dialog.Interval.Label"));
    props.setLook(wlInterval);
    FormData fdlInterval = new FormData();
    fdlInterval.left = new FormAttachment(0, 0);
    fdlInterval.right = new FormAttachment(middle, -margin);
    fdlInterval.top = new FormAttachment(lastControl, margin);
    wlInterval.setLayoutData(fdlInterval);
    wInterval = new TextVar(variables, shell, SWT.SINGLE | SWT.LEFT | SWT.BORDER);
    props.setLook(wInterval);
    wInterval.addModifyListener(lsMod);
    FormData fdInterval = new FormData();
    fdInterval.left = new FormAttachment(middle, 0);
    fdInterval.top = new FormAttachment(lastControl, margin);
    fdInterval.right = new FormAttachment(100, 0);
    wInterval.setLayoutData(fdInterval);
    lastControl = wInterval;
    
    /*
    wlRowTimeField = new Label(shell, SWT.RIGHT);
    wlRowTimeField.setText(BaseMessages.getString(PKG, "Plc4x.Read.Meta.Dialog.RowTimeField.Label"));
    props.setLook(wlRowTimeField);
    FormData fdlRowTimeField = new FormData();
    fdlRowTimeField.left = new FormAttachment(0, 0);
    fdlRowTimeField.right = new FormAttachment(middle, -margin);
    fdlRowTimeField.top = new FormAttachment(lastControl, margin);
    wlRowTimeField.setLayoutData(fdlRowTimeField);
    wRowTimeField = new TextVar(variables, shell, SWT.SINGLE | SWT.LEFT | SWT.BORDER);
    props.setLook(wRowTimeField);
    wRowTimeField.addModifyListener(lsMod);
    FormData fdRowTimeField = new FormData();
    fdRowTimeField.left = new FormAttachment(middle, 0);
    fdRowTimeField.top = new FormAttachment(lastControl, margin);
    fdRowTimeField.right = new FormAttachment(100, 0);
    wRowTimeField.setLayoutData(fdRowTimeField);
    lastControl = wRowTimeField;

    wlLastTimeField = new Label(shell, SWT.RIGHT);
    wlLastTimeField.setText(BaseMessages.getString(PKG, "Plc4x.Read.Meta.Dialog.LastTimeField.Label"));
    props.setLook(wlLastTimeField);
    FormData fdlLastTimeField = new FormData();
    fdlLastTimeField.left = new FormAttachment(0, 0);
    fdlLastTimeField.right = new FormAttachment(middle, -margin);
    fdlLastTimeField.top = new FormAttachment(lastControl, margin);
    wlLastTimeField.setLayoutData(fdlLastTimeField);
    wLastTimeField = new TextVar(variables, shell, SWT.SINGLE | SWT.LEFT | SWT.BORDER);
    props.setLook(wLastTimeField);
    wLastTimeField.addModifyListener(lsMod);
    FormData fdLastTimeField = new FormData();
    fdLastTimeField.left = new FormAttachment(middle, 0);
    fdLastTimeField.top = new FormAttachment(lastControl, margin);
    fdLastTimeField.right = new FormAttachment(100, 0);
    wLastTimeField.setLayoutData(fdLastTimeField);
    lastControl = wLastTimeField;
    */
    wOk = new Button(shell, SWT.PUSH);
    wOk.setText(BaseMessages.getString(PKG, "System.Button.OK"));
    wOk.addListener(SWT.Selection, e -> ok());
    wPreview = new Button(shell, SWT.PUSH);
    wPreview.setText(BaseMessages.getString(PKG, "System.Button.Preview"));
    wPreview.addListener(SWT.Selection, e -> preview());
    wCancel = new Button(shell, SWT.PUSH);
    wCancel.setText(BaseMessages.getString(PKG, "System.Button.Cancel"));
    wCancel.addListener(SWT.Selection, e -> cancel());

    setButtonPositions(new Button[] {wOk, wPreview, wCancel}, margin, null);

    Label wlFields = new Label(shell, SWT.NONE);
    wlFields.setText(BaseMessages.getString(PKG, "Plc4x.Read.Meta.Dialog.Fields.Label"));
    props.setLook(wlFields);
    FormData fdlFields = new FormData();
    fdlFields.left = new FormAttachment(0, 0);
    fdlFields.top = new FormAttachment(lastControl, margin);
    wlFields.setLayoutData(fdlFields);
    lastControl = wlFields;

    final int nrFields = input.getFields().size();

    ColumnInfo[] colinf =
        new ColumnInfo[] {
          new ColumnInfo(
              BaseMessages.getString(PKG, "System.Column.Name"),
              ColumnInfo.COLUMN_TYPE_TEXT,
              false),
          new ColumnInfo(
              BaseMessages.getString(PKG, "Plc4x.Read.Meta.Dialog.Fields.Item"),
              ColumnInfo.COLUMN_TYPE_TEXT,
              false),            
          new ColumnInfo(
              BaseMessages.getString(PKG, "System.Column.Type"),
              ColumnInfo.COLUMN_TYPE_CCOMBO,
              ValueMetaFactory.getValueMetaNames()),
          new ColumnInfo(
              BaseMessages.getString(PKG, "System.Column.Format"),
              ColumnInfo.COLUMN_TYPE_FORMAT,
              2),
          new ColumnInfo(
              BaseMessages.getString(PKG, "System.Column.Length"),
              ColumnInfo.COLUMN_TYPE_TEXT,
              false),
          new ColumnInfo(
              BaseMessages.getString(PKG, "System.Column.Precision"),
              ColumnInfo.COLUMN_TYPE_TEXT,
              false),
          new ColumnInfo(
              BaseMessages.getString(PKG, "System.Column.Currency"),
              ColumnInfo.COLUMN_TYPE_TEXT,
              false),
          new ColumnInfo(
              BaseMessages.getString(PKG, "System.Column.Decimal"),
              ColumnInfo.COLUMN_TYPE_TEXT,
              false),
          new ColumnInfo(
              BaseMessages.getString(PKG, "System.Column.Group"),
              ColumnInfo.COLUMN_TYPE_TEXT,
              false),
          new ColumnInfo(
              BaseMessages.getString(PKG, "System.Column.Value"),
              ColumnInfo.COLUMN_TYPE_TEXT,
              false),
          new ColumnInfo(
              BaseMessages.getString(PKG, "System.Column.SetEmptyString"),
              ColumnInfo.COLUMN_TYPE_CCOMBO,
              new String[] {
                BaseMessages.getString(PKG, "System.Combo.Yes"),
                BaseMessages.getString(PKG, "System.Combo.No")
              })
        };

    wFields =
        new TableView(
            variables,
            shell,
            SWT.BORDER | SWT.FULL_SELECTION | SWT.MULTI,
            colinf,
            nrFields,
            lsMod,
            props);

    FormData fdFields = new FormData();
    fdFields.left = new FormAttachment(0, 0);
    fdFields.top = new FormAttachment(lastControl, margin);
    fdFields.right = new FormAttachment(100, 0);
    fdFields.bottom = new FormAttachment(wOk, -2 * margin);
    wFields.setLayoutData(fdFields);

    lsResize =
        event -> {
          Point size = shell.getSize();
          wFields.setSize(size.x - 10, size.y - 50);
          wFields.table.setSize(size.x - 10, size.y - 50);
          wFields.redraw();
        };
    shell.addListener(SWT.Resize, lsResize);

    getData();
    input.setChanged(changed);

    BaseDialog.defaultShellHandling(shell, c -> ok(), c -> cancel());

    return transformName;      
    };

  @Override
    public void setActive() {
    boolean neverEnding = wNeverEnding.getSelection();

    wlLimit.setEnabled(!neverEnding);
    wLimit.setEnabled(!neverEnding);

    //wlInterval.setEnabled(neverEnding);
    //wInterval.setEnabled(neverEnding);

    //wlRowTimeField.setEnabled(neverEnding);
    //wRowTimeField.setEnabled(neverEnding);

    //wlLastTimeField.setEnabled(neverEnding);
    //wLastTimeField.setEnabled(neverEnding);
  }



  private Image getImage() {
    return SwtSvgImageUtil.getImage( shell.getDisplay(), getClass().getClassLoader(), "plc4x_toddy_read.svg", ConstUi.LARGE_ICON_SIZE,
      ConstUi.LARGE_ICON_SIZE );
  }

  /**
   * Copy information from the meta-data input to the dialog fields.
   * 
   */
  public void getData() {
    if (isDebug()) {
      logDebug("getting fields info...");
    }
    if (input.getConnection() == null) {
        wConnection.setText("");
    } else {
        wConnection.setText(input.getConnection());        
    }

    wLimit.setText(input.getRowLimit());
    wNeverEnding.setSelection(input.isNeverEnding());
    wMaxwait.setText(Const.NVL(input.getMaxwaitInMs(), ""));
    wInterval.setText(Const.NVL(input.getIntervalInMs(), ""));
   // wRowTimeField.setText(Const.NVL(input.getRowTimeField(), ""));
    //wLastTimeField.setText(Const.NVL(input.getLastTimeField(), ""));

    for (int i = 0; i < input.getFields().size(); i++) {
      Plc4xGeneratorField field = input.getFields().get(i);
      TableItem item = wFields.table.getItem(i);
      int col = 1;
      item.setText(col++, Const.NVL(field.getName(), ""));

      item.setText(col++, Const.NVL(field.getItem(), ""));      
      
      String type = field.getType();
      String format = field.getFormat();
      String length = field.getLength() < 0 ? "" : ("" + field.getLength());
      String prec = field.getPrecision() < 0 ? "" : ("" + field.getPrecision());

      String curr = field.getCurrency();
      String group = field.getGroup();
      String decim = field.getDecimal();
      String def = field.getValue();

      item.setText(col++, Const.NVL(type, ""));
      item.setText(col++, Const.NVL(format, ""));
      item.setText(col++, Const.NVL(length, ""));
      item.setText(col++, Const.NVL(prec, ""));
      item.setText(col++, Const.NVL(curr, ""));
      item.setText(col++, Const.NVL(decim, ""));
      item.setText(col++, Const.NVL(group, ""));
      item.setText(col++, Const.NVL(def, ""));
      item.setText(
          col++,
          field.isSetEmptyString()
              ? BaseMessages.getString(PKG, "System.Combo.Yes")
              : BaseMessages.getString(PKG, "System.Combo.No"));
    }

    wFields.setRowNums();
    wFields.optWidth(true);

    setActive();

    wTransformName.selectAll();
    wTransformName.setFocus();
  }

  /**
   * Stores the information from the dialog box in meta-data.
   *
   * @param meta 
   */
  private void getInfo( Plc4xSubsMeta meta)  throws HopException  {
    meta.setConnection(wConnection.getText());
    meta.setRowLimit(wLimit.getText());
    meta.setNeverEnding(wNeverEnding.getSelection());
    meta.setMaxwaitInMs(wMaxwait.getText());
    meta.setIntervalInMs(wInterval.getText());
    //meta.setRowTimeField(wRowTimeField.getText());
    //meta.setLastTimeField(wLastTimeField.getText());

    meta.getFields().clear();

    // CHECKSTYLE:Indentation:OFF
    for (TableItem item : wFields.getNonEmptyItems()) {
        Plc4xGeneratorField field = new Plc4xGeneratorField();
        field.setName(item.getText(1));
        field.setItem(item.getText(2));
        field.setFormat(item.getText(4));
        field.setLength(Const.toInt(item.getText(5), -1));
        field.setPrecision(Const.toInt(item.getText(6), -1));
        field.setCurrency(item.getText(7));
        field.setDecimal(item.getText(8));
        field.setGroup(item.getText(9));
        field.setValue(field.isSetEmptyString() ? "" : item.getText(10));
        field.setSetEmptyString(
            BaseMessages.getString(PKG, "System.Combo.Yes").equalsIgnoreCase(item.getText(11)));
            field.setType(field.isSetEmptyString() ? "String" : item.getText(3));   
            
        meta.getFields().add(field);            
    }
  }

  /**
   * Cancel the dialog.
   */
  private void cancel() {
    transformName = null;
    input.setChanged( changed );
    dispose();
  }

  private void ok() {
    if (Utils.isEmpty(wTransformName.getText())) {
      return;
    }

    transformName = wTransformName.getText(); // return value
    try {
      getInfo(new Plc4xSubsMeta()); // to see if there is an exception
      getInfo(input); // to put the content on the input structure for real if all is well.
      dispose();
    } catch (HopException e) {
      new ErrorDialog(
          shell,
          BaseMessages.getString(PKG, "Plc4x.Read.Meta.Dialog.Illegal.Dialog.Settings.Title"),
          BaseMessages.getString(PKG, "Plc4x.Read.Meta.Dialog.Illegal.Dialog.Settings.Message"),
          e);
    }
  }
  

  /**
   * Preview the data generated by this transform. This generates a pipeline using this transform &
   * a dummy and previews it.
   */
  private void preview() {
    Plc4xSubsMeta oneMeta = new Plc4xSubsMeta();
    try {
      getInfo(oneMeta);
    } catch (HopException e) {
      new ErrorDialog(
          shell,
          BaseMessages.getString(PKG, "Plc4x.Read.Meta.Dialog.Illegal.Dialog.Settings.Title"),
          BaseMessages.getString(PKG, "Plc4x.Read.Meta.Dialog.Illegal.Dialog.Settings.Message"),
          e);
      return;
    }

    PipelineMeta previewMeta =
        PipelinePreviewFactory.generatePreviewPipeline(
            pipelineMeta.getMetadataProvider(), oneMeta, wTransformName.getText());

    EnterNumberDialog numberDialog =
        new EnterNumberDialog(
            shell,
            props.getDefaultPreviewSize(),
            BaseMessages.getString(PKG, "System.Dialog.EnterPreviewSize.Title"),
            BaseMessages.getString(PKG, "System.Dialog.EnterPreviewSize.Message"));
    int previewSize = numberDialog.open();
    if (previewSize > 0) {
      PipelinePreviewProgressDialog progressDialog =
          new PipelinePreviewProgressDialog(
              shell,
              variables,
              previewMeta,
              new String[] {wTransformName.getText()},
              new int[] {previewSize});
      progressDialog.open();

      Pipeline pipeline = progressDialog.getPipeline();
      String loggingText = progressDialog.getLoggingText();

      if (!progressDialog.isCancelled()) {
        if (pipeline.getResult() != null && pipeline.getResult().getNrErrors() > 0) {
          EnterTextDialog etd =
              new EnterTextDialog(
                  shell,
                  BaseMessages.getString(PKG, "System.Dialog.PreviewError.Title"),
                  BaseMessages.getString(PKG, "System.Dialog.PreviewError.Message"),
                  loggingText,
                  true);
          etd.setReadOnly();
          etd.open();
        }
      }

      PreviewRowsDialog prd =
          new PreviewRowsDialog(
              shell,
              variables,
              SWT.NONE,
              wTransformName.getText(),
              progressDialog.getPreviewRowsMeta(wTransformName.getText()),
              progressDialog.getPreviewRows(wTransformName.getText()),
              loggingText);
      prd.open();
    }
  }  
  
}
