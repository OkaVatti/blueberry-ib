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
- specified user roles and permissions (owner, co-owner, admin, moderator, user) 
- users can also have multiple custom roles or "tags" (example: wizard, goblin, commoner, evil, etc) which do not effect permissions, and are purely decorative display features 

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

____

## Component Explanations 

### Frontend

##### Posts

> Posts are defined as chucks of data
> 
>   (i.e: content/string, date, assigned user, 
>     optional content/media, with optional, 
>       mulitple comments, likes, and dislikes, 
>         created by multiple other users, 
>           including the original user who the 
>             post is linked to) 
>
> created by a user, or optionally also, by someone who is anonymous.

##### Comments

> Comments are defined as chucks of data
>
>    (i.e: content/string, date, assigned user, 
>      optional content/media, with optional, 
>        mulitple comments, likes, and dislikes, 
>          created by multiple other users, including 
>            the original user who created the comment 
>              attached to the post.)
>
> created by a user, or optionally also, by someone who is anonymous.

##### Threads

> Threads are defined as 
>
>    (i.e: an array of multiple posts, linked together 
>      by a "thread", stored as json, and referenced 
>        within sequential order by their identifier 
>          `orderInt`)
>
> created by one user, containing posts by 1 users or multiple users
> with the sole purpose of compiling thoughts and ideas linked to one
> another in a much easier way

##### Profiles

> Profiles are defined as
>
>    (i.e: a collection of data relavent to a given user. 
>      containing their display name, username/identifier, 
>        bio/introduction, optional age/birthday, optional 
>          country/region/timezone, optional profile image/gif, 
>            optional profile banner image/gif, their `ink` 
>              score, and the posts and comments they've made;
>                also optional including visibility of their 
>                  likes and dislikes)
>
> Displaying crutial information for indentifying a user that you may like, 
> dislike; enjoy, agree, or disagree with their posts or content.

##### Dashboards

> Dashboards are defined as 
>
>    (i.e: Containers around a given set of functional utilities, based 
>      around user permissions, and roles, along with a viewport, 
>        around which they are able to interact with the site, 
>          and their tools in such a way that allows the owner, 
>            co-owner, administrators, and moderators to do their 
>              jobs effectively.) 
>

##### Boards

> Boards are defined as
>
>    (i.e: A collection of posts and threads, with the option of 
>      comments being attached to posts, pertaining to a specific
>        topic. Boards are identified by a short hand abreviation 
>          consisting of either, a greek letter, english letter, 
>            or arabic numerals. Also containing information about
>              the board, announcements, and filters for posts and 
>                threads.)

### Backend

##### User

> Users are defined as
>
>    (i.e: source objects who create links between one another.
>      defined within the database by; display name, username, 
>        email, password hash, profile image/gif, profile banner 
>          image/gif, bio/introduction string, with the option of
>            being able to become "friends" with other users, 
>              "follow" and be "followed" by other users. can 
>                create posts, threads, comments, likes, dislikes, 
>                  reposts, text-reposts.)

##### Post
> Posts are defined as
>
>   (i.e: content/string, date, created by one user, 
>     optional content/media, with optional, 
>       mulitple comments, likes, and dislikes, 
>         created by multiple other users, 
>           including the original user who the 
>             post is linked to) 
>
> created by a user, or optionally also, by someone who is anonymous.

##### Comments

> Comments are defined as chucks of data
>
>    (i.e: content/string, date, assigned user, 
>      optional content/media, with optional, 
>        mulitple comments, likes, and dislikes, 
>          created by multiple other users, including 
>            the original user who created the comment 
>              attached to the post.)
>
> created by a user, or optionally also, by someone who is anonymous.

##### Threads

> Threads are defined as 
>
>    (i.e: an array of multiple posts, linked together 
>      by a "thread", stored as json, and referenced 
>        within sequential order by their identifier 
>          `orderInt`)
>
> created by one user, containing posts by 1 users or multiple users
> with the sole purpose of compiling thoughts and ideas linked to one
> another in a much easier way.

##### Boards

> Boards are defined as
>
>    (i.e: A collection of posts and threads, with the option of 
>      comments being attached to posts, pertaining to a specific
>        topic. Boards are identified by a short hand abreviation 
>          consisting of either, a greek letter, english letter, 
>            or arabic numerals. Also containing information about
>              the board, announcements, and filters for posts and 
>                threads.)

##### Permissions

> Users are given permissions based on their `role`. A `role`
> is special, specifically in the way in which it can handle permissions
>
> for example: the `Owner` `dev-role` would be assigned to the owner, 
> giving them full control over the more logistical and moderatorial workings 
> of the site. The `co-owner` `dev-role` would essentially do the same thing as 
> the `owner` `dev-role`. the `admin` `dev-role` gives special pemissions 
> pertaining to board creation, board moderation, and post moderation.
> the `moderator` `dev-role` is responsible for board specific moderation.
> These special `dev-role`'s are authorized via pgp key authorization, on top of 
> JWT verification.
>
> regular roles, such as  `n00b13`, `user`, `commoner`, `board-owner` `board-mod`, 
> `server-mod`, `board-admin`, `server-admin`, `co-owner`, and `owner`
> are all regular roles that function according to their permissions.
> 
> For example, the `n00b13` `role`, or aka "noobie", is a role for newcomers,
> its function access and permissions are limited, and users with this role
> are only allowed to make posts in admin/co-owner/owner-approved boards, or 
> otherwise, boards specified in the `magichan.conf` file under `## Boards`.
> as compared to the `user` or `commoner` `role`, which can post freely.
> the `board-ownder`, `board-mod`, `server-admin`, `co-owner`, and `owner`
> `role`s also enjoy this luxery, with the added bonus of having specific 
> permissions raging across forms of moderation, and moderation, tools
> given to each role. 
>
> Tags are different than `role`s in that they don't allow for permissions, 
> and are purely decorative. Primarily used for user expression.
>
> For example, a `user` could have the `Swag` tag, the `Cool` tag, the 
> `cute` tag, the `Gamer` tag, the `LGBTQ+` tag, the `trans` tag. anything the user inputs, 
> it's it's written to that spot in the database.

##### Authorization, Passwords, Communication, and Encryption 

> User authorization is handled via JWT tokens, and agron2id is used for 
> password hashing and ververification. Specific roles, such in the case 
> for `admin`, `co-owner`, and `owner` roles, auth-tokens and encryption 
> keys are required for access to the site and its extended functionality
>
> Database is to be backed-up, and both the live database and backup are 
> encrypted via AES-256, this ensures that all text and images stored
> are done so securely.
> 
> To add another extra layer of security, we layout our trafic through
> up-to-date TLS/SSL encryption algorithms, such that the data tranfered,
> via a modern browser, is done so securely.

##### Moderation

> Moderation is conducted by `moderator`, `admin`, `co-owner`, and `owner`
> `role`s. moderation consist of 4 consequenses, all at the discression of
> the moderation.
>
> 1. Kick
> 1a. A kick consists of a simple "Timeout" from a specific board. If a user
>      is misbehaving, moderation steps in and handles it by giving them a
>       kick, or "timeout" from the board in question.
>
> 2. Temporary Board-Ban
> 2a. A temporary board-ban consists of a temportary ban from a given board.
>      i.e: someone has been misbehaving multiple times, and needs to take a
>      break from a certain board.
>
> 3. Perminent Board-Ban
> 3a. A perminent board-ban consists of a perminent ban from a give board.
>      i.e: someone violated multiple rules on the board, has been deemed
>      "an issue", and should no longer interact with the board.
>
> 4. Termporary Site-Ban
> 4a. A temporary site-ban consists of a temporary ban from the site at large.
>      i.e: someone has been an issue on multiple boards, and is kindly being 
>      asked to take a break from the site, and to reflect on how their actions
>      made other people on the site too, and to hopefully come back, changed
>      for the better.
>
> 5. Perminent Site-Ban
> 5a. A perminent site-ban consists of a perminent ban from the site at large.
> i.e someone has proven to no intention of change, who has violated multiple 
> site and/or board rules on a basis frequent enough for them to be deemed
> "a problem" by more than 2 users.
>
> Moderation and users are given tools to help developers, moderators, and users
> create a community that represents us for who we are as people.
> Users can report posts, provide a reason, and submit it for manual review.
> Moderators can report posts, provide reason, and submit for manuals review.
> Moderators can also review reported posts, provide a reason, and have the post 
> deleted using the tools at their disposal.

> Administrators can report and moderate posts the same way moderators can, on 
> top of all that, they can view all posts made by users, from from their 
> dashboard. They also have the ability to kick, temp-board-ban, perma-board-ban

> users at their discression, given a good reason. 

> Owner and Co-Owner have the ability to do the same as all other roles, but, are
> also able temp-site-ban and perma-site-ban users.

##### Site Rules

> Imageboard / Microblogging Website rules are as follows
> 
> 1. Engage in tomfoolery
> 2. Be whimsical with your presence, always
> 3. don't be a mark-ass trick-ass BITCH
> 4. don't be mean unless it's for a definitely valid reason
> 5. don't be a stranger sweet child :3
> 6. no slurs unless you're allowed to say them, never towards anyone for derogatory purposes or reasons
> 7. no content pertaining to the exploitaion of minors according to any of the definitions givene by the 
>     countries in where the existence and distribution of the exploitative material is prohibated.
> 8. anonymity comes with a price, please use it wisely, do not violate any of the rules listed. 
>     you will be kicked, termorarily or perminently banned from boards or the site, and your name forever
>     listed on the `banned board`, with your user name, reason for ban, the post, or comment that got you 
>     banned, and your last 5 likes, comments, and posts you had mad made at the time, all for users to 
>     comment upo
> 9. Free speech is something we blieve in, but it is not absolute; the things we say have consequenses.
>     It is for that reason alone that we ask our users to engage in their freedom responsibly, and to
>     be sure that we use our words responsibly, as they effect those around us more than we would like
>     to think
 