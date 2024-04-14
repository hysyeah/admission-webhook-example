# admission-webhook-example
admission-webhook-example

系统: ubuntu

部署流程：
1. 构建镜像
`nerdctl build -t hysyeah/pod-webhook:v3 .`
如果使用的容器运行时是docker可以使用下面的命令构建镜像
`docker build -t hysyeah/pod-webhook:v3 .`

2. 配置生成证书
2.1 拷贝/etc/kubernetes/pki/ca.crt,ca.key到./pki目录下(k3s /var/lib/rancher/k3s/server/tls)
2.2 运行gen-certs.sh(会生成对应的证书和创建secret)

3. 部署Deployment或者直接部署二进制文件
`kubectl apply -f deployment.yaml`

或者直接运行服务(保证apiserver在同一个网络(k3s的apiserver有可能是127.0.0.1:6443))，需要修改deployment.yaml中相应的字段，ValidatingWebhookConfiguration，MutatingWebhookConfiguration
`go run *.go -tls-cert-file=./pki/generated/pod-webhook.pem -tls-private-key-file=./pki/generated/pod-webhook-key.pem`

4. 测试验证  
4.1 新建一个`pod`,`kubectl apply -f mycurl.yaml`,查看`pod`的`label`,会发现增了<added-label=yes>
`my-curl                      1/1     Running   0          6s    added-label=yes,app=my-curl`
4.2 新建`configmap`,`kubectl apply -f configmap.yaml`


---
- 安装`cfssl`, `cfssljson`
```shell
wget -q --show-progress --https-only --timestamping \
  https://pkg.cfssl.org/R1.2/cfssl_linux-amd64 \
  https://pkg.cfssl.org/R1.2/cfssljson_linux-amd64
```
`chmod +x cfssl_linux-amd64 cfssljson_linux-amd64`
`sudo mv cfssl_linux-amd64 /usr/local/bin/cfssl`
`sudo mv cfssljson_linux-amd64 /usr/local/bin/cfssljson`

- 安装`nerdctl`
`wget https://github.com/containerd/nerdctl/releases/download/v1.7.5/nerdctl-full-1.7.5-linux-amd64.tar.gz`
`tar -xzvf nerdctl-full-1.7.5-linux-amd64.tar.gz`
`sudo cp bin/nerdctl /usr/local/bin`
`sudo cp bin/buildkitd /usr/local/bin`
```shell
# 启动buildkitd
buildkitd
```
