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

#include "DefaultLocalTimeFieldItem.h"

namespace org
{
	namespace apache
	{
		namespace plc4x
		{
			namespace cpp
			{
				namespace base
				{
					namespace messages
					{
						namespace items
						{
							// ==================================================
							bool DefaultLocalTimeFieldItem::isValidTime(int index)
							{
								return true;
							}

							// ===================================================
							time_t DefaultLocalTimeFieldItem::getTime(int index)
							{								
								return getValue(index);
							}

							// ===================================================
							bool DefaultLocalTimeFieldItem::isValidDateTime(int index)
							{
								return isValidTime(index);
							}

							// ===================================================
							time_t  DefaultLocalTimeFieldItem::getDateTime(int index)
							{
								if (!isValidDateTime(index))
								{
									throw new PlcIncompatibleDatatypeException("time_t", index);
								}

								return getValue(index);
							}

							// ===================================================
							bool DefaultLocalTimeFieldItem::isValidDate(int index)
							{
								return true;
							}

							// ===================================================
							tm DefaultLocalTimeFieldItem::getDate(int index)
							{
								if (!isValidDateTime(index))
								{
									throw new PlcIncompatibleDatatypeException("time_t", index);
								}

								time_t time = getValue(index);
								
								#pragma warning(suppress : 4996)
								tm* date = gmtime(&time);

								return *date;
							}

						
						}
						
					}
				}
			}
		}
	}
}
