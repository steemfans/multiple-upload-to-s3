# MutiplePart Upload To S3

## 简介
支持断点续传的 S3 上传工具。

目前市面上没有支持断点续传的工具，从国内向国外传输超大文件，比如 1TB 大小的文件时，
总是会受到各种断流的影响，于是开发此工具。

原理就是在上传开始的时候，创建一个 sqlite 数据库，来记录具体的进程。
如果出现中断，下次开始的时候，从数据库中搜索上次的中断位置并继续。

## 使用方法

### Docker 模式

1. 首先创建环境变量文件，参照 `.env.example`
```
AWS_ACCESS_KEY=
AWS_SECRET_KEY=
ENDPOINT=
MAX_RETRIES=20
BLOCK_SIZE=50
```
前三项是跟 S3 相关的配置，无须多说。
`MAX_RETRIES` 是配置单个 Part 失败重试次数。
`BLOCK_SIZE` 是配置单个 Part 的上传文件大小，单位 MB。
> 这里需要注意 `BLOCK_SIZE` 的配置。由于 AWS 最多支持 10000 个 Part，
> 所以请根据自己要上传的文件大小来调整该参数。
> 另外目前是直接把当前任务涉及的 Part 读取到内存里，所以过大的 `BLOCK_SIZE` 需要足够的内存。

2. 启动
假设待上传文件的目录位置是 `/upload`，文件名 `test.img`，
环境变量文件名为 `/root/.env`，使用 `wasabi` 的 S3，
要上传到 `test_bucket` 桶。
```
docker run -itd --rm \
    --name s3upload \
    -v /upload:/upload \
    --env-file /root/.env \
    -w /upload \
    steemfans/muts \
    /app/muts s3put -src=/upload/test.img -bucketname=test_bucket -objectname=test.img
```
启动后，通过 `docker logs -f --tail 100 s3upload` 查看上传进度。

3. 重启
如果中间上传失败，重新执行第二步的启动命令即可。

### 普通模式

普通模式需要自行编译。编译后的使用方法类似 Docker 模式。

启动命令例如：

```
env $(tr "\\n" " " < /root/.env) /usr/bin/muts s3put -src=/upload/test.img -bucketname=test_bucket -objectname=test.img
```

## 待完善

1. 降低内存使用
2. 增加 goroutine 支持
3. 增加中转功能
