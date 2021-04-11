<template>
  <div>
    <h1>使用方法</h1>
    <h4>1、文件上传</h4>
    <p>
      将文件放入任意slave节点的/data/oec/目录下，并在该节点上输入命令：/data/conf/filectr.sh write 本地文件名 云端存储路径/存储文件名
      ecid文件大小(单位MB)，其中ecid与纠删码类别的对应关系如下表所示。
    </p>
    <table>
      <tr>
        <th>纠删码类别</th>
        <th>ecid</th>
      </tr>
      <tr>
        <td>DFC(12,6,3)</td>
        <td>dfc_12_6_3</td>
      </tr>
      <tr>
        <td>DFC(8,3,3)</td>
        <td>dfc_8_3_2</td>
      </tr>
      <tr>
        <td>RS(6,3,3)</td>
        <td>rs_6_3_2</td>
      </tr>
      <tr>
        <td>RS(6,4,2)</td>
        <td>rs_6_4_3</td>
      </tr>
      <tr>
        <td>RS(8,6,2)</td>
        <td>rs_8_6_4</td>
      </tr>
      <tr>
        <td>RS(9,6,3)</td>
        <td>rs_9_6_3</td>
      </tr>
      <tr>
        <td>RS(12,6,6)</td>
        <td>rs_12_6_2</td>
      </tr>
      <tr>
        <td>RS(12,8,4)</td>
        <td>rs_12_8_3</td>
      </tr>
    </table>

    <h4>2、文件下载</h4>
    <p>在任意slave节点上输入命令（文件将下载到/data/oec/目录下）：/data/conf/filectr.sh read 云端存储路径/存储文件名 本地文件名。</p>
    <h4>3、开启自动修复功能</h4>
    <p>在任意slave节点上输入命令：/data/conf/filectr.sh startRepair。</p>
    <h4>3、关闭自动修复功能</h4>
    <p>在任意slave节点上输入命令：/data/conf/filectr.sh stopRepair。</p>
    <h4>4、访问文件浏览器</h4>
    <p>在浏览器输入“master节点外网IP:8080/uitest.php”即可。</p>
  </div>
</template>
