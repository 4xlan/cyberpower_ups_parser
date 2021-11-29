# Output parser for CyberPower's UPS

This simple utility parses output of `pwrstat -status` and makes small http page for further parsing.

Maybe in someday it will be able to check status and send alerts/execute commands, based on battery stats.

## Pre-reqs

You should have:
- Installed official utility
- User account with an ability to execute `pwrstat` (please, don't use root)
- UPS, which is supported with utility and plugged with USB as well, ofc.

## Configuring

Just set up port (and IP, optionally) in config in place it with name `config.yaml` in `config` directory near bin file. 

