# golang实现的一个简易的国际象棋

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

