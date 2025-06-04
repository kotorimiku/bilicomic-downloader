<template>
  <div class="p-6 max-w-7xl mx-auto bg-gray-50 rounded-xl">
    <!-- 书籍 ID 输入 -->
    <div class="flex flex-wrap gap-3 mb-4">
      <n-input
        v-model:value="bookId"
        placeholder="请输入书籍 ID"
        class="flex-1"
        @keyup.enter="getChapterList"
      />
      <n-button type="primary" @click="getChapterList">搜索</n-button>
    </div>

    <!-- 操作按钮 -->
    <div class="flex flex-wrap justify-between items-center mb-4 gap-3">
      <div class="flex gap-3">
        <n-button @click="selectAll">全选</n-button>
        <n-button @click="selectInverse">反选</n-button>
        <n-button
          type="success"
          :disabled="!chapterList.length"
          @click="download"
        >
          开始下载
        </n-button>
      </div>
    </div>

    <!-- 章节列表 -->
    <div class="flex flex-col gap-4">
      <chapter-list
        :chapter-list="chapterList"
        v-model:selected-chapters="selectedChapters"
        :title="bookInfo?.Title"
        :cover="bookInfo?.Cover"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import {
  GetDownloader,
  GetBookInfo,
  GetChapter,
  DownloadList,
} from "../../wailsjs/go/bilicomicdownloader/DownloaderManager";
import { bilicomicdownloader as model } from "../../wailsjs/go/models";
import ChapterList from "../components/ChapterList.vue";
import { useRunCommand } from "../composables/RunCommand";
import { useNotify } from "../composables/useNotification";
import { EventsOn } from "../../wailsjs/runtime/runtime";

const runCommand = useRunCommand();

let notify = useNotify() // 延后初始化

EventsOn("message", (message: string) => {
    console.log(message);
    notify.info({
      content: message,
    });
  });

// 数据绑定
const selectedChapters = ref<number[]>([]);
const bookInfo = ref<model.BookInfo | null>(null);
const chapterList = ref<model.Volume[]>([]);
const isDownloading = ref(false); // 是否正在下载
const searchPage = ref<number>(1); // 搜索页码
const keyword = ref<string>(""); // 搜索关键字
const bookId = ref<string>(""); // 书籍 ID
// const searchResult = ref<model.Comic[]>([]); // 搜索结果
const isLoading = ref(false); // 是否正在加载

// 全选
const selectAll = () => {
  selectedChapters.value = chapterList.value.map((_, index) => index);
};

// 反选
const selectInverse = () => {
  chapterList.value.forEach((_, index) => {
    const chapterIndex = selectedChapters.value.indexOf(index);
    if (chapterIndex > -1) {
      selectedChapters.value.splice(chapterIndex, 1);
    } else {
      selectedChapters.value.push(index);
    }
  });
};

// 获取章节列表
const getChapterList = async () => {
  runCommand({
    command: () => GetDownloader(bookId.value),
    onSuccess: () => {
      runCommand({
        command: GetBookInfo,
        onSuccess: (res) => {
          bookInfo.value = res;
        },
      });
      runCommand({
        command: GetChapter,
        onSuccess: (res) => {
          selectedChapters.value = [];
          chapterList.value = res;
        },
      });
    },
  });
};

function debounce(func: Function, delay: number) {
  let timeoutId: any;
  return (...args: any) => {
    clearTimeout(timeoutId); // 清除之前的定时器
    timeoutId = setTimeout(() => {
      func.apply(args); // 延迟执行
    }, delay);
  };
}

// 下载选中卷
const download = debounce(async () => {
  if (selectedChapters.value.length === 0 || isDownloading.value) {
    return; // 没有选中卷，直接返回
  }
  isDownloading.value = true;
  runCommand({
    command: () => DownloadList(selectedChapters.value),
    onSuccess: () => {
      selectedChapters.value = [];
    },
  });
  isDownloading.value = false;
}, 200);
</script>
