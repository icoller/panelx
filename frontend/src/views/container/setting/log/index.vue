<template>
    <div>
        <el-drawer
            v-model="drawerVisiable"
            :destroy-on-close="true"
            :close-on-click-modal="false"
            @close="handleClose"
            size="30%"
        >
            <template #header>
                <DrawerHeader :header="$t('container.cutLog')" :back="handleClose" />
            </template>
            <el-form :model="form" ref="formRef" :rules="rules" v-loading="loading" label-position="top">
                <el-row type="flex" justify="center">
                    <el-col :span="22">
                        <el-form-item prop="logMaxSize" :label="$t('container.maxSize')">
                            <el-input v-model.number="form.logMaxSize">
                                <template #append>
                                    <el-select v-model="form.sizeUnit" style="width: 70px">
                                        <el-option label="B" value="B"></el-option>
                                        <el-option label="KB" value="KB"></el-option>
                                        <el-option label="MB" value="MB"></el-option>
                                        <el-option label="GB" value="GB"></el-option>
                                    </el-select>
                                </template>
                            </el-input>
                        </el-form-item>
                        <el-form-item prop="logMaxFile" :label="$t('container.maxFile')">
                            <el-input v-model.number="form.logMaxFile" />
                        </el-form-item>
                    </el-col>
                </el-row>
            </el-form>
            <template #footer>
                <span class="dialog-footer">
                    <el-button @click="handleClose">{{ $t('commons.button.cancel') }}</el-button>
                    <el-button :disabled="loading" type="primary" @click="onSave(formRef)">
                        {{ $t('commons.button.confirm') }}
                    </el-button>
                </span>
            </template>
        </el-drawer>

        <ConfirmDialog ref="confirmDialogRef" @confirm="onSubmitSave"></ConfirmDialog>
    </div>
</template>
<script lang="ts" setup>
import { reactive, ref } from 'vue';
import { Rules, checkNumberRange } from '@/global/form-rules';
import i18n from '@/lang';
import { MsgSuccess } from '@/utils/message';
import { FormInstance } from 'element-plus';
import { updateLogOption } from '@/api/modules/container';
import DrawerHeader from '@/components/drawer-header/index.vue';

const loading = ref();
const drawerVisiable = ref();
const confirmDialogRef = ref();
const formRef = ref();

interface DialogProps {
    logMaxSize: string;
    logMaxFile: number;
}

const form = reactive({
    logMaxSize: 10,
    logMaxFile: 3,
    sizeUnit: 'MB',
});
const rules = reactive({
    logMaxSize: [checkNumberRange(1, 1024000), Rules.number],
    logMaxFile: [checkNumberRange(1, 100), Rules.number],
});

const emit = defineEmits<{ (e: 'search'): void }>();

const acceptParams = (params: DialogProps): void => {
    form.logMaxFile = params.logMaxFile || 3;
    if (params.logMaxSize) {
        form.logMaxSize = loadSize(params.logMaxSize);
    } else {
        form.logMaxSize = 10;
        form.sizeUnit = 'MB';
    }
    drawerVisiable.value = true;
};

const onSave = async (formEl: FormInstance | undefined) => {
    if (!formEl) return;
    formEl.validate(async (valid) => {
        if (!valid) return;
        let params = {
            header: i18n.global.t('database.confChange'),
            operationInfo: i18n.global.t('database.restartNowHelper'),
            submitInputInfo: i18n.global.t('database.restartNow'),
        };
        confirmDialogRef.value!.acceptParams(params);
    });
};

const onSubmitSave = async () => {
    loading.value = true;
    await updateLogOption(form.logMaxSize + '', form.logMaxFile + '')
        .then(() => {
            loading.value = false;
            drawerVisiable.value = false;
            emit('search');
            MsgSuccess(i18n.global.t('commons.msg.operationSuccess'));
        })
        .catch(() => {
            loading.value = false;
        });
};

const loadSize = (value: string) => {
    if (value.indexOf('b') !== -1 || value.indexOf('B') !== -1) {
        form.sizeUnit = 'B';
        return Number(value.replaceAll('b', '').replaceAll('B', ''));
    }
    if (value.indexOf('k') !== -1 || value.indexOf('K') !== -1) {
        form.sizeUnit = 'KB';
        return Number(value.replaceAll('k', '').replaceAll('K', ''));
    }
    if (value.indexOf('m') !== -1 || value.indexOf('M') !== -1) {
        form.sizeUnit = 'MB';
        return Number(value.replaceAll('m', '').replaceAll('M', ''));
    }
    if (value.indexOf('g') !== -1 || value.indexOf('G') !== -1) {
        form.sizeUnit = 'GB';
        return Number(value.replaceAll('g', '').replaceAll('G', ''));
    }
};

const handleClose = () => {
    emit('search');
    drawerVisiable.value = false;
};

defineExpose({
    acceptParams,
});
</script>
<style scoped>
.help-ul {
    color: #8f959e;
}
</style>