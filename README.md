# scorm-generator-tools
本项目乃Scorm标准课件生成工具, 暂时只支持HTML课件, 因html本身可以包容多种媒体类型,如:html\ppt\excel\word\pdf\MP4\mp3, 所以应该足够了, 后续如果遇到确实需要单独指定媒体类型再完善.

### 使用该工具命令行方式如下:
```
scorm-generator-tools -c 配置文件 -o Scorm课件配置文件输出目录 -d 打包文件目录和文件名称(注意加上.zip扩展名)
```
例子:
```
go get -u -v github.com/wangsying/scorm-generator-tools
scorm-generator-tools -c .\config.yaml -o D:\temps\courses -d d:\temps\package.zip
```

### 配置文件的例子

配置文件的格式采用的yaml标准规范, 因Scorm默认采用XML, 对于普通的课件制作者来说过于复杂, 而借助一些专用的工具又需要付费或者学习如何使用, 所以通过简单地配置将html课件直接打包成Scorm或许更容易一些.

```
courses:
  - title: 课程名称
    organization: 组织名称(建议用英文)
    coursewares:
      - title: "章节名称"
        coursewares:
          - title: "子课程名称(也即课件)"
            href: "课件所在目录"
          - title: "子课程名称2(也即课件)"
            href: "课件2所在目录"
      - title: "eMail Course"
        coursewares:
          - title: "step 1"
            coursewares:
              - title: "Insider Threat Overview"
                href: "D:\\temps\\test2\\eml2"
              - title: "Email Security on Mobile Devices"
                href: "D:\\temps\\test2\\eml6"
          - title: "step 2"
            coursewares:
              - title: "Spear Phishing Threats"
                href: "D:\\temps\\test2\\eml7"
              - title: "Avoiding Dangerous Attachments"
                href: "D:\\temps\\test2\\eml4"
              - title: "Avoiding Dangerous Links"
                href: "D:\\temps\\test2\\eml3"
          - title: "step 3"
            coursewares:
              - title: "Data Entry Phishing"
                href: "D:\\temps\\test2\\eml5"
              - title: "Introduction to Phishing"
                href: "D:\\temps\\test2\\eml1"
```                
