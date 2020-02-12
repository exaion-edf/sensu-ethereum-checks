### Overview

This plugin provides ethereum node basic checks for monitoring node connection to any ethereum rpc able blockchain system.

### Files
 * bin/sensu-ethereum-check

## Usage example

### Help

**sensu-ethereum-check --help**

```
The Sensu Go Ethereum Check plugin

Usage:
  sensu-ethereum-check [flags]

Flags:
  -B, --crit-max-time-without-block float   Critcal max minutes without block (default 20)
  -P, --crit-peers int                      Critical eth pairs amount
  -h, --help                                help for sensu-ethereum-check
  -x, --max-blocks int                      Max blocks to check for address (default 100)
  -a, --miner-addr string                   Miner address (default "0x00")
  -u, --rpc-url string                      Ethereum RPC URL (default "http://127.0.0.1:8545")
  -b, --warn-max-time-without-block float   Warning max minutes without block (default 10)
  -p, --warn-peers int                      Warning eth pairs amount
```

### Check node connectivity

To valide a node connectivity to other peers, you can connect directly to your node and set a "warn-peers" and "crit-peers" to an amount of minimum peers connected to your node

Example: 
```
sensu-ethereum-check -u "127.0.0.1:8545" -p 10 -P 4
```

This will drop a warning if the node (localhost) isn't connected to more than 10 peers, and a critical error if not more than 4.

### Check miner activity

You are able to check your miners activity over the blockchain itself by directly connecting to the miner itself, or by connecting to another node of the same chain.

Example:
```
sensu-ethereum-check -u "127.0.0.1:8545" -a "0x47182729176392917672901" -b 1.5 -B 3 -x 1000
```

This check will drop a warning if no block have been mined by the miner who has the address "0x47182729176392917672901" within the past 1 minutes and 30 seconds. The process will not check more than the 1000 last mined blocks.
