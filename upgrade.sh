tag=$(date '+%m%d%H%M%S')
image_name=base-webhook
#清除以前的镜像
docker exec -it kind-control-plane bash -c "crictl images | grep docker.io/library/${image_name} | awk '{print \$1\":\"\$2}' |  xargs crictl rmi"
#清除以前的镜像，否则通过kind load docker-image加载的镜像无效
docker images | grep ${image_name} | awk '{print $1":"$2}' |  xargs docker rmi
make docker-build IMG=${image_name}:${tag}
kind load docker-image ${image_name}:${tag}
sleep 3
make deploy IMG=${image_name}:${tag}