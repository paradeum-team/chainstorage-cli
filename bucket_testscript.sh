#!/bin/bash

estatus=0

cmdPath='./gcscmd_linux'

testCases() {

  echo "_/_/_/_/_/_/_/_/_/_/_/_/_/ 桶操作 开始 _/_/_/_/_/_/_/_/_/_/_/_/_/"
  # todo: 数据准备，删除所有桶，还是切换APIKEY？
  # 桶操作 => 列出桶对象 => 无桶 => gcscmd ls
  execCmd '桶操作' '列出桶对象' '无桶' 'gcscmd ls' 'ls' ''

  # 桶操作 => 创建桶 => 正常创建 => gcscmd mb cs://bbb
  bucketName="bucket-create-"$(date "+%Y%m%d%H%M%S")"-"$RANDOM
  execCmd '桶操作' '创建桶' '正常创建' 'gcscmd mb cs://'$bucketName 'mb cs://'$bucketName ''

  # 桶操作 => 创建桶 => 非正常创建-桶名重复 => gcscmd mb cs://bbb
  # 此操作应给出错误提示
  execCmdFail '桶操作' '创建桶' '非正常创建-桶名重复' 'gcscmd mb cs://'$bucketName 'mb cs://'$bucketName ''

  # todo: 测试意图？
  # 桶操作 => 创建桶 => 非正常创建-创建多个桶 => gcscmd mb cs://aaa
  execCmdFail '桶操作' '创建桶' '非正常创建-创建多个桶' 'gcscmd mb cs://'$bucketName 'mb cs://'$bucketName ''

  # 桶操作 => 列出桶对象 => 有桶 => gcscmd ls
  # 返回对应bucketName桶数据
  execCmd '桶操作' '列出桶对象' '有桶' 'gcscmd ls' 'ls' ''

  # 桶操作 => 移除桶 => 正常删除-无数据删除 => gcscmd rb cs://bbb
  execCmd '桶操作' '移除桶' '正常删除-无数据删除' 'gcscmd rb cs://'$bucketName 'rb cs://'$bucketName ''

  # 桶操作 => 移除桶 => 重复删除-继续删除已删除的桶 => gcscmd rb cs://bbb
  # 此操作应给出错误提示
  execCmdFail '桶操作' '移除桶' '重复删除-继续删除已删除的桶' 'gcscmd rb cs://'$bucketName 'rb cs://'$bucketName ''

  # 桶操作 => 移除桶 => 正常删除-有数据删除 => gcscmd rb cs://bbb
  # 数据准备
  # 1、创建新桶
  bucketName="bucket-create-"$(date "+%Y%m%d%H%M%S")"-"$RANDOM
  $cmdPath mb cs://$bucketName
  # 2、添加对象
  testDataFileName="testdata_"$(date "+%Y%m%d%H%M%S")"-"$RANDOM".dat"
  #echo $testDataFileName
  dd if=/dev/urandom of=$testDataFileName bs=1024 count=1024
  $cmdPath put $testDataFileName cs://$bucketName
  # 此操作应给出错误提示
  execCmdFail '桶操作' '移除桶' '正常删除-有数据删除' 'gcscmd rb cs://'$bucketName 'rb cs://'$bucketName ''
  # 3、数据清理
  rm -rf $testDataFileName

  # 桶操作 => 移除桶 => 正常删除-有数据强制删除 => gcscmd rb cs://bbb --force
  execCmd '桶操作' '移除桶' '正常删除-有数据强制删除' 'gcscmd rb cs://'$bucketName' --force' 'rb cs://'$bucketName' --force' ''

  # 桶操作 => 清空桶 => 正常清空-有数据清空 => gcscmd rm cs://bbb
  # 数据准备
  # 1、创建新桶
  bucketName="bucket-create-"$(date "+%Y%m%d%H%M%S")"-"$RANDOM
  $cmdPath mb cs://$bucketName
  # 2、添加对象
  testDataFileName="testdata_"$(date "+%Y%m%d%H%M%S")"-"$RANDOM".dat"
  #echo $testDataFileName
  dd if=/dev/urandom of=$testDataFileName bs=1024 count=1024
  $cmdPath put $testDataFileName cs://$bucketName
  execCmd '桶操作' '清空桶' '正常清空-有数据清空' 'gcscmd rm cs://'$bucketName 'rm cs://'$bucketName ''

  # 桶操作测试数据清理
  $cmdPath rb cs://$bucketName --force
  rm -rf $testDataFileName

  echo "_/_/_/_/_/_/_/_/_/_/_/_/_/ 桶操作 结束 _/_/_/_/_/_/_/_/_/_/_/_/_/"
}

execCmd() {
  testModule=$1
  testFunction=$2
  testCase=$3
  testDescription=$4
  testCmd=$5
  testExpectation=$6
  testFail=$7

  echo $testModule"=>"$testFunction"=>"$testCase"=>"$testDescription
  cmdStr=$cmdPath' '$testCmd
  echo 'executing '$cmdStr
  eval $cmdStr
  if [ $? -eq 0 ]; then
    echo -e "\033[32m Success: $cmdStr test pass.  $logs \033[0m"
  else
    echo -e "\033[31m Failure: $cmdStr test fail. \033[0m"
    estatus=$(($etatus + 1))
  fi
  echo ""
}

execCmdFail() {
  testModule=$1
  testFunction=$2
  testCase=$3
  testDescription=$4
  testCmd=$5
  testExpectation=$6
  testFail=$7

  echo $testModule"=>"$testFunction"=>"$testCase"=>"$testDescription
  cmdStr=$cmdPath' '$testCmd
  echo 'executing '$cmdStr
  eval $cmdStr
  if [ $? -eq 0 ]; then
    echo -e "\033[31m Failure: $cmdStr test should be prompt error message. \033[0m"
    estatus=$(($etatus + 1))
  else
    echo -e "\033[32m Success: $cmdStr test pass, return error message and exit code is not equal zero.  $logs \033[0m"
  #  estatus=$?
  fi
  echo ""
}

echo "===========================Chainstorage cli Test start=========================="
testCases
echo "===========================Chainstorage cli Test end=========================="

exit $estatus
