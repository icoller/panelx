#!/bin/bash
###
 # @Author: coller
 # @Date: 2023-12-09 18:04:24
 # @LastEditors: coller
 # @LastEditTime: 2023-12-09 18:05:05
 # @Desc:  
### 

command -v wget >/dev/null || { 
  echo "wget not found, please install it and try again ."
  exit 1
}

if [ ! -f "pxctl" ]; then 
  wget https://ft-resource.oss-cn-hangzhou.aliyuncs.com/installer/pxctl
fi

if [ ! -f "panelx.service" ]; then 
  wget https://ft-resource.oss-cn-hangzhou.aliyuncs.com/installer/panelx.service
fi

if [ ! -f "install.sh" ]; then 
  wget https://ft-resource.oss-cn-hangzhou.aliyuncs.com/installer/install.sh
fi

chmod 755 pxctl install.sh