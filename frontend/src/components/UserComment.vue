<script setup lang="ts">
const props = defineProps([
  'rating',
  'imageUrls',
  'sellerResponse',
  'userName',
  'userAvatar',
  'date',
  'content',
  'variantName',
])
</script>

<template>
  <div style="display: flex; padding: 20px 0 28px">
    <img
      :src="props.userAvatar"
      style="
        padding: 0 10px 0 20px;
        border-radius: 50%;
        width: 40px;
        height: 40px;
        box-sizing: content-box;
      "
    />
    <div>
      <p>{{ props.userName }}</p>
      <el-rate
        :model-value="props.rating"
        disabled
        style="
          --el-rate-fill-color: var(--main-color);
          --el-rate-icon-margin: 2px;
          --el-rate-disabled-void-color: #cecece;
        "
        size="small"
      />
      <p style="color: #9c9c9c; font-size: 13px">
        {{ props.date }} <span v-if="props.variantName">| Option: {{ props.variantName }}</span>
      </p>
      <p style="margin-top: 8px; white-space: pre-wrap">
        {{ props.content }}
      </p>
      <div>
        <el-image
          v-for="(imageUrl, index) in imageUrls"
          :key="index"
          :initial-index="index"
          :src="imageUrl"
          style="
            width: 100px;
            height: 100px;
            border-radius: 2px;
            margin-right: 12px;
            margin-top: 12px;
          "
          show-progress
          :preview-src-list="imageUrls"
          fit="contain"
        />
      </div>
      <div
        v-if="props.sellerResponse"
        style="
          padding: 20px;
          background-color: #d6f2d7;
          border: 1px solid #c2f0c2;
          margin-top: 12px;
        "
      >
        <h4 style="font-weight: 500; font-size: 16px">Seller's Response</h4>
        <p style="margin-top: 4px">
          {{ props.sellerResponse }}
        </p>
      </div>
    </div>
  </div>

  <el-divider style="margin: 0" />
</template>
