<template>
  <div class="p-4">
    <n-space vertical :size="16">
      <n-card
        v-for="item in progressData"
        :key="item.Volume?.Title"
        :class="item.Fail ? 'border-red-500' : ''"
        :style="{ backgroundColor: item.Fail ? '#fff1f0' : undefined }"
        embedded
      >
        <template #header>
          <span class="text-base font-bold">
            {{ item.BookInfo?.Title + ' / ' + item.Volume?.Title }}
            <span v-if="item.Fail" class="text-red-500">（失败）</span>
          </span>
        </template>
        <n-progress
          type="line"
          :percentage="item.Progress"
          :status="item.Fail ? 'error' : 'default'"
          :color="item.Fail ? '#ef4444' : '#18a058'"
          :height="12"
          indicator-placement="inside"
          processing
        />
      </n-card>
    </n-space>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { EventsOn } from '../../wailsjs/runtime/runtime'
import { bilicomicdownloader as model } from '../../wailsjs/go/models'
import { GetDownloaders } from '../../wailsjs/go/bilicomicdownloader/DownloaderManager'

const progressData = ref<model.DownloaderSingle[]>([])

EventsOn('progress', (data: model.DownloaderSingle[]) => {
  progressData.value = data
})

onMounted(async () => {
  progressData.value = await GetDownloaders()
})
</script>
