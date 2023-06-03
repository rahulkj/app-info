# Cloud Foundry App Info Usage CLI Plugin

Cloud Foundry plugin extension to view the buildpacks, instances, capacity assigned to all the applications those are running in a Cloud Foundry deployment.

## Install
* Download the [latest release](https://github.com/rahulkj/app-info/releases/) for your OS 
* Run the command
    `cf install-plugin <Path-To-The-Downloaded-Location>/app-info-<OS>-amd64 -f`

## Usage

```
> cf app-info --help

NAME:
   app-info - Command to view all apps running across all orgs/spaces in the cf deployment

USAGE:
   cf app-info [flags]

OPTIONS:
   --csv or -c            Minimal application details
   --json or -j           All application details in json format
   --manifests or -m      Generate application mainfests in current working directory
```

**SAMPLE OUTPUT**

```
> cf app-info --csv

**** Gathering application metadata from all orgs and spaces ****
**** Following is the csv output ****

ORG,SPACE,APPLICATION,STATE,INSTANCES,MEMORY,DISK,BUILDPACK,DETECTED_BUILDPACK,HEALTH_CHECK
system,system,app-usage-scheduler,STARTED,1,1024 MB,1024 MB,ruby_buildpack,ruby,process
system,system,app-usage-worker,STARTED,1,2048 MB,1024 MB,ruby_buildpack,ruby,process
system,system,app-usage-server,STARTED,2,1024 MB,1024 MB,ruby_buildpack,ruby,http
system,system,search-server-green,STOPPED,2,256 MB,1024 MB,nodejs_buildpack,nodejs,port
system,system,p-invitations-green,STOPPED,2,256 MB,1024 MB,nodejs_buildpack,,port
system,app-metrics-v2,appmetrics,STARTED,1,4096 MB,1024 MB,java_buildpack_offline,java,port
john,sandbox,hello-world,STARTED,5,1024 MB,1024 MB,,nodejs,port
bob,dev,spring-music,STARTED,2,1024 MB,1024 MB,,java,port
bob,dev,spring-music-testing,STARTED,1,1024 MB,1024 MB,,java,port
system,system,search-server-blue,STARTED,2,256 MB,1024 MB,nodejs_buildpack,nodejs,port
system,system,p-invitations-blue,STARTED,2,256 MB,1024 MB,nodejs_buildpack,nodejs,port
system,system,apps-manager-js-green,STOPPED,6,128 MB,1024 MB,staticfile_buildpack,staticfile,port
bob,dev,spring-music-testing1,STARTED,1,1024 MB,1024 MB,,java,port
bob,dev,testing,STARTED,1,1024 MB,1024 MB,,java,port
bob,dev,ra-java-metric,STARTED,1,1024 MB,1024 MB,,java,port
system,system,apps-manager-js-blue,STARTED,6,128 MB,1024 MB,staticfile_buildpack,staticfile,port
system,autoscaling,autoscale,STARTED,3,256 MB,1024 MB,binary_buildpack,binary,port
system,autoscaling,autoscale-api,STARTED,1,1024 MB,1024 MB,java_buildpack_offline,java,port
bob,dev,staticfile,STARTED,1,64 MB,256 MB,staticfile_buildpack,staticfile,port
bob,dev,cf-example-staticfile,STARTED,1,64 MB,256 MB,staticfile_buildpack,staticfile,port
```

```
> cf app-info --manifests

Gathering pplication metadata from all orgs and spaces
Output will be generated in:  /Users/xxxxx/Documents/output
File 'app-usage-scheduler.yml' created successfully.
File 'app-usage-worker.yml' created successfully.
File 'app-usage-server.yml' created successfully.
File 'search-server-green.yml' created successfully.
File 'p-invitations-green.yml' created successfully.
File 'appmetrics.yml' created successfully.
File 'hello-world.yml' created successfully.
File 'spring-music.yml' created successfully.
File 'spring-music-testing.yml' created successfully.
File 'search-server-blue.yml' created successfully.
File 'p-invitations-blue.yml' created successfully.
File 'apps-manager-js-green.yml' created successfully.
File 'spring-music-testing1.yml' created successfully.
File 'testing.yml' created successfully.
File 'ra-java-metric.yml' created successfully.
File 'apps-manager-js-blue.yml' created successfully.
File 'autoscale.yml' created successfully.
File 'autoscale-api.yml' created successfully.
File 'staticfile.yml' created successfully.
File 'cf-example-staticfile.yml' created successfully.
Generate application manifests are located in:  /Users/xxxxx/Documents/output
```

## Uninstall

```
$ cf uninstall-plugin app-info
```
