# ModifyDNSRecord

## 通过Golang解决电信公网IP地址变化的问题

### 背景
家里申请了上海电信的公网IP，但是最近电信的DNS解析不时的更换解析，经常出现公网IP地址变化，导致配置的域名无法正常访问。

### 解决方案
各大云厂商的DNS解析服务都开放有API，可以利用云厂商的API直接实现域名的DNS解析。

我用的是腾讯云，所以基于腾讯云提供实现方案。

### 实现
1. 获取SecretID和SecretKey
获取SecretID和SecretKey的地址：https://console.cloud.tencent.com/cam/capi

2. 获取解析记录ID
获取解析记录ID的接口地址：https://cloud.tencent.com/document/api/1427/56166

3. 修改DNS记录
修改DNS记录的接口地址：https://cloud.tencent.com/document/api/1427/56157

4. 代码实现
```go
//获取IP的方法
func get_external() string {
	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	content, _ := io.ReadAll(resp.Body)
	//buf := new(bytes.Buffer)
	//buf.ReadFrom(resp.Body)
	//s := buf.String()
	return string(content)
}

```

```go
//更新腾讯云的解析记录
credential := common.NewCredential(
	config.SecretId,
	config.SecretKey,
)
// 实例化一个client选项，可选的，没有特殊需求可以跳过
cpf := profile.NewClientProfile()
cpf.HttpProfile.Endpoint = "dnspod.tencentcloudapi.com"
// 实例化要请求产品的client对象,clientProfile是可选的
client, _ := dnspod.NewClient(credential, "", cpf)

// 实例化一个请求对象,每个接口都会对应一个request对象
request := dnspod.NewModifyRecordRequest()
lastestIp := get_external()
//不一致的情况下进行更新，一致就不更新
if ip != lastestIp {
	logger.Println("IP不一致:", lastestIp)
	ip = lastestIp
	request.Domain = common.StringPtr(config.Domain)
	request.SubDomain = common.StringPtr(config.SubDomain)
	request.RecordType = common.StringPtr(config.RecordType)
	request.RecordLine = common.StringPtr(config.RecordLine)
	request.Value = common.StringPtr(ip)
	request.RecordId = common.Uint64Ptr(config.RecordId)

	// 返回的resp是一个ModifyRecordResponse的实例，与请求对象对应
	response, err := client.ModifyRecord(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Println("An API error has returned:", err)
		continue
	}
	if err != nil {
		logger.Println(err)
	}
	// 输出json格式的字符串回包
	logger.Println(response.ToJsonString())
}

```

5. 打包运行

-  linux环境和windows的打包环境不一样，打包命令也不一样，具体参考官方文档。
-  现提供3个可以直接运行的程序包，分别为linux的arm内核和amd内核已经windows内核，可根据自己环境选择。
-  运行前记得修改config.json文件中的配置。
-  linux环境运行命令(当前目录下)：
```shell

./ModifyDNSRecord

```
-  windows环境运行命令(当前目录下)：
```shell

ModifyDNSRecord.exe

```
或直接双击exe运行。

6. config配置

```json
   {
    "SecretId": "abc",//第一步中获取到的SecretID
    "SecretKey": "ABC",//第一步中获取到SecretKey
    "Domain": "baidu.com",//要解析的域名
    "SubDomain": "www",//主机记录，如 www，如果不传，默认为 @。示例值:www
    "RecordType": "A",//记录类型，通过 API 记录类型获得，大写英文，比如：A 。
    "RecordLine": "默认",//记录线路，通过 API 记录线路获得，中文，比如：默认。
    "RecordId": "1999"//记录 ID 。可以通过第二步接口DescribeRecordList查到所有的解析记录列表以及对应的RecordId
  }

```

### 测试
在腾讯云的DNS记录中随意填写了一个值，过几分钟后刷新，发现已经修改成功，网站也可以正常打开。

### 项目地址
https://github.com/DragonJinG/ModifyDNSRecord