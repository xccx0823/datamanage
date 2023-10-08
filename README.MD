## 如何在服务器上使用docker部署kafka

### 创建网络

```shell
docker network create zookeeper_networ
```

### 部署zookeeper

```yaml
version: '3'

networks:
  default:
    external:
      name: zookeeper_network

services:
  zoo1:
    image: zookeeper:3.7.0
    container_name: zoo1
    hostname: zoo1
    ports:
      - 2181:2181
    volumes:
      - "/root/kafka_learn/zookeeper/zoo1/data:/data"
      - "/root/kafka_learn/zookeeper/zoo1/datalog:/datalog"
    environment:
      ZOO_MY_ID: 1
      ZOO_SERVERS: server.1=0.0.0.0:2888:3888;2181 server.2=zoo2:2888:3888;2181 server.3=zoo3:2888:3888;2181

  zoo2:
    image: zookeeper:3.7.0
    container_name: zoo2
    hostname: zoo2
    ports:
      - 2182:2181
    volumes:
      - "/root/kafka_learn/zookeeper/zoo2/data:/data"
      - "/root/kafka_learn/zookeeper/zoo2/datalog:/datalog"
    environment:
      ZOO_MY_ID: 2
      ZOO_SERVERS: server.1=zoo1:2888:3888;2181 server.2=0.0.0.0:2888:3888;2181 server.3=zoo3:2888:3888;2181

  zoo3:
    image: zookeeper:3.7.0
    container_name: zoo3
    hostname: zoo3
    ports:
      - 2183:2181
    volumes:
      - "/root/kafka_learn/zookeeper/zoo3/data:/data"
      - "/root/kafka_learn/zookeeper/zoo3/datalog:/datalog"
    environment:
      ZOO_MY_ID: 3
      ZOO_SERVERS: server.1=zoo1:2888:3888;2181 server.2=zoo2:2888:3888;2181 server.3=0.0.0.0:2888:3888;2181
```

### 部署kafka

```yaml
version: '3'

networks:
  default:
    external:
      name: zookeeper_network

services:
  kafka1:
    image: wurstmeister/kafka:2.12-2.4.1
    restart: unless-stopped
    container_name: kafka1
    hostname: kafka1
    ports:
      - "9092:9092"
    external_links:
      - zoo1
      - zoo2
      - zoo3
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_LISTENERS: PLAINTEXT://0.0.0.0:9092
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://150.158.140.174:9092    ## 宿主机IP
      KAFKA_ADVERTISED_HOST_NAME: kafka1
      KAFKA_ADVERTISED_PORT: 9092
      KAFKA_ZOOKEEPER_CONNECT: "zoo1:2181,zoo2:2181,zoo3:2181"
    volumes:
      - "/root/kafka_learn/kafka/kafka1/data/:/kafka"


  kafka2:
    image: wurstmeister/kafka:2.12-2.4.1
    restart: unless-stopped
    container_name: kafka2
    hostname: kafka2
    ports:
      - "9093:9092"
    external_links:
      - zoo1
      - zoo2
      - zoo3
    environment:
      KAFKA_BROKER_ID: 2
      KAFKA_LISTENERS: PLAINTEXT://0.0.0.0:9092
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://150.158.140.174:9093    ## 宿主机IP
      KAFKA_ADVERTISED_HOST_NAME: kafka2
      KAFKA_ADVERTISED_PORT: 9093
      KAFKA_ZOOKEEPER_CONNECT: "zoo1:2181,zoo2:2181,zoo3:2181"
    volumes:
      - "/root/kafka_learn/kafka/kafka2/data/:/kafka"

  kafka3:
    image: wurstmeister/kafka:2.12-2.4.1
    restart: unless-stopped
    container_name: kafka3
    hostname: kafka3
    ports:
      - "9094:9092"
    external_links:
      - zoo1
      - zoo2
      - zoo3
    environment:
      KAFKA_BROKER_ID: 3
      KAFKA_LISTENERS: PLAINTEXT://0.0.0.0:9092
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://150.158.140.174:9094   ## 宿主机IP
      KAFKA_ADVERTISED_HOST_NAME: kafka3
      KAFKA_ADVERTISED_PORT: 9094
      KAFKA_ZOOKEEPER_CONNECT: "zoo1:2181,zoo2:2181,zoo3:2181"
    volumes:
      - "/root/kafka_learn/kafka/kafka3/data/:/kafka"

  kafka-manager: # Kafka 图形管理界面
    image: sheepkiller/kafka-manager:latest
    restart: unless-stopped
    container_name: kafka-manager
    hostname: kafka-manager
    ports:
      - "9000:9000"
    links: # 连接本compose文件创建的container
      - kafka1
      - kafka2
      - kafka3
    external_links: # 连接外部compose文件创建的container
      - zoo1
      - zoo2
      - zoo3
    environment:
      ZK_HOSTS: zoo1:2181,zoo2:2181,zoo3:2181
      KAFKA_BROKERS: kafka1:9092,kafka2:9093,kafka3:9094
```

