<template>
  <div>
    <!-- 配置切换按钮 -->
    <div class="switch-container">
      <button @click="changePage(Search)" class="btn">搜索页面</button>
      <button @click="changePage(Progress)" class="btn">进度</button>
      <button @click="changePage(Config)" class="btn">配置页面</button>
    </div>
    <KeepAlive :exclude="['Progress']">
      <component :is="currentPage" />
    </KeepAlive>
  </div>
</template>

<script setup lang="ts">
import { shallowRef } from 'vue';
import Search from './views/Search.vue';
import Config from './views/Config.vue';
import Progress from './views/Progress.vue';

const currentPage = shallowRef(Search); // 存储当前组件

// 切换页面
const changePage = (page: any) => {
  currentPage.value = page; // 直接切换组件
};
</script>

<style scoped>
/* 按钮样式 */
.btn {
  padding: 8px 16px;
  font-size: 14px;
  background-color: #fff;
  border: 1px solid #ddd;
  border-radius: 5px;
  cursor: pointer;
  transition: background-color 0.3s ease;
}

.btn:hover {
  background-color: #f1f1f1;
}
</style>
