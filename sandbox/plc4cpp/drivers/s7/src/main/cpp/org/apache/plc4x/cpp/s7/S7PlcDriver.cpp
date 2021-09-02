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

#include "S7PlcDriver.h"
#include "connection/S7PlcConnection.h"
#include <org/apache/plc4x/cpp/api/exceptions/PlcConnectionException.h>
#include <boost/regex.hpp>

using namespace std;
using namespace org::apache::plc4x::cpp::api::exceptions;

namespace org
{
	namespace apache
	{
		namespace plc4x
		{
			namespace cpp
			{
				namespace s7
				{
                    string S7PlcDriver::getProtocolCode() 
                    {
                        return std::string("s7");
                    }

                    
                    string S7PlcDriver::getProtocolName() {
                        return "Siemens S7 (Basic)";
                    }

                    PlcConnection* S7PlcDriver::connect(std::string url)
                    {
                        boost::regex exFilter(S7_URI_PATTERN.c_str());
                        boost::smatch what;                       

                        std::string strHost = "";
                        int iRack = -1;
                        int iSlot = -1; 
                        std::string strParams = "";
                        S7PlcConnection* pS7PlcConnection = NULL;

                        try
                        {
                            if (!boost::regex_search(url, what, exFilter))
                            {
                                BOOST_THROW_EXCEPTION(PlcConnectionException("Connection url doesn't match the format 's7://{host|ip}/{rack}/{slot}'"));
                            }                       

                            strHost = what[1].str(); // Host
                            iRack = std::stoi(what[2].str()); // Rack
                            iSlot = std::stoi(what[3].str()); // Slot
                            strParams = what[4].str(); // Params

                            // Resolve from Hostname to implements PlcConnection
                            pS7PlcConnection = new S7PlcConnection(strHost, iRack, iSlot, strParams);
                        }
                        catch (exception ex) 
                        {
                            BOOST_THROW_EXCEPTION(PlcConnectionException("Error connecting to host", ex));
                        }

                        return pS7PlcConnection;
                    }

                    
                    PlcConnection* S7PlcDriver::connect(std::string url, PlcAuthentication authentication)
                    {
                        BOOST_THROW_EXCEPTION(PlcConnectionException("Basic S7 connections don't support authentication."));
                    }


				}
			}
		}
	}
}
