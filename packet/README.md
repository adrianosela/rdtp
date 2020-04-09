# /packet - rdtp segment header structure

[![Documentation](https://godoc.org/github.com/adrianosela/rdtp/packet?status.svg)](https://godoc.org/github.com/adrianosela/rdtp/packet)
[![license](https://img.shields.io/github/license/adrianosela/rdtp.svg)](https://github.com/adrianosela/rdtp/blob/master/LICENSE)

### Header Format

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