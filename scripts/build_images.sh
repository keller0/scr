#!/bin/bash
set -e

images=("bash:4.4" \
    "gcc:8.1" "gcc:7.3" \
    "python:2.7-slim" "python:3.5" \
    "php:7.2.5" \
    "java:8" \
    "golang:1.8" "golang:1.10")

# build ric
cd ../cmd/ric && make dbuild
cd -
mv ../cmd/ric/run ./run
echo "---------build ric succeed---------"

# make Dockerfile && build image
build_image () {
    if [ $1 = "bash:4.4" ]
    then
        echo "FROM gcc:8.1" > "Dockerfile"
    else
        echo "FROM $1" > "Dockerfile"
    fi
    cat <<EOF >> "Dockerfile"

RUN groupadd ric && useradd -m -d /home/ric -g ric -s /bin/bash ric
COPY ./run /home/ric/run
RUN chmod +x /home/ric/run

USER ric
EOF
}

#build images
echo

for i in "${images[@]}"
do
    echo "---------start building keller0/$i---------"
    build_image $i
    docker build -t "keller0/$i" .
    echo "---------build keller0/$i succeed---------"
done

rm ./run
rm ./Dockerfile