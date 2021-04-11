# jointcloud 云际计算

## 电脑需要安装 node 环境

## 安装依赖命令

```
npm install
```

### 编译及热加载命令

```
npm run dev
```

### 生产环境打包命令

```
npm run build
```

### 项目目录结构 src

```
├── api `接口`
├── assets `图片资源`
│   ├── modules `八大模块图片`
│   ├── partner `合作伙伴图片`
├── components `组件`
│   ├── documents `文件`
│   ├── modules `模块组件`
│   │   ├── cloudBookkeepingAsync `云际异步记账`
│   │   ├── cloudBookkeepingSync `云际同步记账`
│   │   ├── cloudComputation `云际计算`
│   │   ├── cloudStorage `云际存储`
│   │   ├── disasterRecovery `云际容灾`
│   │   ├── serviceRecommendation `云际推荐`
│   │   └── header.vue `模块公共头部`
│   └── .. `其他组件`
├── mixin `复用代码`
├── router `路径声明`
├── scripts/map `列表map（表头）`
├── store `vuex存储`
├── styles `公共样式文件`
├── utils `封装请求及多语言`
├── views `page页面`
│   ├── error-page `模块页面`
│   ├── modules `模块页面`
│   │   ├── cloudBookkeepingAsync `云际异步记账`
│   │   ├── cloudBookkeepingSync `云际同步记账`
│   │   ├── cloudComputation `云际计算`
│   │   ├── cloudStorage `云际存储`
│   │   ├── disasterRecovery `云际容灾`
│   │   └── serviceRecommendation `云际推荐`
└── └── .. `其他页面`
```
