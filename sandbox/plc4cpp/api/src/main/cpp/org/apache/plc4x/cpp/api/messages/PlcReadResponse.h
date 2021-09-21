/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

#ifndef _PLC_READ_RESPONSE
#define _PLC_READ_RESPONSE

#include "PlcFieldResponse.h"
#include "PlcReadRequest.h"

#include <boost/multiprecision/cpp_dec_float.hpp>
#include "../types/ValueTypeObject.h"
#include <ctime>

using namespace boost::multiprecision;
using namespace org::apache::plc4x::cpp::api::types;


namespace org
{
	namespace apache
	{
		namespace plc4x
		{
			namespace cpp
			{
				namespace api
				{
					namespace messages
					{
						/**
						 * Response to a {@link PlcReadRequest}.
						 */
						class PlcReadResponse : PlcFieldResponse
						{
							public:	
								virtual PlcRequest* getRequest() = 0;
								virtual int getNumberOfValues(std::string& strName) = 0;
								virtual ValueTypeObject* getObject(std::string& strName) = 0;
								virtual ValueTypeObject* getObject(std::string& strName, int& iIndex) = 0;
								virtual std::vector<ValueTypeObject*>* getAllObjects(std::string& strName) = 0;
								virtual bool isValidBoolean(std::string& strName) = 0;
								virtual bool isValidBoolean(std::string& strName, int& iIndex) = 0;
								virtual bool getBoolean(std::string& strName) = 0;
								virtual bool getBoolean(std::string& strName, int iIndex) = 0;
								virtual std::vector<bool>* getAllBooleans(std::string& strName) = 0;
								virtual bool isValidByte(std::string& strName) = 0;
								virtual bool isValidByte(std::string& strName, int iIndex) = 0;
								virtual char getByte(std::string& strName) = 0;
								virtual char getByte(std::string& strName, int& iIndex) = 0;
								virtual std::vector<char>* getAllBytes(std::string& strName) = 0;
								virtual bool isValidShort(std::string& strName) =  0;
								virtual bool isValidShort(std::string& strName, int iIndex) = 0;
								virtual short getShort(std::string& strName) = 0;
								virtual short getShort(std::string& strName, int& iIndex) = 0;
								virtual std::vector<short>* getAllShorts(std::string& strName) = 0;
								virtual bool isValidInteger(std::string& strName) = 0;
								virtual bool isValidInteger(std::string& strName, int& iIndex) = 0;
								virtual int getInteger(std::string& strName) = 0;
								virtual int getInteger(std::string& strName, int& iIndex) = 0;
								virtual std::vector<int>* getAllIntegers(std::string& strName) = 0;
								virtual bool isValidBigInteger(std::string& strName) = 0;
								virtual bool isValidBigInteger(std::string& strName, int& iIndex) = 0;
								virtual long long* getBigInteger(std::string& strName) = 0;
								virtual long long* getBigInteger(std::string& strName, int& iIndex) = 0;
								virtual std::vector<long long>* getAllBigIntegers(std::string& strName) = 0;
								virtual bool isValidLong(std::string& strName) = 0;
								virtual bool isValidLong(std::string& strName, int& iIndex) = 0;
								virtual long getLong(std::string& strName) = 0;
								virtual long getLong(std::string& strName, int& iIndex) = 0;
								virtual std::vector<long>* getAllLongs(std::string& strName) = 0;
								virtual bool isValidFloat(std::string& strName) = 0;
								virtual bool isValidFloat(std::string& strName, int& iIndex) = 0;
								virtual float getFloat(std::string& strName) = 0;
								virtual float getFloat(std::string& strName, int& iIndex) = 0;
								virtual std::vector<float>* getAllFloats(std::string& strName) = 0;
								virtual bool isValidDouble(std::string& strName) = 0;
								virtual double getDouble(std::string& strName) = 0;
								virtual double getDouble(std::string& strName, int& iIndex) = 0;
								virtual std::vector<double>* getAllDoubles(std::string& strName) = 0;
								virtual bool isValidBigDecimal(std::string& strName) = 0;
								virtual cpp_dec_float_100 getBigDecimal(std::string& strName) = 0;
								virtual cpp_dec_float_100 getBigDecimal(std::string& strName, int& iIndex) = 0;
								virtual std::vector<cpp_dec_float_100>* getAllBigDecimals(std::string& strName) = 0;
								virtual bool isValidString(std::string& strName) = 0;
								virtual bool isValidString(std::string& strName, int& iIndex) = 0;
								virtual std::string getString(std::string& strName) =  0;
								virtual std::string getString(std::string& strName, int& iIndex) = 0;
								virtual std::vector<std::string>* getAllStrings(std::string& strName) = 0;
								virtual bool isValidTime(std::string& strName) = 0;
								virtual bool isValidTime(std::string& strName, int& iIndex) = 0;
								virtual std::time_t getTime(std::string& strName) = 0;
								virtual std::time_t getTime(std::string& strName, int& iIndex) = 0;
								virtual std::vector<std::time_t>* getAllTimes(std::string& strName) = 0;
								virtual bool isValidDate(std::string& strName) = 0;
								virtual bool isValidDate(std::string& strName, int& iIndex) = 0;
								virtual tm& getDate(std::string& strName) = 0;
								virtual tm& getDate(std::string& strName, int& iIndex) = 0;
								virtual std::vector<tm>* getAllDates(std::string& strName)= 0;
								virtual bool isValidDateTime(std::string& strName) = 0;
								virtual bool isValidDateTime(std::string& strName, int& iIndex) = 0;
								virtual tm* getDateTime(std::string& strName) = 0;
								virtual tm*  getDateTime(std::string& strName, int& iIndex) = 0;
								virtual std::vector<tm>* getAllDateTimes(std::string& strName) = 0;
								virtual bool isValidByteArray(std::string& strName) = 0;
								virtual bool isValidByteArray(std::string& strName, int& iIndex) = 0;
								virtual std::vector<char>* getByteArray(std::string& strName) = 0; // Todo: Pointer to byte-array ???
								virtual std::vector<char>* getByteArray(std::string& strName, int& iIndex) = 0;  // Todo: Pointer to byte-array ???
								virtual std::vector<std::vector<char>*>* getAllByteArrays(std::string& strName) = 0;

							private:
						};
					}
				}
			}
		}
	}
}

#endif

