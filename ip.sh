#!/bin/bash


gip=`curl inet-ip.info`

echo $gip

# Domain : 4d2a.server-on.net
id=mydns984769
pwd=theXaS6kBiu
url="curl www.mydns.jp/directip.html?MID=${id}&PWD=${pwd}&IPV4ADDR=${gip}"

echo $url
curl $url

