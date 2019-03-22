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

package org.apache.plc4x.java.scraper.triggeredscraper.triggerhandler;

import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.s7.model.S7Field;
import org.apache.plc4x.java.scraper.exception.ScraperException;
import org.apache.plc4x.java.scraper.triggeredscraper.TriggeredScrapeJobImpl;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * basic configuration for all available triggers and handling of regarding condition
 */
//ToDo: Improve structure to make it more generic --> PLC4X-89
public class TriggerConfiguration{
    private static final Logger logger = LoggerFactory.getLogger(TriggerConfiguration.class);

    private static final String S_7_TRIGGER_VAR = "S7_TRIGGER_VAR";
    private static final String SCHEDULED = "SCHEDULED";

    private static final double TOLERANCE_FLOATING_EQUALITY = 1e-6;

    private static final Pattern TRIGGER_STRATEGY_PATTERN =
        Pattern.compile("\\((?<strategy>[A-Z_0-9]+),(?<scheduledInterval>\\d+)(,(\\((?<triggerVar>\\S+)\\))((?<comp>[!=<>]{1,2}))(\\((?<compVar>[a-z0-9.\\-]+)\\)))?\\)");

    private final TriggerType triggerType;
    private final Long scrapeInterval;
    private final String triggerVariable;
    private final String comparator;
    private Comparators comparatorType;
    private TriggeredScrapeJobImpl triggeredScrapeJobImpl;

    private final Object compareValue;
    private final PlcField plcField;

    /**
     * default constructor when an S7Field should be used for triggering
     * @param triggerType type of trigger from enum
     * @param scrapeInterval scrape interval of triggered variable
     * @param triggerVariable field that is conditional for trigger comparison
     * @param comparator selected comparator
     * @param compareValue selected ref-value that is comapred against
     * @param triggeredScrapeJobImpl the job which is valid for the configuration
     * @throws ScraperException when something goes wrong with configuration
     */
    public TriggerConfiguration(TriggerType triggerType, String scrapeInterval, String triggerVariable, String comparator, String compareValue, TriggeredScrapeJobImpl triggeredScrapeJobImpl) throws ScraperException {
        this.triggerType = triggerType;
        this.triggeredScrapeJobImpl = triggeredScrapeJobImpl;
        this.scrapeInterval = parseScrapeInterval(scrapeInterval);
        this.triggerVariable = triggerVariable;
        this.comparator = comparator;

        if(this.triggerType.equals(TriggerType.S7_TRIGGER_VAR)) {
            //test for valid field-connection string, on exception quit job and return message to user
            try {
                // TODO: PLC4X-106 - Make the Scraper not depend on S7 directly
                this.plcField = S7Field.of(triggerVariable);
            } catch (PlcInvalidFieldException e) {
                logger.debug(e.getMessage(), e);
                String exceptionMessage = String.format("Invalid trigger Field for Job %s: %s", triggeredScrapeJobImpl.getJobName(), triggerVariable);
                throw new ScraperException(exceptionMessage);
            }
            //ToDo add more and other trigger
        }
        else{
            String exceptionMessage = String.format("TriggerType %s is not yet implemented", this.triggerType);
            throw new ScraperException(exceptionMessage);
        }


        this.compareValue = convertCompareValue(compareValue);
        detectComparatorType();
        matchTypeAndComparator();

    }

    /**
     * default constructor when scheduled trigger shall be performed
     * @param triggerType type of trigger from enum
     * @param scrapeInterval scrape interval of data from block
     * @throws ScraperException when something goes wrong with configuration
     */
    public TriggerConfiguration(TriggerType triggerType, String scrapeInterval) throws ScraperException {
        this.triggerType = triggerType;
        this.scrapeInterval = parseScrapeInterval(scrapeInterval);
        this.triggerVariable = null;
        this.comparator = null;
        this.compareValue = null;
        this.plcField = null;
        this.comparatorType = null;
    }

    /**
     * parses String of scrape interval
     * @param scrapeInterval string extracted from RegEx
     * @return converted value
     * @throws ScraperException if parsing could not be performed
     */
    private long parseScrapeInterval(String scrapeInterval) throws ScraperException {
        try {
            return Long.parseLong(scrapeInterval);
        }
        catch (Exception e){
            handleException(e);
            String exceptionMessage = String.format("No valid numeric for scrapeInterval for Job %s: %s",triggeredScrapeJobImpl.getJobName(),scrapeInterval);
            throw new ScraperException(exceptionMessage);
        }
    }

    /**
     * evaluates the trigger dependent of base type and converts acquired respectively ref-value to the needed datatype
     * @param value acquired value
     * @return true when condition is matched, false otherwise
     * @throws ScraperException when something goes wrong
     */
    boolean evaluateTrigger(Object value) throws ScraperException {
        if(validateDataType().equals(Boolean.class)){
            boolean currentValue;
            boolean refValue;
            try{
                currentValue = (boolean) value;
                refValue = (boolean) compareValue;
            }
            catch (Exception e){
                handleException(e);
                return false;
            }
            if(this.comparatorType.equals(Comparators.EQUAL)){
                return currentValue == refValue;
            }
            else {
                return currentValue != refValue;
            }
        }
        if(validateDataType().equals(Double.class)
            || validateDataType().equals(Integer.class)
            || validateDataType().equals(Long.class)) {
            double currentValue;
            double refValue;
            try{
                refValue = (double) compareValue;
                if(value instanceof Short){
                    currentValue = ((Short) value).doubleValue();
                }
                else {
                    if (value instanceof Integer) {
                        currentValue = ((Integer) value).doubleValue();
                    }
                    else {
                        if (value instanceof Long) {
                            currentValue = ((Long) value).doubleValue();
                        }
                        else{
                            if (value instanceof Double) {
                                currentValue = (Double) value;
                            }else {
                                currentValue = (double) value;
                            }
                        }
                    }

                }

                //

            }
            catch (Exception e){
                handleException(e);
                return false;
            }

            switch (this.comparatorType) {
                case EQUAL:
                    return isApproximately(currentValue,refValue, TOLERANCE_FLOATING_EQUALITY);
                case UNEQUAL:
                    return !isApproximately(currentValue,refValue, TOLERANCE_FLOATING_EQUALITY);
                case SMALLER:
                    return currentValue < refValue;
                case SMALLER_EQUAL:
                    return currentValue <= refValue;
                case GREATER:
                    return currentValue > refValue;
                case GREATER_EQUAL:
                    return currentValue >= refValue;
            }

        }
        //should not happen, as fallback return false which always implies that no data is collected
        return false;
    }

    /**
     * convertes parsed comparator from regex to ComparatorType
     * @throws ScraperException when no valid comparator has been used
     */
    private void detectComparatorType() throws ScraperException {
        switch (this.comparator){
            case "==":
                this.comparatorType= Comparators.EQUAL;
                break;
            case "!=":
                this.comparatorType= Comparators.UNEQUAL;
                break;
            case "<=":
                this.comparatorType= Comparators.SMALLER_EQUAL;
                break;
            case "<":
                this.comparatorType= Comparators.SMALLER;
                break;
            case ">=":
                this.comparatorType= Comparators.GREATER_EQUAL;
                break;
            case ">":
                this.comparatorType= Comparators.GREATER;
                break;
            default:
                throw new ScraperException("Invalid comparator detected!");
        }
    }

    /**
     * matches data-type and comparator for a valid combination
     * @throws ScraperException when invalid combination is detected
     */
    private void matchTypeAndComparator() throws ScraperException {
        if(validateDataType().equals(Boolean.class)
            && !(this.comparatorType.equals(Comparators.EQUAL) || this.comparatorType.equals(Comparators.UNEQUAL))){
            String exceptionMessage = String.format("Trigger-Data-Type (%s) and Comparator (%s) do not match",this.plcField.getDefaultJavaType(),this.comparatorType);
            throw new ScraperException(exceptionMessage);
        }
        //all other combinations are valid
    }

    /**
     * defines the used base type for comparison
     * @return the detected base type
     * @throws ScraperException when an unsupported S7-Type is choosen,which is not (yet) implemented for comparison
     * ToDo check how to handle time-variables if needed
     */
    private Class<?> validateDataType() throws ScraperException {
        if(this.plcField!=null){
            Class<?> javaDataType = this.plcField.getDefaultJavaType();
            if(!javaDataType.equals(Boolean.class)
                && !javaDataType.equals(Integer.class)
                && !javaDataType.equals(Long.class)
                && !javaDataType.equals(Double.class)
            ){
                String exceptionMessage = String.format("Unsupported plc-trigger variable %s with converted data-type %s used",this.plcField,this.plcField.getDefaultJavaType());
                throw new ScraperException(exceptionMessage);
            }
            return javaDataType;
        }
        else{
            String exceptionMessage = String.format("Unsupported plc-trigger variable %s with converted data-type %s used",this.plcField,this.plcField.getDefaultJavaType());
            throw new ScraperException(exceptionMessage);
        }

    }

    /**
     * parses the ref-value to a given value, as well as checking if ref-value matches to the given data-type
     * @param compareValue compare-value extracted by regex
     * @return converted object to needed data-type
     * @throws ScraperException when something does not match or parsing fails
     */
    private Object convertCompareValue(String compareValue) throws ScraperException {
        Class<?> javaDataType =validateDataType();
        if(javaDataType.equals(Boolean.class)){
            switch (compareValue){
                case "1":
                case "true":
                    return true;
                case "0":
                case "false":
                    return false;
                default:
                    String exceptionMessage = String.format("No valid compare Value at DataType Boolean for trigger for Job %s: %s",triggeredScrapeJobImpl.getJobName(),compareValue);
                    throw new ScraperException(exceptionMessage);
            }
        }
        if(javaDataType.equals(Double.class)
            || javaDataType.equals(Integer.class)
            || javaDataType.equals(Long.class)){
            try {
                //everything fits to Double for conversion ... so for first step use only double
                //ToDo if different handling dependent on specific datatype is needed then differ
                return Double.parseDouble(compareValue);
            }
            catch (Exception e){
                logger.debug(e.getMessage(), e);
                String exceptionMessage = String.format("No valid compare Value at DataType Numeric for trigger for Job %s: %s",triggeredScrapeJobImpl.getJobName(),compareValue);
                throw new ScraperException(exceptionMessage);
            }
        }
        String exceptionMessage = "Invalid Datatype detected ... should not happen and be catcht earlier - please report";
        throw new ScraperException(exceptionMessage);
    }

    /**
     * creates the TriggerConfiguration for a given ScrapeJob from triggerConfig-String
     * @param jobTriggerStrategy config-string from file
     * @param triggeredScrapeJob job belonging to the config
     * @return created TriggerConfiguration
     * @throws ScraperException when something goes wrong
     */
    public static TriggerConfiguration createConfiguration(String jobTriggerStrategy,TriggeredScrapeJobImpl triggeredScrapeJob) throws ScraperException {
        Matcher matcher = TRIGGER_STRATEGY_PATTERN.matcher(jobTriggerStrategy);

        if(matcher.matches()){
            String strat = matcher.group("strategy");
            String scheduledMs = matcher.group("scheduledInterval");

            logger.debug("Strategy: {}, scheduled ms: {}",strat,scheduledMs);

            String triggerVar = matcher.group("triggerVar");
            String comparatorString = matcher.group("comp");
            String comparatorVariable = matcher.group("compVar");

            switch (strat){
                case S_7_TRIGGER_VAR:
                    if(triggerVar ==null || comparatorString==null || comparatorVariable==null){
                        throw new ScraperException("S7_TRIGGER_VAR trigger strategy needs the trigger-condition - information missing! given configString: "+jobTriggerStrategy);
                    }
                    return new TriggerConfiguration(TriggerType.S7_TRIGGER_VAR,scheduledMs,triggerVar,comparatorString,comparatorVariable,triggeredScrapeJob);
                case SCHEDULED:
                    if(triggerVar !=null || comparatorString!=null || comparatorVariable!=null){
                        throw new ScraperException("SCHEDULED trigger strategy must only be used with scheduled interval - nothing more!  given configString: "+jobTriggerStrategy);
                    }
                    return new TriggerConfiguration(TriggerType.SCHEDULED,scheduledMs);
                default:
                    throw new ScraperException("Unknown Trigger Strategy "+strat);
            }


        }
        throw new ScraperException("Invalid trigger strategy string description: "+jobTriggerStrategy);
    }

    private void handleException(Exception e){
        //push up if needed
        logger.debug("Exception: ", e);
    }

    TriggerType getTriggerType() {
        return triggerType;
    }

    public long getScrapeInterval() {
        return scrapeInterval;
    }

    String getTriggerVariable() {
        return triggerVariable;
    }

    Comparators getComparatorType() {
        return comparatorType;
    }

    Object getCompareValue() {
        return compareValue;
    }

    /**
     * check for approximate equality to avoid "Floating-point expressions shall not be tested for equality or inequality." Sonar-Bug
     * @param self current value
     * @param other reference value
     * @param within tolerance band
     * @return if approximate equal, false otherwise
     */
    private static boolean isApproximately(double self, double other, double within)
    {
        return Math.abs(self - other) <= within;
    }

    public enum  Comparators{
        EQUAL,
        UNEQUAL,
        GREATER,
        GREATER_EQUAL,
        SMALLER,
        SMALLER_EQUAL
    }

    //ToDo replace constant TriggerType by more generic ones --> PLC4X-89
    public enum TriggerType {
        SCHEDULED,
        S7_TRIGGER_VAR
    }
}