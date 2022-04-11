# 简介
m02中的httpserver改造为优雅启停  

### part1
m08中的httpserver.yaml为相应的deployment创建文件，包含：  
- 优雅启动
- 优雅停止
- 资源限制 Burstable
- 配置分离
- 日志级别（代码未引用）

### part2
ingress发布服务，利用cert-manager申请证书