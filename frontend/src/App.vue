<template>
  <n-notification-provider>
    <n-space class="my-4">
      <n-button @click="changePage(Search)">搜索页面</n-button>
      <n-button @click="changePage(Progress)">进度</n-button>
      <n-button @click="changePage(Config)">配置页面</n-button>
    </n-space>
    <KeepAlive :exclude="['Progress']">
      <component :is="currentPage" />
    </KeepAlive>
  </n-notification-provider>
</template>

<script setup lang="ts">
import { shallowRef } from "vue";
import Search from "./views/Search.vue";
import Config from "./views/Config.vue";
import Progress from "./views/Progress.vue";

import { NNotificationProvider, NButton, NSpace } from "naive-ui";

const currentPage = shallowRef(Search);
const changePage = (page: any) => {
  currentPage.value = page;
};
</script>
