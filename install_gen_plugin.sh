#!/usr/bin/env bash

clear
array=(${GOPATH//:/ }) #按:截取字符串
len=${#array[@]} #获取数组长度
if [ $len -eq 0 ]; then
    echo "请先配置GOPATH"
    exit
fi

# go build 忽略 "_","." 开头的go文件
$(go build cmd/main.go) #编译
# 重命名并移动到 GOPATH/bin
$(mv main "$array/bin/protoc-gen-tmpl")
echo "success mv to  $array/bin/protoc-gen-tmpl"
#生成PB文件
rm -rf rpc/*
#for dir in proto/*; do
#    # -I . 生成含全路径, -I proto/ 则不不包含    不带plugins=grpc默认只生成message
#   $(protoc  -I proto/ ${dir}/*.proto --go_out=plugins=grpc:rpc)
#done

#生成tmpl代码模板
rm -rf tmpl/*
for dir in proto/*; do
    # -I . 生成含全路径, -I proto/ 则不不包含
   $(protoc  -I proto/ ${dir} --tmpl_out=plugins=tmpl:tmpl/)
done

#  protoc  --proto_path=. --mimi_out=. proto/order/order.proto  这样无法生成的原因