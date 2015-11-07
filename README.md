# Rancher Leader Election

This is very simple approach to leader election that uses Rancher Meta-data.  Every
container in a service is given a create_index that increments with each container
that is launched in the service.  The container with the lowest create_index will be
chosen as the master.

## Using

This binary is expected to be ran as an `ENTRYPOINT` to your container.  For example

```Dockerfile
FROM service
ADD https://github.com/rancher/leader/releases/download/v0.1.0/leader /usr/bin/
RUN chmod +x /usr/bin/leader
ENTRYPOINT ["leader"]
CMD ["service", "args"]
```

On startup the `leader` binary will check metadata and see if this container is the leader.
If it is the leader the arguments to `leader` will be executed.  If the container is not the
leader it will just wait until it becomes the leader (ie has the lowest create_index).

## Port proxying

You can pass `--proxy-tcp-port 8080` to the `leader` program and that will make all non-leader
containers just forward traffic to port `8080` of the leader.

## Embed in script

Additionally you can just run `leader -check` to check if this container is the leader or not.
If it is the leader it will exit with status code 0, otherwise 1.

## Usage

```
NAME:
   ./dist/artifacts/leader - Simple leader election with Rancher

USAGE:
   ./dist/artifacts/leader [global options] command [command options] [arguments...]
   
COMMANDS:
   help, h	Shows a list of commands or help for one command
   
GLOBAL OPTIONS:
   --proxy-tcp-port "0"		Port to proxy to the leader
   --check			Check if we are the leader and exit
   --help, -h			show help
   --generate-bash-completion	
   --version, -v		print the version
```

# License
Copyright (c) 2014-2015 [Rancher Labs, Inc.](http://rancher.com)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
