<template>
  <div v-if="chapterList.length" class="flex flex-col gap-5 p-6 bg-white border border-gray-200 rounded-xl shadow-sm">
    <!-- 书籍信息 -->
    <div class="flex gap-5 items-start">
      <img v-if="cover" :src="cover" alt="Book Cover" class="w-36 h-auto rounded-lg shadow-md" />
      <div class="flex-1">
        <h2 class="text-2xl font-semibold text-gray-800">{{ title }}</h2>
      </div>
    </div>

    <!-- 章节列表 -->
    <div class="flex flex-col gap-3 border border-gray-100 rounded-lg bg-gray-50 p-4 overflow-y-auto h-[60vh]">
      <div class="flex flex-wrap gap-3">
        <div
          v-for="(chapter, index) in chapterList"
          :key="index"
          class="w-[28%] min-w-[160px] bg-white border border-gray-200 rounded-md px-3 py-2 shadow-sm cursor-pointer hover:bg-gray-100 transition-colors"
          @mousedown="startSelection(index)"
          @mouseover="handleMouseOver(index)"
          @click="toggleChapter(index)"
        >
          <n-checkbox
          class="h-full w-full"
            :checked="selectedChapters.includes(index)"
          >
            <span class="truncate block max-w-full">{{ chapter.Title }}</span>
          </n-checkbox>
        </div>
      </div>
    </div>
  </div>
</template>


<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { NCheckbox } from 'naive-ui'
import { bilicomicdownloader as model } from '../../wailsjs/go/models'

let lastSelectedIndex = ref<number | null>(null)
let isShiftPressed = ref(false)
let isMouseDown = ref(false)
let startIndex = ref<number | null>(null)

const props = defineProps({
  chapterList: {
    type: Array as () => model.Volume[],
    required: true
  },
  title: String,
  cover: String
})

const selectedChapters = defineModel('selectedChapters', {
  type: Array as () => number[],
  required: true
})

// 键盘 shift 状态
const handleKeydown = (e: KeyboardEvent) => {
  if (e.key === 'Shift') isShiftPressed.value = true
}
const handleKeyup = (e: KeyboardEvent) => {
  if (e.key === 'Shift') isShiftPressed.value = false
}

// 开始选择
const startSelection = (index: number) => {
  isMouseDown.value = true
  startIndex.value = index
  lastSelectedIndex.value = index
}

// 鼠标移动多选
const handleMouseOver = (index: number) => {
  if (isMouseDown.value && startIndex.value !== null) {
    const start = Math.min(startIndex.value, index)
    const end = Math.max(startIndex.value, index)
    selectedChapters.value = Array.from({ length: end - start + 1 }, (_, i) => start + i)
  }
}

// 点击选择/取消
const toggleChapter = (index: number) => {
  if (isShiftPressed.value && lastSelectedIndex.value !== null) {
    const start = Math.min(lastSelectedIndex.value, index)
    const end = Math.max(lastSelectedIndex.value, index)
    selectedChapters.value = Array.from({ length: end - start + 1 }, (_, i) => start + i)
  } else {
    const idx = selectedChapters.value.indexOf(index)
    if (idx > -1) selectedChapters.value.splice(idx, 1)
    else selectedChapters.value.push(index)
    lastSelectedIndex.value = index
  }
}

// 鼠标释放取消标志
const handleMouseUp = () => {
  isMouseDown.value = false
  startIndex.value = null
}

onMounted(() => {
  window.addEventListener('keydown', handleKeydown)
  window.addEventListener('keyup', handleKeyup)
  window.addEventListener('mouseup', handleMouseUp)
})
onUnmounted(() => {
  window.removeEventListener('keydown', handleKeydown)
  window.removeEventListener('keyup', handleKeyup)
  window.removeEventListener('mouseup', handleMouseUp)
})
</script>