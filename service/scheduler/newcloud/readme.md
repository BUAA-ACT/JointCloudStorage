# Scheduler模块newcloud包

## 功能

- 当需要在系统中添加一个新的云时，向任意一个云通过http接口上传新的云信息，该云作为主节点；

- 改云接收到新的云信息后，将信息存入数据库，并推送到其他所有云上，等待其他云的投票；

- 收到推送的云，会在云管理系统里看到投票请求，由管理员确定是否要接收新的云，并将投票信息发送给主节点；

- 主节点收到投票后，若投票信息认可新的云，则将vote数加一，若vote数大于一半，则将其加入到云盘系统。

## 系统接口

- PostNewCloud：
  - 接收新的cloud信息
  - 存入本地的TempCloud表和本地的VoteRequest表
  - 调用PostNewCloudVote接口推送云信息

- PostNewCloudVote：接收主节点的推送

  - 获取新云参数

  - 存入VoteRequest表
  - 等待投票

- GetVoteRequest：用于获取所有等待投票的cloud信息

- PostCloudVote：对相关的cloud进行投票

  - 检查自己是否是主节点，若是则直接投票，否则将投票信息通过PostMasterCloudVote发送给主节点
  - 将VoteRequest表中的相关条目删除

- PostMasterCloudVote：

  - 接收其他云的投票信息
  - 更改相关信息的投票数，若超过一半，调用PostCloudSyn向所有的云推送新云，并将新云写入本地mongo

- PostCloudSyn：主节点向其他云推送新云信息

  - 接收新云信息
  - 写入mongo，并将VoteRequst表中的相关条目删除

