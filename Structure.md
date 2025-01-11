# 工程目录结构

```shell
frog-go/
├── cmd/                           # 项目的主入口目录，包含应用程序的启动代码
│   └── vangogh/                   # main.go 文件所在目录，启动程序入口
│       └── main.go                # 应用的入口文件，配置路由、启动 HTTP 服务等
├── internal/                      # 存放项目内部使用的代码，外部不能引用
│   ├── handlers/                  # HTTP 请求处理函数，包含具体的处理逻辑
│   │   └── image_handler.go       # 处理图片转码相关的 HTTP 请求
│   ├── models/                    # 数据模型层，定义和数据库或应用相关的数据结构
│   │   └── image.go               # 图片相关的数据结构，存储图片信息
│   ├── services/                  # 业务逻辑层，处理转码、队列等业务逻辑
│   │   └── transcoding_service.go # 处理图片转码的核心业务逻辑
│   ├── utils/                     # 工具函数，提供通用的功能
│   │   └── image_utils.go         # 图片处理工具函数，裁剪、缩放等
│   ├── cache/                     # 缓存相关功能，提升性能
│   │   └── cache.go               # 使用缓存加速图片转码的操作
│   ├── repositories/              # 数据库或其他持久化存储交互的封装，同步转码无
│   │   └── image_repository.go    # 与图片存储相关的数据库交互
│   ├── tools/                     # 第三方工具封装（如 ffmpeg、imagicmaker 等）
│   │   ├── ffmpeg_api.go          # 调用 ffmpeg 或 imagicmaker的逻辑封装。
│   │   └── imagicmaker_api.go     # 改为编码云统一提供的api
│   └── routes/                    # 路由配置目录
│       └── router.go              # 配置路由和请求映射的地方
├── pkg/                           # 公共包，项目中可复用的工具、服务等
│   ├── logger/                    # 日志模块，记录系统运行时的日志
│   │   └── logger.go              # 自定义的日志处理函数
│   ├── response/                  # API 响应格式化模块
│   │   └── response.go            # 定义统一的响应结构
│   └── config/                    # 配置模块，读取和管理应用配置
│       └── config.go              # 配置读取、解析逻辑
├── tests/                         # 测试代码目录
│   ├── handlers/                  # 测试控制器相关的代码
│   │   └── image_handler_test.go  # 测试图片转码 HTTP 请求的处理
│   ├── services/                  # 测试服务层代码
│   │   └── transcoding_service_test.go # 测试图片转码的核心业务逻辑
│   └── utils/                     # 测试工具函数相关的代码
│       └── image_utils_test.go    # 测试图片处理工具函数
├── config/                        # 存放配置文件
│   ├── config.yaml                # 系统的配置文件，例如 ffmpeg 的路径、缓存策略等
├── go.mod                         # Go 模块文件，包含依赖管理
├── go.sum                         # Go 模块校验文件
└── README.md                      # 项目说明文档，介绍如何使用该项目
```