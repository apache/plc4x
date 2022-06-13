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

#include "DefaultPlcFieldHandler.h"

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
					namespace connection						
					{
						BaseDefaultFieldItem<bool>* DefaultPlcFieldhandler::encodeBoolean(PlcField field, std::vector<void*> values)
						{
							throw new PlcRuntimeException("Invalid encoder for type "/*Todo:  +std::string((field)*/);
						}

						BaseDefaultFieldItem<char>* DefaultPlcFieldhandler::encodeByte(PlcField field, std::vector<void*>)
						{
							throw new PlcRuntimeException("Invalid encoder for type "/* Todo: + field*/);
						}

						BaseDefaultFieldItem<short>* DefaultPlcFieldhandler::encodeShort(PlcField field, std::vector<void*> values)
						{
							throw new PlcRuntimeException("Invalid encoder for type "/* + field*/);
						}

						BaseDefaultFieldItem<int>* DefaultPlcFieldhandler::encodeInteger(PlcField field, std::vector<void*> values)
						{
							throw new PlcRuntimeException("Invalid encoder for type "/* + field*/);
						}

						BaseDefaultFieldItem<long long>* DefaultPlcFieldhandler::encodeBigInteger(PlcField field, std::vector<void*>)
						{
							throw new PlcRuntimeException("Invalid encoder for type "/* + field*/);
						}

						BaseDefaultFieldItem<long>* DefaultPlcFieldhandler::encodeLong(PlcField field, std::vector<void*>)
						{
							throw new PlcRuntimeException("Invalid encoder for type "/* + field*/);
						}

						BaseDefaultFieldItem<float>* DefaultPlcFieldhandler::encodeFloat(PlcField field, std::vector<void*>)
						{
							throw new PlcRuntimeException("Invalid encoder for type " /*+ field*/);

						}

						BaseDefaultFieldItem<cpp_dec_float_100>* DefaultPlcFieldhandler::encodeBigDecimal(PlcField field, std::vector<void*> values)
						{
							throw new PlcRuntimeException("Invalid encoder for type "/* + field*/);
						}

						BaseDefaultFieldItem<double>* DefaultPlcFieldhandler::encodeDouble(PlcField field, std::vector<void*>)
						{
							throw new PlcRuntimeException("Invalid encoder for type "/* + field*/);
						}

						BaseDefaultFieldItem<std::string>* DefaultPlcFieldhandler::encodeString(PlcField field, std::vector<void*> values)
						{
							throw new PlcRuntimeException("Invalid encoder for type "/* + field*/);
						}

						BaseDefaultFieldItem<time_t>* DefaultPlcFieldhandler::encodeTime(PlcField field, std::vector<void*> values)
						{
							throw new PlcRuntimeException("Invalid encoder for type "/* + field*/);
						}

						BaseDefaultFieldItem<tm>* DefaultPlcFieldhandler::encodeDate(PlcField field, std::vector<void*> values)
						{
							throw new PlcRuntimeException("Invalid encoder for type "/* + field*/);
						}

						BaseDefaultFieldItem<tm>* DefaultPlcFieldhandler::encodeDateTime(PlcField field, std::vector<void*> values)
						{
							throw new PlcRuntimeException("Invalid encoder for type "/* + field*/);
						}

						BaseDefaultFieldItem<char*>* DefaultPlcFieldhandler::encodeByteArray(PlcField field, std::vector<void*> values)
						{
							throw new PlcRuntimeException("Invalid encoder for type "/* + field*/);
						}
					}
				}
			}
		}
	}
}


