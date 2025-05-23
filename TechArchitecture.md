# 技术方案

## 系统开发方案

| 模块         | 方案                                                | 特点                                                                                   |
|--------------|-----------------------------------------------------|----------------------------------------------------------------------------------------|
| 编程语言      | Go                                                  | 高性能并发支持，丰富的工具库和轻量化部署，非常适合高性能图片转码服务的开发。                              |
| Web框架      | Gin / Fiber                                         | Gin 是一个高性能且功能丰富的 Go Web 框架，适合应用系统，Fiber 轻量快速，更适合API构建。转码服务是API服务，而不是应用系统，优先Fiber框架，Gin备选。            |
| 缓存         | 本地缓存: BigCache<br>远程缓存: Redis               | 对频繁访问的资源进行缓存，降低磁盘和存储系统的压力，结合本地缓存和分布式缓存，兼顾高性能和数据一致性。          |
| 任务管理      | Temporal / Asynq / Airworkflow                      | 提供任务调度与优先级管理，支持分布式任务的可靠性和灵活性。                                            |
| 消息队列      | RocketMQ / Redis / EventGateway                     | 支持高并发和延时任务调度，适合分布式任务调度场景。                                                  |
| 负载均衡      | 硬件负载均衡（如 F5/Nginx 增强配置）、Apache APISIX  | 提供动态负载均衡、多插件扩展能力，对请求进行智能分流，适合分布式系统环境。                                |
| 转码工具      | ffmpeg / ImageMagick / webpmux                      | ffmpeg 适合图片与视频混合处理，ImageMagick 和 webpmux 适合纯图片格式转换和优化。                   |
| 监控与告警     | Prometheus + Grafana                                | 实时监控系统性能和任务状态，支持自定义告警策略，确保系统运行稳定。                                       |
| 容错机制      | 多活架构 + 自动健康检查                               | 节点故障后自动切换到备用节点，确保服务的高可用性。                                                   |
| 部署方式      | 容器化（Docker）或虚拟机部署 | 云平台（如：AWS、阿里云、Google Cloud）PaaS（如：Heroku、Google App Engine）	容器化便于快速扩展和跨平台部署；虚拟机部署适合传统环境；| 

## 网络与存储方案

| 模块         | 方案                                                | 特点                                                                                   |
|--------------|-----------------------------------------------------|----------------------------------------------------------------------------------------|
| 存储方式      | 对象存储（如：阿里云OSS、AWS S3）	或自建分布式存储（如：Ceph、HDFS、GlusterFS、OpenStack Swift） | 支持高效的图片上传、下载和管理。易于与其他云服务集成。高度分布式，适合大规模文件存储，支持数据冗余和容错。 |
| 处理器加速     | 高性能服务器（CPU + GPU/TPU 加速）<br>混合架构：FPGA + CPU/GPU<br>FPGA 卡（如 Xilinx Alveo 或 Intel Stratix） | CPU 选择高主频多核处理器（如 AMD EPYC、Intel Xeon）；GPU 选择 NVIDIA A100，用于复杂任务加速；TPU 用于 AI 增强任务。<br>FPGA 负责固定模式的高效处理任务（如编码/解码、特效处理），GPU/CPU 负责动态任务处理，组合实现性能与灵活性的平衡。 |
| 网络接口	| 20GbE / 100GbE / 400GbE 网络接口	| 提供高速数据传输，减少带宽瓶颈，适应大规模数据传输需求。采用光纤连接，支持高并发、高吞吐量的传输。| 
| CDN	| 阿里云、华为云、Cloudflare / Akamai / AWS CloudFront	| 使用 CDN 缓存静态资源至多个边缘节点，减少源服务器的带宽压力，提高全球用户的访问速度和稳定性。优化带宽消耗，减少延迟。| 

## 语言选型

| 编程语言  | 场景                      | 说明                                                                                                                                                       |
|-----------|---------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Go 1.20+   | 同步转码请求，实时返回图片流 | 1. Go是编译型语言，生成的可执行文件速度快、体积小。<br>2. Go语言原生支持goroutines和channels，方便并行处理多个图片任务，提高处理速度。<br>3. 简洁语法和易于维护，静态类型和自动内存管理，轻量级，适合业务简单，性能要求高的场景。 |
| Go 1.20+   | 异步任务转码，回调返回图片地址 | 同样适用于异步任务场景，结合框架实现任务调度和资源管理。                                                                                                      |

## 框架选型

| Web框架            | 场景                      | 说明                                                                                                                                                         |
|---------------------|---------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Web: Fiber 2.52     | 同步转码，实时返回图片流    | 同步转码对于性能要求极高，逻辑较为简单：下载图片，调用SDK命令，返回结果即可，适合使用Fiber框架。                                                      |
| Web: Gin 1.10       | 异步任务转码，回调返回图片地址 | 异步转码对于性能要求高，同时有相关调度层，逻辑相对复杂，更适合使用Gin框架来处理。<br>Workflow部分结合Temporal提供可靠的任务调度与资源管理功能。          |

## 任务框架选型

### go管理框架对比

| 框架   | 语言/生态 | 高性能 | 简单性 | 适用场景                   | 优点                                                         | 缺点                                                         |
|--------|----------|--------|--------|--------------------------|--------------------------------------------------------------|--------------------------------------------------------------|
| Temporal | Go       | 高     | 中     | 分布式任务调度，可靠性强       | 强大的工作流支持，内置重试、状态管理，支持多语言客户端       | 学习曲线较高，依赖性强，资源占用较高                       |
| Argo    | Kubernetes | 高     | 中     | 容器化任务，大规模并发       | 完美集成 Kubernetes，天然支持容器化任务，适合云原生环境       | 配置复杂，对 Kubernetes 依赖重，适合已有云原生架构           |
| Asynq   | Go       | 高     | 高     | 单任务队列，轻量级，适合高吞吐量任务调度 | 极简设计，支持延迟任务与优先级队列，快速部署                   | 不适合复杂工作流场景，缺乏对任务依赖的原生支持                 |
| Cadence | Go       | 高     | 中     | 企业级分布式任务调度，适用于需要复杂工作流的应用 | 企业级功能强大，支持复杂任务依赖管理与工作流状态持久化         | 配置和使用复杂，初始部署需要较高运维成本                       |
| gocraft/work | Go       | 中     | 高     | 轻量任务队列，快速集成         | 高效、轻量级，支持任务优先级，简单易用                         | 功能较为单一，缺乏复杂工作流和多语言支持                       |
| Go-Worker | Go       | 中     | 高     | 多任务并发执行，任务优先级管理 | 轻量级框架，快速实现并发任务处理，支持任务的动态分发           | 缺乏高级功能（如任务依赖、状态持久化）                       |
| Go-Task  | Go       | 中     | 高     | 多任务并发调度，任务依赖管理   | 简单易用，支持任务间的依赖管理                                 | 不适合高并发和复杂的任务调度场景，生态较小                       |

## 缓存选型
### Go语言本地缓存

| **缓存框架**   | **主要特点**                                 | **适用场景**                          | **优点**                                           | **缺点**                                             |
|----------------|--------------------------------------------|--------------------------------------|--------------------------------------------------|----------------------------------------------------|
| **sync.Map**   | Go 原生并发安全的键值存储。                      | 简单缓存，小数据量                       | 内置并发安全，无需额外依赖。                               | 不支持过期时间和自动清理，功能简单。                           |
| **BigCache**   | 高性能内存缓存，针对大量数据量优化，分片管理锁。      | 高并发、大数据量场景                    | 高效的内存管理，支持分片，极低的锁争用。                          | 不支持 TTL 功能，需要额外管理数据过期。                        |
| **FreeCache**  | 高性能缓存，使用环形缓冲区结构，支持自动清理和过期。   | 对内存占用敏感，需高效利用内存的场景      | 支持 TTL、过期自动清理，内存利用率高。                          | 没有高级功能（如分布式支持），不支持复杂的存储逻辑。                |
| **gcache**     | 功能全面，支持 TTL、LRU、LFU 等缓存淘汰策略。         | 中小型缓存需求，需多种缓存策略的场景        | 支持多种淘汰策略（LRU、LFU、ARC），内置过期和自动清理。             | 性能可能不如 BigCache 和 FreeCache。                            |
| **ristretto**  | 高性能缓存框架，支持异步加载和自适应缓存策略。        | 大数据量、高并发和复杂策略的场景         | 高性能，支持并发，自动调整缓存策略（通过频率和权重决定缓存）。       | 相对复杂，需要更高的学习成本。                                  |

### 分布式远程缓存

| **缓存框架**      | **性能**                                      | **分布式支持**                                     | **持久化能力**                                     | **数据一致性**                                    | **管理复杂度**                                    | **扩展性**                                      | **成本**                                          | **适用场景**                                      |
|------------------|---------------------------------------------|--------------------------------------------------|--------------------------------------------------|--------------------------------------------------|--------------------------------------------------|--------------------------------------------------|--------------------------------------------------|--------------------------------------------------|
| **Caffeine**     | 极高（专注于内存操作，高并发性能优异）           | 无（需自行实现或结合其他工具，如 Redis）           | 无（数据重启丢失）                                  | 较低（依赖客户端逻辑和清理策略）                     | 非常低（简单易用，专注于内存缓存）                   | 较好（结合其他工具可实现分布式支持）                | 低（适合内存缓存场景，无需外部依赖）                | 高性能、短生命周期的热点数据缓存，专注于内存操作         |
| **Ehcache**      | 高（单机）                                    | 弱（需扩展），适合轻量级应用                          | 支持多种缓存模型                                    | 较高                                              | 简单                                              | 较差                                              | 低                                                | 适合轻量级应用                                    |
| **Redis**        | 极高（分布式），性能优秀，支持复杂数据结构         | 原生支持，功能丰富，支持 RMI 集群模式                 | 支持 AOF、RDB（可选），功能丰富                        | 高，支持主从复制和事务操作                           | 中等（有一定运维成本）                               | 极高                                              | 中等（依赖内存）                                    | 适合需要复杂数据结构和高可用性的场景，适合缓存热点数据     |
| **Memcached**    | 极高（分布式），性能优秀，专注于缓存功能            | 客户端逻辑实现，功能专一，客户端分片，缺乏高可用机制  | 无（数据重启丢失）                                  | 较高，但需要在数据一致性方面做出权衡                | 简单                                              | 简单，缺乏高级扩展能力                              | 低                                                | 适合简单缓存需求，对性能要求高，但预算受限场景          |
| **Couchbase**    | 较高，分布式缓存，自动分片                        | 原生支持，提供持久化存储，分布式架构，自动复制             | 丰富的功能（无明显限制）                              | 高（原生支持分布式一致性）                             | 较高，分布式架构，运维复杂                           | 极高                                              | 高                                                | 适合需要持久化和高可用性的分布式缓存场景，适合长期存储大文件 |


## 缓存策略
为了防止瞬间高并发造成雪崩效应，服务分为2级缓存。
1. bigcache
  - 最长时效为5分钟
  - 设置灵活的缓存策略（如 LFU 淘汰）
  - 内存占用不用超过容器一半
2. redis
  - 有效期为1-3天，支持自定义时效时间，由请求方确定
3. 定期清理 Redis 未访问的缓存数据
  - 避免数据冗余

