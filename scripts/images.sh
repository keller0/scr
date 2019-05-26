#!/bin/bash
set -e

images=(
    "bash:4.4" \
    "gcc:8.3" "gcc:7.4" "gcc:6.5" "gcc:5.5" \
    "golang:1.11" "golang:1.12" \
    "haskell:8.6"  \
    "openjdk:8" "openjdk:11" "openjdk:12" "openjdk:13" \
    "perl:5.28" \
    "rakudo-star" \
    "php:7.2.5" \
    "python:2.7" "python:3.7" \
    "ruby:2.6" \
    "rust"
    )


# make Dockerfile && build image
build_one_image () {
    if [ $1 = "bash:4.4" ]
    then
        return 0
    else
        echo "FROM $1" > "Dockerfile"
    fi
    cat <<EOF >> "Dockerfile"

RUN groupadd ric && useradd -m -d /home/ric -g ric -s /bin/bash ric
COPY ./run /home/ric/run
RUN chmod +x /home/ric/run

USER ric
WORKDIR /home/ric/
CMD ["/home/ric/run"]
ENTRYPOINT "/home/ric/run"
EOF

    if [ $1 = "bash:4.4" ]
    then
        return 0
    elif [ $1 = "rakudo-star" ]
    then
        docker build -t "yximages/perl6" .
    else
        docker build -t "yximages/$i" .
    fi
}

#build images
build_local () {
    # build ric
    cd ../cmd/ric && make dbuild
    cd -
    mv ../cmd/ric/run ./run
    echo "---------build ric succeed---------"
    # build ric imagws
    for i in "${images[@]}"
    do
        echo "---------start building yximages/$i---------"
        build_one_image $i
        echo "---------build yximages/$i succeed---------"
    done

    rm -f ./run
    rm ./Dockerfile
}

push_to_ali () {
    for i in "${images[@]}"
    do
        if [ $i = "bash:4.4" ]
        then
            continue
        elif [ $i = "rakudo-star" ]
        then
            docker tag "yximages/$i" "registry.cn-shanghai.aliyuncs.com/yxi/perl6" &
            docker push "registry.cn-shanghai.aliyuncs.com/yxi/perl6"
        else
            docker tag "yximages/$i" "registry.cn-shanghai.aliyuncs.com/yxi/$i" &
            docker push "registry.cn-shanghai.aliyuncs.com/yxi/$i"
        fi

    done
    docker tag "yximages/yxi-api" "registry.cn-shanghai.aliyuncs.com/yxi/yxi-api" &
    docker push "registry.cn-shanghai.aliyuncs.com/yxi/yxi-api"
}

push_to_docker_hub () {
    for i in "${images[@]}"
    do
        if [ $i = "bash:4.4" ]
        then
            continue
        elif [ $i = "rakudo-star" ]
        then
            docker push "yximages/perl6"
        else
            docker push "yximages/$i"
        fi

    done
    docker push "yximages/yxi-api"
}

pull_from_docker_hub () {
    for i in "${images[@]}"
    do
        if [ $i = "bash:4.4" ]
        then
            continue
        elif [ $i = "rakudo-star" ]
        then
            docker pull "yximages/perl6"
        else
            docker pull "yximages/$i"
        fi

    done
    docker pull "yximages/yxi-api"
}

pull_from_ali () {
    for i in "${images[@]}"
    do
        if [ $i = "bash:4.4" ]
        then
            continue
        elif [ $i = "rakudo-star" ]
        then
            docker pull "registry.cn-shanghai.aliyuncs.com/yxi/perl6"
            docker tag "registry.cn-shanghai.aliyuncs.com/yxi/perl6" "yximages/perl6"
        else
            docker pull "registry.cn-shanghai.aliyuncs.com/yxi/$i"
            docker tag "registry.cn-shanghai.aliyuncs.com/yxi/$i" "yximages/$i"
        fi

    done

    docker pull "registry.cn-shanghai.aliyuncs.com/yxi/yxi-api"
    docker tag  "registry.cn-shanghai.aliyuncs.com/yxi/yxi-api" "yximages/yxi-api"
}

print_usage() {
  echo "    -b build, -a push to aliyun, -d push to docker hub\n
    -p pull images from docker hub, -pa pull images from aliyun"
}


while getopts 'abddplh' flag; do
  case "${flag}" in
    a) push_to_ali ;;
    b) build_local ;;
    d) push_to_docker_hub ;;
    p) pull_from_docker_hub ;;
    l) pull_from_ali ;;
    h) print_usage
       exit 1 ;;
  esac
  exit 0;
done
print_usage