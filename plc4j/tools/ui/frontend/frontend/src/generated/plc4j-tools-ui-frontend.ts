/*-
 * #%L
 * plc4j-tools-ui-frontend
 * %%
 * Copyright (C) 2017 - 2023 The Apache Software Foundation
 * %%
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * 
 *      http://www.apache.org/licenses/LICENSE-2.0
 * 
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 * #L%
 */
/* tslint:disable */
/* eslint-disable */
// Generated using typescript-generator version 3.2.1263 on 2023-12-26 16:21:15.

export interface PlcBrowseItem {
    name: string;
    readable: boolean;
    writable: boolean;
    children: { [index: string]: PlcBrowseItem };
    tag: PlcTag;
    options: { [index: string]: PlcValue };
    subscribable: boolean;
}

export interface PlcBrowseItemArrayInfo {
    size: number;
    lowerBound: number;
    upperBound: number;
}

export interface PlcBrowseRequest extends PlcRequest {
    queryNames: string[];
}

export interface PlcBrowseResponse extends PlcResponse {
    request: PlcBrowseRequest;
    queryNames: string[];
}

export interface PlcDiscoveryItem {
    name: string;
    attributes: { [index: string]: PlcValue };
    options: { [index: string]: string };
    connectionUrl: string;
    transportUrl: string;
    protocolCode: string;
    transportCode: string;
}

export interface PlcDiscoveryRequest extends PlcRequest {
}

export interface PlcDiscoveryResponse extends PlcResponse {
    responseCode: PlcResponseCode;
    values: PlcDiscoveryItem[];
    request: PlcDiscoveryRequest;
}

export interface PlcConnectionMetadata {
}

export interface PlcDriverMetadata {
}

export interface ArrayInfo {
    size: number;
    lowerBound: number;
    upperBound: number;
}

export interface PlcQuery {
    queryString: string;
}

export interface PlcTag {
    addressString: string;
    arrayInfo: ArrayInfo[];
    plcValueType: PlcValueType;
}

export interface PlcValue {
    bigInteger: number;
    boolean: boolean;
    simple: boolean;
    short: number;
    length: number;
    byte: number;
    int: number;
    long: number;
    float: number;
    double: number;
    object: any;
    integer: number;
    null: boolean;
    string: string;
    bigDecimal: number;
    time: Date;
    duration: Duration;
    date: Date;
    keys: string[];
    list: PlcValue[];
    nullable: boolean;
    metaDataNames: string[];
    dateTime: Date;
    plcValueType: PlcValueType;
    raw: any;
    struct: { [index: string]: PlcValue };
}

export interface PlcRequest extends PlcMessage {
}

export interface PlcResponse extends PlcMessage {
    request: PlcRequest;
}

export interface Duration extends TemporalAmount, Serializable {
}

export interface PlcMessage {
}

export interface TemporalAmount {
    units: TemporalUnit[];
}

export interface Serializable {
}

export interface Driver {
    code: string;
    name: string;
    supportsDiscovery: boolean;
    configurationOptions: { [index: string]: ConfigurationOption };
    transports: { [index: string]: Transport };
}

export interface TemporalUnit {
    durationEstimated: boolean;
    duration: Duration;
    timeBased: boolean;
    dateBased: boolean;
}

export interface ConfigurationOption {
    name: string;
    typeName: string;
    required: boolean;
    defaultValue: any;
}

export interface Transport {
    code: string;
    name: string;
    options: { [index: string]: any };
}

export interface HttpClient<O> {

    request<R>(requestConfig: { method: string; url: string; queryParams?: any; data?: any; copyFn?: (data: R) => R; options?: O; }): RestResponse<R>;
}

export class RestApplicationClient<O> {

    constructor(protected httpClient: HttpClient<O>) {
    }

    /**
     * HTTP GET /api/drivers
     * Java method: org.apache.plc4x.java.tools.ui.controller.DriverController.getDriverList
     */
    getDriverList(options?: O): RestResponse<Driver[]> {
        return this.httpClient.request({ method: "GET", url: uriEncoding`api/drivers`, options: options });
    }
}

export type RestResponse<R> = Promise<Axios.GenericAxiosResponse<R>>;

export type PlcResponseCode = "OK" | "NOT_FOUND" | "ACCESS_DENIED" | "INVALID_ADDRESS" | "INVALID_DATATYPE" | "INVALID_DATA" | "INTERNAL_ERROR" | "REMOTE_BUSY" | "REMOTE_ERROR" | "UNSUPPORTED" | "RESPONSE_PENDING";

export type PlcSubscriptionType = "CYCLIC" | "CHANGE_OF_STATE" | "EVENT";

export type PlcValueType = "NULL" | "BOOL" | "BYTE" | "WORD" | "DWORD" | "LWORD" | "USINT" | "UINT" | "UDINT" | "ULINT" | "SINT" | "INT" | "DINT" | "LINT" | "REAL" | "LREAL" | "CHAR" | "WCHAR" | "STRING" | "WSTRING" | "TIME" | "LTIME" | "DATE" | "LDATE" | "TIME_OF_DAY" | "LTIME_OF_DAY" | "DATE_AND_TIME" | "LDATE_AND_TIME" | "Struct" | "List" | "RAW_BYTE_ARRAY";

function uriEncoding(template: TemplateStringsArray, ...substitutions: any[]): string {
    let result = "";
    for (let i = 0; i < substitutions.length; i++) {
        result += template[i];
        result += encodeURIComponent(substitutions[i]);
    }
    result += template[template.length - 1];
    return result;
}


// Added by 'AxiosClientExtension' extension

import axios from "axios";
import * as Axios from "axios";

declare module "axios" {
    export interface GenericAxiosResponse<R> extends Axios.AxiosResponse {
        data: R;
    }
}

class AxiosHttpClient implements HttpClient<Axios.AxiosRequestConfig> {

    constructor(private axios: Axios.AxiosInstance) {
    }

    request<R>(requestConfig: { method: string; url: string; queryParams?: any; data?: any; copyFn?: (data: R) => R; options?: Axios.AxiosRequestConfig; }): RestResponse<R> {
        function assign(target: any, source?: any) {
            if (source != undefined) {
                for (const key in source) {
                    if (source.hasOwnProperty(key)) {
                        target[key] = source[key];
                    }
                }
            }
            return target;
        }

        const config: Axios.AxiosRequestConfig = {};
        config.method = requestConfig.method as typeof config.method;  // `string` in axios 0.16.0, `Method` in axios 0.19.0
        config.url = requestConfig.url;
        config.params = requestConfig.queryParams;
        config.data = requestConfig.data;
        assign(config, requestConfig.options);
        const copyFn = requestConfig.copyFn;

        const axiosResponse = this.axios.request(config);
        return axiosResponse.then(axiosResponse => {
            if (copyFn && axiosResponse.data) {
                (axiosResponse as any).originalData = axiosResponse.data;
                axiosResponse.data = copyFn(axiosResponse.data);
            }
            return axiosResponse;
        });
    }
}

export class AxiosRestApplicationClient extends RestApplicationClient<Axios.AxiosRequestConfig> {

    constructor(baseURL: string, axiosInstance: Axios.AxiosInstance = axios.create()) {
        axiosInstance.defaults.baseURL = baseURL;
        super(new AxiosHttpClient(axiosInstance));
    }
}
