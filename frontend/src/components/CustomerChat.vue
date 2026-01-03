<script setup lang="ts">
import { ref, computed } from 'vue'
import { ElMessageBox } from 'element-plus'
import {
  Search,
  Close,
  ChatDotRound,
  ArrowDown,
  ArrowLeft,
  ArrowRight,
  MoreFilled,
  Star,
  BellFilled,
  MuteNotification,
  Bell,
  Delete,
  Notification,
  Finished,
  Picture,
  VideoCamera,
} from '@element-plus/icons-vue'
import PinIcon from './icons/PinIcon.vue'
import UnpinIcon from './icons/UnpinIcon.vue'

const props = defineProps({
  isEmbedded: {
    type: Boolean,
    default: false,
  },
})

const isChatVisible = ref(props.isEmbedded)
const searchQuery = ref('')
const filters = ['All', 'Unread', 'Pinned']
const activeFilter = ref('All')
const imageInputRef = ref<HTMLInputElement | null>(null)
const videoInputRef = ref<HTMLInputElement | null>(null)
const isContentHidden = ref(false)

const toggleChat = () => {
  if (props.isEmbedded) return
  isChatVisible.value = !isChatVisible.value
}

const toggleContent = () => {
  isContentHidden.value = !isContentHidden.value
}

const contacts = ref([
  {
    id: 1,
    name: 'Jean.one',
    lastMessage: 'Okay, I will ship it today.',
    time: '13:12',
    unread: 2,
    avatar: 'JO',
    isPinned: true,
    isRead: false,
    isMuted: false,
    messages: [
      {
        id: 1,
        text: 'Hello, is this item available?',
        sender: 'receiver',
        time: '10:05',
        date: 'Yesterday',
      },
      {
        id: 2,
        text: 'Yes, it is! Would you like to order?',
        sender: 'sender',
        time: '10:10',
        date: 'Yesterday',
      },
      {
        id: 3,
        text: 'Hello shop, can you ship this order as soon as possible for me?',
        sender: 'receiver',
        time: '13:12',
        date: 'Today',
      },
      {
        id: 4,
        text: 'Okay, I will ship it today.',
        sender: 'sender',
        time: '13:15',
        date: 'Today',
      },
      { id: 5, text: 'Thank you!', sender: 'receiver', time: '13:20', date: 'Today' },
    ],
  },
  {
    id: 2,
    name: 'Global Tech Store',
    lastMessage: 'Your order has been shipped.',
    time: 'Yesterday',
    unread: 0,
    avatar: 'GT',
    isPinned: false,
    isRead: true,
    isMuted: false,
    messages: [
      {
        id: 1,
        text: 'Your order has been shipped.',
        sender: 'sender',
        time: '09:00',
        date: 'Yesterday',
      },
    ],
  },
  {
    id: 3,
    name: 'ElectroHub',
    lastMessage: 'Thank you for your purchase!',
    time: '20/12/25',
    unread: 0,
    avatar: 'EH',
    isPinned: false,
    isRead: true,
    isMuted: false,
    messages: [
      {
        id: 1,
        text: 'Thank you for your purchase!',
        sender: 'sender',
        time: '15:30',
        date: '20/12/25',
      },
    ],
  },
])

const filteredContacts = computed(() => {
  return contacts.value.filter((contact) => {
    const matchesSearch = contact.name.toLowerCase().includes(searchQuery.value.toLowerCase())
    const matchesFilter =
      activeFilter.value === 'All' ||
      (activeFilter.value === 'Unread' && contact.unread > 0) ||
      (activeFilter.value === 'Pinned' && contact.isPinned)

    return matchesSearch && matchesFilter
  })
})

const activeContact = ref<(typeof contacts.value)[0] | null | undefined>(contacts.value[0])

const selectContact = (contact: any) => {
  activeContact.value = contact
}

const toggleReadStatus = () => {
  if (activeContact.value) {
    activeContact.value.isRead = !activeContact.value.isRead
    activeContact.value.unread = activeContact.value.isRead ? 0 : 1
  }
}

const togglePin = () => {
  if (activeContact.value) {
    activeContact.value.isPinned = !activeContact.value.isPinned
  }
}

const toggleMute = () => {
  if (activeContact.value) {
    activeContact.value.isMuted = !activeContact.value.isMuted
  }
}

const deleteConversation = async () => {
  if (!activeContact.value) return

  try {
    await ElMessageBox.confirm(
      `Are you sure you want to delete your conversation with ${activeContact.value.name}?`,
      'Delete Conversation',
      {
        confirmButtonText: 'Delete',
        cancelButtonText: 'Cancel',
        type: 'warning',
        confirmButtonClass: 'delete-confirm-btn',
        lockScroll: false,
      },
    )

    // Find index of contact to delete
    const index = contacts.value.findIndex((c) => c.id === activeContact.value?.id)
    if (index !== -1) {
      contacts.value.splice(index, 1)

      // Select another contact if available
      if (contacts.value.length > 0) {
        activeContact.value = contacts.value[0]
      } else {
        activeContact.value = null
      }
    }
  } catch {
    // User cancelled
  }
}

const handleImageUpload = () => {
  imageInputRef.value?.click()
}

const handleVideoUpload = () => {
  videoInputRef.value?.click()
}

const onImageSelected = (event: Event) => {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  if (file) {
    // Handle image file upload logic here
    console.log('Image selected:', file.name)
  }
}

const onVideoSelected = (event: Event) => {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  if (file) {
    // Handle video file upload logic here
    console.log('Video selected:', file.name)
  }
}

const getMsgClass = (msg: any) => {
  if (props.isEmbedded) {
    return msg.sender === 'sender' ? 'receiver' : 'sender'
  }
  return msg.sender
}
</script>

<template>
  <div class="customer-chat-container" :class="{ 'is-embedded': isEmbedded }">
    <!-- Floating Button -->
    <button
      v-if="!isEmbedded"
      class="chat-floating-btn"
      :class="{ 'btn-hidden': isChatVisible }"
      @click="toggleChat"
    >
      <el-icon class="chat-icon"><ChatDotRound /></el-icon>
      <span class="chat-text">Chat</span>
      <span class="unread-badge">5</span>
    </button>

    <!-- Chat Window -->
    <div
      v-if="isChatVisible"
      class="chat-window box-shadow"
      :class="{ 'content-hidden': isContentHidden }"
    >
      <div class="chat-header">
        <div class="header-info">
          <span class="header-title">Chat (5)</span>
        </div>
        <div class="header-actions">
          <el-icon v-if="!isEmbedded" class="action-icon" @click="toggleContent">
            <component :is="isContentHidden ? ArrowLeft : ArrowRight" />
          </el-icon>
          <el-icon v-if="!isEmbedded" class="action-icon" @click="toggleChat"><Close /></el-icon>
        </div>
      </div>

      <div class="chat-body">
        <!-- Contact List -->
        <div class="contact-list">
          <div class="search-filter-area">
            <div class="search-bar">
              <el-input v-model="searchQuery" placeholder="Search" :prefix-icon="Search" />
            </div>
            <el-dropdown trigger="click">
              <span class="filter-trigger">
                {{ activeFilter }} <el-icon class="el-icon--right"><ArrowDown /></el-icon>
              </span>
              <template #dropdown>
                <el-dropdown-menu class="customer-chat-filter-dropdown">
                  <el-dropdown-item
                    v-for="filter in filters"
                    :key="filter"
                    @click="activeFilter = filter"
                  >
                    {{ filter }}
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
          <div class="contacts-scroll">
            <div
              v-for="contact in filteredContacts"
              :key="contact.id"
              class="contact-item"
              :class="{ active: activeContact?.id === contact.id }"
              @click="selectContact(contact)"
            >
              <div class="avatar">{{ contact.avatar }}</div>
              <div class="contact-info">
                <div class="name-time">
                  <span class="name one-line-ellipsis">{{ contact.name }}</span>
                  <span class="time">{{ contact.time }}</span>
                </div>
                <div class="last-msg-unread">
                  <span class="last-msg one-line-ellipsis">{{ contact.lastMessage }}</span>
                  <div class="status-indicators">
                    <el-icon v-if="contact.isPinned" size="large" class="pin-icon">
                      <PinIcon />
                    </el-icon>
                    <el-icon
                      size="large"
                      v-if="contact.isMuted"
                      class="mute-icon"
                      style="color: var(--main-color)"
                    >
                      <MuteNotification
                    /></el-icon>
                    <span v-if="contact.unread > 0" class="unread-count">{{
                      contact.unread > 9 ? '9+' : contact.unread
                    }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Chat Area -->
        <div v-if="!isContentHidden" class="chat-area">
          <div v-if="activeContact" class="chat-content">
            <div class="active-contact-header">
              <span class="active-name">{{ activeContact.name }}</span>
              <el-dropdown trigger="click">
                <el-icon class="contact-menu-icon"><MoreFilled /></el-icon>
                <template #dropdown>
                  <el-dropdown-menu class="contact-dropdown-menu">
                    <el-dropdown-item @click="toggleReadStatus">
                      <el-icon v-if="activeContact.isRead"><Notification /></el-icon>
                      <el-icon v-else><Finished /></el-icon>
                      <span>{{ activeContact.isRead ? 'Mark as unread' : 'Mark as read' }}</span>
                    </el-dropdown-item>
                    <el-dropdown-item @click="togglePin">
                      <el-icon v-show="activeContact.isPinned"><UnpinIcon /></el-icon>
                      <el-icon v-show="!activeContact.isPinned"><PinIcon /></el-icon>
                      <span>{{
                        activeContact.isPinned ? 'Unpin conversation' : 'Pin conversation'
                      }}</span>
                    </el-dropdown-item>
                    <el-dropdown-item @click="toggleMute">
                      <el-icon v-show="activeContact.isMuted"><Bell /></el-icon>
                      <el-icon v-show="!activeContact.isMuted"><MuteNotification /></el-icon>
                      <span>{{
                        activeContact.isMuted ? 'Unmute notifications' : 'Mute notifications'
                      }}</span>
                    </el-dropdown-item>
                    <el-dropdown-item divided @click="deleteConversation">
                      <el-icon><Delete /></el-icon>
                      <span>Delete conversation</span>
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
            <div class="messages-container">
              <div class="warning-box">
                <el-icon class="warning-icon"><ChatDotRound /></el-icon>
                <div class="warning-text">
                  <span v-if="isEmbedded">
                    NOTE: Please answer customers politely and avoid conducting transactions outside
                    the platform to ensure your safety and account status.
                  </span>
                  <span v-else>
                    NOTE: Swoo does NOT allow: Deposit/Bank Transfer outside the system, providing
                    personal phone numbers, or other external activities. Please stay alert to avoid
                    scams!
                  </span>
                </div>
              </div>

              <div class="chat-messages">
                <template v-for="(msg, index) in activeContact.messages" :key="msg.id">
                  <div
                    v-if="index === 0 || msg.date !== activeContact.messages?.[index - 1]?.date"
                    class="date-separator"
                  >
                    <span>{{ msg.date }}</span>
                  </div>

                  <!-- Show order card after Yesterday separator -->
                  <div v-if="msg.date === 'Yesterday' && index === 0" class="order-info-card">
                    <div class="order-info-header">
                      <span v-if="isEmbedded">Customer is chatting with you about this order</span>
                      <span v-else>You are chatting with the Seller about this order</span>
                    </div>
                    <div class="order-info-body">
                      <img src="../assets/product-imgs/product1.png" class="order-item-img" />
                      <div class="order-details">
                        <div class="order-id">Order ID: 2506190B9AHV4B</div>
                        <div class="order-total">Total Order: 161,784₫</div>
                        <div class="order-status">Completed</div>
                      </div>
                    </div>
                  </div>

                  <div class="msg-row" :class="getMsgClass(msg)">
                    <div class="msg-bubble">
                      {{ msg.text }}
                      <span class="msg-time">{{ msg.time }}</span>
                    </div>
                  </div>
                </template>
              </div>
            </div>
            <div class="input-area">
              <input
                ref="imageInputRef"
                type="file"
                accept="image/*"
                style="display: none"
                @change="onImageSelected"
              />
              <input
                ref="videoInputRef"
                type="file"
                accept="video/*"
                style="display: none"
                @change="onVideoSelected"
              />
              <input type="text" placeholder="Type a message..." class="msg-input" />

              <button class="upload-btn" @click="handleImageUpload" title="Upload Image">
                <el-icon><Picture /></el-icon>
              </button>
              <button class="upload-btn" @click="handleVideoUpload" title="Upload Video">
                <el-icon><VideoCamera /></el-icon>
              </button>
              <button class="send-btn">
                <svg viewBox="0 0 24 24" width="20" height="20" fill="currentColor">
                  <path d="M2.01 21L23 12 2.01 3 2 10l15 2-15 2z" />
                </svg>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.customer-chat-container {
  position: fixed;
  bottom: 20px;
  right: 20px;
  z-index: 1500;
  font-family: inherit;
}

.customer-chat-container.is-embedded {
  position: relative;
  bottom: 0;
  right: 0;
  z-index: 1;
  width: 100%;
  height: 100%;
}

.chat-floating-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  background-color: #fff;
  border: 1px solid #e5e5e5;
  border-radius: 4px 4px 0 0;
  padding: 12px 16px;
  box-shadow: 0 -2px 10px rgba(0, 0, 0, 0.1);
  color: var(--main-color);
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: absolute;
  bottom: 0;
  right: 0;
  white-space: nowrap;
}

.chat-floating-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.15);
}

.btn-hidden {
  opacity: 0;
  pointer-events: none;
  transform: scale(0.8);
}

.chat-icon {
  font-size: 20px;
}

.chat-text {
  font-size: 14px;
}

.unread-badge {
  background-color: var(--main-color);
  color: white;
  font-size: 10px;
  padding: 3px 6px;
  border-radius: 10px;
  min-width: 18px;
  text-align: center;
}

.chat-window {
  width: 600px;
  height: 450px;
  background: white;
  border-radius: 8px 8px 0 0;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  border: 1px solid #e5e5e5;
  animation: slideIn 0.3s ease-out;
  transition: width 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.is-embedded .chat-window {
  width: 100%;
  height: 100%;
  border-radius: 8px;
  animation: none;
}

.chat-window.content-hidden {
  width: 220px;
}

@keyframes slideIn {
  from {
    transform: translateY(20px);
    opacity: 0;
  }
  to {
    transform: translateY(0);
    opacity: 1;
  }
}

.chat-header {
  background-color: white;
  border-bottom: 1px solid #f0f0f0;
  padding: 10px 16px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-title {
  color: var(--main-color);
  font-weight: 600;
  font-size: 16px;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.action-icon {
  cursor: pointer;
  color: #999;
  font-size: 18px;
  transition: color 0.2s;
}

.action-icon:hover {
  color: var(--main-color);
}

.chat-body {
  display: flex;
  flex: 1;
  overflow: hidden;
}

.contact-list {
  width: 220px;
  border-right: 1px solid #f0f0f0;
  display: flex;
  flex-direction: column;
  background: #fff;
}

.search-filter-area {
  padding: 10px;
  border-bottom: 1px solid #f0f0f0;
  display: flex;
  align-items: center;
  gap: 8px;
}

.search-bar {
  flex: 1;
}

.filter-trigger {
  font-size: 13px;
  color: #666;
  cursor: pointer;
  display: flex;
  align-items: center;
  white-space: nowrap;
}

.filter-trigger:hover {
  color: var(--main-color);
}

.contacts-scroll {
  flex: 1;
  overflow-y: auto;
}

.contact-item {
  display: flex;
  gap: 10px;
  padding: 12px 10px;
  cursor: pointer;
  transition: background 0.2s;
  align-items: center;
}

.contact-item:hover {
  background-color: #f5f5f5;
}

.contact-item.active {
  background-color: #f0fdf4;
}

.avatar {
  width: 40px;
  height: 40px;
  background: var(--main-color);
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 14px;
  flex-shrink: 0;
}

.contact-info {
  flex: 1;
  min-width: 0;
}

.name-time {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 4px;
}

.name {
  font-size: 13px;
  font-weight: 500;
  color: #333;
}

.time {
  font-size: 11px;
  color: #999;
  margin-left: 2px;
}

.last-msg-unread {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 4px;
}

.status-indicators {
  display: flex;
  align-items: center;
  gap: 4px;
  flex-shrink: 0;
}

.last-msg {
  font-size: 12px;
  color: #999;
}

.unread-count {
  background-color: var(--main-color);
  color: white;
  font-size: 10px;
  min-width: 16px;
  height: 16px;
  padding: 0 4px;
  border-radius: 9px;
  text-align: center;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.pin-icon,
.mute-icon {
  font-size: 12px;
  color: #999;
  flex-shrink: 0;
}

.pin-icon {
  color: var(--main-color);
}

.mute-icon {
  color: #999;
}

.chat-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  background-color: #f9f9f9;
}

.chat-content {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.active-contact-header {
  padding: 10px 16px;
  background: white;
  border-bottom: 1px solid #f0f0f0;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.active-name {
  font-weight: 500;
  font-size: 14px;
}

.messages-container {
  flex: 1;
  padding: 10px 16px;
  display: flex;
  flex-direction: column;
  overflow-y: auto;
  gap: 15px;
}

.warning-box {
  background-color: #fffbef;
  border: 1px solid #ffe58f;
  padding: 10px;
  border-radius: 4px;
  display: flex;
  gap: 10px;
  font-size: 11px;
  color: #856404;
}

.warning-icon {
  font-size: 16px;
  flex-shrink: 0;
  margin-top: 2px;
}

.order-info-card {
  background: white;
  border: 1px solid #f0f0f0;
  border-radius: 4px;
  overflow: hidden;
  flex-shrink: 0;
  margin-bottom: 5px;
}

.order-info-header {
  background-color: #f0fdf4;
  padding: 6px 12px;
  border-bottom: 1px solid #f0f0f0;
  font-size: 11px;
  color: #999;
}

.is-embedded .order-info-header {
  background-color: white;
}

.order-info-body {
  background-color: #f0fdf4;
  padding: 8px 10px;
  display: flex;
  gap: 12px;
  min-height: 60px;
  align-items: center;
}

.is-embedded .order-info-body {
  background-color: white;
}

.order-item-img {
  width: 40px;
  height: 40px;
  object-fit: cover;
  border-radius: 2px;
}

.order-details {
  flex: 1;
  font-size: 12px;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.order-id {
  color: #333;
  margin-bottom: 2px;
}

.order-total {
  color: #999;
  margin-bottom: 2px;
}

.order-status {
  color: var(--main-color);
  font-weight: 500;
}

.chat-messages {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.date-separator {
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 10px 0;
  position: relative;
}

.date-separator::before {
  content: '';
  position: absolute;
  left: 0;
  right: 0;
  top: 50%;
  height: 1px;
  background-color: #f0f0f0;
  z-index: 1;
}

.date-separator span {
  background-color: #f9f9f9;
  padding: 0 12px;
  font-size: 11px;
  color: #999;
  position: relative;
  z-index: 2;
}

.msg-row {
  display: flex;
}

.msg-row.sender {
  justify-content: flex-end;
}

.msg-bubble {
  max-width: 80%;
  padding: 8px 12px;
  border-radius: 4px;
  font-size: 13px;
  position: relative;
  padding-bottom: 18px;
}

.receiver .msg-bubble {
  background-color: #f0fdf4;
  border: 1px solid #e5e5e5;
  color: #333;
}

.sender .msg-bubble {
  background-color: white;
  border: 1px solid #d1fae5;
  color: #333;
}

.msg-time {
  position: absolute;
  bottom: 2px;
  right: 8px;
  font-size: 10px;
  color: #999;
}

.input-area {
  padding: 12px 16px;
  background: white;
  border-top: 1px solid #f0f0f0;
  display: flex;
  gap: 10px;
  align-items: center;
}

.upload-btn {
  color: #666;
  background: transparent;
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 6px;
  cursor: pointer;
  transition: all 0.2s;
  border-radius: 4px;
}

.upload-btn:hover {
  color: var(--main-color);
  background-color: #f0fdf4;
}

.upload-btn .el-icon {
  font-size: 18px;
}

.msg-input {
  flex: 1;
  border: 1px solid #e5e5e5;
  border-radius: 4px;
  padding: 8px 12px;
  font-size: 13px;
  outline: none;
}

.msg-input:focus {
  border-color: var(--main-color);
}

.send-btn {
  color: var(--main-color);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 4px;
  transition: transform 0.2s;
}

.send-btn:hover {
  transform: scale(1.1);
}

.one-line-ellipsis {
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.el-dropdown-menu__item {
  --padding: 0 16px !important;
}

.contact-menu-icon {
  cursor: pointer;
  color: #999;
  font-size: 18px;
  transition: color 0.2s;
}

.contact-menu-icon:hover {
  color: var(--main-color);
}

.contact-dropdown-menu .el-dropdown-menu__item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  font-size: 13px;
}

.contact-dropdown-menu .el-dropdown-menu__item .el-icon {
  font-size: 16px;
  color: #666;
}

.contact-dropdown-menu .el-dropdown-menu__item:hover .el-icon {
  color: var(--main-color);
}
</style>

<style>
.delete-confirm-btn {
  background-color: #f56c6c !important;
  border-color: #f56c6c !important;
}

.delete-confirm-btn:hover {
  background-color: #f78989 !important;
  border-color: #f78989 !important;
}
</style>
