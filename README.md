## 四川电信故障分析

限制速度测试以50M上行为列子：

1. IP段固定在某一段171.213.x.x，几乎所有测速网站都达标，网盘上传基本达标单波动较大不均衡。
2. iperf3或者使用本程序连续针对某个IP上传700M-到1G左右以上的数据会马上针对这个IP进行限速到5Mb且长时间不能恢复。
3. iperf3测速100Mps往下走稳定到50Mps，用一段时间后变为1Mbps，与上述第二点吻合。
4. 传输2-3G文件没有传输完成立即断线，经过排查，线路拨号挂起造成了线路重新拨号(IP地址变了）严重影响使用，不得不重头开始上传

所需要工具iperf3，与本程序，二选一也可

### 使用方法

##### 本程序
在build文件夹下载对应的系统版本测试

以LINUX为例服务端运行服务端程序，客户端指定客户端IP

服务端:./upload_server_linux_amd64 默认监听18080端口

客户端:./upload_client_linux_amd64 -ip  114.114.114.114 -port 18080 -n 3

参数说明:

1.  -ip 地址
2.  -port 端口默认18080
3.  -n 1次发送1GB数据，10次发送10GB数据

##### iperf3测试客户端上行

服务器如阿里云运行iperf3 -s

客户端执行iperf3 -c 114.114.114.114 -t 3600

多线程多组合测试效果更好

##### 测试目的

电信是否正对你的IP上传做出了限制影响你正常使用，该投诉工信部投诉工信部。

### 工业和信息化部关于规范电信服务协议有关事项的通知

十一、在电信服务协议有效期间，电信业务经营者不得擅自终止提供服务。未经与用户变更协议，不得擅自撤消任何服务功能或降低服务质量，不得擅自增加收费项目或提高资费标准，不得擅自改变与用户约定的电信业务收费方式。用户因违反相关法律法规的除外。
