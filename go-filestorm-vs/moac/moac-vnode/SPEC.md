## Specification for MOAC implementation ##

This document contains the most up-to-dated specs.

### Subchain implementation ###
[20180302]

To create a new subchain, one would need to do the following
1. Create a SubChainProtocolBase contract. Sharding flag is 0;
2. Create a contract that extend SubChainBase contract. Sharding flag is 0;
3. Activate notifySCS account address by sending 1 sha (more is ok but unnecessary)to that address (0x000000000000000000000000000000000000000d);
4. Call registerOpen function of the SubChainBase contract. This will trigger the notifySCS call. Sharding flag is 0;
5. When SCS got notified, it will initiate a function call RegisterAsSCS in SubChainBase. Sharding flag is 0.;
6. User call SubChainBase with funccode, the funccode to be executed in SCS. Sharding flag is 1.