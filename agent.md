Dev Matcher: Low-Level Design (Monolithic Backend)
.

📁 Go Project Structure (Monolith with Internal Modules)
devmatcher/
├── cmd/
│   └── server/
│       └── main.go            # Main app entry point, init, starts HTTP/WebSocket server
├── internal/
│   ├── api/                   # HTTP handlers, request/response DTOs, routing
│   │   ├── handlers/          # HTTP request handlers
│   │   │   ├── auth_handler.go    # Handles user registration, login, token refresh
│   │   │   ├── user_handler.go    # Handles user profile CRUD, preferences
│   │   │   ├── discovery_handler.go # Handles fetching profiles for swiping
│   │   │   ├── interaction_handler.go # Handles like/pass actions, match creation
│   │   │   └── websocket_handler.go # Manages WebSocket connections for chat/real-time
│   │   ├── middleware/        # HTTP middleware (auth, logging, CORS)
│   │   │   ├── auth.go
│   │   │   └── logger.go
│   │   ├── router.go          # Defines HTTP routes and maps to handlers
│   │   └── dto/               # Data Transfer Objects for API requests/responses
│   │       ├── auth_dto.go
│   │       ├── user_dto.go
│   │       └── interaction_dto.go
│   ├── config/
│   │   └── config.go            # Loads and manages application configuration
│   ├── core/                    # Core business logic services
│   │   ├── auth_service.go      # Authentication logic (token generation, validation)
│   │   ├── user_service.go      # User profile management, preference logic
│   │   ├── data_ingestion_service.go # Fetches/processes data from GitHub, LeetCode etc.
│   │   ├── matching_service.go  # Logic to find and rank compatible users
│   │   ├── interaction_service.go # Processes swipes, creates matches
│   │   └── notification_service.go # Manages sending notifications (push, WebSocket)
│   ├── domain/                  # Core domain models/entities
│   │   ├── user.go              # User struct, validation rules
│   │   ├── interaction.go       # Interaction struct (like/pass)
│   │   ├── match.go             # Match struct
│   │   └── message.go           # Chat message struct
│   ├── store/                   # Data persistence layer (interfaces and implementations)
│   │   ├── user_store.go        # Interface for user data operations
│   │   ├── interaction_store.go # Interface for interaction data operations
│   │   ├── match_store.go       # Interface for match data operations
│   │   ├── message_store.go     # Interface for chat message operations
│   │   └── mongodb/             # MongoDB implementation of store interfaces
│   │       ├── user_mongodb.go
│   │       ├── interaction_mongodb.go
│   │       ├── match_mongodb.go
│   │       ├── message_mongodb.go
│   │       └── db.go              # MongoDB connection setup & helper functions
│   ├── external/                # Clients for third-party services
│   │   ├── github_client.go
│   │   ├── leetcode_client.go
│   │   └── push_notification_client.go # e.g., FCM, APNS
│   ├── eventbus/                # Internal event bus (if used for decoupling services)
│   │   └── local_eventbus.go
│   ├── util/                    # General utility functions
│   │   ├── password_hasher.go
│   │   └── token_generator.go
│   └── ws/                      # WebSocket specific logic (connection hub, message handling)
│       ├── hub.go                 # Manages active WebSocket clients and message broadcasting
│       └── client.go              # Represents a connected WebSocket client
├── go.mod
├── go.sum
└── .env.example                 # Example environment variables


⚙️ Backend Flow Examples
Let's trace the flow for a few key operations. Each step involves specific files and functions.
Components Involved in Flows:
Client: React Native App
Go Application Components:
cmd/server/main.go: Initializes everything.
internal/api/router.go: Routes requests.
internal/api/middleware/: Processes requests before handlers.
internal/api/handlers/: Handles HTTP logic, calls services.
internal/core/: Business logic services.
internal/store/: Database interaction interfaces.
internal/store/mongodb/: MongoDB implementation of stores.
internal/domain/: Data models.
internal/external/: External API clients.
internal/ws/: WebSocket components.
1. User Registration Flow 🙋‍♀️➡️💻
Client: Sends POST /api/v1/auth/register with user details (name, email, password, initial profile data).
cmd/server/main.go / HTTP Server: Receives the request.
internal/api/router.go: Matches the route to auth_handler.Register.
internal/api/handlers/auth_handler.go (Register function):
Decodes request body into dto.UserRegistrationRequest.
Validates the DTO.
Calls core.user_service.RegisterUser(ctx, registrationDetails).
If successful, calls core.auth_service.GenerateTokens(ctx, userID) to get JWTs.
Constructs dto.AuthResponse (tokens, user info) and sends HTTP 201.
internal/core/user_service.go (RegisterUser function):
Checks if email exists via store.user_store.GetByEmail(ctx, email).
If email exists, returns an error.
Hashes password using util.password_hasher.Hash(password).
Creates a domain.User object.
Calls store.user_store.Create(ctx, &user) to save to MongoDB.
If user linked external accounts (e.g., GitHub) during signup:
May trigger core.data_ingestion_service.FetchInitialDataForUser(ctx, userID, linkedAccounts) (can be async via goroutine or event bus).
(Optional) Publishes UserRegisteredEvent via eventbus.
Returns the created user ID or user object.
internal/store/mongodb/user_mongodb.go (GetByEmail, Create functions):
Uses the MongoDB driver to query or insert data into the users collection.
internal/core/data_ingestion_service.go (FetchInitialDataForUser function - if called):
Calls relevant clients in internal/external/ (e.g., github_client.GetProfileData).
Processes fetched data.
Updates the user's profile in the database via store.user_store.Update(ctx, &user).
2. Fetching Profiles for Discovery (Swiping) 🕵️‍♂️🔄
Client: Sends GET /api/v1/discovery/profiles (with JWT in Authorization header).
cmd/server/main.go / HTTP Server: Receives request.
internal/api/router.go: Matches route.
internal/api/middleware/auth.go:
Extracts JWT.
Calls core.auth_service.ValidateToken(ctx, tokenString) to validate.
If valid, extracts userID and adds it to request context. If invalid, returns HTTP 401.
internal/api/handlers/discovery_handler.go (GetProfilesForDiscovery function):
Retrieves userID from context.
Calls core.matching_service.GetRecommendations(ctx, userID, discoveryPreferencesDTO). Discovery preferences might come from query params or user's saved settings.
Maps []domain.User to []dto.UserProfileResponse.
Sends HTTP 200 with the list of profiles.
internal/core/matching_service.go (GetRecommendations function):
Fetches current user's profile and preferences via store.user_store.GetByID(ctx, userID).
Retrieves a pool of candidate users from store.user_store.GetCandidateProfiles(ctx, userID, filters):
Filters exclude already swiped/matched users for the current user.
May apply other pre-filters based on user's basic settings (e.g., active users).
Applies matching algorithm:
Scores candidates based on compatibility (skills, interests, external data).
Considers user's specified preferences.
Returns a ranked list of domain.User objects.
internal/store/mongodb/user_mongodb.go (GetByID, GetCandidateProfiles functions):
Queries MongoDB. GetCandidateProfiles might involve complex queries to exclude users based on interactions collection.
3. User Swipes (Like ❤️ / Pass 🚫)
Client: Sends POST /api/v1/interactions/swipe with {"swiped_user_id": "targetUserID", "action": "like"} (with JWT).
router.go & auth.go middleware: Similar to discovery flow, authenticates user.
internal/api/handlers/interaction_handler.go (RecordSwipe function):
Retrieves swiperUserID from context.
Decodes request body into dto.SwipeRequest.
Calls core.interaction_service.ProcessSwipe(ctx, swiperUserID, targetUserID, action).
Responds with dto.SwipeResponse (e.g., {"match_occurred": true/false, "match_id": "if_true"}) and HTTP 200/201.
internal/core/interaction_service.go (ProcessSwipe function):
Creates domain.Interaction object.
Calls store.interaction_store.Create(ctx, &interaction) to save the swipe.
If action is "like":
Calls store.interaction_store.CheckForMutualLike(ctx, swiperUserID, targetUserID).
If a mutual like exists:
Creates domain.Match object.
Calls store.match_store.Create(ctx, &match) to save the new match.
Triggers notifications: core.notification_service.NotifyNewMatch(ctx, matchDetails). This service might use internal/external/push_notification_client.go and/or send a message via internal/ws/hub.go if users are online.
(Optional) Publishes NewMatchEvent via eventbus.
Returns status indicating a match occurred and the match ID.
Returns status indicating swipe recorded.
internal/store/mongodb/interaction_mongodb.go (Create, CheckForMutualLike functions): Database ops on interactions collection.
internal/store/mongodb/match_mongodb.go (Create function): Database ops on matches collection.
internal/core/notification_service.go (NotifyNewMatch function):
Retrieves device tokens/WebSocket IDs for both users.
Sends push notifications via external.push_notification_client.
Sends real-time messages via ws.hub.BroadcastToUser().
4. Sending a Chat Message 💬 (via WebSocket)
Prerequisite: Users are matched, WebSocket connection established.
Client: User A sends a chat message intended for User B (part of an existing match).
The client's WebSocket connection (established via ws_handler.HandleConnection) sends a JSON message like: {"type": "private_message", "payload": {"match_id": "xyz", "receiver_id": "userB_id", "text": "Hello!"}}.
internal/api/handlers/websocket_handler.go (HandleConnection's read loop for User A):
Receives the JSON message.
Parses the message.
Identifies sender_id (User A) from the authenticated WebSocket connection.
internal/ws/hub.go (or a dedicated chat message handler called by websocket_handler):
Validates if sender_id and receiver_id are part of the given match_id.
Creates a domain.Message object.
Calls store.message_store.Create(ctx, &message) to save the message.
Instructs the hub to route this message to receiver_id (User B).
The hub looks up User B's active WebSocket connection(s).
If User B is online: The hub sends the domain.Message (or a DTO version) to User B's WebSocket client(s).
If User B is offline: The message is stored. core.notification_service.SendChatMessagePushNotification(ctx, messageDetails) could be called to send a push notification.
internal/store/mongodb/message_mongodb.go (Create function): Saves the message to the messages collection.
Client (User B): If online, receives the message via its WebSocket connection and displays it.







React Native App interacts with backend through Nginx (HTTP for REST APIs, WebSocket for real-time).


Nginx routes to the Go Monolith, where the Router & Middleware handle requests.


Requests go to HTTP Handlers (e.g., Auth, User, Discovery, etc.) or WSHandler for WebSocket.


Handlers call Core Services (business logic layer), which interact with:


Stores (MongoDB interfaces)


Redis for caching


External Clients (GitHub, LeetCode, FCM)


Real-time actions (chat, notifications) are sent via WSHub or PushClient.


Internal Event Bus decouples modules (e.g., match creation triggers a notification).




