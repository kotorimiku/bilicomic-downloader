<template>
    <div class="config">
        <div class="form-item">
            <label>保存路径</label>
            <input type="text" v-model="outputPath" placeholder="请输入保存路径" class="input-box" />
        </div>
        <div class="form-item">
            <label>打包方式</label>
            <select v-model="packageType" class="styled-select">
                <option value="cbz" title="会添加元数据ComicInfo.xml">cbz</option>
                <option value="zip">zip</option>
                <option value="epub" title="epub对 avif 图片支持不好，建议使用cbz">epub</option>
                <option value="image">图片</option>
            </select>
        </div>
        <div class="form-item">
            <label>图片格式</label>
            <select v-model="imageFormat" class="styled-select">
                <option value="source">原始格式</option>
                <option value="png">png</option>
                <option value="jpg">jpg</option>
            </select>
        </div>
        <div class="form-item">
            <label>命名风格</label>
            <select v-model="namingStyle" class="styled-select">
                <option value="title" title="第1话">title 第1话</option>
                <option value="index-title" title="1-第1话">index-title 1-第1话</option>
                <option value="02d-index-title" title="01-第1话">02d-index-title 01-第1话</option>
                <option value="03d-index-title" title="001-第1话">03d-index-title 001-第1话</option>
            </select>
        </div>
        <div class="form-item">
            <label>cookie</label>
            <input type="text" v-model="cookie" placeholder="请输入cookie" class="input-box" />
        </div>
        <div class="form-item">
            <button @click="saveConfig" class="btn save-btn">保存配置</button>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useToast } from "vue-toastification";
import { SaveConfig, GetConfig } from '../../wailsjs/go/bilicomicdownloader/Config'
import { bilicomicdownloader as model } from '../../wailsjs/go/models';

const outputPath = ref<string>("")
const packageType = ref<string>("")
const imageFormat = ref<string>("")
const namingStyle = ref<string>("")
const urlBase = ref<string>("")
const cookie = ref<string>("")
const toast = useToast();

const saveConfig = () => {
    SaveConfig({
        urlBase: urlBase.value, outputPath: outputPath.value, packageType: packageType.value, imageFormat: imageFormat.value, namingStyle: namingStyle.value, cookie: cookie.value
    }).then((res: any) => {
        console.log("配置已保存", res)
        toast.success('配置已保存', { timeout: 2000 });
    }).catch((err: any) => {
        console.error("保存配置失败", err)
        toast.error('保存配置失败', { timeout: 2000 });
    })
}

onMounted(() => {
    GetConfig().then((res: model.Config) => {
        if (res) {
            urlBase.value = res.urlBase
            outputPath.value = res.outputPath
            packageType.value = res.packageType
            imageFormat.value = res.imageFormat
            namingStyle.value = res.namingStyle
        }
    }).catch(() => {
        console.log("获取配置失败")
        toast.error('获取配置失败', { timeout: 2000 });
    })
})
</script>

<style scoped>
.config {
    width: 100%;
    max-width: 350px;
    padding: 20px;
    background-color: #ffffff;
    border-radius: 10px;
    box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
    margin: 20px auto;
    font-family: 'Arial', sans-serif;
}

.form-item {
    margin-bottom: 20px;
    display: flex;
    align-items: center;
}

label {
    font-weight: 600;
    font-size: 14px;
    color: #333;
    width: 100px;
    text-align: left;
}

.input-box,
.styled-select {
    width: 100%;
    padding: 10px;
    font-size: 14px;
    border-radius: 5px;
    border: 1px solid #ddd;
    transition: border-color 0.3s ease;
}

.input-box:focus,
.styled-select:focus {
    border-color: #4CAF50;
}

.styled-select {
    appearance: none;
    background-color: #fff;
    cursor: pointer;
}

.styled-select::after {
    content: '\2193';
    position: absolute;
    right: 10px;
    top: 50%;
    transform: translateY(-50%);
    font-size: 18px;
    color: #888;
    pointer-events: none;
}

.btn {
    padding: 10px 20px;
    font-size: 14px;
    background-color: #4CAF50;
    color: white;
    font-weight: 600;
    border: 1px solid #ddd;
    border-radius: 5px;
    cursor: pointer;
    width: 100%;
    transition: background-color 0.3s, transform 0.3s;
}

.btn:hover {
    background-color: #45a049;
    transform: translateY(-2px);
}
</style>
