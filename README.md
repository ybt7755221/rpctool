# rpctool

## 概述

rpctool是一个将专为[grpc-server](https://github.com/ybt7755221/grpc-server)提供的自动生成工具

rpctool是基于[m2p](https://github.com/ybt7755221/m2p)改进而来,如果只需要mysql to proto请使用m2p

命令执行：

    rpctool --mysql user:password@tcp\(host:port\)/database\?charset=utf8 --table tableName --out-file ./