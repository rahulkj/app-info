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

Following is the csv output

ORG,SPACE,APPLICATION,STATE,BUILDPACK,DETECTED_BUILDPACK
system,system,app-usage-scheduler,STARTED,ruby_buildpack,ruby
system,system,app-usage-worker,STARTED,ruby_buildpack,ruby
system,system,app-usage-server,STARTED,ruby_buildpack,ruby
system,system,search-server-green,STOPPED,nodejs_buildpack,nodejs
system,system,p-invitations-green,STOPPED,nodejs_buildpack,
system,app-metrics-v2,appmetrics,STARTED,java_buildpack_offline,java
john,sandbox,hello-world,STARTED,,nodejs
bob,dev,spring-music,STARTED,,java
bob,dev,spring-music-testing,STARTED,,java
system,system,search-server-blue,STARTED,nodejs_buildpack,nodejs
system,system,p-invitations-blue,STARTED,nodejs_buildpack,nodejs
system,system,apps-manager-js-green,STOPPED,staticfile_buildpack,staticfile
bob,dev,spring-music-testing1,STARTED,,java
bob,dev,testing,STARTED,,java
bob,dev,ra-java-metric,STARTED,,java
system,system,apps-manager-js-blue,STARTED,staticfile_buildpack,staticfile
system,autoscaling,autoscale,STARTED,binary_buildpack,binary
system,autoscaling,autoscale-api,STARTED,java_buildpack_offline,java
bob,dev,staticfile,STARTED,staticfile_buildpack,staticfile
bob,dev,cf-example-staticfile,STARTED,staticfile_buildpack,staticfile
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
