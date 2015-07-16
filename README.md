# 邮件群发

这是一个很好的邮件群发的程序，除了模板替换，我想我还需要加上频率控制，我不想刚发几封邮件就被block，特别是在国内的邮件运营条件下。

增加的功能包括:

* 频率控制
* 可以指定模板中哪些字段的值在一个列表中随机选取

## 使用说明

````
./postman
Postman is a utility for sending batch emails.

Usage:

  postman [flags]

Flags:

  -attach      attach a list of comma separated files
  -c           number of concurrent requests to have
  -csv         path to csv of contact list
  -rand        path to json file that descripte which column should pick item from list randomly
  -debug       print emails to stdout instead of sending
  -fmin        number of requests in x minutes, x value 【频率时间间隔】
  -freq        number of requests in x minutes 【频率时间间隔内请求次数】
  -html        html template path
  -password    smtp password
  -port        port of smtp server
  -sender      email to send from
  -server      url of smtp server
  -subject     subject of email
  -text        text template path
  -user        smtp username
````

-rand 选项对应的 json 格式如下：
```
[
    {
        "Name":"Career",
        "Items":["程序员","工程师","学生","设计师","建筑师","电子工程师","分析师"," DBA"]
    }
]
```
- Name 提定 csv 文件中那个列是要随机替换的
- Items 是可待随机替换内容的列表

可以在 json 是指定多个 csv 的文件的列进行替换。被指定要替换的列在 CSV 文件中的内容会被忽略。

## 安装

````
$ go get github.com/liujianping/postman
$ cd $GOHOME/src/github.com/liujianping 
$ go install
$ $GOHOME/bin/postman
````

# Postman ![Analytics](https://ga-beacon.appspot.com/UA-34529482-6/postman/readme?pixel) [![Hack zachlatta/postman on Nitrous.IO](https://d3o0mnbgv6k92a.cloudfront.net/assets/hack-s-v1-7475db0cf93fe5d1e29420c928ebc614.png)](https://www.nitrous.io/hack_button?source=embed&runtime=go&repo=zachlatta%2Fpostman&file_to_open=main.go)

<img src="http://gh.landersbenjamin.com/everything-sloths/svg/mail.svg" width="130" alt="Postman Icon" align="right">
Postman is a command-line utility for batch-sending email.

#### Features

* Fast, templated, bulk emails
* Reads template attributes from CSV
* Works with any SMTP server

##### Why this over `cat | sed | sendmail < bcc distro_list`?

* Supports both text and HTML parts in emails
* All of the power of templates in Go (conditionals, etc)
* Some SMTP providers will complain if there are too many emails in BCC
  (generally >1000)
* Sends emails concurrently

### Installation

    $ go get github.com/zachlatta/postman

### Usage

    $ postman [flags]

#### Example

```
$ postman -html template.html -text template.txt -csv recipients.csv \
    -sender "Zaphod Beeblebrox <zaphod@beeblebrox.com>" \
    -subject "Hello, World!" -server smtp.beeblebrox.com -port 587 \
    -user zaphod -password Betelgeuse123
```

template.html:

```
<h1>Hello, {{.Name}}! You are a {{.Type}}</h1>
```

template.txt:

```
Hello, {{.Name}}! You are a {{.Type}}.
```

recipients.csv:

```
Email,Name,Type
arthur@dent.com,Arthur Dent,Human
ford@prefect.com,Ford Prefect,Alien
martin@gpp.com,Martin,Robot
trillian@mcmillan.com,Trillian,Human
```

Please check out (and contribute to) [the usage page on the
wiki](https://github.com/zachlatta/postman/wiki/Usage) for more details.

## License

[tl;dr](https://tldrlegal.com/license/mit-license)

The MIT License (MIT)

Copyright (c) 2014 Zach Latta

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
of the Software, and to permit persons to whom the Software is furnished to do
so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
