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
import java.util.ArrayList;
import org.apache.hop.core.ICheckResult;
import org.apache.hop.core.annotations.Transform;
import org.apache.hop.core.exception.HopTransformException;
import org.apache.hop.core.exception.HopXmlException;
import org.apache.hop.core.row.IRowMeta;
import org.apache.hop.core.variables.IVariables;
import org.apache.hop.core.xml.XmlHandler;
import org.apache.hop.metadata.api.IHopMetadataProvider;
import org.apache.hop.pipeline.Pipeline;
import org.apache.hop.pipeline.PipelineMeta;
import org.apache.hop.pipeline.transform.BaseTransformMeta;
import org.apache.hop.pipeline.transform.ITransformMeta;
import org.apache.hop.pipeline.transform.TransformMeta;
import org.w3c.dom.Node;

import java.util.List;
import org.apache.hop.core.Const;
import org.apache.hop.core.RowMetaAndData;
import org.apache.hop.core.row.IValueMeta;
import org.apache.hop.metadata.api.HopMetadataProperty;
import org.apache.hop.pipeline.transform.ITransformIOMeta;
import org.apache.hop.pipeline.transform.TransformIOMeta;
import org.apache.plc4x.hop.transforms.util.Plc4xGeneratorField;

/**
 * Meta data for the sample transform.
 */
@Transform(
        id = "plc4x-subs",
        name = "i18n::Plc4x.Read.Name",
        description = "i18n::Plc4x.Read.Description",
        image = "plc4x_subs.svg",
        categoryDescription = "i18n:org.apache.plc4x.hop.transforms.plc4xinput:Plc4x.Category.plc4x",
        documentationUrl = "https://plc4x.apache.org/users/integrations/apache-calcite.html"
)
public class Plc4xSubsMeta extends BaseTransformMeta implements ITransformMeta {

  private static final Class<?> PKG = Plc4xSubsMeta.class; // Needed by Translator

  
  @HopMetadataProperty(
      key = "connection",
      injectionKeyDescription = "Plc4x.Read.Meta.Injection.Connection")
  private String connection;  

  @HopMetadataProperty(
      key = "never_ending",
      injectionKeyDescription = "Plc4x.Read.Meta.Injection.NeverEnding")
  private boolean neverEnding;

  @HopMetadataProperty(
      key = "maxwait_in_ms",
      injectionKeyDescription = "Plc4x.Read.Meta.Injection.MaxwaitInMs")
  private String maxwaitInMs;  
 
  @HopMetadataProperty(
      key = "interval_in_ms",
      injectionKeyDescription = "Plc4x.Read.Meta.Injection.IntervalInMs")
  private String intervalInMs;

  @HopMetadataProperty(
      key = "row_time_field",
      injectionKeyDescription = "Plc4x.Read.Meta.Injection.RowTimeField")
  private String rowTimeField;

  @HopMetadataProperty(
      key = "last_time_field",
      injectionKeyDescription = "Plc4x.Read.Meta.Injection.LastTimeField")
  private String lastTimeField;

  @HopMetadataProperty(
      key = "limit",
      injectionKeyDescription = "Plc4x.Read.Meta.Injection.RowLimit")
  private String rowLimit;

  @HopMetadataProperty(
      groupKey = "fields",
      key = "field",
      injectionGroupDescription = "Plc4x.Read.Meta.Injection.Fields",
      injectionKeyDescription = "Plc4x.Read.Meta.Injection.Fields")
  private List<Plc4xGeneratorField> fields;

  public Plc4xSubsMeta() {
    fields = new ArrayList<>();

    rowLimit = "10";
    neverEnding = false;
    maxwaitInMs = "1000";
    intervalInMs = "5000";
    rowTimeField = "now";
    lastTimeField = "last";
  }  
  
  public Plc4xSubsMeta(Plc4xSubsMeta m) {
    this.connection = m.connection;
    this.neverEnding = m.neverEnding;
    this.intervalInMs = m.intervalInMs;
    this.rowTimeField = m.rowTimeField;
    this.lastTimeField = m.lastTimeField;
    this.rowLimit = m.rowLimit;
    this.fields = new ArrayList<>();
    for (Plc4xGeneratorField field : m.fields) {
      this.fields.add(new Plc4xGeneratorField(field));
    }
  }

  @Override
  public Plc4xSubsMeta clone() {
    return new Plc4xSubsMeta(this);
  }  
  
  

  @Override
  public void getFields( IRowMeta inputRowMeta, String name, IRowMeta[] info, TransformMeta nextTransform,
                         IVariables variables, IHopMetadataProvider metadataProvider ) throws HopTransformException {
    try {
        logBasic("PASO 1");
        List<ICheckResult> remarks = new ArrayList<>();
        RowMetaAndData rowMetaAndData = Plc4xSubs.buildRow(this, remarks, name);
        logBasic("PASO 2");
      if (!remarks.isEmpty()) {
        logBasic("PASO 3");
        StringBuilder stringRemarks = new StringBuilder();
        for (ICheckResult remark : remarks) {
          logBasic("PASO 3x: " + remark.toString());   
          stringRemarks.append(remark.toString()).append(Const.CR);
        }
        logBasic("PASO 3.1");        
        throw new HopTransformException(stringRemarks.toString());
      }
              logBasic("PASO 4");
      for (IValueMeta valueMeta : rowMetaAndData.getRowMeta().getValueMetaList()) {
        valueMeta.setOrigin(name);
      }
        logBasic("PASO 5");
      inputRowMeta.mergeRowMeta(rowMetaAndData.getRowMeta());
    } catch (Exception e) {
      throw new HopTransformException(e);
    }
  }


  public Plc4xSubs createTransform(TransformMeta transformMeta, Plc4xSubsData data, int copyNr,
                                PipelineMeta pipelineMeta, Pipeline pipeline ) {
    return new Plc4xSubs( transformMeta, this, data, copyNr, pipelineMeta, pipeline );
  }

  public Plc4xSubsData getTransformData() {
    return new Plc4xSubsData();
  }
/*
  @Override
  public void loadXml( Node transformNode, IHopMetadataProvider metadataProvider ) throws HopXmlException {
    //load the saved values from the transformnode
    String sampleValue = XmlHandler.getTagValue( transformNode, "sampleValue" );

  }
*/
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
  
    /**
   * Gets intervalInMs
   *
   * @return value of intervalInMs
   */
  public String getMaxwaitInMs() {
    return maxwaitInMs;
  }

  /** @param intervalInMs The intervalInMs to set */
  public void setMaxwaitInMs(String maxwaitInMs) {
    this.maxwaitInMs = maxwaitInMs;
  }

  /**
   * Gets intervalInMs
   *
   * @return value of intervalInMs
   */
  public String getIntervalInMs() {
    return intervalInMs;
  }

  /** @param intervalInMs The intervalInMs to set */
  public void setIntervalInMs(String intervalInMs) {
    this.intervalInMs = intervalInMs;
  }

  /**
   * Gets rowTimeField
   *
   * @return value of rowTimeField
   */
  public String getRowTimeField() {
    return rowTimeField;
  }

  /** @param rowTimeField The rowTimeField to set */
  public void setRowTimeField(String rowTimeField) {
    this.rowTimeField = rowTimeField;
  }

  /**
   * Gets lastTimeField
   *
   * @return value of lastTimeField
   */
  public String getLastTimeField() {
    return lastTimeField;
  }

  /** @param lastTimeField The lastTimeField to set */
  public void setLastTimeField(String lastTimeField) {
    this.lastTimeField = lastTimeField;
  }

  /**
   * Gets rowLimit
   *
   * @return value of rowLimit
   */
  public String getRowLimit() {
    return rowLimit;
  }

  /** @param rowLimit The rowLimit to set */
  public void setRowLimit(String rowLimit) {
    this.rowLimit = rowLimit;
  }

  /**
   * Gets fields
   *
   * @return value of fields
   */
  public List<Plc4xGeneratorField> getFields() {
    return fields;
  }  
  
  /** @param fields The fields to set */
  public void setFields(List<Plc4xGeneratorField> fields) {
    this.fields = fields;
  }  
  
}
