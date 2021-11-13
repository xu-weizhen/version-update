# version-update
字节跳动 版本升级 项目

[营员手册](https://bytedance.feishu.cn/docs/doccn2tYZFh28wRCvBmDqREoNie)


# 启动
运行 `./version-update.exe`
浏览器打开 127.0.0.1:8080
+ 获取新版本页面： 127.0.0.1:8080
+ 配置新规则页面： 127.0.0.1:8080/config


## 获取可升级版本
1. 浏览器访问： `http://127.0.0.1:8080`
2. 填写左侧*当前版本*当前版本信息，点击提交
3. 可在右侧*可升级版本*查看可升级版本信息
![获取可升级版本](./img/homepage.png)


## 添加规则
1. 浏览器访问： `http://127.0.0.1:8080/config`
2. 填写*新规则*信息，点击提交
3. 查看添加结果
![添加规则](./img/configpage.png)



# 测试

1. 到 `test/createWhiteList` 目录下运行main，生成16个白名单列表
2. 到 `test/configModels` 目录下运行main，生成16个config内容，随机生成不论Android还是iOS
3. 把 2. 中的16个不同config 依次填入 "127.0.0.1:8080/config" 中提交 
4. 到test/Loop目录下运行随机生成的用例（大部分不满足版本或者平台）所以只显示返回不是null的链接和返回结果
5. 把 4. 中返回的连接改到 test/oneTimes 下单独再测试 查看返回结果是否符合

