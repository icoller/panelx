<template>
    <el-dialog v-model="open" :title="$t('app.param')" width="30%" :close-on-click-modal="false">
        <el-descriptions border :column="1">
            <el-descriptions-item v-for="(param, key) in params" :label="param.label" :key="key">
                {{ param.value }}
            </el-descriptions-item>
        </el-descriptions>
    </el-dialog>
</template>
<script lang="ts" setup>
import { App } from '@/api/interface/app';
import { GetAppInstallParams } from '@/api/modules/app';
import { ref } from 'vue';

interface ParamProps {
    id: Number;
}
const paramData = ref<ParamProps>({
    id: 0,
});

let open = ref(false);
let loading = ref(false);
const params = ref<App.InstallParams[]>();

const acceptParams = (props: ParamProps) => {
    params.value = [];
    paramData.value.id = props.id;
    get();
    open.value = true;
};

const get = async () => {
    try {
        loading.value = true;
        const res = await GetAppInstallParams(Number(paramData.value.id));
        params.value = res.data;
    } catch (error) {
    } finally {
        loading.value = false;
    }
};

defineExpose({ acceptParams });
</script>