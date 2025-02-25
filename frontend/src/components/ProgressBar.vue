<template>
  <div class="progress-bar-container">
    <div class="progress-bar" :style="progressStyle"></div>
  </div>
</template>

<script lang="ts">
import { defineComponent, computed } from 'vue';

export default defineComponent({
  name: 'ProgressBar',
  props: {
    progress: {
      type: Number,
      required: true,
      validator: (value: number) => value >= 0 && value <= 100,
    },
    color: {
      type: String,
      default: '#42b983',
    },
  },
  setup(props) {
    const progressStyle = computed(() => ({
      width: `${props.progress}%`,
      backgroundColor: props.color,
    }));

    return {
      progressStyle,
    };
  },
});
</script>

<style scoped>
.progress-bar-container {
  width: 100%;
  height: 10px;
  background-color: #e0e0e0;
  border-radius: 5px;
  overflow: hidden;
}

.progress-bar {
  height: 100%;
  transition: width 0.3s ease;
}
</style>