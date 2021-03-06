# rdtp - Reliable Data Transport Protocol

[![Go Report Card](https://goreportcard.com/badge/github.com/adrianosela/rdtp)](https://goreportcard.com/report/github.com/adrianosela/rdtp)
[![Documentation](https://godoc.org/github.com/adrianosela/rdtp?status.svg)](https://godoc.org/github.com/adrianosela/rdtp)
[![license](https://img.shields.io/github/license/adrianosela/rdtp.svg)](https://github.com/adrianosela/rdtp/blob/master/LICENSE)

**[IPPROTO_RDTP = 0x9D]**

Specification of a reliable transport layer protocol to be used over IP networks, along a simplistic and modular implementation in Go.

## To-Dos:
* Reliability
  * Polish socket dialer
  * Implement socket listener
  * Implement selective acknowledgements
* Flow Control
  * Receiver window in header

## Based on:
* UDP - User Datagram Protocol [[RFC]](https://tools.ietf.org/html/rfc768)
* TCP - Transmission Control Protocol [[RFC]](https://tools.ietf.org/html/rfc793)

## Header Format

```
 0      7 8     15 16    23 24    31
+--------+--------+--------+--------+
|     Src. Port   |    Dst. Port    |
+--------+--------+--------+--------+
|      Length     |    Checksum     |
+--------+--------+--------+--------+
|          Sequence Number          |
+--------+-----------------+--------+
|       Acknowledgement Number      |
+--------+-----------------+--------+
|  Flags |                          |
+--------+                          |
|             ( Data )              |
+               ....                +
```

## Important Notes: 

The value for the underlying IP header's "Protocol" field must be set to 0x9D (157 -- currently [Unassigned](https://en.wikipedia.org/wiki/List_of_IP_protocol_numbers))

## Over the Wire

Here's a [Wireshark](https://www.wireshark.org/) capture of an RDTP packet over the wire:

![](./.docs/img/cap0.png)

(The highlighted bytes are the RDTP header + payload)
