# gort
gort is a flexible, fast and concurrent port scanner cli and a port scanning library written in go with extensive 
features and both Windows and Linux support.
  
![Console-Demo-Windows](assets/demo_windows.png)

| Table of Content                                                                                                 |
| -------------                                                                                                    |
| 1. [Features](#features)                                                                                         |
| 2. [Building from source](#building-from-source)                                                                 |
| 3. [Prebuild binaries](#prebuild-binaries)                                                                       |
| 4. [How to include gort as a library in your own program?](#how-to-include-gort-as-a-library-in-your-own-program) |
| 5. [Usage](#usage)                                                                                               |
| 6. [Limitations](#limitations)                                                                                   |
| 7. [Dependencies](#dependencies)                                                                                 |

## Features
- Scanning of a single target or a concurrent scan of multiple targets at once, either provided by host address or IP with flexible ways to specify 
  ranges of hosts either trough ranges denoted by "-" (eg. 192.88.99-100.1-100) or CIDR-formatted subnet ranges (e.g 192.88.99.1/24).
- Reverse-hostname-lookup for targets provided by IP.
- Scanning a given number of ports based on a 
  [list](https://docs.google.com/spreadsheets/d/1r_IriqmkTNPSTiUwii_hQ8Gwl2tfTUz8AGIOIL-wMIE/export?format=csv) 
  of most commonly found open ports and/or scanning a custom list of 
  provided ports with well-known port lookup support based on an automatically updated list provided by 
  [IANA](https://www.iana.org/assignments/service-names-port-numbers/service-names-port-numbers.xhtml).
- ICMP-Ping support
- MAC-Address lookup for hosts in the local network either via ARP-cache lookup (**supported on both Windows and Linux**) 
  or ARP-request (**only supported on Linux and with root privileges**).
- MAC based vendor lookup trough an API provided by [macvendors.co](http://macvendors.co/).
- Target location detection (local or public network)
- Target-status detection: Uses the methods listed above to determine if a target is reachable or not.
  This together with the vendor lookup provides a nice and quick overview over the network structure of a given 
  subnet, and the devices that can be found in it.
- Outputting of scan results to file for later reference.
- Options to filter output to only show hosts confirmed as online or to only display open ports.
- Also usable as port scanning library.

## Building from source
1. Clone the repository and navigate into it.
   ```
   > git clone https://github.com/ElCap1tan/gort.git
     ...
   > cd gort
   ```
2. Make sure you turn on go modules by setting the ```GO111MODULE``` environment variable to ```on```.  
     
   On Windows use
   ```
   SET GO111MODULE=on
   ```
   On Linux use
      ```
      export GO111MODULE=on
      ```
3. Next make sure to ```go get``` the dependencies...
4. ... and build it by running 
   ```
   > go build
   ```
5. The finished binary can be found in the ```gort``` folder either as ```gort``` or ```gort.exe```.
6. If you plan to move gort to another filesystem path or onto another device and are not sure if you will have internet 
   access the first time you run gort make sure to distribute the ```data``` folder, and it's content inside the main 
   ```gort``` folder alongside your binary as it contains crucial data that gort needs to run. If you have internet access
   when running gort for the first time you can skip this as gort will download the newest version of the missing files itself.

## Prebuild binaries
Will be added in the near future. For now you'll have to build yourself.

## How to include gort as a library in your own program?
Using gorts port scanning capabilities in your own project is as easy as running 
```
> go get github.com/ElCap1tan/gort
```
inside your go module root, and you're good to go. For example usage see 
[How to use gort as a library in your own code?](#how-to-use-gort-as-a-library-in-your-own-code).
## Usage
### How to use the gort cli application for scans from the commandline
Depending on the OS you either need run ```gort``` or ```gort.exe```.  
Running ```gort``` without any arguments will display a usage help message.

```
> gort [-p ports] [-mc count] [-closed] [-online] [-file] hosts
```
#### Mandatory arguments: 
**hosts**  
are comma separated values that can either be

| Description                 | Example                              |
| --------------------------- |:------------------------------------:|
| A single host               | 192.88.99.1 or example.com           |
| A range of hosts            | 192.88.99.1-50 or 192.88.99-100.1-50 |
| A CIDR formatted host range | 192.88.99.1/24                       |
#### Optional arguments
| Name          | Description           | Example  |
| ------------- |:---------------------------------------------------------------------------------------------------------:| -------------:|
| -p            | ports are comma separated values that either can be a single port or a range of ports                     | 80 or 100-200 |
| -mc [int]     | Sets the number of most common open ports to scan. If omitted defaults to 1000.                           |               |
| -closed       | If this flag is passed ports with closed and unknown/filtered state are also shown in the console output. |               |
| -online       | If this flag is passed only hosts confirmed as online are shown in the console output.                    |               |
| -file         | If this flag is passed the scan result will be saved to a file.                                           |               |
| -elevated     | **Only important for Linux:** If this flag is passed the ICMP echo requests will be send via raw sockets. You might want to try in unprivileged mode first. **Important:** Must be run as a super-user when this flag is used or else ping tests won't work! |               |

#### Examples:
- scan the 1000 most common open ports of example.com  
  ```gort example.com```  
- scan the 500 most common open ports of example.com and 192.88.99.1  
  ```gort -mc 500 example.com,192.88.99.1```  
- scan a custom list of ports for example.com and also show closed or unknown ports in result  
  ```gort -p 80,443,1000-1024 -closed example.com```  
- scan the subnet 192.88.99.0/24 for the 100 most common open ports and and a custom list of ports
  and only show targets confirmed as online in the scan result.  
  ```
  gort -mc 100 -p 10334,12012 -online 192.88.99.0/24  
  ```
  or 
  ``` 
  gort -mc 100 -p 10334,12012 -online 192.88.99.0-255  
  ```

**IMPORTANT**: If you plan to run gort for the first time **without internet access**, make sure to copy the ```data``` 
folder and it's content into the same location as the binary. For more information take a look [here](#building-from-source).
### How to use gort as a library in your own code?
Will be added soon.

## Limitations
Will be added soon.

## Dependencies
This project uses:
- The [color](https://github.com/fatih/color) library by [fatih](https://github.com/fatih) for the colored console output
- [arp](https://github.com/mdlayher/arp) by [mdlayher](https://github.com/mdlayher) for the ARP-request based mac lookups
- [arp](https://github.com/mostlygeek/arp) by [mostlygeek](https://github.com/mostlygeek) for ARP-cache based mac lookups
- [go-ping](https://github.com/sparrc/go-ping) by [sparrc](https://github.com/sparrc) for the ICMP ping requests
- The MAC vendor-lookup api by [macvendors.co](http://macvendors.co/) for MAC-to-vendor resolution
