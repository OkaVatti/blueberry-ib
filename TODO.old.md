# TODO list :3

ok, as a forewarning, this sits as both a TODO list, a design document, 
and as an explanation for what everything is, and why there are 
so many components with very similar names (blame english for being a 
horrendious language invented by the brain eating aomeba sent from outer
space destined to vanquish humanity as easy as it did the dinosaurs...)

to also better explain this before diving into further detail. here is 
my frontend and backend tech stacks


- Frontend Tech Stack
    - pnpm
    - Vue 3 + Nuxt 4 + TypeScript
    - TailwindCSS + Vanilla CSS
    - Pinia
    - Vite

- Backend Tech Stack
    - Go
    - Gorm
    - Echo Framework
    - Agron2id
    - JWT
    - Postgres
    - Docker

so let's start off simply; these were the things i wanted to get done 
when initially creating this document

## Frontend

- implement universal components
- implement posts, threads, and their layouts, along with their respective views, pages, and components 
- implement users, profiles, pages, and their components
- implement advertisment components
- implement home page, its layout, and its components
- implement boards, their pages, layouts, and components
- implement respective owner/co-owner/admin/moderator dashboards and their utility components for website moderation purposes
- implement search page and its respective components
- implement custom user defined roles/tags (example: wizard, goblin, mage, commoner, goon, evil, etc.)   

## Backend

- organize handlers into proper, seperate files
- specify and handle for user roles and persissions
- work on better postgres database integration
- implement owner and co-owner configurable default boards
- implement low network level logging of traffic, posts made
- implement annonymous posts for logged-in users and for logged out users 
- implement `.conf` file based configuration of database and server functionality 
- implement allow for backend modifications to be made to the frontend
- implement `ink` (similar to reddit karma) in such a way that prevents abuse or "game-ifying" of the algorithm
- implement low level api for administrative and moderation purposes
- implement `client -> server -> database` api for handling custom user defined tags  

____

this quickly devolved into several debugging sessions with AI chat bots poorly trained on outdated code 
and whimsical mysticism about their oh-so superb programming prowess.

"rewrite the following frontend and backend such that they are 100% compatible. rewrite as much code as you need to 
in order to make it compatible, and for the app to handle posts correctly, efficiently, and securely. 
modify the backend files if needed, create new files if deemed necessary, and make sure posts are defined as such within code"

several models obviously failed due to limitations in regards to the financial site of my endeavors. 
so instead i decided it would be best to just lay out every component name and specify what that 
component does, and how it is used in other parts of the project.

here is a "brief" (as brief as i can make it, while still maintaining accuracy) explanation on what everything is

# Blueberry

An Experimental, Neobrutalist, Combination Imageboard + Microblogging PWA

Built using `pnpm`, `Vue 3`, `Nuxt 4`, `TailwindCSS` + `Vanilla CSS`, `Pinia`, `TypeScript`, `Go`, `gorm`, `echo framework`, `postgres`, and `docker`

features:
- users
- posts
- threads
- comments
- reposts & quote reposts
- boards / communities
- ink (read further, it's kinda like karma on reddit)
- accurate real-time search functionality (*insert the emoji for PAIN here*)
- private messages (E2E Encryption, P2P Functionality, built using the hypercore protocol)
- extensive moderation tools
- site-wide role-based user moderation and user permissions
- user created decorative "tags"
- simple, non-invasive advertisments

## Frontend Stuff

- Default Theme: Soft, Pastel, Neobrutalist, light blues and blurples (default, multiple colorschemes supported)

- Pages:
    - indexHome.vue
        - main / home page
    - {...boardName}.vue
        - slug file, loads in content dynamically. you click on a board? this is where that takes you
    - searchPage.vue
        - default user search page. accurate and fast. displays results and suggestions in real time as you're typing
    - settingsPage.vue
        - per user settings page. most settings are saved to localStorage, others are saved directly to the database and are applied once auth requirements have been met.
    - p/{...postViewPage}.vue
        - when clicking on a post, the user is taken to this page.
        - real time updates
        - post content loaded and rendered dynamically
        - comments loaded and rendered with post
    - th/{...threadViewPage}.vue
        - when a clicking on a thread, the user is taken here.
        - real time updates
        - thread / post content is loaded and rendered dynamically in sequential order
        - comments are loaded and rendered with thread posts
    - tags/tagStore.vue
        - main place for users to share their tags with other users. tags be be applied to a users profile from the tagStore
    - tags/tagMaker.vue
        - main place for users to make their tags.
    - u/{...userProfile}.vue
        - slug file, loads in user profile content and components dynamically.

- Layouts:
    - boardLayout.vue
        - default layout for boards 
    - dashLayout.vue
        - default layout for `pages/admin/{...dashBoard}.vue`
    - defaultLayout.vue
        - default layout for root page and pages where the layout is otherwise not specified
    - postLayout.vue
        - default layout for `pages/p/{...postViewPage}.vue`
    - profileLayout.vue
        - default layout for `pages/u/{...userProfile}.vue` 
    - searchLayout.vue
        - default layout for `pages/searchPage.vue`
    - threadLayout.vue
        - default layout for `pages/th/{...threadViewPage}.vue`

- /app/assets/
    - css/
        - style.css
            - main css file. contains neobrutalist styles for everything
        - colours/
            - colours folder
            - contains colorschemes with class names for switching between different styles
        - themes/
            - themes folder
            - contains themes with 1:1 class names for switching between different styles
            - client side only :3
            - nouveau.css
                - focused more on typography
                - borderless look and feel
                - subtle hover animations
                - PT Serif
                - similar to / reference
                    - https://daliwali.neocities.org/
                    - https://bengarbow.com/about
            - bauhaus.css
                - focused on creative uses of geometry
                - looks and feels minimal
                - expressive use of color
                - hover animations
                - PT Mono
                - reference bauhaus design philosophy within the context of web-design
            - system8.css
                - classic
                - focused on iconography
                - elegant serif fonts
                - pretty colors
                - inspired heavily by System7.css, based on the look and feel of Mac OS 8.
                - reference the look and feel Mac OS 8 and its predecesor System 7
            - milleniumxp.css
                - mimics the look and feel of the internet during the early-to-mid 2000s
                - reference Windows ME, Windows 2000, Windows XP, the early years of Windows Vista
                - reference what youtube, 4chan, twitter, looked like circa 2005-2006
            - serial.css
                - late-90s web design / neocities look and feel.
                - inspired by 
                    - https://fauux.neocities.org
                    - https://blackwings.neocities.org
                    - Serial Experiments Lain internet aesthetics (reference previous links)

            - tui
                - based around the designs of terminal user interfaces
                - B/W with ANSI colors on black background
                - pixel art/bitmap font for aesthetic accuracy
    - images/
        - not really needed but still there for when we add branding and country flags

- /routes/
    - routes.ts
        - main file for defining page routes.
        - responsible for all site navigation by users

- /app/stores/
    - auth.ts
        - auth store
    - board.ts
        - board store
    - interactions.ts
        - store for interactions
            - likes & dislikes
            - repost + quote repost (these are technically interactions)
            - comment
            - share
    - moderation.ts
        - store for moderation specific interactions
    - notifications.ts
        - store for notifications
    - posts.ts
        - store for posts
    - privateMessages.ts
        - store for private messages
    - profiles.ts
        - store for user profiles
    - search.ts
        - store for search components and queries
    - threads.ts
        - store for threads
    - ui.ts
        - store for ui specific functionality based on backend functionality?
    - users.ts
        - store for users

- /app/types/
    - api.ts
        - defines errors and responses for the API
    - board.ts
        - defines what boards are as objects within the database
    - comment.ts
        - defines what a comment is (slightly different than a regular interaction (like or dislike))
    - interaction.ts
        - defines what likes, dislikes
    - post.ts
        - defines what a post is
    - privateMessage.ts
        - defines what private messages are
    - profile.ts
        - defines what a profile is
    - repost.ts
        - defines what a repost is (for reposted posts and threads)
    - quoteRepost.ts
        - defines what a quote repost is (for reposted posts and threads)
    - thread.ts
        - defines what a thread is
    - user-shared.ts
        - mostly defines stuff relating to user customization
    - user.ts
        - defines what a user is
        - defines user roles
        - defines user permissions
        - defines user auth tokens

- /public/
    - _robots.txt
        - no bad robots allowed
    - favicon.icon
        - standard website favicon for branding

- /ssl/
    - blueberry.lan.crt
        - websites self-signed (reference `lilith.key`) TLS/SSL certificate
    - lilith.key
        - server owner's ssl key

- /utils/
    - apiClient.ts
        - api handler between client and server

- /app/components/
    - universal/
        - NavBar.vue
            - Navigation bar, contains site branding, SearchBar, and the top 6 boards on the site in a list with links to each board 
                - (updated live based on the number of active users on each board). 
        - SideBar.vue
            - The sidebar contains buttons for navigating the rest of the site 
                - user profile page
                - boards page
                - search page
                - settings page
        - SearchBar.vue
            - A search bar that provides the user with results in real time as they type, used to search for 
                - users
                - posts
                - threads
                - comments
                - reposts
                - quote reposts
                - boards
    - advertisments/
        - AdCard (1 through 4)
            - Advertisment Cards used for advertising goods and services to users of the site 
            - poster template
            - banner template
            - sponsor template
            - ad-network template
    - boards/
        - boardContentDisplay.vue
            - displays all posts, threads, comments, reposts, and quote reposts made to a given board
        - boardInfoDisplay.vue
            - displays information about the board 
                - name
                - topic
                - description
                - number of posts ever made to the board
                - number of threads ever made to the board
                - number of comments ever made to the board
                - number on online users who are currently viewing the board
        - boardPostDisplay.vue
            - displays all posts, post reposts, and quote post reposts made to a given board
        - boardThreadDisplay.vue
            - displays all threads, thread reposts, and quote thread reposts made to a given board 
        - boardCommentDisplay.vue
            - displays all comments made to a given post, thread, repost, quote repost or comment under a given board.
        - utils/
            - boardManager.vue
                - manage site boards (create, modify/edit, delete)
            - databaseManager.vue
                - collection of tools to make database management easier
            - netLogger.vue
                - log and display site network traffic 
            - postManager.vue
                - create, modify, delete, archive, quarantine posts, threads, reposts, quote reposts, and comments
            - postReporter.vue
                - report posts, threads, reposts, quote reposts, and comments
            - priviledgeManager.vue
                - manage user priviledges across boards and site-wide
            - reportedPostsDisplay.vue
                - for moderation use only, display posts that have been flagged/reported for manual review
            - siteSearch.vue
                - search the entire site and database in real time. this time with cool syntax.
            - moderation/
                - kickBoardCard.vue
                    - moderation card for temporarily kicking/timing out a user from a given board for 48 hours
                - modReportPost.vue
                    - moderation component/tool for reporting singular or posts 
                - modDeletePost.vue
                    - moderation component/tool for deleting singular or multiple posts
                - postDeletion.vue
                    - final stage for post deletion on the frontend.
                - postArchivalCard.vue
                    - moderation card for archiving boards, posts, threads, reposts, quote reposts, and comments
                - modReportComment.vue
                    - moderation tool for reporting comments
                - modDeleteComment.vue
                    - moderation tool for deleting comments
                - modReportProfile.vue
                    - moderation tool for reporting a single or multiple user profiles
                - profileLocker.vue
                    - moderation tool for locking down a users profile
                - tempBoardBanCard.vue
                    - moderation card for temporarily banning a user from a given board for a set amount of time 
                        - (min: 2 weeks, max: 6 months)
                - permaBoardBanCard.vue
                    - moderation card for permanently banning a user from a given board for an indefinite amount of time
                - tempSiteBanCard.vue
                    - moderation card for temporarily banning a user from the site at large
                        - (min: 6 Months, Max: 1 Year)
                - permaSiteBanCard.vue
                    - moderation card for permanently banning a user from the site at large for an indefinite amount of time
    - dashboards/
        - moderatorDash.vue
            - Dashboard for moderators, contains tools necessary for reporting and handling reported posts efficiently and effectively 
        - administratorDash.vue
            - Dashboard for administrators, contains moderator tools along with other tools for navigating and search the site for content violating the ToS much easier 
        - coOwnerDash.vue
            - reference ownerDash
        - ownerDash.vue
            - same suite of tools as moderators and admins, with the added bonus of networking and database specific tools.
    - home/
        - header.vue
            - simple page header, specifically for dealing with SEO.
        - homeView.vue
            - more like a proto-layout, components both universal relating to the home page go here
        - announcements.vue
            - simple annoucement component that displays a rotating assortment of "important messages" (mostly used for jokes)
        - boardDisplay.vue
            - displays all default and user created boards
        - popularPostsView.vue
            - displays the current most popular posts on the site (based on likes, dislikes, and comments)
        - footer.vue
            - copyright information, site navigation, important links
    - posts/
        - postCard.vue
            - component card containing post data (text/media, user, interactions)
        - postView.vue
            - clicking on a post takes you to a "View" on top of the current page that displays the content of the given post card, along with the comments linked to that post
        - threadCard.vue
            - component card containing data from multiple posts, displayed in sequential order 
                - (ThreadID, content/text/media, user, interactions)
        - threadView.vue
            - clicking on a post takes you to a "View" on top of the current page that displays the content of the given thread card, along with the comments linked to that thread
        - commentCard.vue
            - component card containing the content of a given comment (text/media, user, interactions)
    - profiles/
        - profileView.vue
            - main file where all profile related components are assembled
        - userProfileDisplay.vue
            - fetches and displays a given users 
                - profile photo/gif
                - display name
                - and username 
            on their posts, threads, reposts, quote reposts likes, comments
        - userInfoDisplay.vue
            - displays a given users information
                - name
                - username
                - DoB (optional)
                - b-day (optional, extrapolated from DoB)
                - bio/intro string
                - profile picture/gif
                - profile banner/gif
                - role (2 per user)
                - tags (max of 4 per user)
        - userPostsView.vue
            - fetches and displays all posts, reposts, and quote reposts (for posts) for a given user
        - userThreadsView.vue
            - fetches and displays all threads, reposts, and quote reposts (for threads) for a given user 
        - userCommentView.vue
            - fetches and displays all comments made by a given user
        - userLikesDisplay.vue
            - fetches and displays all objects that a given user has "liked" 
                - (i.e: like; boolean=true, dislike; boolean=false)
        - userDislikesDisplay.vue
            - fetches and displays all objects that a given user has "disliked" 
                - (i.e: like; boolean=false, dislike; boolean=true)

- /nuxt.config.ts

```ts
import tailwindcss from "@tailwindcss/vite";

// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: "2025-07-15",
  devtools: { enabled: true },
  builder: "vite",
  buildDir: "blueberry-ib-build",
  dev: true,

  devServer: {
    host: "0.0.0.0",
    port: 9876,
    https: {
      key: "ssl/lilith.key",
      cert: "ssl/blueberry.lan.crt",
    },
  },

  dir: {
    assets: "app/assets",
    layouts: "app/layouts",
    pages: "app/pages",
    public: "public",
  },

  modules: [
    "@nuxtjs/sitemap",
    "@nuxtjs/robots",
    "@nuxt/eslint",
    "@nuxt/image",
    "@nuxt/fonts",
    "@nuxtjs/color-mode",
    "nuxt-feather-icons",
    "@pinia/nuxt",
  ],

  imports: {
    dirs: ['stores']
  },

  // Vite-specific server config (covers vite dev server)
  vite: {
    plugins: [tailwindcss()],
  },

  typescript: {
    strict: true,
  },
});
```

## Backend Stuff

- Users are objects with user-defined values attached to them 
    - display name
    - username
    - email
    - password (hashed, stored, and verified with argon2id)
    - birthday (optional)
    - age (calculated from birthday, optional)
    - profile image/gif
    - profile banner/gif
    - bio/intro string (i.e: "hi my name is lily, i'm 20, and i hate my job :3")
- Posts are objects stored within a database, created by a singular user, containing content 
    - strings
    - media (all popular and widely supported MIME-types)
    - with the ability for all other users, including the original post creator, or 'poster', to create comments, 
    - toggle the boolean values of their likes or dislikes according to how the system allows for it. 
- multiple posts created by the same user can be linked together, creating a thread, another object stored within a database that consists of an array of post id's, which gets rendered and displayed their sequential order.
- comments can be made to posts and threads, comments are objects, similar to posts in their function 
    - (text or media to be shared with others), 
    - in relation to a single given post, thread, or comment created by any user.
- Reposts are objects stored within a database, created by a user, containing the content of another users 
    - post
    - comment
    - or thread. 
    - the original post (the one being reposted) is directly embeded, linked to, and referenced.
    - The repost appears under the users profile (the user who reposted the original content)
- Quote Reposts are objects stored within a database. As a combination between posts and reposts, similar to posts in their function 
    - (text or media to be shared with others)
    - in relation to a single given 
        - post
        - thread
        - or comment created by any user 
        - with the added benefit of having more than just the original post linked, referenced, and embeded. 
        - allow users to attach their own text/media to the Reposted post, thread, or comment.
- likes and dislikes are states that can be toggled between. 
    - these interactions are exclusive to posts, threads, and comments.
    - likes are how we indicate we enjoy a post, thread, or comment; a dislike is how we show the opposite. 
    - When the `like` boolean is toggled to true, dislike is toggled to false, and vice-versa. 
    - These interactions are reversible.
- boards are objects representing different communities, topics, or interests. Boards have values and links assigned to them 
    - boardName: string
    - shorthand: string
    - boardDescription: string
    - Boards can have posts and threads made to them, and those posts and threads and all receive 
        - likes
        - dislikes
        - comments from any user with the required permissions
- users, posts, threads, comments, boards, tags, and profiles all have the following table entries 
    - `createdAt`
    - `updatedAt`
    - `deletedAt`, 
    these table entries provide the system with a way to keep track of the time content was created, last edited, and deleted.
- ink, like reddit karma, is a way to easily verify a user by how active they are on the site, based on how many 
    - likes
    - dislikes
    - reposts
    - or quote reposts that their original post, comment, repost, or quote repost received.