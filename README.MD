i'm interested in creating a self-hosted imageboard using pnpm, vue, vite, go, gorm, echo-framework, postgresql.

i want it to look and feel like a new take on classic imageboards.

# Blueberry

a new take on imageboards

Feels like a combination between 4chan, Lainchan, Reddit, and Bluesky

- Admin, Moderation, and Regular Users with verifiable tokens, uuid's, usernames, passwords (hashed with argon2)
- Multiple themes
- Hashing using argon2 for validation and verification of sensitive data
- Encryption using a hard to crack AES-256 variant
- Multiple Default Boards (a board is a group of and for posts relating to a given topic)
- User Created Boards (a board is a group of and for posts relating to a given topic)
- Posts (Created by users, Support for text, images, videos, embeds, links, extended markdown syntax)
- Comments (Created by users, Support for text images, videos, embeds, links, extended markdown syntax)
- Threads (A collection of many posts linked together by a "thread", created by 1 user)
- Likes (toggle on all for posts, comments, and threads, by any and all users)
- Dislikes (toggle on all for posts, comments, and threads, by any and all users)
- Ink (similar to reddit karma, user score based on how many likes, comments, and views they've gotten)
- options for 100% anonymity / privacy options
- Admin / Moderator users
- Moderation / Admin dashboard
- real-time updates 

## Directory Structure

/home/lilith/code-shit/blueberry/
| - apps
| - node_modules
| - packages
| - package.json
| - pnpm-lock.yaml
| - pnpm-workspace.yaml
| - README.MD

./apps/
| - backend
| - frontend

./apps/frontend/
| - node_modules
| - public
| - src
| - index.html
| - package.json
| - README.md
| - tsconfig.app.json
| - tsconfig.json
| - tsconfig.node.json
| - vite.config.ts

apps/backend/
| - go.mod
| - go.sum
| - main.go