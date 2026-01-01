<script setup lang="ts">
import { Splide, SplideSlide, Options } from '@splidejs/vue-splide'
import Header from '@/components/Header.vue'
import { onMounted, ref } from 'vue'
import { quantityFormatNumber } from '@/utils/quantityFormatNumber'
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

const productList = [
  {
    imageUrl: '/src/assets/product-imgs/product1.png',
    name: 'Laptop Dell 15 DC15250 71071928 (Core i5-1334U/16GB/512GB/Intel Graphics/15.6 inch FHD IPS/Win 11/Bạc) - Chính hãng',
    price: 16090000,
    rating: 5.0,
    location: 'Ha Noi',
    discount: 37,
  },
  {
    imageUrl: '/src/assets/product-imgs/product1.png',
    name: 'Laptop Dell 15 DC15250 71071928 (Core i5-1334U/16GB/512GB/Intel Graphics/15.6 inch FHD IPS/Win 11/Bạc) - Chính hãng',
    price: 16090000,
    rating: 5.0,
    location: 'Ha Noi',
    discount: 37,
  },
  {
    imageUrl: '/src/assets/product-imgs/product1.png',
    name: 'Laptop Dell 15 DC15250 71071928 (Core i5-1334U/16GB/512GB/Intel Graphics/15.6 inch FHD IPS/Win 11/Bạc) - Chính hãng',
    price: 16090000,
    rating: 5.0,
    location: 'Ha Noi',
    discount: 37,
  },
  {
    imageUrl: '/src/assets/product-imgs/product1.png',
    name: 'Laptop Dell 15 DC15250 71071928 (Core i5-1334U/16GB/512GB/Intel Graphics/15.6 inch FHD IPS/Win 11/Bạc) - Chính hãng',
    price: 16090000,
    rating: 5.0,
    location: 'Ha Noi',
    discount: 37,
  },
  {
    imageUrl: '/src/assets/product-imgs/product1.png',
    name: 'Laptop Dell 15 DC15250 71071928 (Core i5-1334U/16GB/512GB/Intel Graphics/15.6 inch FHD IPS/Win 11/Bạc) - Chính hãng',
    price: 16090000,
    rating: 5.0,
    location: 'Ha Noi',
    discount: 37,
  },
  {
    imageUrl: '/src/assets/product-imgs/product1.png',
    name: 'Laptop Dell 15 DC15250 71071928 (Core i5-1334U/16GB/512GB/Intel Graphics/15.6 inch FHD IPS/Win 11/Bạc) - Chính hãng',
    price: 16090000,
    rating: 5.0,
    location: 'Ha Noi',
    discount: 37,
  },
  {
    imageUrl: '/src/assets/product-imgs/product1.png',
    name: 'Laptop Dell 15 DC15250 71071928 (Core i5-1334U/16GB/512GB/Intel Graphics/15.6 inch FHD IPS/Win 11/Bạc) - Chính hãng',
    price: 16090000,
    rating: 5.0,
    location: 'Ha Noi',
    discount: 37,
  },
  {
    imageUrl: '/src/assets/product-imgs/product1.png',
    name: 'Laptop Dell 15 DC15250 71071928 (Core i5-1334U/16GB/512GB/Intel Graphics/15.6 inch FHD IPS/Win 11/Bạc) - Chính hãng',
    price: 16090000,
    rating: 5.0,
    location: 'Ha Noi',
    discount: 37,
  },
  {
    imageUrl: '/src/assets/product-imgs/product1.png',
    name: 'Laptop Dell 15 DC15250 71071928 (Core i5-1334U/16GB/512GB/Intel Graphics/15.6 inch FHD IPS/Win 11/Bạc) - Chính hãng',
    price: 16090000,
    rating: 5.0,
    location: 'Ha Noi',
    discount: 37,
  },
  {
    imageUrl: '/src/assets/product-imgs/product1.png',
    name: 'Laptop Dell 15 DC15250 71071928 (Core i5-1334U/16GB/512GB/Intel Graphics/15.6 inch FHD IPS/Win 11/Bạc) - Chính hãng',
    price: 16090000,
    rating: 5.0,
    location: 'Ha Noi',
    discount: 37,
  },
  {
    imageUrl: '/src/assets/product-imgs/product1.png',
    name: 'Laptop Dell 15 DC15250 71071928 (Core i5-1334U/16GB/512GB/Intel Graphics/15.6 inch FHD IPS/Win 11/Bạc) - Chính hãng',
    price: 16090000,
    rating: 5.0,
    location: 'Ha Noi',
    discount: 37,
  },
  {
    imageUrl: '/src/assets/product-imgs/product1.png',
    name: 'Laptop Dell 15 DC15250 71071928 (Core i5-1334U/16GB/512GB/Intel Graphics/15.6 inch FHD IPS/Win 11/Bạc) - Chính hãng',
    price: 16090000,
    rating: 5.0,
    location: 'Ha Noi',
    discount: 37,
  },
  {
    imageUrl: '/src/assets/product-imgs/product1.png',
    name: 'Laptop Dell 15 DC15250 71071928 (Core i5-1334U/16GB/512GB/Intel Graphics/15.6 inch FHD IPS/Win 11/Bạc) - Chính hãng',
    price: 16090000,
    rating: 5.0,
    location: 'Ha Noi',
    discount: 37,
  },
]

const mainOptions: Options = {
  type: 'fade',
  heightRatio: 1,
  pagination: false,
  arrows: false,
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
  // breakpoints: {
  //   640: {
  //     fixedWidth: 66,
  //     fixedHeight: 38,
  //   },
  // },
}

const imageUrls = [
  '/src/assets/product-imgs/product1.png',
  '/src/assets/product-imgs/product1.png',
  '/src/assets/product-imgs/product1.png',
  '/src/assets/product-imgs/product1.png',
]

const main = ref()
const thumbs = ref()

const rating = 4.6
const reviewQuantity = 120
const soldQuantity = 20000

const inStock = 15
const buyQuantity = ref(1)
const location = 'Hanoi, Vietnam'

const activeTab = ref('description')
const totalReviews = ref(256)

onMounted(() => {
  const thumbsSplide = thumbs.value?.splide

  if (thumbsSplide) {
    main.value?.sync(thumbsSplide)
  }
})
</script>

<template>
  <Header />

  <div class="main-container">
    <div class="box-shadow border-radius" style="padding: 20px; margin: 20px 0">
      <el-row :gutter="28">
        <el-col :span="8">
          <Splide :options="mainOptions" ref="main">
            <SplideSlide v-for="(imageUrl, index) in imageUrls" :key="index">
              <img :src="imageUrl" alt="product-image" />
            </SplideSlide>
          </Splide>

          <span style="display: block; height: 5px"></span>

          <Splide
            aria-label="The carousel with thumbnails. Selecting a thumbnail will change the main carousel"
            :options="thumbsOptions"
            ref="thumbs"
          >
            <SplideSlide v-for="(imageUrl, index) in imageUrls" :key="index">
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
            Product Name Product Name Product Name Product Name Product Name Product Name Product
            Name Product NameProduct Name Product Name Product Name Product Name Product Name
            Product Name Product Name
          </h2>
          <div style="display: flex; align-items: center; justify-content: space-between">
            <div style="display: flex; align-items: center">
              <el-rate
                v-model="rating"
                disabled
                show-score
                text-color="#ff9900"
                score-template="{value} points"
              />
              <el-divider direction="vertical" />
              <div>{{ quantityFormatNumber(reviewQuantity) }} reviews</div>
              <el-divider direction="vertical" />
              <div>{{ quantityFormatNumber(soldQuantity) }} sold</div>
            </div>
            <el-button :icon="RedFlagIcon">Report</el-button>
          </div>
          <div style="font-weight: 600; font-size: 28px">$569.00 - $609.00</div>
          <div style="margin-bottom: 16px">
            <el-divider style="margin: 20px 0" />
            <span style="margin: 0 8px">Special Offers</span>
            <div style="border: 1px solid black; display: inline-block; padding: 2px 6px">
              Buy 2 & get 5% off
            </div>
          </div>
          <div style="display: flex; align-items: center; margin-bottom: 16px">
            <span style="margin: 0 8px">Color</span>
            <div style="display: flex">
              <div class="option-wrapper">
                <img
                  :src="imageUrls[0]"
                  alt="product-image"
                  style="height: 100%; object-fit: contain; margin-right: 8px; padding: 4px 0"
                />
                <p>Red</p>
              </div>
            </div>
          </div>
          <div style="display: flex; align-items: center">
            <span style="margin: 0 8px">Size</span>
            <div style="display: flex">
              <div class="option-wrapper option-wrapper--no-image">
                <p>XS</p>
              </div>

              <div class="option-wrapper option-wrapper--no-image">
                <p>S</p>
              </div>
            </div>
          </div>

          <el-divider style="margin: 20px 0" />

          <div>
            <div style="margin-bottom: 4px">
              <span style="font-weight: 700; margin-right: 4px">SKU:</span>
              <span>ABC025168</span>
            </div>

            <div style="margin-bottom: 4px">
              <span style="font-weight: 700; margin-right: 4px">CATEGORY:</span>
              <span>Laptop</span>
            </div>

            <div style="margin-bottom: 4px">
              <span style="font-weight: 700; margin-right: 4px">Brand:</span>
              <span>Dell</span>
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
            <span style="font-weight: 700; font-size: 30px">$609.00</span>
            <div style="display: flex; align-items: center; margin: 12px 0 20px">
              <el-icon
                style="margin-right: 6px; background-color: green; border-radius: 50%; padding: 2px"
                color="#fff"
                ><Check
              /></el-icon>
              <p>In stock: {{ inStock }}</p>
            </div>
            <el-input-number
              v-model="buyQuantity"
              :min="1"
              :max="10"
              size="large"
              style="width: 100%"
            />

            <el-button
              type="primary"
              color="#1aba1a"
              size="large"
              style="width: 100%; margin: 20px 0"
            >
              <el-icon size="large" style="margin-right: 6px"><ShoppingCart /></el-icon>
              Add to Cart
            </el-button>
            <el-button type="primary" color="#f6ab4a" size="large" style="width: 100%; margin: 0">
              <el-icon size="large" style="margin-right: 6px"><Goods /></el-icon>
              Buy Now
            </el-button>
            <div style="display: flex; align-items: center; margin-top: 24px">
              <el-icon style="fill: var(--main-color)"><HeartIcon /></el-icon>
              <el-icon style="fill: var(--main-color); margin-right: 4px"
                ><HeartFilledIcon
              /></el-icon>
              <span style="color: var(--main-color); font-size: 14px">Add to Wishlist</span>
            </div>
          </div>

          <div style="display: flex; align-items: center; margin-left: 8px; margin-top: 8px">
            <el-icon size="large" style="margin-right: 6px"><ShippingIcon /></el-icon>
            <span>
              Ships from <span style="font-weight: 700"> {{ ` ${location}` }}</span>
            </span>
          </div>
        </el-col>
      </el-row>
    </div>

    <div class="box-shadow border-radius" style="padding: 20px; margin-bottom: 20px">
      <el-tabs v-model="activeTab">
        <el-tab-pane label="DESCRIPTION" name="description">
          <div style="padding: 0 20px" class="description-tab">
            <p>
              Built for ultra-fast performance, the thin and lightweight Samsung Galaxy Tab S2 goes
              anywhere you go. Photos, movies and documents pop on a crisp, clear Super AMOLED
              display. Expandable memory lets you enjoy more of your favorite content. And
              connecting and sharing between all your Samsung devices is easier than ever. Welcome
              to life with the reimagined Samsung Galaxy Tab S2. Watch the world come to life on
              your tablet's Super AMOLED display *. With deep contrast, rich colors and crisp
              details, you won't miss a thing.
            </p>
            <img
              src="/src/assets/product-imgs/description-imgs/description-image-1.jpg"
              style="width: 100%; margin-top: 20px; object-fit: contain; border-radius: 20px"
              alt="product-description-image"
            />
            <i style="text-align: center; display: block; font-style: italic"
              >* The Galaxy Tab S2's $4:3$ ratio display provides you with an ideal environment for
              performing office tasks.</i
            >

            <h2>From the manufacturer</h2>
            <p>
              Dive into the blockbuster movies you can't wait to see. Switch between your favorite
              apps quickly and easily. The new and improved octa-core processor gives you the power
              and speed you need to see more and do more. Expand your tablet's memory from 32GB to
              up to an additional 128GB and enjoy more of your favorite music, photos, movies and
              games on the go with a microSD card. With Quick Connect, start a show on your Smart TV
              and, with the touch of a button, take it with you by moving it to your Galaxy Tab S2.
            </p>
            <p>
              Or send videos and photos from your tablet screen to your TV screen to share with
              everyone in the room. Work effortlessly between your Samsung tablet and Samsung
              smartphone with SideSync. Quickly drag and drop photos between devices. And even
              respond to a call from your smartphone right on your tablet screen.
            </p>
          </div>

          <el-row :gutter="20" style="margin-top: 20px">
            <el-col :span="12" v-for="n in 2" :key="n">
              <img
                :src="`/src/assets/product-imgs/description-imgs/description-image-1.jpg`"
                style="width: 100%; object-fit: contain; border-radius: 20px"
                :alt="`product-description-image-${n}`"
              />
            </el-col>
          </el-row>

          <div style="display: flex; align-items: center">
            <el-button plain color="#22c55e" size="large" style="margin: 28px auto 16px">
              Show More
            </el-button>
          </div>
        </el-tab-pane>

        <el-tab-pane :label="`REVIEWS (${totalReviews})`" name="reviews">
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
                ;
              >
                <span style="font-size: 28px; font-weight: 600"
                  >4.6<span style="font-size: 24px"> out of 5</span></span
                >
                <el-rate
                  v-model="rating"
                  disabled
                  style="--el-rate-fill-color: var(--main-color); --el-rate-icon-size: 28px"
                  size="large"
                />
              </el-col>
              <el-col :span="18">
                <button class="rating-filter-btn active">All</button>
                <button class="rating-filter-btn">5 Stars (24)</button>
                <button class="rating-filter-btn">4 Stars (24)</button>
                <button class="rating-filter-btn">3 Stars (24)</button>
                <button class="rating-filter-btn">2 Stars (24)</button>
                <button class="rating-filter-btn">1 Star (24)</button>
                <button class="rating-filter-btn">Has Comment (24)</button>
                <button class="rating-filter-btn">Has Pictures (24)</button>
              </el-col>
            </el-row>
          </div>
          <UserComment
            rating="4"
            :image-urls="imageUrls"
            seller-response="Thank you for sharing a review of your recent experience. Our team at Lovito is committed
          to providing high-quality products at affordable prices, and we are delighted to know that
          we made a positive impression on you. We hope you will visit us again soon."
          />
          <UserComment
            rating="5"
            seller-response="Thank you for sharing a review of your recent experience. Our team at Lovito is committed
          to providing high-quality products at affordable prices, and we are delighted to know that
          we made a positive impression on you. We hope you will visit us again soon."
          />
          <UserComment rating="3" :image-urls="imageUrls" />

          <div style="display: flex">
            <el-pagination
              background
              layout="prev, pager, next"
              :total="1000"
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

      <el-row :gutter="20">
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
          />
        </el-col>
      </el-row>
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
        <SplideSlide v-for="item in productList" :key="item.name">
          <ProductItem
            :image-url="item.imageUrl"
            :name="item.name"
            :price="item.price"
            :rating="item.rating"
            :location="item.location"
            :discount="item.discount"
          />
        </SplideSlide>
      </Splide>
    </div>
  </div>
</template>

<style lang="css" scoped>
.option-wrapper {
  display: flex;
  align-items: center;
  height: 48px;
  border: 1px solid black;
  border-radius: 4px;
  padding: 2px 6px;
  margin-right: 12px;

  &:hover {
    cursor: pointer;
    color: #22c55e;
    border-color: currentColor;
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
