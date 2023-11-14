package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
)

var ip string
var logger *log.Logger

type Config struct {
	SecretId   string
	SecretKey  string
	Domain     string
	SubDomain  string
	RecordType string
	RecordLine string
	Value      string
	RecordId   uint64
}

func init() {
	file := "./logs/ModifyDnsRecord_" + time.Now().Format("2006-01-02 150405") + ".log"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	// goroutineID := runtime.GoID()
	logger = log.New(logFile, "[ModifyDnsRecord|INFO]", log.LstdFlags|log.Lshortfile|8) // 将文件设置为logger作为输出
}
func main() {
	fmt.Println("111")
	getPublicIP()
	fmt.Println("222")
	os.Exit(0)
	config := GetConfig()
	// 实例化一个认证对象，入参需要传入腾讯云账户 SecretId 和 SecretKey，此处还需注意密钥对的保密
	// 代码泄露可能会导致 SecretId 和 SecretKey 泄露，并威胁账号下所有资源的安全性。以下代码示例仅供参考，建议采用更安全的方式来使用密钥，请参见：https://cloud.tencent.com/document/product/1278/85305
	// 密钥可前往官网控制台 https://console.cloud.tencent.com/cam/capi 进行获取
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

	request.Domain = common.StringPtr(config.Domain)
	request.SubDomain = common.StringPtr(config.SubDomain)
	request.RecordType = common.StringPtr(config.RecordType)
	request.RecordLine = common.StringPtr(config.RecordLine)
	request.Value = common.StringPtr(config.Value)
	request.RecordId = common.Uint64Ptr(config.RecordId)

	// 返回的resp是一个ModifyRecordResponse的实例，与请求对象对应
	response, err := client.ModifyRecord(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return
	}
	if err != nil {
		panic(err)
	}
	// 输出json格式的字符串回包
	fmt.Printf("%s", response.ToJsonString())
}

func getPublicIP() string {
	fmt.Println("start1")
	// 获取本地的地址信息
	addrs, err := net.LookupIP("")
	if err != nil {
		fmt.Println("err:", err)
		return ""
	}
	fmt.Println("start")
	// 遍历地址信息,查找公网 IP 地址
	for _, addr := range addrs {
		fmt.Println(addr.String())
		// if addr.Network == "IPv4" && !strings.Contains(addr.String(), "127.") {
		// 	return addr.String()
		// }
	}

	return ""
}

func GetConfig() Config {
	//创建一个空的结构体,将本地文件读取的信息放入
	c := &Config{}
	//创建一个结构体变量的反射12

	file, err := os.OpenFile("./configs/config.json", os.O_RDONLY, 0777)
	if err != nil {
		logger.Println("打开文件失败")
		logger.Println("打开文件失败")
		//	panic(err)
		logger.Println("err:", err)
	}
	defer file.Close()
	contentByte, err2 := io.ReadAll(file)
	if err2 != nil {
		logger.Println("读取configs失败：", err)
		//panic(err)
		logger.Println("err:", err)
	}
	err = json.Unmarshal(contentByte, c)
	if err != nil {
		logger.Println("configs转化json失败：", err)
	}
	return *c

	//返回Config结构体变量
}
