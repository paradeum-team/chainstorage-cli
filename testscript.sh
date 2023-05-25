#!/bin/bash

estatus=0

testCases() {

#  # 帮助 => 帮助功能 => gcscmd -h [--help]
#  execCmd '帮助' '帮助功能' '' 'gcscmd -h [--help]' '-h' ''
#
#  # 帮助 => Basic command => gcscmd ls --help
#  execCmd '帮助' 'Basic command' '' 'gcscmd ls --help' 'ls --help' ''
#
#  # 帮助 => Basic command => gcscmd mb --help
#  execCmd '帮助' 'Basic command' '' 'gcscmd mb --help' 'mb --help' ''
#
#  # 帮助 => Basic command => gcscmd rb --help
#  execCmd '帮助' 'Basic command' '' 'gcscmd rb --help' 'rb --help' ''
#
#  # 帮助 => Basic command => gcscmd get --help
#  execCmd '帮助' 'Basic command' '' 'gcscmd get --help' 'get --help' ''
#
#  # 帮助 => Basic command => gcscmd put --help
#  execCmd '帮助' 'Basic command' '' 'gcscmd put --help' 'put --help' ''
#
#  # 帮助 => Basic command => gcscmd rm --help
#  execCmd '帮助' 'Basic command' '' 'gcscmd rm --help' 'rm --help' ''
#
#  # 帮助 => Basic command => gcscmd rn --help
#  execCmd '帮助' 'Basic command' '' 'gcscmd rn --help' 'rn --help' ''
#
#  # 帮助 => Basic command => gcscmd import --help
#  execCmd '帮助' 'Basic command' '' 'gcscmd import --help' 'import --help' ''
#
#  # 帮助 => Tool Command => config => gcscmd config
#  execCmd '帮助' 'Basic command' '' 'gcscmd config --help' 'config --help' ''
#
#  # 帮助 => Tool Command => version => gcscmd version
#  execCmd '帮助' 'Basic command' '' 'gcscmd version --help' 'version --help' ''
#
#  # 帮助 => Tool Command => log => gcscmd log
#  execCmd '帮助' 'Basic command' '' 'gcscmd log --help' 'log --help' ''

  # todo: 数据准备，删除所有桶，还是切换APIKEY？
  # 桶操作 => 列出桶对象 => 无桶 => gcscmd ls
  execCmd '桶操作' '列出桶对象' '无桶' 'gcscmd ls' 'ls' ''

  # 桶操作 => 创建桶 => 正常创建 => gcscmd mb cs://bbb
  bucketName="bucket"$(date "+%Y%m%d%H%M%S")
  execCmd '桶操作' '创建桶' '正常创建' 'gcscmd mb cs://'$bucketName 'mb cs://'$bucketName ''

  # 桶操作 => 创建桶 => 非正常创建-桶名重复 => gcscmd mb cs://bbb
  execCmd '桶操作' '创建桶' '非正常创建-桶名重复' 'gcscmd mb cs://'$bucketName 'mb cs://'$bucketName ''

  # todo: 测试意图？
  # 桶操作 => 创建桶 => 非正常创建-创建多个桶 => gcscmd mb cs://aaa
  execCmd '桶操作' '创建桶' '非正常创建-创建多个桶' 'gcscmd mb cs://'$bucketName 'mb cs://'$bucketName ''

  # 桶操作 => 列出桶对象 => 有桶 => gcscmd ls
  execCmd '桶操作' '列出桶对象' '无桶' 'gcscmd ls' 'ls' ''

  # 桶操作 => 移除桶 => 正常删除-无数据删除 => gcscmd rb cs://bbb
  execCmd '桶操作' '移除桶' '正常删除-无数据删除' 'gcscmd rb cs://'$bucketName 'rb cs://'$bucketName ''

  # 桶操作 => 移除桶 => 重复删除-继续删除已删除的桶 => gcscmd rb cs://bbb
  execCmd '桶操作' '移除桶' '重复删除-继续删除已删除的桶' 'gcscmd rb cs://'$bucketName 'rb cs://'$bucketName ''

  # 桶操作 => 移除桶 => 正常删除-有数据删除 => gcscmd rb cs://bbb
  # todo:
  # 1、创建新桶
  # 2、添加对象
  execCmd '桶操作' '移除桶' '正常删除-有数据删除' 'gcscmd rb cs://'$bucketName 'rb cs://'$bucketName ''

  # 桶操作 => 移除桶 => 正常删除-有数据强制删除 => gcscmd rb cs://bbb --force
  execCmd '桶操作' '移除桶' '正常删除-有数据强制删除' 'gcscmd rb cs://'$bucketName' --force' 'rb cs://'$bucketName' --force' ''

  # 桶操作 => 清空桶 => 正常清空-有数据清空 => gcscmd rm cs://bbb
  # todo:
  # 1、创建新桶
  # 2、添加对象
  execCmd '桶操作' '清空桶' '正常清空-有数据清空' 'gcscmd rm cs://'$bucketName 'rm cs://'$bucketName ''

  # 对象操作 => 查看对象 => 桶内所有文件查询-有此桶 => gcscmd ls cs://bbb
  execCmd '对象操作' '查看对象' '桶内所有文件查询-有此桶' 'gcscmd ls cs://'$bucketName 'ls cs://'$bucketName ''

  # 对象操作 => 查看对象 => 桶内所有文件查询-无此桶 => gcscmd ls cs://bbb
  execCmd '对象操作' '查看对象' '桶内所有文件查询-无此桶' 'gcscmd ls cs://'$bucketName 'ls cs://'$bucketName ''

  # 对象操作 => 查看对象 => 桶内对应 cid 查询-cid正确 => gcscmd ls cs://bbb --cid QmWgnG7pPjG31w328hZyALQ2BgW5aQrZyKpT47jVpn8CNo
  cid=''
  execCmd '对象操作' '查看对象' '桶内对应 cid 查询-cid正确' 'gcscmd ls cs://'$bucketName' --cid '$cid 'ls cs://'$bucketName' --cid '$cid ''

   # 对象操作 => 查看对象 => 桶内对象名查询-对象名正确 => gcscmd ls cs://bbb --name Tarkov.mp4
  objectName=''
  execCmd '对象操作' '查看对象' '桶内对象名查询-对象名正确' 'gcscmd ls cs://'$bucketName' --name '$objectName 'ls cs://'$bucketName' --name '$objectName ''

  # 对象上传 => 上传文件-当前目录 => 在当前目录上传文件 => gcscmd put ./aaa.mp4 cs://bbb
  dataPath=''
  execCmd '对象上传' '上传文件-当前目录' '在当前目录上传文件' 'gcscmd put '$dataPath' cs://'$bucketName 'put '$dataPath' cs://'$bucketName ''

  # 对象上传 => 上传文件-当前目录 => 绝对路径上传文件 => gcscmd put /home/pz/aaa.mp4 cs://bbb
  execCmd '对象上传' '上传文件-当前目录' '绝对路径上传文件' 'gcscmd put '$dataPath' cs://'$bucketName 'put '$dataPath' cs://'$bucketName ''

  # 对象上传 => 上传文件-当前目录 => 相对路径上传文件 => gcscmd put ../pz/aaa.mp4 cs://bbb
  execCmd '对象上传' '上传文件-当前目录' '相对路径上传文件' 'gcscmd put '$dataPath' cs://'$bucketName 'put '$dataPath' cs://'$bucketName ''

  # 对象上传 => 上传文件-当前目录 => 错误上传-任意方式上传到不存在的桶 => gcscmd put ./aaa.mp4 cs://不存在的桶名
  execCmd '对象上传' '上传文件-当前目录' '错误上传-任意方式上传到不存在的桶' 'gcscmd put '$dataPath' cs://'$bucketName 'put '$dataPath' cs://'$bucketName ''

  # 对象上传 => 上传目录 => 正确上传目录-空目录上传 => gcscmd put ./aaaa cs://bbb
  execCmd '对象上传' '上传文件-当前目录' '正确上传目录-空目录上传' 'gcscmd put '$dataPath' cs://'$bucketName 'put '$dataPath' cs://'$bucketName ''

  # 对象上传 => 上传目录 => 正确上传目录-目录有文件上传 => gcscmd put ./aaaa cs://bbb
  execCmd '对象上传' '上传文件-当前目录' '正确上传目录-目录有文件上传' 'gcscmd put '$dataPath' cs://'$bucketName 'put '$dataPath' cs://'$bucketName ''

  # 对象上传 => 上传carfile => 正确上传car文件 => gcscmd put ./aaa.car cs://bbb --carfile
  execCmd '对象上传' '上传carfile' '正确上传car文件' 'gcscmd put '$dataPath' cs://'$bucketName' --carfile' 'put '$dataPath' cs://'$bucketName' --carfile' ''

  # 对象上传 => 上传carfile => 重复上传-上传已经存在的car文件 => gcscmd put ./aaa.car cs://bbb --carfile
  execCmd '对象上传' '上传carfile' '重复上传-上传已经存在的car文件' 'gcscmd put '$dataPath' cs://'$bucketName' --carfile' 'put '$dataPath' cs://'$bucketName' --carfile' ''

  # 导入 car 文件 => 正确导入car文件 => 当前目录导入 => gcscmd import ./aaa.car cs://bbb
  execCmd '导入 car 文件' '正确导入car文件' '当前目录导入' 'gcscmd import '$dataPath' cs://'$bucketName 'import '$dataPath' cs://'$bucketName ''

  # 导入 car 文件 => 正确导入car文件 => 绝对路径导入 => gcscmd import /home/pz/aaa.car cs://bbb
  execCmd '导入 car 文件' '正确导入car文件' '绝对路径导入' 'gcscmd import '$dataPath' cs://'$bucketName 'import '$dataPath' cs://'$bucketName ''

  # 导入 car 文件 => 正确导入car文件 => 相对路径 => gcscmd import ../pz/aaa.car cs://bbb
  execCmd '导入 car 文件' '正确导入car文件' '相对路径' 'gcscmd import '$dataPath' cs://'$bucketName 'import '$dataPath' cs://'$bucketName ''

  # 下载对象 => 根据cid下载 => cid正确下载 => gcscmd get cs://bbb --cid QmWgnG7pPjG31w328hZyALQ2BgW5aQrZyKpT47jVpn8CNo
  execCmd '下载对象' '根据cid下载' 'cid正确下载' 'gcscmd get cs://'$bucketName' --cid '$cid 'get cs://'$bucketName' --cid '$cid ''

  # 下载对象 => 根据对象名下载 => 对象名正确下载 => gcscmd get cs://bbb --name Tarkov.mp4
  execCmd '下载对象' '根据对象名下载' '对象名正确下载' 'gcscmd get cs://'$bucketName' --name '$objectName 'get cs://'$bucketName' --name '$objectName ''

  # 删除对象 => 清空桶 => 有文件清空桶 => gcscmd rm cs://bbb --force
  execCmd '删除对象' '清空桶' '有文件清空桶' 'gcscmd rm cs://'$bucketName' --force' 'rm cs://'$bucketName' --force' ''

  # 删除对象 => 使用对象名删除单文件 => 正确删除 => gcscmd rm cs://bbb --name Tarkov.mp4
  execCmd '删除对象' '使用对象名删除单文件' '正确删除' 'gcscmd rm cs://'$bucketName' --name '$objectName 'rm cs://'$bucketName' --name '$objectName ''

  # 删除对象 => 使用模糊查询删除对象 => 正常模糊删除 => gcscmd rm cs://bbb --name .mp4 --force
  execCmd '删除对象' '使用模糊查询删除对象' '正常模糊删除' 'gcscmd rm cs://'$bucketName' --name '$objectName' --force' 'rm cs://'$bucketName' --name '$objectName' --force' ''

  # 删除对象 => 使用对象名删除单目录 => 正确删除-目录中无文件 => gcscmd rm cs://bbb --name aaa
  execCmd '删除对象' '使用对象名删除单目录' '正确删除-目录中无文件' 'gcscmd rm cs://'$bucketName' --name '$objectName 'rm cs://'$bucketName' --name '$objectName ''

  # 删除对象 => 使用CID删除单对象 => 正确删除-对应的桶中有此CID删除 => gcscmd rm cs://bbb --cid QmWgnG7pPjG31w328hZyALQ2BgW5aQrZyKpT47jVpn8CNo
  execCmd '删除对象' '使用CID删除单对象' '正确删除-对应的桶中有此CID删除' 'gcscmd rm cs://'$bucketName' --cid '$cid 'rm cs://'$bucketName' --cid '$cid ''

  # 删除对象 => 使用 CID 删除多个对象(命中多个对象时加) => 一个cid有多个对象 => gcscmd rm cs://bbb --cid QmWgnG7pPjG31w328hZyALQ2BgW5aQrZyKpT47jVpn8CNo --force
  execCmd '删除对象' '使用 CID 删除多个对象(命中多个对象时加)' '一个cid有多个对象' 'gcscmd rm cs://'$bucketName' --cid '$cid' --force' 'rm cs://'$bucketName' --cid '$cid' --force' ''

  # 重命名对象 => 使用对象名 => 替换的文件无冲突 => gcscmd rn cs://bbb --name Tarkov.mp4 --rename aaa.mp4
  rename=""
  execCmd '重命名对象' '使用对象名' '替换的文件无冲突' 'gcscmd rn cs://'$bucketName' --name '$objectName' --rename 'rename 'rn cs://'$bucketName' --name '$objectName' --rename 'rename ''

  # 重命名对象 => 使用对象名 => 替换的文件有冲突-有force， => gcscmd rn cs://bbb --name Tarkov.mp4 --rename aaa.mp4 --force
  execCmd '重命名对象' '使用对象名' '替换的文件有冲突-有force' 'gcscmd rn cs://'$bucketName' --name '$objectName' --rename 'rename' --force' 'rn cs://'$bucketName' --name '$objectName' --rename 'rename' --force' ''

  # 重命名对象 => 使用CID => cid单对象-无冲突 => gcscmd rn cs://bbb --cid QmWgnG7pPjG31w328hZyALQ2BgW5aQrZyKpT47jVpn8CNo --rename aaa.mp4
  execCmd '重命名对象' '使用CID' 'cid单对象-无冲突' 'gcscmd rn cs://'$bucketName' --cid '$cid' --rename 'rename 'rn cs://'$bucketName' --cid '$cid' --rename 'rename ''

  # 重命名对象 => 使用CID => cid多对象-无冲突（是否分此情况） =>
  execCmd '重命名对象' '使用CID' 'cid多对象-无冲突（是否分此情况）' 'gcscmd rn cs://'$bucketName' --cid '$cid' --rename 'rename 'rn cs://'$bucketName' --cid '$cid' --rename 'rename ''
}

execCmd() {
  testModule=$1
  testFunction=$2
  testCase=$3
  testDescription=$4
  testCmd=$5
  testExpectation=$6

  echo $testModule"=>"$testFunction"=>"$testCase"=>"$testDescription
  cmdStr='./gcscmd '$testCmd
  echo 'executing '$cmdStr
  eval $cmdStr
  if [ $? -eq 0 ]; then
    echo "Success: "$cmdStr" test pass."
  else
    echo "Failure: "$cmdStr" test fail."
  #  estatus=$?
  fi
  echo ""
}

echo "===========================Chainstorage cli Test start=========================="
testCases
echo "===========================Chainstorage cli Test end=========================="

exit $estatus
