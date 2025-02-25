<template>
  <div class="container">
    <!-- 遍历进度数据 -->
    <div v-for="item in progressData" :key="item.Volume?.Title" class="progress-item">
      <p class="progress-title">{{ item.BookInfo?.Title + "/" + item.Volume?.Title }}</p>
      <ProgressBar :progress="item.Progress" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import ProgressBar from '../components/ProgressBar.vue';
import { EventsOn } from '../../wailsjs/runtime/runtime';
import { bilicomicdownloader as model } from '../../wailsjs/go/models';
import { GetDownloaders } from '../../wailsjs/go/bilicomicdownloader/DownloaderManager';

const progressData = ref<model.DownloaderSingle[]>([]);

// 监听进度事件
EventsOn('progress', (data: model.DownloaderSingle[]) => {
  progressData.value = data;
  console.log(data);

});

// 页面加载时获取下载器数据
onMounted(async () => {
  progressData.value = await GetDownloaders();
});
</script>

<style scoped>
.container {
  padding: 10px;
  font-family: 'Arial', sans-serif;
}

.progress-item {
  margin-bottom: 5px;
  padding: 16px;
  border: 1px solid #ddd;
  border-radius: 8px;
  background-color: #f9f9f9;
}

.progress-title {
  font-size: 16px;
  font-weight: bold;
  color: #333;
  margin-bottom: 8px;
  margin-top: 5px;
}
</style>
