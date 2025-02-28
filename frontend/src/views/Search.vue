<template>
    <div class="container">
        <!-- 书籍 ID 输入框 -->
        <div class="input-group">
            <input type="text" v-model="bookId" placeholder="请输入书籍 ID" class="input-box" @keyup.enter="getChapterList"/>
            <button @click="getChapterList" class="btn">搜索</button>
        </div>

        <!-- 翻页按钮和下载操作按钮 -->
        <div class="action-buttons">

            <!-- 下载操作按钮 -->
            <div class="button-group">
                <button @click="selectAll" class="btn">全选</button>
                <button @click="selectInverse" class="btn">反选</button>
                <button @click="download" :disabled="!chapterList.length" class="btn"
                    :class="{ disabled: !chapterList.length }">开始下载</button>
            </div>
        </div>

        <!-- 书籍信息和卷列表 -->
        <div class="box">

            <div class="chapter-list">
                <chapterList :chapterList="chapterList" v-model:selectedChapters="selectedChapters"
                    :title="bookInfo?.Title" :cover="bookInfo?.Cover" />
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { GetDownloader, GetBookInfo, GetChapter, DownloadList } from '../../wailsjs/go/bilicomicdownloader/DownloaderManager';
import { bilicomicdownloader as model } from '../../wailsjs/go/models';
import ChapterList from '../components/ChapterList.vue';
import { useToast } from 'vue-toastification';

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

const toast = useToast();

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
    try {
        await GetDownloader(bookId.value);
        GetBookInfo().then((bookInfoRes) => {
            bookInfo.value = bookInfoRes;
        });

        GetChapter().then((chapterListRes) => {
            selectedChapters.value = [];
            chapterList.value = chapterListRes;
        });
    } catch (err) {
        console.error(err);
        toast.error(err, {
            timeout: 2000,
            closeOnClick: false,
        });
    }
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
    try {
        isDownloading.value = true;
        DownloadList(selectedChapters.value);
        selectedChapters.value = [];
    } catch (err) {
        console.error(err);
        toast.error(err, {
            timeout: 2000,
            closeOnClick: false,
        });
    } finally {
        isDownloading.value = false;
    }
}, 200);
</script>

<style scoped>
/* 容器布局 */
.container {
    margin: 0 auto;
    padding: 30px;
    background-color: #fafafa;
    border-radius: 10px;
}

/* 输入框与按钮组布局 */
.input-group,
.action-buttons {
    display: flex;
    gap: 15px;
    margin-bottom: 20px;
    flex-wrap: wrap;
    /* 使得在小屏幕下换行 */
}

.action-buttons {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.pagination,
.button-group {
    display: flex;
    gap: 15px;
    align-items: center;
}

/* 输入框样式 */
.input-box {
    flex-grow: 1;
    padding: 12px 15px;
    font-size: 16px;
    border-radius: 8px;
    border: 1px solid #ddd;
    background-color: #fff;
    box-sizing: border-box;
    transition: border-color 0.3s ease;
}

.input-box:focus {
    border-color: #007BFF;
    outline: none;
}

/* 按钮样式 */
.btn {
    padding: 10px 20px;
    font-size: 16px;
    background-color: #fff;
    border: 1px solid #ddd;
    border-radius: 8px;
    cursor: pointer;
    transition: background-color 0.3s ease, transform 0.2s ease;
}

.btn:hover {
    background-color: #f1f1f1;
    transform: translateY(-2px);
}

.btn:disabled,
.disabled {
    background-color: #e0e0e0;
    cursor: not-allowed;
}

.btn:active {
    transform: translateY(1px);
}

/* 容器布局 */
.box {
    display: flex;
    justify-content: space-between;
    gap: 25px;
}

/* 书籍列表 */
.comic-list {
    flex: 1;
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(110px, 1fr));
    /* 自适应列数 */
    height: 100%;
    gap: 15px;
    margin: 0;
}

.comic-item {
    text-align: center;
    background-color: #fff;
    border-radius: 8px;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    overflow: hidden;
    transition: transform 0.3s ease;
}

.comic-item:hover {
    transform: translateY(-5px);
}

.comic-cover {
    width: 100%;
    object-fit: cover;
    border-bottom: 1px solid #ddd;
}

.comic-item p {
    font-size: 14px;
    color: #333;
    margin-top: 10px;
    text-overflow: ellipsis;
    overflow: hidden;
    white-space: nowrap;
}

/* 章节列表 */
.chapter-list {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 15px;
    overflow-y: auto;
    /* 添加滚动条 */
}

/* 响应式布局 */
@media (max-width: 768px) {

    .input-group,
    .action-buttons {
        flex-direction: column;
        gap: 10px;
    }

    .action-buttons {
        flex-direction: column;
        align-items: flex-start;
    }

    .pagination,
    .button-group {
        width: 100%;
        justify-content: space-between;
    }

    .box {
        flex-direction: column;
    }
}
</style>