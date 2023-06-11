# golang实现的一个简易的国际象棋

## 接口

接口地址：https://console-docs.apipost.cn/preview/fe618a73dd31e552/c16399a696fb9411



## 亮点

* 用户部分微服务化

使用go-zero作为微服务框架，etcd为发现中心，调了腾讯云的短信SMS服务

用户部分的四个接口都实现了rpc，注册/登陆都是根据手机号自动实现的，快速注册/登陆

- 聊天室做了心跳检测和自动重连

![image-20230609205635108](https://raw.githubusercontent.com/liuxianloveqiqi/Xian-imagehost/main/image/image-20230609205635108.png)

这里进行心跳检测，如果检测到发送消息失败，就调用重连机制。

![image-20230609205757771](https://raw.githubusercontent.com/liuxianloveqiqi/Xian-imagehost/main/image/image-20230609205757771.png)

这里重连的时候要带jwt，不然连不上

* docker部署

将user的api和rpc和chat的api打成dockerfile，上传到docker hub然后服务器拉镜像，docker-compose启动三个服务分别暴露4001，4002，4003端口

![image-20230611212453019](https://raw.githubusercontent.com/liuxianloveqiqi/Xian-imagehost/main/image/image-20230611212453019.png)

* 聊天室

![image-20230611212655476](https://raw.githubusercontent.com/liuxianloveqiqi/Xian-imagehost/main/image/image-20230611212655476.png)

就只能两个人聊天，好像发快了也会挂掉

* 象棋

![image-20230611212756750](https://raw.githubusercontent.com/liuxianloveqiqi/Xian-imagehost/main/image/image-20230611212756750.png)

双方都要输入start后游戏开始，自动分配黑白方，轮流下棋

操作就是a6a7这种，将a6的棋子移动到a7

象棋的地方有bug，de了半天也没de出来，就是用户输入命令后没给返回了，是Move那里出了问题，最后也没解决问题，虽然我感觉我逻辑部分处理没了，但就是没返回也不知道错哪里了

## 总结

有点舍本逐末了，象棋逻辑那部分没写好，最后也没de出bug来，主要是一开始觉得只能web而不能命令行，而且还要写聊天，象棋那里就没写好。😫😫😫😫😫