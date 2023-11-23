#!/bin/bash

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