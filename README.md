
# Ella’s Corner – A Community-Powered BabyKit Resource

Originally developed as a project for the Go module of the  
[kood/Sisu Full Stack Developer Training](https://koodsisu.fi), this application began as the **Literary Lions Forum** — a simple book club discussion board. It has since evolved into **Ella’s Corner**, a warm, thoughtful platform for new parents to share and discover the most useful baby products, complete with community feedback and donation options.

---

## Project Origins: Literary Lions Forum

This application was built under the following constraints as part of the training requirements:

- ✅ Written **entirely in Go**, using only the **Go standard library** except for the following
- ✅ Uses **SQLite** for data storage
- ✅ Passwords secured using **bcrypt**
- ✅ Session management and **cookie-based login**
- ✅ Public can view all posts and comments
- ✅ Only **logged-in users** can post, comment, or react
- ✅ Includes **categories, filtering, and search**
- ✅ No JavaScript or frontend frameworks (e.g., no Tailwind) unless it is a bonus feature
- ✅ Must be packaged and runnable via **Docker**

---

## Transformation into *Ella’s Corner*

As part of a portfolio refinement effort, the project was reimagined to better showcase real-world product thinking, user-centered design, and practical UX considerations. Changes include:

### Feature Reframing
- Forum **posts** are now **baby item recommendations**
- Each item includes:
  - Image upload
  - Age/stage tags (e.g. Newborn, 3–6 months)
  - Donation flag (users can offer to donate)
- Homepage redesigned to show **top items by age group**
- Submit page for users to contribute items with context
- **Filters** now include:
  - By age group
  - By date

### User Experience Enhancements
- User profiles updated to show:
  - Liked items
  - Submitted items
  - Optional country (used only for donation matching)
- Default **avatars** for users without a profile picture

Users now have the ability to delete or edit their own posts, including changing the image. They are not able to edit or delete anyone else's posts. They can also remove the up for donation tag when the item has been donated. 

Two JavaScript features were added, which could have been originally accepted as bonus features. One is to stop the page from resetting to the top of the page when a user for example likes a post. The second is allowing to use the dropdown menu under the profile picture more smoothly. 

## How to Run the Application

After cloning the repository, run with go run main.go .
Alternatively, run in Docker with docker-compose up .

The website is accessible at localhost:8080 .


## 🔐 Features Summary

- User Registration & Login (cookie sessions)
- Submit, like, and comment on items (only when logged in)
- Browse all items publicly
- Image upload support for items and user profile
- Categories: by baby age/stage
- Filtering: by popularity, age group, and location (country)
- Profiles: editable with liked/submitted items and optional location
- Secure password handling via `bcrypt`
- Clean templating with Go’s `html/template`
- Containerised with Docker

---

## 💻 Tech Stack

- **Language**: Go (net/http, html/template, database/sql)
- **Database**: SQLite (with `github.com/mattn/go-sqlite3`)
- **Authentication**: bcrypt, cookies
- **Frontend**: Custom CSS 
- **Containerization**: Docker

---

## Acknowledgments

Originally built for the [kood/Sisu](https://koodsisu.fi) Full Stack Developer Program.  
Transformed with love into *Ella’s Corner* — a portfolio project focused on usability and real-world user needs.
