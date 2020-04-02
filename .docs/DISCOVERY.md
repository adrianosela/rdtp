## Discovery

The puspose of this document it to write down questions which I have encountered throughout the creation of RDTP, as well as my process in finding out how to answer them.

1) I have chosen the protocol number for RDTP to be 157 (0x9D). Will RDTP packets be dropped by NAT gateways en-route? Will the NAT gateway / router see that the protocol number in the IP header is for a transport protocol which it does not understand and then drop the packet? Otherwise, will it try to look for port numbers where they usually are found (first word of the Transport protocol header)?

// TODO 

