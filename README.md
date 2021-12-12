# Gacha API

This is the API for my future game! But it's made to be a bit more general, so axkeyz/gacha-api it is.

The actual game for players will use TCP instead of REST API, so this API will be more on providing a user-friendly interface for staff to add & manage assets.

When's the actual game going to come out? Haven't got a clue!

## Planned features

- Staff API endpoints: These will allow staff members to login to the game backend and add items directly into the game by adjusting the database.
    - Create new environments + summonable heroes
    - Create new events + manage rewards
    - Manage loot tables & drop rates 
    - Issue resources to players
    - Add associated art
    - Add associated sound
    - Monitor players & player histories
    - Ban players if necessary
    - Chat with other staff

- Usability features 
    - Staff scope of access is limited by role (role-based access)
    - Oauth2 logins
