# scorm-generator-tools
本项目乃Scorm标准课件生成工具, 使用该工具命令行方式如下:
```
scorm-generator-tools -c 配置文件 -o Scorm课件配置文件输出目录 -d 打包文件目录和文件名称(注意加上.zip扩展名)
```
例子:
```
go get -u -v github.com/wangsying/scorm-generator-tools
scorm-generator-tools -c .\config.yaml -o D:\temps\courses -d d:\temps\package.zip
```