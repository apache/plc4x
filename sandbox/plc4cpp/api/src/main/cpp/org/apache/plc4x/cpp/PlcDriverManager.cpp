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

#include "PlcDriverManager.h"
#if defined (_WIN32)
    #include <windows.h>
#elif (__linux__) || (defined (__APPLE__) && defined (__MACH__))
    #include <dlfcn.h>
#endif

#include <boost/filesystem.hpp>
#include <boost/regex.hpp>
#include <iostream>

namespace fs = boost::filesystem;

namespace org
{
	namespace apache
	{
		namespace plc4x
		{
            namespace cpp
            {
                PlcDriverManager::PlcDriverManager()
                {
                    findDrivers();
                }
            
                PlcDriverManager::~PlcDriverManager()
				{
				}

				/**
					* Connects to a PLC using the given plc connection string.
					*
					* @param url plc connection string.
					* @return PlcConnection object.
					* @throws PlcConnectionException an exception if the connection attempt failed.
					*/
				PlcConnection* PlcDriverManager::getConnection(std::string url)
				{
                    PlcDriver* pPlcDriver = NULL;
                    PlcConnection* pPlcConnection = NULL;
                    
                    pPlcDriver = getDriver(url);
                    if (pPlcDriver != NULL)
                    { 
                        pPlcConnection = pPlcDriver->connect(url);
                    }
                    
					return pPlcConnection;
				}

				/**
					* Connects to a PLC using the given plc connection string using given authentication credentials.
					*
					* @param url            plc connection string.
					* @param authentication authentication credentials.
					* @return PlcConnection object.
					* @throws PlcConnectionException an exception if the connection attempt failed.					
					**/
				PlcConnection* PlcDriverManager::getConnection(std::string url, PlcAuthentication& authentication)
				{
                    PlcDriver* pPlcDriver = NULL;
                    PlcConnection* pPlcConnection = NULL;
                    
                    pPlcDriver = getDriver(url);
                    if (pPlcDriver != NULL)
                    {
                        pPlcConnection = pPlcDriver->connect(url, authentication);
                    }
					
					return pPlcConnection;
				}

                PlcDriver* PlcDriverManager::getDriver(std::string url)
                {
                    PlcDriver* pPlcDriver = NULL;
                    std::string strProtocol = "";
                    
                    try
                    {
                        boost::regex exFilter(REGEX_PROTOCOL.c_str());
                        boost::smatch what;

                        if (boost::regex_search(url, what, exFilter))
                        {
                            if (what.size() > 1)
                            {
                                strProtocol = what[1].str();
                                pPlcDriver = _mapDrivers[strProtocol];
                            }
                        }

                        if (pPlcDriver == NULL)
                        {
                            BOOST_THROW_EXCEPTION(PlcConnectionException("Unable to find driver for protocol '" + strProtocol + "'"));
                        }
                    }
                    catch (std::exception& ex)
                    {
                        BOOST_THROW_EXCEPTION(PlcConnectionException("Invalid plc4; connection string '" + url + "'", ex));
                    }

                    return pPlcDriver;
                }

                void PlcDriverManager::findDrivers()
                {
                    typedef PlcDriver* (*pfCreatePlcDriver)();
                    pfCreatePlcDriver CreatePlcDriver;
                        
                    boost::regex exFilter(PLC_DRIVER_TEMPLATE);
                    boost::smatch what;
                        
                    for (fs::recursive_directory_iterator itDirFiles("./"); itDirFiles != fs::recursive_directory_iterator(); itDirFiles++)
                    {
                        if (boost::filesystem::is_regular_file(*itDirFiles))
                        {
                            std::string strFilename = itDirFiles->path().filename().string();

                            if (boost::regex_search(strFilename, what, exFilter))
                            {
                                std::string strDriverName = what[1].str();                                    
                                    
                                try
                                {
#if defined (_WIN32)
                                    HINSTANCE hdll = NULL;
                                    hdll = LoadLibrary((itDirFiles->path().string().c_str()));
                                    if (hdll != NULL)
                                    {
                                        CreatePlcDriver = (pfCreatePlcDriver)GetProcAddress(hdll, PLC_CREATE_DRIVER_INSTANCE.c_str());
#elif defined (__linux__) || (defined (__APPLE__) && defined (__MACH__))
                                    void *hdll = NULL;
                                    hdll = dlopen((itDirFiles->path().string().c_str()),RTLD_NOW);
                                    if (hdll != NULL)
                                    {
                                        CreatePlcDriver = (pfCreatePlcDriver)dlsym(hdll, PLC_CREATE_DRIVER_INSTANCE.c_str());
#endif
                                        if (CreatePlcDriver != NULL)
                                        {
                                            PlcDriver* pPlcDriver = NULL;

                                            pPlcDriver = CreatePlcDriver();
                                            if (pPlcDriver != NULL)
                                            {
                                                _mapDrivers.insert(std::pair<std::string, PlcDriver*>(pPlcDriver->getProtocolCode(), pPlcDriver));
                                            }
                                        }
                                    }
                                }
                                catch (...)
                                {
                                }
                            }
                        }
                    }

                    if (_mapDrivers.size() == 0)
                    {
                        BOOST_THROW_EXCEPTION(PlcConnectionException("Unable to find drivers"));
                    }
                   
                    return;
                }                
			}
		}
	}
}
