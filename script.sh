#!/bin/bash

command -v wget >/dev/null || { 
  echo "wget not found, please install it and try again ."
  exit 1
}

if [ ! -f "pxctl" ]; then 
  wget https://github.com/icoller/installer/raw/main/pxctl
fi

if [ ! -f "panelx.service" ]; then 
  wget https://github.com/icoller/installer/raw/main/panelx.service
fi

if [ ! -f "install.sh" ]; then 
  wget https://github.com/icoller/installer/raw/main/install.sh
fi

chmod 755 pxctl install.sh