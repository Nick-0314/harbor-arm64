# harbor基于arm64平台构建v1.10版本镜像

## 变更

Makefile

```
BUILDBIN=false改为BUILDBIN=true    true是编译registry等二进制文件，false是直接下载固定的文件，一般下载的是基于amd64架构的，arm架构的需要编译   所以此次改为true
```

```
VERSIONTAG=dev 改为 VERSIONTAG=v1.10.1   打包的镜像tag
REGISTRYVERSION=v2.7.1-patch-2819-2553    改为  REGISTRYVERSION=v2.7.1  registry的版本
```

```
CLAIRFLAG=true  打包clair镜像
CHARTFLAG=true  打包chartmuseum镜像
```

删除pushimage: 将第329 330行   改为  去掉上传镜像部分

```
    329                  $(DOCKERBUILD) --pull -f $(MAKEFILEPATH_PHOTON)/$$name/Dockerfile.base -t goharbor/harbor-$$name-base:$(BASEIMAGETAG) . ; \
```

修改build步骤， 支持build clair 和 chartmuseum

```
build:
添加：
	 -e CHARTFLAG=$(CHARTFLAG) -e CLAIRFLAG=$(CLAIRFLAG)
```

修改Dockrfile.base中的基础镜像，photon:2.0镜像不支持arm64

```
sed -i 's|photon:2.0|photon:3.0|g' make/photon/*/Dockerfile.base
```

可选：为build clair添加goproxy

```
ENV GOPROXY="http://172.26.1.9:5000"
ENV GOSUMDB="off"
ENV CGO_ENABLED="0"
```

修改redis镜像Dockerfile和redis.conf,解决Unsupported system page size

```
FROM centos:7

COPY ./make/photon/redis/epel.repo /etc/yum.repos.d/
RUN yum install -y redis && yum clean all

```

修改clair的dumb-init为arm64架构



build基础镜像

```
make build_base_docker
```

build镜像

```
make build
```

打包harbor  

```
make package_offline
```



