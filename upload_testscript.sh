#!/bin/bash

estatus=0

cmdPath='./gcscmd'

testCases() {

  echo "_/_/_/_/_/_/_/_/_/_/_/_/_/ 对象上传 开始 _/_/_/_/_/_/_/_/_/_/_/_/_/"
  echo '对象上传 => 上传文件-当前目录 => 在当前目录上传文件 => gcscmd put ./aaa.mp4 cs://bbb start'
  echo '数据准备'
  echo '1、创建新桶'
  bucketName="bucket-upload-"$(date "+%Y%m%d%H%M%S")"-"$RANDOM
  $cmdPath mb cs://$bucketName
  echo '2、添加对象，10MB'
  testDataFileName="testdata_"$(date "+%Y%m%d%H%M%S")"-"$RANDOM".dat"
  #echo $testDataFileName
  dd if=/dev/urandom of=$testDataFileName bs=10240 count=1024
  # 设置当前目录
  dataPath=''$testDataFileName
  execCmd '对象上传' '上传文件-当前目录' '在当前目录上传文件' 'gcscmd put '$dataPath' cs://'$bucketName 'put '$dataPath' cs://'$bucketName ''
  echo '3、数据清理'
  rm -rf $testDataFileName
  echo '对象上传 => 上传文件-当前目录 => 在当前目录上传文件 => gcscmd put ./aaa.mp4 cs://bbb end'
  echo ''

  echo '对象上传 => 上传文件-当前目录 => 绝对路径上传文件 => gcscmd put /home/pz/aaa.mp4 cs://bbb start'
  echo '数据准备'
  echo '1、添加对象，10MB'
  testDataFileName="testdata_"$(date "+%Y%m%d%H%M%S")"-"$RANDOM".dat"
  #echo $testDataFileName
  # 设置绝对路径
  dataPath=$(pwd)'/'$testDataFileName
  dd if=/dev/urandom of=$dataPath bs=10240 count=1024
  execCmd '对象上传' '上传文件-当前目录' '绝对路径上传文件' 'gcscmd put '$dataPath' cs://'$bucketName 'put '$dataPath' cs://'$bucketName ''
  echo '2、数据清理'
  rm -rf $dataPath
  echo '对象上传 => 上传文件-当前目录 => 绝对路径上传文件 => gcscmd put /home/pz/aaa.mp4 cs://bbb end'
  echo ''

  echo '对象上传 => 上传文件-当前目录 => 相对路径上传文件 => gcscmd put ../pz/aaa.mp4 cs://bbb start'
  echo '数据准备'
  echo '1、添加对象，10MB'
  testDataFileName="testdata_"$(date "+%Y%m%d%H%M%S")"-"$RANDOM".dat"
  #echo $testDataFileName
  # 设置相对路径
  dataPath='tmp/'$testDataFileName
  dd if=/dev/urandom of=$dataPath bs=10240 count=1024
  execCmd '对象上传' '上传文件-当前目录' '相对路径上传文件' 'gcscmd put '$dataPath' cs://'$bucketName 'put '$dataPath' cs://'$bucketName ''
  echo '2、数据清理'
  rm -rf $dataPath
  echo '对象上传 => 上传文件-当前目录 => 相对路径上传文件 => gcscmd put ../pz/aaa.mp4 cs://bbb end'
  echo ''

  echo '对象上传 => 上传文件-当前目录 => 错误上传-任意方式上传到不存在的桶 => gcscmd put ./aaa.mp4 cs://不存在的桶名 start'
  echo '数据准备'
  echo '1、创建新桶'
  wrongBucketName="bucket-upload-"$(date "+%Y%m%d%H%M%S")"-"$RANDOM
  echo '2、添加对象，10MB'
  testDataFileName="testdata_"$(date "+%Y%m%d%H%M%S")"-"$RANDOM".dat"
  #echo $testDataFileName
  dd if=/dev/urandom of=$testDataFileName bs=10240 count=1024
  # 设置当前目录
  dataPath=''$testDataFileName
  execCmdFail '对象上传' '上传文件-当前目录' '错误上传-任意方式上传到不存在的桶' 'gcscmd put '$dataPath' cs://'$wrongBucketName 'put '$dataPath' cs://'$wrongBucketName ''
  echo '3、数据清理'
  rm -rf $testDataFileName
  echo '对象上传 => 上传文件-当前目录 => 错误上传-任意方式上传到不存在的桶 => gcscmd put ./aaa.mp4 cs://不存在的桶名 end'
  echo ''

  # 不考虑空目录的情况，若上传空目录会报错：Uploading folder is empty, or uploading data is invalid in the folder
  #  echo '对象上传 => 上传目录 => 正确上传目录-空目录上传 => gcscmd put ./aaaa cs://bbb start'
  #  echo '数据准备'
  #  # 1、添加目录
  #  testDataFolderame="testfolder_"$(date "+%Y%m%d%H%M%S")"-"$RANDOM
  #  #echo $testDataFolderame
  #  mkdir $testDataFolderame
  #  $cmdPath put $testDataFolderame cs://$bucketName
  #  # 设置对象名称
  #  dataPath=''$testDataFolderame
  #  execCmdFail '对象上传' '上传文件-当前目录' '正确上传目录-空目录上传' 'gcscmd put '$dataPath' cs://'$bucketName 'put '$dataPath' cs://'$bucketName ''
  #  echo '2、数据清理'
  #  rm -rf $dataPath
  #  echo ''

  echo '对象上传 => 上传目录 => 正确上传目录-目录有文件上传 => gcscmd put ./aaaa cs://bbb start'
  echo '数据准备'
  # 1、添加目录
  testDataFolderame="testfolder_"$(date "+%Y%m%d%H%M%S")"-"$RANDOM
  #echo $testDataFolderame
  mkdir $testDataFolderame
  $cmdPath put $testDataFolderame cs://$bucketName
  echo '2、添加对象，10MB'
  testDataFileName="testdata_"$(date "+%Y%m%d%H%M%S")"-"$RANDOM".dat"
  #echo $testDataFileName
  # 设置相对路径
  dataPath=$testDataFolderame'/'$testDataFileName
  dd if=/dev/urandom of=$dataPath bs=10240 count=1024
  execCmd '对象上传' '上传文件-当前目录' '正确上传目录-目录有文件上传' 'gcscmd put '$testDataFolderame' cs://'$bucketName 'put '$testDataFolderame' cs://'$bucketName ''
  echo '3、数据清理'
  rm -rf $testDataFolderame
  echo ''

  # 目前暂时不使用--carfile参数
  #  # 对象上传 => 上传carfile => 正确上传car文件 => gcscmd put ./aaa.car cs://bbb --carfile
  #  echo '数据准备'
  #  echo '1、添加对象，10MB'
  #  testDataFileName="testdata_"$(date "+%Y%m%d%H%M%S")"-"$RANDOM".dat"
  #  #echo $testDataFileName
  #  dd if=/dev/urandom of=$testDataFileName bs=10240 count=1024
  #  # 2、生成CAR文件(测试环境需要安装car命令行程序)
  #  carFilename=$(echo ${testDataFileName/.dat/.car})
  #  car c --version 1 -f $carFilename $testDataFileName
  #  # 设置当前目录
  #  dataPath=''$carFilename
  #  execCmd '对象上传' '上传carfile' '正确上传car文件' 'gcscmd put '$dataPath' cs://'$bucketName' --carfile' 'put '$dataPath' cs://'$bucketName' --carfile' ''
  #  #  echo '3、数据清理'
  #  #  rm -rf $testDataFileName
  #  #  rm -rf $dataPath
  #
  #  # 对象上传 => 上传carfile => 重复上传-上传已经存在的car文件 => gcscmd put ./aaa.car cs://bbb --carfile
  #  execCmd '对象上传' '上传carfile' '重复上传-上传已经存在的car文件' 'gcscmd put '$dataPath' cs://'$bucketName' --carfile' 'put '$dataPath' cs://'$bucketName' --carfile' ''

  # 对象上传测试数据清理
  $cmdPath rb cs://$bucketName --force
  #  rm -rf $testDataFileName
  echo "_/_/_/_/_/_/_/_/_/_/_/_/_/ 对象上传 结束 _/_/_/_/_/_/_/_/_/_/_/_/_/"
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
