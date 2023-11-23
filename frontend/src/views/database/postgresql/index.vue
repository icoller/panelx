<template>
    <div>
        <LayoutContent v-loading="loading" :divider="true"></LayoutContent>
    </div>
</template>

<script lang="ts" setup>
import { ref } from 'vue';
const loading = ref();
</script>

<style lang="scss" scoped></style>
