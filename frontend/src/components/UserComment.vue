<script setup lang="ts">
const props = defineProps(['rating', 'imageUrls', 'sellerResponse'])
</script>

<template>
  <div style="display: flex; padding: 20px 0 28px">
    <img
      src="/src/assets/avatar.jpg"
      style="
        padding: 0 10px 0 20px;
        border-radius: 50%;
        width: 40px;
        height: 40px;
        box-sizing: content-box;
      "
    />
    <div>
      <p>ntduc24</p>
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
      <p style="color: #9c9c9c; font-size: 13px">2024-12-30 22:08 | Option: Gray, S</p>
      <p style="margin-top: 8px">
        The product is really beautiful, everyone. You should buy it quickly while you can because
        it’s super pretty. Hurry while there’s still a subsidy/discount available
      </p>
      <div>
        <el-image
          v-for="(imageUrl, index) in imageUrls"
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
