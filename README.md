### lego demo 代码layout布局
### 这是使用lego脚手架开发的项目模板 
脚手架参考: https://github.com/jeevic/lego

#### 布局格式:
`
----------
    |
    |-- assets  与存储库一起使用的其他资产(图像、徽标等)
    |-- cmd 项目启动主干
    |-- config  配置目录
    |-- deploy 自动化相关
    |-- internal 私有应用程序和库代码 业务逻辑
    |-- pb  定义数据格式
    |-- test 测试项
    |-- main.go 入口
    |

`

####启动:
`
sh -x run.sh 
8500端口是 http server
8501端口是 grpc server
`
布局参考
参考: https://github.com/golang-standards/project-layout