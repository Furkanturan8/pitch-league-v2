# Pitch League Project

## Description

This project was created for football matches.<br>
Users can see and participate in matches nearby.<br>
After the matches are played, the match results are entered, and the league ranking information is saved in the database.<br>
If you find enough people, you can also create a team in your neighborhood and enter the lower leagues.

## API Endpoints

### Authentication
- **POST /api/auth/login** - Logs in a user.
- **POST /api/auth/refresh** - Refreshes the access token.
- **POST /api/auth/logout** - Logs out the authenticated user.

### Teams
- **GET /api/teams/** - Lists all teams.
- **GET /api/teams/:id** - Retrieves a specific team by ID.
- **POST /api/teams/:id/join/:userID** - Adds a new user (player) to a team.

### Fields
- **GET /api/fields/** - Lists all football fields.
- **GET /api/fields/:id** - Retrieves a football field by ID.

### Games
- **GET /api/games/** - Lists all games.
- **GET /api/games/:id** - Retrieves a game by ID.

### Game Participants
- **GET /api/gameParts/** - Retrieves the relationship between teams and games (football field, time, teams, etc.).
- **GET /api/gameParts/:id** - Retrieves the game participants' relationships by their ID.

### Leagues
- **GET /api/leagues/** - Lists all leagues (e.g., Super League, PTT League).
- **GET /api/leagues/:id** - Retrieves a specific league by ID.

### League Teams
- **GET /api/leaguesTeam/** - Lists all teams in leagues ranked by points.
- **GET /api/leaguesTeam/:id** - Retrieves a team by its league team ID.
- **GET /api/leaguesTeam/league/:id** - Retrieves all teams in a league by league ID.

### Matches
- **GET /api/matches/** - Lists all matches.
- **GET /api/matches/:id** - Retrieves a match by its match ID.

## Admin Operations

### Users
- **POST /api/admin/users/** - Admin creates a new user.
- **GET /api/admin/users/** - Admin retrieves all users.
- **GET /api/admin/users/:id** - Admin retrieves a specific user by ID.
- **DELETE /api/admin/users/:id** - Admin deletes a user by ID.
- **PUT /api/admin/users/:id** - Admin updates a user by ID.

### Matches
- **POST /api/admin/matches/** - Admin creates a new match.
- **DELETE /api/admin/matches/:id** - Admin cancels (deletes) a match by match ID.

### Leagues
- **POST /api/admin/leagues/** - Admin creates a new league.

### League Teams
- **POST /api/admin/leagueTeam/** - Admin registers a newly created team in a league.
- **DELETE /api/admin/leagueTeam/:id** - Admin deletes a team from a league by its league team ID.

### Game Participants
- **GET /api/gameParts/users/:id** - Admin retrieves all users participating in a game by game ID.
- **POST /api/admin/gamePart/** - Admin creates a new game participant.
- **DELETE /api/admin/gamePart/:id** - Admin deletes a game participant by ID.

### Additional Information
- The API uses JWT for authentication.
- Admin routes require admin privileges.
- CORS is enabled for cross-origin requests.
- Swagger documentation is available at `/swagger/`.
