# Blueberry Design Document

**An Experimental Neobrutalist Imageboard + Microblogging PWA**

## Overview

Blueberry is a modern web application that combines the best aspects of imageboards (like 4chan) and microblogging platforms (like Twitter/X). It features a neobrutalist design aesthetic with soft pastels and provides comprehensive social interaction capabilities within organized communities called "boards."

## Technology Stack

### Frontend
- **Package Manager**: pnpm
- **Framework**: Vue 3 + Nuxt 4 + TypeScript
- **Styling**: TailwindCSS + Vanilla CSS
- **State Management**: Pinia
- **Build Tool**: Vite

### Backend
- **Language**: Go
- **ORM**: GORM
- **Web Framework**: Echo
- **Password Hashing**: Argon2id
- **Authentication**: JWT
- **Database**: PostgreSQL
- **Containerization**: Docker

### Network

#### IP Addresses

192.168.50.214 - Mono-Repo Development Server

#### Ports
9876 - Frontend Dev Server
5432 - Backend Dev Server
3005 - PostgreSQL Docker Container

## Core Features

- **User Management**: Registration, authentication, profiles with customization
- **Content Creation**: Posts, threads, comments with media support
- **Social Interactions**: Likes, dislikes, reposts, quote reposts
- **Community Organization**: Boards/communities with moderation
- **Reputation System**: "Ink" - a karma-like system preventing gaming
- **Real-time Search**: Accurate, fast search across all content types
- **Private Messaging**: E2E encrypted, P2P using Hypercore protocol
- **Moderation Tools**: Comprehensive role-based moderation system
- **User Customization**: Custom decorative tags and themes
- **Advertising**: Simple, non-invasive advertisement system

---

# Frontend Architecture

## Pages (`/app/pages/`)

### Core Pages

#### `indexHome.vue`
**Purpose**: Main landing page and content feed
**Functionality**:
- Displays popular posts and threads across all boards
- Shows trending boards and communities
- Provides navigation to various site sections
- Renders announcements and site updates
**Components Used**: `homeView.vue`, `popularPostsView.vue`, `boardDisplay.vue`, `announcements.vue`

#### `{...boardName}.vue` (Dynamic Route)
**Purpose**: Individual board/community pages
**Functionality**:
- Dynamically loads content based on board slug
- Displays board-specific posts, threads, and comments
- Shows board information and statistics
- Handles board-specific moderation tools (if user has permissions)
**Components Used**: `boardContentDisplay.vue`, `boardInfoDisplay.vue`, `boardPostDisplay.vue`, `boardThreadDisplay.vue`

#### `searchPage.vue`
**Purpose**: Global search interface
**Functionality**:
- Real-time search results as user types
- Filters by content type (users, posts, threads, comments, boards)
- Advanced search syntax support
- Search history and suggestions
**Components Used**: `SearchBar.vue`, search result components

#### `settingsPage.vue`
**Purpose**: User preference and account management
**Functionality**:
- Theme and appearance settings (stored in localStorage)
- Account settings (stored in database)
- Privacy and notification preferences
- Moderation preferences (if applicable)
**Data Storage**: Mixed localStorage and database persistence

### Content View Pages

#### `p/{...postViewPage}.vue` (Dynamic Route)
**Purpose**: Individual post view with comments
**Functionality**:
- Displays full post content and metadata
- Shows threaded comment discussions
- Real-time updates for new comments/interactions
- Interaction buttons (like, dislike, repost, comment)
**Real-time Features**: WebSocket updates for live interactions

#### `th/{...threadViewPage}.vue` (Dynamic Route)
**Purpose**: Thread view showing sequential posts
**Functionality**:
- Displays all posts in thread chronologically
- Shows comments for each post in thread
- Thread navigation and pagination
- Real-time updates for new posts/comments
**Real-time Features**: Live thread updates and notifications

### User & Tag Pages

#### `u/{...userProfile}.vue` (Dynamic Route)
**Purpose**: User profile pages
**Functionality**:
- Displays user information and statistics
- Shows user's posts, threads, comments
- User interaction history (likes, dislikes)
- Tag collection display
**Privacy**: Respects user privacy settings

#### `tags/tagStore.vue`
**Purpose**: Community tag marketplace
**Functionality**:
- Browse and discover user-created tags
- Apply tags to user profiles
- Tag popularity and trending
- Tag search and filtering

#### `tags/tagMaker.vue`
**Purpose**: Tag creation interface
**Functionality**:
- Visual tag editor with customization options
- Preview functionality
- Tag validation and guidelines
- Publishing and sharing options

## Layouts (`/app/layouts/`)

### `defaultLayout.vue`
**Usage**: Root page and unspecified pages
**Contains**: NavBar, SideBar, main content area, footer
**Responsive**: Mobile-first design with breakpoints

### `boardLayout.vue`
**Usage**: Board-specific pages
**Contains**: Board-specific navigation, board info sidebar, content area
**Features**: Board-specific theming support

### `postLayout.vue`
**Usage**: Individual post view pages
**Contains**: Minimal navigation, focused content area, comment section
**Features**: Optimized for content readability

### `threadLayout.vue`
**Usage**: Thread view pages
**Contains**: Thread navigation, sequential post display area
**Features**: Infinite scroll and pagination support

### `profileLayout.vue`
**Usage**: User profile pages
**Contains**: Profile header, tab navigation, content sections
**Features**: User-specific theming and customization

### `searchLayout.vue`
**Usage**: Search functionality
**Contains**: Prominent search bar, filter options, results area
**Features**: Real-time result updates

### `dashLayout.vue`
**Usage**: Administrative dashboards
**Contains**: Admin navigation, tool panels, main workspace
**Security**: Role-based access control

## Components (`/app/components/`)

### Universal Components (`universal/`)

#### `NavBar.vue`
**Purpose**: Primary site navigation
**Features**:
- Site branding and logo
- Integrated SearchBar component
- Top 6 most active boards (live updates)
- User authentication status
- Quick access buttons
**State Management**: Connected to board and user stores

#### `SideBar.vue`
**Purpose**: Secondary navigation and user actions
**Features**:
- User profile quick access
- Navigation to main sections (boards, search, settings)
- Notification indicators
- Theme switcher
**Responsive**: Collapsible on mobile devices

#### `SearchBar.vue`
**Purpose**: Unified search interface
**Features**:
- Real-time search suggestions
- Multi-type search (users, posts, threads, comments, boards)
- Search history
- Keyboard shortcuts support
**Performance**: Debounced queries and caching

### Advertisement Components (`advertisments/`)

#### `AdCard1.vue` through `AdCard4.vue`
**Purpose**: Non-invasive advertising display
**Templates**:
- **Poster**: Large format promotional content
- **Banner**: Horizontal strip advertisements
- **Sponsor**: Integrated sponsored content
- **Ad-network**: Third-party advertisement integration
**Ethics**: Clear labeling, user privacy respected

### Board Components (`boards/`)

#### Content Display Components

#### `boardContentDisplay.vue`
**Purpose**: Unified content feed for boards
**Features**:
- Mixed content types (posts, threads, comments, reposts)
- Pagination and infinite scroll
- Content filtering options
- Real-time updates
**Performance**: Virtualized scrolling for large feeds

#### `boardInfoDisplay.vue`
**Purpose**: Board information and statistics
**Data Displayed**:
- Board name, topic, description
- Historical statistics (total posts, threads, comments)
- Current activity (online users, recent activity)
- Board rules and guidelines
**Updates**: Real-time user count and activity metrics

#### `boardPostDisplay.vue`
**Purpose**: Post-specific content for boards
**Features**:
- Posts, post reposts, quote post reposts
- Post interaction buttons
- Thumbnail/preview generation
- Content warnings and filtering

#### `boardThreadDisplay.vue`
**Purpose**: Thread-specific content for boards
**Features**:
- Thread previews with post counts
- Thread reposts and quote thread reposts
- Thread status indicators (locked, pinned, etc.)
- Last activity timestamps

#### `boardCommentDisplay.vue`
**Purpose**: Comment display within board context
**Features**:
- Threaded comment display
- Comment on various content types
- Comment interaction tracking
- Moderation indicators

### Board Utility Components (`boards/utils/`)

#### `boardManager.vue`
**Purpose**: Board administration interface
**Capabilities**:
- Create new boards with custom settings
- Modify existing board properties
- Delete or archive boards
- Board permission management
**Access**: Admin and owner roles only

#### `databaseManager.vue`
**Purpose**: Database administration tools
**Features**:
- Database health monitoring
- Query performance metrics
- Data cleanup utilities
- Backup and restore functions
**Security**: Highly restricted access

#### `netLogger.vue`
**Purpose**: Network traffic monitoring
**Capabilities**:
- Real-time traffic analysis
- Request/response logging
- Performance metrics
- Security event detection
**Use Case**: Debugging and security monitoring

#### `postManager.vue`
**Purpose**: Content moderation interface
**Features**:
- Bulk content management (create, modify, delete)
- Content archival and quarantine
- Batch operations on posts/threads/comments
- Content restoration tools
**Efficiency**: Designed for high-volume moderation

#### `postReporter.vue`
**Purpose**: Content reporting system
**Features**:
- Report content for various violations
- Report categorization
- Evidence attachment
- Anonymous reporting options
**Integration**: Feeds into moderation queue

#### `priviledgeManager.vue`
**Purpose**: User permission management
**Capabilities**:
- Board-specific permission assignment
- Site-wide role management
- Permission inheritance
- Temporary privilege grants
**Granularity**: Fine-grained permission control

#### `reportedPostsDisplay.vue`
**Purpose**: Moderation queue interface
**Features**:
- Reported content review
- Report status tracking
- Moderation action logging
- Priority-based sorting
**Workflow**: Streamlined moderation process

#### `siteSearch.vue`
**Purpose**: Advanced administrative search
**Features**:
- Deep database searching
- Advanced syntax support
- Content indexing
- Search analytics
**Power User**: Enhanced search capabilities

### Moderation Components (`boards/utils/moderation/`)

#### User Action Cards
- `kickBoardCard.vue`: 48-hour board timeouts
- `tempBoardBanCard.vue`: 2 weeks to 6 months board bans
- `permaBoardBanCard.vue`: Permanent board bans
- `tempSiteBanCard.vue`: 6 months to 1 year site bans
- `permaSiteBanCard.vue`: Permanent site bans

#### Content Moderation
- `modReportPost.vue`: Post reporting interface
- `modDeletePost.vue`: Post deletion tools
- `postDeletion.vue`: Final deletion confirmation
- `postArchivalCard.vue`: Content archival system
- `modReportComment.vue`: Comment reporting
- `modDeleteComment.vue`: Comment deletion
- `modReportProfile.vue`: Profile reporting
- `profileLocker.vue`: Profile restriction tools

**Design Philosophy**: Progressive discipline system with clear escalation paths

### Dashboard Components (`dashboards/`)

#### `moderatorDash.vue`
**Purpose**: Basic moderation interface
**Tools**:
- Reported content queue
- User interaction history
- Basic moderation actions
- Community guidelines reference
**Scope**: Board-level moderation

#### `administratorDash.vue`
**Purpose**: Enhanced moderation interface
**Tools**: 
- All moderator tools plus:
- Advanced search and filtering
- User behavior analytics
- Site-wide trend analysis
- Automated detection tools
**Scope**: Site-wide moderation capabilities

#### `ownerDash.vue` & `coOwnerDash.vue`
**Purpose**: Full administrative control
**Tools**:
- All admin/moderator tools plus:
- System monitoring and maintenance
- Database administration
- Network diagnostics
- Site configuration management
**Authority**: Complete system control

### Home Page Components (`home/`)

#### `header.vue`
**Purpose**: SEO-optimized page header
**Features**:
- Meta tag management
- Schema markup
- Social media card generation
- Accessibility enhancements

#### `homeView.vue`
**Purpose**: Home page layout controller
**Functionality**:
- Coordinates home page components
- Manages content loading order
- Handles responsive layout
- User-specific content personalization

#### `announcements.vue`
**Purpose**: Site announcements display
**Features**:
- Rotating message carousel
- Announcement prioritization
- Dismissal tracking
- Emergency alert capabilities
**Tone**: Mix of important updates and community humor

#### `boardDisplay.vue`
**Purpose**: Board discovery interface
**Features**:
- Default and user-created boards
- Board popularity metrics
- Category filtering
- Board search functionality

#### `popularPostsView.vue`
**Purpose**: Trending content display
**Algorithm**: 
- Weighted scoring (likes, dislikes, comments)
- Time decay factor
- User engagement quality
- Abuse prevention measures

#### `footer.vue`
**Purpose**: Site footer with essential links
**Contents**:
- Copyright and legal information
- Site navigation shortcuts
- Important policy links
- Contact information

### Post Components (`posts/`)

#### `postCard.vue`
**Purpose**: Individual post display unit
**Contents**:
- Post content (text/media)
- User information and avatar
- Interaction buttons (like, dislike, repost, comment)
- Post metadata (timestamp, board, etc.)
**Design**: Responsive card layout with neobrutalist styling

#### `postView.vue`
**Purpose**: Detailed post view overlay/page
**Features**:
- Full post content display
- Comment thread integration
- Real-time interaction updates
- Social sharing options
**Navigation**: Modal-style or full-page depending on context

#### `threadCard.vue`
**Purpose**: Thread preview display
**Contents**:
- Sequential post previews
- Thread metadata (post count, participants)
- Thread status indicators
- Interaction summaries
**Complexity**: Handles multi-post content efficiently

#### `threadView.vue`
**Purpose**: Full thread display interface
**Features**:
- Sequential post loading
- Thread navigation tools
- Comment integration per post
- Thread-wide interaction tracking

#### `commentCard.vue`
**Purpose**: Individual comment display
**Features**:
- Threaded reply structure
- Comment interactions
- User identification
- Moderation status indicators
**Nesting**: Supports deep comment threads

### Profile Components (`profiles/`)

#### `profileView.vue`
**Purpose**: Profile page coordinator
**Functionality**:
- Assembles all profile components
- Manages tab navigation
- Handles permission-based content display
- Coordinates real-time updates

#### `userProfileDisplay.vue`
**Purpose**: User identification across content
**Display Elements**:
- Profile photo/gif with fallbacks
- Display name and username
- User status indicators
- Role and tag display
**Consistency**: Uniform user representation site-wide

#### `userInfoDisplay.vue`
**Purpose**: Comprehensive user information
**Information Types**:
- Basic info (name, username, bio)
- Optional demographic info (DoB, age)
- Visual customization (profile picture, banner)
- Social elements (roles, tags)
**Privacy**: Respects user privacy settings

#### `userPostsView.vue`
**Purpose**: User's post history display
**Content Types**:
- Original posts
- Post reposts
- Quote post reposts
**Features**: Chronological ordering, search/filter capabilities

#### `userThreadsView.vue`
**Purpose**: User's thread history display
**Content Types**:
- Original threads
- Thread reposts  
- Quote thread reposts
**Organization**: Thread status and engagement metrics

#### `userCommentView.vue`
**Purpose**: User's comment history
**Features**:
- Comment context linking
- Parent content references
- Comment interaction history
**Navigation**: Easy access to original discussions

#### `userLikesDisplay.vue`
**Purpose**: Content user has liked
**Content Types**: Posts, threads, comments, reposts
**Privacy**: May be limited based on user settings
**Purpose**: Showcases user preferences and interests

#### `userDislikesDisplay.vue`
**Purpose**: Content user has disliked
**Privacy**: Typically more restricted than likes
**Use Case**: Personal history, not public showcase

## Assets (`/app/assets/`)

### Styling (`css/`)

#### `style.css`
**Purpose**: Main stylesheet with neobrutalist design
**Features**:
- Core component styles
- Neobrutalist design elements (bold borders, strong shadows)
- Responsive design utilities
- Accessibility enhancements

### Theme System (`css/themes/`)

#### `nouveau.css`
**Aesthetic**: Typography-focused, borderless
**Font**: PT Serif
**Animation**: Subtle hover effects
**Inspiration**: Modern minimalist web design
**Target Audience**: Users preferring refined, readable interfaces

#### `bauhaus.css`
**Aesthetic**: Geometric minimalism
**Font**: PT Mono
**Design**: Creative geometry, expressive color use
**Philosophy**: Bauhaus design principles in web form
**Target Audience**: Users appreciating architectural design

#### `system8.css`
**Aesthetic**: Classic Mac OS 8 nostalgia
**Font**: Elegant serif fonts
**Features**: Retro iconography, classic color schemes
**Inspiration**: System 7 and Mac OS 8 interfaces
**Target Audience**: Retro computing enthusiasts

#### `milleniumxp.css`
**Aesthetic**: Early 2000s internet nostalgia
**Inspiration**: Windows ME/2000/XP era
**Reference Points**: 2005-2006 YouTube, early Twitter, 4chan
**Target Audience**: Users nostalgic for early web aesthetics

#### `serial.css`
**Aesthetic**: Late 90s web design
**Inspiration**: Serial Experiments Lain, Neocities aesthetic
**Features**: Experimental layouts, artistic expression
**Target Audience**: Users interested in alternative web design

#### `tui.css`
**Aesthetic**: Terminal User Interface simulation
**Color Scheme**: Black background with ANSI colors
**Font**: Bitmap/pixel fonts for authenticity
**Target Audience**: Developers and terminal enthusiasts

### Color Schemes (`css/colours/`)
**Purpose**: Switchable color palettes independent of themes
**Implementation**: CSS custom properties for easy switching
**Categories**: Light, dark, high contrast, and seasonal variants

---

# Backend Architecture

## Data Models

### User Model
**Core Fields**:
- `id` (Primary Key)
- `username` (Unique)
- `displayName`
- `email` (Unique, validated)
- `passwordHash` (Argon2id)
- `profileImage` (URL or file reference)
- `profileBanner` (URL or file reference)
- `bio` (Text, length limited)
- `birthday` (Optional, for age calculation)
- `createdAt`, `updatedAt`, `deletedAt`

**Relationships**:
- One-to-many: Posts, Threads, Comments
- Many-to-many: Roles, Tags, Board memberships

### Post Model
**Core Fields**:
- `id` (Primary Key)
- `content` (Text)
- `mediaAttachments` (JSON array of media references)
- `authorId` (Foreign Key to User)
- `boardId` (Foreign Key to Board)
- `threadId` (Optional, Foreign Key to Thread)
- `createdAt`, `updatedAt`, `deletedAt`

**Relationships**:
- Belongs to: User (author), Board
- Optionally belongs to: Thread
- One-to-many: Comments, Interactions

### Thread Model
**Purpose**: Groups related posts together
**Core Fields**:
- `id` (Primary Key)
- `title` (Optional)
- `authorId` (Foreign Key to original creator)
- `boardId` (Foreign Key to Board)
- `postIds` (Array of Post IDs in order)
- `isLocked` (Boolean for moderation)
- `isPinned` (Boolean for board prominence)
- `createdAt`, `updatedAt`, `deletedAt`

### Comment Model
**Core Fields**:
- `id` (Primary Key)
- `content` (Text)
- `authorId` (Foreign Key to User)
- `parentType` (Enum: Post, Thread, Comment)
- `parentId` (Polymorphic reference)
- `threadLevel` (Integer for nesting depth)
- `createdAt`, `updatedAt`, `deletedAt`

**Features**: Supports threaded discussions with configurable depth limits

### Board Model
**Purpose**: Communities/categories for content organization
**Core Fields**:
- `id` (Primary Key)
- `name` (Unique)
- `shorthand` (Unique abbreviation)
- `description` (Text)
- `isDefault` (Boolean for owner-configured defaults)
- `settings` (JSON for board-specific configurations)
- `createdAt`, `updatedAt`, `deletedAt`

**Relationships**:
- One-to-many: Posts, Threads
- Many-to-many: Users (memberships), Moderators

### Interaction Model
**Purpose**: Tracks likes, dislikes, and other engagement
**Core Fields**:
- `id` (Primary Key)
- `userId` (Foreign Key to User)
- `targetType` (Enum: Post, Thread, Comment)
- `targetId` (Polymorphic reference)
- `interactionType` (Enum: Like, Dislike, Repost, QuoteRepost)
- `isActive` (Boolean for toggling)
- `createdAt`, `updatedAt`

**Business Logic**: 
- Exclusive like/dislike states
- Prevents duplicate interactions
- Tracks interaction history

### Repost Models

#### Repost Model
**Purpose**: Shares existing content
**Core Fields**:
- `id` (Primary Key)
- `userId` (Foreign Key to reposter)
- `originalType` (Enum: Post, Thread, Comment)
- `originalId` (Reference to original content)
- `boardId` (Foreign Key to destination board)
- `createdAt`

#### QuoteRepost Model  
**Purpose**: Shares content with additional commentary
**Core Fields**:
- All Repost fields plus:
- `content` (Additional text from reposter)
- `mediaAttachments` (Optional additional media)

### Ink System (Reputation)
**Purpose**: User reputation system preventing gaming
**Implementation**:
- Calculated field, not stored directly
- Based on weighted interactions received
- Time-decay factor for relevance
- Anti-gaming measures:
  - Self-interaction prevention
  - Sockpuppet detection
  - Sudden spike analysis
  - Community voting validation

## State Management (`/app/stores/`)

### `auth.ts`
**Purpose**: Authentication and authorization
**State**:
- Current user information
- JWT token management
- Session persistence
- Login/logout states
**Methods**:
- `login()`, `logout()`, `register()`
- `refreshToken()`, `validateSession()`
- `updateProfile()`, `changePassword()`

### `posts.ts`
**Purpose**: Post management and caching
**State**:
- Post cache by ID
- User's post history
- Feed posts by board
- Post interaction states
**Methods**:
- `createPost()`, `editPost()`, `deletePost()`
- `fetchPosts()`, `fetchUserPosts()`
- `toggleLike()`, `addComment()`

### `threads.ts`
**Purpose**: Thread-specific functionality
**State**:
- Thread data and post sequences
- Thread interaction tracking
- User's thread history
**Methods**:
- `createThread()`, `addToThread()`
- `fetchThread()`, `reorderPosts()`
- `lockThread()`, `pinThread()`

### `interactions.ts`
**Purpose**: User engagement tracking
**State**:
- Like/dislike states by content
- Repost and quote repost data
- Interaction history
**Methods**:
- `toggleLike()`, `toggleDislike()`
- `createRepost()`, `createQuoteRepost()`
- `undoInteraction()`

### `board.ts`
**Purpose**: Board/community management
**State**:
- Board information and settings
- User's board memberships
- Board activity metrics
**Methods**:
- `fetchBoard()`, `joinBoard()`, `leaveBoard()`
- `createBoard()`, `updateBoard()`
- `fetchBoardContent()`

### `search.ts`
**Purpose**: Search functionality and caching
**State**:
- Search queries and results
- Search history
- Filter preferences
**Methods**:
- `search()`, `clearHistory()`
- `saveQuery()`, `getPopularSearches()`
- `advancedSearch()`

### `users.ts` & `profiles.ts`
**Purpose**: User data management
**State**:
- User profiles and information
- Follow relationships
- Profile customization settings
**Methods**:
- `fetchProfile()`, `updateProfile()`
- `followUser()`, `unfollowUser()`
- `blockUser()`, `reportUser()`

### `moderation.ts`
**Purpose**: Moderation tools and actions
**State**:
- Moderation queue
- User permissions
- Moderation history
**Methods**:
- `reportContent()`, `reviewReport()`
- `banUser()`, `deleteContent()`
- `escalateIssue()`

### `notifications.ts`
**Purpose**: User notification system
**State**:
- Notification queue
- Notification preferences
- Read/unread states
**Methods**:
- `fetchNotifications()`, `markAsRead()`
- `updatePreferences()`, `dismissNotification()`

### `privateMessages.ts`
**Purpose**: P2P encrypted messaging
**State**:
- Message threads
- Encryption keys
- Connection status
**Methods**:
- `sendMessage()`, `receiveMessage()`
- `establishConnection()`, `encryptMessage()`
**Technology**: Hypercore protocol integration

### `ui.ts`
**Purpose**: UI state and preferences
**State**:
- Theme selection
- Layout preferences
- Modal/dialog states
**Methods**:
- `setTheme()`, `toggleSidebar()`
- `showModal()`, `hideModal()`
- `savePreferences()`

## Type Definitions (`/app/types/`)

### `user.ts`
**Defines**:
- User interface with all fields
- UserRole enum (Owner, CoOwner, Admin, Moderator, User)
- UserPermissions interface
- AuthToken interface with JWT claims

### `post.ts`
**Defines**:
- Post interface with content and metadata
- MediaAttachment interface for files
- PostStatus enum (Published, Draft, Archived, Deleted)

### `thread.ts`
**Defines**:
- Thread interface with post relationships
- ThreadStatus enum (Open, Locked, Pinned, Archived)

### `comment.ts`
**Defines**:
- Comment interface with threading support
- CommentTree type for nested structures

### `interaction.ts`
**Defines**:
- Interaction interface for all engagement types
- InteractionType enum (Like, Dislike, Repost, QuoteRepost)

### `repost.ts` & `quoteRepost.ts`
**Defines**:
- Repost interfaces for content sharing
- Original content reference types

### `board.ts`
**Defines**:
- Board interface with settings
- BoardSettings interface for configurations
- BoardMembership interface for user-board relationships

### `profile.ts`
**Defines**:
- UserProfile interface for public information
- ProfileSettings interface for privacy controls

### `privateMessage.ts`
**Defines**:
- Message interface with encryption metadata
- MessageThread interface for conversations

### `user-shared.ts`
**Defines**:
- Tag interface for user customization
- Theme interface for appearance settings

### `api.ts`
**Defines**:
- APIResponse<T> generic for consistent responses
- ErrorResponse interface for error handling
- PaginationResponse<T> for paginated data

## Utility Files

### `apiClient.ts` (`/utils/`)
**Purpose**: Centralized API communication
**Features**:
- HTTP client configuration
- Request/response interceptors
- Error handling and retry logic
- Authentication token management
**Implementation**: Axios or Fetch-based with TypeScript support

### `routes.ts` (`/routes/`)
**Purpose**: Application routing definition
**Features**:
- Route definitions with parameters
- Route guards for authentication
- Dynamic route generation
- Navigation utilities

---

# Key System Features

## Ink Reputation System
**Purpose**: Prevent gaming while rewarding quality contributions
**Calculation Factors**:
- Interactions received (likes > reposts > quote reposts > comments)
- Content age (time decay)
- Interaction source credibility
- Community consensus validation

**Anti-Gaming Measures**:
- Self-interaction prevention
- Sockpuppet account detection
- Vote brigade detection
- Organic growth requirements

## Real-Time Features
**Implementation**: WebSocket connections for live updates
**Capabilities**:
- Live post/comment feeds
- Real-time search suggestions
- Notification delivery
- User presence indicators
- Collaborative thread building

## Moderation System
**Philosophy**: Progressive discipline with community involvement
**Levels**:
1. **Moderator**: Board-level content management
2. **Administrator**: Site-wide content oversight  
3. **Co-Owner**: System configuration access
4. **Owner**: Complete system control

**Tools**:
- Automated detection systems
- Community reporting mechanisms
- Appeal and review processes
- Transparent moderation logs

## Privacy & Security
**Features**:
- E2E encrypted private messages
- Granular privacy controls
- Anonymous posting options
- Secure authentication (Argon2id + JWT)
- Rate limiting and abuse prevention

## Theming System
**Architecture**: CSS custom properties + class switching
**Scope**: 
- Site-wide theme consistency
- Per-user theme preferences
- Board-specific theme options
- Accessibility compliance across themes

---

# Development Considerations

## Performance Optimization
- **Frontend**: Virtual scrolling, lazy loading, image optimization
- **Backend**: Database indexing, query optimization, caching layers
- **Real-time**: Efficient WebSocket management, message queuing

## Scalability Planning
- **Database**: Sharding strategies for user and content data
- **Caching**: Redis for session and frequent query caching
- **CDN**: Media file distribution and static asset caching
- **Load Balancing**: Horizontal scaling preparation

## Security Considerations
- **Input Validation**: Comprehensive sanitization and validation
- **Rate Limiting**: Per-user and per-IP restrictions
- **CSRF Protection**: Token-based request validation
- **XSS Prevention**: Content Security Policy and sanitization
- **SQL Injection**: Parameterized queries and ORM usage

## Accessibility
- **WCAG Compliance**: Level AA accessibility standards
- **Keyboard Navigation**: Full keyboard accessibility
- **Screen Reader Support**: Semantic HTML and ARIA labels  
- **Color Contrast**: Sufficient contrast ratios across themes
- **Alternative Text**: Comprehensive alt text for images

---

This enhanced design document provides a comprehensive overview of the Blueberry system architecture, making it clear for both human developers and AI systems to understand the project structure, component relationships, and implementation details.