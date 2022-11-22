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

import javafx.fxml.FXML;
import javafx.scene.control.*;
import javafx.scene.control.cell.TreeItemPropertyValueFactory;
import javafx.scene.input.MouseEvent;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcBrowseItem;
import org.apache.plc4x.java.api.messages.PlcBrowseItemArrayInfo;
import org.apache.plc4x.java.api.messages.PlcBrowseResponse;
import org.apache.plc4x.java.api.model.ArrayInfo;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.kordamp.ikonli.javafx.FontIcon;
import org.kordamp.ikonli.materialdesign.MaterialDesign;

import java.util.ArrayList;
import java.util.Comparator;
import java.util.List;
import java.util.concurrent.ExecutionException;

public class ConnectionTabController {

    @FXML
    public Tab tab;

    @FXML
    public TreeTableView<TreeEntry> resourceTreeView;

    private PlcConnection connection;

    @FXML
    public void initialize() {
        TreeTableColumn<TreeEntry, String> addressColumn = new TreeTableColumn<>("Address");
        addressColumn.setCellValueFactory(new TreeItemPropertyValueFactory<>("address"));
        resourceTreeView.getColumns().add(addressColumn);

        TreeTableColumn<TreeEntry, String> nameColumn = new TreeTableColumn<>("Name");
        nameColumn.setCellValueFactory(new TreeItemPropertyValueFactory<>("name"));
        resourceTreeView.getColumns().add(nameColumn);

        TreeTableColumn<TreeEntry, PlcValueType> typeColumn = new TreeTableColumn<>("Type");
        typeColumn.setCellValueFactory(new TreeItemPropertyValueFactory<>("type"));
        resourceTreeView.getColumns().add(typeColumn);

        TreeTableColumn<TreeEntry, Boolean> readableColumn = new TreeTableColumn<>("Readable");
        readableColumn.setCellValueFactory(new TreeItemPropertyValueFactory<>("readable"));
        resourceTreeView.getColumns().add(readableColumn);

        TreeTableColumn<TreeEntry, Boolean> writableColumn = new TreeTableColumn<>("Writable");
        writableColumn.setCellValueFactory(new TreeItemPropertyValueFactory<>("writable"));
        resourceTreeView.getColumns().add(writableColumn);

        TreeTableColumn<TreeEntry, Boolean> subscribableColumn = new TreeTableColumn<>("Subscribable");
        subscribableColumn.setCellValueFactory(new TreeItemPropertyValueFactory<>("subscribable"));
        resourceTreeView.getColumns().add(subscribableColumn);
    }

    void setConnection(String connectionName, PlcConnection connection) {
        tab.setText(connectionName);

        this.connection = connection;

        if(connection.getMetadata().canBrowse()) {
            try {
                PlcBrowseResponse browseResponse = connection.browseRequestBuilder().build().execute().get();

                TreeItem<ConnectionTabController.TreeEntry> rootItem = new TreeItem<>(
                    new ConnectionTabController.TreeEntry("", "", PlcValueType.NULL,
                        false, false, false));
                rootItem.setGraphic(new FontIcon(MaterialDesign.MDI_FOLDER));
                rootItem.setExpanded(true);

                // Sort the entries first.
                for (String queryName : browseResponse.getQueryNames()) {
                    List<PlcBrowseItem> values = browseResponse.getValues(queryName);
                    values.sort(new PlcBrowseItemComparator());
                    // Then add the elements to the tree.
                    for (PlcBrowseItem value : values) {
                        rootItem.getChildren().add(getTreeItemForBrowseItem(value));
                    }
                }

                resourceTreeView.setRoot(rootItem);
            } catch (ExecutionException | InterruptedException e) {
                throw new RuntimeException(e);
            }
        }
    }

    private TreeItem<ConnectionTabController.TreeEntry> getTreeItemForBrowseItem(PlcBrowseItem browseItem) {
        StringBuilder addressSuffix = new StringBuilder();
        if ((browseItem.getTag().getPlcValueType() == PlcValueType.List) && (browseItem.getTag().getArrayInfo() != null)){
            addressSuffix.append(" ");
            for (ArrayInfo arrayInfo : browseItem.getTag().getArrayInfo()) {
                addressSuffix.append("[").append(arrayInfo.getLowerBound()).append(" .. ").append(arrayInfo.getUpperBound()).append("]");
            }
        }
        TreeItem<ConnectionTabController.TreeEntry> treeItem = new TreeItem<>(new ConnectionTabController.TreeEntry(
            browseItem.getTag().getAddressString() + addressSuffix, browseItem.getName(), browseItem.getTag().getPlcValueType(),
            browseItem.isReadable(), browseItem.isWritable(), browseItem.isSubscribable()));
        if(!browseItem.getChildren().isEmpty()) {
            // Sort the entries first.
            List<PlcBrowseItem> values = new ArrayList<>(browseItem.getChildren().values());
            // Then add the elements to the tree.
            values.sort(new PlcBrowseItemComparator());
            for (PlcBrowseItem child : values) {
                treeItem.getChildren().add(getTreeItemForBrowseItem(child));
            }
        }
        return treeItem;
    }

    @FXML
    public void handleTreeSelectionChanged(MouseEvent mouseEvent) {

    }

    public static class TreeEntry {

        private final String address;

        private final String name;

        private final PlcValueType type;

        private final boolean readable;

        private final boolean writable;

        private final boolean subscribable;

        public TreeEntry(String address, String name, PlcValueType type, boolean readable, boolean writable, boolean subscribable) {
            this.address = address;
            this.name = name;
            this.type = type;
            this.readable = readable;
            this.writable = writable;
            this.subscribable = subscribable;
        }

        public String getAddress() {
            return address;
        }

        public String getName() {
            return name;
        }

        public PlcValueType getType() {
            return type;
        }

        public boolean isReadable() {
            return readable;
        }

        public boolean isWritable() {
            return writable;
        }

        public boolean isSubscribable() {
            return subscribable;
        }
    }

    class PlcBrowseItemComparator implements Comparator<PlcBrowseItem> {
        @Override
        public int compare(PlcBrowseItem o1, PlcBrowseItem o2) {
            return o1.getTag().getAddressString().compareTo(o2.getTag().getAddressString());
        }
    }

}
