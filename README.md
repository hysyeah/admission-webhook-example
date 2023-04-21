# admission-webhook-example
admission-webhook-example

部署流程：
1. 构建镜像
`nerdctl build -t hysyeah/pod-webhook:v3 .`

2. 配置生成证书  
2.1 拷贝/etc/kubernetes/pki/ca.crt,ca.key到./pki目录下  
2.2 运行gen-certs.sh(会生成对应的证书和创建secret)

3. 部署Deployment
`kubectl apply -f deployment.yaml`

4. 测试验证
4.1 新建一个`pod`,`kubectl apply -f mycur.yaml`,查看`pod`的`label`
4.2 新建`configmap`,`kubectl apply -f configmap.yaml`
