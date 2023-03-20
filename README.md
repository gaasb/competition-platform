# competition-platform
### Create and manage competition events.
### Requirement
```text
go above 1.20
Docker
make
```
### Quick Start

```bash
TBD !!! Not working
make build
make run
```
---
### Available routes for REST API:
Default response:
```json
{
  "data": null,
  "message": [
    "some validation info"
  ],
  "status": "error"
}
```
`status` have only two states: `ok` and `error`.
Possibles responses code: `200`, `400`, `401`, `404`
---
### Tournament group:
**:id** is tournament id
> **GET**: `/api/v1/tournament/` - To get all tournaments. **Query** for pagination: **page**=number, **prev**=page_number, **limit**=number. Min and Max value for limit: 1 and 15 by default.

> **GET**: `/api/v1/tournament/:id` - To get data for a specific tournament. **Query** to get brackets for tournament: **brackets**=true

> **POST**: `/api/v1/tournament/:id` - To create a new tournament
>Payload JSON body to create tournament:
> **start_at** and **end_at** must be **greater** than **current real time**.
> Also this data is usable for the **PUT** method which **updates** the data for the tournament.

Payload:
```json
{
	"discipline_name": "Blazing Fast Racing",
	"title": "Awesome Tournament №1",
	"start_at": "2023-01-01T00:00:00.000Z", 
	"end_at": "2024-01-15T20:30:00.000Z",
	"description": {
		"optionally": "description is json type in db to be possible to save wysiwyg data or something.",
		"img": "url",
		"contact": "contact@me.me"
	}
}
```
> **PUT**: `/api/v1/tournament/:id` - To **update** the data for the tournament. Look at **JSON payload** on **POST** method.

> **DELETE**: `/api/v1/tournament/:id` - To delete tournament.
> 
---
### Brackets group:
**:id** is bracket id
**:tid** is tournament id
... To **create** new bracket and **update** status in **unavailable** when tournament is **ended**.
Bracket have 3 states: `pending`, `live` and `finished`. Default status is `pending`.
> **GET**: `/api/v1/bracket/:tid` - To get all brackets from tournament.

> **POST**: `/api/v1/bracket/:tid` - Add new bracket for tournament. The default **lock** for creating brackets in a tournament is no more than **5**. 

Payload:
```json
{
  "type": "ROUND_ROBIN",
  "max_team": 10,
  "max_participants_in_team": 2,
  "playoff_rounds": 1,
  "final_rounds": 1,
  "grand_final_rounds": 3
}
```
Available **types**: `ROUND_ROBIN`, `SINGLE_ELIMINATION`, `DOUBLE_ELIMINATION`

> **PATCH**: `/api/v1/bracket/:id` - To **start** or **finish** bracket. **Query** for change status: **start**=true - set status to `live`, **end**=true - set status to `finished`. There are **no payload** data. To **start** the number of commands in bracket must be **greater** than **2**.

> **DELETE**: `/api/v1/bracket/:id` - To delete bracket.
---
### Participants group:
**:id** is bracket id
**:team** is team alias
> **GET**: `/api/v1/bracket/participants/:id/` - To get all participants and teams from bracket.

> **POST**: `/api/v1/bracket/participants/:id/` - Add a new team with participants. The array of participants **must be equal** to the `max_participants_in_team` value from the bracket. 

Payload:
```json
{
  "team": "Awesome Team",
  "participants": [
    {
      "user_alias": "Player1",
      "contact": "https://t.me/xgaax"
    },
    {
      "user_alias": "Player2",
      "contact": "https://github.com/gaasb"
    }

  ]
}
```
> **DELETE**: `/api/v1/bracket/participants/:id/:team` - To delete team from bracket. 
