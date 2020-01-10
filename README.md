# rdtp - Reliable Datagram Transfer Protocol

Specification and implementation of a reliable transport layer protocol to be used over IP networks.

Goal: Eventually be able to perform HTTP communication over this homemade transport protocol

## Header Format

```
 0      7 8     15 16    23 24    31
+--------+--------+--------+--------+
|     Src. Port   |    Dst. Port    |
+--------+--------+--------+--------+
|      Length     |    Checksum     |
+--------+--------+--------+--------+
|             ( Data )              |
+               ....                +
```
