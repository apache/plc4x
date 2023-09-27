/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package org.apache.plc4x.nifi.address;

import org.apache.nifi.components.AllowableValue;
import org.apache.nifi.components.PropertyDescriptor;
import org.apache.nifi.expression.ExpressionLanguageScope;
import org.apache.nifi.processor.ProcessContext;
import org.apache.nifi.processor.util.JsonValidator;
import org.apache.nifi.processor.util.StandardValidators;
import org.apache.plc4x.java.DefaultPlcDriverManager;

public class AddressesAccessUtils {

	private static DefaultPlcDriverManager manager = new DefaultPlcDriverManager();

    public static DefaultPlcDriverManager getManager() {
        return manager;
    }
	
	public static final AllowableValue ADDRESS_PROPERTY = new AllowableValue(
		"property-address",
		"Use Properties as Addresses",
		"Each property will be treated as tag-address pairs after Expression Language is evaluated.");

	public static final AllowableValue ADDRESS_TEXT = new AllowableValue(
		"text-address",
		"Use 'Address Text' Property",
		"Addresses will be obtained from 'Address Text' Property. It's content must be a valid JSON " +
			"after Expression Language is evaluated. ");

	public static final AllowableValue ADDRESS_FILE = new AllowableValue(
		"file-address",
		"Use 'Address File' Property",
		"Addresses will be obtained from the file in 'Address File' Property. It's content must be a valid JSON " +
			"after Expression Language is evaluated. ");

	public static final PropertyDescriptor PLC_ADDRESS_ACCESS_STRATEGY = new PropertyDescriptor.Builder()
		.name("plc4x-address-access-strategy")
		.displayName("Address Access Strategy")
		.description("Strategy used to obtain the PLC addresses")
		.allowableValues(ADDRESS_PROPERTY, ADDRESS_TEXT, ADDRESS_FILE)
		.defaultValue(ADDRESS_PROPERTY.getValue())
		.required(true)
		.build();

	public static final PropertyDescriptor ADDRESS_TEXT_PROPERTY = new PropertyDescriptor.Builder()
		.name("text-address-property")
		.displayName("Address Text")
		.description("Must contain a valid JSON object after Expression Language is evaluated. "
			+ "Each field-value is treated as tag-address.")
		.expressionLanguageSupported(ExpressionLanguageScope.FLOWFILE_ATTRIBUTES)
		.addValidator(new JsonValidator())
		.addValidator(new TextPropertyAccessStrategy.TagValidator(manager))
		.dependsOn(PLC_ADDRESS_ACCESS_STRATEGY, ADDRESS_TEXT)
		.required(true)
		.build();

	public static final PropertyDescriptor ADDRESS_FILE_PROPERTY = new PropertyDescriptor.Builder()
		.name("file-address-property")
		.displayName("Address File")
		.description("Must contain a valid path after Expression Language is evaluated. "
			+ "The file content must be a valid JSON and each field-value is treated as tag-address.")
		.expressionLanguageSupported(ExpressionLanguageScope.FLOWFILE_ATTRIBUTES)
		.addValidator(StandardValidators.FILE_EXISTS_VALIDATOR)
		.addValidator(new FilePropertyAccessStrategy.TagValidator(manager))
		.dependsOn(PLC_ADDRESS_ACCESS_STRATEGY, ADDRESS_FILE)
		.required(true)
		.build();

	public static AddressesAccessStrategy getAccessStrategy(final ProcessContext context) {
		String value = context.getProperty(PLC_ADDRESS_ACCESS_STRATEGY).getValue();
		if (ADDRESS_PROPERTY.getValue().equalsIgnoreCase(value))
			return new DynamicPropertyAccessStrategy();
		else if (ADDRESS_TEXT.getValue().equalsIgnoreCase(value))
			return new TextPropertyAccessStrategy();
		else if (ADDRESS_FILE.getValue().equalsIgnoreCase(value))
			return new FilePropertyAccessStrategy();
		return null;
	}
}
