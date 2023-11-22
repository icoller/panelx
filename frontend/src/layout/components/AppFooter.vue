<!--
 * @Author: coller
 * @Date: 2023-11-22 00:08:39
 * @LastEditors: coller
 * @LastEditTime: 2023-11-22 00:32:55
 * @Copyright (c) 2023 by www.kuaiqin.com, All Rights Reserved. 
 * @Desc:  
-->
<template>
    <div class="footer">
        <a href="https://www.panelx.cn/" target="_blank">Copyright © 2023 PanelX 面板探索</a>
        <SystemUpgrade :footer="true" />
    </div>
</template>

<script setup lang="ts">
import SystemUpgrade from '@/components/system-upgrade/index.vue';
</script>

<style scoped lang="scss">
.footer {
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: 48px;
    background: #ffffff;
    border-top: 1px solid #e4e7ed;
    box-sizing: border-box;
    padding: 10px 20px;
    a {
        font-size: 12px;
        color: #858585;
        text-decoration: none;
        letter-spacing: 0.5px;
    }
    span {
        font-size: 12px;
        color: #858585;
        text-decoration: none;
        letter-spacing: 0.5px;
    }
}
</style>
