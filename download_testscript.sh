#!/bin/bash

estatus=0

cmdPath='./gcscmd'

testCases() {

  echo "_/_/_/_/_/_/_/_/_/_/_/_/_/ 下载对象 开始 _/_/_/_/_/_/_/_/_/_/_/_/_/"
  echo '下载对象 => 根据cid下载 => cid正确下载 => gcscmd get cs://bbb --cid QmWgnG7pPjG31w328hZyALQ2BgW5aQrZyKpT47jVpn8CNo start'
  echo '数据准备'
  echo '1、创建新桶'
  bucketName="bucket-download-"$(date "+%Y%m%d%H%M%S")"-"$RANDOM
  $cmdPath mb cs://$bucketName
  echo '2、添加对象，10MB'
  testDataFileName="testdata_"$(date "+%Y%m%d%H%M%S")"-"$RANDOM".dat"
  #echo $testDataFileName
  dd if=/dev/urandom of=$testDataFileName bs=10240 count=1024
  resp=$($cmdPath put $testDataFileName cs://$bucketName)
  echo '3、获取CID'
  cid=$(echo "$resp" | awk '/CID:/{print $2}')
  echo '4、清理上传数据'
  rm -rf $testDataFileName
  execCmd '下载对象' '根据cid下载' 'cid正确下载' 'gcscmd get cs://'$bucketName' --cid '$cid 'get cs://'$bucketName' --cid '$cid ''
  echo '5、数据清理'
  rm -rf $testDataFileName
  echo ''

  echo '下载对象 => 根据对象名下载 => 对象名正确下载 => gcscmd get cs://bbb --name Tarkov.mp4 start'
  echo '数据准备'
  echo '1、添加对象，10MB'
  testDataFileName="testdata_"$(date "+%Y%m%d%H%M%S")"-"$RANDOM".dat"
  #echo $testDataFileName
  dd if=/dev/urandom of=$testDataFileName bs=10240 count=1024
  resp=$($cmdPath put $testDataFileName cs://$bucketName)
  echo '2、获取CID'
  cid=$(echo "$resp" | awk '/CID:/{print $2}')
  echo '3、清理上传数据'
  rm -rf $testDataFileName
  # 设置对象名称
  objectName=$testDataFileName
  execCmd '下载对象' '根据对象名下载' '对象名正确下载' 'gcscmd get cs://'$bucketName' --name '$objectName 'get cs://'$bucketName' --name '$objectName ''
  echo '4、数据清理'
  rm -rf $testDataFileName
  echo ''

  echo '下载对象测试数据清理'
  $cmdPath rb cs://$bucketName --force
  #  rm -rf $testDataFileName
  echo "_/_/_/_/_/_/_/_/_/_/_/_/_/ 下载对象 结束 _/_/_/_/_/_/_/_/_/_/_/_/_/"
  echo ''
}

execCmdFail() {
  testModule=$1
  testFunction=$2
  testCase=$3
  testDescription=$4
  testCmd=$5
  testExpectation=$6
  testFail=$7

  #  echo $testModule"=>"$testFunction"=>"$testCase"=>"$testDescription
  cmdStr=$cmdPath' '$testCmd
  echo 'executing '$cmdStr
  eval $cmdStr

  exitCode=$?
  echo ''
  if [ $exitCode -eq 0 ]; then
    echo -e "\033[31mFailure: $cmdStr test should be prompt error message. \033[0m"
    estatus=$(($etatus + 1))
  else
    echo -e "\033[32mSuccess: $cmdStr test pass, return error message and exit code is not equal zero. \033[0m"
  #  estatus=$?
  fi

  echo "exitcode:"$exitCode
  echo ""
}

echo "===========================Chainstorage cli Test start=========================="
echo ''

testCases

echo "===========================Chainstorage cli Test end============================="

echo "Test status code:"$estatus
exit $estatus
