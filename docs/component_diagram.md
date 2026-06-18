
# Components Diagram
This diagram shows the communication between components at an infrastructure level.

```mermaid
graph TD
    subgraph Frontend["Frontend (Vanilla JS)"]
        UI_Admin[Admin Panel]
        UI_Home[Gallery]
    end

    subgraph Backend["Server"]
        API[API REST / GO]
    end

    subgraph BaaS["Firebase"]
        Auth[Firebase Auth]
        DB[(Firestore)]
    end

    subgraph CDN["Storage"]
        Cloud[Cloudinary]
    end

    UI_Admin -- "Login" --> Auth
    UI_Admin -- "HTTP Requests (Token + Data)" --> API
    UI_Home -- "HTTP Requests GET" --> API
    
    API -- "Token Verification" --> Auth
    API -- "Upload Files" --> Cloud
    API -- "Read / Write JSON" --> DB
```

- **Frontend**: Vanilla JavaScript without frameworks or libraries with 2 modules, Admin panel and Gallery.
- **Backend**: API REST with Gin framework in Go as the orchestrator of services.
- **Firebase**:
    - **Auth**: For authentication and authorization.
    - **Firestore**: For data storage (JSON documents).
- **Storage**: Cloudinary for image storage and delivery.
