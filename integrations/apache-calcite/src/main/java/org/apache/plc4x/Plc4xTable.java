package org.apache.plc4x;

import org.apache.calcite.DataContext;
import org.apache.calcite.linq4j.AbstractEnumerable;
import org.apache.calcite.linq4j.Enumerable;
import org.apache.calcite.linq4j.Enumerator;
import org.apache.calcite.rel.type.RelDataType;
import org.apache.calcite.rel.type.RelDataTypeFactory;
import org.apache.calcite.schema.ScannableTable;
import org.apache.calcite.schema.StreamableTable;
import org.apache.calcite.schema.Table;
import org.apache.calcite.schema.impl.AbstractTable;
import org.apache.plc4x.java.scraper.config.JobConfiguration;

import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.concurrent.BlockingQueue;

public class Plc4xTable extends AbstractTable implements StreamableTable, ScannableTable {

    private final BlockingQueue<Object[]> queue;
    private final JobConfiguration conf;
    private Object[] current;

    public Plc4xTable(BlockingQueue<Object[]> queue, JobConfiguration conf) {
        this.queue = queue;
        this.conf = conf;
    }

    @Override
    public RelDataType getRowType(RelDataTypeFactory typeFactory) {
        // Create the table spec
        List<String> names = new ArrayList<>();
        List<RelDataType> types = new ArrayList<>();
        for (Map.Entry<String, String> entry : conf.getFields().entrySet()) {
            names.add(entry.getKey());
            types.add(typeFactory.createJavaType(String.class));
        }
        return typeFactory.createStructType(types, names);
    }

    @Override
    public Enumerable<Object[]> scan(DataContext root) {
        return new AbstractEnumerable<Object[]>() {
            @Override
            public Enumerator<Object[]> enumerator() {
                return new Enumerator<Object[]>() {
                    @Override
                    public Object[] current() {
                        return current;
                    }

                    @Override
                    public boolean moveNext() {
                        try {
                            current = queue.take();
                            return true;
                        } catch (InterruptedException e) {
                            Thread.currentThread().interrupt();
                        }
                        return false;
                    }

                    @Override
                    public void reset() {
                        // Unimplemented
                    }

                    @Override
                    public void close() {
                        // Unimplemented
                    }
                };
            }
        };
    }

    @Override
    public Table stream() {
        return this;
    }
}
