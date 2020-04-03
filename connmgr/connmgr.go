package connmgr

/* This class is in charge of managing a connection:

userLvl: calls read() and write() for messages
connMgr: manages buffers for userLvl read/write
	- user read buff: connMgr receives packets
	from atc and compies payloads to user read buffer
	- user write buff: connMgr reads messages from it and
	sends them to a packetizer whose receiver is the atc,
	which is in charge of eventually transmitting the packet
atc: the air traffic controller (atc), manages retransmissions
     and acking packets
netwk: the interface between program level objects (packets), and
	the host's IP interface

outbound: userLvl -msg-> connMgr(using packetizer) -pck-> atc -pck-> netwk
inbound: nwtwk -pck-> atc -pck-> connMgr -msg-> userLvl

*/
