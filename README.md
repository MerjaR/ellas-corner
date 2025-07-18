
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


## Features Summary

- User Registration & Login (cookie sessions)
- Submit, like, and comment on items (only when logged in)
- Browse all items publicly
- Image upload support for items and user profile
- Categories: by baby age/stage
- Filtering: by popularity and age group
- Profiles: editable with liked/submitted items and optional location
- Secure password handling via `bcrypt`
- Clean templating with Go’s `html/template`
- Containerised with Docker

---

## Tech Stack

- **Language**: Go (net/http, html/template, database/sql)
- **Database**: SQLite (with `github.com/mattn/go-sqlite3`)
- **Authentication**: bcrypt, cookies
- **Frontend**: Custom CSS 
- **Containerization**: Docker

## Tests
Basic unit and integration tests have been added to the project to provide examples of how to test core functionality in a Go web application. These tests use Go’s standard testing package and an in-memory SQLite database to avoid modifying production data.

### Included tests
Unit tests for:
- Post creation logic
- Reaction handling (like/dislike)
- Password hashing with bcrypt

Integration tests for:

- Full user registration and login flow
- Creating a post while authenticated (including session token handling)


Notes
These tests are designed to be minimal but extendable. They demonstrate:

- How to safely isolate logic for testing

- How to test handler behaviour without a full server

- How to simulate form submissions and cookies in test requests

Future tests could expand to cover error cases, comments, and more detailed handler flows.

To run all tests:

go test ./... -v


## Some considerations for future development

This project currently uses a traditional, form-friendly routing style (e.g., /create-post, /edit-post) optimised for server-rendered HTML and Go's net/http standard library.

The codebase is structured in a way that could support RESTful API endpoints which could make it easier to work with.  

Hardcoded paths (e.g., DB location, migration file) could be replaced with environment variables for scaling. 

Images could have UUID based filenames to avoid the same names being used for uploaded images (applied only partly right now).

Hard coded categories could be managed through an admin account, which could also delete or update user accounts. 

In general, there are many features that would add usability like updating your e-mail address. They are entirely possible to add with similar logic already added to the project, but the project has to have some boundaries to free time to learn new things.



---

## Acknowledgments

Originally built for the [kood/Sisu](https://koodsisu.fi) Full Stack Developer Program.  
Transformed with love into *Ella’s Corner* — a portfolio project focused on usability and real-world user needs.
