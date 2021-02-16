# jcsPan Transporter
## 主要功能
 - 暂存用户上传的文件
 - 将用户上传的文件同步到云服务器
 - 不同服务器之间文件的调度

## 使用场景
- 用户上传文件到网盘
    1. 前端构造表单，发起文件上传请求，请求中包括 sessionID、上传目的网盘路径、文件内容
    2. Transporter 请求 DB 对 sessionID 进行校验，验证成功后对用户上传的文件进行暂存并通过路径获取存储方案
    3. Transporter 在 DB 中建立一个同步任务，标记该任务为待同步
    4. Transporter 从任务队列中获取任务，使用 minio 同步本地文件至 s3 存储中
    5. 上传至s3完成后，删除临时文件，标记任务完成
    
- 用户从网盘下载文件
    1. 前端构造表单，发起文件下载请求，请求中包括 sessionID、下载的网盘路径
    2. Transporter 请求 DB 对 sessionID 进行校验，验证成功后通过路径获取存储方案
    3. Transporter 请求 s3 生成一个 presigned URL，返回给用户