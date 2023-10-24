/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.merlot.uns.core;


import java.io.File;
import java.util.HashSet;
import java.util.List;
import java.util.Set;

import org.apache.directory.api.ldap.model.entry.Entry;
import org.apache.directory.api.ldap.model.exception.LdapException;
import org.apache.directory.api.ldap.model.name.Dn;
import org.apache.directory.api.ldap.model.schema.SchemaManager;
import org.apache.directory.api.ldap.model.schema.registries.SchemaLoader;
import org.apache.directory.api.ldap.schema.extractor.SchemaLdifExtractor;
import org.apache.directory.api.ldap.schema.extractor.impl.DefaultSchemaLdifExtractor;
import org.apache.directory.api.ldap.schema.loader.LdifSchemaLoader;
import org.apache.directory.api.ldap.schema.manager.impl.DefaultSchemaManager;
import org.apache.directory.api.util.exception.Exceptions;
import org.apache.directory.server.constants.ServerDNConstants;
import org.apache.directory.server.core.DefaultDirectoryService;
//import org.apache.directory.server.core.api.CacheService;
import org.apache.directory.server.core.api.DirectoryService;
import org.apache.directory.server.core.api.DnFactory;
import org.apache.directory.server.core.api.InstanceLayout;
import org.apache.directory.server.core.api.partition.Partition;
import org.apache.directory.server.core.api.schema.SchemaPartition;
import org.apache.directory.server.core.partition.impl.btree.jdbm.JdbmIndex;
import org.apache.directory.server.core.partition.impl.btree.jdbm.JdbmPartition;
import org.apache.directory.server.core.partition.ldif.LdifPartition;
import org.apache.directory.server.i18n.I18n;
import org.apache.directory.server.ldap.LdapServer;
import org.apache.directory.server.protocol.shared.transport.TcpTransport;


/**
 * A simple example exposing how to embed Apache Directory Server from the bleeding trunk
 * into an application.
 *
 * @author <a href="mailto:dev@directory.apache.org">Apache Directory Project</a>
 * @version $Rev$, $Date$
 */
public class EmbeddedADSVerTrunk
{
    /** The directory service */
    private DirectoryService service;

    /** The LDAP server */
    private LdapServer server;


    /**
     * Add a new partition to the server
     *
     * @param partitionId The partition Id
     * @param partitionDn The partition DN
     * @param dnFactory the DN factory
     * @return The newly added partition
     * @throws Exception If the partition can't be added
     */
    private Partition addPartition( String partitionId, String partitionDn, DnFactory dnFactory ) throws Exception
    {
        // Create a new partition with the given partition id 
        JdbmPartition partition = new JdbmPartition(service.getSchemaManager(), dnFactory);
        partition.setId( partitionId );
        partition.setPartitionPath( new File( service.getInstanceLayout().getPartitionsDirectory(), partitionId ).toURI() );
        partition.setSuffixDn( new Dn( partitionDn ) );
        service.addPartition( partition );
        
        return partition;
    }


    /**
     * Add a new set of index on the given attributes
     *
     * @param partition The partition on which we want to add index
     * @param attrs The list of attributes to index
     */
    private void addIndex( Partition partition, String... attrs )
    {
        // Index some attributes on the apache partition
        Set indexedAttributes = new HashSet();

        for ( String attribute : attrs )
        {
            indexedAttributes.add( new JdbmIndex( attribute, false ) );
        }

        ( ( JdbmPartition ) partition ).setIndexedAttributes( indexedAttributes );
    }

    
    /**
     * initialize the schema manager and add the schema partition to diectory service
     *
     * @throws Exception if the schema LDIF files are not found on the classpath
     */
    private void initSchemaPartition() throws Exception
    {
        InstanceLayout instanceLayout = service.getInstanceLayout();
        
        File schemaPartitionDirectory = new File( instanceLayout.getPartitionsDirectory(), "schema" );

        // Extract the schema on disk (a brand new one) and load the registries
        if ( schemaPartitionDirectory.exists() )
        {
            System.out.println( "schema partition already exists, skipping schema extraction" );
        }
        else
        {
            SchemaLdifExtractor extractor = new DefaultSchemaLdifExtractor( instanceLayout.getPartitionsDirectory() );
            extractor.extractOrCopy();
        }

        SchemaLoader loader = new LdifSchemaLoader( schemaPartitionDirectory );
        SchemaManager schemaManager = new DefaultSchemaManager( loader );

        // We have to load the schema now, otherwise we won't be able
        // to initialize the Partitions, as we won't be able to parse
        // and normalize their suffix Dn
        schemaManager.loadAllEnabled();

        List<Throwable> errors = schemaManager.getErrors();

        if ( errors.size() != 0 )
        {
            throw new Exception( I18n.err( I18n.ERR_317, Exceptions.printErrors( errors ) ) );
        }

        service.setSchemaManager( schemaManager );
        
        // Init the LdifPartition with schema
        LdifPartition schemaLdifPartition = new LdifPartition( schemaManager, service.getDnFactory() );
        schemaLdifPartition.setPartitionPath( schemaPartitionDirectory.toURI() );

        // The schema partition
        SchemaPartition schemaPartition = new SchemaPartition( schemaManager );
        schemaPartition.setWrappedPartition( schemaLdifPartition );
        service.setSchemaPartition( schemaPartition );
    }
    
    
    /**
     * Initialize the server. It creates the partition, adds the index, and
     * injects the context entries for the created partitions.
     *
     * @param workDir the directory to be used for storing the data
     * @throws Exception if there were some problems while initializing the system
     */
    private void initDirectoryService( File workDir ) throws Exception
    {
        // Initialize the LDAP service
        service = new DefaultDirectoryService();
        service.setInstanceLayout( new InstanceLayout( workDir ) );
        
//        CacheService cacheService = new CacheService();
//        cacheService.initialize( service.getInstanceLayout() );
//
//        service.setCacheService( cacheService );
        
        // first load the schema
        initSchemaPartition();
        
        // then the system partition
        // this is a MANDATORY partition
        // DO NOT add this via addPartition() method, trunk code complains about duplicate partition
        // while initializing 
        JdbmPartition systemPartition = new JdbmPartition(service.getSchemaManager(), service.getDnFactory());
        systemPartition.setId( "system" );
        systemPartition.setPartitionPath( new File( service.getInstanceLayout().getPartitionsDirectory(), systemPartition.getId() ).toURI() );
        systemPartition.setSuffixDn( new Dn( ServerDNConstants.SYSTEM_DN ) );
        systemPartition.setSchemaManager( service.getSchemaManager() );
        
        // mandatory to call this method to set the system partition
        // Note: this system partition might be removed from trunk
        service.setSystemPartition( systemPartition );
        
        // Disable the ChangeLog system
        service.getChangeLog().setEnabled( false );
        service.setDenormalizeOpAttrsEnabled( true );

        // Now we can create as many partitions as we need
        // Create some new partitions named 'foo', 'bar' and 'apache'.
        Partition fooPartition = addPartition( "foo", "dc=department,dc=example,dc=com", service.getDnFactory() );
        Partition barPartition = addPartition( "bar", "dc=bar,dc=com", service.getDnFactory() );
        Partition apachePartition = addPartition( "apache", "dc=apache,dc=org", service.getDnFactory() );

        // Index some attributes on the apache partition
        addIndex( apachePartition, "objectClass", "ou", "uid" );

        // And start the service
        service.startup();

        // Inject the context entry for dc=foo,dc=com partition if it does not already exist
        try
        {
            service.getAdminSession().lookup( fooPartition.getSuffixDn() );
        }
        catch ( LdapException lnnfe )
        {
            Dn dnFoo = new Dn( "dc=department,dc=example,dc=com" );
            Entry entryFoo = service.newEntry( dnFoo );
            entryFoo.add( "objectClass", "top", "domain", "extensibleObject" );
            entryFoo.add( "dc", "department" );
            service.getAdminSession().add( entryFoo );
        }

        // Inject the context entry for dc=bar,dc=com partition
        try
        {
            service.getAdminSession().lookup( barPartition.getSuffixDn() );
        }
        catch ( LdapException lnnfe )
        {
            Dn dnBar = new Dn( "dc=bar,dc=com" );
            Entry entryBar = service.newEntry( dnBar );
            entryBar.add( "objectClass", "top", "domain", "extensibleObject" );
            entryBar.add( "dc", "bar" );
            service.getAdminSession().add( entryBar );
        }

        // Inject the context entry for dc=Apache,dc=Org partition
        if ( !service.getAdminSession().exists( apachePartition.getSuffixDn() ) )
        {
            Dn dnApache = new Dn( "dc=Apache,dc=Org" );
            Entry entryApache = service.newEntry( dnApache );
            entryApache.add( "objectClass", "top", "domain", "extensibleObject" );
            entryApache.add( "dc", "Apache" );
            service.getAdminSession().add( entryApache );
        }

        // We are all done !
    }


    /**
     * Creates a new instance of EmbeddedADS. It initializes the directory service.
     *
     * @throws Exception If something went wrong
     */
    public EmbeddedADSVerTrunk( File workDir ) throws Exception
    {
        initDirectoryService( workDir );
    }

    
    /**
     * starts the LdapServer
     *
     * @throws Exception
     */
    public void startServer() throws Exception
    {
        server = new LdapServer();
        int serverPort = 10389;
        server.setTransports( new TcpTransport("0.0.0.0", serverPort ) );     
        server.setDirectoryService( service );
        server.start();
    }
    
    public void stopServer() throws Exception
    {      
        server.stop();
        service.shutdown();
    }    
    
    public DirectoryService getDirectoryService() {
        return service;
    }
    
    public LdapServer getServer() {
        return server;
    }    

    
    /**
     * Main class.
     *
     * @param args Not used. 
     */
    public static void main( String[] args ) 
    {
        try
        {
            File workDir = new File( System.getProperty( "java.io.tmpdir" ) + "/server-work" );
            workDir.mkdirs();
            
            // Create the server
            EmbeddedADSVerTrunk ads = new EmbeddedADSVerTrunk( workDir );

            // Read an entry
            Entry result = ads.service.getAdminSession().lookup( new Dn( "dc=apache,dc=org" ) );

            // And print it if available
            System.out.println( "Found entry : " + result );
            
            // optionally we can start a server too
            ads.startServer();
        }
        catch ( Exception e )
        {
            // Ok, we have something wrong going on ...
            e.printStackTrace();
        }
    }
}
