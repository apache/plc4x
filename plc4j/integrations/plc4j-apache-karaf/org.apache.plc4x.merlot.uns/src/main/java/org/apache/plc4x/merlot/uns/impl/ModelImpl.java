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
package org.apache.plc4x.merlot.uns.impl;

import org.apache.plc4x.merlot.uns.api.Model;
import java.io.IOException;
import java.util.Dictionary;
import java.util.Enumeration;
import java.util.HashMap;
import java.util.Map;
import java.util.UUID;
import java.util.logging.Level;
import org.apache.karaf.config.core.ConfigRepository;
import org.osgi.framework.BundleContext;

public class ModelImpl implements Model {
    
    private static final String CONFIG_PID = "org.apache.plc4x.merlot.uns";    
    private final BundleContext bundleContext;
    private final ConfigRepository configRepository;
    private Dictionary<String, Object> dict = null;
    
    private String enterprise = "root";
    private String plant = "plant";
    private String area = "area";
    private String cell = "cell";
    private String unit = "unit";
    private String domain = null;
    private Node root = null;

    public ModelImpl(BundleContext bundleContext, ConfigRepository configRepository) {
        this.bundleContext = bundleContext;
        this.configRepository = configRepository;
    }
           
    @Override
    public void init() {
        try {
            dict = configRepository
                    .getConfigAdmin()
                    .getConfiguration(CONFIG_PID)
                    .getProperties();
            root = new Node();
            root.id = "root";
            root.nodes = new HashMap();
            reload(root);
        } catch (IOException ex) {
//            LOGGER.info(ex.getMessage());
        }
    }

    @Override
    public void destroy() {
    }

    @Override
    public void start() {
    }

    @Override
    public void stop() {
    }

    @Override
    public void putEnterprise(UUID uuid, String desc) {
//        if (root != null) LOGGER.info("All the 'root' model will be created.");
        root = new Node();
        root.id ="root";
        root.description = desc;
        root.uuid = uuid;
        root.nodes = new HashMap();        
    }

    @Override
    public UUID getEnterpriseUUID() {
        return root.uuid;
    }    
    
    @Override
    public String getEnterpriseId() {
        return root.id;
    }

    @Override
    public String getEnterpriseDescription() {
        return root.description;
    }

    @Override
    public void putSite(UUID plantuuid, String id, String desc) {
        addNode(root.uuid, plantuuid, id, desc);
    }

    @Override
    public void delSite(UUID siteuuid) {
        delNode(root,siteuuid);
    }
            
    @Override
    public UUID getSiteUUID(String id) {
        return getUUID(id);
    }    
    
    @Override
    public String getSiteId(UUID uuid) {
        return getId(uuid);
    }

    @Override
    public String getSiteDescription(UUID uuid) {
        return getDescription(uuid);
    }

    @Override
    public UUID[] getSites() {
        return getUUIDS(root.uuid);
    }

    @Override
    public void putArea(UUID plantuuid, UUID areauuid, String id, String desc) {
        addNode(plantuuid, areauuid, id, desc);
    }

    @Override
    public void delArea(UUID areauuid) {
        delNode(root, areauuid);
    }
        
    @Override
    public UUID getAreaUUID(String id) {
        return getUUID(id);
    }

    @Override
    public String getAreaId(UUID uuid) {
        return getId(uuid);
    }

    @Override
    public String getAreaDescription(UUID uuid) {
        return getDescription(uuid);
    }

    @Override
    public UUID[] getAreas(UUID plantuuid) {
        return getUUIDS(plantuuid);
    }

    @Override
    public void putCell(UUID areauuid, UUID celluuid, String id, String desc) {
        addNode(areauuid, celluuid, id, desc);
    }

    @Override
    public void delCell(UUID celluuid) {
        delNode(root, celluuid);
    }
        
    @Override
    public UUID getCellUUID(String id) {
        return getUUID(id);
    }

    @Override
    public String getCellId(UUID uuid) {
        return getId(uuid);
    }

    @Override
    public String getCellDescription(UUID uuid) {
        return getDescription(uuid);
    }

    @Override
    public UUID[] getCells(UUID uuid) {
        return getUUIDS(uuid);
    }

    @Override
    public void putUnit(UUID celluuid, UUID unituuid, String id, String desc) {
        addNode(celluuid, unituuid, id, desc);
    }

    @Override
    public void delUnit(UUID unituuid) {
        delNode(root, unituuid);
    }
        
    @Override
    public UUID getUnitUUID(String id) {
        return getUUID(id);
    }

    @Override
    public String getUnitId(UUID uuid) {
        return getId(uuid);
    }

    @Override
    public String getUnitDescription(UUID uuid) {
        return getDescription(uuid);
    }

    @Override
    public UUID[] getUnits(UUID uuid) {
        return getUUIDS(uuid);
    }

    @Override
    public void putEmodule(UUID unituuid, UUID emoduleuuid, String id, String desc) {
        addNode(unituuid, emoduleuuid, id, desc);
    }

    @Override
    public void delEmodule(UUID emoduleuuid) {
        delNode(root, emoduleuuid);
    }
        
    @Override
    public UUID getEmoduleUUID(String id) {
        return getUUID(id);
    }
            
    @Override
    public String getEmoduleId(UUID uuid) {
        return getId(uuid);
    }

    @Override
    public String getEmoduleDescription(UUID uuid) {
        return getDescription(uuid);
    }

    @Override
    public UUID[] getEmodules(UUID uuid) {
        return getUUIDS(uuid);
    }

    @Override
    public void putCmodule(UUID emoduleuuid, UUID cmoduleuuId, String id, String desc) {
        addNode(emoduleuuid, cmoduleuuId, id, desc);
    }

    @Override
    public void delCmodule(UUID cmoduleuuid) {
        delNode(root, cmoduleuuid);
    }

    @Override
    public UUID getCmoduleUUID(String id) {
        return getUUID(id);
    }

    @Override
    public String getCmoduleId(UUID uuid) {
        return getId(uuid);
    }

    @Override
    public String getCmoduleDescription(UUID uuid) {
        return getDescription(uuid);
    }

    @Override
    public UUID[] getCmodules(UUID emoduleuuid) {
        return getUUIDS(emoduleuuid);
    }
       
    @Override
    public String getDomain(UUID uuid) {
        return getUuidDomain(uuid);
    }

    @Override
    public String getDomain(String id) {
        return getIdDomain(id);
    }
    
        
    public void reload(){
        try {
            root = new Node();
            root.id = "root";
            root.nodes = new HashMap();
            root.uuid = null;
            reload(root);
        } catch (Exception ex){
            //LOGGER.info(ex.getMessage());
        }        
    }
    
    public void persist(){        
        Enumeration<String> keys = dict.keys();
        String[] delKeys = new String[dict.size()];
        String key = null;
        int i = 0;
        while(keys.hasMoreElements()){  
            key = keys.nextElement();
            delKeys[i] = key;
            i++;
        }
        for (String delKey:delKeys) dict.remove(delKey);
        try {
            persist(root);
            keys = dict.keys();
            Map<String, Object> props = new HashMap();

            while(keys.hasMoreElements()){
                key = keys.nextElement();
                props.put(key, dict.get(key));
            }
            configRepository.update(plant, props);          
        } catch (Exception ex){
//            LOGGER.info(ex.getMessage());            
        }
    }
    
    public void printTree(StringBuilder sb){
        printTree(sb, "", "", root);
    }

    @Override
    public void delTreeNode(UUID uuid) {
        delNode(root, uuid);
    }
    
    
    
    
    /*
    * Support rutines
    */
    
    private void reload(Node node){
        try {
            if (dict != null){  
                String strprops = null;
 
                if ("root".equals(node.id)){
                    strprops = (String) dict.get(node.id);
                } else if (node.uuid != null){
                    strprops = (String) dict.get(node.uuid.toString());
                }

                String[] fields = strprops.split(",", 4);
                node.id  = fields[0];
                node.description = fields[1];
                node.uuid = UUID.fromString(fields[2]);

                if (fields.length > 3){
                    String[] strkeys = fields[3].split(",");
                    for (String key:strkeys){
                        Node newnode = new Node();
                        newnode.uuid = UUID.fromString(key);
                        newnode.nodes = new HashMap();
                        reload(newnode);
                        node.nodes.put(newnode.uuid, newnode);
                    }
                }


            };
        } catch (Exception ex){
//            LOGGER.info(ex.getMessage());
        }           
    }
    
    private void persist(Node node){ 
        try{
            if ("root".equals(node.id)){
                String strprops = node.id+","+node.description+","+node.uuid.toString();
                if (node.nodes.size()>0){
                    UUID[] uuids = getUUIDS(node.uuid);
                    for (UUID uuid:uuids){
                        strprops = strprops.concat(",").concat(uuid.toString());
                        persist(findNode(uuid));
                    }
                }
                dict.put(node.id, strprops);
            } else {
                String strprops = node.id+","+node.description+","+node.uuid.toString();
                if (node.nodes.size()>0){
                    
                    UUID[] uuids = getUUIDS(node.uuid);
                    for (UUID uuid:uuids){
                        strprops = strprops.concat(",").concat(uuid.toString());
                        persist(findNode(uuid));
                    }
                }          
                dict.put(node.uuid.toString(), strprops);                
            }
            
        } catch (Exception ex){
//            LOGGER.info(ex.getMessage());            
        }
    }
    
    private void printTree(StringBuilder sb, String padding, String pointer, Node node){
        if (node != null) {
            sb.append(padding);
	    sb.append(pointer);
            sb.append(node.id+"-{"+node.uuid.toString()+"}");
            sb.append("\n");
            
            StringBuilder paddingBuilder = new StringBuilder(padding);
	    paddingBuilder.append("|  ");
            
            String paddingForBoth = paddingBuilder.toString();
            String pointerForRight = "|---";
	    String pointerForLeft = (node.nodes.size() != 0)? "+--" : "";
            
            for (Node nextnode:node.nodes.values()){
                printTree(sb, paddingForBoth, pointerForLeft, nextnode);
            }
                                	 
        }
    }
            
    private void addNode(UUID upuuid, UUID nodeuuid, String id, String desc) {
        Node upnode = findNode(upuuid);
        if (upnode != null){
            Node node = new Node();
            node.uuid = nodeuuid;
            node.id = id;
            node.description = desc;
            node.nodes =  new HashMap(); 
            upnode.nodes.put(nodeuuid, node);              
        }        
    }
    
    private UUID getUUID(String id){
        Node node = findNode(id);
        if (node != null) return node.uuid;
        return null;        
    }
    
    public String getId(UUID uuid) {
        Node node = findNode(uuid);
        if (node != null) return node.id;
        return null;
    }
    
    public String getDescription(UUID uuid) {
        Node node = findNode(uuid);
        if (node != null) return node.description;
        return null;
    }    
    
    public UUID[] getUUIDS(UUID uuid) {
        Node node = findNode(uuid);        
        if (node != null){
            UUID[] uuids = new UUID[node.nodes.size()];
            int i = 0;
            for (Node emodule:node.nodes.values()){
                uuids[i] = emodule.uuid;
                i++;
            }
            return uuids;
        }
        return null;
    }    
    
    private Node findNode(UUID uuid){
        if (uuid.equals(root.uuid)) return root;
        for (Node site:root.nodes.values()){
            if (uuid.equals(site.uuid)) return site;
            for (Node area:site.nodes.values()){
                if (uuid.equals(area.uuid)) return area;
                for (Node cell:area.nodes.values()){
                    if (uuid.equals(cell.uuid)) return cell;
                    for (Node unit:cell.nodes.values()){
                        if (uuid.equals(unit.uuid)) return unit;
                        for (Node emodule:unit.nodes.values()){
                            if (uuid.equals(emodule.uuid)) return emodule;
                            for (Node cmodule:emodule.nodes.values()){                            
                                if (uuid.equals(cmodule.uuid)) return cmodule;
                            }
                        }
                    }                  
                }
            }
        }
        return null;        
    }
    
    private Node findNode(String id){
        if (id.equals("root")) return root;
        for (Node site:root.nodes.values()){
            if (id.equals(site.id)) return site;
            for (Node area:site.nodes.values()){
                if (id.equals(area.id)) return area;
                for (Node cell:area.nodes.values()){
                    if (id.equals(cell.id)) return cell;
                    for (Node unit:cell.nodes.values()){
                        if (id.equals(unit.id)) return unit;
                        for (Node emodule:unit.nodes.values()){
                            if (id.equals(emodule.id)) return emodule;
                            for (Node cmodule:emodule.nodes.values()){                            
                                if (id.equals(cmodule.id)) return cmodule;
                            }
                        }
                    }                  
                }
            }
        }
        return null;        
    }
    
    private void delNode(Node node, UUID uuid){
        for (Node nextnode:node.nodes.values()){           
            if (uuid.equals(nextnode.uuid)) {             
                node.nodes.remove(nextnode.uuid);
                return;
            }
            delNode(nextnode,uuid);
        }           
    }

    private String getUuidDomain(UUID uuid){
        for (Node site:root.nodes.values()){
            if (uuid.equals(site.uuid)) {
                return root.id+"."+site.id;
            }
            for (Node area:site.nodes.values()){
                if (uuid.equals(area.uuid)) {
                    return root.id+"."+site.id+"."+area.id;
                }
                for (Node cell:area.nodes.values()){
                    if (uuid.equals(cell.uuid)) {
                        return root.id+"."+site.id+"."+area.id+"."+cell.id;                        
                    }
                    for (Node unit:cell.nodes.values()){
                        if (uuid.equals(unit.uuid)) {
                            return root.id+"."+site.id+"."+area.id+"."+cell.id+"."+unit.id;                            
                        }
                        for (Node emodule:unit.nodes.values()){
                            if (uuid.equals(emodule.uuid)) {
                                return root.id+"."+site.id+"."+area.id+"."+cell.id+"."+unit.id+"."+emodule.id;                                 
                            }
                            for (Node cmodule:emodule.nodes.values()){                            
                                if (uuid.equals(cmodule.uuid)) {                                   
                                    return root.id+"."+site.id+"."+area.id+"."+cell.id+"."+unit.id+"."+emodule.id+"."+cmodule.id;                                       
                                }
                            }
                        }
                    }                  
                }
            }
        }
        return null;        
    }    
    
    private String getIdDomain(String id){
        for (Node site:root.nodes.values()){
            if (id.equals(site.id)) {
                return root.id+"."+site.id;
            }
            for (Node area:site.nodes.values()){
                if (id.equals(area.id)) {
                    return root.id+"."+site.id+"."+area.id;
                }
                for (Node cell:area.nodes.values()){
                    if (id.equals(cell.id)) {
                        return root.id+"."+site.id+"."+area.id+"."+cell.id;  
                    }
                    for (Node unit:cell.nodes.values()){
                        if (id.equals(unit.id)) {
                            return root.id+"."+site.id+"."+area.id+"."+cell.id+"."+unit.id;   
                        }
                        for (Node emodule:unit.nodes.values()){
                            if (id.equals(emodule.id)) {
                                return root.id+"."+site.id+"."+area.id+"."+cell.id+"."+unit.id+"."+emodule.id;    
                            }
                            for (Node cmodule:emodule.nodes.values()){                            
                                if (id.equals(cmodule.id)) {
                                    return root.id+"."+site.id+"."+area.id+"."+cell.id+"."+unit.id+"."+emodule.id+"."+cmodule.id;  
                                }
                            }
                        }
                    }                  
                }
            }
        }
        return null;        
    }         
 
    class Node {
        public UUID uuid;
        public String id;
        public String description;
        public HashMap<UUID,Node> nodes;
    }   

    
}
