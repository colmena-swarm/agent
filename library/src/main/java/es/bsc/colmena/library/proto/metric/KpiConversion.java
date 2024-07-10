/*
 *  Copyright 2002-2024 Barcelona Supercomputing Center (www.bsc.es)
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */
package es.bsc.colmena.library.proto.metric;

import es.bsc.colmena.library.RoleDefinition;
import es.bsc.colmena.library.metrics.ThresholdType;
import org.apache.commons.lang3.NotImplementedException;

import java.time.temporal.ChronoUnit;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class KpiConversion {

    //processing_time[1h] < 60.0
    public static RoleDefinition.Metric parse(String metric) {
        Pattern thresholdPattern = Pattern.compile("(.*)\\[(.*)\\]\\s*(<|>=)\\s*(.*)\\s*");
        Matcher matcher = thresholdPattern.matcher(metric);
        matcher.find();
        String key = matcher.group(1);
        String time = matcher.group(2);
        String thresholdType = matcher.group(3);
        String thresholdValue = matcher.group(4);

        Pattern timePattern = Pattern.compile("(\\d*)(\\D*)");
        Matcher timeMatcher = timePattern.matcher(time);
        timeMatcher.find();
        String timeValue = timeMatcher.group(1);
        String timeUnit = timeMatcher.group(2);

        return new RoleDefinition.Metric(
                key,
                Double.parseDouble(thresholdValue),
                ThresholdType.parseRepresentation(thresholdType),
                Integer.parseInt(timeValue),
                parseChronoUnit(timeUnit));
    }

    private static ChronoUnit parseChronoUnit(String timeUnit) {
        switch (timeUnit) {
            case "s" -> { return ChronoUnit.SECONDS; }
            case "h" -> { return ChronoUnit.HOURS; }
            default -> throw new NotImplementedException();
        }
    }

    public static String convert(RoleDefinition.Metric metric) {
        StringBuilder sb = new StringBuilder();
        sb.append(metric.getKey());
        sb.append("[" + metric.getAmountOfTime());
        sb.append(convert(metric.getUnit()));
        sb.append("]");
        sb.append(" " + metric.getThresholdType().getRepresentation() + " ");
        sb.append(metric.getThreshold());
        return sb.toString();
    }

    public static String convert(ChronoUnit chronoUnit) {
        switch (chronoUnit) {
            case HOURS -> { return "h"; }
            case SECONDS -> { return "s"; }
            case DAYS -> { return "d"; }
            default -> throw new NotImplementedException("no conversion for Chronounit");
        }
    }



}
