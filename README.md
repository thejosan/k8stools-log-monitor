
开发初衷： 

1 开发与测试人员经常要看开发与测试环境日志，给他搭ELK他会说ES语法复制，kibana不好用，还不如监控日志然后复现BUG看日志来得快

2 好了他要这样看日志，运维又要给他分配跳板机账号设置响应权限，教他怎么敲命令看容器日志。（what? 开玩笑，连个xshell都不想安装你认为让他搞这么复杂他会愿意吗？）

本小工具为此场景而生。感谢 walk 作者！

目前文件里面为了测试做了ssh 信息输入框，老哥们后续完全删除这些输入框 写死登陆这些信息

已知问题： 

1 用了跳板机 “断开监控” 按钮会很慢才有响应，不知道什么原因，建议直接点右上角X

2 文本框没有检索功能，后续会补上关键词监控（留意更新）

 
![image](https://github.com/thejosan/k8stools-log-monitor/blob/master/img/k8stools-log-monitor.jpg)
