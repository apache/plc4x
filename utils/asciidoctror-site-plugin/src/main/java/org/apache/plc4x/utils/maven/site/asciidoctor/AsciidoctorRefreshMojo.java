/*
 * Licensed under the Apache License, Version 2.0 (the &quot;License&quot;);
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an &quot;AS IS&quot; BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package org.apache.plc4x.utils.maven.site.asciidoctor;

import java.io.File;
import java.util.ArrayList;
import java.util.Collection;
import java.util.Scanner;
import java.util.concurrent.Callable;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.Future;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicBoolean;

import org.apache.commons.io.FileUtils;
import org.apache.commons.io.filefilter.NameFileFilter;
import org.apache.commons.io.filefilter.RegexFileFilter;
import org.apache.commons.io.filefilter.TrueFileFilter;
import org.apache.commons.io.monitor.FileAlterationListener;
import org.apache.commons.io.monitor.FileAlterationListenerAdaptor;
import org.apache.commons.io.monitor.FileAlterationMonitor;
import org.apache.commons.io.monitor.FileAlterationObserver;
import org.apache.maven.plugin.MojoExecutionException;
import org.apache.maven.plugin.MojoFailureException;
import org.apache.maven.plugins.annotations.Mojo;
import org.apache.maven.plugins.annotations.Parameter;
import org.asciidoctor.Asciidoctor;

@Mojo(name = "auto-refresh")
public class AsciidoctorRefreshMojo extends AsciidoctorMojo {
    public static final String PREFIX = AsciidoctorMaven.PREFIX + "refresher.";
    @Parameter(property = PREFIX + "port", required = false)
    protected int port = 2000;

    @Parameter(property = PREFIX + "interval", required = false)
    protected int interval = 2000; // 2s

    private Future<Asciidoctor> asciidoctor = null;
    private Collection<FileAlterationMonitor> monitors = null;

    private final AtomicBoolean needsUpdate = new AtomicBoolean(false);
    private ScheduledExecutorService updaterScheduler = null;

    @Override
    public void execute() throws MojoExecutionException {
        // this is long because of JRuby startup
        createAsciidoctor();

        startPolling();
        startUpdater();

        try {
            doWork();
        } catch (MojoFailureException e) {
            throw new MojoExecutionException("error", e);
        }

        stopUpdater();
        stopMonitor();
    }

    private void stopUpdater() {
        if (updaterScheduler != null) {
            updaterScheduler.shutdown();
        }
    }

    private void startUpdater() {
        updaterScheduler = Executors.newScheduledThreadPool(1);

        // we prevent refreshing more often than all 200ms and we refresh at least once/s
        // NOTE1: it is intended to avoid too much time space between file polling and re-rendering
        // NOTE2: if nothing to refresh it does nothing so all is fine
        updaterScheduler.scheduleAtFixedRate(new Updater(needsUpdate, this), 0, Math.min(1000, Math.max(200, interval / 2)), TimeUnit.MILLISECONDS);
    }

    protected void doWork() throws MojoFailureException, MojoExecutionException {
        getLog().info("Rendered doc in " + executeAndReturnDuration() + "ms");
        doWait();
    }

    protected void doWait() {
        getLog().info("Type [exit|quit] to exit and [refresh] to force a manual re-rendering.");

        String line;
        final Scanner scanner = new Scanner(System.in);
        while ((line = scanner.nextLine()) != null) {
            line = line.trim();
            if ("exit".equalsIgnoreCase(line) || "quit".equalsIgnoreCase(line)) {
                break;
            }

            if ("refresh".equalsIgnoreCase(line)) {
                doExecute();
            } else {
                getLog().warn("'" + line + "' not understood, available commands are [quit, exit, refresh].");
            }
        }
    }

    private void stopMonitor() throws MojoExecutionException {
        if (monitors != null) {
            for (final FileAlterationMonitor monitor : monitors) {
                try {
                    monitor.stop();
                } catch (Exception e) {
                    throw new MojoExecutionException(e.getMessage(), e);
                }
            }
        }
    }

    protected synchronized void doExecute() {
        ensureOutputExists();

        // delete only content files, resources are synchronized so normally up to date
        for (final File f : FileUtils.listFiles(outputDirectory, new RegexFileFilter(ASCIIDOC_REG_EXP_EXTENSION), TrueFileFilter.INSTANCE)) {
            FileUtils.deleteQuietly(f);
        }

        try {
            getLog().info("Re-rendered doc in " + executeAndReturnDuration() + "ms");
        } catch (final MojoExecutionException e) {
            getLog().error(e);
        } catch (final MojoFailureException e) {
            getLog().error(e);
        }
    }

    protected long executeAndReturnDuration() throws MojoExecutionException, MojoFailureException {
        final long start = System.nanoTime();
        super.execute();
        final long end = System.nanoTime();
        return TimeUnit.NANOSECONDS.toMillis(end - start);
    }

    private void startPolling() throws MojoExecutionException {
        monitors = new ArrayList<FileAlterationMonitor>();

        { // content monitor
            final FileAlterationObserver observer;
            if (sourceDocumentName != null) {
                observer = new FileAlterationObserver(sourceDirectory, new NameFileFilter(sourceDocumentName));
            } else if (sourceDirectory != null) {
                observer = new FileAlterationObserver(sourceDirectory, new RegexFileFilter(ASCIIDOC_REG_EXP_EXTENSION));
            } else {
                monitors = null; // no need to start anything because there is no content
                return;
            }

            final FileAlterationMonitor monitor = new FileAlterationMonitor(interval);
            final FileAlterationListener listener = new FileAlterationListenerAdaptor() {
                @Override
                public void onFileCreate(final File file) {
                    getLog().info("File " + file.getAbsolutePath() + " created.");
                    needsUpdate.set(true);
                }

                @Override
                public void onFileChange(final File file) {
                    getLog().info("File " + file.getAbsolutePath() + " updated.");
                    needsUpdate.set(true);
                }

                @Override
                public void onFileDelete(final File file) {
                    getLog().info("File " + file.getAbsolutePath() + " deleted.");
                    needsUpdate.set(true);
                }
            };

            observer.addListener(listener);
            monitor.addObserver(observer);

            monitors.add(monitor);
        }

        { // resources monitors
            if (synchronizations != null) {
                for (final Synchronization s : synchronizations) {
                    final FileAlterationMonitor monitor = new FileAlterationMonitor(interval);
                    final FileAlterationListener listener = new FileAlterationListenerAdaptor() {
                        @Override
                        public void onFileCreate(final File file) {
                            getLog().info("File " + file.getAbsolutePath() + " created.");
                            synchronize(s);
                            needsUpdate.set(true);
                        }

                        @Override
                        public void onFileChange(final File file) {
                            getLog().info("File " + file.getAbsolutePath() + " updated.");
                            synchronize(s);
                            needsUpdate.set(true);
                        }

                        @Override
                        public void onFileDelete(final File file) {
                            getLog().info("File " + file.getAbsolutePath() + " deleted.");
                            FileUtils.deleteQuietly(file);
                            needsUpdate.set(true);
                        }
                    };

                    final File source = s.getSource();

                    final FileAlterationObserver observer;
                    if (source.isDirectory()) {
                        observer = new FileAlterationObserver(source);
                    } else {
                        observer = new FileAlterationObserver(source.getParentFile(), new NameFileFilter(source.getName()));
                    }

                    observer.addListener(listener);
                    monitor.addObserver(observer);

                    monitors.add(monitor);
                }
            }
        }

        for (final FileAlterationMonitor monitor : monitors) {
            try {
                monitor.start();
            } catch (final Exception e) {
                throw new MojoExecutionException(e.getMessage(), e);
            }
        }
    }

    private void createAsciidoctor() {
        final ExecutorService es = Executors.newSingleThreadExecutor();
        asciidoctor = es.submit(new Callable<Asciidoctor>() {
            @Override
            public Asciidoctor call() {
                return Asciidoctor.Factory.create();
            }
        });
        es.shutdown();
    }

    private static class Updater implements Runnable {
        private final AtomicBoolean run;
        private final AsciidoctorRefreshMojo mojo;

        private Updater(final AtomicBoolean run, final AsciidoctorRefreshMojo mojo) {
            this.run = run;
            this.mojo = mojo;
        }

        @Override
        public void run() {
            if (run.get()) {
                run.set(false);
                mojo.doExecute();
            }
        }
    }
}
