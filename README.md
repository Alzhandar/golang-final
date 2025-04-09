# Restaurant Management System

## Project Overview
This is my Go project for the programming language course. It's a restaurant management system backend that handles restaurants, tables, menus, events, and users.

The system lets restaurant managers add their restaurants to the platform, create sections inside each restaurant, add tables with QR codes, and manage special events. Users can view restaurant info and book tables for events.

## ERD Diagram
![ERD Diagram](./docs/erd_diagram.png)

## <img src="https://media-hosting.imagekit.io/3b414a17ce154bb2/Fire.png?Expires=1838019735&Key-Pair-Id=K2ZIVPTIP2VGHC&Signature=XLfIq-eOBgMwCt6DqPEQmUef7~J0PZvWd1glU9X4VuWE-3GwK5i0b8M6Ig7Pj9rm-gkYRR3RUOtlT5~f03HrT96gUAAX7IXuAUUjmKV0uaCouMSA61vLGTyeLdUMfyX4BIlWp5Q7sqmeRrGV9Ac9DfIy0AxSYsSYFQgadfSFG-FsOfcvOV6SYyel-Hny-2YXp8Ut7yLS~GF6~orc05XPvfdSXXLvoy5Np5TlBIf9vWw7v4t6mPKHd3EvZo~gA1a5Vtn297uks3-YI9N2-Z7tYnlooB0u~x1r2oBvJWzdVB4tPDkDUA-X5IaiqLUVtEn~nHH8rxGciHTu7EP8~svdfQ__" alt="Fire" width="30" height="30" />
 Presentation
You can view my presentation [here](./docs/GO.pdf).

## Tech Stack
- Golang with Echo framework for REST API
- PostgreSQL for database
- Swagger for API documentation
- Docker & Docker Compose for containerization
- JWT for authentication

## Main Features
- Restaurant management (add, edit, delete restaurants)
- Table management with QR code generation
- Menu management with multi-language support (Russian/Kazakh)
- Event management and table booking for events
- User registration and authentication
- Use IIKO API and Stripe( or API Kaspi) 

## How to Run
Check the docker-compose.yml file to run the project with all dependencies.