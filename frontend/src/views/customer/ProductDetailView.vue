<script setup lang="ts">
import { Splide, SplideSlide, Options } from '@splidejs/vue-splide'
import Header from '@/components/Header.vue'
import { onMounted, ref, computed, reactive, watch } from 'vue'
import { quantityFormatNumber } from '@/utils/quantityFormatNumber'
import { formatNumberWithDots } from '@/utils/formatNumberWithDots'
import RedFlagIcon from '@/components/icons/RedFlagIcon.vue'
import MessengerIcon from '@/components/icons/MessengerIcon.vue'
import FacebookIcon from '@/components/icons/FacebookIcon.vue'
import PinterestIcon from '@/components/icons/PinterestIcon.vue'
import TwitterXIcon from '@/components/icons/TwitterXIcon.vue'
import { Check, Goods, ShoppingCart } from '@element-plus/icons-vue'
import HeartIcon from '@/components/icons/HeartIcon.vue'
import HeartFilledIcon from '@/components/icons/HeartFilledIcon.vue'
import ShippingIcon from '@/components/icons/ShippingIcon.vue'
import UserComment from '../../components/UserComment.vue'
import ProductItem from '@/components/ProductItem.vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'
import { ElNotification } from 'element-plus'

interface ProductImage {
  id: string
  url: string
  order: number
}

interface OptionGroup {
  key: string
  values: string[]
}

interface Variant {
  id: string
  sku: string
  options: Record<string, string>
  price: number
  stock: number
  image: string
  sold_count: number
}

interface Product {
  id: string
  name: string
  description: string
  images: ProductImage[]
  status: string
  seller_id: string
  rating: number
  rate_count: number
  sold_count: number
  is_active: boolean
  created_at: string
  updated_at: string
  price: {
    min: number
    max: number
  }
  stock: number
  option_groups: OptionGroup[]
  variants: Variant[]
}

const route = useRoute()
const router = useRouter()
const headerRef = ref<InstanceType<typeof Header> | null>(null)
const productId = ref(route.params.id as string)
const product = ref<Product | null>(null)
const isLoading = ref(true)
const isAddingToCart = ref(false)
const productList = ref<any[]>([]) // For related products
const recentlyViewedProducts = ref<any[]>([])

const mainOptions: Options = {
  type: 'fade',
  heightRatio: 1,
  pagination: false,
  arrows: true,
  cover: true,
}

const thumbsOptions: Options = {
  rewind: true,
  fixedWidth: 100,
  fixedHeight: 100,
  isNavigation: true,
  gap: 5,
  focus: 'center',
  pagination: false,
  cover: true,
  dragMinThreshold: {
    mouse: 4,
    touch: 10,
  },
}

const main = ref()
const thumbs = ref()

const selectedOptions = reactive<Record<string, string>>({})
const buyQuantity = ref(1)
const activeTab = ref('description')
const isDescriptionExpanded = ref(false)
const isDescriptionTooLong = ref(false)
const descriptionRef = ref<HTMLElement | null>(null)

const selectedVariant = computed(() => {
  if (!product.value) return null
  return product.value.variants.find((v) => {
    return Object.entries(selectedOptions).every(([key, value]) => v.options[key] === value)
  })
})

const displayImages = computed(() => {
  if (!product.value) return []
  const productImages = [...product.value.images]
    .sort((a, b) => a.order - b.order)
    .map((img) => img.url)
  const variantImages = product.value.variants
    .map((v) => v.image)
    .filter((img) => img && !productImages.includes(img))

  const all = [...productImages, ...new Set(variantImages)]
  return all.length > 0 ? all : ['https://placehold.co/600x600?text=No+Image']
})

const currentPrice = computed(() => {
  if (!product.value) return '0'
  if (selectedVariant.value) {
    return formatNumberWithDots(selectedVariant.value.price)
  }
  if (product.value.price.min === product.value.price.max) {
    return formatNumberWithDots(product.value.price.min)
  }
  return `${formatNumberWithDots(product.value.price.min)} - ${formatNumberWithDots(product.value.price.max)}`
})

const totalPrice = computed(() => {
  if (!product.value) return '0'
  const price = selectedVariant.value ? selectedVariant.value.price : product.value.price.min
  return formatNumberWithDots(price * buyQuantity.value)
})

const currentStock = computed(() => {
  if (selectedVariant.value) return selectedVariant.value.stock
  return product.value?.stock || 0
})

const fetchProduct = async () => {
  if (!productId.value) return
  isLoading.value = true
  try {
    const response = await axios.get(`http://localhost:81/api/product/public/${productId.value}`)
    product.value = response.data

    // Initialize selected options
    if (product.value?.option_groups) {
      product.value.option_groups.forEach((group) => {
        if (group.values.length > 0) {
          selectedOptions[group.key] = group.values[0] as string
        }
      })
    }
  } catch (error) {
    console.error('Error fetching product:', error)
  } finally {
    isLoading.value = false
    checkDescriptionHeight()
  }

  if (product.value) {
    saveToRecentlyViewed(product.value)
    loadRecentlyViewed()
  }
}

const checkDescriptionHeight = () => {
  setTimeout(() => {
    if (descriptionRef.value) {
      isDescriptionTooLong.value = descriptionRef.value.scrollHeight > 200
    }
  }, 100)
}

const saveToRecentlyViewed = (p: Product) => {
  const stored = localStorage.getItem('recently_viewed')
  let list = stored ? JSON.parse(stored) : []

  // Simplify product for storage
  const simplifiedProduct = {
    id: p.id,
    name: p.name,
    imageUrl: p.images.length > 0 ? p.images.sort((a, b) => a.order - b.order)[0]?.url || '' : '',
    price: p.variants.length > 0 ? p.variants[0]?.price || 0 : p.price.min,
    rating: p.rating,
    location: 'Vietnam', // Mock location as it's not in the product data
    discount: 0, // Mock discount
  }

  // Remove existing
  list = list.filter((item: any) => item.id !== simplifiedProduct.id)

  // Add to front
  list.unshift(simplifiedProduct)

  // Limit to 10
  list = list.slice(0, 10)

  localStorage.setItem('recently_viewed', JSON.stringify(list))
}

const loadRecentlyViewed = () => {
  const stored = localStorage.getItem('recently_viewed')
  if (stored) {
    recentlyViewedProducts.value = JSON.parse(stored)
  }
}

watch(
  () => route.params.id,
  (newId) => {
    if (newId && newId !== productId.value) {
      productId.value = newId as string
      product.value = null
      buyQuantity.value = 1
      activeTab.value = 'description'
      isDescriptionExpanded.value = false

      // Clear selected options to avoid stale keys from previous product
      Object.keys(selectedOptions).forEach((key) => delete selectedOptions[key])

      fetchProduct()
    }
  },
)

onMounted(() => {
  fetchProduct()
  loadRecentlyViewed()
  syncSplides()
})

const syncSplides = () => {
  const mainSplide = main.value?.splide
  const thumbsSplide = thumbs.value?.splide

  if (mainSplide && thumbsSplide) {
    mainSplide.sync(thumbsSplide)
  }
}

watch(
  displayImages,
  () => {
    setTimeout(() => {
      main.value?.splide?.refresh()
      thumbs.value?.splide?.refresh()
      syncSplides()
    }, 100)
  },
  { deep: true },
)

watch(selectedVariant, (newVariant) => {
  if (newVariant?.image) {
    const index = displayImages.value.indexOf(newVariant.image)
    if (index !== -1) {
      main.value?.splide?.go(index)
      return
    }
  }
  // Fallback to first image if variant has no image or image not found
  main.value?.splide?.go(0)
})

const selectOption = (key: string, value: string) => {
  selectedOptions[key] = value
}

const addToCart = async () => {
  if (!product.value || !selectedVariant.value) return

  const token = localStorage.getItem('access_token')
  if (!token) {
    ElNotification({
      title: 'Authentication Required',
      message: 'Please login to add items to your cart.',
      type: 'warning',
    })
    return
  }

  isAddingToCart.value = true
  try {
    await axios.post(
      'http://localhost:81/api/order/cart',
      {
        seller_id: product.value.seller_id,
        product_id: product.value.id,
        variant_id: selectedVariant.value.id,
        quantity: buyQuantity.value,
      },
      {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      },
    )

    ElNotification({
      title: 'Success',
      message: 'Product added to cart successfully!',
      type: 'success',
    })
    headerRef.value?.fetchCartCount()
  } catch (error: any) {
    console.error('Error adding to cart:', error)
    ElNotification({
      title: 'Error',
      message: error.response?.data?.message || 'Failed to add product to cart.',
      type: 'error',
    })
  } finally {
    isAddingToCart.value = false
  }
}
</script>

<template>
  <Header ref="headerRef" />

  <div class="main-container" v-if="product">
    <div class="box-shadow border-radius" style="padding: 20px; margin: 20px 0">
      <el-row :gutter="28">
        <el-col :span="8">
          <Splide :options="mainOptions" ref="main">
            <SplideSlide v-for="(imageUrl, index) in displayImages" :key="index">
              <img :src="imageUrl" alt="product-image" />
            </SplideSlide>
          </Splide>

          <span style="display: block; height: 5px"></span>

          <Splide
            aria-label="The carousel with thumbnails. Selecting a thumbnail will change the main carousel"
            :options="thumbsOptions"
            ref="thumbs"
          >
            <SplideSlide v-for="(imageUrl, index) in displayImages" :key="index">
              <img :src="imageUrl" alt="product-image" />
            </SplideSlide>
          </Splide>
        </el-col>
        <el-col :span="10">
          <h2
            style="
              font-weight: 500;
              margin-bottom: 4px;

              --line-height: 36px;
              height: calc(var(--line-height) * 2);
              line-height: var(--line-height);

              font-size: 16px;
              font-weight: 500;
              color: #000000cc;
              display: -webkit-box;
              -webkit-box-orient: vertical;
              line-clamp: 2;
              -webkit-line-clamp: 2;
              overflow: hidden;
              text-overflow: ellipsis;
            "
          >
            {{ product.name }}
          </h2>
          <div style="display: flex; align-items: center; justify-content: space-between">
            <div style="display: flex; align-items: center">
              <el-rate
                v-model="product.rating"
                disabled
                show-score
                style="--el-rate-fill-color: var(--main-color)"
                text-color="var(--main-color)"
                score-template="{value} points"
              />
              <el-divider direction="vertical" />
              <div>{{ quantityFormatNumber(product.rate_count) }} reviews</div>
              <el-divider direction="vertical" />
              <div>{{ quantityFormatNumber(product.sold_count) }} sold</div>
            </div>
            <el-button :icon="RedFlagIcon">Report</el-button>
          </div>
          <div style="font-weight: 600; font-size: 28px">${{ currentPrice }}</div>
          <div style="margin-bottom: 16px" v-if="false">
            <!-- Example of special offers if any -->
            <el-divider style="margin: 20px 0" />
            <span style="margin: 0 8px">Special Offers</span>
            <div style="border: 1px solid black; display: inline-block; padding: 2px 6px">
              Buy 2 & get 5% off
            </div>
          </div>

          <div
            v-for="group in product.option_groups"
            :key="group.key"
            style="display: flex; align-items: center; margin-bottom: 16px"
          >
            <span style="margin: 0 8px; text-transform: capitalize; min-width: 60px">{{
              group.key
            }}</span>
            <div style="display: flex; flex-wrap: wrap; gap: 8px">
              <div
                v-for="val in group.values"
                :key="val"
                class="option-wrapper option-wrapper--no-image"
                :class="{ active: selectedOptions[group.key] === val }"
                @click="selectOption(group.key, val)"
              >
                <p>{{ val }}</p>
              </div>
            </div>
          </div>

          <el-divider style="margin: 20px 0" />

          <div>
            <div style="margin-bottom: 4px" v-if="selectedVariant">
              <span style="font-weight: 700; margin-right: 4px">SKU:</span>
              <span>{{ selectedVariant.sku }}</span>
            </div>

            <div style="margin-bottom: 4px">
              <span style="font-weight: 700; margin-right: 4px">STATUS:</span>
              <span :style="{ color: product.is_active ? 'green' : 'red' }">{{
                product.status
              }}</span>
            </div>
          </div>

          <el-divider style="margin: 20px 0" />

          <div>
            <el-button circle>
              <el-icon size="large"><MessengerIcon /></el-icon>
            </el-button>
            <el-button circle>
              <el-icon size="large"><FacebookIcon /></el-icon>
            </el-button>
            <el-button circle>
              <el-icon size="large"><PinterestIcon /></el-icon>
            </el-button>
            <el-button circle>
              <el-icon size="large"><TwitterXIcon /></el-icon>
            </el-button>
          </div>
        </el-col>

        <el-col :span="6">
          <div
            class="box-sizing border-box"
            style="background-color: #edeff6; padding: 40px 40px 20px"
          >
            <h4>TOTAL PRICE:</h4>
            <span style="font-weight: 700; font-size: 30px">${{ totalPrice }}</span>
            <div style="display: flex; align-items: center; margin: 12px 0 20px">
              <el-icon
                style="margin-right: 6px; background-color: green; border-radius: 50%; padding: 2px"
                color="#fff"
                v-if="currentStock > 0"
                ><Check
              /></el-icon>
              <p v-if="currentStock > 0">In stock: {{ currentStock }}</p>
              <p v-else style="color: red">Out of stock</p>
            </div>
            <el-input-number
              v-model="buyQuantity"
              :min="1"
              :max="currentStock || 1"
              :disabled="currentStock <= 0"
              size="large"
              style="width: 100%"
            />

            <el-button
              type="primary"
              color="#1aba1a"
              size="large"
              :disabled="currentStock <= 0"
              :loading="isAddingToCart"
              style="width: 100%; margin: 20px 0"
              @click="addToCart"
            >
              <el-icon size="large" style="margin-right: 6px"><ShoppingCart /></el-icon>
              Add to Cart
            </el-button>
            <el-button
              type="primary"
              color="var(--main-color)"
              size="large"
              :disabled="currentStock <= 0"
              style="width: 100%; margin: 0"
            >
              <el-icon size="large" style="margin-right: 6px"><Goods /></el-icon>
              Buy Now
            </el-button>
            <div style="display: flex; align-items: center; margin-top: 24px">
              <el-icon style="fill: var(--main-color)"><HeartIcon /></el-icon>
              <!-- <el-icon style="fill: var(--main-color); margin-right: 4px"
                ><HeartFilledIcon
              /></el-icon> -->
              <span style="color: var(--main-color); font-size: 14px; margin-left: 8px"
                >Add to Wishlist</span
              >
            </div>
          </div>

          <div style="display: flex; align-items: center; margin-left: 8px; margin-top: 8px">
            <el-icon size="large" style="margin-right: 6px"><ShippingIcon /></el-icon>
            <span> Ships directly from <span style="font-weight: 700">Seller</span> </span>
          </div>
        </el-col>
      </el-row>
    </div>

    <div class="box-shadow border-radius" style="padding: 20px; margin-bottom: 20px">
      <el-tabs v-model="activeTab">
        <el-tab-pane label="DESCRIPTION" name="description">
          <div
            ref="descriptionRef"
            style="padding: 0 20px; white-space: pre-wrap"
            class="description-tab"
            :class="{ collapsed: !isDescriptionExpanded && isDescriptionTooLong }"
          >
            {{ product.description }}
          </div>

          <div style="display: flex; align-items: center" v-if="isDescriptionTooLong">
            <el-button
              plain
              color="#22c55e"
              size="large"
              style="margin: 28px auto 16px"
              @click="isDescriptionExpanded = !isDescriptionExpanded"
            >
              {{ isDescriptionExpanded ? 'Show Less' : 'Show More' }}
            </el-button>
          </div>
        </el-tab-pane>

        <el-tab-pane :label="`REVIEWS (${product.rate_count})`" name="reviews">
          <div
            style="
              padding: 20px;
              background-color: #e2f0c4;
              border: 1px solid #84f084;
              margin-bottom: 16px;
            "
          >
            <el-row style="align-items: center">
              <el-col
                :span="6"
                style="
                  color: var(--main-color);
                  display: flex;
                  flex-direction: column;
                  align-items: center;
                "
              >
                <span style="font-size: 28px; font-weight: 600"
                  >{{ product.rating }}<span style="font-size: 24px"> out of 5</span></span
                >
                <el-rate
                  v-model="product.rating"
                  disabled
                  style="--el-rate-fill-color: var(--main-color); --el-rate-icon-size: 28px"
                  size="large"
                />
              </el-col>
              <el-col :span="18">
                <button class="rating-filter-btn active">All</button>
                <button class="rating-filter-btn">5 Stars</button>
                <button class="rating-filter-btn">4 Stars</button>
                <button class="rating-filter-btn">3 Stars</button>
                <button class="rating-filter-btn">2 Stars</button>
                <button class="rating-filter-btn">1 Star</button>
                <button class="rating-filter-btn">Has Comment</button>
                <button class="rating-filter-btn">Has Pictures</button>
              </el-col>
            </el-row>
          </div>
          <div v-if="product.rate_count === 0" style="text-align: center; padding: 40px">
            No reviews yet.
          </div>
          <template v-else>
            <UserComment
              rating="4"
              seller-response="Thank you for sharing a review of your recent experience. Our team at Lovito is committed
            to providing high-quality products at affordable prices, and we are delighted to know that
            we made a positive impression on you. We hope you will visit us again soon."
            />
          </template>

          <div style="display: flex" v-if="product.rate_count > 0">
            <el-pagination
              background
              layout="prev, pager, next"
              :total="product.rate_count"
              size="large"
              style="margin: 28px auto 8px"
            />
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>

    <div
      class="box-shadow border-radius"
      style="background-color: #fff; padding: 20px; margin-bottom: 20px"
    >
      <h3 style="font-weight: bold">RELATED PRODUCTS</h3>

      <el-row :gutter="20" v-if="productList.length">
        <el-col
          :span="4.8"
          class="el-col-4-8"
          v-for="(item, index) in productList.slice(0, 10)"
          :key="index"
        >
          <ProductItem
            :image-url="item.imageUrl"
            :name="item.name"
            :price="item.price"
            :rating="item.rating"
            :location="item.location"
            :discount="item.discount"
            :id="item.id"
          />
        </el-col>
      </el-row>
      <div v-else style="text-align: center; color: #999; padding: 40px">
        <el-icon size="40" style="margin-bottom: 12px; opacity: 0.5"><Goods /></el-icon>
        <p style="font-size: 14px">No related products found.</p>
      </div>
    </div>

    <div
      class="box-shadow border-radius recently-viewed"
      style="background-color: #fff; padding: 20px 20px 28px; margin-bottom: 20px"
    >
      <div style="display: flex; margin-bottom: 12px">
        <h3 style="font-weight: bold">YOU RECENTLY VIEWED</h3>
        <RouterLink
          to="/"
          style="position: relative; top: 3px; margin-left: 20px; color: #999; font-size: 13px"
          >View All</RouterLink
        >
      </div>

      <Splide
        v-if="recentlyViewedProducts.length"
        :options="{
          rewind: true,
          perPage: 5,
          gap: '1rem',
          breakpoints: {
            1000: {
              perPage: 1,
            },
          },
        }"
      >
        <SplideSlide v-for="item in recentlyViewedProducts" :key="item.id">
          <ProductItem
            :image-url="item.imageUrl"
            :name="item.name"
            :price="item.price"
            :rating="item.rating"
            :location="item.location"
            :discount="item.discount"
            :id="item.id"
          />
        </SplideSlide>
      </Splide>
      <div v-else style="text-align: center; color: #999; padding: 40px">
        <el-icon size="40" style="margin-bottom: 12px; opacity: 0.5"><ShoppingCart /></el-icon>
        <p style="font-size: 14px">Your recently viewed list is empty.</p>
        <p style="font-size: 13px; color: #bbb">
          Products you visit will appear here for easy access.
        </p>
      </div>
    </div>
  </div>
  <div
    v-else-if="isLoading"
    style="height: 100vh; display: flex; align-items: center; justify-content: center"
  >
    <el-icon class="is-loading" size="40"><Loading /></el-icon>
  </div>
  <div v-else style="height: 100vh; display: flex; align-items: center; justify-content: center">
    <el-result
      icon="error"
      title="Product not found"
      sub-title="Please check the URL or try again later"
    >
      <template #extra>
        <el-button type="primary" @click="$router.push('/')">Back to Home</el-button>
      </template>
    </el-result>
  </div>
</template>

<style lang="css" scoped>
.option-wrapper {
  display: flex;
  align-items: center;
  height: 48px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  padding: 2px 12px;
  margin-right: 12px;
  transition: all 0.2s;
  background-color: #fff;
  user-select: none;

  p {
    margin: 0;
    font-size: 14px;
  }

  &:hover {
    cursor: pointer;
    color: #22c55e;
    border-color: #22c55e;
  }

  &.active {
    cursor: pointer;
    color: #22c55e;
    border-color: #22c55e;
    background-color: rgba(34, 197, 94, 0.05);
    font-weight: 600;
    border-width: 2px;
  }
}

.option-wrapper--no-image {
  justify-content: center;
  height: 40px;
  min-width: 72px;
}

.description-tab {
  p {
    margin-bottom: 16px;
    line-height: 1.6;
    color: #333;
    font-size: 14px;
  }

  h2 {
    margin: 24px 0 16px;
    color: #222;
    font-weight: 700;
    font-size: 30px;
  }

  i {
    color: #666;
    font-size: 12px;
  }

  &.collapsed {
    max-height: 200px;
    overflow: hidden;
    position: relative;

    &::after {
      content: '';
      position: absolute;
      bottom: 0;
      left: 0;
      width: 100%;
      height: 80px;
      background: linear-gradient(transparent, #fff);
    }
  }
}

.rating-filter-btn {
  margin-right: 12px;
  margin-bottom: 6px;
  margin-top: 6px;
  font-size: 16px;
  padding: 8px 16px;
  border: 1px solid #ccc;
  border-radius: 4px;
  cursor: pointer;
  min-width: 80px;
  transition:
    color 0.3s,
    border-color 0.3s;

  &:hover {
    color: #22c55e;
    border-color: currentColor;
  }
}

.rating-filter-btn.active {
  border-color: #22c55e;
  color: #22c55e;
}
</style>

<style>
.el-pagination.is-background .btn-next.is-active,
.el-pagination.is-background .btn-prev.is-active,
.el-pagination.is-background .el-pager li.is-active {
  background-color: var(--main-color);
  color: var(--el-color-white);
}

.el-pager li:hover {
  color: var(--main-color);
  transition: color 0.2s;
}
</style>

<style>
.el-rate__item .el-rate__icon {
  stroke: var(--main-color);
  fill: var(--main-color) !important;
  stroke-width: 32px;
}

.el-rate__item .el-rate__icon.is-active {
  color: var(--main-color) !important;
}
</style>
