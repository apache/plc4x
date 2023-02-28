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

import java.util.ArrayList;
import java.util.logging.Level;
import java.util.logging.Logger;
import org.apache.hop.core.Const;
import org.apache.hop.core.exception.HopException;
import org.apache.hop.core.util.Utils;
import org.apache.hop.core.variables.IVariables;
import org.apache.hop.i18n.BaseMessages;
import org.apache.hop.ui.core.dialog.BaseDialog;
import org.apache.hop.ui.core.gui.WindowProperty;
import org.apache.hop.ui.core.widget.ColumnInfo;
import org.apache.hop.ui.core.widget.TableView;
import org.apache.hop.ui.pipeline.transform.BaseTransformDialog;
import org.apache.hop.ui.workflow.action.ActionDialog;
import org.apache.hop.ui.workflow.dialog.WorkflowDialog;
import org.apache.hop.workflow.WorkflowMeta;
import org.apache.hop.workflow.action.IAction;
import org.apache.hop.workflow.action.IActionDialog;
import org.apache.plc4x.hop.metadata.Plc4xConnection;
import org.eclipse.swt.SWT;
import org.eclipse.swt.events.*;
import org.eclipse.swt.layout.FormAttachment;
import org.eclipse.swt.layout.FormData;
import org.eclipse.swt.layout.FormLayout;
import org.eclipse.swt.widgets.*;

import org.apache.hop.metadata.api.IHopMetadata;
import org.apache.hop.ui.core.dialog.ErrorDialog;
import org.apache.hop.ui.core.metadata.MetadataManager;
import org.apache.hop.ui.core.widget.MetaSelectionLine;


public class Plc4xCheckConnectionsDialog extends ActionDialog implements IActionDialog {
  private static final Class<?> PKG = Plc4xCheckConnectionsDialog.class; // Needed by Translator

  private Shell shell;

  private Text wName;

  private Plc4xCheckConnections action;

  private boolean changed;

  private TableView wFields;
  
  private MetaSelectionLine<Plc4xConnection> wConnection;  

  public Plc4xCheckConnectionsDialog(Shell parent, IAction action, WorkflowMeta workflowMeta, IVariables variables) {
    super( parent, workflowMeta, variables );
    this.action = (Plc4xCheckConnections) action;
    if ( this.action.getName() == null ) {       
      this.action.setName( BaseMessages.getString( PKG, "Plc4xCheckConnections.Label" ) );
    }
  }

  @Override
  public IAction open() {
    Shell parent = getParent();

    shell = new Shell(parent, SWT.DIALOG_TRIM | SWT.MIN | SWT.MAX | SWT.RESIZE);
    props.setLook(shell);
    WorkflowDialog.setShellImage(shell, action);

    ModifyListener lsMod = (ModifyEvent e) -> action.setChanged();
    changed = action.hasChanged();

    FormLayout formLayout = new FormLayout();
    formLayout.marginWidth = Const.FORM_MARGIN;
    formLayout.marginHeight = Const.FORM_MARGIN;

    shell.setLayout(formLayout);
    shell.setText(BaseMessages.getString(PKG, "Plc4xCheckonnections.Title"));

    int middle = props.getMiddlePct();
    int margin = Const.MARGIN;

    // Buttons at the bottom
    //
    Button wOk = new Button(shell, SWT.PUSH);
    wOk.setText(BaseMessages.getString(PKG, "System.Button.OK"));
    wOk.addListener(SWT.Selection, (Event e) -> ok());
    Button wCancel = new Button(shell, SWT.PUSH);
    wCancel.setText(BaseMessages.getString(PKG, "System.Button.Cancel"));
    wCancel.addListener(SWT.Selection, (Event e) -> cancel());
    BaseTransformDialog.positionBottomButtons(shell, new Button[] {wOk, wCancel}, margin, null);

    // Filename line
    Label wlName = new Label(shell, SWT.RIGHT);
    wlName.setText(BaseMessages.getString(PKG, "Plc4xCheckConnections.Name.Label"));
    props.setLook(wlName);
    FormData fdlName = new FormData();
    fdlName.left = new FormAttachment(0, 0);
    fdlName.right = new FormAttachment(middle, -margin);
    fdlName.top = new FormAttachment(0, margin);
    wlName.setLayoutData(fdlName);
    wName = new Text(shell, SWT.SINGLE | SWT.LEFT | SWT.BORDER);
    props.setLook(wName);
    wName.addModifyListener(lsMod);
    FormData fdName = new FormData();
    fdName.left = new FormAttachment(middle, 0);
    fdName.top = new FormAttachment(0, margin);
    fdName.right = new FormAttachment(100, 0);
    wName.setLayoutData(fdName);

    Label wlFields = new Label(shell, SWT.NONE);
    wlFields.setText(BaseMessages.getString(PKG, "Plc4xCheckConnections.Fields.Label"));
    props.setLook(wlFields);
    FormData fdlFields = new FormData();
    fdlFields.left = new FormAttachment(0, 0);
    // fdlFields.right= new FormAttachment(middle, -margin);
    fdlFields.top = new FormAttachment(wName, 2 * margin);
    wlFields.setLayoutData(fdlFields);
    

    wConnection =
        new MetaSelectionLine<>(
            variables,
            metadataProvider,
            Plc4xConnection.class,
            shell,
            SWT.NONE,
            BaseMessages.getString(PKG, "Plc4xCheckConnections.Connection.Label"),
            BaseMessages.getString(PKG, "Plc4xCheckConnections.Connection.Tooltip"));
    props.setLook(wConnection);
    FormData fdConnection = new FormData();
    fdConnection.left = new FormAttachment(0, 0);
    fdConnection.right = new FormAttachment(100, 0);
    //fdConnection.top = new FormAttachment(wTransformName, margin);
    wConnection.setLayoutData(fdConnection);

    try {
      wConnection.fillItems();
    } catch (Exception e) {
      new ErrorDialog(shell, "Error", "Error listing Cassandra connection metadata objects", e);
    }    
    

    // Buttons to the right of the screen...
    Button wbGetConnections = new Button(shell, SWT.PUSH | SWT.CENTER);
    props.setLook(wbGetConnections);
    wbGetConnections.setText(
        BaseMessages.getString(PKG, "Plc4xCheckConnections.GetConnections"));
    wbGetConnections.setToolTipText(
        BaseMessages.getString(PKG, "Plc4xCheckConnections.GetConnections.Tooltip"));
    FormData fdbGetConnections = new FormData();
    fdbGetConnections.right = new FormAttachment(100, -margin);
    fdbGetConnections.top = new FormAttachment(wlFields, margin);
    wbGetConnections.setLayoutData(fdbGetConnections);

    // Buttons to the right of the screen...
    Button wbdSourceFileFolder = new Button(shell, SWT.PUSH | SWT.CENTER);
    props.setLook(wbdSourceFileFolder);
    wbdSourceFileFolder.setText(
        BaseMessages.getString(PKG, "Plc4xCheckConnections.DeleteEntry"));
    wbdSourceFileFolder.setToolTipText(
        BaseMessages.getString(PKG, "Plc4xCheckConnections.DeleteSourceFileButton.Label"));
    FormData fdbdSourceFileFolder = new FormData();
    fdbdSourceFileFolder.right = new FormAttachment(100, -margin);
    fdbdSourceFileFolder.top = new FormAttachment(wbGetConnections, margin);
    wbdSourceFileFolder.setLayoutData(fdbdSourceFileFolder);

    int rows =
        action.getConnections() == null
            ? 1
            : (action.getConnections().length == 0 ? 0 : action.getConnections().length);

    final int FieldsRows = rows;

    ColumnInfo[] colinf =
        new ColumnInfo[] {
          new ColumnInfo(
              BaseMessages.getString(PKG, "Plc4xCheckConnections.Fields.Argument.Label"),
              ColumnInfo.COLUMN_TYPE_CCOMBO,
              this.getWorkflowMeta().getDatabaseNames(),
              false),
          new ColumnInfo(
              BaseMessages.getString(PKG, "Plc4xCheckConnections.Fields.WaitFor.Label"),
              ColumnInfo.COLUMN_TYPE_TEXT,
              false),
          new ColumnInfo(
              BaseMessages.getString(PKG, "Plc4xCheckConnections.Fields.WaitForTime.Label"),
              ColumnInfo.COLUMN_TYPE_CCOMBO,
              Plc4xCheckConnections.unitTimeDesc,
              false),
        };

    colinf[0].setToolTip(BaseMessages.getString(PKG, "Plc4xCheckConnections.Fields.Column"));
    colinf[1].setUsingVariables(true);
    colinf[1].setToolTip(BaseMessages.getString(PKG, "Plc4xCheckConnections.WaitFor.ToolTip"));

    wFields =
        new TableView(
            variables,
            shell,
            SWT.BORDER | SWT.FULL_SELECTION | SWT.MULTI,
            colinf,
            FieldsRows,
            lsMod,
            props);

    FormData fdFields = new FormData();
    fdFields.left = new FormAttachment(0, 0);
    fdFields.top = new FormAttachment(wlFields, margin);
    fdFields.right = new FormAttachment(wbGetConnections, -margin);
    fdFields.bottom = new FormAttachment(wOk, -2 * margin);
    wFields.setLayoutData(fdFields);

    // Delete files from the list of files...
    wbdSourceFileFolder.addListener(
        SWT.Selection,
        e -> {
          int[] idx = wFields.getSelectionIndices();
          wFields.remove(idx);
          wFields.removeEmptyRows();
          wFields.setRowNums();
        });

    // get connections...
    wbGetConnections.addListener(SWT.Selection, e -> getConnections());

    getData();

    BaseDialog.defaultShellHandling(shell, c -> ok(), c -> cancel());

    return action;
  }

  //  public void addDatabases() {
  //    connections = workflowMeta.getDatabaseNames();
  //  }

  public void getConnections() {
      /*
    this.workflowMeta.getMetadataProvider().getMetadataClasses().forEach(
           (c) -> {
               System.out.println("Name: " + c.getName());
           });
    IHopMetadataProvider hmdp = this.workflowMeta.getMetadataProvider();
      */
    Class<IHopMetadata> metadataClass = null;

    java.util.List<Plc4xConnection> databases = new ArrayList<>(); //this.workflowMeta.getMetadataProvider().
      try {
          metadataClass = metadataProvider.getMetadataClassForKey("plc4x-connection");
            MetadataManager<IHopMetadata> manager = new MetadataManager<>(variables, metadataProvider, metadataClass, null);          
            manager.getNames().forEach((s)->{
              try {
                  databases.add((Plc4xConnection) manager.loadElement(s));
              } catch (HopException ex) {
                  Logger.getLogger(Plc4xCheckConnectionsDialog.class.getName()).log(Level.SEVERE, null, ex);
              }
            });
          
      } catch (Exception ex) {
          Logger.getLogger(Plc4xCheckConnectionsDialog.class.getName()).log(Level.SEVERE, null, ex);
      }
      
     
    wFields.removeAll();

    for (Plc4xConnection ci : databases) {
      if (ci != null) {
        wFields.add(new String[] {ci.getName(), "0", Plc4xCheckConnections.unitTimeDesc[0]});
      }
    }
    wFields.removeEmptyRows();
    wFields.setRowNums();
    wFields.optWidth(true);
  }

  public void dispose() {
    WindowProperty winprop = new WindowProperty(shell);
    props.setScreen(winprop);
    shell.dispose();
  }

  /** Copy information from the meta-data input to the dialog fields. */
  public void getData() {
    if (action.getName() != null) {
      wName.setText(action.getName());
    }

    if (action.getConnections() != null) {
      for (int i = 0; i < action.getConnections().length; i++) {
        TableItem ti = wFields.table.getItem(i);
        if (action.getConnections()[i] != null) {
          ti.setText(1, action.getConnections()[i].getName());
          ti.setText(2, "" + Const.toInt(action.getWaitfors()[i], 0));
          ti.setText(3, Plc4xCheckConnections.getWaitTimeDesc(action.getWaittimes()[i]));
        }
      }
      wFields.setRowNums();
      wFields.optWidth(true);
    }
    wName.selectAll();
    wName.setFocus();
  }

  private void cancel() {
    action.setChanged(changed);
    action = null;
    dispose();
  }

  private void ok() {
    if (Utils.isEmpty(wName.getText())) {
      MessageBox mb = new MessageBox(shell, SWT.OK | SWT.ICON_ERROR);
      mb.setText(BaseMessages.getString(PKG, "System.TransformActionNameMissing.Title"));
      mb.setMessage(BaseMessages.getString(PKG, "System.ActionNameMissing.Msg"));
      mb.open();
      return;
    }
    Class<IHopMetadata> metadataClass = null;
    action.setName(wName.getText());

    int nrItems = wFields.nrNonEmpty();
    System.out.println("Numero de items: " + nrItems);
    
    Plc4xConnection[] connections = new Plc4xConnection[nrItems];
    String[] waitfors = new String[nrItems];
    int[] waittimes = new int[nrItems];
    
    
      try {
          metadataClass = metadataProvider.getMetadataClassForKey("plc4x-connection");
      } catch (HopException ex) {
          Logger.getLogger(Plc4xCheckConnectionsDialog.class.getName()).log(Level.SEVERE, null, ex);
      }
    MetadataManager<IHopMetadata> manager = new MetadataManager<>(variables, metadataProvider, metadataClass, null);        

    for (int i = 0; i < nrItems; i++) {
        String arg = wFields.getNonEmpty(i).getText(1);
        Plc4xConnection conn;
        try {
            conn = (Plc4xConnection) manager.loadElement(arg);
            if (conn != null) {
                connections[i] = conn;
                waitfors[i] = "" + Const.toInt(wFields.getNonEmpty(i).getText(2), 0);
                waittimes[i] =
                    Plc4xCheckConnections.getWaitTimeByDesc(wFields.getNonEmpty(i).getText(3));
          }            
        } catch (HopException ex) {
            Logger.getLogger(Plc4xCheckConnectionsDialog.class.getName()).log(Level.SEVERE, null, ex);
        }

    }
    action.setConnections(connections);
    action.setWaitfors(waitfors);
    action.setWaittimes(waittimes);

    dispose();
  }
}
