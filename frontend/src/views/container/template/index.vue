<template>
    <div>
        <Submenu activeName="template" />
        <el-card style="margin-top: 20px">
            <ComplexTable :pagination-config="paginationConfig" v-model:selects="selects" :data="data" @search="search">
                <template #toolbar>
                    <el-button type="primary" @click="onOpenDialog('create')">
                        {{ $t('commons.button.create') }}
                    </el-button>
                    <el-button type="danger" plain :disabled="selects.length === 0" @click="onBatchDelete(null)">
                        {{ $t('commons.button.delete') }}
                    </el-button>
                </template>
                <el-table-column type="selection" fix></el-table-column>
                <el-table-column
                    :label="$t('commons.table.name')"
                    show-overflow-tooltip
                    min-width="100"
                    prop="name"
                    fix
                >
                    <template #default="{ row }">
                        <el-link @click="onOpenDetail(row)" type="primary">{{ row.name }}</el-link>
                    </template>
                </el-table-column>
                <el-table-column :label="$t('container.description')" prop="description" min-width="200" fix />
                <el-table-column :label="$t('commons.table.createdAt')" min-width="80" fix>
                    <template #default="{ row }">
                        {{ dateFromat(0, 0, row.createdAt) }}
                    </template>
                </el-table-column>
                <fu-table-operations :buttons="buttons" :label="$t('commons.table.operate')" />
            </ComplexTable>
        </el-card>

        <el-dialog v-model="detailVisiable" :destroy-on-close="true" :close-on-click-modal="false" width="70%">
            <template #header>
                <div class="card-header">
                    <span>{{ $t('commons.button.view') }}</span>
                </div>
            </template>
            <codemirror
                :autofocus="true"
                placeholder="None data"
                :indent-with-tab="true"
                :tabSize="4"
                style="max-height: 500px"
                :lineWrapping="true"
                :matchBrackets="true"
                theme="cobalt"
                :styleActiveLine="true"
                :extensions="extensions"
                v-model="detailInfo"
                :readOnly="true"
            />
            <template #footer>
                <span class="dialog-footer">
                    <el-button @click="detailVisiable = false">{{ $t('commons.button.cancel') }}</el-button>
                </span>
            </template>
        </el-dialog>

        <OperatorDialog @search="search" ref="dialogRef" />
    </div>
</template>

<script lang="ts" setup>
import ComplexTable from '@/components/complex-table/index.vue';
import Submenu from '@/views/container/index.vue';
import { Codemirror } from 'vue-codemirror';
import { javascript } from '@codemirror/lang-javascript';
import { oneDark } from '@codemirror/theme-one-dark';
import { reactive, onMounted, ref } from 'vue';
import { dateFromat } from '@/utils/util';
import { Container } from '@/api/interface/container';
import OperatorDialog from '@/views/container/template/operator/index.vue';
import { deleteComposeTemplate, searchComposeTemplate } from '@/api/modules/container';
import { useDeleteData } from '@/hooks/use-delete-data';
import i18n from '@/lang';
import { LoadFile } from '@/api/modules/files';

const data = ref();
const selects = ref<any>([]);
const detailVisiable = ref(false);
const detailInfo = ref();
const extensions = [javascript(), oneDark];

const paginationConfig = reactive({
    page: 1,
    pageSize: 10,
    total: 0,
});

const search = async () => {
    let params = {
        page: paginationConfig.page,
        pageSize: paginationConfig.pageSize,
    };
    await searchComposeTemplate(params).then((res) => {
        if (res.data) {
            data.value = res.data.items;
        }
    });
};

const onOpenDetail = async (row: Container.TemplateInfo) => {
    if (row.from === 'edit') {
        detailInfo.value = row.content;
        detailVisiable.value = true;
    } else {
        const res = await LoadFile({ path: row.path });
        detailInfo.value = res.data;
        detailVisiable.value = true;
    }
};

const dialogRef = ref();
const onOpenDialog = async (
    title: string,
    rowData: Partial<Container.TemplateInfo> = {
        name: '',
        from: 'edit',
        description: '',
        path: '',
        content: '',
    },
) => {
    let params = {
        title,
        rowData: { ...rowData },
    };
    dialogRef.value!.acceptParams(params);
};

const onBatchDelete = async (row: Container.RepoInfo | null) => {
    let ids: Array<number> = [];
    if (row) {
        ids.push(row.id);
    } else {
        selects.value.forEach((item: Container.RepoInfo) => {
            ids.push(item.id);
        });
    }
    await useDeleteData(deleteComposeTemplate, { ids: ids }, 'commons.msg.delete', true);
    search();
};

const buttons = [
    {
        label: i18n.global.t('commons.button.edit'),
        disabled: (row: Container.RepoInfo) => {
            return row.downloadUrl === 'docker.io';
        },
        click: (row: Container.RepoInfo) => {
            onOpenDialog('edit', row);
        },
    },
    {
        label: i18n.global.t('commons.button.delete'),
        disabled: (row: Container.RepoInfo) => {
            return row.downloadUrl === 'docker.io';
        },
        click: (row: Container.RepoInfo) => {
            onBatchDelete(row);
        },
    },
];

onMounted(() => {
    search();
});
</script>