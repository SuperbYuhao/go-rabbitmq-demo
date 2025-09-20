## 生成证书
##### 创建根证书颁发机构（CA）
```
openssl genrsa -out ca_key.pem 2048  
openssl req -x509 -new -nodes -key ca_key.pem -days 3650 -out ca_certificate.pem -subj "/CN=MyCA"  
```

##### 创建服务端私钥
```
openssl genrsa -out server_key.pem 2048  
```

##### 创建配置文件 ssl.conf
```
cat > ssl.conf <<EOF  
[req]  
req_extensions = v3_req  
distinguished_name = req_distinguished_name  
  
[req_distinguished_name]  
  
[v3_req]  
basicConstraints = CA:FALSE  
keyUsage = nonRepudiation, digitalSignature, keyEncipherment  
subjectAltName = @alt_names  
  
[alt_names]  
IP.1 = 0.0.0.0  # 这里需要根据实际改地址
DNS.1 = rabbitmq-server  
EOF
```

##### 生成服务器证书请求（包含SAN）
```
openssl req -new -key server_key.pem -out server.csr -subj "/CN=rabbitmq-server" -config ssl.conf  
```

##### 使用CA签名（包含SAN扩展）
```
openssl x509 -req -in server.csr \  
-CA ca_certificate.pem \  
-CAkey ca_key.pem \  
-CAcreateserial \  
-out server_certificate.pem \  
-days 365 \  
-extensions v3_req \  
-extfile ssl.conf
```

##### 生成客户端私钥（client_key.pem）和证书签名请求（client.csr）
```
openssl genrsa -out client_key.pem 2048  
openssl req -new -key client_key.pem -out client.csr -subj "/CN=rabbitmq-client"
```

##### 使用根 CA 对客户端证书请求进行签名，生成客户端证书（client_certificate.pem）
```
openssl x509 -req -in client.csr -CA ca_certificate.pem -CAkey ca_key.pem -CAcreateserial -out client_certificate.pem -days 365
```

## 编写 rabbitmq.conf 文件
```rabbitmq.conf
listeners.ssl.default = 5671  
ssl_options.cacertfile = /etc/rabbitmq/ssl/ca_certificate.pem  
ssl_options.certfile = /etc/rabbitmq/ssl/server_certificate.pem  
ssl_options.keyfile = /etc/rabbitmq/ssl/server_key.pem  
ssl_options.verify = verify_peer  
ssl_options.fail_if_no_peer_cert = false  
  
management.ssl.port = 15671  
management.ssl.cacertfile = /etc/rabbitmq/ssl/ca_certificate.pem  
management.ssl.certfile = /etc/rabbitmq/ssl/server_certificate.pem  
management.ssl.keyfile = /etc/rabbitmq/ssl/server_key.pem
```

## 编写 docker-compose.yml 文件
```docker-compose.yml
services:  
  rabbitmq:  
    image: rabbitmq:3.9-management  # 使用带有管理插件的RabbitMQ镜像  
    container_name: rabbitmq  # 容器名称  
    ports:  
      - "5671:5671"  # RabbitMQ的AMQPs端口  
      - "5672:5672"  # RabbitMQ的AMQP端口  
      - "15671:15671"  # RabbitMQ的HTTPS管理界面端口  
      - "15672:15672"  # RabbitMQ的HTTPS管理界面端口  
    environment:  
      RABBITMQ_DEFAULT_USER: admin  # 设置默认用户名  
      RABBITMQ_DEFAULT_PASS: password  # 设置默认密码  
      RABBITMQ_CONFIG_FILE: /etc/rabbitmq/rabbitmq.conf  
    volumes:  
      - /Users/yuhao/etc/rabbitmq/ssl:/etc/rabbitmq/ssl  
      - ./rabbitmq.conf:/etc/rabbitmq/rabbitmq.conf
```