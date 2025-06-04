<template>
  <n-card title="配置设置" class="max-w-90 w-350 mx-auto my-5">
    <n-form label-placement="left" label-width="100px" :model="form">
      <n-form-item label="保存路径">
        <n-input v-model:value="form.outputPath" placeholder="请输入保存路径" />
      </n-form-item>

      <n-form-item label="打包方式">
        <n-select
          v-model:value="form.packageType"
          :options="packageOptions"
          placeholder="请选择打包方式"
        />
      </n-form-item>

      <n-form-item label="图片格式">
        <n-select
          v-model:value="form.imageFormat"
          :options="imageFormatOptions"
          placeholder="请选择图片格式"
        />
      </n-form-item>

      <n-form-item label="命名风格">
        <n-select
          v-model:value="form.namingStyle"
          :options="namingStyleOptions"
          placeholder="请选择命名风格"
        />
      </n-form-item>

      <n-form-item label="Cookie">
        <n-input v-model:value="form.cookie" placeholder="请输入 cookie" />
      </n-form-item>

      <n-form-item>
        <n-button type="primary" block @click="saveConfig">保存配置</n-button>
      </n-form-item>
    </n-form>
  </n-card>
</template>


<script setup lang="ts">
import { ref, onMounted } from 'vue'
import {
  SaveConfig,
  GetConfig
} from '../../wailsjs/go/bilicomicdownloader/Config'
import { bilicomicdownloader as model } from '../../wailsjs/go/models'
import { useRunCommand } from '../composables/RunCommand'
import { useNotify } from '../composables/useNotification'

const runCommand = useRunCommand()
const notify = useNotify()

const form = ref({
  outputPath: '',
  packageType: '',
  imageFormat: '',
  namingStyle: '',
  cookie: '',
})

const packageOptions = [
  { label: 'cbz（含ComicInfo.xml）', value: 'cbz' },
  { label: 'zip', value: 'zip' },
  { label: 'epub（不支持avif）', value: 'epub' },
  { label: '图片', value: 'image' }
]

const imageFormatOptions = [
  { label: '原始格式', value: 'source' },
  { label: 'png', value: 'png' },
  { label: 'jpg', value: 'jpg' }
]

const namingStyleOptions = [
  { label: 'title 第1话', value: 'title' },
  { label: 'index-title 1-第1话', value: 'index-title' },
  { label: '02d-index-title 01-第1话', value: '02d-index-title' },
  { label: '03d-index-title 001-第1话', value: '03d-index-title' }
]

const saveConfig = () => {
  runCommand({
    command: () =>
      SaveConfig({
        urlBase: '', // 可根据实际需求添加 urlBase 字段
        outputPath: form.value.outputPath,
        packageType: form.value.packageType,
        imageFormat: form.value.imageFormat,
        namingStyle: form.value.namingStyle,
        cookie: form.value.cookie
      }),
    onSuccess: () => notify.success({ content: '保存成功' }),
    errMsg: '保存失败'
  })
}

onMounted(() => {
  runCommand({
    command: GetConfig,
    onSuccess: (res: model.Config) => {
      if (res) {
        form.value.outputPath = res.outputPath
        form.value.packageType = res.packageType
        form.value.imageFormat = res.imageFormat
        form.value.namingStyle = res.namingStyle
        form.value.cookie = res.cookie
      }
    },
    errMsg: '获取配置失败'
  })
})
</script>
