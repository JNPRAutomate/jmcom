# Junos Mass Commander

When operating a large network it is common that you will need to run a massive set of commands across a large set of devices. Examples of this are checking route tables, pulling configurations, or rebooting devices. All of these tasks can easily be accomplished with the Junos Mass Commander or jmcom. Some tasks that can accomplished with Junos Mass Commander:

 - Send commands to multiple hosts at the same time
    - All commands are issued to all connected devices simultaneously
    - It gives you the ability to check the state of multiple devices at the exact same time
 - Reboot devices simultaneously
 - Pull configurations
 - Gather support information

- [Releases](https://github.com/JNPRAutomate/jmcom/releases)
- [Issues](https://github.com/JNPRAutomate/jmcom/issues)



### Limitations

This tool is designed to run the following command types: show, request. Some commands are not supported via NETCONF such as loading licenses or file copies. 

**Since you can use this tool to execute a massive amount of commands extremely quickly please test the command sets you want to run! It is possible to blow up your entire network by rebooting all devices at once. It will only do what you ask it to do.**

If you run into specific commands that seem to hang the session, something that provides strange behavior, or generally stangeness please file an [issue](https://github.com/JNPRAutomate/jmcom/issues). 

Due to the way that jmcom operates it is not designed to do configuration steps. That will be done in a seperate project [Junos Mass Configurer](https://github.com/JNPRAutomate/jmcfg).

# Command line options
jmcom is a command line tool and it can be run from Windows, Mac, or Linux. It requires the use of command line flags to specify three elements: user credentials, commands, and hosts.

- The user credentials consists of a username and either a password or an SSH Key
- Passwords can be specified at the command line or be prompted to enter it
- The same username, password, and ssh host key TODO: are you saying it can use ...
- Hosts can be specified by a hostname or an IP address
- You can log the output of each host to a single file or to a file per devices TODO: not sure on the difference
- Logs are appended to an event they already exists
- Commands are sent to all of the connected hosts at the same time
    - This helps if you are troubleshooting an issue and need to gather the exact state at the exact moment

```
-cmd-file="": File to load commands from
-command="": Commands to run against host: "show version" or for multiple commands "show version","show chassis hardware"
-hosts="": Define hosts to connect to: 1.2.3.3 or 2.3.4.5,1.2.3.4
-hosts-file="": File to load hosts from
-key="": Specify SSH key to use
-log=false: Log output for each host to a separate file
-logdir="": Directory to write logs to. Default is current directory
-password="": Specify password to use with hosts
-prompt=false: Prompts the user to enter a password interactively
-user="": Specify the username to use against hosts
```

## Examples

### Running commands
```
./jmcom -user root -password Juniper -hosts 10.0.1.150 -command "show version"
INFO[2015-05-22T16:57:24-07:00] Connecting to 10.0.1.150
INFO[2015-05-22T16:57:24-07:00] Waiting for connections to establish...
INFO[2015-05-22T16:57:43-07:00] Host: 10.0.1.150 SessionID: 22404 connected  
INFO[2015-05-22T16:57:43-07:00] Issuing commands to hosts...
INFO[2015-05-22T16:57:43-07:00] Host: 10.0.1.150 Sending Command: show version
INFO[2015-05-22T16:57:47-07:00] Host: 10.0.1.150 SessionID: 22404 Command: show version


fpc0:
--------------------------------------------------------------------------
Hostname: EX2200-C
Model: ex2200-c-12p-2g
Junos: 14.1X53-D25.2
JUNOS EX  Software Suite [14.1X53-D25.2]
JUNOS FIPS mode utilities [14.1X53-D25.2]
JUNOS Online Documentation [14.1X53-D25.2]
JUNOS EX 2200 Software Suite [14.1X53-D25.2]
JUNOS Web Management Platform Package [14.1X53-D25.2]
```

### Prompting for a password
```
./jmcom -user root rompt -hosts 10.0.1.150 -command "show version"
Enter password: *******
INFO[2015-05-22T16:58:43-07:00] Connecting to 10.0.1.150
INFO[2015-05-22T16:58:43-07:00] Waiting for connections to establish...
INFO[2015-05-22T16:59:03-07:00] Host: 10.0.1.150 SessionID: 22413 connected  
INFO[2015-05-22T16:59:03-07:00] Issuing commands to hosts...
INFO[2015-05-22T16:59:03-07:00] Host: 10.0.1.150 Sending Command: show version
INFO[2015-05-22T16:59:07-07:00] Host: 10.0.1.150 SessionID: 22413 Command: show version


fpc0:
--------------------------------------------------------------------------
Hostname: EX2200-C
Model: ex2200-c-12p-2g
Junos: 14.1X53-D25.2
JUNOS EX  Software Suite [14.1X53-D25.2]
JUNOS FIPS mode utilities [14.1X53-D25.2]
JUNOS Online Documentation [14.1X53-D25.2]
JUNOS EX 2200 Software Suite [14.1X53-D25.2]
JUNOS Web Management Platform Package [14.1X53-D25.2]


INFO[2015-05-22T16:59:07-07:00] Tasks Complete
```

### Running multiple commands
```
./jmcom -user root -password Juniper -hosts 10.0.1.150 -command "show version,show interfaces terse"
INFO[2015-05-22T16:59:41-07:00] Connecting to 10.0.1.150
INFO[2015-05-22T16:59:41-07:00] Waiting for connections to establish...
INFO[2015-05-22T17:00:00-07:00] Host: 10.0.1.150 SessionID: 22423 connected  
INFO[2015-05-22T17:00:00-07:00] Issuing commands to hosts...
INFO[2015-05-22T17:00:00-07:00] Host: 10.0.1.150 Sending Command: show version
INFO[2015-05-22T17:00:04-07:00] Host: 10.0.1.150 SessionID: 22423 Command: show version


fpc0:
--------------------------------------------------------------------------
Hostname: EX2200-C
Model: ex2200-c-12p-2g
Junos: 14.1X53-D25.2
JUNOS EX  Software Suite [14.1X53-D25.2]
JUNOS FIPS mode utilities [14.1X53-D25.2]
JUNOS Online Documentation [14.1X53-D25.2]
JUNOS EX 2200 Software Suite [14.1X53-D25.2]
JUNOS Web Management Platform Package [14.1X53-D25.2]


INFO[2015-05-22T17:00:04-07:00] Host: 10.0.1.150 Sending Command: show interfaces terse
INFO[2015-05-22T17:00:09-07:00] Host: 10.0.1.150 SessionID: 22423 Command: show interfaces terse


Interface               Admin Link Proto    Local                 Remote
ge-0/0/0                up    up
ge-0/0/0.0              up    up   eth-switch
ge-0/0/1                up    up
ge-0/0/1.0              up    up   eth-switch
ge-0/0/2                up    up
ge-0/0/2.0              up    up   eth-switch
ge-0/0/3                up    down
ge-0/0/3.0              up    down eth-switch
ge-0/0/4                up    up
ge-0/0/4.0              up    up   eth-switch
ge-0/0/5                up    down
ge-0/0/5.0              up    down eth-switch
ge-0/0/6                up    up
ge-0/0/6.0              up    up   eth-switch
ge-0/0/7                up    up
ge-0/0/7.0              up    up   eth-switch
ge-0/0/8                up    down
ge-0/0/8.0              up    down eth-switch
ge-0/0/9                up    down
ge-0/0/9.0              up    down eth-switch
ge-0/0/10               up    down
ge-0/0/10.0             up    down eth-switch
ge-0/0/11               up    down
ge-0/0/11.0             up    down eth-switch
ge-0/1/0                up    down
ge-0/1/0.0              up    down eth-switch
ge-0/1/1                up    down
ge-0/1/1.0              up    down eth-switch
bme0                    up    up
bme0.32768              up    up   inet     128.0.0.1/2
                                            128.0.0.16/2
                                            128.0.0.32/2
                                   tnp      0x10
dsc                     up    up
gre                     up    up
ipip                    up    up
jsrv                    up    up
jsrv.1                  up    up   inet     128.0.0.127/2
lo0                     up    up
lo0.16384               up    up   inet     127.0.0.1           --&gt; 0/0
lsi                     up    up
me0                     up    down
me0.0                   up    down inet
mtun                    up    up
pimd                    up    up
pime                    up    up
tap                     up    up
vlan                    up    up
vlan.0                  up    up   inet     10.0.1.150/24
vme                     up    down


INFO[2015-05-22T17:00:09-07:00] Tasks Complete
```

### CSV Import
Using the flag "-host_list" you can specify a CSV file to load hosts from.

**CSV Format**
```
#host,username,password,keyfile
devops-qfx5100-01-me0.spglab.juniper.net,root,potatoepass,
172.19.100.51,root,!!PROMPT!!,
172.19.100.52,root,!!PROMPT!!,~/.ssh/id_rsa_host
# Special Option !!PROMPT!! prompts user for password
```

```
./jmcom -hosts-file teshosts/host_list.txt -command "show version"
Enter password for devops-qfx5100-01-me0.spglab.juniper.net: *********
Enter password for 172.19.100.51: *********
INFO[2015-05-22T17:06:24-07:00] Connecting to devops-qfx5100-01-me0.spglab.juniper.net
INFO[2015-05-22T17:06:24-07:00] Connecting to 172.19.100.51
INFO[2015-05-22T17:06:24-07:00] Waiting for connections to establish...
INFO[2015-05-22T17:06:25-07:00] Host: 172.19.100.51 SessionID: 63718 connected
INFO[2015-05-22T17:06:42-07:00] Host: devops-qfx5100-01-me0.spglab.juniper.net SessionID: 35767 connected
INFO[2015-05-22T17:06:42-07:00] Issuing commands to hosts...
INFO[2015-05-22T17:06:42-07:00] Host: devops-qfx5100-01-me0.spglab.juniper.net Sending Command: show version
INFO[2015-05-22T17:06:42-07:00] Host: 172.19.100.51 Sending Command: show version
INFO[2015-05-22T17:06:42-07:00] Host: 172.19.100.51 SessionID: 63718 Command: show version


Hostname: srx5800n1
Model: srx5800
JUNOS Software Release [12.3X48-D10.3]
JUNOS AIS Script Suite [1.3R1.0]


INFO[2015-05-22T17:06:44-07:00] Host: devops-qfx5100-01-me0.spglab.juniper.net SessionID: 35767 Command: show version


fpc0:
--------------------------------------------------------------------------
Hostname: devops-qfx5100-01-me0
Model: qfx5100-48t-6q
Junos: 14.1X53-D25.2
JUNOS Base OS Software Suite [14.1X53-D25.2]
JUNOS Base OS boot [14.1X53-D25.2]
JUNOS Crypto Software Suite [14.1X53-D25.2]
JUNOS Online Documentation [14.1X53-D25.2]
JUNOS Kernel Software Suite [14.1X53-D25.2]
JUNOS Packet Forwarding Engine Support (qfx-ex-x86-32) [14.1X53-D25.2]
JUNOS Routing Software Suite [14.1X53-D25.2]
JUNOS Enterprise Software Suite [14.1X53-D25.2]
JUNOS py-base-i386 [14.1X53-D25.2]
JUNOS py-extensions-i386 [14.1X53-D25.2]
Puppet on Junos [2.7.19_1.junos.i386]
Ruby Interpreter [11.10.4_1.junos.i386]
Chef [11.10.4_1.junos.i386]
junos-ez-stdlib [11.10.4_1.junos.i386]
JUNOS Host Software [14.1X53-D25.2]
Junos for Automation Enhancement


INFO[2015-05-22T17:06:44-07:00] Tasks Complete  
```

#### CSV prompt

If you want to be prompted for a specific host's password you can use !!PROMPT!! in the location in the CSV for the password. You will then be prompted for the password.

```csv
10.0.1.150,root,!!PROMPT!!,
```

### Logging

By default the output of the commands are set to the terminal. You can simply redirect this output to a file if you want to collect it.

## Logging to files

Alternatively you can output the logs to a separate file per device. To do this simply use the flag ```-log``` to output the results to a file per device. The file naming will be "hostname.log" or "ip.log" replacing hostname or ip with the hostname or ip of the device.

## Logging files to a specified directory

The logs will be written to the current directory that you are using. You can also specify the output of the logs to go to a specific directory with the flag ```-logdir="/location/to/logs"```.

#Acknowledgements

- Uses [gopass](https://github.com/howeyc/gopass) from Chris Howey/howeyc
- Uses [go-netconf](https://github.com/Juniper/go-netconf) from Juniper Networks
