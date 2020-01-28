# Cloud Foundry App Info Usage CLI Plugin

Cloud Foundry plugin extension to view the buildpacks, instances, capacity assigned to all the applications those are running in a Cloud Foundry deployment.

## Install

```
$ go get github.com/rahulkj/app-info
$ go build github.com/rahulkj/app-info
$ cf install-plugin app-info
```

## Usage

**SAMPLE OUTPUT**

```
$ cf app-info

Following is the csv output

ORG,SPACE,APPLICATION,STATE,BUILDPACK
system,system,apps-manager-js,STARTED,staticfile_buildpack
system,system,app-usage-server,STARTED,ruby_buildpack
system,system,console,STARTED,ruby_buildpack
system,system,app-usage-scheduler,STARTED,ruby_buildpack
system,system,app-usage-worker,STARTED,ruby_buildpack
system,notifications-with-ui,notifications-ui,STARTED,Go
system,autoscaling,autoscale,STARTED,Go
pivotal,dev,swagger-ui,STARTED,staticfile_buildpack
pivotal,dev,autoscaler-test-app,STARTED,java-buildpack=v3.6-offline-https://github.com/cloudfoundry/java-buildpack.git#5194155 java-main open-jdk-like-jre=1.8.0_71 open-jdk-like-memory-calculator=2.0.1_RELEASE spring-auto-reconfiguration=1.10.0_RELEASE
pivotal,dev,log-generator,STARTED,java-buildpack=v3.6-offline-https://github.com/cloudfoundry/java-buildpack.git#5194155 java-main open-jdk-like-jre=1.8.0_71 open-jdk-like-memory-calculator=2.0.1_RELEASE spring-auto-reconfiguration=1.10.0_RELEASE
pivotal,dev,test,STARTED,java-buildpack=v3.6-offline-https://github.com/cloudfoundry/java-buildpack.git#5194155 java-main open-jdk-like-jre=1.8.0_71 open-jdk-like-memory-calculator=2.0.1_RELEASE spring-auto-reconfiguration=1.10.0_RELEASE
```

```
$ cf app-info --verbose

Following is the csv output

ORG,SPACE,APPLICATION,STATE,INSTANCES,MEMORY,DISK
system,system,apps-manager-js,STARTED,6,64 MB,1024 MB
system,system,app-usage-server,STARTED,1,128 MB,1024 MB
system,system,console,STARTED,6,1024 MB,1024 MB
system,system,app-usage-scheduler,STARTED,1,128 MB,1024 MB
system,system,app-usage-worker,STARTED,1,1024 MB,1024 MB
system,notifications-with-ui,notifications-ui,STARTED,1,64 MB,1024 MB
system,autoscaling,autoscale,STARTED,1,256 MB,1024 MB
pivotal,dev,swagger-ui,STARTED,1,1024 MB,1024 MB
pivotal,dev,autoscaler-test-app,STARTED,1,1024 MB,1024 MB
pivotal,dev,log-generator,STARTED,1,1024 MB,1024 MB
pivotal,dev,test,STARTED,1,1024 MB,1024 MB

```

## Uninstall

```
$ cf uninstall-plugin app-info
```
