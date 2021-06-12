module.exports = {
    chainWebpack: config => {
        config
            .plugin('html')
            .tap(args => {
                args[0].title= '云际存储 Joint Cloud Storage —— 导航'
                return args
            })
    }
}