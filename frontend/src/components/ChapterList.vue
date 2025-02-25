<template>
    <div class="book-container" v-if="chapterList.length">
        <!-- 书籍信息 -->
        <div class="book-info">
            <img :src="cover" alt="Book Cover" class="book-cover" v-if="cover" />
            <h2 class="book-title">{{ title, cover }}</h2>
        </div>

        <!-- 章节列表横向排列 -->
        <div class="switch-container">
            <div class="chapter-list">
                <div v-for="(chapter, index) in chapterList" :key="index" class="switch-item"
                     @mousedown="startSelection(index)"
                     @mouseover="handleMouseOver(index)">
                    <input type="checkbox" v-model="selectedChapters" :value="index" :id="'chapter-' + index" />
                    <span :for="'chapter-' + index">{{ chapter.Title }}</span>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
import { bilicomicdownloader as model } from '../../wailsjs/go/models';

let lastSelectedIndex = ref<number | null>(null);
let isShiftPressed = ref(false);
let isMouseDown = ref(false);
let startIndex = ref<number | null>(null);

const props = defineProps({
    chapterList: {
        type: Array as () => model.Volume[],
        required: true
    },
    title: String,
    cover: String
})

const selectedChapters = defineModel('selectedChapters', { type: Array as () => number[], required: true });

// 检测 Shift 键按下状态
const handleKeydown = (event: KeyboardEvent) => {
    if (event.key === "Shift") {
        isShiftPressed.value = true;
    }
};

const handleKeyup = (event: KeyboardEvent) => {
    if (event.key === "Shift") {
        isShiftPressed.value = false;
    }
};

// 开始选择
const startSelection = (index: number) => {
    isMouseDown.value = true;
    startIndex.value = index;
    lastSelectedIndex.value = index;
    togglechapterSelection(index);
};

// 处理鼠标移动事件
const handleMouseOver = (index: number) => {
    if (isMouseDown.value && startIndex.value !== null) {
        const start = Math.min(startIndex.value, index);
        const end = Math.max(startIndex.value, index);
        const range = Array.from({ length: end - start + 1 }, (_, i) => start + i);

        selectedChapters.value = range;
    }
};

// 选择章节
const togglechapterSelection = (index: number) => {
    if (isShiftPressed.value && lastSelectedIndex.value !== null) {
        const start = Math.min(lastSelectedIndex.value, index);
        const end = Math.max(lastSelectedIndex.value, index);
        const range = Array.from({ length: end - start + 1 }, (_, i) => start + i);

        selectedChapters.value = range;
    } else {
        // 单独选择
        const chapterIndex = selectedChapters.value.indexOf(index);
        if (chapterIndex > -1) {
            selectedChapters.value.splice(chapterIndex, 1);
        } else {
            selectedChapters.value.push(index);
        }
        lastSelectedIndex.value = index;
    }
};

// 处理鼠标释放事件
const handleMouseUp = () => {
    isMouseDown.value = false;
    startIndex.value = null;
};

onMounted(() => {
    window.addEventListener('keydown', handleKeydown);
    window.addEventListener('keyup', handleKeyup);
    window.addEventListener('mouseup', handleMouseUp);
});

onUnmounted(() => {
    window.removeEventListener('keydown', handleKeydown);
    window.removeEventListener('keyup', handleKeyup);
    window.removeEventListener('mouseup', handleMouseUp);
});
</script>

<style lang="css" scoped>
/* 书籍信息容器 */
.book-container {
    display: flex;
    flex-direction: column;
    gap: 20px;
    padding: 20px;
    background-color: #f9f9f9;
    border-radius: 8px;
    border: 1px solid #ddd;
}

/* 书籍封面 */
.book-cover {
    width: 150px;
    height: auto;
    border-radius: 8px;
    margin-bottom: 10px;
}

/* 书籍标题 */
.book-title {
    font-size: 24px;
    font-weight: bold;
    color: #333;
    margin: 0;
}

/* 章节列表容器 */
.switch-container {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 15px;
    font-size: 14px;
    padding: 10px;
    border: 1px solid #ddd;
    border-radius: 8px;
    background-color: #fff;
    height: 60vh;
    overflow-y: auto;
    user-select: none;
}

/* 每个章节框 */
.switch-item {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 5px;
    border: 1px solid #ddd;
    border-radius: 8px;
    background-color: #fff;
    width: calc(33.333% - 10px);  /* 三列布局 */
    box-sizing: border-box;
    transition: background-color 0.3s ease;
    cursor: pointer;  /* 改为指针手势，表明可以点击 */
}

/* 为复选框添加自定义样式 */
.switch-item input[type="checkbox"] {
    display: none;  /* 隐藏默认的复选框 */
}

/* 自定义复选框样式 */
.switch-item input[type="checkbox"] + span {
    position: relative;
    padding-left: 30px;
    cursor: pointer;
    font-size: 16px;
    line-height: 20px;
    transition: color 0.3s ease;
}

/* 复选框的背景 */
.switch-item input[type="checkbox"] + span::before {
    content: '';
    position: absolute;
    left: 0;
    top: 0;
    width: 15px;
    height: 15px;
    border: 2px solid #d1d1d1;  /* 更柔和的灰色边框 */
    border-radius: 4px;
    background-color: #fff;  /* 背景白色 */
    transition: background-color 0.3s ease, border-color 0.3s ease;
}

/* 当复选框被选中时 */
.switch-item input[type="checkbox"]:checked + span::before {
    background-color: #7c7c7c;  /* 选中时背景色为浅灰色 */
    border-color: #7c7c7c;  /* 边框变为浅灰色 */
}

/* 鼠标悬停时 */
.switch-item:hover {
    background-color: #f0f0f0;
}

/* 确保复选框和标签的布局 */
.switch-item input[type="checkbox"] {
    margin-right: 10px;
}

.chapter-list {
    display: flex;
    flex-wrap: wrap;
    gap: 10px;
}
</style>