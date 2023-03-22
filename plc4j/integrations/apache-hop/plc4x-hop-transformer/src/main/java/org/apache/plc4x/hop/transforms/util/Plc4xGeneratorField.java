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

package org.apache.plc4x.hop.transforms.util;

import org.apache.hop.metadata.api.HopMetadataProperty;

import java.util.Objects;

public class Plc4xGeneratorField {

  @HopMetadataProperty(injectionKeyDescription = "Plc4x.Read.Meta.Injection.Field.Name")
  private String name;
  
  @HopMetadataProperty(injectionKeyDescription = "Plc4x.Read.Meta.Injection.Field.Item")
  private String item;  

  @HopMetadataProperty(injectionKeyDescription = "Plc4x.Read.Meta.Injection.Field.Type")
  private String type;

  @HopMetadataProperty(injectionKeyDescription = "Plc4x.Read.Meta.Injection.Field.Format")
  private String format;

  @HopMetadataProperty(injectionKeyDescription = "Plc4x.Read.Meta.Injection.Field.Length")
  private int length;

  @HopMetadataProperty(injectionKeyDescription = "Plc4x.Read.Meta.Injection.Field.Precision")
  private int precision;

  @HopMetadataProperty(injectionKeyDescription = "Plc4x.Read.Meta.Injection.Field.Currency")
  private String currency;

  @HopMetadataProperty(injectionKeyDescription = "Plc4x.Read.Meta.Injection.Field.Decimal")
  private String decimal;

  @HopMetadataProperty(injectionKeyDescription = "Plc4x.Read.Meta.Injection.Field.Group")
  private String group;

  // Yeah, it has a key of "nullif", keep it for backward compatibility
  @HopMetadataProperty(
      key = "nullif",
      injectionKeyDescription = "Plc4x.Read.Meta.Injection.Field.Value")
  private String value;

  @HopMetadataProperty(
      key = "set_empty_string",
      injectionKeyDescription = "Plc4x.Read.Meta.Injection.Field.SetEmptyString")
  private boolean setEmptyString;

  public Plc4xGeneratorField() {}

  public Plc4xGeneratorField(Plc4xGeneratorField f) {
    this.name = f.name;
    this.item = f.item;    
    this.type = f.type;
    this.format = f.format;
    this.length = f.length;
    this.precision = f.precision;
    this.currency = f.currency;
    this.decimal = f.decimal;
    this.group = f.group;
    this.value = f.value;
    this.setEmptyString = f.setEmptyString;
  }

  public Plc4xGeneratorField(
      String name,
      String item,
      String type,
      String format,
      int length,
      int precision,
      String currency,
      String decimal,
      String group,
      String value,
      boolean setEmptyString) {
    this.name = name;
    this.item = item;
    this.type = type;
    this.format = format;
    this.length = length;
    this.precision = precision;
    this.currency = currency;
    this.decimal = decimal;
    this.group = group;
    this.value = value;
    this.setEmptyString = setEmptyString;
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    Plc4xGeneratorField that = (Plc4xGeneratorField) o;
    return Objects.equals(name, that.name);
  }

  @Override
  public int hashCode() {
    return Objects.hash(name);
  }

  /**
   * Gets name
   *
   * @return value of name
   */
  public String getName() {
    return name;
  }

  /** @param name The name to set */
  public void setName(String name) {
    this.name = name;
  }
  
    /**
   * Gets Item
   *
   * @return value of name
   */
  public String getItem() {
    return item;
  }

  /** @param name The name to set */
  public void setItem(String item) {
    this.item = item;
  }

  /**
   * Gets type
   *
   * @return value of type
   */
  public String getType() {
    return type;
  }

  /** @param type The type to set */
  public void setType(String type) {
    this.type = type;
  }

  /**
   * Gets format
   *
   * @return value of format
   */
  public String getFormat() {
    return format;
  }

  /** @param format The format to set */
  public void setFormat(String format) {
    this.format = format;
  }

  /**
   * Gets length
   *
   * @return value of length
   */
  public int getLength() {
    return length;
  }

  /** @param length The length to set */
  public void setLength(int length) {
    this.length = length;
  }

  /**
   * Gets precision
   *
   * @return value of precision
   */
  public int getPrecision() {
    return precision;
  }

  /** @param precision The precision to set */
  public void setPrecision(int precision) {
    this.precision = precision;
  }

  /**
   * Gets currency
   *
   * @return value of currency
   */
  public String getCurrency() {
    return currency;
  }

  /** @param currency The currency to set */
  public void setCurrency(String currency) {
    this.currency = currency;
  }

  /**
   * Gets decimal
   *
   * @return value of decimal
   */
  public String getDecimal() {
    return decimal;
  }

  /** @param decimal The decimal to set */
  public void setDecimal(String decimal) {
    this.decimal = decimal;
  }

  /**
   * Gets group
   *
   * @return value of group
   */
  public String getGroup() {
    return group;
  }

  /** @param group The group to set */
  public void setGroup(String group) {
    this.group = group;
  }

  /**
   * Gets value
   *
   * @return value of value
   */
  public String getValue() {
    return value;
  }

  /** @param value The value to set */
  public void setValue(String value) {
    this.value = value;
  }

  /**
   * Gets setEmptyString
   *
   * @return value of setEmptyString
   */
  public boolean isSetEmptyString() {
    return setEmptyString;
  }

  /** @param setEmptyString The setEmptyString to set */
  public void setSetEmptyString(boolean setEmptyString) {
    this.setEmptyString = setEmptyString;
  }
}
