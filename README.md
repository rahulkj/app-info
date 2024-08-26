# AppInfo cli for Cloud Foundry

A command line utility for Cloud Foundry to view the buildpacks, instances, capacity assigned to all the applications those are running in a Cloud Foundry deployment.

## Install
* Download the [latest release](https://github.com/rahulkj/app-info/releases/) for your OS 
* Create a file called `config.yaml` with the following contents
   ```
   token: eyJqzZXIiLCJyb3V0aW5n0NzA5ODAzYzc3ZDBlYjFiNGU0YyIsImVtYWlsIjkVG0oMQVRw2aiksg
   cf_endpoint: https://api.sys.example.com
   ```

   **NOTE:** To fetch the token, please login using the [cf](https://github.com/cloudfoundry/cli/releases) cli, and once authenticated, run the `cf oauth-token` command to fetch the token.
* Run the executable
    `./releases/app-info-darwin-amd64 -config=config.yaml`

## Usage

```
> ./releases/app-info-darwin-amd64 -h

Usage of ./releases/app-info-darwin-amd64:
  -config string
        Absolute path to config file that has the cloud foundry target and bearer token
  -include-env
        Optional flag to include environment variables in json / manifest output. (default false)
  -option string
        csv, json, yaml, packages (default "csv")
```

**SAMPLE OUTPUT**

```
> ./releases/app-info-darwin-amd64 -config=config.yaml 

**** Gathering application metadata from all orgs and spaces ****
Gathering App Data 100% |██████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████| (41/41, 18 it/s)        
**** Following is the csv output ****

ORG,SPACE,APPLICATION,STATE,INSTANCES,MEMORY,DISK,HEALTH_CHECK,STACK,BUILDPACK,DETECTED_BUILDPACK,DETECTED_BUILDPACK_FILENAME
mysql-rep,demo,todo-web-application-mysql,STOPPED,1,1024 MB,1024 MB,port,cflinuxfs4,[https://github.com/cloudfoundry/java-buildpack.git],,[]
dev,dev-group-1,snacks,STARTED,1,1024 MB,1024 MB,port,cflinuxfs4,[java_buildpack_offline],,[java-buildpack-offline-cflinuxfs4-v4.70.0.zip]
p-spring-cloud-services,e172166e-5317-4df2-9384-3480bdf92d85,config-server,STARTED,0,1024 MB,1024 MB,process,cflinuxfs4,[java_buildpack_offline],,[java-buildpack-offline-cflinuxfs4-v4.70.0.zip]
dev,dev-group-2,rabbitmq-demo,STOPPED,1,2048 MB,1024 MB,port,cflinuxfs4,[],,[]
arul,test,cf-hoover,STARTED,1,2048 MB,1024 MB,port,cflinuxfs4,[],,[]
p-spring-cloud-services,90020420-d3ec-44fd-80e9-e4a804e09082,service-registry,STARTED,0,1024 MB,1024 MB,process,cflinuxfs4,[java_buildpack_offline],,[java-buildpack-offline-cflinuxfs4-v4.70.0.zip]
mysql-rep,demo,spring-music,STARTED,1,1024 MB,1024 MB,port,cflinuxfs4,[],,[]
system,system,p-invitations-green,STOPPED,0,0 MB,0 MB,,cflinuxfs4,[nodejs_buildpack],,[nodejs_buildpack-cached-cflinuxfs4-v1.8.26.zip]
ds,sandbox,python-app,STARTED,1,1024 MB,1024 MB,port,cflinuxfs4,[python_buildpack],,[python_buildpack-cached-cflinuxfs4-v1.8.27.zip]
```

```
> ./releases/app-info-darwin-amd64 -config=config.yaml -option=yaml -include-env=true

Gathering pplication metadata from all orgs and spaces
Output will be generated in:  /Users/xxxxx/Documents/output
Gathering App Data 100% |██████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████| (41/41, 12 it/s)        
File 'static-app.yml' created successfully.
File 'p-invitations-blue.yml' created successfully.
File 'nfsbroker.yml' created successfully.
File 'node-app.yml' created successfully.
File 'app-usage-server.yml' created successfully.
File 'service-registry.yml' created successfully.
File 'service-registry.yml' created successfully.
File 'credhub-broker-1.6.4.yml' created successfully.
File 'rabbitmq-demo.yml' created successfully.
File 'python-app.yml' created successfully.
File 'cf-hoover.yml' created successfully.
Generate application manifests are located in:  /Users/xxxxx/Documents/output
```
