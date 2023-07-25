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

package org.apache.plc4x.hop.transforms.plc4xevent;

import java.util.ArrayList;
import org.apache.hop.core.ICheckResult;
import org.apache.hop.core.annotations.Transform;
import org.apache.hop.core.exception.HopTransformException;
import org.apache.hop.core.row.IRowMeta;
import org.apache.hop.core.variables.IVariables;
import org.apache.hop.metadata.api.IHopMetadataProvider;
import org.apache.hop.pipeline.Pipeline;
import org.apache.hop.pipeline.PipelineMeta;
import org.apache.hop.pipeline.transform.BaseTransformMeta;
import org.apache.hop.pipeline.transform.ITransformMeta;
import org.apache.hop.pipeline.transform.TransformMeta;

import java.util.List;
import org.apache.hop.core.Const;
import org.apache.hop.core.RowMetaAndData;
import org.apache.hop.core.row.IValueMeta;
import org.apache.hop.metadata.api.HopMetadataProperty;
import org.apache.hop.pipeline.transform.ITransformIOMeta;
import org.apache.hop.pipeline.transform.TransformIOMeta;

/**
 * Meta data for the sample transform.
 */
@Transform(
        id = "plc4x-event",
        name = "i18n::Plc4x.Read.Name",
        description = "i18n::Plc4x.Read.Description",
        image = "plc4x_event.svg",
        categoryDescription = "i18n:org.apache.plc4x.hop.transforms.plc4xevent:Plc4x.Category.plc4x",
        documentationUrl = "https://plc4x.apache.org/users/integrations/apache-hop.html"
)
public class Plc4xEventMeta extends BaseTransformMeta<Plc4xEvent, Plc4xEventData> {

  private static final Class<?> PKG = Plc4xEventMeta.class; // Needed by Translator

  
  @HopMetadataProperty(
      key = "connection",
      injectionKeyDescription = "Plc4x.Read.Meta.Injection.Connection")
  private String connection;  

  @HopMetadataProperty(
      key = "limit",
      injectionKeyDescription = "Plc4x.Read.Meta.Injection.RowLimit")
  private String rowLimit;  
  
  @HopMetadataProperty(
      key = "never_ending",
      injectionKeyDescription = "Plc4x.Read.Meta.Injection.NeverEnding")
  private boolean neverEnding;
  
  @HopMetadataProperty(
      key = "mode_event",
      injectionKeyDescription = "Plc4x.Read.Meta.Injection.ModeEvent")
  private boolean modeEvent;

  @HopMetadataProperty(
      key = "user_event",
      injectionKeyDescription = "Plc4x.Read.Meta.Injection.UserEvent")
  private boolean userEvent;

  @HopMetadataProperty(
      key = "sys_event",
      injectionKeyDescription = "Plc4x.Read.Meta.Injection.SysEvent")
  private boolean sysEvent;

  @HopMetadataProperty(
      key = "alarm_event",
      injectionKeyDescription = "Plc4x.Read.Meta.Injection.AlarmEvent")
  private boolean alarmEvent;  


  public Plc4xEventMeta() {
    rowLimit = "10";
    neverEnding = false;
    modeEvent = false;
    userEvent = false;
    sysEvent = false;
    alarmEvent = false;

  }  
  
  public Plc4xEventMeta(Plc4xEventMeta m) {
    this.connection = m.connection;
    this.neverEnding = m.neverEnding;
    this.modeEvent = m.modeEvent;
    this.userEvent = m.userEvent;
    this.sysEvent = m.sysEvent;
    this.alarmEvent = m.alarmEvent;
  }

  @Override
  public Plc4xEventMeta clone() {
    return new Plc4xEventMeta(this);
  }  
  
  

  @Override
  public void getFields( IRowMeta inputRowMeta, String name, IRowMeta[] info, TransformMeta nextTransform,
                         IVariables variables, IHopMetadataProvider metadataProvider ) throws HopTransformException {
    try {
        List<ICheckResult> remarks = new ArrayList<>();
        RowMetaAndData rowMetaAndData = Plc4xEvent.buildRow(this, remarks, name);
      if (!remarks.isEmpty()) {
        StringBuilder stringRemarks = new StringBuilder();
        for (ICheckResult remark : remarks) {
          stringRemarks.append(remark.toString()).append(Const.CR);
        }    
        throw new HopTransformException(stringRemarks.toString());
      }
      for (IValueMeta valueMeta : rowMetaAndData.getRowMeta().getValueMetaList()) {
        valueMeta.setOrigin(name);
      }
      inputRowMeta.mergeRowMeta(rowMetaAndData.getRowMeta());
    } catch (Exception e) {
      throw new HopTransformException(e);
    }
  }

  public Plc4xEvent createTransform(TransformMeta transformMeta, Plc4xEventData data, int copyNr,
                                PipelineMeta pipelineMeta, Pipeline pipeline ) {
    return new Plc4xEvent( transformMeta, this, data, copyNr, pipelineMeta, pipeline );
  }

    /**
     *
     * @return
     */
  public Plc4xEventData getTransformData() {
    return new Plc4xEventData();
  }

  @Override
  public void setDefault() {
    //default values when creating a new transform
  }
  
  
  /**
   * Returns the Input/Output metadata for this transform. 
   * The generator transform only produces output, does not accept input!
   * TransformIOMeta(inputAcceptor, outputProducer, inputOptional, outputDynamic, inputDynamic)
   */
  @Override
  public ITransformIOMeta getTransformIOMeta() {
    return new TransformIOMeta(false, true, false, false, false, false);
  }

  /**
   * Gets Plc4xConnection metadata name.
   *
   * @return value of intervalInMs
   */
  public String getConnection() {
    return connection;
  }

  /** @param connection  */
  public void setConnection(String connection) {
    this.connection = connection;
  }  
  
  public String getRowLimit() {
    return rowLimit;
  }

  /** @param connection  */
  public void setRowLimit(String rowLimit) {
    this.rowLimit = rowLimit;
  }  

  /**
   * Gets neverEnding
   *
   * @return value of neverEnding
   */
  public boolean isNeverEnding() {
    return neverEnding;
  }

  /** @param neverEnding The neverEnding to set */
  public void setNeverEnding(boolean neverEnding) {
    this.neverEnding = neverEnding;
  }
  
  public boolean isModeEvent() {
    return modeEvent;
  }

  public void setModeEvent(boolean modeEvent) {
    this.modeEvent = modeEvent;
  }
  
  public boolean isUserEvent() {
    return userEvent;
  }

  public void setUserEvent(boolean userEvent) {
    this.userEvent = userEvent;
  }

  public boolean isSysEvent() {
    return sysEvent;
  }

  public void setSysEvent(boolean sysEvent) {
    this.sysEvent = sysEvent;
  }  
  
  public boolean isAlarmEvent() {
    return alarmEvent;
  }

  public void setAlarmEvent(boolean alarmEvent) {
    this.alarmEvent = alarmEvent;
  } 
  
}
