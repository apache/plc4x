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

#ifndef _SYSTEMCONFIGURATION_H_
#define _SYSTEMCONFIGURATION_H_

#include <boost/property_tree/ptree.hpp>
#include <boost/property_tree/xml_parser.hpp>
#include <boost/foreach.hpp>
#include <string>
#include <set>
#include <exception>
#include <iostream>

// TODO: Implementation SystemConfiguration

namespace org
{
	namespace apache
	{
		namespace plc4x
		{
			namespace cpp
			{
				namespace utils
				{
                    namespace pt = boost::property_tree;

					class SystemConfiguration : public pt::ptree
					{

					public:

                        SystemConfiguration();

						~SystemConfiguration();

						void Load();

                        long getLong(std::string strPath, long lDefaultValue);

					private:

                        const std::string CONFIGFILE = "./ Config/Config.xml";

                        pt::ptree initpTree;

					};
				}
			}
		}
	}
}
#endif