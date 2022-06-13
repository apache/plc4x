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

#include "ProxyDriver.h"
#include "connection/ProxyConnection.h"
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
				namespace proxy
				{
                    string ProxyDriver::getProtocolCode()
                    {
                        return std::string("proxy");
                    }

                    
                    string ProxyDriver::getProtocolName() {
                        return "Proxy";
                    }

                    PlcConnection* ProxyDriver::connect(std::string url)
                    {
                        boost::regex exFilter(PROXY_URI_PATTERN.c_str());
                        boost::smatch what;                       

                        std::string strHost = "";
                        std::string strParams = "";
                        ProxyConnection* pProxyConnection = NULL;

                        try
                        {
                            if (!boost::regex_search(url, what, exFilter))
                            {
                                BOOST_THROW_EXCEPTION(PlcConnectionException("Connection url doesn't match the format 'proxy://{host|ip}'"));
                            }                       

                            strHost = what[1].str(); // Host
                            strParams = what[2].str(); // Params

                            // Resolve from Hostname to implements PlcConnection
                            pProxyConnection = new ProxyConnection(strHost, strParams);
                        }
                        catch (exception ex) 
                        {
                            BOOST_THROW_EXCEPTION(PlcConnectionException("Error connecting to host", ex));
                        }

                        return pProxyConnection;
                    }

                    
                    PlcConnection* ProxyDriver::connect(std::string url, PlcAuthentication authentication)
                    {
                        BOOST_THROW_EXCEPTION(PlcConnectionException("Proxy connections don't support authentication."));
                    }


				}
			}
		}
	}
}
