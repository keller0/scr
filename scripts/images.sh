#!/bin/bash
set -e

images=(
    "bash:4.4" \
    "gcc:10" \
    "golang:1.18" \
    "haskell:8.10"  \
    "openjdk:14" \
    "perl:5.28" \
    "rakudo-star" \
    "php:7.4" \
    "python:3.8" \
    "ruby:2.7" \
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

print_usage() {
#  echo "    -b build, -a push to aliyun, -d push to docker hub\n
#    -p pull images from docker hub, -pa pull images from aliyun"

  echo "    -b build"
}


while getopts 'bddph' flag; do
  case "${flag}" in
    b) build_local ;;
#    d) push_to_docker_hub ;;
#    p) pull_from_docker_hub ;;
    h) print_usage
       exit 1 ;;
  esac
  exit 0;
done
print_usage