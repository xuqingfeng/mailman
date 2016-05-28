## mailman :cop:
> 可自定义 SMTP 服务器的本地:love_letter:客户端

![Travis CI](https://img.shields.io/travis/xuqingfeng/mailman/master.svg)
![Coverage Status](https://img.shields.io/coveralls/xuqingfeng/mailman/master.svg)
![Go Report Card](https://goreportcard.com/badge/github.com/xuqingfeng/mailman)

### 安装

```bash
go get github.com/xuqingfeng/mailman

cd $GOPATH/src/github.com/xuqingfeng/mailman

make run
```

### 贡献代码

1.write code
2.make test pass (`make test`)
3.format code (`make format`)
4.send pull request

#### 默认支持的邮箱类型

- QQ Mail
- Hotmail
- Outlook
- Gmail

关于 Gmail 无法发送邮件问题,查看这里: [https://support.google.com/mail/answer/14257](https://support.google.com/mail/answer/14257)
