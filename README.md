# M2P

## 概述

M2P是一个将mysql table转化为proto的小工具

命令执行：

    m2p --mysql user:password@tcp\(host:port\)/database\?charset=utf8 --table tableName --out-file ./