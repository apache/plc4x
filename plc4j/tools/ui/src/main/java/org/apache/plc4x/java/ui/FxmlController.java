/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.ui;

import javafx.event.ActionEvent;
import javafx.fxml.FXML;
import javafx.fxml.FXMLLoader;
import javafx.scene.control.*;
import javafx.scene.control.cell.TextFieldTreeCell;
import javafx.scene.input.MouseButton;
import javafx.scene.input.MouseEvent;
import javafx.scene.paint.Paint;

import javafx.util.StringConverter;
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.PlcDriver;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.messages.PlcDiscoveryItem;
import org.apache.plc4x.java.api.messages.PlcDiscoveryResponse;
import org.kordamp.ikonli.javafx.FontIcon;
import org.kordamp.ikonli.materialdesign.MaterialDesign;

import java.io.IOException;
import java.util.Objects;
import java.util.concurrent.CompletableFuture;

public class FxmlController {

    private final PlcDriverManager driverManager;

    @FXML
    public TreeView<TreeEntry> plcTreeView;

    @FXML
    public Button browseButton;

    @FXML
    public TabPane connectionTabs;

    public FxmlController() {
        driverManager = new PlcDriverManager();
    }

    @FXML
    public void initialize() throws Exception {
        plcTreeView.setCellFactory(treeEntryTreeView -> new TextFieldTreeCell<>(new StringConverter<TreeEntry>(){
            @Override
            public String toString(TreeEntry treeEntry) {
                return treeEntry.getName();
            }

            @Override
            public TreeEntry fromString(String string) {
                return null;
            }
        }));
        TreeItem<TreeEntry> rootItem = new TreeItem<>(new TreeEntry(
            TreeEntryType.ROOT, "", "Available Drivers"));
        rootItem.setGraphic(new FontIcon(MaterialDesign.MDI_FOLDER));
        for (String protocolCode : driverManager.listDrivers()) {
            PlcDriver driver = driverManager.getDriver(protocolCode);
            TreeItem<TreeEntry> driverItem = new TreeItem<>(new TreeEntry(
                TreeEntryType.DRIVER, driver.getProtocolCode(), driver.getProtocolName()));
            FontIcon icon;
            if(driver.getMetadata().canDiscover()) {
                icon = new FontIcon(MaterialDesign.MDI_CHECK_CIRCLE);
                icon.setIconColor(Paint.valueOf("green"));
            } else {
                icon = new FontIcon(MaterialDesign.MDI_MINUS_CIRCLE);
                icon.setIconColor(Paint.valueOf("red"));
            }
            driverItem.setGraphic(icon);
            rootItem.getChildren().add(driverItem);
        }
        rootItem.setExpanded(true);
        plcTreeView.setRoot(rootItem);

        browseButton.setDisable(true);
    }

    @FXML
    public void handleTreeSelectionChanged(MouseEvent mouseEvent) {
        MultipleSelectionModel<TreeItem<TreeEntry>> selectionModel = plcTreeView.getSelectionModel();
        boolean buttonEnabled = false;
        if(!selectionModel.getSelectedItems().isEmpty()) {
            TreeEntry selectedItem = selectionModel.getSelectedItems().get(0).getValue();
            switch (selectedItem.getType()) {
                case DRIVER:
                    try {
                        buttonEnabled = driverManager.getDriver(selectedItem.getCode()).getMetadata().canDiscover();
                    } catch (Exception e) {
                        // Ignore ...
                    }
                    // If the item was double-clicked, start the scan right away.
                    if(buttonEnabled && mouseEvent.getButton().equals(MouseButton.PRIMARY) && mouseEvent.getClickCount() == 2) {
                        try {
                            handleBrowseButtonClicked(null);
                        } catch (Exception e) {
                            // Ignore ...
                        }
                    }
                    break;
                case PLC:
                    if(mouseEvent.getButton().equals(MouseButton.PRIMARY) && mouseEvent.getClickCount() == 2) {
                        String connectionString = selectedItem.getCode();
                        try {
                            PlcConnection connection = driverManager.getConnection(connectionString);
                            FXMLLoader loader = new FXMLLoader(Objects.requireNonNull(getClass().getResource("connection-tab.fxml")));

                            Tab connectionTab = loader.load();
                            ConnectionTabController controller = loader.getController();
                            controller.setConnection(selectedItem.getName(), connection);

                            connectionTabs.getTabs().add(connectionTab);
                        } catch (PlcConnectionException e) {
                            throw new RuntimeException(e);
                        } catch (IOException e) {
                            throw new RuntimeException(e);
                        }
                    }
                    break;
                case ADDRESS:
                    buttonEnabled = true;
                    break;
            }
        }
        browseButton.setDisable(!buttonEnabled);
    }

    @FXML
    public void handleBrowseButtonClicked(ActionEvent actionEvent) throws Exception {
        MultipleSelectionModel<TreeItem<TreeEntry>> selectionModel = plcTreeView.getSelectionModel();
        if(!selectionModel.getSelectedItems().isEmpty()) {
            TreeItem<TreeEntry> selectedTreeItem = selectionModel.getSelectedItems().get(0);
            TreeEntry selectedItem = selectedTreeItem.getValue();
            selectedTreeItem.setExpanded(true);
            PlcDriver driver = driverManager.getDriver(selectedItem.getCode());
            CompletableFuture<? extends PlcDiscoveryResponse> browseFuture = driver.discoveryRequestBuilder().build().execute();
            browseFuture.whenComplete((response, throwable) -> {
                if(throwable == null) {
                    for (PlcDiscoveryItem discoveredPlc : response.getValues()) {
                        TreeItem<TreeEntry> plcItem = new TreeItem<>(new TreeEntry(
                            TreeEntryType.PLC, discoveredPlc.getConnectionUrl(), discoveredPlc.getName()));
                        plcItem.setGraphic(new FontIcon(MaterialDesign.MDI_LAN_CONNECT));
                        selectedTreeItem.getChildren().add(plcItem);
                    }
                }
            });
        }
    }

    public enum TreeEntryType {
        ROOT,
        DRIVER,
        PLC,
        ADDRESS
    }

    public static class TreeEntry {

        private final TreeEntryType type;

        private final String code;

        private final String name;

        public TreeEntry(TreeEntryType type, String code, String name) {
            this.type = type;
            this.code = code;
            this.name = name;
        }

        public TreeEntryType getType() {
            return type;
        }

        public String getCode() {
            return code;
        }

        public String getName() {
            return name;
        }

    }

}
