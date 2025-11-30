<template>
  <view class="recipe-container">
    <!-- 搜索栏 -->
    <view class="search-section">
      <uni-easyinput
        v-model="searchText"
        placeholder="搜索菜谱..."
        prefixIcon="search"
        trim="both"
      ></uni-easyinput>
    </view>

    <!-- 菜谱列表 -->
    <view class="recipe-list">
      <view 
        v-for="(recipe, index) in recipeList" 
        :key="index" 
        class="recipe-item"
        @click="viewRecipe(recipe)"
      >
        <view class="recipe-image">
          <image
            :src="recipe.image"
            mode="aspectFill"
            class="image"
          ></image>
        </view>
        <view class="recipe-info">
          <text class="recipe-title">{{ recipe.title }}</text>
          <text class="recipe-description">{{ recipe.description }}</text>
          <view class="recipe-meta">
            <view class="meta-item">
              <uni-icons type="clock" size="14" color="#999"></uni-icons>
              <text class="meta-text">{{ recipe.time }}</text>
            </view>
            <view class="meta-item">
              <uni-icons type="star" size="14" color="#FFD700"></uni-icons>
              <text class="meta-text">{{ recipe.rating }}</text>
            </view>
          </view>
        </view>
        <view class="recipe-actions">
          <view class="action-btn" @click.stop="addToCart(recipe)">
            <uni-icons type="cart" size="20" color="#FF6B9D"></uni-icons>
          </view>
        </view>
      </view>
      
      <!-- 空状态 -->
      <view v-if="recipeList.length === 0" class="empty-state">
        <uni-icons type="folder-add" size="80" color="#ccc"></uni-icons>
        <text class="empty-text">还没有菜谱</text>
        <text class="empty-desc">快来添加第一个菜谱吧</text>
      </view>
    </view>

    <!-- 悬浮添加按钮 -->
    <view class="floating-add-btn" @click="addRecipe">
      <uni-icons type="plus" size="24" color="#fff"></uni-icons>
    </view>
  </view>
</template>

<script setup>
import { ref, computed } from 'vue'

// 响应式数据
const searchText = ref('')
const selectedCategory = ref('all')

const categories = ref([
  { id: 'all', name: '全部' },
  { id: 'chinese', name: '中餐' },
  { id: 'western', name: '西餐' },
  { id: 'dessert', name: '甜点' },
  { id: 'drink', name: '饮品' },
  { id: 'snack', name: '小吃' }
])

const recipes = ref([
  {
    id: 1,
    name: '红烧肉',
    description: '经典美味的红烧肉，肥而不腻',
    image: '/static/images/recipe1.jpg',
    category: 'chinese',
    cookTime: 60,
    difficulty: '中等',
    ingredients: ['五花肉', '生抽', '老抽', '冰糖', '葱姜蒜'],
    steps: ['步骤1', '步骤2', '步骤3']
  },
  {
    id: 2,
    name: '提拉米苏',
    description: '意大利经典咖啡甜点',
    image: '/static/images/recipe2.jpg',
    category: 'dessert',
    cookTime: 30,
    difficulty: '简单',
    ingredients: ['手指饼干', '马斯卡彭奶酪', '咖啡', '鸡蛋', '糖'],
    steps: ['步骤1', '步骤2', '步骤3']
  },
  {
    id: 3,
    name: '意式咖啡',
    description: '浓郁香醇的意式浓缩咖啡',
    image: '/static/images/recipe3.jpg',
    category: 'drink',
    cookTime: 5,
    difficulty: '简单',
    ingredients: ['咖啡豆', '水'],
    steps: ['步骤1', '步骤2']
  }
])

// 计算属性
const filteredRecipes = computed(() => {
  let filtered = recipes.value
  
  // 按分类筛选
  if (selectedCategory.value !== 'all') {
    filtered = filtered.filter(recipe => recipe.category === selectedCategory.value)
  }
  
  // 按关键词搜索
  if (searchText.value) {
    filtered = filtered.filter(recipe => 
      recipe.name.includes(searchText.value) || 
      recipe.description.includes(searchText.value)
    )
  }
  
  return filtered
})

// 方法
const onSearch = () => {
  // 搜索逻辑
  console.log('搜索关键词:', searchText.value)
}

const selectCategory = (categoryId) => {
  selectedCategory.value = categoryId
}

const viewRecipe = (recipe) => {
  uni.navigateTo({
    url: `/pages/recipe/detail?id=${recipe.id}`
  })
}

const orderRecipe = (recipe) => {
  uni.navigateTo({
    url: `/pages/order/order?recipeId=${recipe.id}`
  })
}

const createRecipe = () => {
  uni.navigateTo({
    url: '/pages/recipe/create'
  })
}
</script>

<style lang="scss" scoped>
.recipe-container {
  min-height: 100vh;
  background: #f5f5f5;
  padding-bottom: 120rpx;
}

.search-section {
  background: #fff;
  padding: 20rpx 30rpx;
  
  .search-box {
    background: #f8f8f8;
    border-radius: 20rpx;
    overflow: hidden;
  }
}

.category-section {
  background: #fff;
  padding: 20rpx 0;
  margin-bottom: 20rpx;
  
  .category-scroll {
    white-space: nowrap;
    
    .category-list {
      display: flex;
      padding: 0 30rpx;
      
      .category-item {
        display: inline-block;
        padding: 16rpx 32rpx;
        margin-right: 20rpx;
        background: #f8f8f8;
        border-radius: 25rpx;
        transition: all 0.3s;
        
        .category-name {
          font-size: 28rpx;
          color: #666;
        }
        
        &.active {
          background: #FF6B9D;
          
          .category-name {
            color: #fff;
          }
        }
      }
    }
  }
}

.recipe-list {
  padding: 0 30rpx;
  
  .recipe-item {
    display: flex;
    align-items: center;
    background: #fff;
    padding: 30rpx;
    border-radius: 20rpx;
    margin-bottom: 20rpx;
    box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.05);
    
    .recipe-image {
      margin-right: 30rpx;
      
      .image {
        width: 120rpx;
        height: 120rpx;
        border-radius: 12rpx;
        object-fit: cover;
      }
    }
    
    .recipe-info {
      flex: 1;
      
      .recipe-title {
        display: block;
        font-size: 32rpx;
        font-weight: bold;
        color: #333;
        margin-bottom: 10rpx;
      }
      
      .recipe-description {
        display: block;
        font-size: 26rpx;
        color: #999;
        margin-bottom: 20rpx;
        line-height: 1.4;
      }
      
      .recipe-meta {
        display: flex;
        gap: 30rpx;
        
        .meta-item {
          display: flex;
          align-items: center;
          gap: 8rpx;
          
          .meta-text {
            font-size: 24rpx;
            color: #999;
          }
        }
      }
    }
    
    .recipe-actions {
      .action-btn {
        width: 60rpx;
        height: 60rpx;
        background: #f8f8f8;
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
        transition: all 0.3s;
        cursor: pointer;
        
        &:active {
          background: #FF6B9D;
          transform: scale(0.9);
        }
      }
    }
  }
}

.empty-state {
  text-align: center;
  padding: 100rpx 0;
  
  .empty-text {
    display: block;
    font-size: 32rpx;
    color: #999;
    margin: 30rpx 0 20rpx;
  }
  
  .empty-desc {
    font-size: 26rpx;
    color: #ccc;
  }
}

.floating-add-btn {
  position: fixed;
  right: 40rpx;
  bottom: 120rpx;
  width: 100rpx;
  height: 100rpx;
  background: linear-gradient(135deg, #FF6B9D 0%, #FF8E9E 100%);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 8rpx 30rpx rgba(255, 107, 157, 0.3);
  z-index: 999;
  transition: all 0.3s ease;
  
  &:active {
    transform: scale(0.95);
    box-shadow: 0 4rpx 15rpx rgba(255, 107, 157, 0.4);
  }
  
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0% {
    box-shadow: 0 8rpx 30rpx rgba(255, 107, 157, 0.3);
  }
  50% {
    box-shadow: 0 8rpx 30rpx rgba(255, 107, 157, 0.6);
  }
  100% {
    box-shadow: 0 8rpx 30rpx rgba(255, 107, 157, 0.3);
  }
}
</style> 