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
package org.apache.plc4x.merlot.uns.api;

import java.util.UUID;


public interface Model {
        
    public void init();
    
    public void destroy();
    
    public void start();
    
    public void stop(); 
    
    public void putEnterprise(UUID uuid, String desc);
    public UUID getEnterpriseUUID();
    public String getEnterpriseId();
    public String getEnterpriseDescription();
    
    public void putSite(UUID siteuuid, String id, String desc);
    public void delSite(UUID siteuuid);
    public UUID getSiteUUID(String id);            
    public String getSiteId(UUID uuid);
    public String getSiteDescription(UUID uuid);
    public UUID[] getSites();

    public void putArea(UUID siteuuid, UUID areauuid, String id, String desc);
    public void delArea(UUID areauuid);           
    public UUID getAreaUUID(String id);    
    public String getAreaId(UUID uuid);
    public String getAreaDescription(UUID uuid);
    public UUID[] getAreas(UUID siteuuid);
    
    public void putCell(UUID areauuid, UUID celluuid, String id, String desc);
    public void delCell(UUID celluuid);        
    public UUID getCellUUID(String id);    
    public String getCellId(UUID uuid);
    public String getCellDescription(UUID uuid);
    public UUID[] getCells(UUID areauuid);

    public void putUnit(UUID celluuid, UUID unituuid, String id, String desc);
    public void delUnit(UUID unituuid);     
    public UUID getUnitUUID(String id);
    public String getUnitId(UUID uuid);
    public String getUnitDescription(UUID uuid);
    public UUID[] getUnits(UUID celluuid);   
    
    public void putEmodule(UUID unituuid, UUID emoduleuuid, String id, String desc);
    public void delEmodule(UUID emoduleuuid);     
    public UUID getEmoduleUUID(String id);
    public String getEmoduleId(UUID uuid);
    public String getEmoduleDescription(UUID uuid);
    public UUID[] getEmodules(UUID uuid);

    public void putCmodule(UUID emoduleuuid, UUID cmoduleuuId, String id, String desc);
    public void delCmodule(UUID cmoduleuuid);     
    public UUID getCmoduleUUID(String id);    
    public String getCmoduleId(UUID uuid);
    public String getCmoduleDescription(UUID uuid);
    public UUID[] getCmodules(UUID emoduleuuid);    
    
    public void reload();
    public void persist();
    public void printTree(StringBuilder sb);
    public void delTreeNode(UUID uuid);
    
    /*
    * Retrieves the Domino of the specified object. 
    * The result is a String that follows the style defined in RFC-1034.
    * The String Domain is rebuilt based on the domain tree defined according 
    * to ISA-S88 / Pack-ML.
    * @param id: unique id identification of the object.
    * @return: Domain to which the object belongs.
    */
    public String getDomain(UUID uuid);
    
    public String getDomain(String id);    
    
}
