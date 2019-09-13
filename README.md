# Rancher Wins

`wins` is a way to operate the Windows host inside the Windows container.

## How to use 

### Modules
```
> wins.exe -h
NAME:
   rancher-wins - A way to operate the Windows host inside the Windows container

USAGE:
   wins.exe [global options] command [command options] [arguments...]

VERSION:
   ...

DESCRIPTION:
   Rancher Wins Component (...)

COMMANDS:
     srv, server
     cli, client
     help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

#### Server (run on Windows host)
```
> wins.exe srv -h
NAME:
   rancher-wins srv - The server side commands of Rancher Wins

USAGE:
   rancher-wins srv command [command options] [arguments...]

COMMANDS:
     app, application  Manage Rancher Wins Application

OPTIONS:
   --help, -h  show help
```

#### Client (run inside Windows container)
```
> wins.exe cli -h
NAME:
   rancher-wins cli - The client side commands of Rancher Wins

USAGE:
   rancher-wins cli command [command options] [arguments...]

COMMANDS:
     hns               Manage Host Networking Service
     hst, host         Manage Host
     net, network      Manage Network Adapter
     prc, process      Manage Processes
     route             Manage Routes
     app, application  Manage Rancher Wins Application

OPTIONS:
   --help, -h  show help
```

### Examples

``` powershell
# [host] start the wins server
> wins.exe srv app run --log-level debug --listen rancher_wins

# [host] verify the created npipe
> Get-ChildItem //./pipe/ | Where-Object Name -eq "rancher_wins"
```

#### Query the host network adapter
``` powershell
# [host] start a container
> $WINS_BIN_PATH=<...>; docker run --rm -it -v //./pipe/rancher_wins://./pipe/rancher_wins -v "$($WINS_BIN_PATH):c:\host\wins" -w c:\host\wins --entrypoint powershell mcr.microsoft.com/windows/servercore:ltsc2019

# [inside container] query the host network adapter
>> .\wins.exe cli network get
{"InterfaceIndex":"7","GatewayAddress":"10.170.0.1","SubnetCIDR":"10.170.0.0/20","HostName":"frank-wins-dev","AddressCIDR":"10.170.15.229/32"}
```

#### Start a process on the host
``` powershell
# [host] download nginx 
> Invoke-WebRequest -UseBasicParsing -OutFile nginx.zip -Uri http://nginx.org/download/nginx-1.15.9.zip

# [host] expand nginx in the current directory
> Expand-Archive -Force -Path nginx.zip -Destination c:\nginx

# [host] start a container
> $WINS_BIN_PATH=<...>; echo "`$NGINX_BIND_DIR=$NGINX_BIND_DIR"; docker run --rm -it -v //./pipe/rancher_wins://./pipe/rancher_wins -v "$($WINS_BIN_PATH):c:\host\wins" -v "c:\nginx:c:\nginx" -w c:\host\wins --entrypoint powershell mcr.microsoft.com/windows/servercore:ltsc2019

# [inside container] start nginx and receive the running output
>> .\wins.exe cli app run --path c:\nginx\nginx-1.15.9\nginx.exe --exposes TCP:80

# [host] verify the process
> Get-Process rancher-wins-*
> curl.exe 127.0.0.1
```

## Build
``` powershell
> make.bat build
```

## Testing
There are not any Docker-in-Docker supported Windows images for now, `rancher/wins` has to separate the validation test and integration test.

For validation test, which could be embedded into a containerized CI flow, please run the below command in `PowerShell`:

``` powershell
> make.bat
```

For integration test, please run the below command in `PowerShell`:

``` powershell
> make.bat integration
```

> Note: Don't use `bin/wins.exe` after integration testing. Please `make.bat build` again.

If want both of them, please run the below command in `PowerShell`:

``` powershell
> make.bat all
```

## License

Copyright (c) 2014-2019 [Rancher Labs, Inc.](http://rancher.com)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.