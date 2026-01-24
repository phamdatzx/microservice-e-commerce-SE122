<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessageBox, ElMessage } from 'element-plus'
import axios from 'axios'
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
  Loading,
} from '@element-plus/icons-vue'
import PinIcon from './icons/PinIcon.vue'
import UnpinIcon from './icons/UnpinIcon.vue'
import { socketService, SOCKET_EVENTS } from '@/utils/socket'

const props = defineProps({
  isEmbedded: {
    type: Boolean,
    default: false,
  },
})

const route = useRoute()

const isChatVisible = ref(props.isEmbedded)
const searchQuery = ref('')
const filters = ['All', 'Unread', 'Pinned']
const activeFilter = ref('All')
const imageInputRef = ref<HTMLInputElement | null>(null)
const isContentHidden = ref(false)
const contacts = ref<any[]>([])
const activeContact = ref<any>(null)
const newMessage = ref('')
const otherUserTyping = ref(false)
const messagesContainerRef = ref<HTMLElement | null>(null)
const selectedImageFile = ref<File | null>(null)
const imagePreviewUrl = ref<string | null>(null)
const isSending = ref(false)
const isLoadingContacts = ref(false)
const isLoadingMessages = ref(false)

const fetchConversations = async () => {
  const token = localStorage.getItem('access_token')
  if (!token) return

  isLoadingContacts.value = true
  try {
    const response = await axios.get(
      `${import.meta.env.VITE_BE_API_URL}/notification/chat/conversations`,
      {
        params: { limit: 1000 },
        headers: { Authorization: `Bearer ${token}` },
      },
    )

    if (response.data && response.data.data) {
      contacts.value = response.data.data.map((conv: any) => {
        const otherParty = props.isEmbedded ? conv.customer : conv.seller
        return {
          id: conv._id,
          name: otherParty.name || otherParty.username,
          lastMessage: conv.lastMessage || 'No messages yet',
          time: formatTime(conv.lastUpdated),
          unread: conv.unreadCount || 0,
          avatarUrl: otherParty.image,
          avatarInitials: getInitials(otherParty.name || otherParty.username),
          isPinned: false,
          isRead: conv.unreadCount === 0,
          isMuted: false,
          messages: [],
          order: conv.order || null,
          raw: conv,
        }
      })
      hydrateLastMessagePrefix()
    }
  } catch (error) {
    console.error('Error fetching conversations:', error)
  } finally {
    isLoadingContacts.value = false
  }
}

const openChatWithSeller = async (sellerId: string) => {
  if (!sellerId) return
  isChatVisible.value = true
  isContentHidden.value = false

  const existing = contacts.value.find((c) => {
    const targetId = props.isEmbedded ? c.raw.customerId : c.raw.sellerId
    return String(targetId) === String(sellerId)
  })

  if (existing) {
    selectContact(existing)
    return
  }

  const token = localStorage.getItem('access_token')
  if (!token) return

  try {
    const response = await axios.post(
      `${import.meta.env.VITE_BE_API_URL}/notification/chat/conversations`,
      { sellerId },
      {
        headers: { Authorization: `Bearer ${token}` },
      },
    )

    if (response.data && response.data.data) {
      const conv = response.data.data
      const otherParty = props.isEmbedded ? conv.customer : conv.seller
      const newContact = {
        id: conv._id,
        name: otherParty.name || otherParty.username,
        lastMessage: conv.lastMessage || 'No messages yet',
        time: formatTime(conv.lastUpdated),
        unread: conv.unreadCount || 0,
        avatarUrl: otherParty.image,
        avatarInitials: getInitials(otherParty.name || otherParty.username),
        isPinned: false,
        isRead: conv.unreadCount === 0,
        isMuted: false,
        messages: [],
        order: conv.order || null,
        raw: conv,
      }
      contacts.value.unshift(newContact)
      selectContact(newContact)
    }
  } catch (error) {
    console.error('Error starting new conversation:', error)
    ElMessage.error('Failed to start chat')
  }
}

const handleOpenChatEvent = (event: any) => {
  if (event.detail && event.detail.sellerId) {
    openChatWithSeller(event.detail.sellerId)
  }
}

const scrollToBottom = async () => {
  await nextTick()
  if (messagesContainerRef.value) {
    messagesContainerRef.value.scrollTop = messagesContainerRef.value.scrollHeight
  }
}

const formatMessage = (msg: any, contact: any) => {
  const isMe =
    String(msg.senderId) === String(props.isEmbedded ? contact.raw.sellerId : contact.raw.userId)
  return {
    id: msg._id,
    text: msg.content === '[Image]' ? '' : msg.content,
    sender: isMe ? 'receiver' : 'sender',
    time: new Date(msg.createdAt).toLocaleTimeString('en-US', {
      hour: '2-digit',
      minute: '2-digit',
      hour12: false,
    }),
    date: formatDateSeparator(msg.createdAt),
    image: msg.image,
  }
}

const markAsReadAPI = async (conversationId: string) => {
  const token = localStorage.getItem('access_token')
  if (!token) return

  try {
    await axios.patch(
      `${import.meta.env.VITE_BE_API_URL}/notification/chat/conversations/${conversationId}/read`,
      {},
      {
        headers: { Authorization: `Bearer ${token}` },
      },
    )
    window.dispatchEvent(new CustomEvent('trigger-chat-badge-update'))
  } catch (error) {
    console.error('Error marking conversation as read:', error)
  }
}

const updateContactInList = (message: any) => {
  const index = contacts.value.findIndex((c) => String(c.id) === String(message.conversationId))
  if (index !== -1) {
    const contact = contacts.value[index]

    // Check if the message was sent by the current user
    const currentUserId = props.isEmbedded ? contact.raw.sellerId : contact.raw.userId
    const isMe = String(message.senderId) === String(currentUserId)
    const prefix = isMe ? 'You: ' : ''

    contact.lastMessage = prefix + (message.content || '')
    if (message.image && (!contact.lastMessage || contact.lastMessage.trim() === '')) {
      // Or if content is empty but has image
      // If content is empty/null but has image, or if just image
      if (!message.content || message.content.trim() === '' || message.content === '[Image]') {
        contact.lastMessage = prefix + '[Image]'
      } else {
        contact.lastMessage = prefix + message.content
      }
    } else if (message.content === '[Image]') {
      contact.lastMessage = prefix + '[Image]'
    }

    contact.time = formatTime(message.createdAt)
    contact.raw.lastUpdated = message.createdAt

    contacts.value.splice(index, 1)
    contacts.value.unshift(contact)
    return contact
  }
  return null
}

const fetchMessages = async (conversationId: string) => {
  const token = localStorage.getItem('access_token')
  if (!token) return

  isLoadingMessages.value = true
  try {
    const response = await axios.get(
      `${import.meta.env.VITE_BE_API_URL}/notification/chat/conversations/${conversationId}/messages`,
      {
        params: { limit: 1000 },
        headers: { Authorization: `Bearer ${token}` },
      },
    )

    if (response.data && response.data.data) {
      const contact = contacts.value.find((c) => c.id === conversationId)
      if (contact) {
        contact.messages = response.data.data.map((msg: any) => formatMessage(msg, contact))
        scrollToBottom()
      }
    }
  } catch (error) {
    console.error('Error fetching messages:', error)
  } finally {
    isLoadingMessages.value = false
  }
}

const hydrateLastMessagePrefix = async () => {
  const token = localStorage.getItem('access_token')
  if (!token) return

  // Limit concurrency if needed, but for now simple loop
  for (const contact of contacts.value) {
    try {
      const response = await axios.get(
        `${import.meta.env.VITE_BE_API_URL}/notification/chat/conversations/${contact.id}/messages`,
        {
          params: { limit: 1 },
          headers: { Authorization: `Bearer ${token}` },
        },
      )
      if (response.data && response.data.data && response.data.data.length > 0) {
        const lastMsg = response.data.data[0]
        const currentUserId = props.isEmbedded ? contact.raw.sellerId : contact.raw.userId
        const isMe = String(lastMsg.senderId) === String(currentUserId)

        if (isMe) {
          const content =
            lastMsg.content === '[Image]' || (!lastMsg.content && lastMsg.image)
              ? '[Image]'
              : lastMsg.content
          contact.lastMessage = `You: ${content}`
        }
      }
    } catch (e) {
      // ignore error for hydration
    }
  }
}

const formatDateSeparator = (dateStr: string) => {
  const date = new Date(dateStr)
  const now = new Date()
  const yesterday = new Date(now)
  yesterday.setDate(now.getDate() - 1)

  if (date.toDateString() === now.toDateString()) return 'Today'
  if (date.toDateString() === yesterday.toDateString()) return 'Yesterday'
  return date.toLocaleDateString('en-US', { day: 'numeric', month: 'short' })
}

const formatTime = (dateStr: string) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  const now = new Date()

  if (date.toDateString() === now.toDateString()) {
    return date.toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit', hour12: false })
  }
  return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
}

const getInitials = (name: string) => {
  return name
    ? name
        .split(' ')
        .map((n) => n[0])
        .join('')
        .toUpperCase()
        .slice(0, 2)
    : '?'
}

const filteredContacts = computed(() => {
  const filtered = contacts.value.filter((contact) => {
    const matchesSearch = contact.name.toLowerCase().includes(searchQuery.value.toLowerCase())
    const matchesFilter =
      activeFilter.value === 'All' ||
      (activeFilter.value === 'Unread' && contact.unread > 0) ||
      (activeFilter.value === 'Pinned' && contact.isPinned)

    return matchesSearch && matchesFilter
  })

  return [...filtered].sort((a, b) => {
    const dateA = new Date(a.raw.lastUpdated || 0).getTime()
    const dateB = new Date(b.raw.lastUpdated || 0).getTime()
    return dateB - dateA
  })
})

const totalUnreadCount = computed(() => {
  return contacts.value.reduce((total, contact) => total + contact.unread, 0)
})

onMounted(() => {
  fetchConversations()

  const token = localStorage.getItem('access_token')
  if (token) {
    socketService.connect(token)
    setupSocketListeners()
  }

  window.addEventListener('open-chat', handleOpenChatEvent)

  // Handle focus query param for sellers or direct links
  const focusId = route.query.focus
  if (focusId) {
    // Wait for conversations to load first
    watch(
      isLoadingContacts,
      (loading) => {
        if (!loading && focusId) {
          openChatWithSeller(focusId as string)
        }
      },
      { immediate: true },
    )
  }
})

onUnmounted(() => {
  socketService.disconnect()
  window.removeEventListener('open-chat', handleOpenChatEvent)
})

const setupSocketListeners = () => {
  socketService.on(SOCKET_EVENTS.NEW_MESSAGE, (message: any) => {
    const existingContact = contacts.value.find(
      (c) => String(c.id) === String(message.conversationId),
    )
    if (existingContact) {
      const isDuplicate = existingContact.messages.some(
        (m: any) =>
          String(m.id) === String(message._id) ||
          (m.text === message.content &&
            m.image === message.image &&
            Math.abs(new Date(m.createdAt).getTime() - new Date(message.createdAt).getTime()) <
              5000),
      )
      if (isDuplicate) return
    }

    const updatedContact = updateContactInList(message)
    if (updatedContact) {
      const formattedMsg = formatMessage(message, updatedContact)
      updatedContact.messages.push(formattedMsg)

      if (
        activeContact.value &&
        String(activeContact.value.id) === String(message.conversationId)
      ) {
        scrollToBottom()
      }

      // Always increment unread count for new messages, even if active
      // User must click to mark as read
      updatedContact.unread++
      updatedContact.isRead = false
    } else {
      fetchConversations()
    }
  })

  socketService.on(SOCKET_EVENTS.MESSAGE_SENT, (message: any) => {
    const contact = updateContactInList(message)
    if (contact) {
      const isDuplicate = contact.messages.some(
        (m: any) =>
          String(m.id) === String(message._id) ||
          (m.text === message.content &&
            m.image === message.image &&
            Math.abs(new Date(m.createdAt).getTime() - new Date(message.createdAt).getTime()) <
              5000),
      )

      if (!isDuplicate) {
        const formattedMsg = formatMessage(message, contact)
        contact.messages.push(formattedMsg)
        scrollToBottom()
      }
    }
  })

  socketService.on(SOCKET_EVENTS.USER_TYPING, (data: any) => {
    if (activeContact.value?.id === data.conversationId) {
      otherUserTyping.value = data.isTyping
    }
  })

  socketService.on(SOCKET_EVENTS.MESSAGES_UPDATED, (data: any) => {
    const contact = contacts.value.find((c) => c.id === data.conversationId)
    if (contact) {
      contact.unread = 0
      contact.isRead = true
    }
  })

  socketService.on(SOCKET_EVENTS.NEW_NOTIFICATION, (notification: any) => {
    if (notification.type === 'chat' && notification.data) {
      const messageData = {
        conversationId: notification.data.conversationId,
        content: notification.message,
        createdAt: notification.createdAt,
        image: notification.data.hasImage ? '[Image]' : null,
      }

      const updatedContact = updateContactInList(messageData)
      if (updatedContact) {
        updatedContact.unread++
        updatedContact.isRead = false
      } else {
        fetchConversations()
      }
    }
  })

  socketService.on(SOCKET_EVENTS.ERROR, (error: any) => {
    console.error('Socket error event:', error)
    ElMessage.error(error.message || 'Socket error occurred')
  })
}

watch(activeContact, (newVal, oldVal) => {
  if (oldVal) {
    socketService.emit(SOCKET_EVENTS.LEAVE_CONVERSATION, { conversationId: oldVal.id })
  }
  if (newVal) {
    socketService.emit(SOCKET_EVENTS.JOIN_CONVERSATION, { conversationId: newVal.id })
    // Removed auto-read logic from here
    otherUserTyping.value = false
  }
})

watch(newMessage, (val) => {
  if (activeContact.value) {
    socketService.emit(SOCKET_EVENTS.TYPING, {
      conversationId: activeContact.value.id,
      isTyping: val.length > 0,
    })
  }
})

const toggleChat = () => {
  if (props.isEmbedded) return
  isChatVisible.value = !isChatVisible.value
}

const toggleContent = () => {
  isContentHidden.value = !isContentHidden.value
}

const shouldShowDate = (index: number, messages: any[]) => {
  if (index === 0) return true
  return messages[index].date !== messages[index - 1].date
}

const selectContact = (contact: any) => {
  activeContact.value = contact

  if (isContentHidden.value) {
    isContentHidden.value = false
  }

  // Mark as read when explicitly selected (clicked)
  if (contact) {
    if (contact.unread > 0) {
      contact.unread = 0
      contact.isRead = true
      socketService.emit(SOCKET_EVENTS.MESSAGE_READ, { conversationId: contact.id })
      markAsReadAPI(contact.id)
    }
  }

  if (contact && contact.messages.length === 0) {
    fetchMessages(contact.id)
  } else if (contact) {
    scrollToBottom()
  }
}

const toggleReadStatus = () => {
  if (activeContact.value) {
    activeContact.value.isRead = !activeContact.value.isRead
    activeContact.value.unread = activeContact.value.isRead ? 0 : 1
    if (activeContact.value.isRead) {
      markAsReadAPI(activeContact.value.id)
    }
  }
}

const markActiveConversationAsRead = () => {
  if (activeContact.value && !activeContact.value.isRead) {
    activeContact.value.isRead = true
    activeContact.value.unread = 0
    socketService.emit(SOCKET_EVENTS.MESSAGE_READ, { conversationId: activeContact.value.id })
    markAsReadAPI(activeContact.value.id)
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

const sendMessage = async () => {
  if (
    isSending.value ||
    (!newMessage.value.trim() && !selectedImageFile.value) ||
    !activeContact.value
  ) {
    return
  }

  isSending.value = true
  try {
    if (selectedImageFile.value) {
      await sendImageMessage()
    } else {
      const payload = {
        conversationId: activeContact.value.id,
        content: newMessage.value.trim(),
      }

      socketService.emit(SOCKET_EVENTS.SEND_MESSAGE, payload)
      newMessage.value = ''
    }
  } catch (err) {
    console.error('❌ Error in sendMessage:', err)
  } finally {
    isSending.value = false
  }
}

const sendImageMessage = async () => {
  if (!activeContact.value || !selectedImageFile.value) return

  const token = localStorage.getItem('access_token')
  if (!token) return

  const formData = new FormData()
  formData.append('conversationId', activeContact.value.id)
  formData.append('content', newMessage.value.trim() || '[Image]')
  formData.append('image', selectedImageFile.value)

  try {
    const response = await axios.post(
      `${import.meta.env.VITE_BE_API_URL}/notification/chat/messages/with-image`,
      formData,
      {
        headers: {
          Authorization: `Bearer ${token}`,
          'Content-Type': 'multipart/form-data',
        },
      },
    )

    if (response.data && response.data.success) {
      const message = response.data.data

      socketService.emit(SOCKET_EVENTS.SEND_MESSAGE, {
        conversationId: activeContact.value.id,
        content: message.content,
        image: message.image,
      })

      removeSelectedImage()
      newMessage.value = ''
    }
  } catch (error) {
    console.error('Error sending image message:', error)
    ElMessage.error('Failed to send image')
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

    const index = contacts.value.findIndex((c) => c.id === activeContact.value?.id)
    if (index !== -1) {
      contacts.value.splice(index, 1)

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

const onImageSelected = (event: Event) => {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  if (file) {
    if (file.size > 5 * 1024 * 1024) {
      ElMessage.warning('Image size must be less than 5MB')
      return
    }
    selectedImageFile.value = file
    imagePreviewUrl.value = URL.createObjectURL(file)
  }
}

const removeSelectedImage = () => {
  selectedImageFile.value = null
  if (imagePreviewUrl.value) {
    URL.revokeObjectURL(imagePreviewUrl.value)
    imagePreviewUrl.value = null
  }
  if (imageInputRef.value) {
    imageInputRef.value.value = ''
  }
}

const getMsgClass = (msg: any) => {
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
      <span v-if="totalUnreadCount > 0" class="unread-badge">{{
        totalUnreadCount > 99 ? '99+' : totalUnreadCount
      }}</span>
    </button>

    <!-- Chat Window -->
    <div
      v-if="isChatVisible"
      class="chat-window box-shadow"
      :class="{ 'content-hidden': isContentHidden }"
    >
      <div class="chat-header">
        <div class="header-info">
          <span class="header-title">Chat ({{ totalUnreadCount }})</span>
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
            <div v-if="isLoadingContacts" class="loading-container">
              <el-icon class="is-loading"><Loading /></el-icon>
            </div>
            <template v-else>
              <div
                v-for="contact in filteredContacts"
                :key="contact.id"
                class="contact-item"
                :class="{ active: activeContact?.id === contact.id }"
                @click="selectContact(contact)"
              >
                <div v-if="contact.avatarUrl" class="avatar-img-wrapper">
                  <img :src="contact.avatarUrl" class="avatar-image" />
                </div>
                <div v-else class="avatar">{{ contact.avatarInitials }}</div>
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
            </template>
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
            <div class="messages-container" ref="messagesContainerRef">
              <div v-if="isLoadingMessages" class="loading-container">
                <el-icon class="is-loading"><Loading /></el-icon>
              </div>
              <template v-else>
                <div class="warning-box">
                  <el-icon class="warning-icon"><ChatDotRound /></el-icon>
                  <div class="warning-text">
                    <span v-if="isEmbedded">
                      NOTE: Please answer customers politely and avoid conducting transactions
                      outside the platform to ensure your safety and account status.
                    </span>
                    <span v-else>
                      NOTE: Swoo does NOT allow: Deposit/Bank Transfer outside the system, providing
                      personal phone numbers, or other external activities. Please stay alert to
                      avoid scams!
                    </span>
                  </div>
                </div>

                <div class="chat-messages">
                  <template v-for="(msg, index) in activeContact.messages" :key="msg.id">
                    <div
                      v-if="shouldShowDate(Number(index), activeContact.messages)"
                      class="date-separator"
                    >
                      <span>{{ msg.date }}</span>
                    </div>

                    <div v-if="index === 0 && activeContact.order" class="order-info-card">
                      <div class="order-info-header">
                        <span v-if="isEmbedded"
                          >Customer is chatting with you about this order</span
                        >
                        <span v-else>You are chatting with the Seller about this order</span>
                      </div>
                      <div class="order-info-body">
                        <img
                          :src="
                            activeContact.order.image ||
                            activeContact.order.item_image ||
                            '../assets/product-imgs/product1.png'
                          "
                          class="order-item-img"
                        />
                        <div class="order-details">
                          <div class="order-id">
                            Order ID: {{ activeContact.order.id || activeContact.order._id }}
                          </div>
                          <div class="order-total">
                            Total Order:
                            {{ activeContact.order.total_price || activeContact.order.total }}đ
                          </div>
                          <div class="order-status" style="text-transform: capitalize">
                            {{ activeContact.order.status }}
                          </div>
                        </div>
                      </div>
                    </div>

                    <div class="msg-row" :class="getMsgClass(msg)">
                      <div class="msg-bubble">
                        <div v-if="msg.image" class="msg-image-container">
                          <el-image
                            :src="msg.image"
                            class="msg-image"
                            @load="scrollToBottom"
                            :preview-src-list="[msg.image]"
                            :preview-teleported="true"
                            fit="cover"
                          />
                        </div>
                        <span v-if="msg.text">{{ msg.text }}</span>
                        <span class="msg-time">{{ msg.time }}</span>
                      </div>
                    </div>
                  </template>

                  <div v-if="otherUserTyping" class="msg-row sender">
                    <div class="msg-bubble typing-indicator">
                      <span>.</span><span>.</span><span>.</span>
                    </div>
                  </div>
                </div>
              </template>
            </div>

            <div v-if="imagePreviewUrl" class="image-preview-bar">
              <div class="preview-item">
                <img :src="imagePreviewUrl" class="preview-img" />
                <div class="preview-remove" @click="removeSelectedImage">
                  <el-icon><Close /></el-icon>
                </div>
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
                v-model="newMessage"
                type="text"
                placeholder="Type a message..."
                class="msg-input"
                :disabled="isSending"
                @keyup.enter="sendMessage"
                @click="markActiveConversationAsRead"
                @focus="markActiveConversationAsRead"
              />

              <button
                class="upload-btn"
                :disabled="isSending"
                @click="handleImageUpload"
                title="Upload Image"
              >
                <el-icon><Picture /></el-icon>
              </button>

              <button class="send-btn" :disabled="isSending" @click="sendMessage">
                <el-icon v-if="isSending" class="is-loading"><Loading /></el-icon>
                <svg v-else viewBox="0 0 24 24" width="20" height="20" fill="currentColor">
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

.avatar-img-wrapper {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  overflow: hidden;
  flex-shrink: 0;
}

.avatar-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
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
  min-width: 60px;
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

.msg-input:disabled {
  background-color: #f5f5f5;
  cursor: not-allowed;
}

.send-btn {
  color: var(--main-color);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 4px;
  transition: transform 0.2s;
  background: transparent;
  border: none;
  cursor: pointer;
}

.send-btn:hover:not(:disabled) {
  transform: scale(1.1);
}

.send-btn:disabled,
.upload-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
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

.typing-indicator {
  display: flex;
  gap: 2px;
  padding: 8px 15px !important;
  min-height: 32px;
  align-items: center;
}

.typing-indicator span {
  height: 8px;
  width: 8px;
  background: #999;
  border-radius: 50%;
  display: block;
  animation: typing 1s infinite;
}

.typing-indicator span:nth-child(2) {
  animation-delay: 0.2s;
}

.typing-indicator span:nth-child(3) {
  animation-delay: 0.4s;
}

@keyframes typing {
  0% {
    transform: translateY(0px);
  }
  50% {
    transform: translateY(-5px);
  }
  100% {
    transform: translateY(0px);
  }
}

.msg-image-container {
  margin-bottom: 8px;
  border-radius: 4px;
  overflow: hidden;
  max-width: 180px;
}

.msg-image {
  display: block;
  max-width: 100%;
  max-height: 120px;
  object-fit: contain;
  cursor: pointer;
}

.image-preview-bar {
  padding: 8px 16px;
  background-color: #f9f9f9;
  border-top: 1px solid #f0f0f0;
  display: flex;
  gap: 10px;
}

.preview-item {
  position: relative;
  width: 60px;
  height: 60px;
  border-radius: 4px;
  overflow: hidden;
  border: 1px solid #e5e5e5;
}

.preview-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.preview-remove {
  position: absolute;
  top: 2px;
  right: 2px;
  background: rgba(0, 0, 0, 0.5);
  color: white;
  border-radius: 50%;
  width: 16px;
  height: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  font-size: 10px;
}

.preview-remove:hover {
  background: rgba(0, 0, 0, 0.7);
}

.loading-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 50%;
  width: 100%;
  color: var(--main-color);
  font-size: 24px;
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
