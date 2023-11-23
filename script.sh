#!/bin/bash

command -v wget >/dev/null || { 
  echo "wget not found, please install it and try again ."
  exit 1
}

if [ ! -f "pxctl" ]; then 
  wget http://resources.panelx.cn/installer/pxctl
fi

if [ ! -f "panelx.service" ]; then 
  wget http://resources.panelx.cn/installer/panelx.service
fi

if [ ! -f "install.sh" ]; then 
  wget http://resources.panelx.cn/installer/install.sh
fi

chmod 755 pxctl install.sh