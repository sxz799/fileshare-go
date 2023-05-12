# fileshare-go

使用 gin+vue 实现临时文件共享服务程序



## Docker部署说明

```
docker build -t filshare-go .
```
```
docker run -p 4000:4000  -v /your/conf/conf.yaml:/home/conf/conf.yaml -v /your/cert/path:/home/cert -v /your/files:/home/files filshare-go
```

**使用自定义的配置文件可以复制映射conf.yaml文件到容器的`/home/conf/conf.yaml`**

**获取上传的文件可映射目录files(以时间戳命名,对应关系在db.db文件内,需自行进行容器终端备份到files目录后查看)**

**SSL证书 文件名修改为 ***key.key*** 和 ***pem.pem*** 放置在`/your/cert/path/`即可**


