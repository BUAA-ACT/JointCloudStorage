<template>
  <div>
    <h1>手动部署流程（购买资源时已部署完成）</h1>
    <h4>1、在所有节点上安装Docker及Weave（以Ubuntu16.04.6 LTS为例）</h4>
    <p>
      a) curl -sSL https://get.daocloud.io/docker | sh <br />
      b) wget -O /usr/local/bin/weave https://raw.githubusercontent.com/zettio/weave/master/weave<br />
      c) chmod a+x /usr/local/bin/weave<br />
      d) weave launch
    </p>

    <h4>2、将文件夹images文件夹conf复制到给所有节点的/data/目录下，并按以下说明修改conf中的配置文件和脚本</h4>
    <p>
      (1) startmaster.sh（仅master节点需修改）<br />
      a)
      修改startmaster.sh里的第3步和第5步中的IP，其中第3步中使用云主机的物理IP，第5步中使用用户自行为各个容器分配的虚拟IP（所有节点的所有容器的虚拟IP处于同一内网即可）。<br />

      (2) startslave.sh（仅slave节点需修改）<br />
      a)
      修改startslave.sh里的第3步和第5步中的IP，其中第3步中使用云主机的物理IP，第5步中使用用户自行为各个容器分配的虚拟IP（所有节点的所有容器的虚拟IP处于同一内网即可）。<br />

      (3) core-site.xml<br />
      a) 将fs.defaultFS标签的value改为“hdfs://master上的bh-namenode容器的虚拟IP:9000”。<br />

      (4) hdfs-site.xml<br />
      a) 将oec.controller.addr标签的value改为master上的bh-redis容器的虚拟IP。<br />
      b) 将oec.local.addr标签的value改为本机上的bh-redis容器的虚拟IP。<br />

      (5) workers<br />
      a) 填入所有bh-datanode容器的虚拟IP，每个IP占一行。<br />

      (6) sysSetting.xml<br />
      a) 将controller.addr标签下的oecaddr、hdfsaddr、redisaddr子标签内的值分别改为master上的bh-coordinator、bh-namenode、bh-redis容器的虚拟IP。<br />
      b)
      将agents.addr标签下的oecaddr、hdfsaddr、redisaddr子标签内的值分别改为各个slave上的bh-worker、bh-datanode、bh-redis容器的虚拟IP（每个slave对应一组oecaddr、hdfsaddr、redisaddr子标签）。<br />
      c) 将mysql.addr标签的value改为master上bh-lamp容器的虚拟IP。<br />
      d)
      将local.addr标签下的oecaddr、hdfsaddr、redisaddr子标签内的值分别改为本机上的bh-coordinator/bh-worker、bh-namenode/bh-datanode、bh-redis容器的虚拟IP。<br />
      e) 将dss.parameter标签的value改为master上的bh-namenode容器的虚拟IP。<br />
    </p>
    <h4>3、在master节点执行startmaster.sh</h4>
    <p>
      a) chmod +x /data/conf/startmaster.sh<br />
      b) /data/conf/startmaster.sh
    </p>

    <h4>4、在slave节点执行startslave.sh</h4>
    <p>
      a) chmod +x /data/conf/startslave.sh<br />
      b) /data/conf/startslave.sh
    </p>
  </div>
</template>
