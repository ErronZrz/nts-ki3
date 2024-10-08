## 环境依赖

需要 Go 1.21 及以上版本。

## 部署步骤

1. 在 Unix 系统中，执行 `make build` 命令编译项目。
2. 执行 `./scanner` 命令启动扫描器。
3. 执行 `./keyid` 命令启动仅探测 Key ID 的扫描器。

## 命令行参数

参数说明：
- `--date`：扫描的日期，格式为 `YYYY-MM-DD`。必须指定。
- `--index`：扫描的索引，指定读取的输入文件。若不指定则默认为 `0`。
- `--maxCoroutines`：最大并行执行的任务数量，每个 IP 地址对应一个任务。若不指定则默认为 `100`。

示例：

```shell
./scanner --date 2024-07-01 --index 0 --maxCoroutines 10
```

上述命令表示对 `{user_home}/.nts/2024-07-01_ntske_ip_0.txt` 中的 IP 地址执行探测，至多同时执行 10 个任务。

## 输出结果


扫描器的输出结果保存在 `{user_home}/.nts/` 目录下，文件名格式为 `{date}_ntske_{index}.txt`。文件的每一行均表示对一个 NTS-KE 服务器的有效探测结果。每行包含的字段有：
- IP 地址
- 证书中的通用名称
- NTPv4 主机，`Default` 表示由当前主机提供 NTPv4 服务
- NTPv4 端口
- 支持的 AEAD 算法，多个算法之间用逗号分隔，每个算法后用括号内的数字表示该算法的 Cookie 长度，单位字节
- 状态字符串
- 服务器当前的 key ID，长度为 4 字节，用 8 个十六进制数字表示
- 证书生效时间
- 证书过期时间
- 证书颁发者组织名称
- 证书颁发者通用名称
- 生成本条记录的时间
- 不使用 NTS 的 offset
- 算法为 SIV_CMAC_256 的 offset
- 算法为 SIV_CMAC_256 的 offset，使用真实发包时间作为 T1
- 算法为 SIV_CMAC_384 的 offset
- 算法为 SIV_CMAC_384 的 offset，使用真实发包时间作为 T1
- 算法为 SIV_CMAC_512 的 offset
- 算法为 SIV_CMAC_512 的 offset，使用真实发包时间作为 T1

上面的状态字符串由 4 个 `Y` 或者 `N` 组成，每个 `Y` 分别表示（`N` 相反）：
- 证书的通用名称的解析结果与 IP 地址匹配
- 证书未过期
- 证书非自签名证书
- 该服务器能够响应 NTS 二阶段（认证时间同步）请求

相邻字段间用制表符 `\t` 进行分隔。例如：

```
44.9.16.66	a.lx1duc.ampr.org	Default	123	SIV_CMAC_256(104),SIV_CMAC_384(136),SIV_CMAC_512(168)	YYYY	37EDFC29	2024-06-15 15:50:31	2024-09-13 15:50:30	Let's Encrypt	E5	807.443	792.228	791.963	786.756	786.756	788.358	788.358
```