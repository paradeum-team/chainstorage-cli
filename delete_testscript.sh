#!/bin/bash

estatus=0

cmdPath='./gcscmd_linux'

testCases() {

  echo "_/_/_/_/_/_/_/_/_/_/_/_/_/ 删除对象 开始 _/_/_/_/_/_/_/_/_/_/_/_/_/"
  # 删除对象 => 清空桶 => 有文件清空桶 => gcscmd rm cs://bbb --force
  # 数据准备
  # 1、创建新桶
  bucketName="bucket-delete-"$(date "+%Y%m%d%H%M%S")"-"$RANDOM
  $cmdPath mb cs://$bucketName
  # 2、添加对象，10MB
  testDataFileName="testdata_"$(date "+%Y%m%d%H%M%S")"-"$RANDOM".dat"
  #echo $testDataFileName
  dd if=/dev/urandom of=$testDataFileName bs=10240 count=1024
  $cmdPath put $testDataFileName cs://$bucketName
  execCmd '删除对象' '清空桶' '有文件清空桶' 'gcscmd rm cs://'$bucketName' --force' 'rm cs://'$bucketName' --force' ''
  # 3、数据清理
  rm -rf $testDataFileName

  # 删除对象 => 使用对象名删除单文件 => 正确删除 => gcscmd rm cs://bbb --name Tarkov.mp4
  # 数据准备
  # 1、添加对象，10MB
  testDataFileName="testdata_"$(date "+%Y%m%d%H%M%S")"-"$RANDOM".dat"
  #echo $testDataFileName
  dd if=/dev/urandom of=$testDataFileName bs=10240 count=1024
  $cmdPath put $testDataFileName cs://$bucketName
  # 设置对象名称
  objectName=$testDataFileName
  execCmd '删除对象' '使用对象名删除单文件' '正确删除' 'gcscmd rm cs://'$bucketName' --name '$objectName 'rm cs://'$bucketName' --name '$objectName ''
  # 2、数据清理
  rm -rf $testDataFileName

  # 删除对象 => 使用模糊查询删除对象 => 正常模糊删除 => gcscmd rm cs://bbb --name .mp4 --force
  # 1、添加对象1，10MB
  testDataFileName="testdata_"$(date "+%Y%m%d%H%M%S")"-"$RANDOM".dat"
  #echo $testDataFileName
  dd if=/dev/urandom of=$testDataFileName bs=10240 count=1024
  $cmdPath put $testDataFileName cs://$bucketName
  # 2、添加对象2，5MB
  testDataFileName2="testdata_"$(date "+%Y%m%d%H%M%S")"-"$RANDOM".dat"
  #echo $testDataFileName
  dd if=/dev/urandom of=$testDataFileName2 bs=5120 count=1024
  $cmdPath put $testDataFileName2 cs://$bucketName
  # 设置对象名称
  objectName="testdata"
  execCmd '删除对象' '使用模糊查询删除对象' '正常模糊删除' 'gcscmd rm cs://'$bucketName' --name '$objectName' --force' 'rm cs://'$bucketName' --name '$objectName' --force' ''
  # 3、数据清理
  rm -rf $testDataFileName
  rm -rf $testDataFileName2

  # 删除对象 => 使用对象名删除单目录 => 正确删除-目录中无文件 => gcscmd rm cs://bbb --name aaa
  # 1、添加对象
  testDataFolderame="testfolder_"$(date "+%Y%m%d%H%M%S")"-"$RANDOM
  #echo $testDataFolderame
  mkdir $testDataFolderame
  $cmdPath put $testDataFolderame cs://$bucketName
  # 设置对象名称
  objectName=$testDataFolderame
  execCmd '删除对象' '使用对象名删除单目录' '正确删除-目录中无文件' 'gcscmd rm cs://'$bucketName' --name '$objectName 'rm cs://'$bucketName' --name '$objectName ''
  # 2、数据清理
  rm -rf $testDataFolderame

  # 删除对象 => 使用CID删除单对象 => 正确删除-对应的桶中有此CID删除 => gcscmd rm cs://bbb --cid QmWgnG7pPjG31w328hZyALQ2BgW5aQrZyKpT47jVpn8CNo
  # 数据准备
  # 1、添加对象，10MB
  testDataFileName="testdata_"$(date "+%Y%m%d%H%M%S")"-"$RANDOM".dat"
  #echo $testDataFileName
  dd if=/dev/urandom of=$testDataFileName bs=10240 count=1024
  resp=$($cmdPath put $testDataFileName cs://$bucketName)
  # 2、获取CID
  cid=$(echo "$resp"|awk '/CID:/{print $2}')
  execCmd '删除对象' '使用CID删除单对象' '正确删除-对应的桶中有此CID删除' 'gcscmd rm cs://'$bucketName' --cid '$cid 'rm cs://'$bucketName' --cid '$cid ''
  # 3、数据清理
  rm -rf $testDataFileName

  # 删除对象 => 使用 CID 删除多个对象(命中多个对象时加) => 一个cid有多个对象 => gcscmd rm cs://bbb --cid QmWgnG7pPjG31w328hZyALQ2BgW5aQrZyKpT47jVpn8CNo --force
  # 数据准备
  # 1、添加对象，10MB
  testDataFileName="testdata_"$(date "+%Y%m%d%H%M%S")"-"$RANDOM".dat"
  #echo $testDataFileName
  dd if=/dev/urandom of=$testDataFileName bs=10240 count=1024
  $cmdPath put $testDataFileName cs://$bucketName
  # 2、复制对象，10MB
  testDataFileName2="testdata_"$(date "+%Y%m%d%H%M%S")"-"$RANDOM".dat"
  #echo $testDataFileName2
  cp $testDataFileName $testDataFileName2
  resp=$($cmdPath put $testDataFileName2 cs://$bucketName)
  # 3、获取CID
  cid=$(echo "$resp"|awk '/CID:/{print $2}')
  execCmd '删除对象' '使用 CID 删除多个对象(命中多个对象时加)' '一个cid有多个对象' 'gcscmd rm cs://'$bucketName' --cid '$cid' --force' 'rm cs://'$bucketName' --cid '$cid' --force' ''
  # 4、数据清理
  rm -rf $testDataFileName
  rm -rf $testDataFileName2

  # 删除对象测试数据清理
  $cmdPath rb cs://$bucketName --force

  echo "_/_/_/_/_/_/_/_/_/_/_/_/_/ 删除对象 结束 _/_/_/_/_/_/_/_/_/_/_/_/_/"

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
