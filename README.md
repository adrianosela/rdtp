# rdtp - Reliable Datagram Transfer Protocol

Specification and implementation of a reliable datagram transfer protocol (Transport Layer) to be used on IP networks


```
              0      7 8     15 16    23 24    31
             +--------+--------+--------+--------+
             |     Src. Port   |    Dst. Port    |
             +--------+--------+--------+--------+
             |      Length     |    Checksum     |
             +--------+--------+--------+--------+
             |             ( Data )              |
             +               ....                +
             
       Reliable Datagram Transfer Protocol Header Format
```
