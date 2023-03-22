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

import org.apache.hop.core.Const;
import org.apache.hop.ui.core.PropsUi;
import org.apache.hop.ui.core.dialog.ErrorDialog;
import org.apache.hop.ui.core.gui.GuiCompositeWidgets;
import org.apache.hop.ui.core.gui.GuiCompositeWidgetsAdapter;
import org.apache.hop.ui.core.metadata.IMetadataEditor;
import org.apache.hop.ui.core.metadata.MetadataEditor;
import org.apache.hop.ui.core.metadata.MetadataManager;
import org.apache.hop.ui.hopgui.HopGui;
import org.apache.plc4x.java.DefaultPlcDriverManager;
import org.apache.plc4x.java.api.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.eclipse.swt.widgets.Composite;
import org.eclipse.swt.SWT;
import org.eclipse.swt.layout.FormAttachment;
import org.eclipse.swt.layout.FormData;
import org.eclipse.swt.widgets.*;


public class Plc4xConnectionEditor  extends MetadataEditor<Plc4xConnection>
    implements IMetadataEditor<Plc4xConnection>{

  private static final Class<?> PKG = Plc4xConnectionEditor.class; // For Translator

  public static final String PARENT_WIDGET_ID = "Plc4xConnectionEditor.Widgets.ParentId";

  private Composite parent;
  private Text wName;
  private GuiCompositeWidgets widgets; 
    
    public Plc4xConnectionEditor(HopGui hopGui, MetadataManager<Plc4xConnection> manager, Plc4xConnection metadata) {
        super(hopGui, manager, metadata);
    }

    
    @Override
    public void createControl(Composite parent) {
    this.parent = parent;
    PropsUi props = PropsUi.getInstance();
    int margin = props.getMargin();
    int middle = props.getMiddlePct();
    // Name...
    //
    // What's the name
    Label wlName = new Label(parent, SWT.RIGHT);
    props.setLook(wlName);
    wlName.setText("Plc4x url name");
    FormData fdlName = new FormData();
    fdlName.top = new FormAttachment(0, margin * 2);
    fdlName.left = new FormAttachment(0, 0);
    fdlName.right = new FormAttachment(middle, 0);
    wlName.setLayoutData(fdlName); 
    wName = new Text(parent, SWT.SINGLE | SWT.LEFT | SWT.BORDER);
    props.setLook(wName);
    
    FormData fdName = new FormData();
    fdName.top = new FormAttachment(wlName, 0, SWT.CENTER);
    fdName.left = new FormAttachment(middle, margin);
    fdName.right = new FormAttachment(100, 0);
    wName.setLayoutData(fdName);

    // Rest of the widgets...
    //
    widgets = new GuiCompositeWidgets(manager.getVariables());
    widgets.createCompositeWidgets(getMetadata(), null, parent, PARENT_WIDGET_ID, wName);

    // Set content on the widgets...
    //
    setWidgetsContent();

    // Add changed listeners
    wName.addListener(SWT.Modify, e -> setChanged());
    widgets.setWidgetsListener(
        new GuiCompositeWidgetsAdapter() {
          @Override
          public void widgetModified(
              GuiCompositeWidgets compositeWidgets, Control changedWidget, String widgetId) {
            setChanged();
          }
        });
    }    
    
    
    @Override
    public void setWidgetsContent() {
        Plc4xConnection meta = getMetadata();
        wName.setText(Const.NVL(meta.getName(), ""));
        widgets.setWidgetsContents(meta, parent, PARENT_WIDGET_ID);
    }

    @Override
    public void getWidgetsContent(Plc4xConnection meta) {
        meta.setName(wName.getText());
        widgets.getWidgetsContents(meta, PARENT_WIDGET_ID);
    }

  @Override
  public Button[] createButtonsForButtonBar(Composite parent) {
    PropsUi props = PropsUi.getInstance();

    Button wbTest = new Button(parent, SWT.PUSH | SWT.CENTER);
    props.setLook(wbTest);
    wbTest.setText("Test");
    wbTest.addListener(SWT.Selection, e -> test());

    return new Button[] {wbTest};
  }
  
  public void test() {   
    PlcConnection plcConnection;      
    try {      
        Plc4xConnection meta = getMetadata();
        plcConnection = new DefaultPlcDriverManager().getConnection(meta.getUrl()); 
        plcConnection.connect();       
        Thread.sleep(100);
        plcConnection.close();         
        MessageBox box = new MessageBox(parent.getShell(), SWT.ICON_INFORMATION | SWT.OK);
        box.setText("Success!");
        box.setMessage("It's possible to connect to Device with this metadata!");
        box.open();

    } catch (Exception e) {
      new ErrorDialog(parent.getShell(), "Error", "We couldn't connect using this information", e);
    } 
  }
    
    
}
